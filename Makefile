.ONESHELL:
.SHELL := /bin/bash

GO_FMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
# GO_VERSION?=$(shell go version)
GO_VERSION?=1.14.1
GO_VERSION_NUMBER?=$(word 3, $(GO_VERSION))
GO_BUILD_PLATFORM?=$(subst /,-,$(lastword $(GO_VERSION)))
BUILD_PLATFORM:=$(subst /,-,$(lastword $(GO_VERSION)))
BUILD_BUILDER:=$(shell whoami)
BUILD_HOST:=$(shell hostname)
BUILD_DATE:=$(shell /bin/date -u)
PRERELEASE?=
PROJECT_DIR?=$(PWD)
BIN_DIR?="$(PROJECT_DIR)/bin"
SERVICE_BINARY?=$(shell basename $(CURDIR))
SERVICE_CMD?="$(PROJECT_DIR)/cmd/$(SERVICE_BINARY)/"
BUILD_SCRIPTS_DIR?="$(PROJECT_DIR)/scripts"
VERSION?=$(shell $(BUILD_SCRIPTS_DIR)/version.sh)
FEATURE_TAG?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
GIT_TAG?="$(shell git describe)"
GIT_COMMIT?="$(shell git rev-list -1 HEAD)"
# GIT_BRANCH?=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
CI_COMMIT_REF_SLUG?=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
GIT_AUTHOR=$(shell git show -s --format='%ae' $(GIT_COMMIT}))
DOCKER_IMAGE_NAME:=unanet-docker.jfrog.io/eve-bot
DOCKER_IMAGE_TAG?=$(shell $(BUILD_SCRIPTS_DIR)/docker-image-tag.sh $(VERSION) $(FEATURE_TAG))
TIMESTAMP_UTC:=$(shell /bin/date -u "+%Y%m%d%H%M%S")
TS:=$(shell /bin/date "+%Y%m%d%H%M%S")

export CI_COMMIT_REF_SLUG

default: build

.PHONY: git-tag
git-tag:
	@echo "===> Git Tag: ${VERSION}"
	git tag -a ${VERSION} -m "${VERSION}"

.PHONY: show-version
show-version:
	@echo "$(VERSION)"

.PHONY: git-details
git-details:
	@echo
	@echo "===> Git Details..."
	@echo "	sha: $(GIT_COMMIT)"
	@echo "	branch: $(CI_COMMIT_REF_SLUG)"
	@echo "	tag: $(GIT_TAG)"
	@echo "	author: $(GIT_AUTHOR)"
	@echo

.PHONY: build-details
build-details:
	@echo
	@echo "===> Build Details..."
	@echo "	golang: $(GO_VERSION_NUMBER)"
	@echo "	platform: $(BUILD_PLATFORM)"
	@echo "	builder: $(BUILD_BUILDER)"
	@echo "	host: $(BUILD_HOST)"
	@echo "	date: $(BUILD_DATE)"
	@echo

## Used for local development. Detects OS/ARCH (good when on Mac or not linux_ amd64)
.PHONY: build
build:
	@echo
	@echo "===> Building Docker Image..."
	@docker build \
		--build-arg BUILD_HOST="${BUILD_HOST}" \
		--build-arg BUILDER="${BUILD_BUILDER}" \
		--build-arg GIT_COMMIT="${GIT_COMMIT}" \
		--build-arg GIT_COMMIT_AUTHOR="${GIT_AUTHOR}" \
		--build-arg BUILD_DATE="${TIMESTAMP_UTC}" \
		--build-arg VERSION="${VERSION}" \
		--build-arg GIT_BRANCH="${CI_COMMIT_REF_SLUG}" \
		--build-arg PRERELEASE="${PRERELEASE}" \
		. -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	@docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${DOCKER_IMAGE_NAME}:latest
	@echo "Docker Image(s) built"

.PHONY: push
push: build
	@echo
	@echo "===> Pushing Docker Image..."
	@docker push ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	@docker push ${DOCKER_IMAGE_NAME}:latest

# Code Cyclomatic Complexity
cyclomatic-top:
	gocyclo -top 10 $(GO_FMT_FILES)

cyclomatic-over:
	gocyclo -over 10 $(GO_FMT_FILES)