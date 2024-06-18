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



## Define function to reconnect nodes
#reconnect_nodes() {
#    local nodes=("$@")  # 接收传入的节点数组
#
#    # Reconnect nodes
##    local first_node=${nodes[0]}
#    for (( i = 0; i < ${#nodes[@]}; i++ )); do
#        local node=${nodes[i]}
#        number="${node:1:1}"
#        if [ number -n 0 ];then
#          cat /tmp/out_${node} | tee -a /tmp/in_${first_node} >> /tmp/in_C2 &
#        else
#          cat /tmp/out_${first_node} | tee -a /tmp/in_${node} &
#        fi
#    done
#
#    for (( i = 0; i < ${#nodes[@]}; i++ )); do
#            local node=${nodes[i]}
#            number="${node:1:1}"
#            if [ number -n 0 ];then
#              cat /tmp/out_${node} | tee -a /tmp/in_${first_node} >> /tmp/in_C2 &
#            else
#              cat /tmp/out_${first_node} | tee -a /tmp/in_${node} &
#
#            if
#            cat /tmp/out_A2 | tee -a /tmp/in_C2 >> /tmp/in_Debug &
#            cat /tmp/out_C2 | tee -a /tmp/in_A2 >> /tmp/in_${NODES[0]} &
#
#            cat /tmp/out_A3 | tee -a /tmp/in_C3 >> /tmp/in_Debug &
#            cat /tmp/out_C3 | tee -a /tmp/in_A3 >> /tmp/in_${NODES[1]} &
#        done
#}

NODE_TO_REMOVE=$1
CONNECTED_NODES=("${@:2}")


# Step 1: Send delete message to connected nodes
echo ";:snd:${NODE_TO_REMOVE};:rec:${NODE_TO_REMOVE};:typ:del;:msg:Hello, may I leave the network ?" > /tmp/in_${NODE_TO_REMOVE}
sleep 15 # Wait for the delete messages to propagate

## Step 2: Kill the node process
#kill_node_process $NODE_TO_REMOVE
#sleep 1
## Step 3: Remove named pipes
#remove_named_pipes $NODE_TO_REMOVE
#sleep 1

echo "fin"

## Step 4: Reconnect remaining nodes if any; pour tester que le
## Clear old cats and tees
#echo "Clearing old cats and tees..."
#killall tee 2> /dev/null
#killall cat 2> /dev/null

#./Reconnect.sh
