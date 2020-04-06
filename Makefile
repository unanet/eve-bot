.ONESHELL:
.SHELL := /bin/bash

TIMESTAMP_UTC:=$(shell /bin/date -u "+%Y%m%d%H%M%S")

# CI Variables: Use CI when on CI Server, otherwise set explicitly when running locally
CI_COMMIT_BRANCH?=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
CI_PROJECT_NAME?=$(shell basename $(CURDIR))
CI_COMMIT_SHA?=$(shell git rev-list -1 HEAD)
CI_COMMIT_SHORT_SHA?=$(shell git rev-parse --short=8 HEAD)
CI_PROJECT_URL?=https://gitlab.unanet.io/devops/eve-bot

# Export the CI variables
export CI_COMMIT_BRANCH
export CI_PROJECT_NAME
export CI_COMMIT_SHA
export CI_COMMIT_SHORT_SHA

PRERELEASE?=
GO_FILES:=$$(find . -name '*.go' | grep -v vendor)
GIT_TAG:=$(shell git describe --abbrev=0)

GIT_AUTHOR:=$(shell git show -s --format='%ae' $(CI_COMMIT_SHA}))
BUILD_BUILDER:=$(shell whoami)
BUILD_HOST:=$(shell hostname)
BUILD_DATE:=$(shell /bin/date -u)
VERSION:=$(shell $(PWD)/scripts/version.sh $(CI_COMMIT_SHORT_SHA) $(CI_COMMIT_BRANCH) $(GIT_TAG))
FEATURE_TAG:=$(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
DOCKER_IMAGE_TAG:=$(shell $(PWD)/scripts/docker-image-tag.sh $(VERSION) $(FEATURE_TAG) $(CI_COMMIT_BRANCH) $(GIT_TAG))
DOCKER_IMAGE_NAME:=unanet-docker.jfrog.io/$(CI_PROJECT_NAME)
DOCKER_BUILD_IMAGE:=unanet-docker.jfrog.io/golang-base

DOCKER_IMAGE_ARGS := \
	--build-arg BUILD_HOST="${BUILD_HOST}" \
	--build-arg BUILDER="${BUILD_BUILDER}" \
	--build-arg GIT_COMMIT="${CI_COMMIT_SHA}" \
	--build-arg GIT_COMMIT_AUTHOR="${GIT_AUTHOR}" \
	--build-arg BUILD_DATE="${TIMESTAMP_UTC}" \
	--build-arg VERSION="${VERSION}" \
	--build-arg GIT_BRANCH="${CI_COMMIT_BRANCH}" \
	--build-arg PRERELEASE="${PRERELEASE}" \
	--build-arg BUILD_IMAGE="${DOCKER_BUILD_IMAGE}"


LABEL_PREFIX := com.unanet
DOCKER_IMAGE_LABELS := \
	--label "${LABEL_PREFIX}.git_commit_sha=${CI_COMMIT_SHORT_SHA}" \
	--label "${LABEL_PREFIX}.gitlab_project_id=${CI_PROJECT_ID}" \
	--label "${LABEL_PREFIX}.build_number=${BUILD_NUMBER}" \
	--label "${LABEL_PREFIX}.build_date=${TIMESTAMP_UTC}" \
	--label "${LABEL_PREFIX}.version=${VERSION}" 


default: details build


[Slack App Creds](https://api.slack.com/apps/A011B3L27P1/general?)

.PHONY: changelog
changelog:
	@echo
	@echo "===> Changelog..."
	@git log v0.4.0...v0.5.0 --pretty=format:'1. [view commit](${CI_PROJECT_URL}/-/commit/%H)	%cn	`%s`	(%ci)' --reverse | tee CHANGELOG.md
	# @git log v0.4.0...v0.5.0 --pretty=format:'<li> <a href="${CI_PROJECT_URL}/-/commit/%H">view commit &bull;</a> %s</li> ' --reverse | tee CHANGELOG.md
	# @git log ${GIT_TAG}...${VERSION} --pretty=format:'<li> <a href="${CI_PROJECT_URL}/-/commit/%H">view commit &bull;</a> %s</li> ' --reverse | tee CHANGELOG.md
	# @curl --header 'Content-Type: application/json' --header "PRIVATE-TOKEN: ${BUILD_ADMIN_KEY}" \
    #  	--data '{ "name": "${VERSION} Release", "tag_name": "${VERSION}", "description": "Release ${VERSION} from ${GIT_AUTHOR}", "assets": { "links": [{ "name": "hoge", "url": "https://google.com", "filepath": "/binaries/linux-amd64" }] } }' \
    #  	--request POST https://gitlab.example.com/api/v4/projects/24/releases

.PHONY: tag
tag:
	@echo
	@echo "===> Git Tag Version: ${VERSION}"
	@git remote remove origin
	@git remote add origin https://${BUILD_ADMIN_USER}:${BUILD_ADMIN_KEY}@${CI_SERVER_HOST}/${CI_PROJECT_PATH}.git
	@git config user.email "${BUILD_ADMIN_EMAIL}"
	@git config user.name "${BUILD_ADMIN_USER}"
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
	@echo "	full sha: $(CI_COMMIT_SHA)"
	@echo "	short sha: $(CI_COMMIT_SHORT_SHA)"	
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
	@docker build ${DOCKER_IMAGE_ARGS} . ${DOCKER_IMAGE_LABELS} -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	@docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${DOCKER_IMAGE_NAME}:latest
	@echo "Docker Image(s) built"
	@echo

.PHONY: publish
publish:
	@echo
	@echo "===> Publish Docker Image..."
	@docker push ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	@docker push ${DOCKER_IMAGE_NAME}:latest
	@echo "===> Docker Image Pushed..."
	@echo
