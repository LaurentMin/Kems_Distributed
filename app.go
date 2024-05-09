package main

import (
	"fmt"
	"math/rand"
	"time"
)

///////////////////////////
// GAMESTATE DECLARATION //
///////////////////////////
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
	NumPlayers   int // Number of players
	HandSize     int // Number of cards a player can hold
	DrawPileSize int // Number of cards in the draw pile
}

/*
	Global game state structure (shared data)
	Contains deck, players (and their hands, score), draw pile
*/
type GameState struct {
	Settings GameSettings
	Deck     []Card
	Players  []Player
	DrawPile []Card
}

//////////////////////////////
// GAMESTATE INITIALISATION //
//////////////////////////////
/*
	Returns initialised Settings
*/
func getInitSettings() GameSettings {
	logMessage("getInitSettings", "Initialising settings")
	settings := GameSettings{
		NumPlayers:   4,
		HandSize:     4,
		DrawPileSize: 4,
	}
	return settings
}

/*
	Returns an initialised Deck
	Builds the deck and shuffles it
*/
func getInitDeck() []Card {
	logMessage("getInitDeck", "Initialising deck")
	var deck []Card

	// Setting the deck building parameters
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}

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
*/
func getInitState() GameState {
	logMessage("getInitState", "Initialising game state")
	game := GameState{
		Settings: getInitSettings(),
		Deck:     getInitDeck(),
		Players:  getInitPlayers(),
		DrawPile: []Card{},
	}
	return game
}

////////////////
// GAME UTILS //
////////////////
/*
	Returns a game state with a new draw pile
*/
func renewDrawPile(game GameState) GameState {
	logMessage("renewDrawPile", "Drawing cards from the deck to renew the draw pile.")
	game.DrawPile = []Card{}
	// Drawing cards from the deck until draw pile is filled
	for i:=0; i<game.Settings.DrawPileSize; i++ {
		// Adding card to the pile
		game.DrawPile = append(game.DrawPile, game.Deck[len(game.Deck)-1])
		// Removing added card from the deck
		game.Deck = game.Deck[:len(game.Deck)-1]
	}
	
	return game
}

/////////////////////////
// GAME INITIALISATION //
/////////////////////////
func app() {
	// APP TESTS
	gameState := getInitState()
	fmt.Println(gameState)
	gameState = renewDrawPile(gameState)
	fmt.Println(gameState)
	 
}
