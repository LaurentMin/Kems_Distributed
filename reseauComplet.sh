#!/bin/bash
echo "Making named pipes"
mkfifo /tmp/in_Debug

mkfifo /tmp/in_A1 /tmp/out_A1
mkfifo /tmp/in_C1 /tmp/out_C1

mkfifo /tmp/in_A2 /tmp/out_A2
mkfifo /tmp/in_C2 /tmp/out_C2

mkfifo /tmp/in_A3 /tmp/out_A3
mkfifo /tmp/in_C3 /tmp/out_C3
 
 echo "Starting A1"
./app -n A1 < /tmp/in_A1 >> /tmp/out_A1 &
sleep 1
echo "Starting C1"
./ctl -n C1 < /tmp/in_C1 >> /tmp/out_C1 &
sleep 1

echo "Starting A2"
./app -n A2 < /tmp/in_A2 >> /tmp/out_A2 &
sleep 1
echo "Starting C2"
./ctl -n C2 < /tmp/in_C2 >> /tmp/out_C2 &
sleep 1

echo "Starting A3"
./app -n A3 < /tmp/in_A3 >> /tmp/out_A3 &
sleep 1
echo "Starting C3"
./ctl -n C3 < /tmp/in_C3 >> /tmp/out_C3 &
sleep 1

echo "Starting Network"
cat /tmp/out_A1 | tee -a /tmp/in_C1 >> /tmp/in_Debug &
cat /tmp/out_C1 | tee -a /tmp/in_A1 /tmp/in_C3 >> /tmp/in_C2 &

cat /tmp/out_A2 | tee -a /tmp/in_C2 >> /tmp/in_Debug &
cat /tmp/out_C2 | tee -a /tmp/in_A2 /tmp/in_C1 >> /tmp/in_C3 &

cat /tmp/out_A3 | tee -a /tmp/in_C3 >> /tmp/in_Debug &
cat /tmp/out_C3 | tee -a /tmp/in_A3 /tmp/in_C2 >> /tmp/in_C1 &

echo "Everything running. (start a display and a player to begin)"