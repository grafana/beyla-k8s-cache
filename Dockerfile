# Build the autoinstrumenter binary
FROM golang:1.23 AS builder

ARG TARGETARCH
ENV GOARCH=$TARGETARCH

WORKDIR /opt/app-root

# Copy the go manifests and source
COPY .git/ .git/
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY vendor/ vendor/
COPY go.mod go.mod
COPY go.sum go.sum
#TODO COPY Makefile Makefile
#TODO COPY LICENSE LICENSE
#TODO COPY NOTICE NOTICE
#TODO COPY third_party_licenses.csv third_party_licenses.csv

# Build
RUN go build -o bin/k8s-cache ./cmd/k8s-cache/main.go

# Create final image from minimal + built binary
FROM debian:bookworm-slim

LABEL maintainer="Grafana Labs <hello@grafana.com>"

WORKDIR /

COPY --from=builder /opt/app-root/bin/k8s-cache .
#TODO COPY --from=builder /opt/app-root/LICENSE .
#TODO COPY --from=builder /opt/app-root/NOTICE .
#TODO COPY --from=builder /opt/app-root/third_party_licenses.csv .

COPY --from=builder /etc/ssl/certs /etc/ssl/certs

USER 0:0

ENTRYPOINT [ "/k8s-cache" ]