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
	Formats a message before sending it to other ctls
	Works in 3 steps :
	1. Determining a key val sep for each pair
	2. Determining a global field seperator
	3. Building the global message
*/
func formatMessage(keyTab []string, valTab []string) string {
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

// fonctions pour parser un message reçu avec les différentes valeurs d'une autre application
func parse_keyval(msg string) []string {
	if len(msg) < 4 {
		return []string{msg}
	}

	sep := msg[0:1]

	return strings.Split(msg, sep)
}

// fonction pour trouver une valeur dans un tableau de clefvaleur
func findval(tab_keyval []string, key string) string {

	var val string = ""

	for _, keyval := range tab_keyval {
		if len(keyval) < 4 {
			continue
		}

		tab_key_val := strings.Split(keyval[1:], keyval[0:1])
		if tab_key_val[0] == key {
			//val = return val
		}
	}
	return val
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
	fmt.Println(formatMessage([]string{"snd", "hlg", "msg"}, []string{"elouan", "23", "coucou"}))

	////////////////

	var rcvmsg string
	var tab_keyval []string // tableau de clefvaleur
	var h int = 0           // horloge entière

	for {
		fmt.Scanln(&rcvmsg)
		//fmt.Printf("message controler : %s \n", rcvmsg)

		tab_keyval = parse_keyval(rcvmsg)

		/*
			for _, keyval := range tab_keyval {
				tab_key_val := strings.Split(keyval[1:], keyval[0:1])
				   fmt.Printf("  %q\n", tab_key_val)
				   fmt.Printf("  key : %s  val : %s\n", tab_key_val[0], tab_key_val[1])
			}
		*/
		// traitement de l'horloge
		s_hrcv := findval(tab_keyval, "hlg")
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
			fmt.Printf(findval(tab_keyval, "msg") + "\n")
		} else {
			fmt.Printf(formatMessage([]string{"msg", "hlg"}, []string{rcvmsg, strconv.Itoa(h)}) + "\n")
		}
	}
}
