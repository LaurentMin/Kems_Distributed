#!/bin/bash

echo "Making named pipes..."
mkfifo /tmp/in_N1 /tmp/out_N1

echo "Adding pipes to network..."
cat /tmp/out_N0 | tee -a /tmp/in_N1 &
cat /tmp/out_N1 | tee -a /tmp/in_N0 &

echo "Starting node..."
./net -n N1 -a N0 < /tmp/in_N1 >> /tmp/out_N1 &

echo "Network built."