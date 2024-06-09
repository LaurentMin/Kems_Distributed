#!/bin/bash

echo "Making named pipes..."
mkfifo /tmp/in_N0 /tmp/out_N0

echo "Starting first node..."
./net -n N0 -a N0 < /tmp/in_N0 >> /tmp/out_N0 &

echo "Network created."