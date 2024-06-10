#!/bin/bash

echo "Clearing old cats and tees..."
killall tee 2> /dev/null
killall cat 2> /dev/null

echo "Making named pipes..."
mkfifo /tmp/in_N2 /tmp/out_N2

echo "Adding pipes to network..."
cat /tmp/out_N1 | tee -a /tmp/in_N2 &
cat /tmp/out_N2 | tee -a /tmp/in_N1 &

echo "Starting node..."
./net -n N2 -a N1 < /tmp/in_N2 >> /tmp/out_N2 &

echo "Network built."
