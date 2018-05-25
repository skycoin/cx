%{
	package cxgo0
	import (
		// "fmt"
		"bytes"
		. "github.com/skycoin/cx/cx"
	)

	var PRGRM *CXProgram
	var DataOffset int

	var lineNo int = 0
	var replMode bool = false
	var inREPL bool = false
	var inFn bool = false
	var fileName string

	func WritePrimary (typ int, byts []byte) []*CXExpression {
		if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
			arg := MakeArgument("")
			arg.AddType(TypeNames[typ])
                        arg.AddValue(&byts)
			arg.MemoryRead = MEM_DATA
			arg.MemoryWrite = MEM_DATA
			arg.Offset = DataOffset
			arg.Package = pkg
			arg.Program = PRGRM
			size := len(byts)
			arg.Size = size
			arg.TotalSize = size
			arg.PointeeSize = size
			DataOffset += size
			PRGRM.Data = append(PRGRM.Data, Data(byts)...)
			expr := MakeExpression(nil)
			expr.Package = pkg
			expr.Outputs = append(expr.Outputs, arg)
			return []*CXExpression{expr}
		} else {
			panic(err)
		}
	}

	func TotalLength (lengths []int) int {
		var total int = 1
		for _, i := range lengths {
			total *= i
		}
		return total
	}

	func FunctionDeclaration (fn *CXFunction, inputs []*CXArgument, outputs []*CXArgument, exprs []*CXExpression) {
		// adding inputs, outputs
		for _, inp := range inputs {
			fn.AddInput(inp)
		}
		for _, out := range outputs {
			fn.AddOutput(out)
		}
	}
	
	func Parse (code string) int {
		codeBuf := bytes.NewBufferString(code)
		return yyParse(NewLexer(codeBuf))
	}
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

        function *CXFunction
}

%token  <byt>           BYTENUM
%token  <i32>           INT BOOLEAN
%token  <i64>           LONG
%token  <f32>           FLOAT
%token  <f64>           DOUBLE
%token  <byt>           BYTE_LITERAL
%token  <i32>           INT_LITERAL BOOLEAN_LITERAL
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
                                                
%type   <function>      function_header

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
        ;

global_declaration:
                VAR declarator declaration_specifiers SEMICOLON
                {
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				expr := WritePrimary($3.Type, make([]byte, $3.TotalSize))
				exprOut := expr[0].Outputs[0]
				$3.Name = $2.Name
				$3.MemoryRead = MEM_DATA
				$3.MemoryWrite = MEM_DATA
				$3.Offset = exprOut.Offset
				$3.Package = exprOut.Package
				pkg.AddGlobal($3)
			} else {
				panic(err)
			}
                }
        |       VAR declarator declaration_specifiers ASSIGN initializer SEMICOLON
                {
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				expr := WritePrimary($3.Type, make([]byte, $3.Size))
				exprOut := expr[0].Outputs[0]
				$3.Name = $2.Name
				// fmt.Println("hou", $5[0].Outputs[0].Value)
				// $3.Value = $5[0].Outputs[0].Value
				$3.MemoryRead = MEM_DATA
				$3.MemoryWrite = MEM_DATA
				$3.Offset = exprOut.Offset
				$3.Size = exprOut.Size
				$3.TotalSize = exprOut.TotalSize
				$3.Package = exprOut.Package
				pkg.AddGlobal($3)
			} else {
				panic(err)
			}
                }
                ;

struct_declaration:
                TYPE IDENTIFIER STRUCT struct_fields
                {
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				strct := MakeStruct($2)
				pkg.AddStruct(strct)

				var size int
                                for _, fld := range $4 {
                                        strct.AddField(fld)
					size += fld.TotalSize
				}
				strct.Size = size
			} else {
				panic(err)
			}
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
			pkg := MakePackage($2)
			pkg.AddImport(pkg)
			PRGRM.AddPackage(pkg)
                }
                ;

import_declaration:
                IMPORT STRING_LITERAL SEMICOLON
                {
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				if _, err := pkg.GetImport($2); err != nil {
					if imp, err := PRGRM.GetPackage($2); err == nil {
						pkg.AddImport(imp)
					} else {
						panic(err)
					}
				}
			} else {
				panic(err)
			}
                }
                ;

function_header:
                FUNC IDENTIFIER
                {
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				fn := MakeFunction($2)
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
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				fn := MakeFunction($5)
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
			FunctionDeclaration($1, $2, nil, nil)
                }
        |       function_header function_parameters function_parameters compound_statement
                {
			FunctionDeclaration($1, $2, $3, nil)
                }
        ;

// parameter_type_list
parameter_type_list:
                //parameter_list COMMA ELLIPSIS
		parameter_list
                ;

parameter_list:
                parameter_declaration
                {
			// if $1.IsArray {
			// 	$1.TotalSize = $1.Size * TotalLength($1.Lengths)
			// } else {
			// 	$1.TotalSize = $1.Size
			// }
			$$ = []*CXArgument{$1}
                }
	|       parameter_list COMMA parameter_declaration
                {
			// if $3.IsArray {
			// 	$3.TotalSize = $3.Size * TotalLength($3.Lengths)
			// } else {
			// 	$3.TotalSize = $3.Size
			// }
			// lastPar := $1[len($1) - 1]
			// $3.Offset = lastPar.Offset + lastPar.TotalSize
			$$ = append($1, $3)
                }
                ;

parameter_declaration:
                declarator declaration_specifiers
                {
			$2.Name = $1.Name
			$2.Package = $1.Package
			// $2.MemoryType = MEM_STACK
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
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				arg := MakeArgument("")
				arg.AddType(TypeNames[TYPE_UNDEFINED])
				arg.Name = $1
				arg.Package = pkg
				$$ = arg
			} else {
				panic(err)
			}
                }
	|       LPAREN declarator RPAREN
                { $$ = $2 }
                ;


declaration_specifiers:
                MUL_OP declaration_specifiers
                {
			$2.DeclarationSpecifiers = append($2.DeclarationSpecifiers, DECL_POINTER)
			if !$2.IsPointer {
				$2.IsPointer = true
				// $2.MemoryRead = MEM_HEAP
				$2.PointeeSize = $2.Size
				$2.Size = TYPE_POINTER_SIZE
				$2.TotalSize = TYPE_POINTER_SIZE
				$2.IndirectionLevels++
			} else {
				pointer := $2

				for c := $2.IndirectionLevels - 1; c > 0 ; c-- {
					pointer = pointer.Pointee
					pointer.IndirectionLevels = c
					pointer.IsPointer = true
				}

				pointee := MakeArgument("")
				pointee.AddType(TypeNames[pointer.Type])
				// pointee.Size = pointer.Size
				// pointee.TotalSize = pointer.TotalSize
				pointee.IsPointer = true

				$2.IndirectionLevels++

				// pointer.Type = TYPE_POINTER
				pointer.Size = TYPE_POINTER_SIZE
				pointer.TotalSize = TYPE_POINTER_SIZE
				pointer.Pointee = pointee
			}
			
			$$ = $2
                }
        |       LBRACK RBRACK declaration_specifiers
                {
			$3.DeclarationSpecifiers = append($3.DeclarationSpecifiers, DECL_SLICE)
			arg := $3
                        arg.IsArray = true
			arg.Lengths = append([]int{SLICE_SIZE}, arg.Lengths...)
			arg.TotalSize = arg.Size * TotalLength(arg.Lengths)
			$$ = arg
                }
        |       LBRACK INT_LITERAL RBRACK declaration_specifiers
                {
			$4.DeclarationSpecifiers = append($4.DeclarationSpecifiers, DECL_ARRAY)
			arg := $4
                        arg.IsArray = true
			arg.Lengths = append([]int{int($2)}, arg.Lengths...)
			arg.TotalSize = arg.Size * TotalLength(arg.Lengths)
			// arg.Size = GetArgSize($4.Type)
			$$ = arg
                }
        |       type_specifier
                {
			arg := MakeArgument("")
			arg.AddType(TypeNames[$1])
			arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_BASIC)
			arg.Type = $1
			arg.Size = GetArgSize($1)
			arg.TotalSize = arg.Size
			$$ = arg
                }
        |       IDENTIFIER
                {
			// custom type in the current package
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				if strct, err := PRGRM.GetStruct($1, pkg.Name); err == nil {
					arg := MakeArgument("")
					arg.AddType(TypeNames[TYPE_CUSTOM])
					arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)
					arg.CustomType = strct
					arg.Size = strct.Size
					arg.TotalSize = strct.Size

					$$ = arg
				} else {
					panic("type '" + $1 + "' does not exist")
				}
			} else {
				panic(err)
			}
                }
        |       IDENTIFIER PERIOD IDENTIFIER
                {
			// custom type in an imported package
			if pkg, err := PRGRM.GetCurrentPackage(); err == nil {
				if imp, err := pkg.GetImport($1); err == nil {
					if strct, err := PRGRM.GetStruct($3, imp.Name); err == nil {
						arg := MakeArgument("")
						arg.AddType(TypeNames[TYPE_CUSTOM])
						arg.CustomType = strct
						arg.Size = strct.Size
						arg.TotalSize = strct.Size
						arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)

						$$ = arg
					} else {
						panic("type '" + $1 + "' does not exist")
					}
				} else {
					panic(err)
				}
			} else {
				panic(err)
			}
			
			// if pkg, err := PRGRM.GetPackage($1); err == nil {
			// 	if strct, err := PRGRM.GetStruct($3, pkg.Name); err == nil {
			// 		arg := MakeArgument(TYPE_CUSTOM)
			// 		arg.CustomType = strct
			// 		arg.Size = strct.Size
			// 		arg.TotalSize = strct.Size

			// 		$$ = arg
			// 	} else {
			// 		panic("type '" + $1 + "' does not exist")
			// 	}
			// } else {
			// 	panic(err)
			// }
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
        |       IDENTIFIER COLON constant_expression
        |       struct_literal_fields COMMA IDENTIFIER COLON constant_expression
                ;


// expressions
array_literal_expression:
                LBRACK INT_LITERAL RBRACK IDENTIFIER LBRACE argument_expression_list RBRACE
        |       LBRACK INT_LITERAL RBRACK IDENTIFIER LBRACE RBRACE
        |       LBRACK INT_LITERAL RBRACK type_specifier LBRACE argument_expression_list RBRACE
        |       LBRACK INT_LITERAL RBRACK type_specifier LBRACE RBRACE
        |       LBRACK INT_LITERAL RBRACK array_literal_expression
        ;

slice_literal_expression:
                LBRACK RBRACK IDENTIFIER LBRACE argument_expression_list RBRACE
        |       LBRACK RBRACK IDENTIFIER LBRACE RBRACE
        |       LBRACK RBRACK type_specifier LBRACE argument_expression_list RBRACE
        |       LBRACK RBRACK type_specifier LBRACE RBRACE
        |       LBRACK RBRACK slice_literal_expression
                ;

primary_expression:
                IDENTIFIER
        |       IDENTIFIER LBRACE struct_literal_fields RBRACE
        |       STRING_LITERAL
        |       BOOLEAN_LITERAL
        |       BYTE_LITERAL
        |       INT_LITERAL
        |       FLOAT_LITERAL
        |       DOUBLE_LITERAL
        |       LONG_LITERAL
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
                ;

relational_expression:
                shift_expression
	|       relational_expression LT_OP shift_expression
	|       relational_expression GT_OP shift_expression
	|       relational_expression LTEQ_OP shift_expression
	|       relational_expression GTEQ_OP shift_expression
                ;

equality_expression:
                relational_expression
	|       equality_expression EQ_OP relational_expression
	|       equality_expression NE_OP relational_expression
                ;

and_expression: equality_expression
	|       and_expression REF_OP equality_expression
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
        |       postfix_expression PERIOD IDENTIFIER LBRACE struct_literal_fields RBRACE
        ;

assignment_expression:
                /* conditional_expression */
                struct_literal_expression
	|       unary_expression assignment_operator assignment_expression
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

initializer_list:
                designation initializer
	|       initializer
	|       initializer_list COMMA designation initializer
	|       initializer_list COMMA initializer
                ;

designation:    designator_list ASSIGN
                ;

designator_list:
                designator
	|       designator_list designator
                ;

designator:
                LBRACK constant_expression RBRACK
	|       PERIOD IDENTIFIER
                ;






// statements
statement:      compound_statement
	|       expression_statement
	|       selection_statement
	|       iteration_statement
                ;

/* labeled_statement: */
/*                 IDENTIFIER COLON statement */
/* 	|       CASE constant_expression COLON statement */
/* 	|       DEFAULT COLON statement */
/*                 ; */

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
        |       IF conditional_expression LBRACE block_item_list RBRACE elseif_list SEMICOLON
        |       IF conditional_expression compound_statement
	|       SWITCH LPAREN conditional_expression RPAREN statement
                ;

elseif:         ELSE IF expression LBRACE block_item_list RBRACE
        ;

elseif_list:
                elseif
        |       elseif_list elseif
        ;

else_statement:
                ELSE LBRACE block_item_list RBRACE
        ;

iteration_statement:
                FOR expression compound_statement
        |       FOR expression_statement expression_statement compound_statement
        |       FOR expression_statement expression_statement expression compound_statement
        |       FOR declaration expression_statement compound_statement
        |       FOR declaration expression_statement expression compound_statement
                ;

/* jump_statement: GOTO IDENTIFIER SEMICOLON */
/* 	|       CONTINUE SEMICOLON */
/* 	|       BREAK SEMICOLON */
/* 	|       RETURN SEMICOLON */
/* 	|       RETURN expression SEMICOLON */
/*                 ; */

%%
