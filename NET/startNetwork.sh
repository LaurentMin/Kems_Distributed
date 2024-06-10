#!/bin/bash

echo "Making named pipes..."
mkfifo /tmp/in_N1 /tmp/out_N1

echo "Starting first node..."
./net -n N1 -a N1 < /tmp/in_N1 >> /tmp/out_N1 &

echo "Network created."
