```
kind create cluster
docker buildx build . -t grafana/beyla-k8s-cache:dev
kind load docker-image grafana/beyla-k8s-cache:dev
kubectl apply -f deployments/app.yml
```