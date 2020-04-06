#!/usr/bin/env bash


url=$(curl -v --request POST --header "PRIVATE-TOKEN: ${BUILD_ADMIN_KEY}" --data-urlencode file@CHANGELOG.md ${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/uploads | jq .url | sed -e 's/^"//' -e 's/"$//')

echo $url