#!/bin/bash

export TARBALL_URL=
export RUN_TESTS=
export ALLOW_TESTS_FAILURES=
export CI_ENV_REPO=
export CI_ENV_BRANCH=

# Configure test run parameters
export TARBALL_URL=https://dev.mysql.com/get/Downloads/MySQL-8.4/mysql-8.4.0-linux-glibc2.17-x86_64-minimal.tar.xz
export RUN_TESTS=YES
export ALLOW_TESTS_FAILURES=YES

# It is also possible to specify custom ci-env
#export CI_ENV_REPO=https://github.com/kamil-holubicki/orchestrator-ci-env.git
#export CI_ENV_BRANCH=DISTMYSQL-140

script/dock system
