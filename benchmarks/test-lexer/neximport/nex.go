package neximport

import (
	"io"

	. "github.com/SkycoinProject/cx/tests/test-lexer/oldnex"
)

func NewLexer_(r io.Reader) *Lexer {
	return NewLexer(r)
}
