package base

import (
	//"errors"
	//"os"
)

func os_Open(filename *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("os.Open", "str", filename); err == nil {
		//file := string(*filename.Value)
		
		// if file, err := os.Open(file); err == nil {
			
		// }

		return nil
	} else {
		return err
	}
}
