# Main binary configuration
CMD ?= k8s-cache
MAIN_GO_FILE ?= cmd/$(CMD)/main.go
GOOS ?= linux
GOARCH ?= amd64

# todo: upload to a grafana artifact
PROTOC_IMAGE = docker.io/mariomac/protoc-go:latest

# regular expressions for excluded file patterns
EXCLUDE_COVERAGE_FILES="(/grafana/beyla-k8s-cache/pkg/informer/)"

# DRONE_TAG is set from Drone. Required for building container images.
RELEASE_VERSION := $(shell git describe --tags --always)
RELEASE_REVISION := $(shell git rev-parse --short HEAD )
#BUILDINFO_PKG ?= github.com/grafana/beyla/pkg/buildinfo

# go-install-tool will 'go install' any package $2 and install it locally to $1.
# This will prevent that they are installed in the $USER/go/bin folder and different
# projects ca have different versions of the tools
PROJECT_DIR := $(shell dirname $(abspath $(firstword $(MAKEFILE_LIST))))

TOOLS_DIR ?= $(PROJECT_DIR)/bin
TEST_OUTPUT ?= ./testoutput

# $(1) command name
# $(2) repo URL
# $(3) version
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Removing any outdated version of $(1)";\
rm -f $(1)*;\
echo "Downloading $(2)@$(3)" ;\
GOBIN=$(TOOLS_DIR) GOFLAGS="-mod=mod" go install "$(2)@$(3)" ;\
touch "$(1)-$(3)";\
rm -rf $$TMP_DIR ;\
}
endef

# gomod-version returns the version number of the go.mod dependency
define gomod-version
$(shell sh -c "echo $$(grep $(1) go.mod | awk '{print $$2}')")
endef

# Check that given variables are set and all have non-empty values,
# die with an error otherwise.
#
# Params:
#   1. Variable name(s) to test.
#   2. (optional) Error message to print.
check_defined = \
	$(strip $(foreach 1,$1, \
		$(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
	$(if $(value $1),, \
	  $(error Undefined $1$(if $2, ($2))))

# prereqs binary dependencies
GOLANGCI_LINT = $(TOOLS_DIR)/golangci-lint
GOIMPORTS_REVISER = $(TOOLS_DIR)/goimports-reviser
GO_LICENSES = $(TOOLS_DIR)/go-licenses
ENVTEST = $(TOOLS_DIR)/setup-envtest
ENVTEST_K8S_VERSION = 1.30.0

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

GOIMPORTS_REVISER_ARGS = -company-prefixes github.com/grafana -project-name github.com/grafana/beyla-k8s-cache/

define check_format
	$(shell $(foreach FILE, $(shell find . -name "*.go" -not -path "**/vendor/*"), \
		$(GOIMPORTS_REVISER) $(GOIMPORTS_REVISER_ARGS) -list-diff -output stdout $(FILE);))
endef


.PHONY: prereqs
prereqs:
	@echo "### Check if prerequisites are met, and installing missing dependencies"
	mkdir -p $(TEST_OUTPUT)/run
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,v1.60.3)
	$(call go-install-tool,$(GOIMPORTS_REVISER),github.com/incu6us/goimports-reviser/v3,v3.6.4)
	$(call go-install-tool,$(GO_LICENSES),github.com/google/go-licenses,v1.6.0)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,latest)

.PHONY: fmt
fmt: prereqs
	@echo "### Formatting code and fixing imports"
	@$(foreach FILE, $(shell find . -name "*.go" -not -path "**/vendor/*"), \
		$(GOIMPORTS_REVISER) $(GOIMPORTS_REVISER_ARGS) $(FILE);)

.PHONY: checkfmt
checkfmt:
	@echo '### check correct formatting and imports'
	@if [ "$(strip $(check_format))" != "" ]; then \
		echo "$(check_format)"; \
		echo "Above files are not properly formatted. Run 'make fmt' to fix them"; \
		exit 1; \
	fi

.PHONY: lint
lint: prereqs checkfmt
	@echo "### Linting code"
	$(GOLANGCI_LINT) run ./... --timeout=6m

# As generated artifacts are part of the code repo (pkg/ebpf packages), you don't have
# to run this target for each build. Only when you change the C code inside the bpf folder.
# You might want to use the docker-generate target instead of this.

.PHONY: verify
verify: prereqs lint test

# TODO: embed software version in executable, as well as other ldflags
.PHONY: compile
compile:
	@echo "### Compiling project"
	CGO_ENABLED=0 go build -mod vendor -a -o bin/$(CMD) $(MAIN_GO_FILE)

.PHONY: dev
dev: prereqs generate compile-for-coverage

# Generated binary can provide coverage stats according to https://go.dev/blog/integration-test-coverage
.PHONY: compile-for-coverage
compile-for-coverage:
	@echo "### Compiling project to generate coverage profiles"
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -mod vendor -cover -a -o bin/$(CMD) $(MAIN_GO_FILE)

.PHONY: test
test: integration-test
	@echo "### Testing code"
	go test -race -mod vendor -a ./pkg/... -coverpkg=./... -coverprofile $(TEST_OUTPUT)/cover.all.txt

# TODO: merge coverage with test
.PHONY: envtest
integration-test: prereqs
	@echo "### Integration testing"
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./envtest/... -coverpkg=./... -coverprofile cover.envtest.out

.PHONY: cov-exclude-generated
cov-exclude-generated:
	grep -vE $(EXCLUDE_COVERAGE_FILES) $(TEST_OUTPUT)/cover.all.txt > $(TEST_OUTPUT)/cover.txt

.PHONY: coverage-report
coverage-report: cov-exclude-generated
	@echo "### Generating coverage report"
	go tool cover --func=$(TEST_OUTPUT)/cover.txt

.PHONY: coverage-report-html
coverage-report-html: cov-exclude-generated
	@echo "### Generating HTML coverage report"
	go tool cover --html=$(TEST_OUTPUT)/cover.txt

# TODO: add a test to check third_party_licenses is updated
.PHONY: update-licenses
update-licenses: prereqs
	@echo "### Updating third_party_licenses.csv"
	$(GO_LICENSES) report --include_tests ./... > third_party_licenses.csv

.PHONY: clean-testoutput
clean-testoutput:
	@echo "### Cleaning ${TEST_OUTPUT} folder"
	rm -rf ${TEST_OUTPUT}/*


.PHONY: protoc-gen
protoc-gen:
	docker run --rm -v $(PWD):/work -w /work $(PROTOC_IMAGE) protoc --go_out=pkg --go-grpc_out=pkg proto/informer.proto
