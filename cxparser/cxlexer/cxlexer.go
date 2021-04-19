package cxlexer

type token int

const (
	ILLEGAL token = iota
	EOF

	IDENT
)
