package main

import (
	"fmt"
	"log"
	"os"
	"flag"
	"strconv"
	"sync"
	"time"
	"math/rand"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
//#region Struct Appli
type Card struct {
	Value string
	Suit  string
}

type Player struct {
	id  int
	Hand []Card
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//#region Func Appli
func newDeck() []Card {
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}
	var deck []Card
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Value: value, Suit: suit})
		}
	}
	return deck
}

func shuffleDeck(deck []Card) []Card {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

// Pour échanger les cartes sur le plateau (hand) et les cartes de la main (board)
func dealCards(hand *[]Card, board *[]Card, card_hand int, card_board int) {
	(*hand)[card_hand], (*board)[card_board] = (*board)[card_board], (*hand)[card_hand]
}

func pickCard(deck *[]Card) Card {
	card := (*deck)[0]
	*deck = (*deck)[1:]
	return card
}

func addCardTo(pack *[]Card, card Card) {
	*pack = append(*pack, card)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//#region Affichage
func sendToCtl(msg string, nom string) {
	fmt.Println(msg, nom)
}

func sendToTerm(msg string, nom string) {
	l := log.New(os.Stderr, "", 0)
	l.Println(msg, nom)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//#region Func Comm
func sendperiodic() {
	var sndmsg string
	var i int

	i = 0

	for {
		mutex.Lock()
		i = i + 1
		sndmsg = "message_" + strconv.Itoa(i) + "\n"
		fmt.Print(sndmsg)
		mutex.Unlock()
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func receive() {
	var rcvmsg string
	l := log.New(os.Stderr, "", 0)

	for {
		fmt.Scanln(&rcvmsg)
		mutex.Lock()
		l.Println("reception <", rcvmsg, ">")
		mutex.Unlock()
		rcvmsg = ""
	}
}

var mutex = &sync.Mutex{}

////////////////////////////////////////////////////////////////////////////////////////////////////
//#region Main
func main() {

	stderr := log.New(os.Stderr, "", 0)

	// require flag
	p_nom := flag.String("n", "default", "nom")
    flag.Parse()
	if *p_nom == "" {
		stderr.Println("Le nom est obligatoire")
		os.Exit(1)
	}

    nom := *p_nom + "-" + strconv.Itoa(os.Getpid())
	fmt.Println("Joueur "+ nom)

	deck := newDeck()
	deck = shuffleDeck(deck)
	fmt.Println(deck)
	fmt.Println(len(deck))

	board := [4]Card{}
	for i := 0; i < 4; i++ {
		board[i] = pickCard(&deck)
	}
	fmt.Println(board)

	player1 := Player{id: 1}
	player2 := Player{id: 2}
	player3 := Player{id: 3}
	player4 := Player{id: 4}

	for i := 0; i < 4; i++ {
		addCardTo(&(player1.Hand), pickCard(&deck))
		addCardTo(&(player2.Hand), pickCard(&deck))
		addCardTo(&(player3.Hand), pickCard(&deck))
		addCardTo(&(player4.Hand), pickCard(&deck))
	}

	fmt.Println(player1.Hand)
	fmt.Println(player2.Hand)
	fmt.Println(player3.Hand)
	fmt.Println(player4.Hand)

	fmt.Println(len(deck))

	/* comm
	go receive()
	for {
		time.Sleep(time.Duration(60) * time.Second)
	}
	*/
}
