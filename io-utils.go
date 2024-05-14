package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
	Setting chans as global so functions can ouput
	(needed for saving message)
*/
var inChan chan string
var outChan chan string

/*
	Reading go routine (sends read data from sdtin through channel)
*/
func read(ch chan<- string) {
	for {
		ch <- scanUntilNewline()
	}
}

/*
	Writing go routine (writes data from channel to stdout)
	Waits 1 sec before after printing
*/
func write(ch <-chan string) {
	// Init seed on program startup (programs must not be started at the same time)
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case message := <-ch:
			// Wait a random time between 10 and 40 ms and print next message
			duration := time.Duration(rand.Intn(51) + 50)
			time.Sleep(duration * time.Millisecond)
			fmt.Printf(message)
		}
	}
}
