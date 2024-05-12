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
		messageReceived = scanUntilNewline()
		logInfo("main", "Message received "+messageReceived)

		// Determine message type for processing
		keyValTable = decodeMessage(messageReceived)
		sender := findValue(keyValTable, "snd")

		// Filter out random messages (Display messages for example)
		if len(sender) != 2 || len(name) != 2 || (sender != "A"+name[1:2] && sender[:1] != "C") {
			logError("main", "Display message OR invalid sender OR wrong ctl name (ignored) - CAN BE FATAL!")
			messageReceived = ""
			continue
		}

		// Clock updating

		// logInfo("main", "Clock updating...")
		clockReceivedStr := findValue(keyValTable, "hlg")
		if clockReceivedStr != "" && sender[:1] == "C" { // Filters out messages from an app with a clock (should never happen)
			// Clock adjustment if message received from other controller
			clockReceived, err := strconv.Atoi(clockReceivedStr)
			if err != nil {
				logError("main", "Error converting string to int : "+err.Error())
				continue
			}
			clock = clockAdjustment(clock, clockReceived)
			// logInfo("main", "Clock updated, message received from other controller.")

		} else if clockReceivedStr == "" && sender == "A"+name[1:2] { // Filters out messages without a clock from the wrong app or a controller.
			// Incremented if message received from app
			clock = clock + 1
			// logInfo("main", "Clock updated, message received from local app.")
		} else { // Filters out messages from other controller to their own app or other errors
			// ERROR, ignoring
			logError("main", "Message from another controller to it's own app (IGNORED) OR UNEXPECTED ERROR.")
			messageReceived = ""
			continue
		}

		// Message emission
		// getting message
		messageReceived = findValue(keyValTable, "msg")
		// Filter out wrong messages
		if len(messageReceived) < 11 {
			// logInfo("main", "Wrong message type for app received "+messageReceived+" (ignoring).")
			logInfo("main", "Wrong message type for app received (too short) (ignoring).")
			messageReceived = ""
			continue
		}

		// Controller sent message (sending to app)
		// logInfo("main", "Sending message...")
		if clockReceivedStr != "" && sender[:1] == "C" {
			// Filter out wrong messages
			if messageReceived[:11] != "[GAMESTATE]" && messageReceived[:11] != "[ACRITICAL]" && messageReceived[:11] != "[VCRITICAL]" {
				// logInfo("main", "Wrong message type for app received "+messageReceived+" (ignoring).")
				logInfo("main", "Wrong message type for app received (controller sent wrong type) (ignoring).")
				messageReceived = ""
				continue
			}

			// Message from controller, sending to base app
			fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{name, findValue(keyValTable, "msg")}) + "\n")
			logInfo("main", "Message sent to local app.")

			messageReceived = ""
			continue
		}

		// App sent message (sending to controllers)
		if clockReceivedStr == "" && sender == "A"+name[1:2] {
			// Filter out wrong messages
			if messageReceived[:11] != "[GAMESTATE]" && messageReceived[:11] != "[BCRITICAL]" && messageReceived[:11] != "[ECRITICAL]" {
				// logInfo("main", "Wrong message type for app received "+messageReceived+" (ignoring).")
				logInfo("main", "Wrong message type for app received (app sent wrong type) (ignoring).")
				messageReceived = ""
				continue
			}

			fmt.Printf(encodeMessage([]string{"snd", "hlg", "msg"}, []string{name, strconv.Itoa(clock), findValue(keyValTable, "msg")}) + "\n")
			logInfo("main", "Message sent to other controller.")

			messageReceived = ""
			continue
		}

		logError("main", "CRITICAL ERROR, MESSAGE TREATMENT WAS NOT IMPLEMENTED (should never happen)")
		messageReceived = ""
	}
}
