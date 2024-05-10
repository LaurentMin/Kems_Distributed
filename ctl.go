package main

import (
	"flag"
	"fmt"
	"strconv"
)

///////////
// CLOCK //
///////////
/*
	Clock adjustment
*/
func clockAdjustment(x, y int) int {
	logMessage("clockAdjustment", "Adjusting clock to max(local,received) + 1.")
	if x < y {
		return y + 1
	}
	return x + 1
}

//////////
// MAIN //
//////////
func main() {
	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "controller", "name")
	flag.Parse()
	name = *pName

	// Starting Controller
	logInfo("main", "Launching controller...")
	// Initialising key variables for controller
	messageReceived := ""
	keyValTable := []string{}
	clock := 0

	// Main loop of the controller, manages message reception and emission and processing
	for {
		logInfo("main", "Waiting for message.")
		// Message reception
		// fmt.Scanln(&messageReceived)
		// ReadString until '\n' delimiter (instead of Scanln)
		messageReceived = scanUntilNewline()
		messageReceived = messageReceived[:len(messageReceived)-1]
		logInfo("main", "Message received. "+messageReceived)

		// Defining local clock depending on received message
		logInfo("main", "Clock updating...")
		keyValTable = decodeMessage(messageReceived)
		clockReceivedStr := findValue(keyValTable, "hlg")
		if clockReceivedStr != "" {
			// Clock adjustment if message received from other controller
			clockReceived, err := strconv.Atoi(clockReceivedStr)
			if err != nil {
				logError("main", "Error converting string to int : "+err.Error())
				continue
			}
			clock = clockAdjustment(clock, clockReceived)
			logInfo("main", "Clock updated, message recived from other controller.")
		} else {
			// Incremented if message received from base app
			clock = clock + 1
			logInfo("main", "Clock updated, message recived from local app.")

		}

		// Message emission
		logInfo("main", "Sending message...")
		if clockReceivedStr != "" {
			// Sending to base app
			fmt.Printf(findValue(keyValTable, "msg") + "\n")
			logInfo("main", "Message sent to local app.")
		} else {
			// Sending to other controller
			fmt.Printf(encodeMessage([]string{"msg", "hlg"}, []string{messageReceived, strconv.Itoa(clock)}) + "\n")
			logInfo("main", "Message sent to other controller.")
		}

		messageReceived = ""
	}
}
