package main

import (
	"flag"
	"fmt"
)

func main() {
	// Getting name from commandline (needs to ba the same as the app)
	pName := flag.String("n", "player", "name (1,2,3)")
	flag.Parse()
	name = *pName
	// Cheking name
	if name != "1" && name != "2" && name != "3" {
		logError("main", "Wrong name for anneauCtl structure, change code if needed.")
		return
	}

	// Starting Player
	logInfo("main", "Launching player...")
	// Initialising key variables for player
	inputLocation := "/tmp/in_A" + name
	playerInput := ""

	// Main loop of the player, manages user input and sending messages
	for {
		logInfo("main", "Input action.")
		// Message reception
		playerInput = scanUntilNewline()
		logInfo("main", playerInput+" action received.")

		logInfo("main", "Sending action to "+inputLocation)
		fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{"P" + name, playerInput}) + "\n")

		playerInput = ""
	}
}
