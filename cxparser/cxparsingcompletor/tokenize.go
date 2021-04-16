package parsingcompletor

import (
	"fmt"
	"io"
	"strconv"
)

func Tokenize(r io.Reader, w io.Writer) {
	var sym yySymType

	lex := NewLexer(r)
	token := lex.Lex(&sym)
	for token > 0 {
		fmt.Fprintln(w, TokenName(token), TokenValue(token, &sym))
		token = lex.Lex(&sym)
	}
}

func TokenValue(token int, sym *yySymType) interface{} {
	switch token {
	case BOOLEAN_LITERAL:
		return sym.bool
	case BYTE_LITERAL:
		return sym.i8
	case SHORT_LITERAL:
		return sym.i16
	case INT_LITERAL:
		return sym.i32
	case LONG_LITERAL:
		return sym.i64
	case UNSIGNED_BYTE_LITERAL:
		return sym.ui8
	case UNSIGNED_SHORT_LITERAL:
		return sym.ui16
	case UNSIGNED_INT_LITERAL:
		return sym.ui32
	case UNSIGNED_LONG_LITERAL:
		return sym.ui64
	case FLOAT_LITERAL:
		return sym.f32
	case DOUBLE_LITERAL:
		return sym.f64
	case AFF, BOOL, F32, F64, I8, I16, I32, I64,
		UI8, UI16, UI32, UI64, REF_OP, ADD_OP, SUB_OP, MUL_OP, DIV_OP, MOD_OP,
		GT_OP, LT_OP, GTEQ_OP, LTEQ_OP, RIGHT_ASSIGN, LEFT_ASSIGN, ADD_ASSIGN,
		SUB_ASSIGN, MUL_ASSIGN, DIV_ASSIGN, MOD_ASSIGN, AND_ASSIGN, XOR_ASSIGN,
		OR_ASSIGN, NEG_OP, ASSIGN, CASSIGN, STRING_LITERAL, IDENTIFIER:
		return sym.tok
	}
	return ""
}

func TokenName(token int) string {
	switch token {
	case ADDR:
		return "  ADDR"
	case ADD_ASSIGN:
		return "ADDSET"
	case ADD_OP:
		return " ADDOP"
	case AFF:
		return "   AFF"
	case AFFVAR:
		return "AFFVAR"
	case AND:
		return "   AND"
	case AND_ASSIGN:
		return "ANDSET"
	case AND_OP:
		return " ANDOP"
	case ASSIGN:
		return "  ASGN"
	case BASICTYPE:
		return " BASIC"
	case BITANDEQ:
		return "BANDEQ"
	case BITCLEAR_OP:
		return "BCLROP"
	case BITOREQ:
		return " BOREQ"
	case BITOR_OP:
		return " BOROP"
	case BITXOREQ:
		return "BXOREQ"
	case BITXOR_OP:
		return "BXOROP"
	case BOOL:
		return "  BOOL"
	case BOOLEAN_LITERAL:
		return "BOOLLT"
	case BREAK:
		return " BREAK"
	case BYTE_LITERAL:
		return "BYTELT"
	case CAFF:
		return "  CAFF"
	case CASE:
		return "  CASE"
	case CASSIGN:
		return "CASSGN"
	case CLAUSES:
		return "CLAUSE"
	case COLON:
		return " COLON"
	case COMMA:
		return " COMMA"
	case COMMENT:
		return "COMMNT"
	case CONST:
		return " CONST"
	case CONTINUE:
		return "CONTNU"
	case DEC_OP:
		return " DECOP"
	case DEF:
		return "   DEF"
	case DEFAULT:
		return "DFAULT"
	case DIVEQ:
		return " DIVEQ"
	case DIV_ASSIGN:
		return "DIVSET"
	case DIV_OP:
		return " DIVOP"
	case DOUBLE_LITERAL:
		return "DBLLIT"
	case DPROGRAM:
		return " DPROG"
	case DSTACK:
		return "DSTACK"
	case DSTATE:
		return "DSTATE"
	case ELSE:
		return "  ELSE"
	case ENUM:
		return "  ENUM"
	case EQUAL:
		return " EQUAL"
	case EQUALWORD:
		return "EQWORD"
	case EQ_OP:
		return "  EQOP"
	case EXP:
		return "   EXP"
	case EXPEQ:
		return " EXPEQ"
	case EXPR:
		return "  EXPR"
	case F32:
		return "   F32"
	case F64:
		return "   F64"
	case FIELD:
		return " FIELD"
	case FLOAT_LITERAL:
		return "FLOATL"
	case FOR:
		return "   FOR"
	case FUNC:
		return "  FUNC"
	case GE_OP:
		return "  GEOP"
	case GOTO:
		return "  GOTO"
	case GTEQ_OP:
		return "GTEQOP"
	case GTHANEQ:
		return "GTHNEQ"
	case GTHANWORD:
		return " GTHNW"
	case GT_OP:
		return "  GTOP"
	case I16:
		return "   I16"
	case I32:
		return "   I32"
	case I64:
		return "   I64"
	case I8:
		return "    I8"
	case IDENTIFIER:
		return " IDENT"
	case IF:
		return "    IF"
	case IMPORT:
		return "IMPORT"
	case INC_OP:
		return " INCOP"
	case INFER:
		return " INFER"
	case INT_LITERAL:
		return "INTLIT"
	case LBRACE:
		return "LBRACE"
	case LBRACK:
		return "LBRACK"
	case LEFTSHIFT:
		return "LSHIFT"
	case LEFTSHIFTEQ:
		return " LSHEQ"
	case LEFT_ASSIGN:
		return "  LSET"
	case LEFT_OP:
		return "LEFTOP"
	case LE_OP:
		return "  LEOP"
	case LONG_LITERAL:
		return "LONGLT"
	case LPAREN:
		return "LPAREN"
	case LTEQ_OP:
		return "LTEQOP"
	case LTHANEQ:
		return "LTHNEQ"
	case LTHANWORD:
		return "LTHANW"
	case LT_OP:
		return "  LTOP"
	case MINUSEQ:
		return "MNUSEQ"
	case MINUSMINUS:
		return "MINUS2"
	case MOD_ASSIGN:
		return "MODSET"
	case MOD_OP:
		return " MODOP"
	case MULTEQ:
		return "MULTEQ"
	case MUL_ASSIGN:
		return "MULSET"
	case MUL_OP:
		return " MULOP"
	case NEG_OP:
		return " NEGOP"
	case NEW:
		return "   NEW"
	case NEWLINE:
		return "NEWLIN"
	case NE_OP:
		return "  NEOP"
	case NOT:
		return "   NOT"
	case OBJECT:
		return "OBJECT"
	case OBJECTS:
		return "OBJCTS"
	case OP:
		return "    OP"
	case OR:
		return "    OR"
	case OR_ASSIGN:
		return " ORSET"
	case OR_OP:
		return "  OROP"
	case PACKAGE:
		return "PACKAG"
	case PERIOD:
		return "PERIOD"
	case PLUSEQ:
		return "PLUSEQ"
	case PLUSPLUS:
		return " PLUS2"
	//case PSTEP:
	//	return " PSTEP"
	//case STEP:
	//	return "  STEP"
	//case TSTEP:
	//	return " TSTEP"
	case PTR_OP:
		return " PTROP"
	case RBRACE:
		return "RBRACE"
	case RBRACK:
		return "RBRACK"
	case REF_OP:
		return " REFOP"
	case REM:
		return "   REM"
	case REMAINDER:
		return "REMNDR"
	case REMAINDEREQ:
		return "RMDREQ"
	case RETURN:
		return "RETURN"
	case RIGHTSHIFT:
		return "RSHIFT"
	case RIGHTSHIFTEQ:
		return " RSHEQ"
	case RIGHT_ASSIGN:
		return "  RSET"
	case RIGHT_OP:
		return "RGHTOP"
	case RPAREN:
		return "RPAREN"
	case SEMICOLON:
		return "SCOLON"
	case SHORT_LITERAL:
		return "SHRTLT"
	//case SPACKAGE:
	//	return "SPACKG"
	//case SSTRUCT:
	//	return "SSTRCT"
	//case SFUNC:
	//	return " SFUNC"
	case STR:
		return "   STR"
	case STRING_LITERAL:
		return "STRLIT"
	case STRUCT:
		return "STRUCT"
	case SUB_ASSIGN:
		return "SUBSET"
	case SUB_OP:
		return " SUBOP"
	case SWITCH:
		return "SWITCH"
	case TAG:
		return "   TAG"
	case TYPE:
		return "  TYPE"
	case TYPSTRUCT:
		return "TSTRCT"
	case UI16:
		return "  UI16"
	case UI32:
		return "  UI32"
	case UI64:
		return "  UI64"
	case UI8:
		return "   UI8"
	case UNEQUAL:
		return "  UNEQ"
	case UNION:
		return " UNION"
	case UNSIGNED_BYTE_LITERAL:
		return " UBLIT"
	case UNSIGNED_INT_LITERAL:
		return " UILIT"
	case UNSIGNED_LONG_LITERAL:
		return " ULLIT"
	case UNSIGNED_SHORT_LITERAL:
		return " USLIT"
	case VALUE:
		return " VALUE"
	case VAR:
		return "   VAR"
	case XOR_ASSIGN:
		return "XORSET"
	}
	return "UNK" + strconv.Itoa(token)
}
