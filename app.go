package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

////////////////////////////
// GAME STATE DECLARATION //
////////////////////////////
/*
	Single Card structure
*/
type Card struct {
	Value string
	Suit  string
}

/*
	Single Player structure
*/
type Player struct {
	Name  string
	Score int
	Hand  []Card
}

/*
	Game parameters structure
	Allows to define the rules for the game
*/
type GameSettings struct {
	HandSize     int // Number of cards a player can hold
	DrawPileSize int // Number of cards in the draw pile
}

/*
	Global game state structure (shared data)
	Contains deck, players (and their hands, score), draw pile
*/
type GameState struct {
	Settings    GameSettings
	Deck        []Card
	Players     []Player
	DrawPile    []Card
	DiscardPile []Card
}

///////////////////////////////
// GAME STATE INITIALISATION //
///////////////////////////////
/*
	Returns initialised Settings
*/
func getInitSettings() GameSettings {
	logMessage("getInitSettings", "Initialising settings")
	settings := GameSettings{
		HandSize:     2,
		DrawPileSize: 2,
	}
	return settings
}

/*
	Returns an initialised Deck
	Builds the deck and shuffles it
*/
func getInitDeck() []Card {
	logMessage("getInitDeck", "Initialising deck")
	deck := []Card{}

	// Setting the deck building parameters
	//values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	//suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10"}
	suits := []string{"Clubs", "Diamonds"}

	// Building the deck
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Value: value, Suit: suit})
		}
	}

	// Shuffling the deck
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

/*
	Returns an initialised list of Players
*/
func getInitPlayers() []Player {
	logMessage("getInitPlayers", "Initialising players")
	players := []Player{
		{Name: "Newbie", Score: 0, Hand: []Card{}},
		{Name: "Mexican", Score: 0, Hand: []Card{}},
		{Name: "Convict", Score: 0, Hand: []Card{}},
		{Name: "Tomatoe", Score: 0, Hand: []Card{}},
	}
	return players
}

/*
	Returns an initialised Game state (used when there is none)
	It has initialised settings, deck (shuffled) and players
	Draw pile and discard pile are empty as well as player hands
*/
func getInitState() GameState {
	logMessage("getInitState", "Initialising game state")
	game := GameState{
		Settings:    getInitSettings(),
		Deck:        getInitDeck(),
		Players:     getInitPlayers(),
		DrawPile:    []Card{},
		DiscardPile: []Card{},
	}
	return game
}

////////////////
// GAME UTILS //
////////////////
/*
	Returns a gamestate with a new deck and empty discard pile
	Puts the cards of the discard pile back into the deck and shuffles the deck
*/
func reshuffleDiscard(game GameState) GameState {
	logMessage("reshuffleDiscard", "Putting cards back from discard pile to deck and shuffling.")

	// Putting cards of the discard pile back in the deck
	for i := 0; i < len(game.DiscardPile); i++ {
		// Putting discard pile cards back in the deck
		game.Deck = append(game.Deck, game.DiscardPile[i])
	}

	// Emptying discard pile
	game.DiscardPile = []Card{}

	// Shuffling the deck
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(game.Deck), func(i, j int) {
		game.Deck[i], game.Deck[j] = game.Deck[j], game.Deck[i]
	})

	return game
}

/*
	Returns a game state with a new draw pile
	Works both with a renewing and a first draw of the draw pile
*/
func renewDrawPile(game GameState) GameState {
	// Bool allows to differentiate the first draw from the renewing of the pile
	drawPileEmpty := len(game.DrawPile) == 0
	if drawPileEmpty {
		logMessage("renewDrawPile", "Drawing cards from the deck to fill the draw pile (first draw).")
	} else {
		logMessage("renewDrawPile", "Drawing cards from the deck to renew the draw pile.")
	}

	// Drawing cards from the deck until draw pile is filled
	for i := 0; i < game.Settings.DrawPileSize; i++ {
		// Reshuffling deck if it is empty
		if len(game.Deck) == 0 {
			game = reshuffleDiscard(game)
		}

		if drawPileEmpty == false {
			// Renewing drawpile
			// Adding old card to the discard pile
			game.DiscardPile = append(game.DiscardPile, game.DrawPile[i])
			// Replacing old card by new one in the draw pile
			game.DrawPile[i] = game.Deck[len(game.Deck)-1]
		} else {
			// Setting empty draw pile
			// Adding new card to the draw pile
			game.DrawPile = append(game.DrawPile, game.Deck[len(game.Deck)-1])
		}
		// Removing added card from the deck
		game.Deck = game.Deck[:len(game.Deck)-1]
	}

	return game
}

/*
	Returns a game state with a new player hands
	Works both with a renewing and a first draw of the player hands
*/
func renewPlayerHands(game GameState) GameState {
	// Bool allows to differentiate the first draw from the renewing of the pile
	playerHandsEmpty := len(game.Players[0].Hand) == 0
	if playerHandsEmpty {
		logMessage("renewPlayerHands", "Drawing cards from the deck to fill the players hands (first draw).")
	} else {
		logMessage("renewPlayerHands", "Drawing cards from the deck to renew players hands.")
	}

	// Error returns same game state
	if len(game.Players) == 0 {
		logError("renewPlayerHands", "Game does not have players.")
		return game
	}

	// Drawing cards from the deck until player hands are filled
	for i := 0; i < game.Settings.HandSize; i++ {
		for playerIndex := 0; playerIndex < len(game.Players); playerIndex++ {
			// Reshuffling deck if it is empty
			if len(game.Deck) == 0 {
				game = reshuffleDiscard(game)
			}

			if playerHandsEmpty == false {
				// Renewing hands
				// Adding old card to the discard pile
				game.DiscardPile = append(game.DiscardPile, game.Players[playerIndex].Hand[i])
				// Replacing old card by new one in the hand
				game.Players[playerIndex].Hand[i] = game.Deck[len(game.Deck)-1]
			} else {
				// Setting empty hands
				// Adding new card to the hand
				game.Players[playerIndex].Hand = append(game.Players[playerIndex].Hand, game.Deck[len(game.Deck)-1])
			}
			// Removing added card from the deck
			game.Deck = game.Deck[:len(game.Deck)-1]
		}
	}

	return game
}

/*
	Returns a game state after a player exchanged a card with the draw pile
*/
func swapCard(playerCard Card, drawPileCard Card, player Player, game GameState) GameState {
	logMessage("swapCard", "Swapping player card "+toStringCard(playerCard)+" with draw pile card "+toStringCard(drawPileCard))
	// Getting usefull variables and checking params
	indexPlayer := findIndexPlayer(player, game.Players)
	if indexPlayer == -1 {
		logError("swapCard", "Player "+toStringPlayer(player)+" does not exist.")
		return game
	}

	playerHand := game.Players[indexPlayer].Hand
	if contains(playerCard, playerHand) == false {
		logError("swapCard", "Player "+toStringPlayer(player)+" does not have card "+toStringCard(playerCard))
		return game
	}

	if contains(drawPileCard, game.DrawPile) == false {
		logError("swapCard", "Can't find card "+toStringCard(drawPileCard)+" in draw pile.")
		return game
	}

	// The player hand and the drawpile both have the corresponding cards
	// Changing players card
	game.Players[indexPlayer].Hand[findIndexCard(playerCard, playerHand)] = drawPileCard
	// Changing draw pile card
	game.DrawPile[findIndexCard(drawPileCard, game.DrawPile)] = playerCard

	logSuccess("swapCard", "Cards exchanged "+toStringCard(drawPileCard)+" with "+toStringCard(playerCard))
	return game
}

//////////
// GAME //
//////////
func main() {
	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "app", "name")
	flag.Parse()
	name = *pName

	// Starting App
	logInfo("main", "Launching app...")
	// Initialising key variables for app
	messageReceived := ""
	game := getInitState()
	game = renewDrawPile(game)
	game = renewPlayerHands(game)

	// Main loop of the app, manages message reception and emission and processing
	for {
		if name == "A1" {
			// One of the apps plays the game (for testing)
			logInfo("main", "A1 swapping cards.")
			game = swapCard(game.Players[0].Hand[0], game.DrawPile[0], game.Players[0], game)
			// game = renewDrawPile(game)
			// game = renewPlayerHands(game)
			fmt.Printf(gameStateToString(game) + "\n")
			time.Sleep(time.Duration(10) * time.Second)
		} else {
			// Standard app behaviour
			logInfo("main", "Waiting for message.")
			// Message reception
			messageReceived = scanUntilNewline()
			messageReceived = messageReceived[:len(messageReceived)-1]
			logInfo("main", "Message received. "+messageReceived)

			// Message is not a game state (ignore)
			if len(messageReceived) <= 11 || messageReceived[:11] != "[GAMESTATE]" {
				logInfo("main", "Wrong message type for app received "+messageReceived+" (ignoring).")
				messageReceived = ""
				continue
			}

			// Message is a game state (process)
			logInfo("main", "Processing game state... "+messageReceived)
			// Replace game state if an update was received
			if gameStateToString(game) != messageReceived {
				game = stringToGameState(messageReceived)
				// Sending update to next app (through controller)
				fmt.Printf(messageReceived + "\n")
				logInfo("main", "Sent updated game state to next app through controller.")
			} else {
				logSuccess("main", "Game state is already up to date, all apps up to date.")
			}

			messageReceived = ""
		}
	}
}
