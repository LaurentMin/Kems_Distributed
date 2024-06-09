package main

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
func getDiffIdIndexOrCreateIfNotExists(table *[]Diffusion, id string, numNeighbours int) int {
	for i, diff := range *table {
		if diff.diffIndex == id {
			return i
		}
	}
	newDiff := getDiffusioni(id, numNeighbours)
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

func printDiffusion(diff Diffusion) string {
	return diff.diffIndex + "|" + string(diff.color) + "|" + string(diff.nbNeighbours) + "|" + diff.parent
}
