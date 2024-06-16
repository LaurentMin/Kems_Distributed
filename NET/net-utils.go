package main

import (
	"strconv"
	"strings"
)

func diffusionToString(diff DiffusionMessage) string {
	str := "[DIFFUSION]"
	str += encodeMessage([]string{"diffIndex", "color", "value"}, []string{diff.diffIndex, string(diff.color), diff.value})
	return str
}

func stringToDiffusion(diffString string) DiffusionMessage {
	if diffString[:11] != "[DIFFUSION]" {
		logError("stringToDiffusion", "String is not a diffusion "+diffString+" (can be fatal for program).")
		return DiffusionMessage{}
	}
	diffString = diffString[11:]
	tabString := decodeMessage(diffString)
	diff := DiffusionMessage{}
	diff.diffIndex = findValue(tabString, "diffIndex")
	diff.color = Color(findValue(tabString, "color"))
	diff.value = findValue(tabString, "value")
	return diff
}

/*
Returns -1 if no value found
*/
func getDiffIdIndexOrCreateIfNotExists(table *[]Diffusion, id string, numNeighbours int, val string) int {
	for i, diff := range *table {
		if diff.diffIndex == id {
			return i
		}
	}
	newDiff := getDiffusioni(id, numNeighbours, val)
	*table = append(*table, newDiff)
	return len(*table) - 1
}

/*
Custom append to add string in array only if does not exist
*/
func addNeighbour(neighbours *[]string, addMe string) {
	for _, id := range *neighbours {
		if id == addMe {
			return
		}
	}
	*neighbours = append(*neighbours, addMe)
}

func deleteNeighbour(neighbours *[]string, deleteMe string) {
	for i, id := range *neighbours {
		if id == deleteMe {
			*neighbours = append((*neighbours)[:i], (*neighbours)[i+1:]...)
			return
		}
	}
}

func printDiffusion(diff Diffusion) string {
	return diff.diffIndex + "|" + string(diff.color) + "|" + strconv.Itoa(diff.nbNeighbours) + "|" + diff.parent + "|" + diff.value
}

func printDiffusionMessage(diff DiffusionMessage) string {
	return diff.diffIndex + "|" + string(diff.color) + "|" + diff.value
}

/*
Checks if node can start a wave for election
*/
func canParticipateToElection(tab []Diffusion) bool {
	for i := len(tab) - 1; i >= 0; i-- {
		if tab[i].value == "new" { // Last election encountered is finished => no election yet started
			return true
		} else if len(tab[i].value) > 1 && tab[i].value[:1] == "N" { // Last election encoutered is still ongoing (note : node name is passed as value of diffusion for elections)
			return false
		}
	}
	return true
}

/*
Checks if an election wave must be stoped
*/
func stopElecWave(tab []Diffusion, diff DiffusionMessage) bool {
	for i := len(tab) - 1; i >= 0; i-- {
		// logError("stopElecWave", printDiffusion(tab[i]))
		if tab[i].value == "new" || diff.value == "new" { // Last election is finished => no election yet started
			return false
		} else if len(tab[i].value) > 1 && tab[i].value[:1] == "N" { // Election is still ongoing
			currElec, err1 := strconv.Atoi(tab[i].value[1:])
			currWave, err2 := strconv.Atoi(diff.value[1:])
			if err1 != nil {
				logError("stopElecWave", "FATAL ERROR while converting currElec to int, should never happen.")
				return true
			}
			if err2 != nil {
				logError("stopElecWave", "FATAL ERROR while converting currWave to int, should never happen.")
				return true
			}
			if currWave > currElec {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func isDiffCtlMsg(value string) bool {
	return (value != "new" && value != "del") || value[0] == 'N'
}

func getOriginIndex(diffIndex string) string {
	if len(diffIndex) < 4 {
		logError("getOriginIndex", "FATAL, could not get index from a diffusion index !!")
		return ""
	}
	diffIndex = diffIndex[1:]
	index := strings.IndexRune(diffIndex, 'D')
	if index == -1 {
		logError("getOriginIndex", "FATAL, could not get index from a diffusion index !!")
		return diffIndex
	}
	return diffIndex[:index]
}
