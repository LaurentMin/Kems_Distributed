package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
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

/*
	estampille structure
*/
type Request struct {
	Type  string
	Clock int
}

//////////////////
// HELPER FUNCS //
//////////////////
/*
	Returns true if siteN has smallest date in estampilles
	(called to start a critical section for example)
*/
func canGoCritical(estampilles []Request, site int) bool {
	for i := 0; i < len(estampilles); i++ {
		logError("canGoCritical", "num site : "+ strconv.Itoa(i+1) + " clock : " + strconv.Itoa(estampilles[i].Clock))
		if estampilles[site].Clock > estampilles[i].Clock || (estampilles[site].Clock == estampilles[i].Clock && site > i) {
			return false
		}
	}
	return true
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
	estampilles := []Request{Request{"[ECRITICAL]", 0}, Request{"[ECRITICAL]", 0}, Request{"[ECRITICAL]", 0}} // Index 0..2 corresponds to controllers 1..3
	siteNum, _ := strconv.Atoi(name[1:2])                                                                     // Ok if this makes app crash (name must be defined)
	siteNum -= 1

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
			logWarning("main", "Display message OR invalid sender OR wrong ctl name (ignored) - CAN BE FATAL!")
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
				logError("main", "Error converting string to int : "+err.Error()+" (FATAT, clock corruption)")
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
			logWarning("main", "Message from another controller to it's own app (IGNORED) OR UNEXPECTED ERROR.")
			messageReceived = ""
			continue
		}

		// Message processing
		// getting message
		messageReceived = findValue(keyValTable, "msg")
		// Filter out wrong messages
		if len(messageReceived) < 11 {
			// logInfo("main", "Wrong message type for app received "+messageReceived+" (ignoring).")
			logInfo("main", "Wrong message type for app received (too short) (ignoring).")
			messageReceived = ""
			continue
		}

		// Receive from controller
		// logInfo("main", "Sending message...")
		if clockReceivedStr != "" && sender[:1] == "C" {
			otherSiteNumber, _ := strconv.Atoi(sender[1:2])
			otherSiteNumber -= 1
			switch messageReceived[:11] {
			case "[GAMESTATE]":
				fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{name, messageReceived}) + "\n")
				logInfo("main", "Gamestate message sent to local app.")
				time.Sleep(1 * time.Second)

			case "[ACRITICAL]": // Other controller asks for access restriction
				estampilles[otherSiteNumber].Type = "[ACRITICAL]"
				estampilles[otherSiteNumber].Clock = clock
				fmt.Printf(encodeMessage([]string{"snd", "hlg", "msg"}, []string{name, strconv.Itoa(clock), "[VCRITICAL]"}) + "\n")
				logInfo("main", "Answered to other controller restriction access demand.")
				time.Sleep(1 * time.Second)
				// Check if can start own critical
				if estampilles[siteNum].Type == "[ACRITICAL]" && canGoCritical(estampilles, siteNum) {
					fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{name, "[BCRITICAL]"}) + "\n")
					logInfo("main", "Begin critical section sent to base app.")
					time.Sleep(1 * time.Second)
				}

			case "[VCRITICAL]": // Other controller validates request reception
				// Do not replace an ask by a reception
				if estampilles[otherSiteNumber].Type != "[ACRITICAL]" {
					estampilles[otherSiteNumber].Type = "[VCRITICAL]"
					estampilles[otherSiteNumber].Clock = clock
				}
				logInfo("main", "Critical section reception was confirmed by other controller.")
				// Check if can start own critical
				if estampilles[siteNum].Type == "[ACRITICAL]" && canGoCritical(estampilles, siteNum) {
					fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{name, "[BCRITICAL]"}) + "\n")
					logInfo("main", "Begin critical section sent to base app.")
					time.Sleep(1 * time.Second)
				}

			case "[ECRITICAL]": // Other controller liberates access restriction
				estampilles[otherSiteNumber].Type = "[ECRITICAL]"
				estampilles[otherSiteNumber].Clock = clock
				logInfo("main", "Other controller ended restriction access.")
				// Check if can start own critical
				if estampilles[siteNum].Type == "[ACRITICAL]" && canGoCritical(estampilles, siteNum) {
					fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{name, "[BCRITICAL]"}) + "\n")
					logInfo("main", "Begin critical section sent to base app.")
					time.Sleep(1 * time.Second)
				}

			default:
				logError("main", "Wrong message type received (controller sent wrong type) (ignoring) (could be critical for clock).")
			}

			messageReceived = ""
			continue
		}

		// Received from app
		if clockReceivedStr == "" && sender == "A"+name[1:2] {
			switch messageReceived[:11] {
			case "[GAMESTATE]":
				fmt.Printf(encodeMessage([]string{"snd", "hlg", "msg"}, []string{name, strconv.Itoa(clock), messageReceived}) + "\n")
				logInfo("main", "Gamestate message sent to other controller.")
				time.Sleep(1 * time.Second)

			case "[ACRITICAL]": // Base app asks critical (asking other controllers)
				estampilles[siteNum].Type = "[ACRITICAL]"
				estampilles[siteNum].Clock = clock
				fmt.Printf(encodeMessage([]string{"snd", "hlg", "msg"}, []string{name, strconv.Itoa(clock), "[ACRITICAL]"}) + "\n")
				logInfo("main", "Asked other controllers for access restriction.")
				time.Sleep(1 * time.Second)

			case "[ECRITICAL]": // Base app stops critical (liberating other controllers)
				estampilles[siteNum].Type = "[ECRITICAL]"
				estampilles[siteNum].Clock = clock
				fmt.Printf(encodeMessage([]string{"snd", "hlg", "msg"}, []string{name, strconv.Itoa(clock), "[ECRITICAL]"}) + "\n")
				logInfo("main", "Liberated other controllers from access restriction.")
				time.Sleep(1 * time.Second)

			default:
				logError("main", "Wrong message type received (app sent wrong type) (ignoring) (could be critical for clock).")
			}

			messageReceived = ""
			continue
		}

		logError("main", "CRITICAL ERROR, MESSAGE TREATMENT WAS NOT IMPLEMENTED (should never happen) (bad for clock)")
		messageReceived = ""
	}
}
