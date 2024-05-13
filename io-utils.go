package main

import (
	"fmt"
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
			fmt.Printf(message)
			time.Sleep(1 * time.Second)
		}
	}
}
