"""
Load test for billing-service.

Configuration via CLI args:
  --url          base URL (default: http://localhost:8080/api/v1)
  --requests     total number of requests (default: 15000)
  --concurrency  number of concurrent coroutines (default: 50)
  --email        admin email for auth
  --password     admin password

Example:
  python tests/load_test.py --requests 15000 --concurrency 50
"""

import argparse
import asyncio
import time
import uuid
from datetime import datetime, timezone

import aiohttp


class LoadTester:
    def __init__(self, base_url: str, concurrency: int, email: str, password: str):
        self.base_url = base_url.rstrip('/')
        self.concurrency = concurrency
        self.email = email
        self.password = password
        self.token = None
        self.customer_id = None
        self.subscription_id = None
        self._semaphore = asyncio.Semaphore(concurrency)

    async def _request(self, session: aiohttp.ClientSession, method: str, path: str, **kwargs):
        url = f"{self.base_url}{path}"
        async with self._semaphore:
            async with session.request(method, url, **kwargs) as resp:
                body = await resp.json()
                if not resp.ok:
                    raise RuntimeError(f"{method} {path} -> {resp.status}: {body}")
                return body

    async def register(self, session: aiohttp.ClientSession):
        try:
            return await self._request(
                session, 'POST', '/auth/register',
                json={"email": self.email, "password": self.password, "name": "Load Tester", "role": "admin"}
            )
        except RuntimeError as e:
            # email already registered is acceptable
            if "CONFLICT" not in str(e) and "23505" not in str(e):
                raise

    async def login(self, session: aiohttp.ClientSession):
        resp = await self._request(
            session, 'POST', '/auth/login',
            json={"email": self.email, "password": self.password}
        )
        self.token = resp['access_token']

    async def setup_plan(self, session: aiohttp.ClientSession):
        plan = await self._request(
            session, 'POST', '/plans',
            headers={"Authorization": f"Bearer {self.token}"},
            json={
                "name": "Load Test Plan",
                "currency": "USD",
                "billing_interval": "monthly",
                "base_price_cents": 1000,
                "trial_days": 0,
                "is_active": True,
                "prices": [
                    {"model": "usage_based", "unit_price_cents": 1, "included_units": 0}
                ]
            }
        )
        return plan['id']

    async def setup_customer(self, session: aiohttp.ClientSession):
        customer = await self._request(
            session, 'POST', '/customers',
            headers={"Authorization": f"Bearer {self.token}"},
            json={"email": "load@example.com", "name": "Load Customer", "country": "US", "currency": "USD"}
        )
        return customer['id']

    async def setup_subscription(self, session: aiohttp.ClientSession, customer_id: str, plan_id: str):
        sub = await self._request(
            session, 'POST', '/subscriptions',
            headers={"Authorization": f"Bearer {self.token}"},
            json={"customer_id": customer_id, "plan_id": plan_id, "quantity": 1, "seats": 1}
        )
        return sub['id']

    async def send_usage(self, session: aiohttp.ClientSession, idx: int):
        recorded_at = datetime.now(timezone.utc).isoformat()
        idempotency_key = f"load-{uuid.uuid4()}-{idx}"
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

    async def run(self, total_requests: int):
        timeout = aiohttp.ClientTimeout(total=300)
        async with aiohttp.ClientSession(timeout=timeout) as session:
            print("[setup] register / login / plan / customer / subscription")
            await self.register(session)
            await self.login(session)
            plan_id = await self.setup_plan(session)
            self.customer_id = await self.setup_customer(session)
            self.subscription_id = await self.setup_subscription(session, self.customer_id, plan_id)

            print(f"[load] starting {total_requests} requests with concurrency={self.concurrency}")
            success = 0
            failed = 0
            start = time.perf_counter()

            async def worker(idx: int):
                nonlocal success, failed
                try:
                    await self.send_usage(session, idx)
                    success += 1
                except Exception as e:
                    failed += 1
                    if failed <= 5 or idx % 1000 == 0:
                        print(f"[error] request {idx}: {e}")

            await asyncio.gather(*(worker(i) for i in range(total_requests)))
            elapsed = time.perf_counter() - start

            print("\n=== RESULTS ===")
            print(f"Total requests : {total_requests}")
            print(f"Success        : {success}")
            print(f"Failed         : {failed}")
            print(f"Elapsed        : {elapsed:.2f} s")
            print(f"RPS            : {success / elapsed:.2f}")

            # verify DB state
            summary = await self._request(
                session, 'GET',
                f"/usage/summary?subscription_id={self.subscription_id}&metric=api_requests&from=2020-01-01T00:00:00Z&to=2099-01-01T00:00:00Z",
                headers={"Authorization": f"Bearer {self.token}"}
            )
            print(f"Usage total in DB: {summary['total_quantity']} (expected {success})")
            if summary['total_quantity'] != success:
                print("WARNING: mismatch between success count and persisted usage summary!")


async def main():
    parser = argparse.ArgumentParser(description="Billing service load test")
    parser.add_argument("--url", default="http://localhost:8080/api/v1")
    parser.add_argument("--requests", type=int, default=15000)
    parser.add_argument("--concurrency", type=int, default=50)
    parser.add_argument("--email", default="loadtest@example.com")
    parser.add_argument("--password", default="loadtest123")
    args = parser.parse_args()

    tester = LoadTester(args.url, args.concurrency, args.email, args.password)
    await tester.run(args.requests)


if __name__ == "__main__":
    asyncio.run(main())
