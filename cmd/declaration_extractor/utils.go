package declaration_extractor

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode"
)

// Modified ScanLines to include "\r\n" or "\n" in line
func scanLinesWithLineTerminator(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		advance = i + 1 // i + 1 includes the line termninator
		if data[i] == '\n' {
			// We have a line terminated by single newline.
			return advance, data[0:advance], nil
		}

		if len(data) > i+1 && data[i+1] == '\n' {
			advance += 1
		}
		return advance, data[0:advance], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func ReplaceCommentsWithWhitespaces(source []byte) []byte {

	reComment := regexp.MustCompile(`//.*|/\*[\s\S]*?\*/|\"//.*\"`)

	comments := reComment.FindAllIndex(source, -1)

	// Loops through each character and replaces with whitespcace
	for i := range comments {
		for loc := comments[i][0]; loc < comments[i][1]; loc++ {
			if unicode.IsSpace(rune(source[loc])) {
				continue
			}
			source[loc] = byte(' ')
		}
	}

	return source
}

func ReplaceStringContentsWithWhitespaces(source []byte) ([]byte, error) {

	var sourceWithoutStringContents []byte
	sourceWithoutStringContents = append(sourceWithoutStringContents, source...)
	var inStdString bool
	var inRawString bool
	var lineno int

	for i, char := range source {

		if char == '\n' {
			lineno++
		}

		//if end of line and quote not terminated
		if char == '\n' && inStdString {
			return sourceWithoutStringContents, fmt.Errorf("%v: syntax error: quote not terminated", lineno)
		}

		if char == '"' && !inStdString && !inRawString {
			inStdString = true
			continue
		}

		if char == '"' && inStdString {
			inStdString = false
			continue
		}

		if char == '`' && !inRawString && !inStdString {
			inRawString = true
			continue
		}

		if char == '`' && inRawString {
			inRawString = false
			continue
		}

		if (inStdString || inRawString) && !unicode.IsSpace(rune(char)) {
			sourceWithoutStringContents[i] = byte(' ')
		}
	}

	return sourceWithoutStringContents, nil
}

func ContainsTokenByte(tokenList [][]byte, tokenByte []byte) bool {
	for _, token := range tokenList {
		if bytes.Equal(token, tokenByte) {
			return true
		}
	}
	return false
}

func ContainsTokenByteInToken(tokenList [][]byte, tokenByte []byte) bool {
	for _, token := range tokenList {
		if bytes.Contains(token, tokenByte) {
			return true
		}
	}
	return false
}
