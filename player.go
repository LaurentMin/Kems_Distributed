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
	actionParamsNames := []string{"playerIndex", "playerCardIndex", "drawPileCardIndex"}
	actionParamsValues := []string{name, "0", "0"}
	params := encodeMessage(actionParamsNames, actionParamsValues)
	if len(actionParamsNames) != len(actionParamsValues) {
		logError("main", "Bad parameter setting, modify code.")
		return
	}

	// Main loop of the player, manages user input and sending messages
	for {
		logInfo("main", "Input action.")
		// Message reception
		playerInput = scanUntilNewline()
		logInfo("main", playerInput+" action received.")

		// Building action
		action := encodeMessage([]string{"typ", "prm"}, []string{playerInput, params})
		// Sending action
		logInfo("main", "Sending action "+action+" to "+inputLocation)
		// logInfo("main", "Sending action to "+inputLocation)
		fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{"P" + name, action}) + "\n")

		action = ""
		playerInput = ""
	}
}
