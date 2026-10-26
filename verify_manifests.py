from pathlib import Path
import yaml

objects=[]
for path in sorted(Path("k8s").glob("*.yaml")):
    objects.extend(x for x in yaml.safe_load_all(path.read_text(encoding="utf-8")) if x)
api=next(x for x in objects if x["kind"]=="Deployment" and x["metadata"]["name"]=="cs411-api")
service=next(x for x in objects if x["kind"]=="Service" and x["metadata"]["name"]=="cs411-api")
container=api["spec"]["template"]["spec"]["containers"][0]
labels=api["spec"]["template"]["metadata"]["labels"]
hpa=next((x for x in objects if x["kind"]=="HorizontalPodAutoscaler"),None)
checks={
    "Service selects API pods":all(labels.get(k)==v for k,v in service["spec"]["selector"].items()),
    "CPU and memory requests":all(k in container.get("resources",{}).get("requests",{}) for k in ("cpu","memory")),
    "readiness probe":container.get("readinessProbe",{}).get("httpGet",{}).get("path")=="/ready",
    "liveness probe":container.get("livenessProbe",{}).get("httpGet",{}).get("path")=="/health",
    "HPA targets API":hpa is not None and hpa["spec"]["scaleTargetRef"].get("name")=="cs411-api",
    "HPA CPU target":hpa is not None and any(m.get("resource",{}).get("name")=="cpu" for m in hpa["spec"].get("metrics",[])),
}
for label,ok in checks.items():print(f"{'PASS' if ok else 'FAIL'} {label}")
if not all(checks.values()):raise SystemExit(1)
