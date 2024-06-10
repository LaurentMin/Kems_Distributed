#!/bin/bash

echo "Making named pipes..."
mkfifo /tmp/in_Debug

mkfifo /tmp/in_A0 /tmp/out_A0
mkfifo /tmp/in_C0 /tmp/out_C0

mkfifo /tmp/in_N0 /tmp/out_N0

echo "Adding pipes to network..."
cat /tmp/out_N0 | tee -a /tmp/in_C0 &

cat /tmp/out_A0 | tee -a /tmp/in_C0 >> /tmp/in_Debug &
cat /tmp/out_C0 | tee -a /tmp/in_A0 >> /tmp/in_N0 &

echo "Starting App..."
./app -n A0 < /tmp/in_A0 >> /tmp/out_A0 &
sleep 0
echo "Starting Controller..."
./ctl -n C0 < /tmp/in_C0 >> /tmp/out_C0 &
sleep 0
echo "Starting first network node..."
./net -n N0 -a N0 < /tmp/in_N0 >> /tmp/out_N0 &

echo "Network created."