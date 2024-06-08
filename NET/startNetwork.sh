#!/bin/bash

echo "Making named pipes..."
mkfifo /tmp/in_admin /tmp/out_admin
mkfifo /tmp/in_N0 /tmp/out_N0
echo "Starting first node..."
./netadmin < /tmp/in_admin >> /tmp/out_admin &
sleep 1
./net -n N0 -a N0 < /tmp/in_N0 >> /tmp/out_N0 &

echo "Adding node to network..."
cat /tmp/out_admin | tee -a /tmp/in_N0 &
cat /tmp/out_N0 | tee -a /tmp/in_admin &
echo "Network created."