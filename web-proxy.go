package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

func cardsToJSON(cards []Card) string {
	jsonString := "["
	for i, card := range cards {
		if i != 0 {
			jsonString += ", "
		}
		jsonString += "{\"value\": \"" + card.Value + "\", \"suit\": \"" + card.Suit + "\"}"
	}
	jsonString += "]"
	return jsonString
}

func scoresToJSON(players []Player) string {
	jsonString := "["
	for i, player := range players {
		if i != 0 {
			jsonString += ", "
		}
		jsonString += strconv.Itoa(player.Score)
	}
	jsonString += "]"
	return jsonString
}

func do_webserver(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bonjour depuis le serveur web en Go !")
}

func do_websocket(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	cnx, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logError("Web proxy", "Websocket upgrade failed"+err.Error())
		return
	}
	logSuccess("Web proxy", "Websocket connected")
	ws = cnx
	sendAction("InitPlayer", []string{"newPlayer"}, []string{name})
	for {
		_, message, err := cnx.ReadMessage()
		if err != nil {
			logError("Web proxy", "Error reading message from websocket: "+err.Error())
			return
		}
		logInfo("Web proxy", "Received message "+string(message))
		// ------------------------------- WEB to APP -------------------------------

		// Declared an empty interface
		var result map[string]interface{}

		// Unmarshal or Decode the JSON to the interface
		err = json.Unmarshal([]byte(message), &result)
		if err != nil {
			logError("Web proxy", "Error unmarshalling JSON: "+err.Error())
			continue
		}

		actionType := result["action"].(string)

		switch actionType {
		case "NextTurn":
			sendAction("NextTurn", []string{"playerIndex"}, []string{name})
		case "SwapCards":
			handCardValue := result["handCardValue"].(string)
			handCardSuit := result["handCardSuit"].(string)
			drawCardValue := result["drawCardValue"].(string)
			drawCardSuit := result["drawCardSuit"].(string)

			handCardIndex := -1
			drawCardIndex := -1
			for i, card := range state.Players[playerId-1].Hand {
				if card.Value == handCardValue && card.Suit == handCardSuit {
					handCardIndex = i
				}
			}
			for i, card := range state.DrawPile {
				if card.Value == drawCardValue && card.Suit == drawCardSuit {
					drawCardIndex = i
				}
			}
			sendAction("SwapCards", []string{"playerIndex", "playerCardIndex", "drawPileCardIndex"}, []string{name, strconv.Itoa(handCardIndex), strconv.Itoa(drawCardIndex)})
		case "Kems":
			sendAction("Kems", []string{"playerIndex"}, []string{name})
		case "ContreKems":
			otherPlayerId := -1
			for otherPlayerId == -1 {
				otherPlayerId = checkIfKems(state)
			}
			otherPlayerId++

			sendAction("ContreKems", []string{"playerIndex"}, []string{strconv.Itoa(otherPlayerId)})
		case "ResetGame":
			sendAction("ResetGame", []string{"playerIndex"}, []string{name})
		default:
			logError("Web proxy", "Unknown action type: "+actionType)
		}

	}
}

func ws_send(msg string) {
	if ws == nil {
		logError("Web proxy", "Websocket not opened")
	} else {
		err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			logError("Web proxy", "Error sending message to websocket: "+err.Error())
		} else {
			logInfo("Web proxy", "Sent message "+msg)
			fmt.Println(msg)
		}
	}
}

func handleBackUpdate(update string, playerId int) {

	messageNotForClient := len(update) < 11 || update[:11] != "[GAMESTATE]"
	if messageNotForClient {
		return
	}

	logInfo("Web proxy", "Handling update "+update)
	newState := stringToGameState(update)

	if len(newState.Players) == 0 {
		logInfo("main", "Wrong state received ignoring.")
		return
	}

	// ------------------------------- APP to WEB -------------------------------
	sendGameStateToWeb(newState, playerId)
	return
}

func sendGameStateToWeb(newState GameState, playerId int) {
	// Send the player's hand, the draw pile, the last discard pile and the players' scores in JSON format stringified

	numberRound := 1
	for i := 0; i < len(newState.Players); i++ {
		numberRound += newState.Players[i].Score
	}

	jsonGameState := "{\"playerId\":" + name + ", \"hand\": " + cardsToJSON(newState.Players[playerId-1].Hand) + ", \"drawPile\": " + cardsToJSON(newState.DrawPile) + ", \"scores\": " + scoresToJSON(newState.Players) + ", \"round\": " + strconv.Itoa(numberRound)

	//Discard can be empty at the beginning of the game
	if len(newState.DiscardPile) > 0 {
		lastDiscard := newState.DiscardPile[len(newState.DiscardPile)-1]
		jsonLastDiscard := ", \"discardPile\": " + cardsToJSON([]Card{lastDiscard})
		jsonGameState += jsonLastDiscard
	}

	potentialWinner := checkIfKems(newState)
	if potentialWinner != -1 {
		jsonGameState += ", \"potentialWinner\": " + strconv.Itoa(potentialWinner+1)
	}

	if len(state.Players) > 0 {
		isEndRound := false
		for i := 0; i < len(newState.Players); i++ {
			if newState.Players[i].Score != state.Players[i].Score {
				isEndRound = true
			}
		}
		if isEndRound {
			jsonGameState += ", \"newRound\": true"
		}
	}
	jsonGameState += "}"

	state = newState
	ws_send(jsonGameState)
}

func listenAppUpdateAndTransmit(playerId int) {
	scanner := bufio.NewScanner(os.Stdin)
	logInfo("Web proxy", "Listening for updates from the app")
	for scanner.Scan() {
		update := scanner.Text()
		handleBackUpdate(update, playerId)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
	}
}

var ws *websocket.Conn = nil
var state GameState
var playerId int

func main() {
	var port = flag.String("p", "4444", "n° de port")
	var addr = flag.String("a", "localhost", "nom/adresse machine")
	var playerName = flag.String("n", "1", "nom du joueur (1,2,3)")
	flag.Parse()
	name = *playerName

	var err error
	playerId, err = strconv.Atoi(name)
	if err != nil {
		logError("Web proxy", "Error converting player name to int: "+err.Error())
		return
	}

	go listenAppUpdateAndTransmit(playerId)

	logInfo("Web proxy", "Starting web server on "+*addr+":"+*port)
	http.HandleFunc("/", do_webserver)
	http.HandleFunc("/ws", do_websocket)
	http.ListenAndServe(*addr+":"+*port, nil) // Listen on port 4444
}
