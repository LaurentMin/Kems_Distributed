package main

import (
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
	// ASCII ranges in which to look for seperators (includes numbers)
	beginRangeASCII := [5]int{58, 33, 91, 123, 48}
	endRangeASCII := [5]int{64, 47, 96, 126, 57}

	// Error returns ""
	if len(beginRangeASCII) != len(endRangeASCII) {
		fmt.Println("Impossible to find seperation caracter (incorrect ASCII range).")
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
	fmt.Println("Seperation caracter not found.")
	return ""
}

/*
	Formats a message (before sending to other controlers)
	Works in 3 steps :
	1. Determining a key val sep for each pair
	2. Determining a global field seperator
	3. Building the global message
*/
func encodeMessage(keyTab []string, valTab []string) string {
	// Error returns ""
	if len(keyTab) != len(valTab) {
		fmt.Println("Wrong parity for formatting.")
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
		fmt.Println("Ivalid message for parsing.")
		return []string{}
	}

	// Getting seperator and returning splitted string
	sep := msg[0:1]
	// msg[1:] is to avoid that split returns an first empty element
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
		fmt.Println("No value to find in empty table.")
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
	fmt.Println("No value found.")
	return ""
}

// fonction pour recaler l'horloge
func recaler(x, y int) int {
	if x < y {
		return y + 1
	}
	return x + 1
}

func main() {

	//////////////// TEST
	// fmt.Println(encodeMessage([]string{"key1", "key2", "key3"}, []string{"val1", "val2", "val3"}))
	test := encodeMessage([]string{"snd", "hlg", "msg"}, []string{"elouan", "23", "coucou"})
	fmt.Println(test)
	decodedTest := decodeMessage(test)
	fmt.Println(decodedTest)
	fmt.Println(findValue(decodedTest,"snd"))
	////////////////

	var rcvmsg string
	var tab_keyval []string // tableau de clefvaleur
	var h int = 0           // horloge entière

	for {
		fmt.Scanln(&rcvmsg)
		//fmt.Printf("message controler : %s \n", rcvmsg)

		tab_keyval = decodeMessage(rcvmsg)

		/*
			for _, keyval := range tab_keyval {
				tab_key_val := strings.Split(keyval[1:], keyval[0:1])
				   fmt.Printf("  %q\n", tab_key_val)
				   fmt.Printf("  key : %s  val : %s\n", tab_key_val[0], tab_key_val[1])
			}
		*/
		// traitement de l'horloge
		s_hrcv := findValue(tab_keyval, "hlg")
		if s_hrcv != "" {
			hrcv, err := strconv.Atoi(s_hrcv)
			if err != nil {
				fmt.Println("Error converting string to int: ", err)
				continue
			}
			h = recaler(h, hrcv)
		} else {
			h = h + 1
		}

		// traitement du message
		if s_hrcv != "" {
			fmt.Printf(findValue(tab_keyval, "msg") + "\n")
		} else {
			fmt.Printf(encodeMessage([]string{"msg", "hlg"}, []string{rcvmsg, strconv.Itoa(h)}) + "\n")
		}
	}
}
