#!/usr/bin/env bash

short_git_commit=$(git rev-parse --short=10 HEAD)
git_branch=$(git branch | sed -n -e 's/^\* \(.*\)/\1/p')
#git_tag=$(git tag -l --merged master --sort='-*authordate' | head -n1)
git_tag=$(git describe --abbrev=0)

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
