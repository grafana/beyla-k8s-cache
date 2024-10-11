PROTOC_IMAGE = docker.io/mariomac/protoc-go:latest

.PHONY: protoc-image-build
protoc-image-build:
	docker build . -f protoc.Dockerfile -t $(PROTOC_IMAGE)

.PHONY: protoc-gen
protoc-gen:
	docker run --rm -v $(PWD):/work -w /work $(PROTOC_IMAGE) protoc --go_out=pkg --go-grpc_out=pkg proto/informer.proto
