package main

import (
	"strconv"
	"strings"
)

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
