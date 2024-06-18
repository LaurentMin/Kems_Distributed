package main

import (
	"flag"
	"strconv"
	"strings"
	"time"
)

///// TEST
func testDiffusion(table *[]Diffusion, neighbours *[]string) {
	if name != "N1" {
		return
	}
	for i := 0; i < 20; i++ {
		logError("test", "Begin in "+strconv.Itoa(20-i))
		time.Sleep(time.Second)
	}
	for i := 0; i < 100; i++ {
		startDiffusion(i+100, "new", table, len(*neighbours))
	}
}

func testRemoving(table *[]Diffusion, neighbours *[]string) {
	if name != "N0" {
		return
	}
	for i := 0; i < 20; i++ {
		logError("test", "Removing in "+strconv.Itoa(20-i))
		time.Sleep(time.Second)
	}
	startDiffusion(6969, "del", table, len(*neighbours))
}

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

	time.Sleep(1 * time.Second)

	for {
		select {
		case <-stop:
			logMessage("connect", "Connection go routine stopped.")
			return
		default:
			outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, askNode, "con", string(askToConnect)}) + "\n"
			logInfo("connect", "Asked to join a network through : "+askNode)
			time.Sleep(30 * time.Second)
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
	value        string
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
func getDiffusioni(index string, neighbours int, val string) Diffusion {
	return Diffusion{
		diffIndex:    index,
		color:        blanc,
		parent:       "",
		nbNeighbours: neighbours,
		value:        val,
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
	diffID := name + "D" + strconv.Itoa(counter)

	// Create diffusion
	newDiff := getDiffusioni(diffID, nbNeighbours, val)
	newDiff.color = bleu
	newDiff.parent = name
	*table = append(*table, newDiff)

	// Create message
	diff := getDiffusionMessagei(diffID)
	diff.color = bleu
	diff.value = val
	outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, "all", "net", diffusionToString(diff)}) + "\n"

	logInfo("startDiffusion", "Diffused message to all neighbours : "+diffID)
}

//////////////////////////////////////////
////////// NET MESSAGE HANDLING //////////
//////////////////////////////////////////
/*
DIFFUSION
NET is connected to network and receives a net message (diffusion messages are net messages)
*/

func handleDiffusionMessage(sender string, recipient string, msgcontent string, table *[]Diffusion, neighbours *[]string, zombie *bool) {
	if len(msgcontent) < 11 || msgcontent[:11] != "[DIFFUSION]" {
		logError("handleDiffusionMessage", "Fatal error, net message content is corrupted (ignored).")
		return
	}
	numNeighbours := len(*neighbours)
	diffMessage := stringToDiffusion(msgcontent)
	if stopElecWave(*table, diffMessage) {
		logWarning("handleDiffusionMessage", "Stopped election wave for "+diffMessage.diffIndex)
		return
	}
	tabIndex := getDiffIdIndexOrCreateIfNotExists(table, diffMessage.diffIndex, numNeighbours, diffMessage.value)

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
			logInfo("handleDiffusionMessage", "Replied red message to sender.")
		}
	case rouge:
		(*table)[tabIndex].nbNeighbours -= 1
		if (*table)[tabIndex].nbNeighbours == 0 {
			if (*table)[tabIndex].parent == name {
				if len((*table)[tabIndex].value) > 1 && (*table)[tabIndex].value[:1] == "N" { // Had asked for election to add a node
					addNeighbour(neighbours, (*table)[tabIndex].value)
					outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, (*table)[tabIndex].value, "con", string(acceptConnection)}) + "\n"
					logSuccess("handleDiffusionMessage", "Election ended, connection accepted for "+(*table)[tabIndex].value)
				} else if diffMessage.value == "del" { // NODE has diffused del message, can deactivate
					*zombie = true
					logSuccess("handleDiffusionMessage", "Node successfully deactivated : "+diffMessage.diffIndex)
				} else { // NODE has diffused new message
					logSuccess("handleDiffusionMessage", "Diffusion terminée : "+diffMessage.diffIndex)
				}
			} else { // Not the diffusion initiator
				outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, (*table)[tabIndex].parent, "net", diffusionToString(diffMessage)}) + "\n"
				logInfo("handleDiffusionMessage", "Passing red message to parent.")
				// If diffusion was a controller message, send it to own controller here own controller message already ignored because initiator doesn't get here.

				if isDiffCtlMsg(diffMessage.value) {
					outChan <- diffMessage.value + "\n"
					logInfo("handleDiffusionMessage", "Transmitted message to controller.")

				} else if diffMessage.value == "new" || diffMessage.value == "del" {
					outChan <- encodeMessage([]string{"snd", "msg"}, []string{"C" + getOriginIndex(diffMessage.diffIndex), diffMessage.value}) + "\n"
					logInfo("handleDiffusionMessage", "Sent new or del message to controller.")
				}
			}
		} else {
			logWarning("handleDiffusionMessage", "Decremented node neighbour count.")
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

	inChan = make(chan string, 100)
	outChan = make(chan string, 100)
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

	lastCtlMsgHlg := -1
	zombie := false // when true, net does not send messages to controller anymore, continues to passively transmit net messages

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

	///////// tests /////
	// go testDiffusion(&diffTable, &neighbours)
	go testRemoving(&diffTable, &neighbours)

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
		hlg := findValue(keyValTable, "hlg")
		notControllerMessage := sender == "C"+name[1:] && hlg == "" // filter out ctl msg to app
		invalidSender := len(sender) < 2 || len(name) < 2 || (sender != "C"+name[1:] && sender[0] != 'N')
		messageForMe := strings.EqualFold(recipient, "all") || recipient == name || recipient == ""

		if invalidSender || !messageForMe || notControllerMessage {
			logWarning("main", "Message not for node (ignored) wrong controller or Sj message OR unexpected message - COULD BE FATAL!")
			messageReceived = ""
			continue
		}

		/* HANDLE CONTROLLER MESSAGE */
		if sender[0] == 'C' && connected && !zombie {
			hlgInt, err := strconv.Atoi(hlg)
			if hlgInt <= lastCtlMsgHlg || err != nil {
				logWarning("main", "This message was meant for a controller (ignored).")
				messageReceived = ""
				continue
			}
			lastCtlMsgHlg = hlgInt
			startDiffusion(counter, messageReceived, &diffTable, len(neighbours))
			counter++

			messageReceived = ""
			continue
		}

		/* HANDLE NETWORK MESSAGE */
		msgcontent := findValue(keyValTable, "msg")

		if sender[0] == 'N' && connected { // still goes here whether is zombie node
			switch msgtype {
			case "con":
				if len(neighbours) == 0 { // No other nodes in the network, can accept without election
					addNeighbour(&neighbours, sender)
					outChan <- encodeMessage([]string{"snd", "rec", "typ", "msg"}, []string{name, sender, "con", string(acceptConnection)}) + "\n"

					logSuccess("handleConnectionMessage", "Connection accepted for "+sender)
				} else if msgcontent == string(askToConnect) && canParticipateToElection(diffTable) {
					startDiffusion(counter, sender, &diffTable, len(neighbours))
					counter++

					logInfo("main", "Asked network to add new node.")
				} else {
					logWarning("main", "Can't participate to election or unexpected connection message (ignored).")
				}
			case "net":

				handleDiffusionMessage(sender, recipient, msgcontent, &diffTable, &neighbours, &zombie)

			default:
				logError("main", "Ignored network message.")
			}
			messageReceived = ""
			continue
		}

		// HANDLE NETWORK CONNECTION
		if sender[0] == 'N' && !connected && !zombie {
			switch msgcontent {

			case string(acceptConnection):
				stop <- true // channel initialised only if connected is false when program begins
				connected = true
				addNeighbour(&neighbours, sender) // Adds neighbour if does not exist
				startDiffusion(counter, "new", &diffTable, len(neighbours))
				counter++
				logSuccess("main", "Successfully connected to network.")

			case string(refuseConnection):
				logWarning("main", "Connection to network was not accepted.")

			default:
				if msgtype == "net" {
					logWarning("main", "Node not yet connected to network (ignored)")
				} else {
					logWarning("main", "Unexpected connection message, ignored.")
				}
			}
			messageReceived = ""
			continue
		}

		logWarning("main", "(ignored) Node certainly in zombie mode OR received unexpected message while not connected to network => FATAL.")
		messageReceived = ""
	}
}
