#!/bin/bash

# Function to kill the node process
kill_node_process() {
    NODE_NAME=$1
    echo "Killing node process for $NODE_NAME..."
    NODE_PID=$(ps aux | grep "./net -n $NODE_NAME" | grep -v grep | awk '{print $2}')
    if [ -n "$NODE_PID" ]; then
        kill $NODE_PID
        echo "Node process $NODE_PID killed."
    else
        echo "Node process not found."
    fi
}

# Function to remove named pipes
remove_named_pipes() {
    NODE_NAME=$1
    echo "Removing named pipes for $NODE_NAME..."
    rm -f /tmp/in_$NODE_NAME /tmp/out_$NODE_NAME
    echo "Named pipes removed."
}

# Function to reconnect remaining nodes
reconnect_nodes() {
    NODE_TO_REMOVE=$1
    CONNECTED_NODES=("${@:2}")

    echo "Reconnecting nodes..."

    # Reconnect all connected nodes among themselves
    for ((i = 0; i < ${#CONNECTED_NODES[@]}; i++)); do
        for ((j = i + 1; j < ${#CONNECTED_NODES[@]}; j++)); do
            NODE_A=${CONNECTED_NODES[$i]}
            NODE_B=${CONNECTED_NODES[$j]}
            echo "Connecting $NODE_A and $NODE_B"
            cat /tmp/out_$NODE_A | tee -a /tmp/in_$NODE_B &
            cat /tmp/out_$NODE_B | tee -a /tmp/in_$NODE_A &
        done
    done

    echo "Nodes reconnected."
}

# Check if at least one argument is provided
if [ $# -lt 1 ]; then
    echo "Error: Invalid number of arguments"
    echo "Usage: $0 nodeName [connectedNode1 connectedNode2 ...]"
    exit 1
fi

NODE_TO_REMOVE=$1
CONNECTED_NODES=("${@:2}")

# Step 1: Kill the node process
kill_node_process $NODE_TO_REMOVE

# Step 2: Remove named pipes
remove_named_pipes $NODE_TO_REMOVE

# Step 3: Reconnect remaining nodes if any
if [ ${#CONNECTED_NODES[@]} -gt 1 ]; then
    reconnect_nodes $NODE_TO_REMOVE "${CONNECTED_NODES[@]}"
elif [ ${#CONNECTED_NODES[@]} -eq 1 ]; then
    # If only one connected node, just remove the node without reconnecting
    echo "Node $NODE_TO_REMOVE was a leaf node. No need to reconnect."
else
    echo "Node $NODE_TO_REMOVE had no connected nodes."
fi

echo "Node $NODE_TO_REMOVE removed."
