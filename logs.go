package main

import (
	"flag"
	"log"
	"os"
)

// Terminal codes
var red string = "\033[1;31m"
var orange string = "\033[1;33m"
var green string = "\033[1;32m"
var blue string = "\033[1;34m"
var reset string = "\033[0;00m"


var pid = os.Getpid()
var name = ""
var stderr = log.New(os.Stderr, "", 0)


func logMessage(where string, what string) {
    stderr.Printf(" + [%.6s %d] %-8.8s : %s\n", name, pid, where, what)
}

func logWarning(where string, what string) {

    stderr.Printf("%s * [%.6s %d] %-8.8s : %s\n%s", orange, name, pid, where, what, reset)
}

func logError(where string, what string) {
    stderr.Printf("%s ! [%.6s %d] %-8.8s : %s\n%s", red, name, pid, where, what, reset)
}

func main() {
    pName := flag.String("n", "test", "nom")
    name = *pName
	flag.Parse()
    
	logMessage("1234567891011121314", "123456789101112131412345678910111213141234567891011121314123456789101112131412345678910111213141234567891011121314123456789101112131412345678910111213141234567891011121314123456789101112131412345678910111213141234567891011121314")
	logWarning("hello", "worddddld")
	logError("hello", "world")
}