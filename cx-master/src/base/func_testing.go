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
			fmt.Println(fmt.Sprintf("%s: %d: an error was expected and did not occur", expr.FileName, expr.FileLine))
		} else {
			fmt.Println(fmt.Sprintf("%s: %d: %s", expr.FileName, expr.FileLine, _message))
		}
		
		return nil
	} else {
		return nil
	}
}

func test_value (result *CXArgument, expected *CXArgument, message *CXArgument, expr *CXExpression) error {
	if result.Type != expected.Type {
		fmt.Println(fmt.Sprintf("%s: %d: result and expected value are not of the same type", expr.FileName, expr.FileLine))
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
			fmt.Printf("%s: %d: result was not equal to the expected value\n", expr.FileName, expr.FileLine)
			return errors.New("")
		} else {
			fmt.Println(fmt.Sprintf("%s: %d: result was not equal to the expected value; %s\n", expr.FileName, expr.FileLine, _message))
			return errors.New("")
		}
	}
	
	return nil
}
