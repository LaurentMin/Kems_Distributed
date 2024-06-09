package main

import (
	"flag"
	"time"
)

type MessageContent string

const (
	askToConnect     MessageContent = "Hello, may I join your awesome network ?"
	acceptConnection MessageContent = "Hello, of course you can join our network ?"
	refuseConnection MessageContent = "Hello, sorry but you'll have to wait ?"
)

/*
Connect go routine sends asks to connect to network untill stoped (waits a certain amount of time in between pings)
*/
func connect(stop <-chan bool, askNode string) {
	time.Sleep(5 * time.Second)
	for {
		select {
		case <-stop:
			logMessage("connect", "Connection go routine stopped.")
			return
		default:
			outChan <- encodeMessage([]string{"snd", "typ", "msg"}, []string{name, "con", string(askToConnect)}) + "\n"
			logInfo("connect", "Asked to join a network through : "+askNode)
			time.Sleep(5 * time.Second)
		}
	}
}

/*
NET is connected to network and a new node asks to join
*/
func handleConnectionMessage(sender string, msgcontent string) {
	// Connection message, accept
	if msgcontent == string(askToConnect) {
		outChan <- encodeMessage([]string{"snd", "typ", "msg"}, []string{name, "con", string(acceptConnection)}) + "\n"
		logInfo("handleConnectionMessage", "Connection accepted.")
	} else {
		logInfo("handleConnectionMessage", "Unexpected connection message, ignored.")
	}
}

/*
NET is connected to network and receives message from network
*/
func handleNetMessage(sender string, msgcontent string) {
	logInfo("handleNetMessage", "Function not implemented.")
}

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
	msgtype := ""
	keyValTable := []string{}

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
		// Filter out random messages
		invalidSender := len(sender) < 2 || (sender[0] != 'C' && sender[0] != 'N')
		if len(name) < 2 || invalidSender || msgtype == "" {
			logWarning("main", "NET has bad name or received wrong message (ignored) - CAN BE FATAL!")
			messageReceived = ""
			continue
		}

		/* HANDLE CONTROLLER MESSAGE */
		if sender[0] == 'C' && connected {
			outChan <- encodeMessage([]string{"snd", "typ", "msg"}, []string{name, "net", messageReceived}) + "\n"
			logInfo("main", "Controller message sent to network.")
			messageReceived = ""
			continue
		}

		/* HANDLE NETWORK MESSAGE */
		msgcontent := findValue(keyValTable, "msg")
		if sender[0] == 'N' && connected {
			switch msgtype {
			case "con":
				handleConnectionMessage(sender, msgcontent) // must log action
			case "net":
				handleNetMessage(sender, msgcontent) // must log action
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
