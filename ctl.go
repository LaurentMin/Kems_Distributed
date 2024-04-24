package main

import (
	"fmt"
	"strings"
	"strconv"
)


var fieldsep = "/"
var keyvalsep = "="

// fonction pour formater un message destiné à une autre application
func msg_format(key string, val string) string {
    return fieldsep + keyvalsep + key + keyvalsep + val
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
            val = tab_key_val[1]
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
	var rcvmsg string
	var tab_keyval []string // tableau de clefvaleur
	var h int = 0 // horloge entière

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
			hrcv, _ := strconv.Atoi(s_hrcv)
			h = recaler(h, hrcv)
		} else {
			h = h + 1
		}
	
		// traitement du message 
		if  s_hrcv != "" {
			fmt.Printf(findval(tab_keyval, "msg"))
		} else {
			fmt.Printf(msg_format("msg", rcvmsg) + msg_format("hlg", strconv.Itoa(h)))
		}
    }
}