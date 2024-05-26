#!/bin/bash
# Retrieve command line argument

echo "Making named pipes"

# mkfifo /tmp/in_D1
# mkfifo /tmp/in_D2
# mkfifo /tmp/in_D3

mkfifo /tmp_$1 /tmp/out_$1
 
echo "Starting N1"
./net -n $1 < /tmp/in_$1 >> /tmp/out_$1 &
sleep 1

echo "Everything running. (start a display and a player to begin)"