.ONESHELL:
.SHELL := /bin/bash

TIMESTAMP_UTC:=$(shell /bin/date -u "+%Y%m%d%H%M%S")

# CI Variables: Use CI when on CI Server, otherwise set explicitly when running locally
CI_COMMIT_BRANCH?=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
CI_PROJECT_NAME?=$(shell basename $(CURDIR))

PRERELEASE?=
GO_FILES:=$$(find . -name '*.go' | grep -v vendor)
GIT_TAG:=$(shell git describe --abbrev=0)
GIT_COMMIT:=$(shell git rev-list -1 HEAD)
GIT_SHORT_SHA:=$(shell git rev-parse --short=10 HEAD)
GIT_AUTHOR:=$(shell git show -s --format='%ae' $(GIT_COMMIT}))
BUILD_BUILDER:=$(shell whoami)
BUILD_HOST:=$(shell hostname)
BUILD_DATE:=$(shell /bin/date -u)
VERSION:=$(shell $(PWD)/scripts/version.sh $(GIT_SHORT_SHA) $(CI_COMMIT_BRANCH) $(GIT_TAG))
FEATURE_TAG:=$(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
DOCKER_IMAGE_TAG:=$(shell $(PWD)/scripts/docker-image-tag.sh $(VERSION) $(FEATURE_TAG) $(CI_COMMIT_BRANCH) $(GIT_TAG))
DOCKER_IMAGE_NAME:=unanet-docker.jfrog.io/$(CI_PROJECT_NAME)
DOCKER_BUILD_IMAGE:=golang:1.14.1-alpine

# Export and CI variables
export CI_COMMIT_BRANCH
export CI_PROJECT_NAME

default: details build

.PHONY: tag
tag:
	@echo
	@echo "===> Git Tag Version: ${VERSION}"
	@git tag -a ${VERSION} -m "${VERSION}"
	@git push origin ${VERSION}
	@echo

.PHONY: details
details:
	@echo
	@echo "===> Build Details..."
	@echo "	builder: $(BUILD_BUILDER)"
	@echo "	host: $(BUILD_HOST)"
	@echo "	date: $(BUILD_DATE)"
	@echo
	@echo "===> Git Details..."
	@echo "	full sha: $(GIT_COMMIT)"
	@echo "	short sha: $(GIT_SHORT_SHA)"	
	@echo "	branch: $(CI_COMMIT_BRANCH)"
	@echo "	tag: $(GIT_TAG)"
	@echo "	author: $(GIT_AUTHOR)"
	@echo
	@echo "===> Version Details..."
	@echo "	build: $(VERSION)"
	@echo "	feature: $(FEATURE_TAG)"
	@echo
	@echo "===> Docker Details..."
	@echo "	builder: $(DOCKER_BUILD_IMAGE)"	
	@echo "	tag: $(DOCKER_IMAGE_TAG)"	
	@echo "	image: $(DOCKER_IMAGE_NAME)"	
	@echo

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
		--build-arg GIT_BRANCH="${CI_COMMIT_BRANCH}" \
		--build-arg PRERELEASE="${PRERELEASE}" \
		--build-arg BUILD_IMAGE="${DOCKER_BUILD_IMAGE}" \
		. -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	@docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${DOCKER_IMAGE_NAME}:latest
	@echo "Docker Image(s) built"
	@echo

.PHONY: publish
publish: build
	@echo
	@echo "===> Publish Docker Image..."
	@docker push ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	@docker push ${DOCKER_IMAGE_NAME}:latest
	@echo "===> Docker Image Pushed..."
	@echo
