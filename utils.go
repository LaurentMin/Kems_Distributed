package main

import (
	"bufio"
	"os"
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
	// logMessage("determineSep", "Determining seperator for "+msg)
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
			// logSuccess("determineSep", asciiVal+" found as a seperator for "+msg)
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
	// logMessage("encodeMessage", "Encoding message with "+strconv.Itoa(len(keyTab))+" key value pairs.")
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
	// logSuccess("encodeMessage", "Message formatted correctly : "+msg)
	return msg
}

///////////////////////
// DECODING MESSAGES //
///////////////////////
/*
	Parses a message (received from another controller)
*/
func decodeMessage(msg string) []string {
	// logMessage("decodeMessage", "Parsing : "+msg)
	// Error returns empty table
	if len(msg) < 4 {
		// logWarning("decodeMessage", "Message too short for parsing : "+msg)
		return []string{}
	}

	// Getting seperator and returning splitted string
	sep := msg[0:1]
	// msg[1:] is to avoid that split returns a first empty element
	// logSuccess("decodeMessage", msg+" parsed with seperator "+sep)
	return strings.Split(msg[1:], sep)
}

/*
	Finds the FIRST value that matches a specific key in []string
	This function can be used only with message parsed with decodeMessage()
	Returns "" if value is "" or if error (no value found or other)
*/
func findValue(table []string, key string) string {
	// logMessage("findValue", "Finding value of key "+key)
	// Error returns ""
	if len(table) == 0 {
		// logWarning("findValue", "No value to find in empty table, key : "+key)
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
			// logSuccess("findValue", pair[1]+" value found for key "+key)
			return pair[1]
		}
	}

	// Error returns ""
	// logMessage("findValue", "No value found for key : "+key)
	return ""
}

///////////////////
// READING STDIN //
///////////////////
/*
	Scans stdin until the first newline '\n' character is reached
*/
func scanUntilNewline() string {
	reader := bufio.NewReader(os.Stdin)
	input := ""

	for {
		line, err := reader.ReadString('\n')
		// ERROR returns what could be read
		if err != nil {
			logError("scanUntilNewline:", "Error while using ReadString('\n')"+err.Error())
			return input
		}
		input += line

		if strings.Contains(line, "\n") {
			break
		}
	}

	return strings.ReplaceAll(input, "\n", "")
}
