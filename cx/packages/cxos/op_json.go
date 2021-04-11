// +build cxos

package cxos

import (
	"bufio"
	"encoding/json"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/util"
	"io"
	"os"
)

const (
	JSON_TOKEN_INVALID = -1
	JSON_TOKEN_NULL    = iota
	JSON_TOKEN_DELIM
	JSON_TOKEN_BOOL
	JSON_TOKEN_F64
	JSON_TOKEN_NUMBER
	JSON_TOKEN_STR
	JSON_DELIM_SQUARE_LEFT  = 91
	JSON_DELIM_SQUARE_RIGHT = 93
	JSON_DELIM_CURLY_LEFT   = 123
	JSON_DELIM_CURLY_RIGHT  = 125
)

type JSONFile struct {
	file        *os.File
	reader      *bufio.Reader
	decoder     *json.Decoder
	token       interface{}
	tokenType   int32
	tokenDelim  json.Delim
	tokenBool   bool
	tokenF64    float64
	tokenNumber json.Number
	tokenStr    string
}

var jsons []JSONFile
var freeJsons []int32

// Open the named json file for reading, returns an i32 identifying the json cxgo.
func opJsonOpen(inputs []ast.CXValue, outputs []ast.CXValue) {
	handle := int32(-1)

	file, err := util.CXOpenFile(inputs[0].Get_str())
	if err == nil {
		freeCount := len(freeJsons)
		if freeCount > 0 {
			freeCount--
			handle = int32(freeJsons[freeCount])
			freeJsons = freeJsons[:freeCount]
		} else {
			handle = int32(len(jsons))
			jsons = append(jsons, JSONFile{})
		}

		if handle < 0 || handle >= int32(len(jsons)) {
			panic("internal error")
		}

		var jsonFile JSONFile
		jsonFile.file = file
		jsonFile.reader = bufio.NewReader(file)
		jsonFile.decoder = json.NewDecoder(jsonFile.reader)
		jsonFile.decoder.UseNumber()

		jsons[handle] = jsonFile
	}

	outputs[0].Set_i32(int32(handle))
}

// Close json cxgo (and all underlying resources) idendified by it's i32 handle.
func opJsonClose(inputs []ast.CXValue, outputs []ast.CXValue) {
	success := false
    handle := inputs[0].Get_i32()
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		if err := jsonFile.file.Close(); err != nil {
			panic(err)
		}

		jsons[handle] = JSONFile{}
		freeJsons = append(freeJsons, handle)
		success = true
	}

	outputs[0].Set_bool(success)
}

// More return true if there is another element in the current array or object being parsed.
func opJsonTokenMore(inputs []ast.CXValue, outputs []ast.CXValue) {
	more := false
	success := false

	if jsonFile := validJsonFile(inputs[0].Get_i32()); jsonFile != nil {
		more = jsonFile.decoder.More()
		success = true
	}

	outputs[0].Set_bool(more)
	outputs[1].Set_bool(success)
}

// Token parses the next token.
func opJsonTokenNext(inputs []ast.CXValue, outputs []ast.CXValue) {
	tokenType := int32(JSON_TOKEN_INVALID)
	success := false

	if jsonFile := validJsonFile(inputs[0].Get_i32()); jsonFile != nil {
		token, err := jsonFile.decoder.Token()
		if err == io.EOF {
			tokenType = JSON_TOKEN_NULL
			success = true
		} else if err == nil {
			jsonFile.token = token
			switch value := token.(type) {
			case json.Delim:
				tokenType = JSON_TOKEN_DELIM
				jsonFile.tokenDelim = value
				success = true
			case bool:
				tokenType = JSON_TOKEN_BOOL
				jsonFile.tokenBool = value
				success = true
			case float64:
				tokenType = JSON_TOKEN_F64
				jsonFile.tokenF64 = value
				success = true
			case json.Number:
				tokenType = JSON_TOKEN_NUMBER
				jsonFile.tokenNumber = value
				success = true
			case string:
				tokenType = JSON_TOKEN_STR
				jsonFile.tokenStr = value
				success = true
			default:
				if value == nil {
					tokenType = JSON_TOKEN_NULL
					success = true
				}
			}
		}
		jsonFile.tokenType = tokenType
	}

	outputs[0].Set_i32(tokenType)
	outputs[1].Set_bool(success)
}

// Type returns the type of the current token.
func opJsonTokenType(inputs []ast.CXValue, outputs []ast.CXValue) {
	tokenType := int32(JSON_TOKEN_INVALID)
	success := false

	if jsonFile := validJsonFile(inputs[0].Get_i32()); jsonFile != nil {
		tokenType = jsonFile.tokenType
		success = true
	}

	outputs[0].Set_i32(tokenType)
	outputs[1].Set_bool(success)
}

// Delim returns current token as an int32 delimiter.
func opJsonTokenDelim(inputs []ast.CXValue, outputs []ast.CXValue) {
	tokenDelim := int32(JSON_TOKEN_INVALID)
	success := false

	if jsonFile := validJsonFile(inputs[0].Get_i32()); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_DELIM {
			tokenDelim = int32(jsonFile.tokenDelim)
			success = true
		}
	}

	outputs[0].Set_i32(tokenDelim)
	outputs[1].Set_bool(success)
}

// Bool returns current token as a bool value.
func opJsonTokenBool(inputs []ast.CXValue, outputs []ast.CXValue) {
	tokenBool := false
	success := false

	if jsonFile := validJsonFile(inputs[0].Get_i32()); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_BOOL {
			tokenBool = jsonFile.tokenBool
			success = true
		}
	}

	outputs[0].Set_bool(tokenBool)
	outputs[1].Set_bool(success)
}

// Float64 returns current token as float64 value.
func opJsonTokenF64(inputs []ast.CXValue, outputs []ast.CXValue) {
	var tokenF64 float64
	success := false

	if jsonFile := validJsonFile(inputs[0].Get_i32()); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_F64 {
			tokenF64 = jsonFile.tokenF64
			success = true
		} else if jsonFile.tokenType == JSON_TOKEN_NUMBER {
			var err error
			if tokenF64, err = jsonFile.tokenNumber.Float64(); err == nil {
				success = true
			}
		}
	}

	outputs[0].Set_f64(tokenF64)
	outputs[1].Set_bool(success)
}

// Int64 returns current token as int64 value.
func opJsonTokenI64(inputs []ast.CXValue, outputs []ast.CXValue) {
	var tokenI64 int64
	success := false

	if jsonFile := validJsonFile(inputs[0].Get_i32()); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_NUMBER {
			var err error
			if tokenI64, err = jsonFile.tokenNumber.Int64(); err == nil {
				success = true
			}
		}
	}

	outputs[0].Set_i64(tokenI64)
	outputs[1].Set_bool(success)
}

// Str returns current token as string value.
func opJsonTokenStr(inputs []ast.CXValue, outputs []ast.CXValue) {
	var tokenStr string
	success := false

	if jsonFile := validJsonFile(inputs[0].Get_i32()); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_STR {
			tokenStr = jsonFile.tokenStr
			success = true
		}
	}

    outputs[0].Set_str(tokenStr)
	outputs[1].Set_bool(success)
}

// helper function used to validate json handle from i32
func validJsonFile(handle int32) *JSONFile {
	if handle >= 0 && handle < int32(len(jsons)) && jsons[handle].file != nil {
		return &jsons[handle]
	}
	return nil
}
