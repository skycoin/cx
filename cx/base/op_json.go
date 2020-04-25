// +build base

package cxcore

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	. "github.com/SkycoinProject/cx/cx"
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
func opJsonOpen(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	handle := int32(-1)

	file, err := CXOpenFile(ReadStr(fp, expr.Inputs[0]))
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

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(handle))
}

// Close json parser (and all underlying resources) idendified by it's i32 handle.
func opJsonClose(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	success := false

	handle := ReadI32(fp, expr.Inputs[0])
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		if err := jsonFile.file.Close(); err != nil {
			panic(err)
		}

		jsons[handle] = JSONFile{}
		freeJsons = append(freeJsons, handle)
		success = true
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), success)
}

// More return true if there is another element in the current array or object being parsed.
func opJsonTokenMore(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	more := false
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		more = jsonFile.decoder.More()
		success = true
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), more)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

// Token parses the next token.
func opJsonTokenNext(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), tokenType)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

// Type returns the type of the current token.
func opJsonTokenType(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	tokenType := int32(JSON_TOKEN_INVALID)
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		tokenType = jsonFile.tokenType
		success = true
	}

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), tokenType)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

// Delim returns current token as an int32 delimiter.
func opJsonTokenDelim(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	tokenDelim := int32(JSON_TOKEN_INVALID)
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_DELIM {
			tokenDelim = int32(jsonFile.tokenDelim)
			success = true
		}
	}

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), tokenDelim)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

// Bool returns current token as a bool value.
func opJsonTokenBool(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	tokenBool := false
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_BOOL {
			tokenBool = jsonFile.tokenBool
			success = true
		}
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), tokenBool)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

// Float64 returns current token as float64 value.
func opJsonTokenF64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), tokenF64)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

// Int64 returns current token as int64 value.
func opJsonTokenI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), tokenI64)
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

// Str returns current token as string value.
func opJsonTokenStr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var tokenStr string
	success := false

	if jsonFile := validJsonFileExpr(expr, fp); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_STR {
			tokenStr = jsonFile.tokenStr
			success = true
		}
	}

	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(tokenStr))
	WriteBool(GetFinalOffset(fp, expr.Outputs[1]), success)
}

// helper function used to validate json handle from expr
func validJsonFileExpr(expr *CXExpression, fp int) *JSONFile {
	handle := ReadI32(fp, expr.Inputs[0])
	return validJsonFile(handle)
}

// helper function used to validate json handle from i32
func validJsonFile(handle int32) *JSONFile {
	if handle >= 0 && handle < int32(len(jsons)) && jsons[handle].file != nil {
		return &jsons[handle]
	}
	return nil
}
