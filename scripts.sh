#!/usr/bin/env bash

cmd=$1

if [[ -z "$cmd" ]];then
  echo "Usage scripts.sh [vendor]"
  exit 1
fi


vendor() {
  mods=$(find . -name "go.mod" | grep -v "vendor")

  root=$PWD
  for mod in $mods;do
    version=$(cat ${mod} | grep -e "^go " | awk -F' ' '{print $2}')
    echo "mod ${mod} version=go:${version}"
    dir=$(dirname "$mod")
    cd "${dir:2}" && rm -fr vendor && rm -fr go.sum && go mod tidy -compat=${version} && go mod vendor
    cd "${root}"
  done
}

case $cmd in
vendor)
  vendor
  ;;
*)
  echo "Usage scripts.sh [tag|vendor]"
  exit 1
  ;;
esac