package main

import (
	"flag"
	"fmt"
)

func displayCard(card Card) string {
	suit := ""
	if card.Suit == "Hearts" {
		suit = "♥"
	} else if card.Suit == "Diamonds" {
		suit = "♦"
	} else if card.Suit == "Clubs" {
		suit = "♣"
	} else if card.Suit == "Spades" {
		suit = "♠"
	}

	if card.Suit == "Hearts" || card.Suit == "Diamonds" {
		return red + "|" + card.Value + " " + suit + "|" + reset + "\t"
	} else {
		return "|" + card.Value + " " + suit + "|\t"
	}
}

func displayDrawPile(game GameState) {
	drawCards := "Draw pile:\t"
	numbersCard := "          \t"
	for i := 0; i < len(game.DrawPile); i++ {
		drawCards += displayCard(game.DrawPile[i])
		numbersCard += "  " + fmt.Sprint(i+1) + "  \t"
	}
	stderr.Printf(drawCards + "\n" + numbersCard + "\n\n")
}

func displayDiscardPile(game GameState) {
	discardCards := "Discard pile:\t"
	numbersCard := "          \t"
	for i := 0; i < len(game.DiscardPile); i++ {
		discardCards += displayCard(game.DiscardPile[i])
		numbersCard += "  " + fmt.Sprint(i+1) + "  \t"
	}
	stderr.Printf(discardCards + "\n" + numbersCard + "\n\n")
}

func displayDeck(game GameState) {
	deckCards := "Deck pile:\t"
	numbersCard := "          \t"
	for i := 0; i < len(game.Deck); i++ {
		deckCards += displayCard(game.Deck[i])
		numbersCard += "  " + fmt.Sprint(i+1) + "  \t"
	}
	stderr.Printf(deckCards + "\n" + numbersCard + "\n\n")
}

func displayPlayerHand(player Player) {
	numbersCard := "          \t"
	cards := player.Name + "'s hand:\t"
	for i := 0; i < len(player.Hand); i++ {
		cards += displayCard(player.Hand[i])
		numbersCard += "  " + fmt.Sprint(i+1) + "\t"
	}
	stderr.Printf(cards + "\n" + numbersCard + "\n\n")
}

func displayGameBoard(game GameState) {
	displayDeck(game)
	displayDiscardPile(game)
	displayDrawPile(game)
	for i := 0; i < len(game.Players); i++ {
		displayPlayerHand(game.Players[i])
	}
}

/*
	This program displays the full state of the game (usefull for debugging and testing)
	Update is done when app that sent update receives its own state, setting the end of the update
*/
func main() {
	// Getting input file from commandline
	pfile := flag.String("f", "/tmp/in_Debug", "input file")
	flag.Parse()

	// Setting app name (usefull for debug)
	name = "Display"

	messageReceived := ""
	state := GameState{}
	inputFile := *pfile
	logInfo("main", "Displaying with input file "+inputFile)

	// Main loop, displays game state when receives it
	for {
		logInfo("main", "Waiting for next state...")
		// Message reception
		messageReceived = scanUntilNewline()

		// Ignore message not for display
		if len(messageReceived) < 11 || messageReceived[:11] != "[GAMESTATE]" {
			logInfo("main", "Message received not destinated to display.")
			state = GameState{}
			messageReceived = ""
			continue
		}

		state = stringToGameState(messageReceived)
		// Check if message is state
		if len(state.Players) == 0 {
			logInfo("main", "Wrong state received ignoring.")
		} else {
			logInfo("main", "State received, displaying.")
			displayGameBoard(state)
		}

		state = GameState{}
		messageReceived = ""
	}
}
