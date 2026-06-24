"""
Latency benchmark for billing-service usage endpoint.

Usage:
  python tests/bench_usage.py --requests 15000 --concurrency 50
"""

import argparse
import asyncio
import statistics
import time
from collections import Counter
from datetime import datetime, timezone

import aiohttp


class UsageBenchmark:
    def __init__(self, base_url: str, concurrency: int, email: str, password: str):
        self.base_url = base_url.rstrip('/')
        self.concurrency = concurrency
        self.email = email
        self.password = password
        self.token = None
        self.customer_id = None
        self.subscription_id = None
        self._semaphore = asyncio.Semaphore(concurrency)
        self.latencies = []
        self.statuses = Counter()
        self.errors = []

    async def _request(self, session, method, path, track=True, **kwargs):
        url = f"{self.base_url}{path}"
        async with self._semaphore:
            start = time.perf_counter()
            try:
                async with session.request(method, url, **kwargs) as resp:
                    body = await resp.text()
                    elapsed = (time.perf_counter() - start) * 1000
                    if track:
                        self.statuses[resp.status] += 1
                        self.latencies.append(elapsed)
                        if not resp.ok:
                            self.errors.append((path, resp.status, body[:200]))
                    return resp.status, body
            except Exception as exc:
                elapsed = (time.perf_counter() - start) * 1000
                if track:
                    self.latencies.append(elapsed)
                    self.statuses['exception'] += 1
                    self.errors.append((path, 0, str(exc)))
                raise

    async def setup(self, session):
        import json
        try:
            await self._request(
                session, 'POST', '/auth/register', track=False,
                json={"email": self.email, "password": self.password, "name": "Bench", "role": "admin"}
            )
        except Exception:
            pass
        _, body = await self._request(
            session, 'POST', '/auth/login', track=False,
            json={"email": self.email, "password": self.password}
        )
        self.token = json.loads(body)['access_token']

        _, body = await self._request(
            session, 'POST', '/plans', track=False,
            headers={"Authorization": f"Bearer {self.token}"},
            json={
                "name": "Bench Plan",
                "currency": "USD",
                "billing_interval": "monthly",
                "base_price_cents": 0,
                "trial_days": 0,
                "is_active": True,
                "prices": [
                    {"model": "usage_based", "unit_price_cents": 1, "included_units": 0}
                ]
            }
        )
        plan_id = json.loads(body)['id']

        _, body = await self._request(
            session, 'POST', '/customers', track=False,
            headers={"Authorization": f"Bearer {self.token}"},
            json={"email": "bench@example.com", "name": "Bench Customer", "country": "US", "currency": "USD"}
        )
        self.customer_id = json.loads(body)['id']

        _, body = await self._request(
            session, 'POST', '/subscriptions', track=False,
            headers={"Authorization": f"Bearer {self.token}"},
            json={"customer_id": self.customer_id, "plan_id": plan_id, "quantity": 1, "seats": 1}
        )
        self.subscription_id = json.loads(body)['id']

    async def run(self, total_requests: int):
        import json
        import uuid as uuid_mod
        timeout = aiohttp.ClientTimeout(total=300)
        async with aiohttp.ClientSession(timeout=timeout) as session:
            await self.setup(session)

            print(f"[bench] starting {total_requests} requests, concurrency={self.concurrency}\n")
            overall_start = time.perf_counter()

            async def worker(idx):
                recorded_at = datetime.now(timezone.utc).isoformat()
                idempotency_key = f"bench-{idx}-{uuid_mod.uuid4()}"
                return await self._request(
                    session, 'POST', '/usage',
                    headers={
                        "Authorization": f"Bearer {self.token}",
                        "Idempotency-Key": idempotency_key
                    },
                    json={
                        "customer_id": self.customer_id,
                        "subscription_id": self.subscription_id,
                        "metric": "api_requests",
                        "quantity": 1,
                        "idempotency_key": idempotency_key,
                        "recorded_at": recorded_at
                    }
                )

            await asyncio.gather(*(worker(i) for i in range(total_requests)))
            elapsed = time.perf_counter() - overall_start

            success = self.statuses[200] + self.statuses[201]
            failed = total_requests - success

            lat = sorted(self.latencies)
            p50 = statistics.median(lat)
            p95_index = int(len(lat) * 0.95)
            p95 = lat[p95_index] if p95_index < len(lat) else lat[-1]
            p99_index = int(len(lat) * 0.99)
            p99 = lat[p99_index] if p99_index < len(lat) else lat[-1]
            avg = statistics.mean(lat)
            mx = max(lat) if lat else 0

            print("=== RESULTS ===")
            print(f"Запросов: {total_requests}")
            print(f"Конкурентность: {self.concurrency}")
            print(f"Ошибок: {failed}")

            status_str = ", ".join(f"{code}={count}" for code, count in sorted(self.statuses.items(), key=lambda x: str(x[0])))
            print(f"Статусы: {status_str}")
            print(f"Производительность: {success / elapsed:.0f} req/s")
            print(f"Средняя задержка: {avg:.3f} ms")
            print(f"p50: {p50:.3f} ms")
            print(f"p95: {p95:.3f} ms")
            print(f"p99: {p99:.3f} ms")
            print(f"Максимум: {mx:.3f} ms")

            if failed and self.errors:
                print(f"\nFirst errors:")
                for path, code, body in self.errors[:5]:
                    print(f"  {path}: {code} {body[:200]}")

            # DB consistency check
            _, body = await self._request(
                session, 'GET',
                f"/usage/summary?subscription_id={self.subscription_id}&metric=api_requests&from=2020-01-01T00:00:00Z&to=2099-01-01T00:00:00Z",
                headers={"Authorization": f"Bearer {self.token}"}, track=False
            )
            total_db = json.loads(body)['total_quantity']
            print(f"\nUsage total in DB: {total_db} (expected {success})")
            if total_db != success:
                print("WARNING: DB mismatch")


async def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--url", default="http://localhost:8080/api/v1")
    parser.add_argument("--requests", type=int, default=15000)
    parser.add_argument("--concurrency", type=int, default=50)
    parser.add_argument("--email", default="bench@example.com")
    parser.add_argument("--password", default="bench12345")
    args = parser.parse_args()
    bench = UsageBenchmark(args.url, args.concurrency, args.email, args.password)
    await bench.run(args.requests)


if __name__ == "__main__":
    asyncio.run(main())
