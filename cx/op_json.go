// +build base extra full

package cxcore

import (
	"bufio"
	"encoding/json"
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

// Open the named json file for reading, returns an i32 identifying the json parser.
func opJSONOpen(expr *CXExpression, fp int) {
	handle := int32(-1)

	file, err := os.Open(ReadStr(fp, expr.Inputs[0]))
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

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(handle)))
}

// Close json parser (and all underlying resources) idendified by it's i32 handle.
func opJSONClose(expr *CXExpression, fp int) {
	success := false

	handle := ReadI32(fp, expr.Inputs[0])
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		jsonFile.file.Close()
		jsons[handle] = JSONFile{}
		freeJsons = append(freeJsons, handle)
		success = true
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromBool(success))
}

// More return true if there is another element in the current array or object being parsed.
func opJSONTokenMore(expr *CXExpression, fp int) {
	more := false
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		more = jsonFile.decoder.More()
		success = true
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromBool(more))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromBool(success))
}

// Token parses the next token.
func opJSONTokenNext(expr *CXExpression, fp int) {
	tokenType := int32(JSON_TOKEN_INVALID)
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
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

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(tokenType))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromBool(success))
}

// Type returns the type of the current token.
func opJSONTokenType(expr *CXExpression, fp int) {
	tokenType := int32(JSON_TOKEN_INVALID)
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		tokenType = jsonFile.tokenType
		success = true
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(tokenType))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromBool(success))
}

// Delim returns current token as an int32 delimiter.
// Panics if token type is not JSON_TOKEN_DELIM.
func opJSONTokenDelim(expr *CXExpression, fp int) {
	tokenDelim := int32(JSON_TOKEN_INVALID)
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_DELIM {
			tokenDelim = int32(jsonFile.tokenDelim)
			success = true
		}
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(tokenDelim))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromBool(success))
}

// Bool returns current token as a bool value.
// Panics if token type is not JSON_TOKEN_BOOL.
func opJSONTokenBool(expr *CXExpression, fp int) {
	tokenBool := false
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_BOOL {
			tokenBool = jsonFile.tokenBool
			success = true
		}
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromBool(tokenBool))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromBool(success))
}

// Float64 returns current token as float64 value.
// Panics if token can't be interpreted as float64 value.
func opJSONTokenF64(expr *CXExpression, fp int) {
	var tokenF64 float64
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
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

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromF64(tokenF64))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromBool(success))
}

// Int64 returns current token as int64 value.
// Panics if  token can't be interpreted as int64 value.
func opJSONTokenI64(expr *CXExpression, fp int) {
	var tokenI64 int64
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_NUMBER {
			var err error
			if tokenI64, err = jsonFile.tokenNumber.Int64(); err == nil {
				success = true
			}
		}
	}

	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI64(tokenI64))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromBool(success))
}

// Str returns current token as string value.
// Panics if token type is not JSON_TOKEN_STR.
func opJSONTokenStr(expr *CXExpression, fp int) {
	var tokenStr string
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_STR {
			tokenStr = jsonFile.tokenStr
			success = true
		}
	}

	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(tokenStr))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[1]), FromBool(success))
}

// helper function used to validate json handle from expr
func validJsonFileExpr(expr *CXExpression, fp int) *JSONFile {
	handle := ReadI32(fp, expr.Inputs[0])
	return validJsonFile(handle)
}

// helper function used to valid json handle from i32
func validJsonFile(handle int32) *JSONFile {
	if handle >= 0 && handle < int32(len(jsons)) && jsons[handle].file != nil {
		return &jsons[handle]
	}
	return nil
}
