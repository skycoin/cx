package base

import (
	"fmt"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func test_error (message *CXArgument, isErrorPresent bool, expr *CXExpression) error {
	if !isErrorPresent {
		var _message string
		encoder.DeserializeRaw(*message.Value, &_message)
		if _message == "" {
			fmt.Println(fmt.Sprintf("%d: an error was expected and did not occur", expr.FileLine))
		} else {
			fmt.Println(fmt.Sprintf("%d: %s", expr.FileLine, _message))
		}
		
		return nil
	} else {
		return nil
	}
}

func test_value (result *CXArgument, expected *CXArgument, message *CXArgument, expr *CXExpression) error {
	if result.Typ != expected.Typ {
		fmt.Println(fmt.Sprintf("%d: result and expected value are not of the same type", expr.FileLine))
		return nil
	}
	
	equal := true
	var _message string
	encoder.DeserializeRaw(*message.Value, &_message)

	if len(*result.Value) != len(*expected.Value) {
		equal = false
	}

	if equal {
		for i, byt := range *result.Value {
			if byt != (*expected.Value)[i] {
				equal = false
				break
			}
		}
	}

	if !equal {
		if _message == "" {
			//fmt.Println(fmt.Sprintf("%d: result was not equal to the expected value", expr.FileLine))
			fmt.Printf("%d: result was not equal to the expected value\n", expr.FileLine)
			return errors.New("")
		} else {
			fmt.Println(fmt.Sprintf("%d: result was not equal to the expected value; %s\n", expr.FileLine, _message))
			return errors.New("")
		}
	}
	
	return nil
}
