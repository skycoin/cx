%{
	package partialparsing
	import (
		"bytes"
	    "github.com/skycoin/cx/cx/ast"
	    "github.com/skycoin/cx/cx/constants"
		"github.com/skycoin/cx/cxparser/actions"
	)

        /*
        This is a computer generated file
        - partialparsing.go is generated from partialparsing.y via 
        - ./bin/goyacc -o cxparser/cxpartialparsing/partialparsing.go cxparser/cxpartialparsing/partialparsing.y
        */

	var Program *ast.CXProgram

	var lineNo int = -1
	var replMode bool = false
	var inREPL bool = false
	var inFn bool = false
	var fileName string

	// Parse() is the function that is called from main().
	// It is needed because yyParse is not exported.
	func Parse(code string) int {
		codeBuf := bytes.NewBufferString(code)
		return yyParse(NewLexer(codeBuf))
	}

	func PreFunctionDeclaration (fn *ast.CXFunction, inputs []*ast.CXArgument, outputs []*ast.CXArgument) {
		// adding inputs, outputs
		for _, inp := range inputs {
			fn.AddInput(inp)
		}
		for _, out := range outputs {
			fn.AddOutput(out)
		}
	}
	
%}

%union {
	i int
	i8 int8
	i16 int16
	i32 int32
	i64 int64
	ui8 uint8
	ui16 uint16
	ui32 uint32
	ui64 uint64
	f32 float32
	f64 float64
	tok string
	bool bool
	string string
	stringA []string
	ints    []int

	line int

	argument *ast.CXArgument
	arguments []*ast.CXArgument

	expression *ast.CXExpression
	expressions []*ast.CXExpression

        function *ast.CXFunction
}

%token  <bool>          BOOLEAN_LITERAL
%token  <i8>            BYTE_LITERAL
%token  <i16>           SHORT_LITERAL
%token  <i32>           INT_LITERAL
%token  <i64>           LONG_LITERAL
%token  <ui8>           UNSIGNED_BYTE_LITERAL
%token  <ui16>          UNSIGNED_SHORT_LITERAL
%token  <ui32>          UNSIGNED_INT_LITERAL
%token  <ui64>          UNSIGNED_LONG_LITERAL
%token  <f32>           FLOAT_LITERAL
%token  <f64>           DOUBLE_LITERAL
%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK IDENTIFIER
                        VAR COMMA PERIOD COMMENT STRING_LITERAL PACKAGE IF ELSE FOR TYPSTRUCT STRUCT
                        SEMICOLON NEWLINE
                        ASSIGN CASSIGN IMPORT RETURN GOTO GT_OP LT_OP GTEQ_OP LTEQ_OP EQUAL COLON NEW
                        EQUALWORD GTHANWORD LTHANWORD
                        GTHANEQ LTHANEQ UNEQUAL AND OR
                        ADD_OP SUB_OP MUL_OP DIV_OP MOD_OP REF_OP NEG_OP AFFVAR
                        PLUSPLUS MINUSMINUS REMAINDER LEFTSHIFT RIGHTSHIFT EXP
                        NOT
                        BITXOR_OP BITOR_OP BITCLEAR_OP
                        PLUSEQ MINUSEQ MULTEQ DIVEQ REMAINDEREQ EXPEQ
                        LEFTSHIFTEQ RIGHTSHIFTEQ BITANDEQ BITXOREQ BITOREQ

                        DEC_OP INC_OP PTR_OP LEFT_OP RIGHT_OP
                        GE_OP LE_OP EQ_OP NE_OP AND_OP OR_OP
                        ADD_ASSIGN AND_ASSIGN LEFT_ASSIGN MOD_ASSIGN
                        MUL_ASSIGN DIV_ASSIGN OR_ASSIGN RIGHT_ASSIGN
                        SUB_ASSIGN XOR_ASSIGN
                        BOOL F32 F64
                        I8 I16 I32 I64
                        STR
                        UI8 UI16 UI32 UI64
                        UNION ENUM CONST CASE DEFAULT SWITCH BREAK CONTINUE
                        TYPE
                        
                        /* Types */
                        BASICTYPE

                        /* Removers */
                        REM DEF EXPR FIELD CLAUSES OBJECT OBJECTS

                        /* Debugging */
                        DSTACK DPROGRAM DSTATE
                        /* Affordances */
                        AFF CAFF TAG INFER VALUE
                        /* Pointers */
                        ADDR

%type   <i32>           int_value

%type   <tok>           unary_operator
%type   <i>             type_specifier
%type   <argument>      declaration_specifiers
%type   <argument>      declarator
%type   <argument>      direct_declarator
%type   <argument>      parameter_declaration
%type   <arguments>     parameter_type_list
%type   <arguments>     function_parameters
%type   <arguments>     parameter_list
%type   <arguments>     fields
%type   <arguments>     struct_fields
%type   <arguments>     id_list
%type   <arguments>     types_list
                                                
%type   <function>      function_header

%type   <ints>          indexing_literal
%type   <ints>          indexing_slice_literal

                        // for struct literals
%right                  IDENTIFIER LBRACE

//%start
                        
%%

translation_unit:
                external_declaration
        |       translation_unit external_declaration
        ;

external_declaration:
                package_declaration
        |       global_declaration
        |       function_declaration
        |       import_declaration
        |       struct_declaration

       // |       stepping
        ;

/*
stepping:       TSTEP int_value int_value
        |       STEP int_value
        ;
*/

global_declaration:
                VAR declarator declaration_specifiers SEMICOLON
                {
			actions.DeclareGlobal($2, $3, nil, false)
                }
        |       VAR declarator declaration_specifiers ASSIGN initializer SEMICOLON
                {
			actions.DeclareGlobal($2, $3, nil, false)
                }
                ;

struct_declaration:
                TYPE IDENTIFIER STRUCT struct_fields
                {
			actions.DeclareStruct($2, $4)
                }
                ;

struct_fields:
                LBRACE RBRACE SEMICOLON
                { $$ = nil }
        |       LBRACE fields RBRACE SEMICOLON
                { $$ = $2 }
        ;

fields:         parameter_declaration SEMICOLON
                {
			$$ = []*ast.CXArgument{$1}
                }
        |       fields parameter_declaration SEMICOLON
                {
			$$ = append($1, $2)
                }
        ;

package_declaration:
                PACKAGE IDENTIFIER SEMICOLON
                {
			actions.DeclarePackage($2)
                }
                ;

import_declaration:
                IMPORT STRING_LITERAL SEMICOLON
                {
			actions.DeclareImport($2, CurrentFileName, lineNo)
                }
                ;

function_header:
                FUNC IDENTIFIER
                {
			if pkg, err := Program.GetCurrentPackage(); err == nil {
				fn := ast.MakeFunction($2, CurrentFileName, lineNo)
				pkg.AddFunction(fn)

                                $$ = fn
			} else {
				panic(err)
			}
                }
        |       FUNC LPAREN parameter_type_list RPAREN IDENTIFIER
                {
			if len($3) > 1 {
				panic("method has multiple receivers")
			}

			fnName := $3[0].CustomType.Name + "." + $5

			if pkg, err := Program.GetCurrentPackage(); err == nil {
				fn := ast.MakeFunction(fnName, CurrentFileName, lineNo)
				pkg.AddFunction(fn)

                                fn.AddInput($3[0])

                                $$ = fn
			} else {
				panic(err)
			}
                }
        ;

function_parameters:
                LPAREN RPAREN
                { $$ = nil }
        |       LPAREN parameter_type_list RPAREN
                { $$ = $2 }
                ;

function_declaration:
                function_header function_parameters compound_statement
                {
			PreFunctionDeclaration($1, $2, nil)
                }
        |       function_header function_parameters function_parameters compound_statement
                {
			PreFunctionDeclaration($1, $2, $3)
                }
        ;

parameter_type_list:
                //parameter_list COMMA ELLIPSIS
		parameter_list
                ;

parameter_list:
                parameter_declaration
                {
			$$ = []*ast.CXArgument{$1}
                }
	|       parameter_list COMMA parameter_declaration
                {
			$$ = append($1, $3)
                }
                ;

parameter_declaration:
                declarator declaration_specifiers
                {
			$2.ArgDetails.Name = $1.ArgDetails.Name
			$2.ArgDetails.Package = $1.ArgDetails.Package
			$2.IsLocalDeclaration = true
			$$ = $2
                }
                ;

identifier_list:
                IDENTIFIER
	|       identifier_list COMMA IDENTIFIER
                ;

declarator:     direct_declarator
                ;

direct_declarator:
                IDENTIFIER
                {
			if pkg, err := Program.GetCurrentPackage(); err == nil {
				arg := ast.MakeArgument("", actions.CurrentFile, actions.LineNo)
				arg.AddType(constants.TypeNames[constants.TYPE_UNDEFINED])
				arg.ArgDetails.Name = $1
				arg.ArgDetails.Package = pkg
				$$ = arg
			} else {
				panic(err)
			}
                }
	|       LPAREN declarator RPAREN
                { $$ = $2 }
                ;

id_list:	IDENTIFIER
		{
			arg := actions.DeclarationSpecifiersStruct($1, "", false, actions.CurrentFile, actions.LineNo)
			$$ = []*ast.CXArgument{arg}
		}
	|	type_specifier
		{
			arg := actions.DeclarationSpecifiersBasic($1)			 
			$$ = []*ast.CXArgument{arg}
		}
	|	id_list COMMA IDENTIFIER
		{
			arg := actions.DeclarationSpecifiersStruct($3, "", false, actions.CurrentFile, actions.LineNo)
			$$ = append($1, arg)
		}
	|	id_list COMMA type_specifier
		{
			arg := actions.DeclarationSpecifiersBasic($3)
			$$ = append($1, arg)
		}
	;

types_list:
		LPAREN id_list RPAREN
		{
			$$ = $2
		}
	|	LPAREN RPAREN
		{
			$$ = nil
		}
	;

declaration_specifiers:
		FUNC types_list types_list
		{
			arg := ast.MakeArgument("", actions.CurrentFile, actions.LineNo).AddType("func")
			arg.Inputs = $2
			arg.Outputs = $3
			$$ = actions.DeclarationSpecifiers(arg, []int{0}, constants.DECL_FUNC)
		}
        |       MUL_OP declaration_specifiers
                {
			$$ = actions.DeclarationSpecifiers($2, []int{0}, constants.DECL_POINTER)
                }
        // |       indexing_literal declaration_specifiers
        //         {
	// 		$$ = actions.DeclarationSpecifiers($2, $1, constants.DECL_ARRAY)
        //         }
        |       LBRACK RBRACK declaration_specifiers
                {
			$$ = actions.DeclarationSpecifiers($3, []int{0}, constants.DECL_SLICE)
                }
        |       type_specifier
                {
			$$ = actions.DeclarationSpecifiersBasic($1)
                }
        |       IDENTIFIER
                {
			$$ = actions.DeclarationSpecifiersStruct($1, "", false, CurrentFileName, lineNo)
                }
        |       indexing_literal type_specifier
                {
			basic := actions.DeclarationSpecifiersBasic($2)
			$$ = actions.DeclarationSpecifiers(basic, $1, constants.DECL_ARRAY)
                }
        |       indexing_literal IDENTIFIER
                {
			strct := actions.DeclarationSpecifiersStruct($2, "", false, actions.CurrentFile, actions.LineNo)
			$$ = actions.DeclarationSpecifiers(strct, $1, constants.DECL_ARRAY)
                }
        |       IDENTIFIER PERIOD IDENTIFIER
                {
			$$ = actions.DeclarationSpecifiersStruct($3, $1, true, CurrentFileName, lineNo)
                }
	|       type_specifier PERIOD IDENTIFIER
		{
			$$ = actions.DeclarationSpecifiersStruct($3, constants.TypeNames[$1], true, CurrentFileName, lineNo)
		}
		/* type_specifier declaration_specifiers */
	/* |       type_specifier */
	/* |       type_qualifier declaration_specifiers */
	/* |       type_qualifier */
                ;

type_specifier:
                AFF
                { $$ = constants.TYPE_AFF }
        |       BOOL
                { $$ = constants.TYPE_BOOL }
        |       STR
                { $$ = constants.TYPE_STR }
        |       F32
                { $$ = constants.TYPE_F32 }
        |       F64
                { $$ = constants.TYPE_F64 }
        |       I8
                { $$ = constants.TYPE_I8 }
        |       I16
                { $$ = constants.TYPE_I16 }
        |       I32
                { $$ = constants.TYPE_I32 }
        |       I64
                { $$ = constants.TYPE_I64 }
        |       UI8
                { $$ = constants.TYPE_UI8 }
        |       UI16
                { $$ = constants.TYPE_UI16 }
        |       UI32
                { $$ = constants.TYPE_UI32 }
        |       UI64
                { $$ = constants.TYPE_UI64 }
	/* |       struct_or_union_specifier */
        /*         { */
        /*             $$ = "struct" */
        /*         } */
	/* |       enum_specifier */
        /*         { */
        /*             $$ = "enum" */
        /*         } */
	/* |       TYPEDEF_NAME // check */
                ;



struct_literal_fields:
                // empty
        |       IDENTIFIER COLON constant_expression
        |       struct_literal_fields COMMA IDENTIFIER COLON constant_expression
                ;

array_literal_expression_list:
		assignment_expression
        |       LBRACE array_literal_expression_list RBRACE
	|       array_literal_expression_list COMMA assignment_expression
                ;

indexing_literal:
		LBRACK INT_LITERAL RBRACK
		{
			$$ = []int{int($2)}
		}
        |       indexing_literal LBRACK INT_LITERAL RBRACK
		{
			$$ = append($1, int($3))
		}
		;

indexing_slice_literal:
		LBRACK RBRACK
		{
			$$ = []int{0}
		}
        |       indexing_slice_literal LBRACK RBRACK
		{
			$$ = append($1, 0)
		}
		;

// expressions
array_literal_expression:
                indexing_literal IDENTIFIER LBRACE array_literal_expression_list RBRACE
        |       indexing_literal IDENTIFIER LBRACE RBRACE
        |       indexing_literal type_specifier LBRACE array_literal_expression_list RBRACE
        |       indexing_literal type_specifier LBRACE RBRACE
        ;

slice_literal_expression:
                LBRACK RBRACK IDENTIFIER LBRACE argument_expression_list RBRACE
        |       LBRACK RBRACK IDENTIFIER LBRACE RBRACE
        |       LBRACK RBRACK type_specifier LBRACE argument_expression_list RBRACE
        |       LBRACK RBRACK type_specifier LBRACE RBRACE
        |       LBRACK RBRACK slice_literal_expression
                ;



/* infer_action_arg: */
/*                 MUL_OP GT_OP assignment_expression */
/*         |       MUL_OP GT_OP MUL_OP */
/*         ; */

infer_action_arg:
                IDENTIFIER
        |       INT_LITERAL
	|	type_specifier PERIOD IDENTIFIER
        ;

infer_action:
                IDENTIFIER LPAREN infer_action_arg COMMA IDENTIFIER RPAREN
	|	IDENTIFIER LPAREN infer_action_arg RPAREN
	|	IDENTIFIER LPAREN infer_action RPAREN
	|	IDENTIFIER LPAREN infer_action COMMA infer_action RPAREN
        ;

infer_actions:
                infer_action SEMICOLON
        |       infer_actions infer_action SEMICOLON
                ;

/* infer_target: */
/*                 IDENTIFIER LPAREN IDENTIFIER RPAREN SEMICOLON */
/*         ; */

/* infer_targets: */
/*                 infer_target */
/*         |       infer_targets infer_target */
/*         ; */

infer_clauses:
        |       infer_actions
        /* |       infer_targets */
                ;

int_value:
		INT_LITERAL
		{
			$$ = $1
		}
        |       SUB_OP INT_LITERAL
		{
			$$ = -$2
		}

primary_expression:
                IDENTIFIER
        /* |       IDENTIFIER LBRACE struct_literal_fields RBRACE */
        |       INFER LBRACE infer_clauses RBRACE
        |       STRING_LITERAL
        |       BOOLEAN_LITERAL
        |       BYTE_LITERAL
        |       SHORT_LITERAL
        |       INT_LITERAL
        |       LONG_LITERAL
        |       UNSIGNED_BYTE_LITERAL
        |       UNSIGNED_SHORT_LITERAL
        |       UNSIGNED_INT_LITERAL
        |       UNSIGNED_LONG_LITERAL
        |       FLOAT_LITERAL
        |       DOUBLE_LITERAL
        |       LPAREN expression RPAREN
        |       array_literal_expression
        |       slice_literal_expression
                ;

after_period:   type_specifier
        |       IDENTIFIER
        ;

postfix_expression:
                primary_expression
	|       postfix_expression LBRACK expression RBRACK
        |       type_specifier PERIOD after_period
	|       postfix_expression LPAREN RPAREN
	|       postfix_expression LPAREN argument_expression_list RPAREN
	|       postfix_expression INC_OP
        |       postfix_expression DEC_OP
        |       postfix_expression PERIOD IDENTIFIER
        /* |       postfix_expression PERIOD IDENTIFIER LBRACE struct_literal_fields RBRACE */
                ;

argument_expression_list:
                assignment_expression
	|       argument_expression_list COMMA assignment_expression
                ;

unary_expression:
                postfix_expression
	|       INC_OP unary_expression
	|       DEC_OP unary_expression
	|       unary_operator unary_expression // check
                ;

unary_operator:
                REF_OP
	|       MUL_OP
	|       ADD_OP
	|       SUB_OP
	|       NEG_OP
                ;

multiplicative_expression:
                unary_expression
	|       multiplicative_expression MUL_OP unary_expression
	|       multiplicative_expression DIV_OP unary_expression
	|       multiplicative_expression MOD_OP unary_expression
                ;

additive_expression:
                multiplicative_expression
	|       additive_expression ADD_OP multiplicative_expression
	|       additive_expression SUB_OP multiplicative_expression
                ;

shift_expression:
                additive_expression
	|       shift_expression LEFT_OP additive_expression
	|       shift_expression RIGHT_OP additive_expression
        |       shift_expression BITCLEAR_OP additive_expression
                ;

relational_expression:
                shift_expression
	|       relational_expression EQ_OP shift_expression
	|       relational_expression NE_OP shift_expression
	|       relational_expression LT_OP shift_expression
	|       relational_expression GT_OP shift_expression
	|       relational_expression LTEQ_OP shift_expression
	|       relational_expression GTEQ_OP shift_expression
                ;

and_expression: relational_expression
	|       and_expression REF_OP relational_expression
                ;

exclusive_or_expression:
                and_expression
	|       exclusive_or_expression BITXOR_OP and_expression
                ;

inclusive_or_expression:
                exclusive_or_expression
	|       inclusive_or_expression BITOR_OP exclusive_or_expression
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
	|       logical_or_expression '?' expression COLON conditional_expression
                ;

struct_literal_expression:
                conditional_expression
	|       IDENTIFIER LBRACE struct_literal_fields RBRACE
	|       unary_operator IDENTIFIER LBRACE struct_literal_fields RBRACE
        |       postfix_expression PERIOD IDENTIFIER LBRACE struct_literal_fields RBRACE
        ;

assignment_expression:
                /* conditional_expression */
                struct_literal_expression
	|       unary_expression assignment_operator assignment_expression
                ;

assignment_operator:
                ASSIGN
        |       CASSIGN
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
	|       expression COMMA assignment_expression
                ;

constant_expression:
                conditional_expression
                ;

declaration:
                VAR declarator declaration_specifiers SEMICOLON
        |       VAR declarator declaration_specifiers ASSIGN initializer SEMICOLON
                ;

initializer:
                assignment_expression
                ;

// statements
statement:      labeled_statement
        |       compound_statement
	|       expression_statement
	|       selection_statement
	|       iteration_statement
        |       jump_statement
                ;

labeled_statement:
                IDENTIFIER COLON block_item
	|       CASE constant_expression COLON statement
	|       DEFAULT COLON statement
                ;

compound_statement:
                LBRACE RBRACE SEMICOLON
	|       LBRACE block_item_list RBRACE SEMICOLON
                ;

block_item_list:
                block_item
	|       block_item_list block_item
                ;

block_item:     declaration
        |       statement
                ;

expression_statement:
                SEMICOLON
	|       expression SEMICOLON
                ;

selection_statement:
                IF conditional_expression LBRACE block_item_list RBRACE elseif_list else_statement SEMICOLON
        |       IF conditional_expression LBRACE block_item_list RBRACE else_statement SEMICOLON
        |       IF conditional_expression LBRACE RBRACE else_statement SEMICOLON //
        |       IF conditional_expression LBRACE block_item_list RBRACE elseif_list SEMICOLON
        |       IF conditional_expression LBRACE RBRACE elseif_list SEMICOLON
        |       IF conditional_expression LBRACE RBRACE elseif_list else_statement SEMICOLON
        |       IF conditional_expression compound_statement
	|       SWITCH LPAREN conditional_expression RPAREN statement
                ;

elseif:         ELSE IF conditional_expression LBRACE block_item_list RBRACE
        |       ELSE IF conditional_expression LBRACE RBRACE
        ;

elseif_list:
                elseif
        |       elseif_list elseif
        ;

else_statement:
                ELSE LBRACE block_item_list RBRACE
	|      	ELSE LBRACE RBRACE
        ;

iteration_statement:
                FOR expression compound_statement
        |       FOR expression_statement expression_statement compound_statement
        |       FOR expression_statement expression_statement expression compound_statement
        |       FOR declaration expression_statement compound_statement
        |       FOR declaration expression_statement expression compound_statement
                ;

jump_statement: GOTO IDENTIFIER SEMICOLON
	|       CONTINUE SEMICOLON
	|       BREAK SEMICOLON
	|       RETURN SEMICOLON
	|       RETURN expression SEMICOLON
                ;

%%
