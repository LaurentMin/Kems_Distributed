package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
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
	Vector clock adjustment
*/
func vClockAdjustment(x, y []int, ind int) []int {
	// logMessage("vClockAdjustment", "Adjusting vector clock to max(local,received) + 1.")
	for i := 0; i < len(x); i++ {
		if x[i] < y[i] {
			x[i] = y[i]
		} else {
			x[i] = x[i]
		}
	}
	x[ind] = x[ind] + 1
	return x
}

/*
	Cast string to vector clock*
*/
func castStringToVClock(strVlg string) []int {
    strVlg = strings.Trim(strVlg, "[]")
    elements := strings.Split(strVlg, " ")

    var vlg []int

	for _, element := range elements {
        num, err := strconv.Atoi(element)
        if err != nil {
            panic(err)
        }
        vlg = append(vlg, num)
    }

	return vlg
}

/*
	Cast vector clock to string
*/
func castVClockToString(vlg []int) string {
	var strVlg string

	for _, element := range vlg {
		strVlg += strconv.Itoa(element) + " "
	}

	return strVlg
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
	vClock := []int{0, 0, 0}
	// Find the controller number in vClock
	idVClock, err := strconv.Atoi(name[len(name)-1:])
	logInfo("main", "idVClock: "+strconv.Itoa(idVClock))
	if err != nil {
		logError("main", "Error converting string to int for idVClock: "+err.Error())
	}

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
		sender := findValue(keyValTable, "snd")
		// Filter out random messages (Display messages for instance)
		if len(sender) != 2 || len(name) != 2 {
			logError("main", "Display message OR invalid sender OR wrong ctl name (ignored) - CAN BE FATAL!")
			messageReceived = ""
			continue
		}

		// Defining local clock depending on received message (ignores messages from other controllers to their own apps)
		// logInfo("main", "Clock updating...")
		clockReceivedStr := findValue(keyValTable, "hlg")
		vClockReceivedStr := findValue(keyValTable, "vlg")
		if clockReceivedStr != "" && vClockReceivedStr != "" && sender[:1] == "C" { // Filters out messages from an app with a clock
			// Clock adjustment if message received from other controller
			// In this case the message is from a controller to another controller
			clockReceived, err := strconv.Atoi(clockReceivedStr)
			if err != nil {
				logError("main", "Error converting string to int : "+err.Error())
				continue
			}
			clock = clockAdjustment(clock, clockReceived)

			// Vector clock adjustment if message received from other controller
			vClockReceived := castStringToVClock(vClockReceivedStr)
			vClock = vClockAdjustment(vClock, vClockReceived, idVClock)
			// logInfo("main", "Clock updated, message received from other controller.")

		} else if clockReceivedStr == "" && sender == "A"+name[1:2] { // Filters out messages without a clock from the wrong app.
			// Incremented if message received from base app
			clock = clock + 1
			vClock[idVClock] = vClock[idVClock] + 1
			// logInfo("main", "Clock updated, message received from local app.")

		} else { // Filters out messages from other controller to their own app
			// ERROR, ignoring
			logError("main", "Message from another controller to it's own app (IGNORED) OR UNEXPECTED ERROR.")
			messageReceived = ""
			continue
		}

		// Message emission
		// logInfo("main", "Sending message...")
		if clockReceivedStr != "" {
			// Sending to base app
			fmt.Printf(encodeMessage([]string{"snd", "msg"}, []string{name, findValue(keyValTable, "msg")}) + "\n")
			logInfo("main", "Message sent to local app.")
		} else {
			// Sending to other controller
			fmt.Printf(encodeMessage([]string{"snd", "hlg", "vlg", "msg"}, []string{name, strconv.Itoa(clock), castVClockToString(vClock), findValue(keyValTable, "msg")}) + "\n")
			logInfo("main", "Message sent to other controller.")
		}

		messageReceived = ""
	}
}
