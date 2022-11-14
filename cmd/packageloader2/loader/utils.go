package loader

func Contains(list []string, element string) bool {
	for _, elem := range list {
		if elem == element {
			return true
		}
	}
	return false
}

// Removes Duplicates from list
func removeDuplicates(list []string) []string {
	var newList []string
	for _, elem := range list {
		if !Contains(newList, elem) {
			newList = append(newList, elem)
		}
	}
	return newList
}
