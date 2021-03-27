package cxparser

import "regexp"

type RegularExpression struct {
	// comments
	reComment           *regexp.Regexp
	reMultiCommentOpen  *regexp.Regexp
	reMultiCommentClose *regexp.Regexp

	// packages and structs
	rePackage     *regexp.Regexp
	rePackageName *regexp.Regexp
	reStruct      *regexp.Regexp
	reStructName  *regexp.Regexp

	// globals
	reGlobal     *regexp.Regexp
	reGlobalName *regexp.Regexp

	// body open/close
	reBodyOpen  *regexp.Regexp
	reBodyClose *regexp.Regexp

	// imports
	reImport     *regexp.Regexp
	reImportName *regexp.Regexp
}

func newRegulaEexpression() RegularExpression {

	re := RegularExpression{

		reMultiCommentOpen: regexp.MustCompile(`/\*`),

		reMultiCommentClose: regexp.MustCompile(`\*/`),

		reComment: regexp.MustCompile("//"),

		rePackage: regexp.MustCompile("package"),

		rePackageName: regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`),

		reStruct: regexp.MustCompile("type"),

		reStructName: regexp.MustCompile(`(^|[\s])type\s+([_a-zA-Z][_a-zA-Z0-9]*)?\s`),

		reGlobal: regexp.MustCompile("var"),

		reGlobalName: regexp.MustCompile(`(^|[\s])var\s([_a-zA-Z][_a-zA-Z0-9]*)`),

		reBodyOpen: regexp.MustCompile("{"),

		reBodyClose: regexp.MustCompile("}"),

		reImport: regexp.MustCompile("import"),

		reImportName: regexp.MustCompile(`(^|[\s])import\s+"([_a-zA-Z][_a-zA-Z0-9/-]*)"`),
	}

	return re
}
