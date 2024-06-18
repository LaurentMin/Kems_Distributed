package main

import (
	"flag"
	"fmt"
	"strconv"
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
	drawCards := green + "Draw pile:" + reset + "\t"
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
	fmt.Println("Save game : i")
	fmt.Println("Swap cards : s <drawPileCardIndex> <playerCardIndex>")
	fmt.Println("Next turn : n")
	fmt.Println("Kems: kems")
	fmt.Println("ContreKems: c <playerIndex>\n\n\n")
}

func displayGameBoard(game GameState) {
	clearScreen()
	displayCommands()
	// displayDeck(game)
	fmt.Printf("Deck has %d cards left\n\n", len(game.Deck))
	// displayDiscardPile(game)
	displayDrawPile(game)
	fmt.Println()
	for i := 0; i < len(game.Players); i++ {
		if len(game.Players[i].Hand) == 0 {
			fmt.Println(game.Players[i].Name + " is disabled \n\n")
			continue
		}
		displayPlayerHand(game.Players[i])
	}
}

func displayWarningKems(playerIndex int) {
	fmt.Print(orange + "Player " + strconv.Itoa(playerIndex) + " has KEMS !!\n" + reset)
}

/*
	This program displays the full state of the game (usefull for debugging and testing)
	Update is done when app that sent update receives its own state, setting the end of the update
*/
func main() {
	// Getting input file from commandline
	pfile := flag.String("f", "/tmp/in_Debug", "input file")
	flag.Parse()
	name = "display"

	// Initialising important variables
	messageReceived := ""
	state := GameState{}
	inputFile := *pfile
	logInfo("main", "Displaying with input file "+inputFile)
	// Go routines to read input
	inChan := make(chan string, 10)
	go read(inChan)

	displayKemsRules()
	// time.Sleep(5 * time.Second)
	for {
		logInfo("main", "Waiting for next state...")
		messageReceived = <-inChan
		logInfo("main", "Message received: "+messageReceived)

		// Ignore message not for display
		if len(messageReceived) < 11 || messageReceived[:11] != "[GAMESTATE]" {
			logWarning("main", "Message received not destinated to display.")
			messageReceived = ""
			continue
		}

		newState := stringToGameState(messageReceived)
		// Check if message is state
		if len(newState.Players) == 0 {
			logInfo("main", "Wrong state received ignoring.")
		} else {
			winner := checkIfWinner(state, newState)
			if winner != -1 {
				clearScreen()
				fmt.Println("KEMS !!\n")
				fmt.Println("Player " + strconv.Itoa(winner) + " has won a point!\n")
				displayPlayerHand(state.Players[winner])
				fmt.Println()
				displayScore(newState)
				time.Sleep(10 * time.Second) // Pas ouf ça faudrait peut-être une variable pour être sûr que tt le monde est prêt
				displayGameBoard(newState)

				state = newState
				messageReceived = ""
				continue
			}
			loser := checkIfLoser(state, newState)
			if loser != -1 {
				clearScreen()
				fmt.Println("KEMS for the player " + strconv.Itoa(loser) + " ... Almost! Counter KEMS!\n")
				fmt.Println("Player " + strconv.Itoa(loser) + " has lost a point!\n")
				displayPlayerHand(state.Players[loser])
				fmt.Println()
				displayScore(newState)
				time.Sleep(10 * time.Second) // Pas ouf ça faudrait peut-être une variable pour être sûr que tt le monde est prêt
				displayGameBoard(newState)

				state = newState
				messageReceived = ""
				continue
			}
			potentialWinner := checkIfKems(newState)
			logInfo("main", "State received, displaying.")
			displayGameBoard(newState)
			if potentialWinner != -1 {
				displayWarningKems(potentialWinner)
			}
		}

		state = newState
		messageReceived = ""
	}
}
