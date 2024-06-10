#!/bin/bash

echo "Clearing old cats and tees..."
killall tee 2> /dev/null
killall cat 2> /dev/null

echo "Making named pipes..."
mkfifo /tmp/in_A3 /tmp/out_A3
mkfifo /tmp/in_C3 /tmp/out_C3

mkfifo /tmp/in_N3 /tmp/out_N3

echo "Adding pipes to network..."
cat /tmp/out_N1 | tee -a /tmp/in_N2 /tmp/in_N3 >> /tmp/in_C1 &
cat /tmp/out_N2 | tee -a /tmp/in_N1 >> /tmp/in_C2 &
cat /tmp/out_N3 | tee -a /tmp/in_N1 >> /tmp/in_C3 &

cat /tmp/out_A1 | tee -a /tmp/in_C1 >> /tmp/in_Debug &
cat /tmp/out_C1 | tee -a /tmp/in_A1 >> /tmp/in_N1 &

cat /tmp/out_A2 | tee -a /tmp/in_C2 >> /tmp/in_Debug &
cat /tmp/out_C2 | tee -a /tmp/in_A2 >> /tmp/in_N2 &

cat /tmp/out_A3 | tee -a /tmp/in_C3 >> /tmp/in_Debug &
cat /tmp/out_C3 | tee -a /tmp/in_A3 >> /tmp/in_N3 &

echo "Starting App..."
../app -n A3 < /tmp/in_A3 >> /tmp/out_A3 &
sleep 1
echo "Starting Controller..."
../ctl -n C3 < /tmp/in_C3 >> /tmp/out_C3 &
sleep 1
echo "Starting network node..."
./net -n N3 -a N1 < /tmp/in_N3 >> /tmp/out_N3 &

echo "Network built."