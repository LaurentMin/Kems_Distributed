#!/bin/bash

scripts=(
    "startNetwork.sh"
    "add1.sh"
    "add2.sh"
)

for script in "${scripts[@]}"; do
    ./$script
    sleep 20
done
