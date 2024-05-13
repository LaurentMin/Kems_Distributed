package main

import (
	"fmt"
	"time"
)

/*
	Reading go routine (sends read data from sdtin through channel)
*/
func read(ch chan<- string) {
	ch <- scanUntilNewline()
}

/*
	Writing go routine (writes data from channel to stdout)
	Waits 1 sec before after printing
*/
func write(ch <-chan string) {
	for {
		fmt.Printf(<-ch)
		time.Sleep(1 * time.Second)
	}
}
