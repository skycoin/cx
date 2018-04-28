package base

import (
	"os"
	"errors"
	"fmt"
	"io/ioutil"
	//"path/filepath"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var openFiles map[string]*os.File = make(map[string]*os.File, 0)

func os_ReadFile (fileName *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("os.ReadFile", "str", fileName); err == nil {
		var name string
		encoder.DeserializeRaw(*fileName.Value, &name)
		
		if byts, err := ioutil.ReadFile(name); err == nil {
			sByts := encoder.Serialize(byts)
			assignOutput(0, sByts, "[]byte", expr, call)
		} else {
			return err
		}
		return nil
	} else {
		return err
	}
}

func os_Create (fileName *CXArgument) error {
	if err := checkType("os.Create", "str", fileName); err == nil {
		var name string
		encoder.DeserializeRaw(*fileName.Value, &name)
		
		if file, err := os.Create(name); err == nil {
			openFiles[name] = file
		} else {
			return err
		}

		return nil
	} else {
		return err
	}
}

func os_Write (fileName, byts *CXArgument) error {
	if err := checkTwoTypes("os.Write", "str", "[]byte", fileName, byts); err == nil {
		var _fileName string
		var _byts []byte
		encoder.DeserializeRaw(*fileName.Value, &_fileName)
		encoder.DeserializeRaw(*byts.Value, &_byts)

		if file, ok := openFiles[_fileName]; ok {
			if _, err := file.Write(_byts); err == nil {
				return nil
			} else {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("file '%s' is not currently open", _fileName))
		}
	} else {
		return err
	}
}

func os_WriteFile (fileName, byts *CXArgument) error {
	if err := checkTwoTypes("os.WriteFile", "str", "[]byte", fileName, byts); err == nil {
		var _fileName string
		var _byts []byte
		encoder.DeserializeRaw(*fileName.Value, &_fileName)
		encoder.DeserializeRaw(*byts.Value, &_byts)

		if err := ioutil.WriteFile(_fileName, _byts, os.FileMode(644)); err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func os_Open (fileName *CXArgument) error {
	if err := checkType("os.Open", "str", fileName); err == nil {
		var name string
		encoder.DeserializeRaw(*fileName.Value, &name)
		
		if file, err := os.Open(name); err == nil {
			openFiles[name] = file
		} else {
			return err
		}

		return nil
	} else {
		return err
	}
}

func os_Close (fileName *CXArgument) error {
	if err := checkType("os.Close", "str", fileName); err == nil {
		var name string
		encoder.DeserializeRaw(*fileName.Value, &name)
		
		if file, ok := openFiles[name]; ok {
			if err := file.Close(); err != nil {
				return err
			}
		}
		return nil
	} else {
		return err
	}
}

func os_GetWorkingDirectory (expr *CXExpression, call *CXCall) error {
	path := encoder.Serialize(call.Program.Path)
	assignOutput(0, path, "str", expr, call)
	return nil
}
