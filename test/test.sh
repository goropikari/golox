#!/bin/bash

set -e

TEST_DIR=$(cd $(dirname $0); pwd)
for i in $(ls ${TEST_DIR}/*.lox); do
    echo $i
    ./golox $i
    echo pass
done
