#!/bin/bash

echo "Clearing old cats and tees..."
killall tee 3> /dev/null
killall cat 3> /dev/null

echo "Making named pipes..."
mkfifo /tmp/in_N6 /tmp/out_N6
mkfifo /tmp/in_N7 /tmp/out_N7

echo "Adding pipes to network..."
cat /tmp/out_N1 | tee -a /tmp/in_N2 >> /tmp/in_N4 &
cat /tmp/out_N2 | tee -a /tmp/in_N1 /tmp/in_N5 >> /tmp/in_N3 &
cat /tmp/out_N3 | tee -a /tmp/in_N2 &
cat /tmp/out_N4 | tee -a /tmp/in_N1 >> /tmp/in_N6 &
cat /tmp/out_N5 | tee -a /tmp/in_N2 >> /tmp/in_N7 &
cat /tmp/out_N6 | tee -a /tmp/in_N4 &
cat /tmp/out_N7 | tee -a /tmp/in_N5 &

echo "Starting node..."
./net -n N7 -a N5 < /tmp/in_N7 >> /tmp/out_N7 &
./net -n N6 -a N4 < /tmp/in_N6 >> /tmp/out_N6 &

echo "Network built."