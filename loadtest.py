"""Controlled client load against a port-forwarded Service."""
import argparse
import concurrent.futures
import csv
from pathlib import Path
import sys
import time
import urllib.error
import urllib.request

parser = argparse.ArgumentParser()
parser.add_argument("--url", default="http://localhost:18080/api/work?iterations=2000000")
parser.add_argument("--requests", type=int, default=100)
parser.add_argument("--seconds", type=float, default=0, help="sustained-load duration; overrides --requests")
parser.add_argument("--workers", type=int, default=8)
parser.add_argument("--interval-ms", type=float, default=100, help="per-worker pause during sustained load")
parser.add_argument("--output", help="CSV file to append to; writes a header when new")
args = parser.parse_args()

def request():
    started = time.perf_counter()
    try:
        with urllib.request.urlopen(args.url, timeout=10) as response:
            status = response.status
    except urllib.error.HTTPError as exc:
        status = exc.code
    except Exception:
        status = 0
    return status, (time.perf_counter() - started) * 1000

def p95(values):
    ordered = sorted(values)
    return ordered[max(0, int(.95 * len(ordered)) - 1)] if ordered else 0

def sustained(deadline):
    rows = []
    while time.perf_counter() < deadline:
        rows.append(request())
        time.sleep(max(0, args.interval_ms) / 1000)
    return rows

began = time.perf_counter()
with concurrent.futures.ThreadPoolExecutor(max_workers=args.workers) as pool:
    if args.seconds > 0:
        deadline = began + args.seconds
        results = [row for batch in pool.map(lambda _: sustained(deadline), range(args.workers)) for row in batch]
    else:
        results = list(pool.map(lambda _: request(), range(args.requests)))
elapsed = time.perf_counter() - began
successes = [ms for status, ms in results if status == 200]
errors = sum(status == 0 or status >= 500 for status, _ in results)
limited = sum(status == 429 for status, _ in results)
header = ["requests", "workers", "requests_per_sec", "p95_success_ms", "p95_all_ms", "successes", "errors", "rate_limited"]
row = [len(results), args.workers, f"{len(results)/elapsed:.2f}", f"{p95(successes):.2f}", f"{p95([ms for _, ms in results]):.2f}", len(successes), errors, limited]
if args.output:
    path = Path(args.output)
    new_file = not path.exists() or path.stat().st_size == 0
    with path.open("a", newline="", encoding="utf-8") as stream:
        writer = csv.writer(stream)
        if new_file:
            writer.writerow(header)
        writer.writerow(row)
else:
    writer = csv.writer(sys.stdout)
    writer.writerow(header)
    writer.writerow(row)
