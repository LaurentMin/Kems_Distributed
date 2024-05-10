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
	// logMessage("clockAdjustment", "Adjusting clock to max(local,received) + 1.")
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
	// logInfo("main", "Launching controller...")
	// Initialising key variables for controller
	messageReceived := ""
	keyValTable := []string{}
	clock := 0

	// Main loop of the controller, manages message reception and emission and processing
	for {
		// logInfo("main", "Waiting for message.")
		// Message reception
		// fmt.Scanln(&messageReceived)
		// ReadString until '\n' delimiter (instead of Scanln)
		messageReceived = scanUntilNewline()
		// logInfo("main", "Message received. "+messageReceived)
		logInfo("main", "Message received "+messageReceived)

		// Determine message type for processing
		keyValTable = decodeMessage(messageReceived)
		sender := findValue(keyValTable, "sender")
		// Filter out random messages
		if len(sender) != 2 || len(name) != 2 {
			logError("main", "Message invalid sender OR wrong ctl name (ignored) - CAN BE FATAL!")
			messageReceived = ""
			continue
		}

		// Defining local clock depending on received message (ignores messages from other controllers to their own apps)
		// logInfo("main", "Clock updating...")
		clockReceivedStr := findValue(keyValTable, "hlg")
		if clockReceivedStr != "" && sender[:1] == "C" { // Filters out messages from an app with a clock
			// Clock adjustment if message received from other controller
			// In this case the message is from a controller to another controller
			clockReceived, err := strconv.Atoi(clockReceivedStr)
			if err != nil {
				logError("main", "Error converting string to int : "+err.Error())
				continue
			}
			clock = clockAdjustment(clock, clockReceived)
			// logInfo("main", "Clock updated, message received from other controller.")

		} else if clockReceivedStr == "" && sender == "A"+name[1:2] { // Filters out messages without a clock from the wrong app.
			// Incremented if message received from base app
			clock = clock + 1
			// logInfo("main", "Clock updated, message received from local app.")
		} else { // Filters out messages from other controller to it's own app
			// ERROR, ignoring
			logError("main", "Unexpected ERROR, message was not supposed to be built this way.")
			messageReceived = ""
			continue
		}

		// Message emission
		// logInfo("main", "Sending message...")
		if clockReceivedStr != "" {
			// Sending to base app
			fmt.Printf(encodeMessage([]string{"msg", "sender"}, []string{findValue(keyValTable, "msg"), name}) + "\n")
			logInfo("main", "Message sent to local app.")
		} else {
			// Sending to other controller
			fmt.Printf(encodeMessage([]string{"msg", "hlg", "sender"}, []string{messageReceived, strconv.Itoa(clock), name}) + "\n")
			logInfo("main", "Message sent to other controller.")
		}

		messageReceived = ""
	}
}
