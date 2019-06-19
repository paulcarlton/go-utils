#!/bin/bash
# (c) Copyright 2019 Hewlett Packard Enterprise Development LP

set -euo pipefail

exit_handler(){
   local ret="${?}"
   echo "Caught test Exit, code: ${ret}"
   tests_teardown
   exit ${ret}
}

# Setup signal handlers
trap 'exit_handler' EXIT

is_pod_running() {
    local pod_name="${1}"
    kubectl get pods --namespace "${K8S_NAMESPACE}" --field-selector status.phase=Running \
        -l "app.kubernetes.io/name=${pod_name},app.kubernetes.io/instance=${pod_name}"
}

get_pod() {
    local pod_name="${1}"
    if [ -z "$(is_pod_running "${pod_name}")" ]; then
        echo "no pods running with name: ${pod_name}"
        exit 1
    fi
    kubectl get pods --namespace "${K8S_NAMESPACE}" --field-selector status.phase=Running \
        -l "app.kubernetes.io/name=${pod_name},app.kubernetes.io/instance=${pod_name}" \
        -o jsonpath="{.items[0].metadata.name}"
}

wait_for_pod() {
    local pod_name="${1}"
    local time_out="${2:-300}"
    echo "waiting for ${pod_name} to be ready"
    while (( time_out > 0 )); do
        if [ -n "$(is_pod_running "${pod_name}")" ]; then
            echo "${pod_name} running"
            return
        fi
        echo "waiting for pod ${pod_name}"
        sleep 10
        (( time_out-=10 ))
    done
    echo "${pod_name} not ready"
    exit 1
}

kill_portforwards() {
    if [ ! -e "${WORK_DIR}/port-forwards.txt" ]; then
        return
    fi
    local pid
    while read -r pid
    do
        if kill -0 "${pid}"; then
            kill -SIGKILL "${pid}"
        fi
    done < "${WORK_DIR}/port-forwards.txt"
    rm -f "${WORK_DIR}/port-forwards.txt"
}

setup_port_forward() {
    local pod_name="${1}"
    local pod_id
    pod_id="$(get_pod "${pod_name}")"
    if [ -z "${pod_id}" ]; then
        echo "no pod named ${pod_name} found"
        exit 1
    fi
    local port="${2}"
    if [ -z "${port}" ]; then
        echo "no port number provided"
        exit 1
    fi
    kubectl port-forward -n "${K8S_NAMESPACE}" "${pod_id}" "${port}" &
    echo "${!}" >> "${WORK_DIR}/port-forwards.txt"
    sleep 5
}

test_init() {
    # Initialise result file
    results_file_name="results-$(cat /proc/sys/kernel/random/uuid).txt"
}

report_results () {
    echo "======================== Results ============================="
    if [ -n "${DEV_MODE:-}" ]; then
        echo "not running in a collie test, no subunit update"
        cat "${WORK_DIR}/${results_file_name}"
        return
    fi
    subunit_file_name=${subunit_file_name:-subunit-file} # Should be set by collie
    echo "subunit_file_name=${subunit_file_name}"
    while read -r line
    do
        test_name=$(echo "${line}" | awk -F ':' '{print$1}')
        result=$(echo "${line}" | awk -F ':' '{print$2}')
        echo "processing line ${line}"
        echo "${test_name}"
        echo "${result}"
        if [ "${result}" == "pass" ]; then
            subunit-output --success "${test_name}" >> "${subunit_file_name}"
        elif [ "${result}" == "fail" ]; then
            subunit-output --fail "${test_name}"  >> "${subunit_file_name}"
        elif [ "${result}" == "skip" ]; then
            subunit-output --skip "${test_name}"  >> "${subunit_file_name}"
        elif [ "${result}" == "exists" ]; then
            bin/subunit-output --exists "${test_name}"  >> "${subunit_file_name}"
        else
            echo "result of the test case is neither pass or fail or skip"
        fi
    done < "${WORK_DIR}/${results_file_name}"
}

set_test_result() {
    local test_name="${1}"
    local result="${2:-fail}"
    echo "${test_name}:${result}" >> "${WORK_DIR}/${results_file_name}"
}

test_list() {
    declare -F | awk '{print $NF}' | sort | grep "^test_.[0-9]"
}

set_outcome() {
   if grep -q ":fail$" "${WORK_DIR}/${results_file_name}"; then
      exit 1
   fi
}

run_tests() {
    test_init
    tests_setup
    for t in $(test_list)
    do
        local result=fail
        if ${t}; then
            result=pass
        fi
        set_test_result "${t}" "${result}"
    done
    tests_teardown
    trap - EXIT
    report_results
    set_outcome
}

# The following functions are default implementations designed to be overridden by the user
tests_setup() {
    echo "no user defined tests setup"
}

tests_teardown() {
    echo "no user defined tests teardown"
}
