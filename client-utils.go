package main

import "fmt"

func sendAction(actionType string, actionParamsNames []string, actionParamsValues []string) {
	params := encodeMessage(actionParamsNames, actionParamsValues)
	action := encodeMessage([]string{"typ", "prm"}, []string{actionType, params})

	logInfo("sendAction", "Sending action "+action)
	fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{"P" + name, action}) + "\n")
}

func checkIfKems(game GameState) int {
	for i := 0; i < len(game.Players); i++ {
		if hasKems(game, i) {
			return i
		}
	}
	return -1
}

func checkIfWinner(oldGame GameState, newGame GameState) int {
	if len(oldGame.Players) == 0 {
		return -1
	}
	for i := 0; i < len(newGame.Players); i++ {
		if newGame.Players[i].Score > oldGame.Players[i].Score {
			return i
		}
	}
	return -1
}

func checkIfLoser(oldGame GameState, newGame GameState) int {
	if len(oldGame.Players) == 0 {
		return -1
	}
	for i := 0; i < len(newGame.Players); i++ {
		if newGame.Players[i].Score < oldGame.Players[i].Score {
			return i
		}
	}
	return -1
}
