package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
)

///////////
// CLOCK //
///////////
/*
	Clock adjustment
*/
func clockAdjustment(x, y int) int {
	// logMessage("clockAdjustment", "Adjusting clock to max(local,received) + 1.")
	if x < y {
		return y + 1
	}
	return x + 1
}

/*
	Vector clock adjustment
*/
func vClockAdjustment(x []int, y []int, ind int) []int {
	// logMessage("vClockAdjustment", "Adjusting vector clock to max(local,received) + 1.")
	minRange := len(x)
	if len(y) < len(x) {
		minRange = len(y)
	}

	for i := 0; i < minRange; i++ {
		if x[i] < y[i] {
			x[i] = y[i]
		} else {
			x[i] = x[i]
		}
	}
	x[ind] = x[ind] + 1
	return x
}

/*
	Cast string to vector clock*
*/
func castStringToVClock(strVlg string) []int {
	strVlg = strings.Trim(strVlg, "[]")
	elements := strings.Split(strVlg, " ")

	var vlg []int

	for _, element := range elements {
		//logInfo("castStringToVClock", "Element: "+element)
		if element != "" {
			num, err := strconv.Atoi(element)
			if err != nil {
				logError("castStringToVClock", "Error converting string to int: "+err.Error())
			}
			vlg = append(vlg, num)
		}
	}

	return vlg
}

/*
	Cast vector clock to string
*/
func castVClockToString(vlg []int) string {
	var strVlg string = "["

	for _, element := range vlg {
		strVlg += strconv.Itoa(element) + " "
	}

	strVlg += "]"

	return strVlg
}

///////////
// FILES //
///////////
/*
	Save game in file
*/
func saveGame(gameSave string, name string, vClock []int) {
	logInfo("saveGame", "Saving game in file.")
	// Open file
	filename := "gameSave" + name + castVClockToString(vClock) + ".txt"
	file, err := os.Create(filename)
	if err != nil {
		logError("saveGame", "Error creating file: "+err.Error())
		return
	}
	defer file.Close()

	// Write in file
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(gameSave)
	if err != nil {
		logError("saveGame", "Error writing in file: "+err.Error())
		return
	}

	err = writer.Flush()
	if err != nil {
		logError("saveGame", "Error flushing writer: "+err.Error())
		return
	}

	logSuccess("saveGame", "Game saved in file.")
}

/*
	estampille structure
*/
type Request struct {
	Type  string
	Clock int
}

/////////////////
// ESTAMPILLES //
/////////////////
/*
	Returns true if siteN has smallest date in estampilles
	(called to start a critical section for example)
*/
func canGoCritical(estampilles []Request, site int) bool {
	for i := 0; i < len(estampilles); i++ {
		logError("canGoCritical", "num site : "+strconv.Itoa(i)+" clock : "+strconv.Itoa(estampilles[i].Clock))
		if estampilles[i].Clock > -1 && (estampilles[site].Clock > estampilles[i].Clock || (estampilles[site].Clock == estampilles[i].Clock && site > i)) {
			return false
		}
	}
	return true
}

/*
Adds a controller to estampilles table
*/
func addController(ctl string, table *[]Request, vClock *[]int) {
	if len(ctl) < 2 || ctl[0] != 'C' {
		logError("addController", "Error adding controller: "+ctl)
		return
	}

	ctlNumber, err := strconv.Atoi(ctl[1:])
	if err != nil {
		logError("addController", "Error converting string to int: "+err.Error())
	}

	for i := 0; i <= ctlNumber; i++ {
		if i >= len(*vClock) {
			*vClock = append(*vClock, 0)
		}
		if i >= len(*table) {
			*table = append(*table, Request{"[ECRITICAL]", -1})
		}
	}
}

/*
Removes a controller from estampilles table
*/
func removeController(ctl string, table *[]Request, vClock *[]int) {
	if len(ctl) < 2 || ctl[0] != 'C' {
		logError("addController", "Error removing controller: "+ctl)
		return
	}

	ctlNumber, err := strconv.Atoi(ctl[1:])
	if err != nil || ctlNumber >= len(*table) {
		logError("addController", "Controller out of range or Error converting string to int: "+err.Error())
	}

	(*table)[ctlNumber].Clock = -1
}

//////////
// MAIN //
//////////
func main() {
	// Getting name from commandline (usefull for logging)
	pName := flag.String("n", "controller", "name")
	flag.Parse()
	name = *pName

	// Starting Controller
	// logInfo("main", "Launching controller...")
	// Initialising key variables for controller
	messageReceived := ""
	keyValTable := []string{}
	clock := 0

	vClock := []int{}
	// Save state
	saveState := false

	//	estampilles := []Request{Request{"[ECRITICAL]", 0}, Request{"[ECRITICAL]", 0}, Request{"[ECRITICAL]", 0}} // Index 0..2 corresponds to controllers 1..3
	estampilles := []Request{}
	siteNum, _ := strconv.Atoi(name[1:2]) // Ok if this makes app crash (name must be defined)

	// Find the controller number in vClock
	idVClock := siteNum

	// Go routines to read and write input / output
	inChan := make(chan string, 10)
	outChan := make(chan string, 10)
	go read(inChan)
	go write(outChan)

	// init vclock and estampilles
	for i := 0; i <= siteNum; i++ {
		estampilles = append(estampilles, Request{"[ECRITICAL]", -1})
		vClock = append(vClock, 0)
	}

	// Main loop of the controller, manages message reception and emission and processing
	for {
		// logInfo("main", "Waiting for message.")
		// Message reception
		messageReceived = <-inChan
		logInfo("main", "Message received "+messageReceived)

		// Determine message type for processing
		keyValTable = decodeMessage(messageReceived)
		sender := findValue(keyValTable, "snd")

		// Filter out random messages (Display messages for example)
		if len(sender) != 2 || len(name) != 2 || (sender != "A"+name[1:2] && sender[:1] != "C") {
			logWarning("main", "Display message OR invalid sender OR wrong ctl name (ignored) - CAN BE FATAL!")
			messageReceived = ""
			continue
		}

		// Check if message from net to add or delete nodes from mutual exclusion table
		switch findValue(keyValTable, "msg") {
		case "new":
			// add node to mutual exclusion table
			addController(sender, &estampilles, &vClock)
			numPlayers := strconv.Itoa(len(estampilles))
			outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[UPDATEPLA]" + numPlayers}) + "\n"
			logInfo("main", "Controller was added to table, num players sent to app : "+sender)
		case "del":
			// remove node from mutual exclusion table
			removeController(sender, &estampilles, &vClock)
			logInfo("main", "Controller was removed from table : "+sender)
		}

		// Clock updating
		// logInfo("main", "Clock updating...")
		clockReceivedStr := findValue(keyValTable, "hlg")
		vClockReceivedStr := findValue(keyValTable, "vlg")
		if clockReceivedStr != "" && vClockReceivedStr != "" && sender[:1] == "C" { // Filters out messages from an app with a clock (should never happen)

			// Clock adjustment if message received from other controller
			clockReceived, err := strconv.Atoi(clockReceivedStr)
			if err != nil {
				logError("main", "Error converting string to int : "+err.Error()+" (FATAT, clock corruption)")
				continue
			}
			clock = clockAdjustment(clock, clockReceived)

			// Vector clock adjustment if message received from other controller
			vClockReceived := castStringToVClock(vClockReceivedStr)
			vClock = vClockAdjustment(vClock, vClockReceived, idVClock)
			// logInfo("main", "Clock updated, message received from other controller.")

		} else if clockReceivedStr == "" && sender == "A"+name[1:2] { // Filters out messages without a clock from the wrong app or a controller.
			// Incremented if message received from app
			clock = clock + 1
			vClock[idVClock] = vClock[idVClock] + 1
			// logInfo("main", "Clock updated, message received from local app.")

		} else { // Filters out messages from other controller to their own app or other errors
			// ERROR, ignoring
			logWarning("main", "Message from another controller to it's own app (IGNORED) OR UNEXPECTED ERROR.")
			messageReceived = ""
			continue
		}

		// Message processing
		// getting message
		messageReceived = findValue(keyValTable, "msg")
		// Filter out wrong messages
		if len(messageReceived) < 11 {
			// logInfo("main", "Wrong message type for app received "+messageReceived+" (ignoring).")
			logInfo("main", "Wrong message type for app received (too short) (ignoring).")
			messageReceived = ""
			continue
		}

		// Receive from controller
		// logInfo("main", "Sending message...")
		if clockReceivedStr != "" && sender[:1] == "C" {
			otherSiteNumber, _ := strconv.Atoi(sender[1:2])
			switch messageReceived[:11] {
			case "[GAMESTATE]":
				// Do not replace an ask by a gamestate
				if estampilles[otherSiteNumber].Type != "[ACRITICAL]" {
					estampilles[otherSiteNumber].Type = "[GAMESTATE]"
					estampilles[otherSiteNumber].Clock = clock
				}
				// Check if can start own critical
				if estampilles[siteNum].Type == "[ACRITICAL]" && canGoCritical(estampilles, siteNum) {
					outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[BCRITICAL]"}) + "\n"
					logInfo("main", "Begin critical section sent to base app.")
				}
				// Send gamestate to app
				outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, messageReceived}) + "\n"
				logInfo("main", "Gamestate message sent to local app.")

			case "[SAVEORDER]":
				// Do not replace an ask by a reception
				if estampilles[otherSiteNumber].Type != "[ACRITICAL]" {
					estampilles[otherSiteNumber].Type = "[SAVEORDER]"
					estampilles[otherSiteNumber].Clock = clock
				}

				// Check if can start own critical
				if estampilles[siteNum].Type == "[ACRITICAL]" && canGoCritical(estampilles, siteNum) {
					outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[BCRITICAL]"}) + "\n"
					logInfo("main", "Begin critical section sent to base app.")
				}

				// Save order received from other controller
				if strconv.FormatBool(saveState) != messageReceived[11:] {
					outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, messageReceived}) + "\n"
					// outChan <- encodeMessage([]string{"snd", "hlg", "vlg", "msg"}, []string{name, strconv.Itoa(clock), castVClockToString(vClock), messageReceived}) + "\n"
					saveState = !saveState
					logInfo("main", "Save order received from other controller and send to local app.")
				}

			case "[ACRITICAL]": // Other controller asks for access restriction
				estampilles[otherSiteNumber].Type = "[ACRITICAL]"
				estampilles[otherSiteNumber].Clock = clock
				outChan <- encodeMessage([]string{"snd", "hlg", "vlg", "msg"}, []string{name, strconv.Itoa(clock), castVClockToString(vClock), "[VCRITICAL]" + sender}) + "\n"
				logInfo("main", "Answered to other controller restriction access demand.")
				// Check if can start own critical
				if estampilles[siteNum].Type == "[ACRITICAL]" && canGoCritical(estampilles, siteNum) {
					outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[BCRITICAL]"}) + "\n"
					logInfo("main", "Begin critical section sent to base app.")
				}

			case "[VCRITICAL]": // Other controller validates request reception
				// Reject validations that are not meant for this controller
				if messageReceived[11:] != name {
					logWarning("main", "Validation not for this controller (ignored).")
					messageReceived = ""
					continue
				}

				// Do not replace an ask by a reception
				if estampilles[otherSiteNumber].Type != "[ACRITICAL]" {
					estampilles[otherSiteNumber].Type = "[VCRITICAL]"
					estampilles[otherSiteNumber].Clock = clock
				}
				logInfo("main", "Critical section reception was confirmed by other controller.")
				// Check if can start own critical
				if estampilles[siteNum].Type == "[ACRITICAL]" && canGoCritical(estampilles, siteNum) {
					outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[BCRITICAL]"}) + "\n"
					logInfo("main", "Begin critical section sent to base app.")
				}

			case "[ECRITICAL]": // Other controller liberates access restriction
				estampilles[otherSiteNumber].Type = "[ECRITICAL]"
				estampilles[otherSiteNumber].Clock = clock
				logInfo("main", "Other controller ended restriction access.")
				// Check if can start own critical
				if estampilles[siteNum].Type == "[ACRITICAL]" && canGoCritical(estampilles, siteNum) {
					outChan <- encodeMessage([]string{"snd", "msg"}, []string{name, "[BCRITICAL]"}) + "\n"
					logInfo("main", "Begin critical section sent to base app.")
				}

			default:
				logError("main", "Wrong message type received (controller sent wrong type) (ignoring) (could be critical for clock).")
			}

			messageReceived = ""
			continue
		}

		// Received from app
		if clockReceivedStr == "" && sender == "A"+name[1:2] {
			switch messageReceived[:11] {
			case "[GAMESTATE]":
				// Do not replace an ask by a gamestate
				if estampilles[siteNum].Type != "[ACRITICAL]" {
					estampilles[siteNum].Type = "[GAMESTATE]"
					estampilles[siteNum].Clock = clock
				}
				outChan <- encodeMessage([]string{"snd", "hlg", "vlg", "msg"}, []string{name, strconv.Itoa(clock), castVClockToString(vClock), messageReceived}) + "\n"
				logInfo("main", "Gamestate message sent to other controller.")

			case "[SAVEORDER]":
				// Save order received from base app
				// Do not replace an ask by a gamestate
				if estampilles[siteNum].Type != "[ACRITICAL]" {
					estampilles[siteNum].Type = "[SAVEORDER]"
					estampilles[siteNum].Clock = clock
				}

				if findValue(keyValTable, "saveOrder") == "1" {
					gameSave := findValue(keyValTable, "msg")
					gameSave = gameSave[11:]
					// logInfo("main", "Order saved: "+gameSave)
					// logInfo("main", "Save local game.")

					// Save game in file
					saveGame(gameSave, name, vClock)
					saveState = !saveState

					logInfo("main", "Sent save order to other controllers.")
					outChan <- encodeMessage([]string{"snd", "hlg", "vlg", "msg"}, []string{name, strconv.Itoa(clock), castVClockToString(vClock), "[SAVEORDER]" + strconv.FormatBool(saveState)}) + "\n"

				} else if findValue(keyValTable, "saveOrder") == "0" {
					// made save
					gameSave := findValue(keyValTable, "msg")
					gameSave = gameSave[11:]

					// Save game in file
					saveGame(gameSave, name, vClock)
				}

			case "[ACRITICAL]": // Base app asks critical (asking other controllers)
				estampilles[siteNum].Type = "[ACRITICAL]"
				estampilles[siteNum].Clock = clock
				outChan <- encodeMessage([]string{"snd", "hlg", "vlg", "msg"}, []string{name, strconv.Itoa(clock), castVClockToString(vClock), "[ACRITICAL]"}) + "\n"
				logInfo("main", "Asked other controllers for access restriction.")

			case "[ECRITICAL]": // Base app stops critical (liberating other controllers)
				estampilles[siteNum].Type = "[ECRITICAL]"
				estampilles[siteNum].Clock = clock
				outChan <- encodeMessage([]string{"snd", "hlg", "vlg", "msg"}, []string{name, strconv.Itoa(clock), castVClockToString(vClock), "[ECRITICAL]"}) + "\n"
				logInfo("main", "Liberated other controllers from access restriction.")

			default:
				logError("main", "Wrong message type received (app sent wrong type) (ignoring) (could be critical for clock).")
			}

			messageReceived = ""
			continue
		}

		logError("main", "CRITICAL ERROR, MESSAGE TREATMENT WAS NOT IMPLEMENTED (should never happen) (bad for clock)")
		messageReceived = ""
	}
}
