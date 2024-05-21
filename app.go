package main

import (
	"flag"
	"math/rand"
	"strconv"
	"time"
)

//#region ACT HANDLING UTILS
///////////////////////////
// ACTION HANDLING UTILS //
///////////////////////////
/*
	Returns a gamestate with a new deck and empty discard pile
	Puts the cards of the discard pile back into the deck and shuffles the deck
*/
func reshuffleDiscard(game GameState) GameState {
	// logMessage("reshuffleDiscard", "Putting cards back from discard pile to deck and shuffling.")

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
		// logMessage("renewDrawPile", "Drawing cards from the deck to fill the draw pile (first draw).")
	} else {
		// logMessage("renewDrawPile", "Drawing cards from the deck to renew the draw pile.")
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
		// logMessage("renewPlayerHands", "Drawing cards from the deck to fill the players hands (first draw).")
	} else {
		// logMessage("renewPlayerHands", "Drawing cards from the deck to renew players hands.")
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
	// logMessage("swapCard", "Swapping player card "+toStringCard(playerCard)+" with draw pile card "+toStringCard(drawPileCard))
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

	// logSuccess("swapCard", "Cards exchanged "+toStringCard(drawPileCard)+" with "+toStringCard(playerCard))
	return game
}

//#region ACT HANDLING
/////////////////////
// ACTION HANDLING //
/////////////////////
/*
	Global variable to know which player controls this app instance
*/
var lastConnectedPlayer string = ""

/*
	Handle player action
*/
func handleAction(fullAction string, game GameState) GameState {
	// Get action type and parameters
	actionTab := decodeMessage(fullAction)
	actionType := findValue(actionTab, "typ")
	actionParams := findValue(actionTab, "prm")

	// Handle player definition
	if actionType == "InitPlayer" {
		newPlayerTab := decodeMessage(actionParams)
		newPlayer := findValue(newPlayerTab, "newPlayer")
		// Check new player validity
		if newPlayer != "1" && newPlayer != "2" && newPlayer != "3" {
			// Player not valid (ignore)
			logWarning("handleAction", "Player Disconnected OR Player not valid (player reset)")
			lastConnectedPlayer = ""
			return game
		} else {
			// Valid player : updates
			lastConnectedPlayer = newPlayer
			// logInfo("handleAction", "App player set to "+lastConnectedPlayer+" game was reset.")
			// Game resets when player connects
			game = GameState{}
			game = getInitState()
			game = renewDrawPile(game)
			game = renewPlayerHands(game)
			return game
		}
	}

	// Check if app controls valid player
	if lastConnectedPlayer != "1" && lastConnectedPlayer != "2" && lastConnectedPlayer != "3" {
		logError("handleAction", "No player defined or player not recognized! Impossible to do other actions than player initialisation.")
		return game
	}

	// Process action
	switch actionType {
	case "ResetGame": // NO CONTROLS -> Resets the whole game (players, scores, decks, ...)
		game = GameState{}
		game = getInitState()
		game = renewDrawPile(game)
		game = renewPlayerHands(game)

	case "NewRound": // NO CONTROLS -> Starts a new game sleeve (deals new hands and new draw pile)
		game = renewPlayerHands(game)
		game = renewDrawPile(game)

	case "NextTurn": // NO CONTROLS -> Goes to the next sleeve when no other player wants to trade (deals new draw pile)
		game = renewDrawPile(game)

	case "Kems": // CONTROLS -> Increments player score if won (or does nothing)
		// Getting app player index
		appPlayerIndex, _ := strconv.Atoi(lastConnectedPlayer)
		appPlayerIndex -= 1

		// Player won
		if hasKems(game, appPlayerIndex) {
			// Add score to player
			game.Players[appPlayerIndex].Score += 1
			// Start new sleeve
			game = renewPlayerHands(game)
			game = renewDrawPile(game)
		}

	case "ContreKems": // CONTROLS -> Player counters another players win
		otherPlayerIndexString := decodeMessage(actionParams)
		otherPlayerIndex, err := strconv.Atoi(findValue(otherPlayerIndexString, "playerIndex"))
		// Check params
		if err != nil {
			logError("handleAction", "Error converting action params to integers for Contre Kems "+err.Error()+" action, (ignored) "+actionType)
			return game
		}
		otherPlayerIndex -= 1
		if otherPlayerIndex < 0 || otherPlayerIndex >= len(game.Players) {
			logError("handleAction", "Wrong params values, action (ignored) "+actionType)
			return game
		}

		// Player countered
		if hasKems(game, otherPlayerIndex) {
			game.Players[otherPlayerIndex].Score -= 1
			game = renewPlayerHands(game)
			game = renewDrawPile(game)
		}

	case "SwapCards": // CONTROLS -> Plater swaps one card of his hand with one of the draw pile
		// Get params for card swapping
		cardsIndexes := decodeMessage(actionParams)
		playerIndex, err1 := strconv.Atoi(findValue(cardsIndexes, "playerIndex"))
		playerCardIndex, err2 := strconv.Atoi(findValue(cardsIndexes, "playerCardIndex"))
		drawPileCardIndex, err3 := strconv.Atoi(findValue(cardsIndexes, "drawPileCardIndex"))
		// Check params
		if err1 != nil || err2 != nil || err3 != nil {
			logError("handleAction", "Error converting action params to integers for card swapping "+err1.Error()+err2.Error()+err3.Error()+" action, (ignored) "+actionType)
			return game
		}
		playerIndex -= 1
		if playerIndex < 0 || playerCardIndex < 0 || drawPileCardIndex < 0 || playerIndex >= len(game.Players) || playerCardIndex >= len(game.Players[playerIndex].Hand) || drawPileCardIndex >= len(game.DrawPile) {
			logError("handleAction", "Wrong params values, action (ignored) "+actionType)
			return game
		}
		// Update gamestate
		game = swapCard(game.Players[playerIndex].Hand[playerCardIndex], game.DrawPile[drawPileCardIndex], game.Players[playerIndex], game)

	case "SavePoint": // CONTROLS -> Saves the current game state
		// Save the game state

		logInfo("handleAction", "Send order received from base app, gamestate saved, sending to controller.")
		outChan <- encodeMessage([]string{"snd", "msg", "saveOrder"}, []string{name, "[SAVEORDER]" + gameStateToString(game), "1"}) + "\n"
		return game

	default: // Uknown action, ERROR
		// Action not recognized, send same game state (app should not share it)
		logError("handleAction", "Uknown action, (ignored) "+actionType)
		return game
	}

	// Sends updated (or not) game state (if not updated, app should not share it)
	return game
}

//#region GAME
//////////
// GAME //
//////////
func main() {
	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "app", "name")
	flag.Parse()
	name = *pName

	// Starting App
	// logInfo("main", "Launching app...")
	// Initialising key variables for app
	messageReceived := ""
	keyValTable := []string{}
	actionToDo := ""
	game := getInitState()
	game = renewPlayerHands(game)
	game = renewDrawPile(game)
	// Go routines to read and write input / output
	inChan = make(chan string, 10)
	outChan = make(chan string, 10)
	go read(inChan)
	go write(outChan)

	// Main loop of the app, manages message reception and emission and processing
	for {
		// logInfo("main", "Waiting for message.")
		// Message reception
		messageReceived = <-inChan
		logInfo("main", "Message received. "+messageReceived)

		// Determine message type for processing
		keyValTable = decodeMessage(messageReceived)
		sender := findValue(keyValTable, "snd")

		
		// Filter out obviously wrong messages that an app should not receive
		if len(sender) != 2 || len(name) != 2 || (sender != "C"+name[1:2] && sender[:1] != "P") {
			logError("main", "Message invalid sender OR invalid app name (ignored) - CAN BE FATAL!")
			messageReceived = ""
			continue
		}

		// CONTROLLER sent message
		// Filter out messages from our controller to other controllers
		if findValue(keyValTable, "hlg") != "" {
			logWarning("main", "Message from own controller to other controllers, (ignored).")
			messageReceived = ""
			continue
		}

		// PLAYER sent message (ask for exclusive access)
		if sender == "P"+lastConnectedPlayer || (sender[:1] == "P" && lastConnectedPlayer == "") {
			actionToDo = findValue(keyValTable, "msg")
			// Ask for exclusive access
			outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[ACRITICAL]"}) + "\n"
			logInfo("main", "Asked for exclusive access.")
			messageReceived = ""
			continue
		}
		


		// Getting message
		messageReceived = findValue(keyValTable, "msg")


		// Filter out wrong messages (just in case)
		if len(messageReceived) < 11 || (messageReceived[:11] != "[GAMESTATE]" && messageReceived[:11] != "[BCRITICAL]" && messageReceived[:11] != "[SAVEORDER]") {
			// logInfo("main", "Wrong message type for app received "+messageReceived+" (ignoring).")
			logInfo("main", "Wrong message type for app received (ignoring).")
			messageReceived = ""
			continue
		}


		// Message is an exclusive access grant => handle action
		if messageReceived[:11] == "[BCRITICAL]" {
			// Error if app is not trying to handle an action
			if actionToDo == "" {
				logError("main", "App received access but did not need it anymore (liberating)")
				outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[ECRITICAL]"}) + "\n"
				messageReceived = ""
				continue
			}

			oldGame := gameStateToString(game)
			game = handleAction(actionToDo, game)
			if oldGame == gameStateToString(game) {
				logWarning("main", "Action did not change game state, no update required. (Ended critical access)")
				outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[ECRITICAL]"}) + "\n"
			} else {
				logSuccess("main", "Gamestate updated, sending game update. (Ended critical access) + (Sent update to display)")
				outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, gameStateToString(game)}) + "\n"
				logInfo("main", "Ended critical access message sent.")
				outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[ECRITICAL]"}) + "\n"
				logInfo("main", "Sent update to display.")
				outChan <- gameStateToString(game) + "\n"

			}
			// Reset action (it has been processed)
			actionToDo = ""

			messageReceived = ""
			continue
		}


		// Message is a game state (process)
		// logInfo("main", "Processing game state... "+messageReceived)
		// Replace game state if an update was received
		if messageReceived[:11] == "[GAMESTATE]" {
			if gameStateToString(game) != messageReceived {
				game = stringToGameState(messageReceived)
				// Updated game state not sent anymore when update is received
				// outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, gameStateToString(game)}) + "\n"
				logInfo("main", "Updated game state (but did not diffuse the update) (update not sent to display).")
				outChan <- gameStateToString(game) + "\n" // Updating interface just in case went wrong last time
			} else {
				logSuccess("main", "Game state is already up to date, all apps up to date. (should not happen anymore)")
			}

			messageReceived = ""
			continue
		}

		// Message of save order from controller
		if messageReceived[:11] == "[SAVEORDER]" {
			logInfo("main", "Save order received, saving game state.")
			outChan <- encodeMessage([]string{"snd", "msg", "saveOrder"}, []string{name, "[SAVEORDER]" + gameStateToString(game), "0"}) + "\n"
		}

		logError("main", "CRITICAL ERROR, MESSAGE TREATMENT WAS NOT IMPLEMENTED (should never happen)")
		messageReceived = ""
	}
}
