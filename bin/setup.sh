#!/usr/bin/env bash
# (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

# Set versions of software required
metalinter_version=2.0.12
golang_version=1.12.7

function usage()
{
    echo "USAGE: ${0##*/}"
    echo "Install software required for golang project"
}

function args() {
    while [ $# -gt 0 ]
    do
        case "$1" in
            "--help") usage; exit;;
            "-?") usage; exit;;
            *) usage; exit;;
        esac
    done
}

function install_glide() {
    echo "Installing glide"
    curl https://glide.sh/get | sh
}

function install_gometalinter() {
    echo "Installing gometalinter version: ${metalinter_version}"
    set -e
    pushd /tmp
    curl -qL -o gometalinter.tar.gz "https://github.com/alecthomas/gometalinter/releases/download/v${metalinter_version}/gometalinter-${metalinter_version}-linux-amd64.tar.gz" &&
     mkdir -p gometalinter &&
     tar xvzf gometalinter.tar.gz --strip-components=1 -C gometalinter &&
     rm ./gometalinter/COPYING ./gometalinter/README.md &&
     cp ./gometalinter/* "${PROJECT_BIN_ROOT}" &&
     rm -f gometalinter.tar.gz &&
     rm -rf ./gometalinter &&
    popd >/dev/null
    set +e
    "${PROJECT_BIN_ROOT}"/gometalinter install
}

function install_golang() {
    echo "Installing golang version: ${golang_version}"
    pushd /tmp >/dev/null
    # shellcheck disable=SC1090
    curl -qL -O "https://storage.googleapis.com/golang/go${golang_version}.linux-amd64.tar.gz" &&
      tar xfa go${golang_version}.linux-amd64.tar.gz &&
      rm -rf "${PROJECT_BIN_ROOT}/go" &&
      mv go "${PROJECT_BIN_ROOT}" &&
      source "${SCRIPT_DIR}/env.sh" &&
    popd >/dev/null

    pushd "${GOROOT}/src/go/types" > /dev/null
    echo "Installing gotype linter"
    go build gotype.go
    cp gotype "${GOBIN}"
    popd >/dev/null
}

function install_godocdown() {
    echo "installing godocdown"
    go get github.com/robertkrimen/godocdown/godocdown
}

function make_local() {
    if [ ! -d "${PROJECT_BIN_ROOT}" ] ; then
        echo "Creating directory for ${PROJECT_NAME} software in ${PROJECT_BIN_ROOT}"
        mkdir -p "${PROJECT_BIN_ROOT}"
    fi
    # shellcheck disable=SC1090
    source "${SCRIPT_DIR}/env.sh"
}

SCRIPT_DIR="$(readlink -f "$(dirname "${0}")")"
# shellcheck disable=SC1090
if ! source "${SCRIPT_DIR}/env.sh"; then
    exit 1
fi

args "${@}"

echo "Running setup script to setup software for ${PROJECT_NAME}"

# Remove any legacy installs
rm -rf "${PROJECT_DIR}/bin/local"

make_local

gometalinter --version 2>&1 | grep $metalinter_version >/dev/null
ret_code="${?}"
if [[ "${ret_code}" != "0" || ! -e "${PROJECT_BIN_ROOT}/gometalinter" ]] ; then
    install_gometalinter
    gometalinter --version 2>&1 | grep $metalinter_version >/dev/null
    ret_code="${?}"
    if [ "${ret_code}" != "0" ] ; then
        echo "Failed to install gometalinter"
        exit 1
    fi
fi


go version 2>&1 | grep $golang_version >/dev/null
ret_code="${?}"
if [[ "${ret_code}" != "0"  || "${GOROOT}" != "${PROJECT_BIN_ROOT}/go" ]] ; then
    install_golang
    go version 2>&1 | grep $golang_version >/dev/null
    ret_code="${?}"
    if [ "${ret_code}" != "0" ] ; then
        echo "Failed to install golang"
        exit 1
    fi
fi

godocdown >/dev/null 2>&1
ret_code="${?}"
if [[ "${ret_code}" == "127" || "${GOBIN}" != "${PROJECT_BIN_ROOT}" ]] ; then
    install_godocdown
    godocdown >/dev/null 2>&1
    if [ "$?" == "127" ] ; then
        echo "Failed to install godocdown"
        exit 1
    fi
fi

glide >/dev/null 2>&1
ret_code="${?}"
if [[ "${ret_code}" == "127" || "${GOBIN}" != "${PROJECT_BIN_ROOT}" ]] ; then
    install_glide
    glide >/dev/null 2>&1
    if [ "$?" == "127" ] ; then
        echo "failed to install glide"
        exit 1
    fi
fi
