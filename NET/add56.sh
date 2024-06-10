#!/bin/bash

echo "Clearing old cats and tees..."
killall tee 2> /dev/null
killall cat 2> /dev/null

echo "Making named pipes..."
mkfifo /tmp/in_N5 /tmp/out_N5
mkfifo /tmp/in_N6 /tmp/out_N6

echo "Adding pipes to network..."
cat /tmp/out_N0 | tee -a /tmp/in_N1 >> /tmp/in_N3 &
cat /tmp/out_N1 | tee -a /tmp/in_N0 /tmp/in_N4 >> /tmp/in_N2 &
cat /tmp/out_N2 | tee -a /tmp/in_N1 &
cat /tmp/out_N3 | tee -a /tmp/in_N0 >> /tmp/in_N5 &
cat /tmp/out_N4 | tee -a /tmp/in_N1 >> /tmp/in_N6 &
cat /tmp/out_N5 | tee -a /tmp/in_N3 &
cat /tmp/out_N6 | tee -a /tmp/in_N4 &

echo "Starting node..."
./net -n N6 -a N4 < /tmp/in_N6 >> /tmp/out_N6 &
./net -n N5 -a N3 < /tmp/in_N5 >> /tmp/out_N5 &

echo "Network built."