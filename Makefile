# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
# SPDX-FileCopyrightText: 2024 metal-stack Authors
#
# SPDX-License-Identifier: Apache-2.0

ENSURE_GARDENER_MOD         := $(shell go get github.com/gardener/gardener@$$(go list -m -f "{{.Version}}" github.com/gardener/gardener))
GARDENER_HACK_DIR    		:= $(shell go list -m -f "{{.Dir}}" github.com/gardener/gardener)/hack
IMAGE_TAG                   := $(or ${GITHUB_TAG_NAME}, latest)
REGISTRY                    := ghcr.io/metal-stack
IMAGE_PREFIX                := $(REGISTRY)
REPO_ROOT                   := $(shell dirname "$(realpath $(lastword $(MAKEFILE_LIST)))")
HACK_DIR                    := $(REPO_ROOT)/hack
HOSTNAME                    := $(shell hostname)
LD_FLAGS                    := "-w -X github.com/metal-stack/gardener-extension-backup-s3/pkg/version.Version=$(IMAGE_TAG)"
LEADER_ELECTION             := false
IGNORE_OPERATION_ANNOTATION := false
WEBHOOK_CONFIG_URL          := localhost
GO_VERSION                  := 1.22
GOLANGCI_LINT_VERSION       := v1.56.2

ifeq ($(CI),true)
  DOCKER_TTY_ARG=""
else
  DOCKER_TTY_ARG=t
endif

export GO111MODULE := on

TOOLS_DIR := $(HACK_DIR)/tools
include $(GARDENER_HACK_DIR)/tools.mk


#################################################################
# Rules related to binary build, Docker image build and release #
#################################################################

.PHONY: build
build:
	go build -ldflags $(LD_FLAGS) -tags netgo ./cmd/gardener-extension-backup-s3

.PHONY: install
install: revendor $(HELM)
	@LD_FLAGS="-w -X github.com/gardener/$(EXTENSION_PREFIX)-$(NAME)/pkg/version.Version=$(VERSION)" \
	bash $(GARDENER_HACK_DIR)/install.sh ./...

.PHONY: docker-image
docker-image:
	@docker build --no-cache \
		--tag $(IMAGE_PREFIX)/gardener-extension-backup-s3:$(IMAGE_TAG) \
		--file Dockerfile --memory 6g .

.PHONY: docker-push
docker-push:
	@docker push $(IMAGE_PREFIX)/gardener-extension-backup-s3:$(IMAGE_TAG)

#####################################################################
# Rules for verification, formatting, linting, testing and cleaning #
#####################################################################

.PHONY: tidy
tidy:
	@GO111MODULE=on go mod tidy
	@mkdir -p $(REPO_ROOT)/.ci/hack && cp $(GARDENER_HACK_DIR)/.ci/* $(REPO_ROOT)/.ci/hack/ && chmod +xw $(REPO_ROOT)/.ci/hack/*

.PHONY: clean
clean:
	@$(shell find ./example -type f -name "controller-registration.yaml" -exec rm '{}' \;)
	@bash $(GARDENER_HACK_DIR)/clean.sh ./cmd/... ./pkg/...

.PHONY: check-generate
check-generate:
	@bash $(GARDENER_HACK_DIR)/check-generate.sh $(REPO_ROOT)

.PHONY: check
check: $(GOIMPORTS) $(GOLANGCI_LINT) $(HELM)
	@REPO_ROOT=$(REPO_ROOT) bash $(GARDENER_HACK_DIR)/check.sh --golangci-lint-config=./.golangci.yaml ./cmd/... ./pkg/...
	@REPO_ROOT=$(REPO_ROOT) bash $(GARDENER_HACK_DIR)/check-charts.sh ./charts

.PHONY: generate
generate: $(VGOPATH) $(HELM) $(YQ)
	@REPO_ROOT=$(REPO_ROOT) VGOPATH=$(VGOPATH) GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) bash $(GARDENER_HACK_DIR)/generate-sequential.sh ./charts/... ./cmd/... ./pkg/...

.PHONY: generate-in-container
generate-in-container: tidy $(HELM)
	echo $(shell git describe --abbrev=0 --tags) > VERSION
	docker run --rm -i$(DOCKER_TTY_ARG) \
		--env GOCACHE=/gocache \
		--mount type=tmpfs,destination=/gocache,tmpfs-mode=1777 \
		--user $$(id -u):$$(id -g) \
		--volume $(PWD):/go/src/github.com/metal-stack/gardener-extension-backup-s3:z \
		--workdir /go/src/github.com/metal-stack/gardener-extension-backup-s3 \
		golang:$(GO_VERSION) \
		/usr/bin/make generate

.PHONY: format
format: $(GOIMPORTS) $(GOIMPORTSREVISER)
	@bash $(GARDENER_HACK_DIR)/format.sh ./cmd ./pkg

.PHONY: verify
verify: check format
