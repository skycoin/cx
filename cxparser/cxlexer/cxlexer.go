package cxlexer

type TokenType string

const (
	ILLEGAL = "ILLEGAL"

	EOF = "EOF"

	IDENT = "IDENT"

	INT = "INT"

	STRING = "STRING"

	ASSIGN = "="

	PLUS = "+"

	MINUS = "-"

	BANG = "!"

	ASTERISK = "*"

	SLASH = "/"

	LT = "<"

	GT = ">"

	EQ = "=="

	NOT_EQ = "!="

	COMMA = ","

	SEMICON = ";"

	SECON = ":"

	LPAREN = "("

	RPAREN = ")"

	LBRACKET = "{"

	LBRACKET = "}"

	LBRACKET = "["

	LBRACKET = "]"

	FUNCTION = "func"

	VAR = "VAR"

	TRUE = "TRUE"

	FALSE = "FALSE"

	IF = "IF"

	ELSE = "ELSE"

	RETURN = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {

	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
