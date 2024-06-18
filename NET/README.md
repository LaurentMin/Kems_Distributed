# Testing instructions

# Using NET To delete
1. Once in NET/ (where you found this readme) execute `./build.sh`
2. Then `cd ..`
3. Then `./builScript.sh`
4. Go back in NET/ `cd NET/`
5. Start network (takes 30 seconds) `./fullyStartScript.sh`
6. Then `./Delete.sh N2(noeud for delete) N1(noeud for reconnection)` N2 N1 is just an exemple.

# Using NET To add 

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
4. starts N4 on N1
5. starts N6 on N4 AND N5 on N3

**Use `fullStartScript.sh` (Starts every script with an interval of 10 seconds).**
