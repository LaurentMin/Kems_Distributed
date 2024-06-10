#!/bin/bash

echo "Making named pipes..."
mkfifo /tmp/in_A1 /tmp/out_A1
mkfifo /tmp/in_C1 /tmp/out_C1

mkfifo /tmp/in_N1 /tmp/out_N1

echo "Adding pipes to network..."
cat /tmp/out_N0 | tee -a /tmp/in_N1 >> /tmp/in_C0 &
cat /tmp/out_N1 | tee -a /tmp/in_N0 >> /tmp/in_C1&

cat /tmp/out_A1 | tee -a /tmp/in_C1 >> /tmp/in_Debug &
cat /tmp/out_C1 | tee -a /tmp/in_A1 >> /tmp/in_N1 &

echo "Starting App..."
./app -n A1 < /tmp/in_A1 >> /tmp/out_A1 &
sleep 0
echo "Starting Controller..."
./ctl -n C1 < /tmp/in_C1 >> /tmp/out_C1 &

echo "Starting node..."
./net -n N1 -a N0 < /tmp/in_N1 >> /tmp/out_N1 &

echo "Network built."