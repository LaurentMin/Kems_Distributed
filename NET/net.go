package main

import (
	"fmt"
	"strings"
	"flag"
)

func main () {
	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "net", "name")
	flag.Parse()
	name := *pName

	messageReceived := ""
	inChan = make(chan string, 10)
	outChan = make(chan string, 10)
	// Reading go routine (sends read data from sdtin through channel)
	go read(inChan)

	// Writing go routine (writes data from channel to stdout)
	go write(outChan)

	for {
		messageReceived = <- inChan
		fmt.Println("Message received ", name, " : ", messageReceived)
		if ! strings.Contains(messageReceived, "test") {
			outChan <- encodeMessage([]string{"typ", "msg"}, []string{"test", messageReceived})
		}
	}
}
