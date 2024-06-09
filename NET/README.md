# Using NET

1. `./build.sh` after editing
2. `./clearScript.sh` to stop execution
3. `./startNetwork.sh` to start a network (first node is N0)
4. `./addNode.sh newNode nodeToConnectTo` to add a node

The script addNode does not work because it deletes old cat...

So scripts add scripts built `addn.sh` :
0. `./startNetwork.sh` Starts N0
1. `add1.sh` Starts N1 on N0
2. starts N2 on N1
3. starts N3 on N0
