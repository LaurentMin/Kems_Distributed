#!/bin/bash

echo "Clearing old cats and tees..."
killall tee 2> /dev/null
killall cat 2> /dev/null

echo "Adding pipes to network..."
cat /tmp/out_N0 | tee -a /tmp/in_N2 /tmp/in_N3 &
cat /tmp/out_N2 | tee -a /tmp/in_N0 /tmp/in_N3 >> /tmp/in_C2 &
cat /tmp/out_N3 | tee -a /tmp/in_N0 /tmp/in_N2 >> /tmp/in_C3 &

cat /tmp/out_A2 | tee -a /tmp/in_C2 >> /tmp/in_Debug &
cat /tmp/out_C2 | tee -a /tmp/in_A2 >> /tmp/in_N2 &

cat /tmp/out_A3 | tee -a /tmp/in_C3 >> /tmp/in_Debug &
cat /tmp/out_C3 | tee -a /tmp/in_A3 >> /tmp/in_N3 &


echo "Network built."