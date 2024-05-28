#!/bin/bash

export TARBALL_URL=
export RUN_TESTS=
export ALLOW_TESTS_FAILURES=
export CI_ENV_REPO=
export CI_ENV_BRANCH=

# Configure test run parameters
export TARBALL_URL=https://dev.mysql.com/get/Downloads/MySQL-8.4/mysql-8.4.0-linux-glibc2.17-x86_64-minimal.tar.xz
# export TARBALL_URL=https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-8.0.37-linux-glibc2.17-x86_64-minimal.tar.xz

#export TARBALL_URL=https://downloads.percona.com/downloads/Percona-Server-8.0/Percona-Server-8.0.36-28/binary/tarball/Percona-Server-8.0.36-28-Linux.x86_64.glibc2.17-minimal.tar.gz

export RUN_TESTS=YES
export ALLOW_TESTS_FAILURES=YES

# It is also possible to specify custom ci-env
#export CI_ENV_REPO=https://github.com/kamil-holubicki/orchestrator-ci-env.git
#export CI_ENV_BRANCH=DISTMYSQL-140

script/dock system
