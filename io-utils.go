package main

import (
	"fmt"
	"math/rand"
	"time"
)

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
	for {
		select {
		case message := <-ch:
			// Wait a random time between 10 and 40 ms and print next message
			rand.Seed(time.Now().UnixNano())
			duration := time.Duration(rand.Intn(31) + 10)
			time.Sleep(duration * time.Millisecond)
			fmt.Printf(message)
		}
	}
}
