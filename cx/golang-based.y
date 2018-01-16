%{
	package main
	import (
		// "strings"
"fmt"
		// "github.com/skycoin/skycoin/src/cipher/encoder"
		. "github.com/skycoin/cx/src/base"
                )

	var prgrm = MakeProgram(1024, 1024, 1024)

	var lineNo int = 0
	var webMode bool = false
	var baseOutput bool = false
	var replMode bool = false
	var helpMode bool = false
	var compileMode bool = false
	var replTargetFn string = ""
	var replTargetStrct string = ""
	var replTargetMod string = ""
	var dStack bool = false
	var inREPL bool = false
	var inFn bool = false
	var tag string = ""
	var asmNL = "\n"
	var fileName string
%}

%union {
	i int
	byt byte
	i32 int32
	i64 int64
	f32 float32
	f64 float64
	tok string
	bool bool
	string string
	stringA []string

	line int
}

%token	<tok>	LLITERAL
%token	<tok>	LASOP LCOLAS
%token	<tok>	LBREAK LCASE LCHAN LCONST LCONTINUE LDDD
%token	<tok>	LDEFAULT LDEFER LELSE LFALL LFOR LFUNC LGO LGOTO
%token	<tok>	LIF LIMPORT LINTERFACE LMAP LNAME
%token	<tok>	LPACKAGE LRANGE LRETURN LSELECT LSTRUCT LSWITCH
%token	<tok>	LTYPE LVAR






%token  <byt>           BYTE_LITERAL
%token  <i32>           INT_LITERAL BOOLEAN_LITERAL
%token  <i64>           LONG_LITERAL
%token  <f32>           FLOAT_LITERAL
%token  <f64>           DOUBLE_LITERAL
%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK IDENTIFIER
                        VAR COMMA COMMENT STRING_LITERAL PACKAGE IF ELSE FOR TYPSTRUCT STRUCT
                        ASSIGN CASSIGN IMPORT RETURN GOTO GTHAN LTHAN EQUAL COLON NEW
                        EQUALWORD GTHANWORD LTHANWORD
                        GTHANEQ LTHANEQ UNEQUAL AND OR
                        PLUS MINUS MULT DIV AFFVAR
                        PLUSPLUS MINUSMINUS REMAINDER LEFTSHIFT RIGHTSHIFT EXP
                        NOT
                        BITAND BITXOR BITOR BITCLEAR
                        PLUSEQ MINUSEQ MULTEQ DIVEQ REMAINDEREQ EXPEQ
                        LEFTSHIFTEQ RIGHTSHIFTEQ BITANDEQ BITXOREQ BITOREQ

                        DEC_OP INC_OP PTR_OP LEFT_OP RIGHT_OP
                        GE_OP LE_OP EQ_OP NE_OP AND_OP OR_OP
                        ADD_ASSIGN AND_ASSIGN LEFT_ASSIGN MOD_ASSIGN
                        MUL_ASSIGN DIV_ASSIGN OR_ASSIGN RIGHT_ASSIGN
                        SUB_ASSIGN XOR_ASSIGN
                        BOOL BYTE F32 F64
                        I8 I16 I32 I64
                        STR
                        UI8 UI16 UI32 UI64
                        UNION ENUM CONST CASE DEFAULT SWITCH BREAK CONTINUE
                        TYPE

                        /* Types */
                        BASICTYPE
                        /* Selectors */
                        SPACKAGE SSTRUCT SFUNC
                        /* Removers */
                        REM DEF EXPR FIELD INPUT OUTPUT CLAUSES OBJECT OBJECTS
                        /* Stepping */
                        STEP PSTEP TSTEP
                        /* Debugging */
                        DSTACK DPROGRAM DSTATE
                        /* Affordances */
                        AFF TAG INFER VALUE
                        /* Pointers */
                        ADDR



                        

%token		LANDAND LANDNOT LBODY LCOMM LDEC LEQ LGE LGT
%token		LIGNORE LINC LLE LLSH LLT LNE LOROR LRSH

/* %type	<i>	lbrace import_here */
/* %type	<sym>	sym packname */
/* %type	<val>	oliteral */

/* %type	<node>	stmt ntype */
/* %type	<node>	arg_type */
/* %type	<node>	case caseblock */
/* %type	<node>	compound_stmt dotname embed expr complitexpr bare_complitexpr */
/* %type	<node>	expr_or_type */
/* %type	<node>	fndcl hidden_fndcl fnliteral */
/* %type	<node>	for_body for_header for_stmt if_header if_stmt non_dcl_stmt */
/* %type	<node>	interfacedcl keyval labelname name */
/* %type	<node>	name_or_type non_expr_type */
/* %type	<node>	new_name dcl_name oexpr typedclname */
/* %type	<node>	onew_name */
/* %type	<node>	osimple_stmt pexpr pexpr_no_paren */
/* %type	<node>	pseudocall range_stmt select_stmt */
/* %type	<node>	simple_stmt */
/* %type	<node>	switch_stmt uexpr */
/* %type	<node>	xfndcl typedcl start_complit */

/* %type	<list>	xdcl fnbody fnres loop_body dcl_name_list */
/* %type	<list>	new_name_list expr_list keyval_list braced_keyval_list expr_or_type_list xdcl_list */
/* %type	<list>	oexpr_list caseblock_list elseif elseif_list else stmt_list oarg_type_list_ocomma arg_type_list */
/* %type	<list>	interfacedcl_list vardcl vardcl_list structdcl structdcl_list */
/* %type	<list>	common_dcl constdcl constdcl1 constdcl_list typedcl_list */

/* %type	<node>	convtype comptype dotdotdot */
/* %type	<node>	indcl interfacetype structtype ptrtype */
/* %type	<node>	recvchantype non_recvchantype othertype fnret_type fntype */

/* %type	<sym>	hidden_importsym hidden_pkg_importsym */

/* %type	<node>	hidden_constant hidden_literal hidden_funarg */
/* %type	<node>	hidden_interfacedcl hidden_structdcl */

/* %type	<list>	hidden_funres */
/* %type	<list>	ohidden_funres */
/* %type	<list>	hidden_funarg_list ohidden_funarg_list */
/* %type	<list>	hidden_interfacedcl_list ohidden_interfacedcl_list */
/* %type	<list>	hidden_structdcl_list ohidden_structdcl_list */

/* %type	<typ>	hidden_type hidden_type_misc hidden_pkgtype */
/* %type	<typ>	hidden_type_func */
/* %type	<typ>	hidden_type_recv_chan hidden_type_non_recv_chan */

%left		LCOMM	/* outside the usual hierarchy; here for good error messages */

%left		LOROR
%left		LANDAND
%left		LEQ LNE LLE LGE LLT LGT
%left		'+' '-' '|' '^'
%left		'*' '/' '%' '&' LLSH LRSH LANDNOT

/*
 * manual override of shift/reduce conflicts.
 * the general form is that we assign a precedence
 * to the token being shifted and then introduce
 * NotToken with lower precedence or PreferToToken with higher
 * and annotate the reducing rule accordingly.
 */
%left		NotPackage
%left		LPACKAGE

%left		NotParen
%left		'('

%left		')'
%left		PreferToRightParen

%error-verbose

%%
file:
	loadsys
	package
	imports
	xdcl_list

package:
                %prec NotPackage
                {
                    fmt.Println("huh")
                }
        |	LPACKAGE LNAME ';'
                {
                    fmt.Println("hi")
                }

/*
 * this loads the definitions for the low-level runtime functions,
 * so that the compiler can generate calls to them,
 * but does not make the name "runtime" visible as a package.
 */
loadsys:
	import_package
	import_there

imports:
|	imports import ';'

import:
	LIMPORT import_stmt
|	LIMPORT '(' import_stmt_list osemi ')'
|	LIMPORT '(' ')'

import_stmt:
	import_here import_package import_there
|	import_here import_there

import_stmt_list:
	import_stmt
|	import_stmt_list ';' import_stmt

import_here:
	LLITERAL
|	sym LLITERAL
|	'.' LLITERAL

import_package:
	LPACKAGE LNAME import_safety ';'

import_safety:
|	LNAME

import_there:
	hidden_import_list '$' '$'

/*
 * declarations
 */
xdcl:
|	common_dcl
|	xfndcl
|	non_dcl_stmt
|	error

common_dcl:
	LVAR vardcl
|	LVAR '(' vardcl_list osemi ')'
|	LVAR '(' ')'
|	lconst constdcl
|	lconst '(' constdcl osemi ')'
|	lconst '(' constdcl ';' constdcl_list osemi ')'
|	lconst '(' ')'
|	LTYPE typedcl
|	LTYPE '(' typedcl_list osemi ')'
|	LTYPE '(' ')'

lconst:
	LCONST

vardcl:
	dcl_name_list ntype
|	dcl_name_list ntype '=' expr_list
|	dcl_name_list '=' expr_list

constdcl:
	dcl_name_list ntype '=' expr_list
|	dcl_name_list '=' expr_list

constdcl1:
	constdcl
|	dcl_name_list ntype
|	dcl_name_list

typedclname:
	sym

typedcl:
	typedclname ntype

simple_stmt:
	expr
|	expr LASOP expr
|	expr_list '=' expr_list
|	expr_list LCOLAS expr_list
|	expr LINC
|	expr LDEC

case:
	LCASE expr_or_type_list ':'
|	LCASE expr_or_type_list '=' expr ':'
|	LCASE expr_or_type_list LCOLAS expr ':'
|	LDEFAULT ':'

compound_stmt:
	'{'
	stmt_list '}'

caseblock:
	case
	stmt_list

caseblock_list:
|	caseblock_list caseblock

loop_body:
	LBODY
	stmt_list '}'

range_stmt:
	expr_list '=' LRANGE expr
|	expr_list LCOLAS LRANGE expr
|	LRANGE expr

for_header:
	osimple_stmt ';' osimple_stmt ';' osimple_stmt
|	osimple_stmt
|	range_stmt

for_body:
	for_header loop_body

for_stmt:
	LFOR
	for_body

if_header:
	osimple_stmt
|	osimple_stmt ';' osimple_stmt

/* IF cond body (ELSE IF cond body)* (ELSE block)? */
if_stmt:
	LIF
	if_header
	loop_body
	elseif_list else

elseif:
	LELSE LIF 
	if_header loop_body

elseif_list:
|	elseif_list elseif

else:
|	LELSE compound_stmt

switch_stmt:
	LSWITCH
	if_header
	LBODY caseblock_list '}'

select_stmt:
	LSELECT
	LBODY caseblock_list '}'

/*
 * expressions
 */
expr:
	uexpr
|	expr LOROR expr
|	expr LANDAND expr
|	expr LEQ expr
|	expr LNE expr
|	expr LLT expr
|	expr LLE expr
|	expr LGE expr
|	expr LGT expr
|	expr '+' expr
|	expr '-' expr
|	expr '|' expr
|	expr '^' expr
|	expr '*' expr
|	expr '/' expr
|	expr '%' expr
|	expr '&' expr
|	expr LANDNOT expr
|	expr LLSH expr
|	expr LRSH expr
	/* not an expression anymore, but left in so we can give a good error */
|	expr LCOMM expr

uexpr:
	pexpr
|	'*' uexpr
|	'&' uexpr
|	'+' uexpr
|	'-' uexpr
|	'!' uexpr
|	'~' uexpr
|	'^' uexpr
|	LCOMM uexpr

/*
 * call-like statements that
 * can be preceded by 'defer' and 'go'
 */
pseudocall:
	pexpr '(' ')'
|	pexpr '(' expr_or_type_list ocomma ')'
|	pexpr '(' expr_or_type_list LDDD ocomma ')'

pexpr_no_paren:
	LLITERAL
|	name
|	pexpr '.' sym
|	pexpr '.' '(' expr_or_type ')'
|	pexpr '.' '(' LTYPE ')'
|	pexpr '[' expr ']'
|	pexpr '[' oexpr ':' oexpr ']'
|	pexpr '[' oexpr ':' oexpr ':' oexpr ']'
|	pseudocall
|	convtype '(' expr ocomma ')'
|	comptype lbrace start_complit braced_keyval_list '}'
|	pexpr_no_paren '{' start_complit braced_keyval_list '}'
|	'(' expr_or_type ')' '{' start_complit braced_keyval_list '}'
|	fnliteral

start_complit:

keyval:
	expr ':' complitexpr

bare_complitexpr:
	expr
|	'{' start_complit braced_keyval_list '}'

complitexpr:
	expr
|	'{' start_complit braced_keyval_list '}'

pexpr:
	pexpr_no_paren
|	'(' expr_or_type ')'

expr_or_type:
	expr
|	non_expr_type	%prec PreferToRightParen

name_or_type:
	ntype

lbrace:
	LBODY
|	'{'

/*
 * names and types
 *	newname is used before declared
 *	oldname is used after declared
 */
new_name:
	sym

dcl_name:
	sym

onew_name:
|	new_name

sym:
	LNAME
|	hidden_importsym
|	'?'

hidden_importsym:
	'@' LLITERAL '.' LNAME
|	'@' LLITERAL '.' '?'

name:
	sym	%prec NotParen

labelname:
	new_name

/*
 * to avoid parsing conflicts, type is split into
 *	channel types
 *	function types
 *	parenthesized types
 *	any other type
 * the type system makes additional restrictions,
 * but those are not implemented in the grammar.
 */
dotdotdot:
	LDDD
|	LDDD ntype

ntype:
	recvchantype
|	fntype
|	othertype
|	ptrtype
|	dotname
|	'(' ntype ')'

non_expr_type:
	recvchantype
|	fntype
|	othertype
|	'*' non_expr_type

non_recvchantype:
	fntype
|	othertype
|	ptrtype
|	dotname
|	'(' ntype ')'

convtype:
	fntype
|	othertype

comptype:
	othertype

fnret_type:
	recvchantype
|	fntype
|	othertype
|	ptrtype
|	dotname

dotname:
	name
|	name '.' sym

othertype:
	'[' oexpr ']' ntype
|	'[' LDDD ']' ntype
|	LCHAN non_recvchantype
|	LCHAN LCOMM ntype
|	LMAP '[' ntype ']' ntype
|	structtype
|	interfacetype

ptrtype:
	'*' ntype

recvchantype:
	LCOMM LCHAN ntype

structtype:
	LSTRUCT lbrace structdcl_list osemi '}'
|	LSTRUCT lbrace '}'

interfacetype:
	LINTERFACE lbrace interfacedcl_list osemi '}'
|	LINTERFACE lbrace '}'

/*
 * function stuff
 * all in one place to show how crappy it all is
 */
xfndcl:
	LFUNC fndcl fnbody

fndcl:
	sym '(' oarg_type_list_ocomma ')' fnres
|	'(' oarg_type_list_ocomma ')' sym '(' oarg_type_list_ocomma ')' fnres

hidden_fndcl:
	hidden_pkg_importsym '(' ohidden_funarg_list ')' ohidden_funres
|	'(' hidden_funarg_list ')' sym '(' ohidden_funarg_list ')' ohidden_funres

fntype:
	LFUNC '(' oarg_type_list_ocomma ')' fnres

fnbody:
|	'{' stmt_list '}'

fnres:
	%prec NotParen
|	fnret_type
|	'(' oarg_type_list_ocomma ')'

fnlitdcl:
	fntype

fnliteral:
	fnlitdcl lbrace stmt_list '}'
|	fnlitdcl error

/*
 * lists of things
 * note that they are left recursive
 * to conserve yacc stack. they need to
 * be reversed to interpret correctly
 */
xdcl_list:
|	xdcl_list xdcl ';'

vardcl_list:
	vardcl
|	vardcl_list ';' vardcl

constdcl_list:
	constdcl1
|	constdcl_list ';' constdcl1

typedcl_list:
	typedcl
|	typedcl_list ';' typedcl

structdcl_list:
	structdcl
|	structdcl_list ';' structdcl

interfacedcl_list:
	interfacedcl
|	interfacedcl_list ';' interfacedcl

structdcl:
	new_name_list ntype oliteral
|	embed oliteral
|	'(' embed ')' oliteral
|	'*' embed oliteral
|	'(' '*' embed ')' oliteral
|	'*' '(' embed ')' oliteral

packname:
	LNAME
|	LNAME '.' sym

embed:
	packname

interfacedcl:
	new_name indcl
|	packname
|	'(' packname ')'

indcl:
	'(' oarg_type_list_ocomma ')' fnres

/*
 * function arguments.
 */
arg_type:
	name_or_type
|	sym name_or_type
|	sym dotdotdot
|	dotdotdot

arg_type_list:
	arg_type
|	arg_type_list ',' arg_type

oarg_type_list_ocomma:
|	arg_type_list ocomma

/*
 * statement
 */
stmt:
|	compound_stmt
|	common_dcl
|	non_dcl_stmt
|	error

non_dcl_stmt:
	simple_stmt
|	for_stmt
|	switch_stmt
|	select_stmt
|	if_stmt
|	labelname ':'
	stmt
|	LFALL
|	LBREAK onew_name
|	LCONTINUE onew_name
|	LGO pseudocall
|	LDEFER pseudocall
|	LGOTO new_name
|	LRETURN oexpr_list

stmt_list:
	stmt
|	stmt_list ';' stmt

new_name_list:
	new_name
|	new_name_list ',' new_name

dcl_name_list:
	dcl_name
|	dcl_name_list ',' dcl_name

expr_list:
	expr
|	expr_list ',' expr

expr_or_type_list:
	expr_or_type
|	expr_or_type_list ',' expr_or_type

/*
 * list of combo of keyval and val
 */
keyval_list:
	keyval
|	bare_complitexpr
|	keyval_list ',' keyval
|	keyval_list ',' bare_complitexpr

braced_keyval_list:
|	keyval_list ocomma

/*
 * optional things
 */
osemi:
|	';'

ocomma:
|	','

oexpr:
|	expr

oexpr_list:
|	expr_list

osimple_stmt:
|	simple_stmt

ohidden_funarg_list:
|	hidden_funarg_list

ohidden_structdcl_list:
|	hidden_structdcl_list

ohidden_interfacedcl_list:
|	hidden_interfacedcl_list

oliteral:
|	LLITERAL

/*
 * import syntax from package header
 */
hidden_import:
	LIMPORT LNAME LLITERAL ';'
|	LVAR hidden_pkg_importsym hidden_type ';'
|	LCONST hidden_pkg_importsym '=' hidden_constant ';'
|	LCONST hidden_pkg_importsym hidden_type '=' hidden_constant ';'
|	LTYPE hidden_pkgtype hidden_type ';'
|	LFUNC hidden_fndcl fnbody ';'

hidden_pkg_importsym:
	hidden_importsym

hidden_pkgtype:
	hidden_pkg_importsym

/*
 *  importing types
 */

hidden_type:
	hidden_type_misc
|	hidden_type_recv_chan
|	hidden_type_func

hidden_type_non_recv_chan:
	hidden_type_misc
|	hidden_type_func

hidden_type_misc:
	hidden_importsym
|	LNAME
|	'[' ']' hidden_type
|	'[' LLITERAL ']' hidden_type
|	LMAP '[' hidden_type ']' hidden_type
|	LSTRUCT '{' ohidden_structdcl_list '}'
|	LINTERFACE '{' ohidden_interfacedcl_list '}'
|	'*' hidden_type
|	LCHAN hidden_type_non_recv_chan
|	LCHAN '(' hidden_type_recv_chan ')'
|	LCHAN LCOMM hidden_type

hidden_type_recv_chan:
	LCOMM LCHAN hidden_type

hidden_type_func:
	LFUNC '(' ohidden_funarg_list ')' ohidden_funres

hidden_funarg:
	sym hidden_type oliteral
|	sym LDDD hidden_type oliteral

hidden_structdcl:
	sym hidden_type oliteral

hidden_interfacedcl:
	sym '(' ohidden_funarg_list ')' ohidden_funres
|	hidden_type

ohidden_funres:
|	hidden_funres

hidden_funres:
	'(' ohidden_funarg_list ')'
|	hidden_type

/*
 *  importing constants
 */

hidden_literal:
	LLITERAL
|	'-' LLITERAL
|	sym

hidden_constant:
	hidden_literal
|	'(' hidden_literal '+' hidden_literal ')'

hidden_import_list:
|	hidden_import_list hidden_import

hidden_funarg_list:
	hidden_funarg
|	hidden_funarg_list ',' hidden_funarg

hidden_structdcl_list:
	hidden_structdcl
|	hidden_structdcl_list ';' hidden_structdcl

hidden_interfacedcl_list:
	hidden_interfacedcl
|	hidden_interfacedcl_list ';' hidden_interfacedcl

%%

/* static void */
/* fixlbrace(int lbr) */
/* { */
/* 	// If the opening brace was an LBODY, */
/* 	// set up for another one now that we're done. */
/* 	// See comment in lex.c about loophack. */
/* 	if(lbr == LBODY) */
/* 		loophack = 1; */
/* } */
