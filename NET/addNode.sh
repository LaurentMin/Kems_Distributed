#!/bin/bash

# Add a NET to another node that already exists
if [ $# -eq 2 ] && [ "$1" != "$2" ]; then
    echo "Making named pipes..."
    mkfifo /tmp/in_$1 /tmp/out_$1
    echo "Starting node..."
    ./net -n $1 -a $2 < /tmp/in_$1 >> /tmp/out_$1 &
    echo "Adding node to network..."
    cat /tmp/out_$2 | tee -a /tmp/in_$1 &
    cat /tmp/out_$1 | tee -a /tmp/in_$2 &
    echo "Node added to network."
else
    echo "Error: Invalid number of arguments (2 and must be different)"
    echo "Usage: $0 nodeName nodeToConnectTo"
    exit 1
fi