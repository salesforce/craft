#!/bin/bash

# Build craft

set -e

BUILD_USER=${BUILD_USER:-"${USER}@${HOSTNAME}"}
BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}
VERBOSE=${VERBOSE:-}

repo_path="craft/cmd"

version=`git tag --points-at HEAD`
revision=$( git rev-parse --short HEAD 2> /dev/null || echo 'unknown' )
branch=$( git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown' )
go_version=$( go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/' )


# go 1.4 requires ldflags format to be "-X key value", not "-X key=value"
# ldseparator here is for cross compatibility
ldseparator="="
if [ "${go_version:0:3}" = "1.4" ]; then
	ldseparator=" "
fi

ldflags="
  -X ${repo_path}/base.Version${ldseparator}${version}
  -X ${repo_path}/base.Revision${ldseparator}${revision}
  -X ${repo_path}/base.Branch${ldseparator}${branch}
  -X ${repo_path}/base.BuildUser${ldseparator}${BUILD_USER}
  -X ${repo_path}/base.BuildDate${ldseparator}${BUILD_DATE}
  -X ${repo_path}/base.GoVersion${ldseparator}${go_version}"

echo ">>> Building craft..."

if [ -n "$VERBOSE" ]; then
  echo "Building with -ldflags $ldflags"
fi

GOBIN=$PWD go build -ldflags "${ldflags}" -o bin/craft main.go
GOBIN=$PWD env GOOS=linux GOARCH=amd64 go build -ldflags "${ldflags}" -o bin/craft_linux main.go
GOBIN=$PWD env GOOS=darwin GOARCH=amd64 go build -ldflags "${ldflags}" -o bin/craft_darwin main.go
GOBIN=$PWD env GOOS=windows GOARCH=amd64 go build -ldflags "${ldflags}" -o bin/craft_windows main.go
exit 0
