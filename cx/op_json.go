// +build base extra full

package cxcore

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

const (
	JSON_TOKEN_NULL = iota
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

// Open the named json file for reading, returns an i32 identifying the json parser.
func opJSONOpen(expr *CXExpression, fp int) {
	file, err := os.Open(ReadStr(fp, expr.Inputs[0]))
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	var jsonFile JSONFile
	jsonFile.file = file
	jsonFile.reader = reader
	jsonFile.decoder = json.NewDecoder(jsonFile.reader)
	jsonFile.decoder.UseNumber()
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(len(jsons))))
	jsons = append(jsons, jsonFile)
}

// Close json parser (and all underlying resources) idendified by it's i32 handle.
func opJSONClose(expr *CXExpression, fp int) {
	jsonFile := jsons[ReadI32(fp, expr.Inputs[0])]
	jsonFile.file.Close()
}

// More return true if there is another element in the current array or object being parsed.
func opJSONTokenMore(expr *CXExpression, fp int) {
	jsonFile := jsons[ReadI32(fp, expr.Inputs[0])]
	more := jsonFile.decoder.More()
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromBool(more))
}

// Token parses the next token.
func opJSONTokenNext(expr *CXExpression, fp int) {
	index := ReadI32(fp, expr.Inputs[0])
	var tokenType int32 = JSON_TOKEN_NULL

	token, err := jsons[index].decoder.Token()
	if err == io.EOF {
	} else if err != nil {
		panic(err)
	} else {
		jsons[index].token = token
		switch value := token.(type) {
		case json.Delim:
			tokenType = JSON_TOKEN_DELIM
			jsons[index].tokenDelim = value
		case bool:
			tokenType = JSON_TOKEN_BOOL
			jsons[index].tokenBool = value
		case float64:
			tokenType = JSON_TOKEN_F64
			jsons[index].tokenF64 = value
		case json.Number:
			tokenType = JSON_TOKEN_NUMBER
			jsons[index].tokenNumber = value
		case string:
			tokenType = JSON_TOKEN_STR
			jsons[index].tokenStr = value
		default:
			tokenType = JSON_TOKEN_NULL
		}
	}
	jsons[index].tokenType = tokenType
}

// Type returns the type of the current token.
func opJSONTokenType(expr *CXExpression, fp int) {
	jsonFile := jsons[ReadI32(fp, expr.Inputs[0])]
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(jsonFile.tokenType))
}

// Delim returns current token as an int32 delimiter.
// Panics if token type is not JSON_TOKEN_DELIM.
func opJSONTokenDelim(expr *CXExpression, fp int) {
	jsonFile := jsons[ReadI32(fp, expr.Inputs[0])]
	if jsonFile.tokenType != JSON_TOKEN_DELIM {
		panic("json : not a delim value")
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(jsonFile.tokenDelim)))
}

// Bool returns current token as a bool value.
// Panics if token type is not JSON_TOKEN_BOOL.
func opJSONTokenBool(expr *CXExpression, fp int) {
	jsonFile := jsons[ReadI32(fp, expr.Inputs[0])]
	if jsonFile.tokenType != JSON_TOKEN_BOOL {
		panic("json : not a bool value")
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromBool(jsonFile.tokenBool))
}

// Float64 returns current token as float64 value.
// Panics if token can't be interpreted as float64 value.
func opJSONTokenF64(expr *CXExpression, fp int) {
	jsonFile := jsons[ReadI32(fp, expr.Inputs[0])]
	var value float64
	if jsonFile.tokenType == JSON_TOKEN_F64 {
		value = jsonFile.tokenF64
	} else if jsonFile.tokenType == JSON_TOKEN_NUMBER {
		var err error
		value, err = jsonFile.tokenNumber.Float64()
		if err != nil {
			panic(err)
		}
	} else {
		panic("json : not a f64 value")
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromF64(value))
}

// Int64 returns current token as int64 value.
// Panics if  token can't be interpreted as int64 value.
func opJSONTokenI64(expr *CXExpression, fp int) {
	jsonFile := jsons[ReadI32(fp, expr.Inputs[0])]
	var value int64
	if jsonFile.tokenType == JSON_TOKEN_NUMBER {
		var err error
		value, err = jsonFile.tokenNumber.Int64()
		if err != nil {
			panic(err)
		}
	} else {
		panic("json : not an int64 value")
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI64(value))
}

// Str returns current token as string value.
// Panics if token type is not JSON_TOKEN_STR.
func opJSONTokenStr(expr *CXExpression, fp int) {
	jsonFile := jsons[ReadI32(fp, expr.Inputs[0])]
	if jsonFile.tokenType != JSON_TOKEN_STR {
		panic("json : not a string value")
	}
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(jsonFile.tokenStr))
}
