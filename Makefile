################################################################################
# Khan
#
# This Makefile contains various helpful commands as well as actual dependency
# targets.  The top section contains commands, while the bottom contains the
# dependencies for those commands.

# Ensure everything is ready to go
.PHONY: default
default: pre-commit-install
	@echo Ready to go!

.PHONY: clean
clean:
	rm -rf bin

.PHONY: nomad-test-server
nomad-test-server: ./bin/nomad
	nomad agent -dev

.PHONY: pre-commit-install
pre-commit-install: ./.git/hooks/pre-commit

.PHONY: build
build: ./bin/khan pre-commit-install

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
./bin/khan: ./cmd/khan/main.go ./internal/app/*.go
	go build -o ./bin/khan ./cmd/khan/main.go

./.git/hooks/pre-commit: ./bin/pre-commit .pre-commit-config.yaml
	pre-commit install -t pre-commit
