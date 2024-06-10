#!/bin/bash

echo "Clearing old cats and tees..."
killall tee 3> /dev/null
killall cat 3> /dev/null

echo "Making named pipes..."
mkfifo /tmp/in_N3 /tmp/out_N3

echo "Adding pipes to network..."
cat /tmp/out_N1 | tee -a /tmp/in_N2 &
cat /tmp/out_N2 | tee -a /tmp/in_N1 >> /tmp/in_N3 &
cat /tmp/out_N3 | tee -a /tmp/in_N2 &

echo "Starting node..."
./net -n N3 -a N2 < /tmp/in_N3 >> /tmp/out_N3 &

echo "Network built."