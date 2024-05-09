package main

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
	for i := 0; i < len(cardsString); i++ {
		cardsString += toStringCard(cards[i])
	}
	return cardsString
}

/*
	type1ToType2 functions (for changing the type of a data)
	Used for transforming struc into string before sending and oppisit when received
*/
func gameStateToString(game GameState) string {
	gameString := ""

	return gameString
}

func stringToGameState(gameString string) GameState {
	game := GameState{}

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
