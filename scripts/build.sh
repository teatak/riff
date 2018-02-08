#!/usr/bin/env bash
set -e

cd `dirname $0`
cd ..

BIN="riff"
Type="release"
VERSION="$(cat version)"
GITSHA="$(git rev-parse HEAD)"
GITBRANCH="$(git rev-parse --abbrev-ref HEAD)"

# Determine the arch/os combos we're building for
XC_OS=${XC_OS:-"linux darwin windows freebsd openbsd solaris"}
XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
LDFLAGS="-X github.com/gimke/riff/common.Type=${Type} -X github.com/gimke/riff/common.GitSha=${GITSHA} -X github.com/gimke/riff/common.GitBranch=${GITBRANCH} -X github.com/gimke/riff/common.Version=${VERSION}"

# Delete the old dir
rm -rf bin/*
rm -rf pkg/*
mkdir -p bin/

if [ "${DEV}x" != "x" ]; then
    XC_OS=("$(go env GOOS)")
    XC_ARCH=("$(go env GOARCH)")
fi

# Build!

for OS in ${XC_OS}; do
    for ARCH in ${XC_ARCH}; do
        if [ ${OS}/${ARCH} == "darwin/arm" ]; then
            continue
        fi
        if [ ${OS}/${ARCH} == "windows/arm" ]; then
            continue
        fi
        if [ ${OS}/${ARCH} == "solaris/arm" ]; then
            continue
        fi
        if [ ${OS}/${ARCH} == "solaris/386" ]; then
            continue
        fi
        echo "Building ${OS}/${ARCH}"
        NAME="${BIN}"
        if [ ${OS} == "windows" ]; then
            NAME="${BIN}.exe"
        fi
        if ! CGO_ENABLED=0 GOOS="${OS}" GOARCH="${ARCH}" go build -ldflags "${LDFLAGS}" -o ./pkg/"${OS}"_"${ARCH}"/"${NAME}" ./cmd/ > /dev/null 2>&1; then
            echo -e "\033[31;1mBuilding ${OS}/${ARCH} error\033[0m"
            continue
        fi
    done
done

DEV_PLATFORM="./pkg/$(go env GOOS)_$(go env GOARCH)"
for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f); do
    cp ${F} bin/
done

# Zip and copy to the dist dir
for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "Packaging ${OSARCH}"

    pushd $PLATFORM >/dev/null 2>&1
    zip ../${OSARCH}.zip ./* >/dev/null 2>&1
    popd >/dev/null 2>&1
done

# Done!
echo -e "\033[32;1mBuild Done\033[0m"
