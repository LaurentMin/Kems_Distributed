package main

import (
	"flag"
	"fmt"
	"strconv"
)

func handleUserInput(input string, playerIndex int) {
	logInfo("Player terminal", "Handling input "+input)
	input = strings.ToLower(strings.TrimSpace(input))	
	switch input[0] {
		case 's':
			drawPileIndex := strconv.Atoi(input[1])
			playerCardIndex := strconv.Atoi(input[2])
			sendAction("SwapCards", []string{"playerIndex", "playerCardIndex", "drawPileCardIndex"}, []string{name, playerCardIndex, drawPileIndex})
		case 'n':
			sendAction("NextTurn", []string{"playerIndex"}, []string{playerIndex})
		case 'k':
			if input == "kems" {
				sendAction("Kems", []string{"playerIndex"}, []string{playerIndex})
			} 
		case 'c':
			//test is input[1] is a number
			if _, err := strconv.Atoi(input[1]); err == nil {
				sendAction("ContreKems", []string{"playerIndex"}, []string{input[1]})
			
		default:
			logError("Player terminal", "Invalid input "+input)
		}
}

func sendAction(actionType string, actionParamsNames []string, actionParamsValues []string) {
	params := encodeMessage(actionParamsNames, actionParamsValues)
	action := encodeMessage([]string{"typ", "prm"}, []string{actionType, params})

	logInfo("main", "Sending action "+action)
	fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{"P" + name, action}) + "\n")
}

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
		handleUserInput(playerInput, name)

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
