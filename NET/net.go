package main

import (
	"fmt"
	"strings"
	"flag"
)

func main () {
	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "net", "name")
	pAskNode := flag.String("a", "N1", "name of other node")
	flag.Parse()

	name := *pName
	fmt.Println("Starting ", name)

	askNode := *pAskNode
	fmt.Println("Ask to ", askNode)

	if askNode {
		// TODO : executer ./reseauNet.sh name askNode
	
	}

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
