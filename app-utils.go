package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
	// logMessage("getInitSettings", "Initialising settings")
	settings := GameSettings{
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
	// logMessage("getInitDeck", "Initialising deck")
	deck := []Card{}

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
	// logMessage("getInitPlayers", "Initialising players")
	players := []Player{
		{Name: "Newbie", Score: 0, Hand: []Card{}},
		{Name: "Mexican", Score: 0, Hand: []Card{}},
		{Name: "Convict", Score: 0, Hand: []Card{}},
	}
	return players
}

/*
	Returns an initialised Game state (used when there is none)
	It has initialised settings, deck (shuffled) and players
	Draw pile and discard pile are empty as well as player hands
*/
func getInitState() GameState {
	// logMessage("getInitState", "Initialising game state")
	game := GameState{
		Settings:    getInitSettings(),
		Deck:        getInitDeck(),
		Players:     getInitPlayers(),
		DrawPile:    []Card{},
		DiscardPile: []Card{},
	}
	return game
}

////////////////////////
// ToSTRING FUNCTIONS //
////////////////////////
/*
	toString functions (for printing structures)
*/
func toStringCard(card Card) string {
	return card.Value + " " + card.Suit
}

func toStringPlayer(player Player) string {
	return player.Name
}

func toStringCards(cards []Card) string {
	cardsString := ""
	for i := 0; i < len(cards); i++ {
		cardsString += "," + toStringCard(cards[i])
	}
	return cardsString
}

///////////////////////
// CASTING FUNCTIONS //
///////////////////////
/*
	type1ToType2 functions (for changing the type of a data)
	Used for transforming struc into string before sending and opposit when received
*/
func gameStateToString(game GameState) string {
	// ERROR returns ""
	if game.Settings == (GameSettings{}) {
		logError("main", "Game state is empty (can be fatal for program).")
		return ""
	}

	// Prepare Players
	playersString := ""
	for i := 0; i < len(game.Players); i++ {
		playersString +=
			encodeMessage(
				[]string{"Player" + strconv.Itoa(i)},
				[]string{
					encodeMessage(
						[]string{"Name", "Score", "Hand"},
						[]string{game.Players[i].Name, strconv.Itoa(game.Players[i].Score), toStringCards(game.Players[i].Hand)})})
	}

	gameString := "[GAMESTATE]"
	// Encode all at once (or can generate problems)
	gameString += encodeMessage(
		[]string{
			"DrawPileSize",
			"HandSize",
			"Deck",
			"DrawPile",
			"DiscardPile",
			"numPlayers",
			"Players"},
		[]string{
			strconv.Itoa(game.Settings.DrawPileSize),
			strconv.Itoa(game.Settings.HandSize),
			toStringCards(game.Deck),
			toStringCards(game.DrawPile),
			toStringCards(game.DiscardPile),
			strconv.Itoa(len(game.Players)),
			playersString})

	return gameString
}

// Helper func to retrieve cards for Deck, DrawPile, DiscardPile and player Hand
func getCardsFromString(cardsString string) []Card {
	cards := []Card{}
	// Table contains all cards as "value suit"
	cardsTab := decodeMessage(cardsString)
	for i := 0; i < len(cardsTab); i++ {
		// cardsTab[i] format "value suit" splitted in {value,suit}
		card := strings.Split(cardsTab[i], " ")
		// Add card to cards list
		cards = append(cards, Card{card[0], card[1]})
	}

	return cards
}
func stringToGameState(gameString string) GameState {
	// ERROR returns empty game state
	if gameString[:11] != "[GAMESTATE]" {
		logError("main", "String is not a game state "+gameString+" (can be fatal for program).")
		return GameState{}
	}
	// Remove game state header
	gameString = gameString[11:]
	// Decode game state string
	tabString := decodeMessage(gameString)
	game := GameState{}
	//game.Players = []Player{}

	// Retrieve Settings
	game.Settings.DrawPileSize, _ = strconv.Atoi(findValue(tabString, "DrawPileSize"))
	game.Settings.HandSize, _ = strconv.Atoi(findValue(tabString, "HandSize"))
	// Retrieve Deck
	game.Deck = getCardsFromString(findValue(tabString, "Deck"))
	// Retrieve DrawPile
	game.DrawPile = getCardsFromString(findValue(tabString, "DrawPile"))
	// Retrieve DiscardPile
	game.DiscardPile = getCardsFromString(findValue(tabString, "DiscardPile"))
	// Retrieve number of players (used to decode players)
	numPlayers, _ := strconv.Atoi(findValue(tabString, "numPlayers"))
	// Retrieve Players
	players := findValue(tabString, "Players")
	playersTab := decodeMessage(players)
	for i := 0; i < numPlayers; i++ {
		// Gett player string
		player := findValue(playersTab, "Player"+strconv.Itoa(i))
		// Split the player in a tab for each value
		playerTab := decodeMessage(player)
		// Set each field of the player
		score, _ := strconv.Atoi(findValue(playerTab, "Score"))
		game.Players = append(game.Players, Player{findValue(playerTab, "Name"), score, getCardsFromString(findValue(playerTab, "Hand"))})
	}

	return game
}

//////////////////////
// LOOKUP FUNCTIONS //
//////////////////////
/*
	contains functions (to find a given element in an array)
*/
func contains(card Card, cards []Card) bool {
	for i := 0; i < len(cards); i++ {
		if card.Suit == cards[i].Suit && card.Value == cards[i].Value {
			return true
		}
	}
	return false
}

/*
	findIndex functions (to find the index of a given element in an array)
	Returns -1 if none found => dangerous to use directly in []
*/
// Finds a player index by name only
func findIndexPlayer(player Player, players []Player) int {
	for i := 0; i < len(players); i++ {
		if player.Name == players[i].Name {
			return i
		}
	}
	return -1
}

// Finds a card index
func findIndexCard(card Card, cards []Card) int {
	for i := 0; i < len(cards); i++ {
		if card.Suit == cards[i].Suit && card.Value == cards[i].Value {
			return i
		}
	}
	return -1
}

////////////////////////////////
// GAMESTATE HELPER FUNCTIONS //
////////////////////////////////
/*
	Checks if given player has won
*/
func hasKems(game GameState, playerIndex int) bool {
	// Check params, Errors return false
	player := game.Players[playerIndex]
	handSize := game.Settings.HandSize
	if len(player.Hand) < handSize {
		logError("hasKems", "Player has insufficient cards in hand. (FATAL ERROR)")
		return false
	}
	// Getting first card
	value := player.Hand[0].Value
	suit := player.Hand[0].Suit
	if value == "" || suit == "" {
		logError("hasKems", "Player cards are undefined. (FATAL ERROR)")
		return false
	}

	// Returns false if any card is different from first one
	for i := 1; i < len(player.Hand); i++ {
		if value != player.Hand[i].Value {
			return false
		}
	}

	// returns true if all similar
	return true
}

func sendGameStateToPLayer(game GameState) {
	fmt.Printf(gameStateToString(game) + "\n")
}
