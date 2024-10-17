```
kind create cluster
docker buildx build . -t grafana/beyla-k8s-cache:dev
kind load docker-image grafana/beyla-k8s-cache:dev
kubectl apply -f deployments/app.yml
```


## Some optimizations to consider

* Let passing slog instance to informers
* Add synchronization complete signal in "welcome message" of subscribe
* Don't send UPDATE messages unless there is a useful change in the object properties
* Add a timestamp to the events, and let Beyla query for events younger than a given timestamp, to optimize reconnections
