#!/bin/bash

scripts=(
    "startNetwork.sh"
    "add1.sh"
    "add2.sh"
    "add3.sh"
    "add4.sh"
    "add56.sh"
)

for script in "${scripts[@]}"; do
    ./$script
    sleep 10
done
