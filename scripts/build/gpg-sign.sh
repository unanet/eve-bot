#!/usr/bin/env bash

PROJECT_DIR=$(PWD)
BIN_DIR="$PROJECT_DIR/bin"

if [ "$(ls -A $BIN_DIR)" ]; then
	cd $BIN_DIR
	find $BIN_DIR -mindepth 1 -maxdepth 1 -type f -exec basename {} \; | grep '.zip\|.tar\|.tar.gz'| grep -v 'SHA256SUM'|grep -v 'asc' | while read ARTIFACT; do
        shasum -a 256 $ARTIFACT >./$ARTIFACT.SHA256SUM
		gpg -ab ./$ARTIFACT.SHA256SUM
	done
fi

exit 0
