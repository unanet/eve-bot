#!/usr/bin/env bash

PROJECT_DIR=$(PWD)
BIN_DIR="$PROJECT_DIR/bin"

if [ "$(ls -A $BIN_DIR)" ]; then
	cd $BIN_DIR
	find $BIN_DIR -mindepth 1 -maxdepth 1 -type f -exec basename {} \; | grep '.zip\|.tar\|.tar.gz'| grep -v 'SHA256SUM'|grep -v 'asc' | while read ARTIFACT; do
		echo
		echo "Verify Checksum ==>"		
		shasum -a 256 --check ./$ARTIFACT.SHA256SUM
		echo
		echo "Verify GPG Sig ==>"
        gpg --verify ./$ARTIFACT.SHA256SUM.asc
	done
	echo
fi

exit 0
