package main

import (
	"log"
	"os"
)

/*
	Message logging
*/
var stderr = log.New(os.Stderr, "", 0)
var pid = os.Getpid()
var name = "default"

/*
	Terminal color codes
*/
var red string = "\033[1;31m"
var orange string = "\033[1;33m"
var green string = "\033[1;32m"
var blue string = "\033[1;34m"
var reset string = "\033[0;00m"

/*
	Logging functions (they all print to stderr with program details)
*/
func logMessage(where string, what string) {
	stderr.Printf(" + [%-10.10s %d] %-16.16s : %s\n", name, pid, where, what)
}

func logSuccess(where string, what string) {
	stderr.Printf("%s # [%-10.10s %d] %-16.16s : %s\n%s", green, name, pid, where, what, reset)
}

func logInfo(where string, what string) {
	stderr.Printf("%s ? [%-10.10s %d] %-16.16s : %s\n%s", blue, name, pid, where, what, reset)
}

func logWarning(where string, what string) {
	stderr.Printf("%s * [%-10.10s %d] %-16.16s : %s\n%s", orange, name, pid, where, what, reset)
}

func logError(where string, what string) {
	stderr.Printf("%s ! [%-10.10s %d] %-16.16s : %s\n%s", red, name, pid, where, what, reset)
}
