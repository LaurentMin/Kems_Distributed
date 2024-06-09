package main

import (
	"flag"
	"strconv"
	"strings"
	"time"
)

////////////////////////////////
////////// CONNECTION //////////
////////////////////////////////
/*
Connection messages
*/
type MessageContent string

const (
	askToConnect     MessageContent = "Hello, may I join your awesome network ?"
	acceptConnection MessageContent = "Hello, of course you can join our network ?"
	refuseConnection MessageContent = "Hello, sorry but you'll have to wait ?"
)

/*
Connect go routine asks to connect to network until stopped (waits a certain amount of time in between pings)
*/
func connect(stop <-chan bool, askNode string) {
	time.Sleep(5 * time.Second)
	for {
		select {
		case <-stop:
			logMessage("connect", "Connection go routine stopped.")
			return
		default:
			outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, askNode, "con", string(askToConnect)}) + "\n"
			logInfo("connect", "Asked to join a network through : "+askNode)
			time.Sleep(5 * time.Second)
		}
	}
}

/*
Connection message handling
NET is connected to network and a new node asks to join
*/
func handleConnectionMessage(sender string, msgcontent string, neighbours *[]string) {
	// Connection message, accept
	if msgcontent == string(askToConnect) {
		addNeighbour(neighbours, sender)
		outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, sender, "con", string(acceptConnection)}) + "\n"
		logInfo("handleConnectionMessage", "Connection accepted.")
	} else {
		logInfo("handleConnectionMessage", "Unexpected connection message, ignored.")
	}
}

///////////////////////////////
////////// DIFFUSION //////////
///////////////////////////////
type Color string

const (
	blanc Color = "blanc"
	bleu  Color = "bleu"
	rouge Color = "rouge"
)

/*
Type to have multiple diffusions at the same time
A diffusion index is the name of a node concatenated to the id of it's diffusion
*/
type Diffusion struct {
	diffIndex    string
	color        Color
	parent       string
	nbNeighbours int
}

/*
Used for messages
*/
type DiffusionMessage struct {
	diffIndex string
	color     Color
	value     string
}

/*
Diffusion constructor
*/
func getDiffusioni(index string, neighbours int) Diffusion {
	return Diffusion{
		diffIndex:    index,
		color:        blanc,
		parent:       "",
		nbNeighbours: neighbours,
	}
}

/*
Diffusion message
*/
func getDiffusionMessagei(index string) DiffusionMessage {
	return DiffusionMessage{
		diffIndex: index,
		color:     blanc,
		value:     "",
	}
}

/*
Starts a diffusion from node
*/
func startDiffusion(counter int, val string, table *[]Diffusion, nbNeighbours int) {
	diffID := name + strconv.Itoa(counter)

	// Create diffusion
	newDiff := getDiffusioni(diffID, nbNeighbours)
	newDiff.color = bleu
	newDiff.parent = name
	*table = append(*table, newDiff)

	// Create message
	diff := getDiffusionMessagei(diffID)
	diff.color = bleu
	diff.value = val
	outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, "all", "net", diffusionToString(diff)}) + "\n"
	logInfo("startDiffusion", "Diffused message to all neighbours.")
}

//////////////////////////////////////////
////////// NET MESSAGE HANDLING //////////
//////////////////////////////////////////
/*
DIFFUSION
NET is connected to network and receives a net message (only diffusion messages are net messages => for now)
*/
func handleDiffusionMessage(sender string, recipient string, msgcontent string, table *[]Diffusion, numNeighbours int) {
	if len(msgcontent) < 11 || msgcontent[:11] != "[DIFFUSION]" {
		logError("handleDiffusionMessage", "Fatal error, net message content is corrupted (ignored).")
		return
	}

	diffMessage := stringToDiffusion(msgcontent) // must ignore numNeighbours and parent from this object
	tabIndex := getDiffIdIndexOrCreateIfNotExists(table, diffMessage.diffIndex, numNeighbours)
	logError("handleDiffusionMessage", printDiffusion((*table)[tabIndex]))
	switch diffMessage.color {
	case bleu:
		if (*table)[tabIndex].color == blanc {
			(*table)[tabIndex].color = bleu
			(*table)[tabIndex].parent = sender
			outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, "all", "net", diffusionToString(diffMessage)}) + "\n"
			logInfo("handleDiffusionMessage", "Sent blue message to neighbours.")
		} else {
			diffMessage.color = rouge
			outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, sender, "net", diffusionToString(diffMessage)}) + "\n"
			logInfo("handleDiffusionMessage", "Sent red message to sender.")
		}
	case rouge:
		(*table)[tabIndex].nbNeighbours -= 1
		if (*table)[tabIndex].nbNeighbours == 0 {
			if (*table)[tabIndex].parent == name {
				logSuccess("handleDiffusionMessage", "Diffusion terminée : "+diffMessage.diffIndex)
			} else {
				outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, (*table)[tabIndex].parent, "net", diffusionToString(diffMessage)}) + "\n"
				logInfo("handleDiffusionMessage", "Passing red message to parent.")
			}
		}
	default:
		logError("handleDiffusionMessage", "Fatal error, diffusion message has unexpected color (ignored).")
		return
	}
}

//////////////////////////
////////// MAIN //////////
//////////////////////////
func main() {
	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "default", "name")
	pAskNode := flag.String("a", "default", "name of node to connect to")
	flag.Parse()
	name = *pName
	askNode := *pAskNode

	inChan = make(chan string, 10)
	outChan = make(chan string, 10)
	// Reading go routine (sends read data from stdin through channel)
	go read(inChan)
	// Writing go routine (writes data from channel to stdout)
	go write(outChan)

	// Program variables
	var stop chan bool
	messageReceived := ""
	sender := ""
	recipient := ""
	msgtype := ""
	keyValTable := []string{}
	counter := 0
	diffTable := []Diffusion{}
	neighbours := []string{}

	// Ask to join network
	connected := false
	if name != askNode { // First node  of the network has itself as askNode
		stop = make(chan bool, 10)
		go connect(stop, askNode)
		outChan <- "ping" + "\n"
		logInfo("main", "Started new node on existing "+askNode)
	} else {
		connected = true // First node of the network
		outChan <- "ping" + "\n"
		logInfo("main", "Started a new network.")
	}

	// Main message handling loop
	for {
		logInfo("main", "Waiting for message.")
		// Message reception
		messageReceived = <-inChan
		logInfo("main", "Message received : "+messageReceived)

		// "easter egg"
		if messageReceived == "ping" || messageReceived == "pong" {
			if messageReceived == "ping" {
				outChan <- "pong" + "\n"
				logInfo("main", "Replied to ping.")
			}
			messageReceived = ""
			continue
		}

		// Determine message type for processing
		keyValTable = decodeMessage(messageReceived)
		sender = findValue(keyValTable, "snd")
		msgtype = findValue(keyValTable, "typ")
		recipient = findValue(keyValTable, "rec")
		// Filter out random messages
		invalidSender := len(sender) < 2 || (sender[0] != 'C' && sender[0] != 'N')
		messageForMe := strings.EqualFold(recipient, "all") || recipient == name
		if len(name) < 2 || invalidSender || !messageForMe || msgtype == "" {
			logWarning("main", "Message not for node (ignored) OR unexpected message - COULD BE FATAL!")
			messageReceived = ""
			continue
		}

		/* HANDLE CONTROLLER MESSAGE */
		if sender[0] == 'C' && connected {
			// outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, "all", "net", messageReceived}) + "\n"
			logInfo("main", "Controller message sent to network.")
			messageReceived = ""
			continue
		}

		/* HANDLE NETWORK MESSAGE */
		msgcontent := findValue(keyValTable, "msg")
		if sender[0] == 'N' && connected {
			switch msgtype {
			case "con":
				handleConnectionMessage(sender, msgcontent, &neighbours) // must log action
			case "net":
				handleDiffusionMessage(sender, recipient, msgcontent, &diffTable, len(neighbours))
			default:
				logError("main", "Ignored network message.")
			}
			messageReceived = ""
			continue
		}

		// HANDLE NETWORK CONNECTION
		if sender[0] == 'N' && !connected {
			switch msgcontent {
			case string(acceptConnection):
				stop <- true // channel initialised only if connected is false when program begins
				connected = true
				addNeighbour(&neighbours, sender) // Adds neighbour if does not exist
				startDiffusion(counter, name, &diffTable, len(neighbours))
				logSuccess("main", "Successfully connected to network.")
			case string(refuseConnection):
				logWarning("main", "Connection to network was not accepted.")
			default:
				logWarning("main", "Unexpected connection message, ignored.")
			}
			messageReceived = ""
			continue
		}

		logError("main", "FATAL ! Node received unexpected message while not connected to network. (usually from controller)")
		messageReceived = ""
	}
}
