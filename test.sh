#!/usr/bin/env bash

set -e

test_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/tools

. "${test_dir}"/test_compatibility.sh
