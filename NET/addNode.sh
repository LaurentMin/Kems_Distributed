#!/bin/bash

# Create first NET (doesn't connect to anyone)
if [ $# -eq 1 ]; then
    echo "Making named pipes..."
    mkfifo /tmp/in_$1 /tmp/out_$1
    echo "Starting node..."
    ./net -n $1 < /tmp/in_$1 >> /tmp/out_$1 &
    echo "New network created."
# Add a NET to another node that already exists
elif [ $# -eq 2 ]; then
    echo "Making named pipes..."
    mkfifo /tmp/in_$1 /tmp/out_$1
    echo "Starting node..."
    ./net -n $1 -a $2 < /tmp/in_$1 >> /tmp/out_$1 &
    echo "Connecting named pipes..."
    cat /tmp/out_$2 | tee -a /tmp/in_$1
    cat /tmp/out_$1 | tee -a /tmp/in_$2
    echo "Node added to network."
else
    echo "Error: Invalid number of arguments"
    echo "Usage: $0 nodeName nodeToConnectTo (starting new network => 1 agument)"
    exit 1
fi