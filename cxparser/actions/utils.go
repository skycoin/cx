package actions

import "strings"

func StripNameNumber(name string) string {
	index := strings.Index(name, ":")
	if index == -1 {
		return name
	}

	return name[index+1:]
}
