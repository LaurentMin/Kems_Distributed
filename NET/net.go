package main

import (
	"flag"
	"fmt"
	"regexp"
)

func main() {
	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "N1", "name")
	pAskNode := flag.String("a", "N1", "name of node to connect to")
	flag.Parse()

	name := *pName
	fmt.Println("Starting ", name)

	askNode := *pAskNode
	fmt.Println("Ask to ", askNode)

	if askNode != "" {
		// TODO : executer ./reseauNet.sh name askNode
	}

	inChan = make(chan string, 10)
	outChan = make(chan string, 10)
	// Reading go routine (sends read data from sdtin through channel)
	go read(inChan)

	// Writing go routine (writes data from channel to stdout)
	go write(outChan)

	messageReceived := ""
	sender := ""
	msgtype := ""
	keyValTable := []string{}

	for {
		logInfo("main", "Waiting for message.")
		// Message reception
		messageReceived = <-inChan
		logInfo("main", "Message received : "+messageReceived)

		// Determine message type for processing
		keyValTable = decodeMessage(messageReceived)
		sender = findValue(keyValTable, "snd")
		msgtype = findValue(keyValTable, "typ")
		logMessage("main", sender)
		logMessage("main", msgtype)
		// Filter out random messages
		validSender, _ := regexp.MatchString("(N|C)[0-9]+", sender)
		if len(name) < 2 || len(sender) < 2 || !validSender || msgtype == "" {
			logWarning("main", "NET has bad name or received wrong message (ignored) - CAN BE FATAL!")
			messageReceived = ""
			continue
		}

		// Send controller message to network
		if sender[0] == 'C' {
			outChan <- encodeMessage([]string{"snd", "typ", "msg"}, []string{name, "net", messageReceived})

			messageReceived = ""
			continue
		}

		// Handle natwork message
		if sender[0] == 'N' {
			msgcontent := findValue(keyValTable, "msg")
			outChan <- msgcontent

			messageReceived = ""
			continue
		}

		logError("main", "FATAL this should never be reached!")
		messageReceived = ""
	}
}
