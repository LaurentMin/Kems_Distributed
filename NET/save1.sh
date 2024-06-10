#!/bin/bash

echo "Making named pipes..."
mkfifo /tmp/in_A2 /tmp/out_A2
mkfifo /tmp/in_C2 /tmp/out_C2

mkfifo /tmp/in_N2 /tmp/out_N2

echo "Adding pipes to network..."
cat /tmp/out_N1 | tee -a /tmp/in_N2 >> /tmp/in_C1 &
cat /tmp/out_N2 | tee -a /tmp/in_N1 >> /tmp/in_C2 &

cat /tmp/out_A1 | tee -a /tmp/in_C1 >> /tmp/in_Debug &
cat /tmp/out_C1 | tee -a /tmp/in_A1 >> /tmp/in_N1 &

cat /tmp/out_A2 | tee -a /tmp/in_C2 >> /tmp/in_Debug &
cat /tmp/out_C2 | tee -a /tmp/in_A2 >> /tmp/in_N2 &

echo "Starting App..."
# ../app -n A2 < /tmp/in_A2 >> /tmp/out_A2 &
sleep 1
echo "Starting Controller..."
# ../ctl -n C2 < /tmp/in_C2 >> /tmp/out_C2 &

echo "Starting network node..."
./net -n N2 -a N1 < /tmp/in_N2 >> /tmp/out_N2 &

echo "Network built."