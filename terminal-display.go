package main

import (
	"flag"
	"fmt"
	"time"
)

func clearScreen() {
	stderr.Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
}

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

func displayScore(game GameState) {
	scoreString := "Actual score:\t "
	scoreMax := 0
	var winner string
	for i := 0; i < len(game.Players); i++ {
		scoreString += game.Players[i].Name + ": " + fmt.Sprint(game.Players[i].Score) + "\t"
		if game.Players[i].Score > scoreMax {
			scoreMax = game.Players[i].Score
			winner = game.Players[i].Name
		}
	}
	stderr.Printf(scoreString + "\n The current winner is " + winner + "\n\n")
}

func displayKemsRules() {
	fmt.Println("KEMS Card Game Rules:")
	fmt.Println("KEMS is a card game played with a standard deck of 52 cards.")
	fmt.Println("The objective is to be the first player to form a hand with the same value and shout 'KEMS' to win.")
	fmt.Println()
	fmt.Println("Rules:")
	fmt.Println("1. Players are dealt 4 cards at the beginning.")
	fmt.Println("2. Each turn there are 4 cards in the draw pile.")
	fmt.Println("3. Any player can exchange a card from their hand with a card from the draw pile at any moment.")
	fmt.Println("4. When nobody is exchanging cards, a player replace the draw pile.")
	fmt.Println("5. The first player to form a combination of cards of the exact same value quickly needs to type 'KEMS' to win 1 point.")
	fmt.Println("6. When one player has a 'KEMS' combination, the other players can type 'C' followed by the player index BEFORE the winner type 'KEMS' to counter the 'KEMS', therefore giving a malus of 1 point to the winner.")
}

func displayCommands() {
	fmt.Println("Commands (maj and spaces are not necessary):")
	fmt.Println("Swap cards : s <drawPileCardIndex> <playerCardIndex>")
	fmt.Println("Next turn : n")
	fmt.Println("Kems: kems")
	fmt.Println("ContreKems: c <playerIndex>\n\n\n")
}

func displayGameBoard(game GameState) {
	clearScreen()
	displayCommands()
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

	messageReceived := ""
	state := GameState{}
	inputFile := *pfile
	logInfo("main", "Displaying with input file "+inputFile)

	displayKemsRules()
	time.Sleep(5 * time.Second)
	for {
		logInfo("main", "Waiting for next state...")
		messageReceived = scanUntilNewline()

		// Ignore message not for display
		if len(messageReceived) < 11 || messageReceived[:11] != "[GAMESTATE]" {
			logInfo("main", "Message received not destinated to display.")
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
