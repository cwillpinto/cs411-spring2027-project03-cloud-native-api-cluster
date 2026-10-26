# Project 3: Cloud Native API Cluster

**CS 411: Computing Paradigms - Spring 2027**  
**Week 16: April 26-30, 2027**  
**Due:** April 30, 2027  
**Value:** 100 points

## Purpose

Deploy the supplied Go API and Redis dependency as a Kubernetes service. Make the API image safe to run, add health and resource controls, implement request protection, and gather evidence of scaling and recovery under load. Explain what the evidence actually shows, including experimental limits.

## Learning objectives

- Connect image design, probes, resource requests, Services, and autoscaling to observed behavior.
- Explain liveness versus readiness during dependency failure and pod recovery.
- Measure throughput, latency, rejection, and replica behavior during a controlled failure campaign, then connect the deployed system to evidence from earlier course projects.

## Supplied starter and local preflight

The starter includes the API, Redis counter, Go tests, Kubernetes manifests, a manifest verifier, and a load generator. Complete the marked TODOs in `src/rate_limit.go`, `Dockerfile`, `k8s/api.yaml`, and `k8s/hpa.yaml`. Preserve `/health` (process liveness), `/ready` (Redis connectivity), and `/api/work?iterations=...` (CPU work and Redis increment). A rejected request returns HTTP 429.

```sh
go test ./...
python -m pip install -r requirements.txt
python verify_manifests.py
docker build -t cs411-project3-api:local .
```

The unchanged starter intentionally fails limiter tests and manifest verification. Save that initial output. Your completed package must pass before cluster testing. Do not include credentials or private cluster configuration in your repository.

## Milestone A: deployable service (25 points)

Build a multi-stage image that runs the API as a non-root user. Deploy the supplied Redis and API Services and Deployment. In `k8s/api.yaml`, add realistic CPU/memory requests and limits, a liveness probe on `/health`, and a readiness probe on `/ready`. Explain why the probes differ. The local image name assumes Docker Desktop Kubernetes shares its image store. For another cluster, load or publish your own image and adjust `image` and `imagePullPolicy`.

## Milestone B: scaling and protection (40 points)

Implement the per-process rate limiter in `src/rate_limit.go` and document its limit/window behavior. Complete `k8s/hpa.yaml` with CPU-based autoscaling, two to five replicas, and a stated utilization target. CPU-based HPA needs resource requests and metrics-server. Explain why a per-replica limiter permits more aggregate requests when replicas grow and cannot enforce a global quota. In your report, distinguish requests rejected by the limiter from server failures or connection errors; each suggests a different corrective action.

## Milestone C: load and failure campaign (35 points)

Apply the manifests and check rollout, endpoints, and metrics before load:

```sh
kubectl apply -f k8s/
kubectl get deploy,svc,pods,hpa
kubectl top pods
kubectl port-forward svc/cs411-api 18080:80
```

In another terminal, run a baseline and sustained campaign. Tune `--workers`, `--interval-ms`, and `iterations` to create measurable CPU load without overwhelming your machine. Save CSV output as `Project3_Load_Test.csv`; include exact commands and environment in your report.

```sh
python loadtest.py --requests 100 --output Project3_Load_Test.csv
python loadtest.py --seconds 120 --workers 8 --output Project3_Load_Test.csv
```

The CSV records requests/second, p95 latency for successful and all requests, successes, server/network errors, and HTTP 429 responses. During sustained load, record timestamps and Ready versus desired replicas from `kubectl get hpa,pods -w`. Delete one named serving API pod while traffic continues; record replacement, endpoint readiness, and any request disruption in `Project3_Recovery_Transcript.txt`. A final CSV row alone does not prove HPA scaling: include the replica timeline. If metrics-server is unavailable, report that failed check and limitation rather than claiming autoscaling.

## Analysis and submission

Submit your repository URL plus `Project3_Architecture.pdf` or `.md`, `Project3_Load_Test.csv`, `Project3_Recovery_Transcript.txt`, and `Project3_Report.pdf` or `.md`. The architecture artifact should show client, Service, API replicas, Redis, probes, and HPA metrics. The report should explain image and manifest decisions, test setup, baseline versus load, scaling delay, failure recovery, bottleneck, rate-limit behavior, and one concrete improvement. Add a concise systems-showcase section that cites one concrete artifact or measurement from Project 1 and one from Project 2, then explains how parallel, distributed, and cloud-native patterns work together in your deployed, load-tested system. Include `Dockerfile`, `k8s/`, and `src/` in the repository. Screenshots alone are insufficient.

## Rubric

| Criterion | Points |
| --- | ---: |
| Container and deployment correctness | 20 |
| Health, resources, and Service design | 20 |
| Autoscaling and rate limiting | 25 |
| Load and failure evidence | 20 |
| Systems integration and cross-project evidence | 15 |
| **Total** | **100** |

## AI usage and academic integrity

Add an **AI Usage** section to `Project3_Report`. Name each tool and list the prompts or interactions that affected the work. If you did not use AI, write `AI Usage: None.` You must be able to explain and defend every submitted design and line of code.

## Late work

Late work may be accepted for up to three calendar days with a 25 percent deduction per day, following the course syllabus.
