#!/usr/bin/env bash

short_git_commit=$1
git_branch=$2
git_tag=$3

semver_parts=(${git_tag//./ })
major=${semver_parts[0]}
minor=${semver_parts[1]}
patch=${semver_parts[2]}

count=$(git rev-list HEAD ^${git_tag} --ancestry-path ${git_tag} --count)

case $git_branch in
master)
	version=${major}.$((minor + 1)).0
	;;
patch)
	version=${major}.${minor}.$((patch + 1))
	;;	
*)
	version=${major}.${minor}.${patch}-${short_git_commit}-${count}
	;;
esac

echo ${version}
exit 0
