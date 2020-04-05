.ONESHELL:
.SHELL := /bin/bash

GO?= go
GO_OPTS?=
GO_VERSION_MIN?=1.14
GO_FMT?=gofmt
GO_FMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GO_VERSION?= $(shell $(GO) version)
GO_VERSION_NUMBER?= $(word 3, $(GO_VERSION))
GO_BUILD_PLATFORM?= $(subst /,-,$(lastword $(GO_VERSION)))
GO_PATH:= $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
GO_PRE_111?= $(shell echo $(GO_VERSION_NUMBER) | grep -E 'go1\.(10|[0-9])\.')
BUILD_DEV?=1
BUILD_PLATFORM:=$(subst /,-,$(lastword $(GO_VERSION)))
BUILD_BUILDER:=$(shell whoami)
BUILD_HOST:=$(shell hostname)
BUILD_DATE:=$(shell /bin/date -u)
PRERELEASE?=
PROJECT_DIR?=$(PWD)
BIN_DIR?="$(PROJECT_DIR)/bin"
SERVICE_BINARY?=$(shell basename $(CURDIR))
SERVICE_CMD?="$(PROJECT_DIR)/cmd/$(SERVICE_BINARY)/"
SCRIPTS_DIR?="$(PROJECT_DIR)/scripts"
BUILD_SCRIPTS_DIR?="$(SCRIPTS_DIR)/build"
VERSION?=$(shell $(BUILD_SCRIPTS_DIR)/version.sh)
TEST_TIMEOUT?=60s
FEATURE_TAG?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
GIT_TAG?="$(shell git describe)"
GIT_COMMIT?="$(shell git rev-list -1 HEAD)"
GIT_BRANCH?= $(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
GIT_AUTHOR=$(shell git show -s --format='%ae' $(GIT_COMMIT}))
DOCKER_IMAGE_NAME := unanet-docker.jfrog.io/eve-bot
DOCKER_IMAGE_TAG?=$(shell $(BUILD_SCRIPTS_DIR)/docker-image-tag.sh $(VERSION) $(FEATURE_TAG))
TIMESTAMP_UTC:=$(shell /bin/date -u "+%Y%m%d%H%M%S")
TS:=$(shell /bin/date "+%Y%m%d%H%M%S")
pkgs?=./...

.PHONY: all
all: clean show-version version-check unused fmt fmtcheck vet style build test

.PHONY: clean
clean:
	GO111MODULE=$(GO111MODULE) $(GO) clean $(GO_OPTS) $(SERVICE_CMD)
	rm -rf $(BIN_DIR)

.PHONY: show-version
show-version:
	@echo "$(VERSION)"

.PHONY: show-tag
show-tag:
	@echo "${DOCKER_IMAGE_TAG}"	

.PHONY: show-git-details
show-git-details:
	@echo
	@echo "===> Git Details..."
	@echo "	sha: $(GIT_COMMIT)"
	@echo "	branch: $(GIT_BRANCH)"
	@echo "	tag: $(GIT_TAG)"
	@echo "	author: $(GIT_AUTHOR)"
	@echo

.PHONY: show-build-details
show-build-details:
	@echo
	@echo "===> Build Details..."
	@echo "	golang: $(GO_VERSION_NUMBER)"
	@echo "	platform: $(BUILD_PLATFORM)"
	@echo "	builder: $(BUILD_BUILDER)"
	@echo "	host: $(BUILD_HOST)"
	@echo "	date: $(BUILD_DATE)"
	@echo

## Used for local development. Detects OS/ARCH (good when on Mac or not linux_ amd64)
.PHONY: docker
docker:
	@echo
	@echo "===> Building Docker Image..."
	@docker build \
		--build-arg BUILD_HOST="${BUILD_HOST}" \
		--build-arg BUILDER="${BUILD_BUILDER}" \
		--build-arg GIT_COMMIT="${GIT_COMMIT}" \
		--build-arg GIT_COMMIT_AUTHOR="${GIT_AUTHOR}" \
		--build-arg BUILD_DATE="${TIMESTAMP_UTC}" \
		--build-arg VERSION="${VERSION}" \
		--build-arg GIT_BRANCH="${GIT_BRANCH}" \
		--build-arg PRERELEASE="${PRERELEASE}" \
		. -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	@docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${DOCKER_IMAGE_NAME}:latest
	@echo "Docker Image(s) built"

.PHONY: docker-push
docker-push: docker
	@echo
	@echo "===> Pushing Docker Image..."
	@docker push ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	@docker push ${DOCKER_IMAGE_NAME}:latest

.PHONY: build
build:
	@echo "===> Dev Build..."
	@BUILD_DEV=${BUILD_DEV} \
	GIT_BRANCH=${GIT_BRANCH} \
	VERSION=${VERSION} \
	OUTPUT_DIR=${BIN_DIR} \
	PRERELEASE=${PRERELEASE} \
	sh -c "'$(BUILD_SCRIPTS_DIR)/build.sh'"

.PHONY: git-tag
git-tag:
	@echo "===> Git Tag: ${VERSION}"
	git tag -a ${VERSION} -m "${VERSION}"

.PHONY: sign
sign:
	@echo "===> Signing Release Artifacts"
	@sh -c "'$(BUILD_SCRIPTS_DIR)/gpg-sign.sh'"

.PHONY: verify
verify:
	@echo "===> Verify Release Signatures"
	@sh -c "'$(BUILD_SCRIPTS_DIR)/gpg-verify.sh'"


.PHONY: prep
prep: clean show-version version-check unused fmtcheck vet

default: all

unexport GOVENDOR
ifeq (, $(GO_PRE_111))
	ifneq (,$(wildcard go.mod))
		# Enforce Go modules support just in case the directory is inside GOPATH
		GO111MODULE := on

		ifneq (,$(wildcard vendor))
			# Always use the local vendor/ directory to satisfy the dependencies.
			GO_OPTS := $(GO_OPTS) -mod=vendor
		endif
	endif
else
	ifneq (,$(wildcard go.mod))
		ifneq (,$(wildcard vendor))
$(warning service requires Go >= '$(GO_VERSION_MIN)' because of Go modules)
$(warning Current Go runtime is '$(GO_VERSION_NUMBER)')
		endif
	else
		# This repository isn't using Go modules (yet).
		GOVENDOR := $(GO_PATH)/bin/govendor
	endif

	unexport GO111MODULE
endif


.PHONY: fmtcheck
fmtcheck:
	@echo "===> checking code format..."
	@sh -c "'$(BUILD_SCRIPTS_DIR)/gofmtcheck.sh'"

.PHONY: test
test: fmtcheck
	@echo "===> running all tests"
	GO111MODULE=$(GO111MODULE) $(GO) test -timeout=$(TEST_TIMEOUT) -parallel=4 -race $(GO_OPTS) $(pkgs)

.PHONE: test-cover
test-cover: fmtcheck
	@echo "===> running test coverage"
	GO111MODULE=$(GO111MODULE) $(GO) test -timeout=$(TEST_TIMEOUT) -coverprofile=coverage.out -parallel=4 $(GO_OPTS) $(pkgs)
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: fmt
fmt:
	@echo "===> formatting code..."
	$(GO) fmt $(pkgs)

.PHONY: vet
vet:
	@echo "===> vetting code..."
	GO111MODULE=$(GO111MODULE) $(GO) vet $(GO_OPTS) $(pkgs)

.PHONY: run
run: 
	@echo "===> running server..."
	$(GO) run $(SERVICE_CMD)

.PHONY: version-check
version-check: 
	@sh -c "'$(BUILD_SCRIPTS_DIR)/goversioncheck.sh' '$(GO_VERSION_MIN)'"

.PHONY: style
style:
	@echo "===> server code style check..."
	@fmtRes=$$($(GO_FMT) -d $$(find . -path ./vendor -prune -o -name '*.go' -print)); \
	if [ -n "$${fmtRes}" ]; then \
		echo "gofmt checking failed!"; echo "$${fmtRes}"; echo; \
		echo "Please ensure you are using $$($(GO) version) for formatting code."; \
		exit 1; \
	fi

.PHONY: tidy
tidy:
	@echo "===> tidy go modules"
	$(GO) mod tidy

.PHONY: unused
unused: $(GOVENDOR)
ifdef GOVENDOR
	@echo "===> running check for unused packages..."
	@$(GOVENDOR) list +unused | grep . && exit 1 || echo 'No unused packages'
else
ifdef GO111MODULE
	@echo "===> running check for unused/missing packages in go.mod..."
	GO111MODULE=$(GO111MODULE) $(GO) mod tidy
	@git diff --exit-code -- go.sum go.mod
ifneq (,$(wildcard vendor))
	@echo "===> running check for unused packages in vendor/..."
	GO111MODULE=$(GO111MODULE) $(GO) mod vendor
	@git diff --exit-code -- go.sum go.mod vendor/
endif
endif
endif


# Code Cyclomatic Complexity
cyclomatic-top:
	gocyclo -top 10 $(GO_FMT_FILES)

cyclomatic-over:
	gocyclo -over 10 $(GO_FMT_FILES)