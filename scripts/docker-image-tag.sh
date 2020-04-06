#!/usr/bin/env bash

version=$1
feature=$2
git_branch=$3
git_tag=$4
count=$(git rev-list HEAD ^${git_tag} --ancestry-path ${git_tag} --count)

rtnval=""
case $git_branch in
master)
	rtnval="${version}"
	;;
*)
	rtnval=${git_tag}-${count}-${feature}
	;;
esac

echo ${rtnval}
exit 0