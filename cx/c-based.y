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

	argument *CXArgument
	arguments []*CXArgument

	expression *CXExpression
	expressions []*CXExpression
}

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
%type   <tok>           type_specifier

%start translation_unit
%%

translation_unit:
                external_declaration
	|       translation_unit external_declaration
                ;

external_declaration:
                package_definition
	|       declaration
                ;

primary_expression:
                IDENTIFIER
                STRING_LITERAL
                BOOLEAN_LITERAL
                BYTE_LITERAL
                INT_LITERAL
                FLOAT_LITERAL
                DOUBLE_LITERAL
                LONG_LITERAL
        | '('   expression ')'
                ;

enumeration_constant:
                IDENTIFIER
                ;

postfix_expression:
                primary_expression
	|       postfix_expression '[' expression ']'
	|       postfix_expression '(' ')'
	|       postfix_expression '(' argument_expression_list ')'
	|       postfix_expression '.' IDENTIFIER
	|       postfix_expression PTR_OP IDENTIFIER // check
	|       postfix_expression INC_OP
	|       postfix_expression DEC_OP
                ;

argument_expression_list:
                assignment_expression
	|       argument_expression_list ',' assignment_expression
                ;

unary_expression:
                postfix_expression
	|       INC_OP unary_expression
	|       DEC_OP unary_expression
	|       unary_operator unary_expression // check
                ;

unary_operator:
                '&'
	|       '*'
	|       '+'
	|       '-'
	|       '~' // check
	|       '!'
                ;

multiplicative_expression:
                unary_expression
	|       multiplicative_expression '*' unary_expression
	|       multiplicative_expression '/' unary_expression
	|       multiplicative_expression '%' unary_expression
                ;

additive_expression:
                multiplicative_expression
	|       additive_expression '+' multiplicative_expression
	|       additive_expression '-' multiplicative_expression
                ;

shift_expression:
                additive_expression
	|       shift_expression LEFT_OP additive_expression
	|       shift_expression RIGHT_OP additive_expression
                ;

relational_expression:
                shift_expression
	|       relational_expression '<' shift_expression
	|       relational_expression '>' shift_expression
	|       relational_expression LE_OP shift_expression
	|       relational_expression GE_OP shift_expression
                ;

equality_expression:
                relational_expression
	|       equality_expression EQ_OP relational_expression
	|       equality_expression NE_OP relational_expression
                ;

and_expression: equality_expression
	|       and_expression '&' equality_expression
                ;

exclusive_or_expression:
                and_expression
	|       exclusive_or_expression '^' and_expression
                ;

inclusive_or_expression:
                exclusive_or_expression
	|       inclusive_or_expression '|' exclusive_or_expression
                ;

logical_and_expression:
                inclusive_or_expression
	|       logical_and_expression AND_OP inclusive_or_expression
                ;

logical_or_expression:
                logical_and_expression
	|       logical_or_expression OR_OP logical_and_expression
                ;

conditional_expression:
                logical_or_expression
	|       logical_or_expression '?' expression ':' conditional_expression
                ;

assignment_expression:
                conditional_expression
	|       unary_expression assignment_operator assignment_expression
                ;

assignment_operator:
                '='
	|       MUL_ASSIGN
	|       DIV_ASSIGN
	|       MOD_ASSIGN
	|       ADD_ASSIGN
	|       SUB_ASSIGN
	|       LEFT_ASSIGN
	|       RIGHT_ASSIGN
	|       AND_ASSIGN
	|       XOR_ASSIGN
	|       OR_ASSIGN
                ;

expression:     assignment_expression
	|       expression ',' assignment_expression
                ;

constant_expression:
                conditional_expression
                ;






// includes function declaration
declaration:    declaration_specifiers ';'
	|       declaration_specifiers init_declarator_list ';'
                {fmt.Println("declaration")}
                ;

declaration_specifiers:
		type_specifier declaration_specifiers
	|       type_specifier
	|       type_qualifier declaration_specifiers
	|       type_qualifier
                ;

// I think that a function declaration is declared in here
init_declarator_list:
                init_declarator
	|       init_declarator_list ',' init_declarator
                ;

init_declarator:
                declarator '=' initializer
	|       declarator
                ;

type_specifier:
                BOOL
        |       BYTE
        |       F32
        |       F64
        |       I8
        |       I16
        |       I32
        |       I64
        |       STR
        |       UI8
        |       UI16
        |       UI32
        |       UI64
	|       struct_or_union_specifier
                {
                    $$ = "struct"
                }
	|       enum_specifier
                {
                    $$ = "enum"
                }
	/* |       TYPEDEF_NAME // check */
                ;

struct_or_union_specifier:
                struct_or_union '{' struct_declaration_list '}'
	|       struct_or_union IDENTIFIER '{' struct_declaration_list '}'
	|       struct_or_union IDENTIFIER
                ;

struct_or_union:
                STRUCT
	|       UNION // check
                ;

struct_declaration_list:
                struct_declaration
	|       struct_declaration_list struct_declaration
                ;

struct_declaration:
                specifier_qualifier_list ';'	/* for anonymous struct/union */
	|       specifier_qualifier_list struct_declarator_list ';'
                ;

specifier_qualifier_list:
                type_specifier specifier_qualifier_list
	|       type_specifier
	|       type_qualifier specifier_qualifier_list
	|       type_qualifier
                ;

struct_declarator_list:
                struct_declarator
	|       struct_declarator_list ',' struct_declarator
                ;

struct_declarator:
                ':' constant_expression
	|       declarator ':' constant_expression
	|       declarator
                ;

enum_specifier: ENUM '{' enumerator_list '}'
	|       ENUM '{' enumerator_list ',' '}'
	|       ENUM IDENTIFIER '{' enumerator_list '}'
	|       ENUM IDENTIFIER '{' enumerator_list ',' '}'
	|       ENUM IDENTIFIER
                ;

enumerator_list:enumerator
	|       enumerator_list ',' enumerator
                ;

enumerator:     enumeration_constant '=' constant_expression
	|       enumeration_constant
                ;






type_qualifier: CONST
        /* |       VAR */
                ;

declarator:     pointer direct_declarator
	|       direct_declarator
                ;

direct_declarator:
                IDENTIFIER
	| '('   declarator ')'
	|       direct_declarator '[' ']'
	|       direct_declarator '[' '*' ']'
	|       direct_declarator '[' type_qualifier_list '*' ']'
	|       direct_declarator '[' type_qualifier_list assignment_expression ']'
	|       direct_declarator '[' type_qualifier_list ']'
	|       direct_declarator '[' assignment_expression ']'
	|       direct_declarator '(' parameter_type_list ')'
	|       direct_declarator '(' ')'
	|       direct_declarator '(' identifier_list ')'
                ;

pointer:  '*'   type_qualifier_list pointer
	| '*'   type_qualifier_list
	| '*'   pointer
	| '*'
	;

type_qualifier_list:
                type_qualifier
	|       type_qualifier_list type_qualifier
                ;







parameter_type_list:
		parameter_list
                ;

parameter_list:
                parameter_declaration
	|       parameter_list ',' parameter_declaration
                ;

parameter_declaration:
                declaration_specifiers declarator
	|       declaration_specifiers
                ;

identifier_list:
                IDENTIFIER
	|       identifier_list ',' IDENTIFIER
                ;









initializer:
                '{' initializer_list '}'
	| '{'   initializer_list ',' '}'
	|       assignment_expression
                ;

initializer_list:
                designation initializer
	|       initializer
	|       initializer_list ',' designation initializer
	|       initializer_list ',' initializer
                ;

designation:    designator_list '='
                ;

designator_list:
                designator
	|       designator_list designator
                ;

designator:
                '[' constant_expression ']'
	| '.'   IDENTIFIER
                ;





statement:      labeled_statement
	|       compound_statement
	|       expression_statement
	|       selection_statement
	|       iteration_statement
	|       jump_statement
                ;

labeled_statement:
                IDENTIFIER ':' statement
	|       CASE constant_expression ':' statement
	|       DEFAULT ':' statement
                ;

compound_statement:
                '{' '}'
	| '{'   block_item_list '}'
                ;

block_item_list:
                block_item
	|       block_item_list block_item
                ;

block_item:     declaration
	|       statement
                ;

expression_statement:
                ';'
	|       expression ';'
                ;

selection_statement:
                IF '(' expression ')' statement ELSE statement
	|       IF '(' expression ')' statement
	|       SWITCH '(' expression ')' statement
                ;

iteration_statement:
                FOR '(' expression_statement expression_statement ')' statement
	|       FOR '(' expression_statement expression_statement expression ')' statement
	|       FOR '(' declaration expression_statement ')' statement
	|       FOR '(' declaration expression_statement expression ')' statement
                ;

jump_statement: GOTO IDENTIFIER ';'
	|       CONTINUE ';'
	|       BREAK ';'
	|       RETURN ';'
	|       RETURN expression ';'
                ;








package_definition:
                function_definition
        |       PACKAGE IDENTIFIER
                {
			pkg := MakeModule($2)
			prgrm.AddModule(pkg)
                }
                ;

function_definition:
                declaration_specifiers declarator declaration_list compound_statement
                {
                    fmt.Println("again")
                }
	/* |       FUNC declarator compound_statement */
        /*         { */
        /*             fmt.Println("again") */
        /*         } */
                ;

declaration_list:
                declaration
	|       declaration_list declaration
                ;

%%
