#!/bin/bash

# Function to kill the node process
kill_node_process() {
    NODE_NAME=$1
    NODE_PID=$(ps aux | grep "./net -n $NODE_NAME" | grep -v grep | awk '{print $2}')
    if [ -n "$NODE_PID" ]; then
        kill $NODE_PID
    fi
}

# Function to remove named pipes
remove_named_pipes() {
    NODE_NAME=$1
    rm -f /tmp/in_$NODE_NAME /tmp/out_$NODE_NAME
}

# Function to reconnect remaining nodes
reconnect_nodes() {
    NODE_TO_REMOVE=$1
    CONNECTED_NODES=("${@:2}")

    # Reconnect all connected nodes among themselves
    for ((i = 0; i < ${#CONNECTED_NODES[@]}; i++)); do
        for ((j = i + 1; j < ${#CONNECTED_NODES[@]}; j++)); do
            NODE_A=${CONNECTED_NODES[$i]}
            NODE_B=${CONNECTED_NODES[$j]}
            cat /tmp/out_$NODE_A | tee -a /tmp/in_$NODE_B &
            cat /tmp/out_$NODE_B | tee -a /tmp/in_$NODE_A &
        done
    done
}

# Function to send delete message
send_delete_message() {
    NODE_NAME=$1
    CONNECTED_NODE=$2
    echo "del" > /tmp/in_$CONNECTED_NODE
}

## Check if at least one argument is provided
#if [ $# -lt 1 ]; then
#    echo "Error: Invalid number of arguments"
#    echo "Usage: $0 nodeName [connectedNode1 connectedNode2 ...]"
#    exit 1
#fi

NODE_TO_REMOVE=$1
CONNECTED_NODES=("${@:2}")

# Step 1: Send delete message to connected nodes
if [ ${#CONNECTED_NODES[@]} -gt 0 ]; then
    for NODE in "${CONNECTED_NODES[@]}"; do
        send_delete_message $NODE_TO_REMOVE $NODE
    done
    sleep 5 # Wait for the delete messages to propagate
fi

# Step 2: Kill the node process
kill_node_process $NODE_TO_REMOVE

# Step 3: Remove named pipes
remove_named_pipes $NODE_TO_REMOVE

# Step 4: Reconnect remaining nodes if any
if [ ${#CONNECTED_NODES[@]} -gt 1 ]; then
    reconnect_nodes $NODE_TO_REMOVE "${CONNECTED_NODES[@]}"
fi

# Call net program with leave flag and connected nodes
./net -d $NODE_TO_REMOVE -f "${CONNECTED_NODES[@]}"
