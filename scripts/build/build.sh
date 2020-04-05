#!/usr/bin/env bash

# Set initial project level variables
START_TIME=$SECONDS
TS=`/bin/date "+%Y-%m-%d-%H:%M:%S"`
PROJECT_DIR=$(PWD)
OUTPUT_DIR=${OUTPUT_DIR:=${PROJECT_DIR}/bin}

# Delete the old build output
rm -rf ${OUTPUT_DIR}/*
mkdir -p ${OUTPUT_DIR}/

# Check/Install Build Deps
if ! which gox >/dev/null; then
	echo "==> Installing gox..."
	go get -u github.com/mitchellh/gox
fi

if ! which gocyclo >/dev/null; then
	echo "==> Installing gocyclo..."
	go get -u github.com/fzipp/gocyclo
fi

# Set build variables
BUILD_USER=$(whoami)
BUILD_HOST=$(hostname)
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "$(git status --porcelain)" && echo "+CHANGES" || true)
GIT_COMMIT_AUTHOR=$(git show -s --format='%ae' ${GIT_COMMIT})
VERSION="${VERSION:-0.0.0}"
XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:-"linux darwin windows freebsd openbsd solaris"}
XC_EXCLUDE_OSARCH=${XC_EXCLUDE_OSARCH:-"!darwin/arm !darwin/386 !freebsd/386 !freebsd/amd64 !freebsd/arm !openbsd/386 !openbsd/amd64 !linux/arm !solaris/amd64 !windows/386 !windows/amd64"}
PRERELEASE="${PRERELEASE:-}"

# force statically linked binaries
export CGO_ENABLED=0

# If its dev mode, only build for the local dev environment
if [ "${BUILD_DEV}x" != "x" ]; then
	XC_OS=$(go env GOOS)
	XC_ARCH=$(go env GOARCH)
fi

# Allow LD_FLAGS to be appended during development compilations
LD_FLAGS="${LD_FLAGS} -X main.VersionPrerelease=${PRERELEASE} -X main.BuildHost=${BUILD_HOST} -X main.GitBranch=${GIT_BRANCH} -X main.Builder=${BUILD_USER} -X main.Version=${VERSION} -X main.BuildDate=${TS} -X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.GitCommitAuthor=${GIT_COMMIT_AUTHOR}"

function buildBinary() {
	gox \
		-os="${XC_OS}" \
		-arch="${XC_ARCH}" \
		-osarch="${XC_EXCLUDE_OSARCH}" \
		-ldflags "${LD_FLAGS}" \
		-output "bin/{{.OS}}_{{.Arch}}/${PWD##*/}" \
		${PROJECT_DIR}/cmd/eve-bot/
}

function printBuildResults() {
	echo
	echo "	VERSION: ${VERSION}"
	echo "	START_TIME: ${TS}"
	echo "	END_TIME: `/bin/date "+%Y-%m-%d-%H:%M:%S"`"
	echo "	BUILD_TIME (seconds): $(($SECONDS - $START_TIME))"
	echo "	BUILD_USER: ${BUILD_USER}"
	echo "	BUILD_HOST: ${BUILD_HOST}"
	echo "	PRE_RELEASE: ${PRERELEASE}"
	echo "	GIT_SHA: ${GIT_COMMIT}${GIT_DIRTY}"
	echo "	GIT_COMMIT_AUTHOR: ${GIT_COMMIT_AUTHOR}"
	echo "	GIT_BRANCH: ${GIT_BRANCH}"
	echo "	XC_ARCH: ${XC_ARCH}"
	echo "	XC_OS: ${XC_OS}"
	echo "	OUTPUT_DIR: ${OUTPUT_DIR}/"
	echo 
}

function packageRelease() {
	if [ "${BUILD_DEV}x" = "x" ]; then
		# Zip and copy to the dist dir
		echo "==> Packaging..."
		for PLATFORM in $(find ${OUTPUT_DIR}/ -mindepth 1 -maxdepth 1 -type d); do
			OSARCH=$(basename ${PLATFORM})
			echo "--> ${OSARCH}"

			pushd $PLATFORM >/dev/null 2>&1
			zip ../${OSARCH}_${VERSION}.zip ./*
			popd >/dev/null 2>&1
		done
	fi
}

function moveBinaryBin() {
	# Move all the compiled things to the $GOPATH/bin
	GOPATH=${GOPATH:-$(go env GOPATH)}
	case $(uname) in
	CYGWIN*)
		GOPATH="$(cygpath $GOPATH)"
		;;
	esac
	OLDIFS=$IFS
	IFS=: MAIN_GOPATH=($GOPATH)
	IFS=$OLDIFS

	# Create GOPATH/bin if it's doesn't exists
	if [ ! -d $MAIN_GOPATH/bin ]; then
		echo "==> Creating GOPATH/bin directory..."
		mkdir -p $MAIN_GOPATH/bin
	fi

	# Copy our OS/Arch to the bin/ directory
	DEV_PLATFORM="${OUTPUT_DIR}/$(go env GOOS)_$(go env GOARCH)"
	if [[ -d "${DEV_PLATFORM}" ]]; then
		for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f); do
			cp ${F} ${OUTPUT_DIR}/
			cp ${F} ${MAIN_GOPATH}/bin/
		done
	fi
}

# Entrypoint ento the script
if buildBinary; then
	moveBinaryBin
	packageRelease
	printBuildResults
	ls -hl ${OUTPUT_DIR}/
fi
