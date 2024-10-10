PROTOC_IMAGE = docker.io/mariomac/protoc-go:latest

.PHONY: protoc-image-build
protoc-image-build:
	docker build . -f protoc.Dockerfile -t $(PROTOC_IMAGE)

.PHONY: protoc-gen
protoc-gen:
	docker run --rm -v $(PWD):/work -w /work $(PROTOC_IMAGE) protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto