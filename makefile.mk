# (c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

# Makes a recipe passed to a single invocation of the shell.
.ONESHELL:

MAKEFILE_PATH:=$(abspath $(dir $(lastword $(MAKEFILE_LIST))))
SHELLCHECK_VERSION?=v0.6.0

GO_SOURCES:=$(wildcard *.go)
GO_TEST_SOURCES:=$(wildcard *test.go)
BASH_SOURCES:=$(wildcard *.sh)
COVERAGE_DIR:=$(CURDIR)/coverage
COVERAGE_HTML_DIR:=$(COVERAGE_DIR)/html

COVERAGE_ARTIFACT:=${COVERAGE_HTML_DIR}/main.html
LINT_ARTIFACT:=._gometalinter
TEST_ARTIFACT:=${COVERAGE_DIR}/coverage.out
BASH_ARTIFACT:=._shellcheck

YELLOW:=\033[0;33m
GREEN:=\033[0;32m
RED:=\033[0;31m
NC:=\033[0m
NC_DIR:=: $(CURDIR)$(NC)

.PHONY: goimports gofmt clean-lint lint clean-test test clean-coverage \
	coverage clean-docker-shellcheck docker-shellcheck
# Stop prints each line of the recipe.
.SILENT:


goimports: ${GO_SOURCES}
	echo "${YELLOW}Running goimports${NC_DIR}" && \
	goimports -w $^


gofmt: ${GO_SOURCES}
	echo "${YELLOW}Running gofmt${NC_DIR}" && \
	gofmt -w -s $^


clean-test:
	rm -rf $(dir ${TEST_ARTIFACT})

test: ${TEST_ARTIFACT}
${TEST_ARTIFACT}: ${GO_SOURCES}
	if [ -n "${GO_TEST_SOURCES}" ]; then
		{ echo "${YELLOW}Running go test${NC_DIR}" && \
		  mkdir -p $(dir ${TEST_ARTIFACT}) && \
		  go test -coverprofile=$@ -v && \
		  echo "${GREEN}TEST PASSED${NC}"; } || \
		{ $(MAKE) --makefile=$(lastword $(MAKEFILE_LIST)) clean-test && \
          echo "${RED}TEST FAILED${NC}" && \
		  exit 1; }
	fi


clean-coverage: clean-test
	rm -rf $(dir ${COVERAGE_ARTIFACT})

coverage: ${COVERAGE_ARTIFACT}
${COVERAGE_ARTIFACT}: ${TEST_ARTIFACT}
	if [ -e "$<" ]; then
		echo "${YELLOW}Running go tool cover${NC_DIR}" && \
		mkdir -p $(dir ${COVERAGE_ARTIFACT}) && \
		go tool cover -html=$< -o $@ && \
		echo "${GREEN}Generated: $@${NC}"
	fi


clean-lint:
	rm -f ${LINT_ARTIFACT}

lint: ${LINT_ARTIFACT}
${LINT_ARTIFACT}: ${MAKEFILE_PATH}/gometalinter.json ${GO_SOURCES}
	echo "${YELLOW}Running go lint${NC_DIR}" && \
    (cd ${MAKEFILE_PATH} && \
	 procs=$$(expr $$(grep -c ^processor /proc/cpuinfo) '*' 2 '-' 1) && \
	 gometalinter \
		--config gometalinter.json \
		--concurrency=$${procs} \
		"$$(realpath --relative-to ${MAKEFILE_PATH} ${CURDIR})/.") && \
    touch $@


clean-docker-shellcheck:
	rm -f ${BASH_ARTIFACT}

docker-shellcheck: ${BASH_ARTIFACT}
${BASH_ARTIFACT}: ${BASH_SOURCES}
	echo "${YELLOW}Running shellcheck${NC_DIR}" && \
	docker run --rm \
		-v "${CURDIR}:/mnt" \
		koalaman/shellcheck:${SHELLCHECK_VERSION} -x $^ && \
	touch $@
