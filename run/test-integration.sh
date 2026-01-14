#!/bin/bash

export TARBALL_URL=
export RUN_TESTS=
export ALLOW_TESTS_FAILURES=
export CI_ENV_REPO=
export CI_ENV_BRANCH=

# Configure test run parameters
export TARBALL_URL=https://downloads.percona.com/downloads/Percona-Server-8.4/Percona-Server-8.4.7-7/binary/tarball/Percona-Server-8.4.7-7-Linux.x86_64.glibc2.35-minimal.tar.gz
export RUN_TESTS=YES
export ALLOW_TESTS_FAILURES=YES

script/dock test
