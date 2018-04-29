%{
	package main
	import (
		"github.com/skycoin/skycoin/src/cipher/encoder"
		. "github.com/skycoin/cx/src/base"
	)
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

	selectStatement selectStatement
	selectStatements []selectStatement

	arrayArguments [][]*CXExpression

        function *CXFunction
}

%token  <byt>           BYTE_LITERAL
%token  <bool>          BOOLEAN_LITERAL
%token  <i32>           INT_LITERAL
%token  <i64>           LONG_LITERAL
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

%type   <tok>           after_period
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
                                                
%type   <expressions>   assignment_expression
%type   <expressions>   constant_expression
%type   <expressions>   conditional_expression
%type   <expressions>   logical_or_expression
%type   <expressions>   logical_and_expression
%type   <expressions>   exclusive_or_expression
%type   <expressions>   inclusive_or_expression
%type   <expressions>   and_expression
%type   <expressions>   equality_expression
%type   <expressions>   relational_expression
%type   <expressions>   shift_expression
%type   <expressions>   additive_expression
%type   <expressions>   multiplicative_expression
%type   <expressions>   unary_expression
%type   <expressions>   argument_expression_list
%type   <expressions>   postfix_expression
%type   <expressions>   primary_expression
                        
%type   <expressions>   array_literal_expression_list
%type   <expressions>   array_literal_expression

%type   <expressions>   struct_literal_fields
%type   <selectStatement>   elseif
%type   <selectStatements>   elseif_list

%type   <expressions>   declaration
//                      %type   <expressions>   init_declarator_list
//                      %type   <expressions>   init_declarator

%type   <expressions>   initializer
%type   <expressions>   designation
%type   <expressions>   designator_list
%type   <expressions>   designator

%type   <expressions>   expression
%type   <expressions>   block_item
%type   <expressions>   block_item_list
%type   <expressions>   compound_statement
%type   <expressions>   else_statement
%type   <expressions>   labeled_statement
%type   <expressions>   expression_statement
%type   <expressions>   selection_statement
%type   <expressions>   iteration_statement
%type   <expressions>   jump_statement
%type   <expressions>   statement

%type   <function>      function_header

                        // for struct literals
%right                   IDENTIFIER LBRACE
// %right                  IDENTIFIER
                        
/* %start                  translation_unit */
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
        ;

global_declaration:
                VAR declarator declaration_specifiers SEMICOLON
                {
			DeclareGlobal($2, $3, nil, false)
                }
        |       VAR declarator declaration_specifiers ASSIGN initializer SEMICOLON
                {
			DeclareGlobal($2, $2, $5, true)
                }
                ;

struct_declaration:
                TYPE IDENTIFIER STRUCT struct_fields
                {
			DeclareStruct($2, $4)
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
			$$ = []*CXArgument{$1}
                }
        |       fields parameter_declaration SEMICOLON
                {
			$$ = append($1, $2)
                }
        ;

package_declaration:
                PACKAGE IDENTIFIER SEMICOLON
                {
			DeclarePackage($2)
                }
                ;

import_declaration:
                IMPORT STRING_LITERAL SEMICOLON
                {
			DeclareImport($2)
                }
        ;

function_header:
                FUNC IDENTIFIER
                {
			$$ = FunctionHeader($2, nil, false)
                }
        |       FUNC LPAREN parameter_type_list RPAREN IDENTIFIER
                {
			$$ = FunctionHeader($5, $3, true)
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
			FunctionDeclaration($1, $2, nil, $3)
                }
        |       function_header function_parameters function_parameters compound_statement
                {
			FunctionDeclaration($1, $2, $3, $4)
                }
        ;

parameter_type_list:
		parameter_list
                ;

parameter_list:
                parameter_declaration
                {
			$$ = []*CXArgument{$1}
                }
	|       parameter_list COMMA parameter_declaration
                {
			$$ = append($1, $3)
                }
                ;

parameter_declaration:
                declarator declaration_specifiers
                {
			$2.Name = $1.Name
			$2.Package = $1.Package
			$$ = $2
                }
                ;

declarator:     direct_declarator
                ;

direct_declarator:
                IDENTIFIER
                {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				arg := MakeArgument(TYPE_UNDEFINED)
				arg.Name = $1
				arg.Package = pkg
				$$ = arg
			} else {
				panic(err)
			}
                }
	|       LPAREN declarator RPAREN
                { $$ = $2 }
	// |       direct_declarator '[' ']'
        //         {
	// 		$1.IsArray = true
	// 		$$ = $1
        //         }
        //	|direct_declarator '[' MUL_OP ']'
        //              	|direct_declarator '[' type_qualifier_list MUL_OP ']'
        //              	|direct_declarator '[' type_qualifier_list assignment_expression ']'
        //              	|direct_declarator '[' type_qualifier_list ']'
        //              	|direct_declarator '[' assignment_expression ']'
	// |    direct_declarator LPAREN parameter_type_list RPAREN
	// |    direct_declarator LPAREN RPAREN
	// |    direct_declarator LPAREN identifier_list RPAREN
                ;

// check
/* pointer:        /\* MUL_OP   type_qualifier_list pointer // check *\/ */
/*         /\* |       MUL_OP   type_qualifier_list // check *\/ */
/*         /\* |       MUL_OP   pointer *\/ */
/*         /\* |        *\/MUL_OP */
/*                 ; */

/* type_qualifier_list: */
/*                 type_qualifier */
/* 	|       type_qualifier_list type_qualifier */
/*                 ; */








declaration_specifiers:
                MUL_OP declaration_specifiers
                {
			$$ = DeclarationSpecifiers($2, 0, DECL_POINTER)
                }
        |       LBRACK INT_LITERAL RBRACK declaration_specifiers
                {
			$$ = DeclarationSpecifiers($4, int($2), DECL_ARRAY)
                }
        |       LBRACK RBRACK declaration_specifiers
                {
			$$ = DeclarationSpecifiers($3, 0, DECL_SLICE)
                }
        |       type_specifier
                {
			arg := MakeArgument($1)
			arg.Type = $1
			arg.Size = GetArgSize($1)
			$$ = DeclarationSpecifiers(arg, 0, DECL_BASIC)
                }
        |       IDENTIFIER
                {
			$$ = DeclarationSpecifiersStruct($1, "", false)
                }
        |       IDENTIFIER PERIOD IDENTIFIER
                {
			$$ = DeclarationSpecifiersStruct($3, $1, true)
                }
		/* type_specifier declaration_specifiers */
	/* |       type_specifier */
	/* |       type_qualifier declaration_specifiers */
	/* |       type_qualifier */
                ;

type_specifier:
                BOOL
                { $$ = TYPE_BOOL }
        |       BYTE
                { $$ = TYPE_BYTE }
        |       STR
                { $$ = TYPE_STR }
        |       F32
                { $$ = TYPE_F32 }
        |       F64
                { $$ = TYPE_F64 }
        |       I8
                { $$ = TYPE_I8 }
        |       I16
                { $$ = TYPE_I16 }
        |       I32
                { $$ = TYPE_I32 }
        |       I64
                { $$ = TYPE_I64 }
        |       UI8
                { $$ = TYPE_UI8 }
        |       UI16
                { $$ = TYPE_UI16 }
        |       UI32
                { $$ = TYPE_UI32 }
        |       UI64
                { $$ = TYPE_UI64 }
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
                { $$ = nil }
        |       IDENTIFIER COLON constant_expression
                {
			$$ = Assignment([]*CXExpression{StructLiteralFields($1)}, $3)
                }
        |       struct_literal_fields COMMA IDENTIFIER COLON constant_expression
                {
			$$ = append($1, Assignment([]*CXExpression{StructLiteralFields($3)}, $5)...)
                }
                ;

array_literal_expression_list:
                assignment_expression
                {
			$1[len($1) - 1].IsArrayLiteral = true
			$$ = $1
                }
	|       array_literal_expression_list COMMA assignment_expression
                {
			$3[len($3) - 1].IsArrayLiteral = true
			$$ = append($1, $3...)
                }
                ;

// expressions
array_literal_expression:
                LBRACK INT_LITERAL RBRACK IDENTIFIER LBRACE array_literal_expression_list RBRACE
                {
			$$ = $6
                }
        |       LBRACK INT_LITERAL RBRACK IDENTIFIER LBRACE RBRACE
                {
			$$ = nil
                }
        |       LBRACK INT_LITERAL RBRACK type_specifier LBRACE array_literal_expression_list RBRACE
                {
			$$ = ArrayLiteralExpression(int($2), $4, $6)
                }
        |       LBRACK INT_LITERAL RBRACK type_specifier LBRACE RBRACE
                {
			$$ = nil
                }
        |       LBRACK INT_LITERAL RBRACK array_literal_expression
                {
			for _, expr := range $4 {
				if expr.Outputs[0].Name == $4[len($4) - 1].Inputs[0].Name {
					expr.Outputs[0].Lengths = append([]int{int($2)}, expr.Outputs[0].Lengths[:len(expr.Outputs[0].Lengths) - 1]...)
					expr.Outputs[0].TotalSize = expr.Outputs[0].Size * TotalLength(expr.Outputs[0].Lengths)
				}
			}

			$$ = $4
                }
        ;

primary_expression:
                IDENTIFIER
                {
			$$ = PrimaryIdentifier($1)
                }
        |       IDENTIFIER LBRACE struct_literal_fields RBRACE
                {
			$$ = PrimaryStructLiteral($1, $3)
                }
        |       STRING_LITERAL
                {
			$$ = WritePrimary(TYPE_STR, encoder.Serialize($1))
                }
        |       BOOLEAN_LITERAL
                {
			$$ = WritePrimary(TYPE_BOOL, encoder.Serialize($1))
                }
        |       BYTE_LITERAL
                {
			$$ = WritePrimary(TYPE_BYTE, encoder.Serialize($1))
                }
        |       INT_LITERAL
                {
			$$ = WritePrimary(TYPE_I32, encoder.Serialize($1))
                }
        |       FLOAT_LITERAL
                {
			$$ = WritePrimary(TYPE_F32, encoder.Serialize($1))
                }
        |       DOUBLE_LITERAL
                {
			$$ = WritePrimary(TYPE_F64, encoder.Serialize($1))
                }
        |       LONG_LITERAL
                {
			$$ = WritePrimary(TYPE_I64, encoder.Serialize($1))
                }
        |       LPAREN expression RPAREN
                { $$ = $2 }
        |       array_literal_expression
                {
			$$ = $1
                }
                ;

after_period:   type_specifier
                {
			$$ = TypeNames[$1]
                }
        |       IDENTIFIER
        ;

postfix_expression:
                primary_expression
	|       postfix_expression LBRACK expression RBRACK
                {
			$$ = PostfixExpressionArray($1, $3)
                }
        |       type_specifier PERIOD after_period
                {
			$$ = PostfixExpressionNative(int($1), $3)
                }
	|       postfix_expression LPAREN RPAREN
                {
			$$ = PostfixExpressionEmptyFunCall($1)
                }
	|       postfix_expression LPAREN argument_expression_list RPAREN
                {
			$$ = PostfixExpressionFunCall($1, $3)
                }
	|       postfix_expression INC_OP
                {
			$$ = PostfixExpressionIncDec($1, true)
                }
        |       postfix_expression DEC_OP
                {
			$$ = PostfixExpressionIncDec($1, false)
                }
        |       postfix_expression PERIOD IDENTIFIER
                {
			PostfixExpressionField($1, $3)
                }
                ;

argument_expression_list:
                assignment_expression
	|       argument_expression_list COMMA assignment_expression
                {
			$$ = append($1, $3...)
                }
                ;

unary_expression:
                postfix_expression
	|       INC_OP unary_expression
                {
			// TODO
			$$ = $2
                }
	|       DEC_OP unary_expression
                {
			// TODO
			$$ = $2
                }
	|       unary_operator unary_expression
                {
			$$ = UnaryExpression($1, $2)
                }
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
                {
			$$ = ShorthandExpression($1, $3, OP_MUL)
                }
        |       multiplicative_expression DIV_OP unary_expression
                {
			$$ = ShorthandExpression($1, $3, OP_DIV)
                }
        |       multiplicative_expression MOD_OP unary_expression
                {
			$$ = ShorthandExpression($1, $3, OP_MOD)
                }
                ;

additive_expression:
                multiplicative_expression
        |       additive_expression ADD_OP multiplicative_expression
                {
			$$ = ShorthandExpression($1, $3, OP_ADD)
                }
	|       additive_expression SUB_OP multiplicative_expression
                {
			$$ = ShorthandExpression($1, $3, OP_SUB)
                }
                ;

shift_expression:
                additive_expression
        |       shift_expression LEFT_OP additive_expression
                {
			$$ = ShorthandExpression($1, $3, OP_BITSHL)
                }
        |       shift_expression RIGHT_OP additive_expression
                {
			$$ = ShorthandExpression($1, $3, OP_BITSHR)
                }
                ;

relational_expression:
                shift_expression
        |       relational_expression LT_OP shift_expression
                {
			$$ = ShorthandExpression($1, $3, OP_LT)
                }
        |       relational_expression GT_OP shift_expression
                {
			$$ = ShorthandExpression($1, $3, OP_GT)
                }
        |       relational_expression LTEQ_OP shift_expression
                {
			$$ = ShorthandExpression($1, $3, OP_LTEQ)
                }
        |       relational_expression GTEQ_OP shift_expression
                {
			$$ = ShorthandExpression($1, $3, OP_GTEQ)
                }
                ;

equality_expression:
                relational_expression
        |       equality_expression EQ_OP relational_expression
                {
			$$ = ShorthandExpression($1, $3, OP_EQUAL)
                }
        |       equality_expression NE_OP relational_expression
                {
			$$ = ShorthandExpression($1, $3, OP_UNEQUAL)
                }
                ;

and_expression: equality_expression
        |       and_expression REF_OP equality_expression
                {
			$$ = ShorthandExpression($1, $3, OP_BITAND)
                }
                ;

exclusive_or_expression:
                and_expression
        |       exclusive_or_expression BITXOR_OP and_expression
                {
			$$ = ShorthandExpression($1, $3, OP_BITXOR)
                }
                ;

inclusive_or_expression:
                exclusive_or_expression
        |       inclusive_or_expression BITOR_OP exclusive_or_expression
                {
			$$ = ShorthandExpression($1, $3, OP_BITOR)
                }
                ;

logical_and_expression:
                inclusive_or_expression
	|       logical_and_expression AND_OP inclusive_or_expression
                {
			$$ = ArithmeticOperation($1, $3, Natives[OP_BOOL_AND])
                }
                ;

logical_or_expression:
                logical_and_expression
	|       logical_or_expression OR_OP logical_and_expression
                {
			$$ = ArithmeticOperation($1, $3, Natives[OP_BOOL_OR])
                }
                ;

conditional_expression:
                logical_or_expression
	|       logical_or_expression '?' expression COLON conditional_expression
                ;

assignment_expression:
                conditional_expression
	|       unary_expression assignment_operator assignment_expression
                {
			if $3[0].IsArrayLiteral {
				$$ = ArrayLiteralAssignment($1, $3)
			} else if $3[0].IsStructLiteral {
				$$ = StructLiteralAssignment($1, $3)
			} else {
				$$ = Assignment($1, $3)
			}
                }
                ;

assignment_operator:
                ASSIGN
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
                {
			$3[len($3) - 1].Outputs = append($1[len($1) - 1].Outputs, $3[len($3) - 1].Outputs...)
			$$ = append($1, $3...)
                }
                ;

constant_expression:
                conditional_expression
                ;

declaration:
                VAR declarator declaration_specifiers SEMICOLON
                {
			$$ = DeclareLocal($2, $3, nil, false)
                }
        |       VAR declarator declaration_specifiers ASSIGN initializer SEMICOLON
                {
			$$ = DeclareLocal($2, $3, $5, true)
                }
                ;

initializer:    assignment_expression
                ;

designation:    designator_list ASSIGN
                {
			$$ = nil
                }
                ;

designator_list:
                designator
                {
			$$ = nil
                }
	|       designator_list designator
                {
			$$ = nil
                }
                ;

designator:
                LBRACK constant_expression RBRACK
                {
			$$ = nil
                }
	|       PERIOD IDENTIFIER
                {
			$$ = nil
                }
                ;






// statements
statement:      /* labeled_statement */
	/* |        */compound_statement
	|       expression_statement
	|       selection_statement
	|       iteration_statement
	/* |       jump_statement */
                ;

labeled_statement:
                IDENTIFIER COLON statement
                { $$ = nil }
	|       CASE constant_expression COLON statement
                { $$ = nil }
	|       DEFAULT COLON statement
                { $$ = nil }
                ;

compound_statement:
                LBRACE RBRACE SEMICOLON
                { $$ = nil }
	|       LBRACE block_item_list RBRACE SEMICOLON
                {
                    $$ = $2
                }
                ;

block_item_list:
                block_item
	|       block_item_list block_item
                {
			$$ = append($1, $2...)
                }
                ;

block_item:     declaration
        |       statement
                ;

expression_statement:
                SEMICOLON
                { $$ = nil }
	|       expression SEMICOLON
                {
			if $1[len($1) - 1].Operator == nil {
				$$ = nil
			} else {
				$$ = $1
			}
                }
                ;

selection_statement:
                IF expression LBRACE block_item_list RBRACE elseif_list else_statement SEMICOLON
                {
			$$ = SelectionStatement($2, $4, $6, $7, SEL_ELSEIFELSE)
                }
        |       IF expression LBRACE block_item_list RBRACE else_statement SEMICOLON
                {
			$$ = SelectionExpressions($2, $4, $6)
                }
        |       IF expression LBRACE block_item_list RBRACE elseif_list SEMICOLON
                {
			$$ = SelectionStatement($2, $4, $6, nil, SEL_ELSEIF)
                }
        |       IF expression compound_statement
                {
			$$ = SelectionExpressions($2, $3, nil)
                }
	|       SWITCH LPAREN expression RPAREN statement
                { $$ = nil }
                ;

elseif:         ELSE IF expression LBRACE block_item_list RBRACE
                {
			$$ = selectStatement{
				Condition: $3,
				Then: $5,
			}
                }
                ;

elseif_list:    elseif
                {
			$$ = []selectStatement{$1}
                }
        |       elseif_list elseif
                {
			$$ = append($1, $2)
                }
        ;

else_statement:
                ELSE LBRACE block_item_list RBRACE
                {
			$$ = $3
                }
        ;

iteration_statement:
                FOR expression compound_statement
                {
			$$ = IterationExpressions(nil, $2, nil, $3)
                }
        |       FOR expression_statement expression_statement compound_statement
                {			
			$$ = IterationExpressions($2, $3, nil, $4)
                }
        |       FOR expression_statement expression_statement expression compound_statement
                {
			$$ = IterationExpressions($2, $3, $4, $5)
                }
                ;

jump_statement: GOTO IDENTIFIER SEMICOLON
                { $$ = nil }
	|       CONTINUE SEMICOLON
                { $$ = nil }
	|       BREAK SEMICOLON
                { $$ = nil }
	|       RETURN SEMICOLON
                { $$ = nil }
	|       RETURN expression SEMICOLON
                { $$ = nil }
                ;

%%
