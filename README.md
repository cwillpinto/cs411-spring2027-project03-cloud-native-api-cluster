# CS411 Spring 2027 · Project 3: Cloud Native API Cluster

Starter template: https://github.com/cwillpinto/cs411-spring2027-project03-cloud-native-api-cluster

Read the included handout. The Go API and Redis counter protocol are supplied. Complete the rate limiter in `src/rate_limit.go`, harden the `Dockerfile`, add API resource/probe settings in `k8s/api.yaml`, and define the HPA in `k8s/hpa.yaml`. The public verifier checks manifest structure; live scaling and recovery require a Kubernetes cluster with metrics-server.

```sh
go test ./...
python -m pip install -r requirements.txt
python verify_manifests.py
docker build -t cs411-project3-api:local .
kubectl apply -f k8s/
```

The default manifest assumes a local cluster that shares the Docker Desktop image store. For other clusters, load or publish your own image and change `image` and `imagePullPolicy` accordingly. Port-forward the Service with `kubectl port-forward svc/cs411-api 18080:80` and run `python loadtest.py --seconds 120 --workers 8 --output Project3_Load_Test.csv`; the script records requests per second, p95 latency, network/server errors, and HTTP 429 counts. The unchanged starter builds but fails the rate-limit and manifest checks.

## Work in your own GitHub repository

The instructor starter is an upstream source. Keep your completed work in a repository under your own GitHub account, then submit the repository URL and the evidence listed in the assignment handout through Blackboard.

If Blackboard links a template repository, choose **Use this template**, create the new repository under your account, and clone that new repository. Check `git remote -v` before pushing; `origin` must be your repository.

For a direct clone, first create an empty repository under your account. Replace the placeholders below with the starter link supplied in Blackboard and your own empty repository URL:

```sh
git clone https://github.com/cwillpinto/cs411-spring2027-project03-cloud-native-api-cluster
cd <cloned-directory>
git remote rename origin upstream
git remote add origin <your-empty-repository-url>
git remote -v
git push -u origin HEAD
```

Commit and push your completed work to your own `origin`. Do not push to the instructor's upstream repository. Use the exact submission files and verification commands in the assignment handout.
