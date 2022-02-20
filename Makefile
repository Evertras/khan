################################################################################
# Khan
#
# This Makefile contains various helpful commands as well as actual dependency
# targets.  The top section contains commands, while the bottom contains the
# dependencies for those commands.

# Build and run Khan
.PHONY: default
default: ./bin/khan
	@./bin/khan

# Clean temporary files
.PHONY: clean
clean: clean-nomad clean-bin

.PHONY: clean-nomad
clean-nomad:
	rm -rf nomad/data

.PHONY: clean-bin
clean-bin:
	rm -rf bin

# Run a Nomad server for testing purposes
.PHONY: nomad-test-server
nomad-test-server: ./bin/nomad ./nomad/dev-server.hcl
	nomad agent -config ./nomad/dev-server.hcl

# Build everything
.PHONY: build
build: ./bin/khan pre-commit-install

.PHONY: test
test: pre-commit-install
	@go test ./internal/...

# Format our Go code
.PHONY: fmt
fmt:
	@go fmt ./cmd/...
	@go fmt ./internal/...

################################################################################
# Local bin files
#
# This section contains local tools that we download on demand, so that the
# developer doesn't need to download global versions.

NOMAD_VERSION := 1.2.6

# For now we only support Linux 64 bit and MacOS
ifeq ($(shell uname), Darwin)
OS_URL := darwin
else
OS_URL := linux
endif

./bin/nomad:
	@mkdir -p bin
	curl -o bin/nomad.zip \
		https://releases.hashicorp.com/nomad/$(NOMAD_VERSION)/nomad_$(NOMAD_VERSION)_$(OS_URL)_amd64.zip
	@cd bin && unzip nomad.zip
	@rm bin/nomad.zip

./bin/pre-commit.pyz:
	@mkdir -p bin
	curl -Lo bin/pre-commit.pyz https://github.com/pre-commit/pre-commit/releases/download/v2.17.0/pre-commit-2.17.0.pyz

./bin/pre-commit: ./bin/pre-commit.pyz
	@echo '#!/bin/bash\npython3 bin/pre-commit.pyz "$$@"' > ./bin/pre-commit
	@chmod +x ./bin/pre-commit

################################################################################
# Local dependencies and builds
INTERNAL_GO_SOURCES := $(shell find internal/ -name '*.go')

./bin/khan: ./cmd/khan/*.go $(INTERNAL_GO_SOURCES) go.*
	go build -o ./bin/khan ./cmd/khan/*.go

./.git/hooks/pre-commit: ./bin/pre-commit .pre-commit-config.yaml
	./bin/pre-commit install -t pre-commit

./.git/hooks/pre-push: ./bin/pre-commit .pre-commit-config.yaml
	./bin/pre-commit install -t pre-push

# Doesn't depend on template to avoid overwriting unexpectedly
./nomad/dev-server.hcl:
	@sed 's|DATADIR|$(shell pwd)/nomad/data|g' ./nomad/dev-server.tpl.hcl > ./nomad/dev-server.hcl

.PHONY: pre-commit-install
pre-commit-install: ./.git/hooks/pre-commit ./.git/hooks/pre-push

