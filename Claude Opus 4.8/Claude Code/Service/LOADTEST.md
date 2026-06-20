# Load Test Report

Load generator: `loadtest/main.go` (Go, stdlib only) — logs in, seeds a
customer/plan/subscription, then drives N concurrent workers against a scenario
for a fixed duration and reports throughput, error rate and latency percentiles.

Run it yourself:

```bash
docker compose up --build -d
go run ./loadtest -scenario read    -c 100 -d 15s
go run ./loadtest -scenario usage   -c 100 -d 15s
go run ./loadtest -scenario invoice -c 100 -d 15s
go run ./loadtest -scenario mixed   -c 100 -d 20s
```

## Environment

- Single host, **12 logical CPUs**, 8 GiB available to Docker.
- All processes (API, Postgres, Redis, worker, scheduler **and** the load
  generator) share this one host, so they compete for CPU. Absolute numbers
  would be higher on separated hosts; **relative** before/after comparisons are valid.
- Scenarios:
  - `read` — GET plans / customers / invoices / subscriptions / healthz (70/30 auth-vs-health mix).
  - `usage` — `POST /usage` with a unique idempotency key (single-row insert).
  - `invoice` — `POST /invoices/generate` (heavy: ~6 reads + invoice + line inserts in one tx).
  - `mixed` — 70% read, 25% usage write, 5% invoice generate.
- Auth `POST /login` is deliberately **not** in the hot loop — it is rate-limited
  by design (10/min/IP); the token is obtained once and reused.

## Results — baseline (`DB_MAX_CONNS=10`)

| Scenario | Conc | Throughput | Errors | p50 | p90 | p99 | max |
|----------|-----:|-----------:|:------:|----:|----:|----:|----:|
| read     |  50  | **16,422/s** | 0% | 3.0ms | 4.8ms | 7.8ms | 132ms |
| read     | 100  | 15,081/s | 0% | 6.8ms | 10.2ms | 17.4ms | 203ms |
| read     | 200  | 15,468/s | 0% | 14.1ms | 19.0ms | 28.4ms | 562ms |
| usage    |  50  | 3,570/s | 0% | 12.9ms | 15.3ms | 43.7ms | 348ms |
| usage    | 100  | 3,794/s | 0% | 25.8ms | 28.2ms | 34.5ms | 101ms |
| usage    | 200  | 3,817/s | 0% | 51.1ms | 55.8ms | 69.1ms | 311ms |
| invoice  |  50  | 988/s | 0% | 49.1ms | 56.5ms | 72.3ms | 147ms |
| invoice  | 100  | 1,005/s | 0% | 97.5ms | 106ms | 130ms | 184ms |
| mixed    | 100  | 3,774/s | 0% | 19.7ms | 27.5ms | 200ms (p99) | 312ms |
| mixed    | 200  | 3,624/s | 0% | 40.6ms | 53.1ms | 426ms (p99) | 723ms |

**Observation:** write throughput is flat across concurrency (usage ≈3.8k,
invoice ≈1.0k) while latency grows linearly with concurrency — the textbook
signature of **connection-pool saturation**. The pool was hardcoded at 10.

## Tuning: connection pool 10 → 40

The pool size was made configurable (`DB_MAX_CONNS`, default 10; compose sets
api=40, worker=30). Re-running the write scenarios:

| Scenario | Conc | Pool 10 | Pool 40 | Change |
|----------|-----:|--------:|--------:|:------:|
| usage    | 100  | 3,794/s | **8,705/s** | **+129%** |
| usage    | 200  | 3,817/s | 7,999/s | +109% |
| invoice  |  50  | 988/s   | 2,122/s | +115% |
| invoice  | 100  | 1,005/s | **2,192/s** | **+118%** |

Pure write throughput roughly **doubled**, p99 for invoice dropped 130ms → 80ms.
This confirms the pool was the write bottleneck.

## The tradeoff: bottleneck shifts to Postgres CPU

With the larger pool, the **mixed** workload did *not* improve — it regressed
(c=100: 3,774/s → ~2,200/s; tail p99 ~300ms). Sampling resource use during a
mixed run with pool=40:

```
postgres   CPU ≈ 860%   (8.6 of 12 cores)   <-- saturated
api        CPU ≈ 136%   (1.4 cores)
worker     CPU ≈  23%
```

With only 10 connections, write transactions were throttled, leaving CPU for
reads. With 40, many heavy invoice transactions run at once and **Postgres
becomes CPU-bound**, so reads in the mixed workload now contend with a saturated
database and the tail latency grows. **Pool sizing is a tradeoff, not a free
win** — it should be tuned to the database's core count and the read/write ratio.
On a dedicated multi-core DB host, pool=40 is the better setting; on this shared
12-core box, ~20–25 would balance the mixed workload better. The value is now an
env var precisely so it can be tuned per deployment.

## Background workers under load

During the runs the worker pool stayed ahead of the firehose of emitted jobs:

```
jobs completed : 108,908
jobs queued    : 97        (live backlog, draining)
jobs running   : 2
jobs dead      : 0
```

Retries/backoff/DLQ were never triggered — **0 dead jobs** under sustained load.

## Verdict

- **Reads:** ~15–16k req/s, p99 < 30ms, **0 errors** up to 200 concurrent clients.
- **Writes:** ~8.7k usage inserts/s and ~2.2k full invoice generations/s after pool tuning, **0 errors**.
- **Stability:** zero 4xx/5xx, zero timeouts, zero crashes, zero dead jobs across
  every run; latency degrades gracefully (linearly) as load exceeds capacity.
- **Identified bottleneck:** the DB connection pool (fixed in this pass) and,
  beyond it, Postgres CPU — both expected for a write-heavy billing workload and
  both addressable by sizing `DB_MAX_CONNS` and scaling Postgres.

### Caveats
- Load generator shares the host with the server and database → real-world
  numbers on separated hardware would be higher.
- `go run` recompiles per invocation (minor noise between runs).
- Single Postgres instance, no read replicas; no HTTP keep-alive tuning beyond defaults.
