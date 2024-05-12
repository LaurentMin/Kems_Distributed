package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func handleUserInput(input string, playerIndex string) {
	logInfo("handleUserInput", "Handling input "+input)
	input = strings.ToLower(strings.TrimSpace(input))
	// Error empty input
	if len(input) == 0 {
		logError("Input terminal", "Input empty!")
		return
	}

	switch input[0] {
	case 's':
		numberStr := input[1:]
		if _, err := strconv.Atoi(numberStr); err == nil {
			drawPileIndex := int(input[1]) - '0' - 1
			playerCardIndex := int(input[2]) - '0' - 1
			sendAction("SwapCards", []string{"playerIndex", "playerCardIndex", "drawPileCardIndex"}, []string{name, strconv.Itoa(playerCardIndex), strconv.Itoa(drawPileIndex)})
		}
	case 'n':
		sendAction("NextTurn", []string{"playerIndex"}, []string{playerIndex})
	case 'k':
		if input == "kems" {
			sendAction("Kems", []string{"playerIndex"}, []string{playerIndex})
		}
	case 'c':
		if _, err := strconv.Atoi(string(input[1])); err == nil {
			otherPlayerIndex := int(input[1]) - '0'
			sendAction("ContreKems", []string{"playerIndex"}, []string{strconv.Itoa(otherPlayerIndex)})
		}

	default:
		logError("Input terminal", "Invalid input "+input)
	}
}

func sendAction(actionType string, actionParamsNames []string, actionParamsValues []string) {
	params := encodeMessage(actionParamsNames, actionParamsValues)
	action := encodeMessage([]string{"typ", "prm"}, []string{actionType, params})

	logInfo("sendAction", "Sending action "+action)
	fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{"P" + name, action}) + "\n")
}

func main() {
	// Getting name from commandline (needs to be the same as the app)
	pName := flag.String("n", "player", "name (1,2,3)")
	flag.Parse()
	name = *pName

	if name != "1" && name != "2" && name != "3" {
		logError("main", "Wrong playerId for anneauCtl structure, change code if needed.")
		return
	}

	// Starting Player
	sendAction("InitPlayer", []string{"newPlayer"}, []string{name})
	// logInfo("main", "Launching player...")
	playerInput := ""

	for {
		// Message reception
		playerInput = scanUntilNewline()
		handleUserInput(playerInput, name)

		playerInput = ""
	}
}
