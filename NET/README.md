# Testing instructions

1. Once in NET/ (where you found this readme) execute `./build.sh`
2. Then `cd ..`
3. Then `./builScript.sh`
4. Go back in NET/ `cd NET/`
5. Start network (takes 30 seconds) `./fullyStartScript.sh`
6. In another terminal run terminal display and then players 1 to 3.

<br>

In this scenario, nodes 0, 1, 2 and 3 are created. Then after 30 seconds, node 0 removes itself. The player is still displayed and cards are given to him. However, it is technically disabled, it's NET controller is in `ZOMBIE` mode.