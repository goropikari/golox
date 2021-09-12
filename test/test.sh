#!/bin/bash

set -e

TEST_DIR=$(cd $(dirname $0); pwd)
for i in $(ls ${TEST_DIR}/*.tlps); do
    echo $i
    ./tlps $i
    echo pass
done
