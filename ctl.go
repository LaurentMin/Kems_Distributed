package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

///////////////////////
// ENCODING MESSAGES //
///////////////////////
/*
	Returns the first ASCII "seperator" character non present in the string received as a parameter
	Letters are not included (lookup order can be modified with beginRangeASCII and endRangeASCII)
*/
func determineSep(msg string) string {
	logMessage("determineSep", "Determining seperator for "+msg)
	// ASCII ranges in which to look for seperators (includes numbers)
	beginRangeASCII := [5]int{58, 33, 91, 123, 48}
	endRangeASCII := [5]int{64, 47, 96, 126, 57}

	// Error returns ""
	if len(beginRangeASCII) != len(endRangeASCII) {
		logError("determineSep", "Incorrect ASCII range (correct function code).")
		return ""
	}

	// lookup loop (looks in the ascii ranges one by one if msg contains character)
	for i := 0; i < len(beginRangeASCII); i++ {
		for asciiCode := beginRangeASCII[i]; asciiCode <= endRangeASCII[i]; asciiCode++ {
			asciiVal := string(rune(asciiCode))
			if strings.Contains(msg, asciiVal) {
				continue
			}
			logSuccess("determineSep", asciiVal+" found as a seperator for "+msg)
			return asciiVal
		}
	}

	// Error returns ""
	logError("determineSep", "Seperation caracter not found for "+msg)
	return ""
}

/*
	Formats a message (before sending to other controllers)
	Works in 3 steps :
	1. Determining a key val sep for each pair
	2. Determining a global field seperator
	3. Building the global message
*/
func encodeMessage(keyTab []string, valTab []string) string {
	logMessage("encodeMessage", "Encoding message with "+strconv.Itoa(len(keyTab))+" key value pairs.")
	// Error returns ""
	if len(keyTab) != len(valTab) {
		logError("encodeMessage", "Wrong parity for formatting.")
		return ""
	}

	// 1.
	// Formatting each key value pair in the tables
	for i := 0; i < len(keyTab); i++ {
		pairSep := determineSep(keyTab[i] + valTab[i])
		// Error occurs returns ""
		if pairSep == "" {
			return ""
		}
		// Updating values with seperator
		keyTab[i] = pairSep + keyTab[i]
		valTab[i] = pairSep + valTab[i]
	}

	// 2.
	// Getting the field sep
	tmp := ""
	for i := 0; i < len(keyTab); i++ {
		tmp += keyTab[i] + valTab[i]
	}
	fieldSep := determineSep(tmp)
	// Error occurs returns ""
	if fieldSep == "" {
		return ""
	}

	// 3.
	// Formatting the full message with all key val pairs and field sep
	msg := ""
	for i := 0; i < len(keyTab); i++ {
		msg += fieldSep + keyTab[i] + valTab[i]
	}
	logSuccess("encodeMessage", "Message formatted correctly : "+msg)
	return msg
}

///////////////////////
// DECODING MESSAGES //
///////////////////////
/*
	Parses a message (received from another controller)
*/
func decodeMessage(msg string) []string {
	logMessage("decodeMessage", "Parsing : "+msg)
	// Error returns empty table
	if len(msg) < 4 {
		logWarning("decodeMessage", "Message too short for parsing : "+msg)
		return []string{}
	}

	// Getting seperator and returning splitted string
	sep := msg[0:1]
	// msg[1:] is to avoid that split returns a first empty element
	logSuccess("decodeMessage", msg+" parsed with seperator "+sep)
	return strings.Split(msg[1:], sep)
}

/*
	Finds the FIRST value that matches a specific key in []string
	This function can be used only with message parsed with decodeMessage()
	Returns "" if value is "" or if error (no value found or other)
*/
func findValue(table []string, key string) string {
	logMessage("findValue", "Finding value of key "+key)
	// Error returns ""
	if len(table) == 0 {
		logWarning("findValue", "No value to find in empty table, key : "+key)
		return ""
	}

	// Loop on the table to find key
	for i := 0; i < len(table); i++ {
		pair := decodeMessage(table[i])
		// Invalid pair, goes to next pair
		if len(pair) == 0 {
			continue
		}

		// Trying to match key
		if pair[0] == key {
			logSuccess("findValue", pair[1]+" value found for key "+key)
			return pair[1]
		}
	}

	// Error returns ""
	logMessage("findValue", "No value found for key : "+key)
	return ""
}

///////////
// OTHER //
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

	// Main loop of the controller, manages message reception and emission as well as treatment
	for {
		logInfo("main", "Waiting for message.")
		// Message reception
		fmt.Scanln(&messageReceived)
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
