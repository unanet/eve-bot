#!/usr/bin/env bash

version=$1
feature=$2
git_branch=$(git branch | sed -n -e 's/^\* \(.*\)/\1/p')
git_tag=$(git tag -l --merged master --sort='-*authordate' | head -n1)
count=$(git rev-list HEAD ^${git_tag} --ancestry-path ${git_tag} --count)
curr_tag=$(git describe --abbrev=0)

rtnval=""
case $git_branch in
master)
	rtnval="${version}"
	;;
*)
	rtnval=${curr_tag}-${count}-${feature}
	;;
esac

echo ${rtnval}
exit 0