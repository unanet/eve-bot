.ONESHELL:
.SHELL := /bin/bash

TIMESTAMP_UTC:=$(shell /bin/date -u "+%Y%m%d%H%M%S")

# CI Variables: Use CI when on CI Server, otherwise set explicitly when running locally
CI_COMMIT_BRANCH?=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
CI_COMMIT_SHORT_SHA?=$(shell git rev-parse --short=8 HEAD)
CI_COMMIT_SHA?=$(shell git rev-list -1 HEAD)
CI_PROJECT_NAME?=$(shell basename $(CURDIR))

# CI_PROJECT_URL?=https://gitlab.unanet.io/devops/eve-bot
CI_PROJECT_URL?=""
CI_API_V4_URL?=""
CI_PROJECT_ID?=0

	# uploadedURL=$$($(docker-exec) curl --request POST --header "PRIVATE-TOKEN: ${BUILD_ADMIN_KEY}" --form "file=@CHANGELOG.md" ${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/uploads | jq .url | sed -e 's/^"//' -e 's/"$$//')") && \
	# curl --request POST --header "PRIVATE-TOKEN: ${BUILD_ADMIN_KEY}" \
	# 	--form "description=Changelog File: [CHANGELOG.md]($$uploadedURL)" ${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/repository/tags/${GIT_TAG}/release


# Export the CI variables
export CI_COMMIT_BRANCH
export CI_PROJECT_NAME
export CI_COMMIT_SHA
export CI_COMMIT_SHORT_SHA

PRERELEASE?=
UPLOAD_URL?=
GO_FILES:=$$(find . -name '*.go' | grep -v vendor)
GIT_TAG:=$(shell git describe --abbrev=0)
CUR_DIR := $(shell pwd)
GOOS?=linux
GOARCH?=amd64
GIT_AUTHOR:=$(shell git show -s --format='%ae' $(CI_COMMIT_SHA}))
BUILD_BUILDER:=$(shell whoami)
BUILD_HOST:=$(shell hostname)
BUILD_DATE:=$(shell /bin/date -u)
VERSION:=$(shell $(CUR_DIR)/scripts/version.sh $(CI_COMMIT_SHORT_SHA) $(CI_COMMIT_BRANCH) $(GIT_TAG))
FEATURE_TAG:=$(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
DOCKER_IMAGE_TAG:=$(shell $(CUR_DIR)/scripts/docker-image-tag.sh $(VERSION) $(FEATURE_TAG) $(CI_COMMIT_BRANCH) $(GIT_TAG))
DOCKER_IMAGE_BASE:=unanet-docker.jfrog.io
DOCKER_IMAGE_NAME:=$(DOCKER_IMAGE_BASE)/$(CI_PROJECT_NAME)
DOCKER_BUILD_IMAGE:=$(DOCKER_IMAGE_BASE)/golang-base
DOCKER_UID = $(shell id -u)
DOCKER_GID = $(shell id -g)
BUILD_ADMIN_USER?=$(shell whoami)
BUILD_ADMIN_EMAIL?=$(GIT_AUTHOR)

DOCKER_IMAGE_ARGS := \
	--build-arg BUILD_IMAGE="${DOCKER_BUILD_IMAGE}"

GOLANG_BUILD_LDFLAGS := \
	-X main.BuildHost="${BUILD_HOST}" \
    -X main.GitBranch="${CI_COMMIT_BRANCH}" \
    -X main.Builder="${BUILD_BUILDER}" \
    -X main.Version="${VERSION}" \
    -X main.BuildDate="${TIMESTAMP_UTC}" \
    -X main.GitCommit="${CI_COMMIT_SHA}" \
    -X main.GitCommitAuthor="${GIT_AUTHOR}" \
    -X main.VersionPrerelease="${PRERELEASE}"

LABEL_PREFIX := com.unanet
DOCKER_IMAGE_LABELS := \
	--label "${LABEL_PREFIX}.git_commit_sha=${CI_COMMIT_SHORT_SHA}" \
	--label "${LABEL_PREFIX}.gitlab_project_id=${CI_PROJECT_ID}" \
	--label "${LABEL_PREFIX}.build_date=${BUILD_DATE}" \
	--label "${LABEL_PREFIX}.version=${VERSION}" \
	--label "${LABEL_PREFIX}.maintainer=${BUILD_ADMIN_USER} <${BUILD_ADMIN_EMAIL}>" \

docker-exec = docker run --rm \
	-e DOCKER_UID=${DOCKER_UID} \
	-e DOCKER_GID=${DOCKER_GID} \
	-e GOOS=${GOOS} \
	-e GOARCH=${GOARCH} \
	-v ${CUR_DIR}:/src \
	-w /src \
	${DOCKER_BUILD_IMAGE}

default: version test build

.PHONY: build
build:
	rm -rf build
	$(docker-exec) go build -ldflags "${GOLANG_BUILD_LDFLAGS}" -o ./build/eve-bot ./cmd/eve-bot/
	tar -czvf build/eve-bot.tar.gz build/eve-bot
	docker build $(DOCKER_IMAGE_ARGS) . ${DOCKER_IMAGE_LABELS} -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	docker tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ${DOCKER_IMAGE_NAME}:latest

.PHONY: test
test:	
	$(docker-exec) go test -v ./...

.PHONY: version
version:	
	@echo "${VERSION}"


.PHONY: tag
tag:
	@echo
	@echo "===> Git Tag Version: ${VERSION}"
	@git remote remove origin
	@git remote add origin https://${BUILD_ADMIN_USER}:${BUILD_ADMIN_KEY}@${CI_SERVER_HOST}/${CI_PROJECT_PATH}.git
	@git config user.email "${GITLAB_USER_EMAIL}"
	@git config user.name "${GITLAB_USER_NAME}"
	@git tag -a ${VERSION} -m "${VERSION}"
	@git push origin ${VERSION}
	@echo

.PHONY: publish
publish:
	docker push ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	docker push ${DOCKER_IMAGE_NAME}:latest
	sha1=$$($(docker-exec) sha1sum build/eve-bot.tar.gz | awk '{ print $$1 }') && \
	md5=$$($(docker-exec) md5sum build/eve-bot.tar.gz | awk '{ print $$1 }') && \
	curl --fail -H 'X-JFrog-Art-Api:${JFROG_API_KEY}' \
		-H "X-Checksum-Sha1:$$sha1" \
		-H "X-Checksum-MD5:$$md5" \
		-T build/eve-bot.tar.gz \
		"https://unanet.jfrog.io/unanet/generic/eve-bot/eve-bot-${VERSION}.tar.gz"	

# .PHONY: release
# release:
# 	sha1=$$()
# 	@echo
# 	@echo "===> Generate ${GIT_TAG} Changelog..."
# 	@git log v0.4.0...${GIT_TAG} --pretty=format:'1. [view commit](${CI_PROJECT_URL}/-/commit/%H)	%cn	`%s`	(%ci)' --reverse | tee CHANGELOG.md
# 	@echo
# 	@echo "===> Uploading ${GIT_TAG} Changelog..."
# 	uploadedURL=$$($(docker-exec) curl --request POST --header "PRIVATE-TOKEN: ${BUILD_ADMIN_KEY}" --form "file=@CHANGELOG.md" ${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/uploads | jq .url | sed -e 's/^"//' -e 's/"$$//')") && \
# 	curl --request POST --header "PRIVATE-TOKEN: ${BUILD_ADMIN_KEY}" \
# 		--form "description=Changelog File: [CHANGELOG.md]($$uploadedURL)" ${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/repository/tags/${GIT_TAG}/release
