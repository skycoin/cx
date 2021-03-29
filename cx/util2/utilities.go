package util2

import (
	"fmt"
	"github.com/skycoin/cx/cx/constants"
)

// Debug ...
func Debug(args ...interface{}) {
	fmt.Println(args...)
}

// IsTempVar ...
func IsTempVar(name string) bool {
	if len(name) >= len(constants.LOCAL_PREFIX) && name[:len(constants.LOCAL_PREFIX)] == constants.LOCAL_PREFIX {
		return true
	}
	return false
}

func CheckForEscapedChars(str string) []byte {
	var res []byte
	var lenStr = int(len(str))
	for c := 0; c < len(str); c++ {
		var nextCh byte
		ch := str[c]
		if c < lenStr-1 {
			nextCh = str[c+1]
		}
		if ch == '\\' {
			switch nextCh {
			case '%':
				c++
				res = append(res, nextCh)
				continue
			case 'n':
				c++
				res = append(res, '\n')
				continue
			default:
				res = append(res, ch)
				continue
			}

		} else {
			res = append(res, ch)
		}
	}

	return res
}

