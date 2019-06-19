#!/bin/bash
# (c) Copyright 2019 Hewlett Packard Enterprise Development LP

set -euo pipefail

if [ -n "${VERBOSE_MODE:-}" ];then
    export PS4='+(`basename ${BASH_SOURCE}`:${LINENO}): ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'
    set -x
fi

run_docker_tests () {
    docker run --rm --net=host \
        -v "${WORK_DIR}":"${WORK_DIR}" \
        -v "${HOME}"/.minikube:"${HOME}"/.minikube \
        -e HTTP_PROXY -e HTTPS_PROXY -e NO_PROXY -e VERBOSE_MODE -e WORK_DIR \
        docker.artifactory.zing.hpelabs.net/ncs-qe/rdaclient-tests:latest
    ret="${?}"
    exit ${ret}
}

run_collie_tests () {
    # pull is neccessary because we are using latest and when using latests
    # docker does not verify that the local copy of latest is upto date with
    # the remote version of latest.
    docker pull docker.release.zing.hpelabs.net/ncs-qe/colliectl:latest

    docker run --rm --net=host \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v "${WORK_DIR}"/collie_logs:/tmp/collie_logs \
        -v "${WORK_DIR}":/tmp/rda \
        -v "${HOME}"/.minikube:"${HOME}"/.minikube \
        -e http_proxy -e https_proxy -e no_proxy \
        -e COLLIE_LOGGING_LEVEL=DEBUG \
        --env-file "${WORK_DIR}/env.list" \
        staging-docker.artifactory.zing.hpelabs.net/ncs-qe/colliectl:latest \
        -i "${test_image}" \
        -s minikube -t minikube -r smoke --skip_validation -d
    ret="${?}"
    if [ -e "${WORK_DIR}/collie_logs/collie_run_rdaclient-tests_smoke.log" ]; then
        echo "collie logs..."
        cat "${WORK_DIR}/collie_logs/collie_run_rdaclient-tests_smoke.log"
    fi
    exit ${ret}
}

test_run() {
    if [ -z "${DEV_MODE:-}" ]; then
         if [ -n "${NO_COLLIE:-}" ]; then
            run_docker_tests
         else
            test_image="${1}"
            run_collie_tests
         fi
    else
        # shellcheck disable=SC1090
        source "${TESTER_SCRIPT_DIR}/test-utils.sh"
        # shellcheck disable=SC1090
        source "${1}"
        # shellcheck disable=SC2091
        if [ -e "${WORK_DIR}/env.list" ]; then
            $(awk '!/^#/{print "export "$0}' "${WORK_DIR}/env.list")
        fi
        run_tests
    fi
}

main() {
    arg="${1:-}"
    if [ -z "${arg}" ]; then
        usage
        exit 1
    fi

    test_run "${arg}"
}

usage() {
    echo "usage: ${0} <docker image for collie to run> | <shell script defining tests>"
}

main "${@}"