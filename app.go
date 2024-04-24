package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func sendperiodic() {
	var sndmsg string
	var i int

	i = 0

	for {
		mutex.Lock()
		i = i + 1
		sndmsg = "message_" + strconv.Itoa(i) + "\n"
		fmt.Print(sndmsg)
		mutex.Unlock()
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func receive() {
	var rcvmsg string
	l := log.New(os.Stderr, "", 0)

	for {
		fmt.Scanln(&rcvmsg)
		mutex.Lock()
		l.Println("reception <", rcvmsg, ">")
		/*
		for i := 1; i < 6; i++ {
			l.Println("traitement message", i)
			time.Sleep(time.Duration(1) * time.Second)
		}
		*/
		mutex.Unlock()
		rcvmsg = ""
	}
}

var mutex = &sync.Mutex{}

func main() {

	go sendperiodic()
	go receive()
	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
}
