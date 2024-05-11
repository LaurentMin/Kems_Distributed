package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func handleUserInput(input string, playerIndex string) {
	logInfo("Input terminal", "Handling input "+input)
	input = strings.ToLower(strings.TrimSpace(input))
	switch input[0] {
	case 's':
		numberStr := input[1:]
		if _, err := strconv.Atoi(numberStr); err == nil {
			drawPileIndex := int(input[1]) - '0' - 1
			playerCardIndex := int(input[2]) - '0' - 1
			sendAction("SwapCards", []string{"playerIndex", "playerCardIndex", "drawPileCardIndex"}, []string{playerId, strconv.Itoa(playerCardIndex), strconv.Itoa(drawPileIndex)})
		}
	case 'n':
		sendAction("NextTurn", []string{"playerIndex"}, []string{playerIndex})
	case 'k':
		if input == "kems" {
			sendAction("Kems", []string{"playerIndex"}, []string{playerIndex})
		}
	case 'c':
		if _, err := strconv.Atoi(string(input[1])); err == nil {
			sendAction("ContreKems", []string{"playerIndex"}, []string{string(input[1])})
		}

	default:
		logError("Input terminal", "Invalid input "+input)
	}
}

func sendAction(actionType string, actionParamsNames []string, actionParamsValues []string) {
	params := encodeMessage(actionParamsNames, actionParamsValues)
	action := encodeMessage([]string{"typ", "prm"}, []string{actionType, params})

	logInfo("Input terminal", "Sending action "+action)
	fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{"P" + playerId, action}) + "\n")
}

func main() {
	// Getting name from commandline (needs to be the same as the app)
	pName := flag.String("n", "player", "name (1,2,3)")
	flag.Parse()
	playerId = *pName

	if playerId != "1" && playerId != "2" && playerId != "3" {
		logError("Input terminal", "Wrong playerId for anneauCtl structure, change code if needed.")
		return
	}

	// Starting Player
	sendAction("InitPlayer", []string{}, []string{})
	logInfo("Input terminal", "Launching player...")
	playerInput := ""

	for {
		// Message reception
		playerInput = scanUntilNewline()
		handleUserInput(playerInput, playerId)

		playerInput = ""
	}
}
