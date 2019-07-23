(c) Copyright 2018-2019 Hewlett Packard Enterprise Development LP

# Developers Guide

This repository is intended to contain testing and development utilities as well as library code for use
with golang and python. To run tests and build any executables (non at present) type 'make' in top level
directory of the repository. This will also install any software you need on your workstation.

## Setup

clone into $GOPATH/src/github.com/paulcarlton/go-utils:

    cd $GOPATH/src/github.com/paulcarlton/go-utils
    git clone git@github.com:paulcarlton/go-utils.git
    cd utils

Optionally install required software versions in project's bin directory:

    . bin/env.sh
    setup.sh

This project requires the following software:

    glide version >= v0.13.2
    metalinter version = 2.0.12
    golang version = 1.11.3
    godocdown version = head

You can install these in the project bin directory using the 'setup.sh' script:

    . bin/env.sh
    setup.sh

The setup.sh script can safely be run at any time. It installs the required software in the <project-dir>bin/local.

## Development

The Makefile in the project's top level directory will compile, build and test all components.

    make check build

To run the build and test in a docker container, type:

    make

If changes are made to go sources you may need to perform a glide update, type:

    make glide-update

## Golang Utilities

The 'goutils' directory contains golang utility functions. To test:

    make glide
    make -C pkg/goutils --makefile=${PWD}/makefile.mk

## Golang Kubernetes Utilities

The 'k8sutils' directory contains golang utility functions related to Kubernetes. To test:

    make glide
    make -C pkg/k8sutils --makefile=${PWD}/makefile.mk
