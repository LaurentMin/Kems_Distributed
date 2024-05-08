package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	// ASCII ranges in which to look for seperators (includes numbers)
	beginRangeASCII := [5]int{58, 33, 91, 123, 48}
	endRangeASCII := [5]int{64, 47, 96, 126, 57}

	// Error returns ""
	if len(beginRangeASCII) != len(endRangeASCII) {
		stderr.Printf("Incorrect ASCII range (correct function code).\n")
		return ""
	}

	// lookup loop (looks in the ascii ranges one by one if msg contains character)
	for i := 0; i < len(beginRangeASCII); i++ {
		for asciiCode := beginRangeASCII[i]; asciiCode <= endRangeASCII[i]; asciiCode++ {
			asciiVal := string(rune(asciiCode))
			if strings.Contains(msg, asciiVal) {
				continue
			}
			return asciiVal
		}
	}

	// Error returns ""
	stderr.Printf("Seperation caracter not found for %s\n", msg)
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
	// Error returns ""
	if len(keyTab) != len(valTab) {
		stderr.Printf("Wrong parity for formatting.\n")
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
	return msg
}

///////////////////////
// DECODING MESSAGES //
///////////////////////
/*
	Parses a message (received from another controller)
*/
func decodeMessage(msg string) []string {
	// Error returns empty table
	if len(msg) < 4 {
		stderr.Printf("Ivalid message for parsing %s\n", msg)
		return []string{}
	}

	// Getting seperator and returning splitted string
	sep := msg[0:1]
	// msg[1:] is to avoid that split returns a first empty element
	return strings.Split(msg[1:], sep)
}

/*
	Finds the FIRST value that matches a specific key in []string
	This function can be used only with message parsed with decodeMessage()
	Returns "" if value is "" or if error (no value found or other)
*/
func findValue(table []string, key string) string {
	// Error returns ""
	if len(table) == 0 {
		stderr.Printf("No value to find in empty table, key %s\n", key)
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
			return pair[1]
		}
	}

	// Error returns ""
	stderr.Printf("No value found for key %s\n", key)
	return ""
}

///////////
// OTHER //
///////////
/*
	Clock adjustment
*/
func recaler(x, y int) int {
	if x < y {
		return y + 1
	}
	return x + 1
}

/*
	Message logging
*/
var pid = os.Getpid()
var name = "default"
var stderr = log.New(os.Stderr, "", 0)

func main() {
	//////////////// TESTING
	/*
		
		fmt.Println(encodeMessage([]string{"key1", "key2", "key3"}, []string{"val1", "val2", "val3"}))
		test := encodeMessage([]string{"snd", "hlg", "msg"}, []string{"elouan", "23", "coucou"})
		fmt.Println(test)
		decodedTest := decodeMessage(test)
		fmt.Println(decodedTest)
		fmt.Println(findValue(decodedTest,"snd"))
	*/
	
	logMessage("hello", "world")
	logSuccess("hello", "world")
	logInfo("hello", "world")
	logWarning("hello", "world")
	logError("hello", "world")
	//////////////// BEGINNING OF PROGRAM
	

	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "ecrivain", "nom")
    flag.Parse()
	name = *pName

	// Initialising key variables for controller
	var messageReceived string
	var keyValTable []string
	var clock int = 0

	// Main loop of the controller, manages message reception and emission
	for {
		// Message reception
		fmt.Scanln(&messageReceived)
		keyValTable = decodeMessage(messageReceived)

		// Defining local clock depending on received message
		clockReceivedStr := findValue(keyValTable, "hlg")
		// Adjustment if message received from other controller
		if clockReceivedStr != "" {
			clockReceived, err := strconv.Atoi(clockReceivedStr)
			if err != nil {
				stderr.Printf("Error converting string to int: ", err)
				continue
			}
			clock = recaler(clock, clockReceived)
			// Incremented if message received from base app
		} else {
			clock = clock + 1
		}

		// Message emission
		// Sending to base app
		if clockReceivedStr != "" {
			fmt.Printf(findValue(keyValTable, "msg") + "\n")
		// Sending to other controller
		} else {
			fmt.Printf(encodeMessage([]string{"msg", "hlg"}, []string{messageReceived, strconv.Itoa(clock)}) + "\n")
		}
	}
}
