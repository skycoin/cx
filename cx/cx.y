%{
	package main
	import (
		// "strings"
		// "fmt"
		// "os"
		// "time"

		// "github.com/skycoin/cx/cx/cx0"
		"github.com/skycoin/skycoin/src/cipher/encoder"
		. "github.com/skycoin/cx/src/base"
	)

	var prgrm = MakeProgram(1024, 1024, 1024)
	var data Data
	var dataOffset int

	var lineNo int = 0
	var webMode bool = false
	var baseOutput bool = false
	var replMode bool = false
	var helpMode bool = false
	var compileMode bool = false
	var replTargetFn string = ""
	var replTargetStrct string = ""
	var replTargetMod string = ""
	// var dStack bool = false
	var inREPL bool = false
	// var inFn bool = false
	// var tag string = ""
	// var asmNL = "\n"
	var fileName string

	// Primary expressions (literals) are saved in the MEM_DATA segment at compile-time
	// This function writes those bytes to prgrm.Data
	func WritePrimary (typ int, byts []byte) []*CXExpression {
		if pkg, err := prgrm.GetCurrentPackage(); err == nil {
			arg := MakeArgument(typ)
			arg.MemoryType = MEM_DATA
			arg.Offset = dataOffset
			arg.Package = pkg
			arg.Program = prgrm
			size := len(byts)
			arg.Size = size
			dataOffset += size
			prgrm.Data = append(prgrm.Data, Data(byts)...)
			expr := MakeExpression(nil)
			expr.Outputs = append(expr.Outputs, arg)
			return []*CXExpression{expr}
		} else {
			panic(err)
		}
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

	/* parameter *CXParameter */
	/* parameters []*CXParameter */

	argument *CXArgument
	arguments []*CXArgument

        /* definition *CXDefinition */
	/* definitions []*CXDefinition */

	expression *CXExpression
	expressions []*CXExpression

        function *CXFunction

	/* field *CXField */
	/* fields []*CXField */

	/* name string */
	/* names []string */
}

%token  <byt>           BYTE_LITERAL
%token  <i32>           INT_LITERAL BOOLEAN_LITERAL
%token  <i64>           LONG_LITERAL
%token  <f32>           FLOAT_LITERAL
%token  <f64>           DOUBLE_LITERAL
%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK IDENTIFIER
                        VAR COMMA PERIOD COMMENT STRING_LITERAL PACKAGE IF ELSE FOR TYPSTRUCT STRUCT
                        SEMICOLON EXCL
                        ASSIGN CASSIGN IMPORT RETURN GOTO GTHAN LTHAN EQUAL COLON NEW
                        EQUALWORD GTHANWORD LTHANWORD
                        GTHANEQ LTHANEQ UNEQUAL AND OR
                        ADD_OP SUB_OP MUL_OP DIV_OP MOD_OP AFFVAR
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

%type   <i>             type_specifier
%type   <argument>      declaration_specifiers
%type   <argument>      declarator
%type   <argument>      direct_declarator
%type   <argument>      parameter_declaration
%type   <arguments>     parameter_type_list
%type   <arguments>     parameter_list
                                                
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

%type   <expressions>   expression
%type   <expressions>   block_item_list
%type   <expressions>   compound_statement
%type   <expressions>   labeled_statement
%type   <expressions>   expression_statement
%type   <expressions>   selection_statement
%type   <expressions>   iteration_statement
%type   <expressions>   jump_statement
%type   <expressions>   statement

%type   <function>      function_header

/* %start                  translation_unit */
%%

translation_unit:
                external_declaration
        |       translation_unit external_declaration
        ;

external_declaration:
                package_declaration
        // |       global_declaration
        |       function_declaration
        /* |       method_declaration */
        /* |       struct_declaration */
        ;

package_declaration:
                PACKAGE IDENTIFIER SEMICOLON
                {
			pkg := MakePackage($2)
			prgrm.AddPackage(pkg)
                }
                ;

function_header:
                FUNC IDENTIFIER
                {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				fn := MakeFunction($2)
				pkg.AddFunction(fn)

                                $$ = fn
			} else {
				panic(err)
			}
                }
        ;

function_declaration:
                function_header LPAREN parameter_type_list RPAREN compound_statement
        |       function_header LPAREN parameter_type_list RPAREN LPAREN parameter_type_list RPAREN compound_statement
                {
			// fixing output parameters' offsets (they need to be + last input offset)
			offset := $3[len($3) - 1].Offset + $3[len($3) - 1].Size
			for _, out := range $6 {
				out.Offset += offset
				offset += out.Size
			}

			// adding all the inputs, outputs and expressions
			for _, inp := range $3 {
				$1.AddInput(inp)
			}
			for _, out := range $6 {
				$1.AddOutput(out)
			}
			for _, expr := range $8 {
				$1.AddExpression(expr)
			}

			$1.Length = len($1.Expressions)

			var symbols map[string]*CXArgument = make(map[string]*CXArgument, 0)
			for _, inp := range $1.Inputs {
				if inp.Name != "" {
					symbols[inp.Name] = inp
				}
			}
			for _, out := range $1.Outputs {
				if out.Name != "" {
					symbols[out.Name] = out
				}
			}

			for _, expr := range $1.Expressions {
				for _, inp := range expr.Inputs {
					if inp.Name != "" {
						if arg, found := symbols[inp.Name]; !found {
							// it should exist. error
							panic("identifier '" + inp.Name + "' does not exist")
						} else {
							// inp.Offset = off
							// arg.Offset = offset
							inp.Offset = arg.Offset
							inp.Size = arg.Size
						}
					}
				}
				for _, out := range expr.Outputs {
					if out.Name != "" {
						if arg, found := symbols[out.Name]; !found {
							out.Offset = offset
							symbols[out.Name] = out
							offset += out.Size
						} else {
							out.Offset = arg.Offset
							out.Size = arg.Size
						}
					}
				}
			}
			$1.Size = offset
                }
        ;

/* method_declaration: */
/*                 FUNC */
/*         ; */


// parameter_type_list
parameter_type_list:
                //parameter_list COMMA ELLIPSIS
		parameter_list
                ;

parameter_list:
                parameter_declaration
                {
			$$ = []*CXArgument{$1}
                }
	|       parameter_list COMMA parameter_declaration
                {
			lastPar := $1[len($1) - 1]
			$3.Offset = lastPar.Offset + lastPar.Size
			$$ = append($1, $3)
                }
                ;

parameter_declaration:
                declarator declaration_specifiers
                {
			$2.Name = $1.Name
			$2.IsArray = $1.IsArray
			// input and output parameters are always in the stack
			$2.MemoryType = MEM_STACK
			$$ = $2
                }
        //                      |declaration_specifiers abstract_declarator
	/* |    declaration_specifiers */
                ;

identifier_list:
                IDENTIFIER
	|       identifier_list COMMA IDENTIFIER
                ;

// declarator
declarator:     /* pointer direct_declarator */
        /*         { */
        /*             $2.IsPointer = true */
        /*             $$ = $2 */
        /*         } */
	/* |        */direct_declarator
                ;

direct_declarator:
                IDENTIFIER
                {
			arg := MakeArgument(TYPE_UNDEFINED)
			arg.Name = $1
			$$ = arg
                }
	|       LPAREN   declarator RPAREN
                { $$ = $2 }
	|       direct_declarator '[' ']'
                {
			$1.IsArray = true
			$$ = $1
                }
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
pointer:        MUL_OP   type_qualifier_list pointer // check
        |       MUL_OP   type_qualifier_list // check
        |       MUL_OP   pointer
        |       MUL_OP
                ;

type_qualifier_list:
                type_qualifier
	|       type_qualifier_list type_qualifier
                ;









// declaration_specifiers
declaration_specifiers:
                pointer type_specifier
                {
			arg := MakeArgument($2)
                        arg.IsPointer = true
			arg.Size = GetArgSize($2)
			$$ = arg
                }
        |       type_specifier
                {
			arg := MakeArgument($1)
			arg.Size = GetArgSize($1)
			$$ = arg
                }
		/* type_specifier declaration_specifiers */
	/* |       type_specifier */
	/* |       type_qualifier declaration_specifiers */
	/* |       type_qualifier */
                ;

type_qualifier: CONST
        |       VAR
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
        /* |       pointer type_specifier */
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









// expressions
primary_expression:
                IDENTIFIER
                {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				arg := MakeArgument(TYPE_IDENTIFIER)
				arg.Name = $1
				arg.Package = pkg
				$$ = []*CXExpression{&CXExpression{Outputs: []*CXArgument{arg}}}
			} else {
				panic(err)
			}
                }
        |       type_specifier PERIOD IDENTIFIER
                {
			// these will always be native functions
			if opCode, ok := OpCodes[TypeNames[$1] + "." + $3]; ok {
				$$ = []*CXExpression{MakeExpression(Natives[opCode])}
			} else {
				panic(ok)
			}
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
                ;

postfix_expression:
                primary_expression
	// |       postfix_expression '[' expression ']'
	|       postfix_expression LPAREN RPAREN
                {
			// expr := $1[0]
			$$ = $1
			// fmt.Println("yoyo", expr)

			// arg := $1[0].Outputs[0]
			// opName := arg.Name

			// if op, err := prgrm.GetFunction(opName, arg.Package.Name); err == nil {
			// 	$$ = []*CXExpression{MakeExpression(op)}
			// } else {
			// 	panic(err)
			// }
                }
	|       postfix_expression LPAREN argument_expression_list RPAREN
                {
			/*
			  i32.add(5, 5), foo(10, 20)
			*/
			expr := $1[0]
			if expr.Operator == nil {
				if op, err := prgrm.GetFunction($1[0].Outputs[0].Name, $1[0].Outputs[0].Package.Name); err == nil {
					expr.Operator = op
				} else {
					panic(err)
				}
			}

			var nestedExprs []*CXExpression
			for _, inpExpr := range $3 {
				if inpExpr.Operator == nil {
					// then it's a literal
					expr.AddInput(inpExpr.Outputs[0])
				} else {
					// then it's a function call
					if len(inpExpr.Outputs) < 1 {
						out := MakeParameter(MakeGenSym(LOCAL_PREFIX), inpExpr.Operator.Outputs[0].Type)
						inpExpr.AddOutput(out)
						expr.AddInput(out)
					}
					nestedExprs = append(nestedExprs, inpExpr)
				}
			}
			
			nestedExprs = append(nestedExprs, $1...)
			$$ = nestedExprs
                }
	/* |       postfix_expression PERIOD IDENTIFIER */
	/* |       postfix_expression PERIOD IDENTIFIER // check */
        /*         { */

        /*         } */
	|       postfix_expression INC_OP
                {
			$$ = $1
                }
        |       postfix_expression DEC_OP
                {
			$$ = $1
                }
        |       postfix_expression PERIOD IDENTIFIER
                {
			fullName := $1[len($1) - 1].Outputs[0].Name + "." + $3
			if opCode, ok := OpCodes[fullName]; ok {
				$$ = []*CXExpression{MakeExpression(Natives[opCode])}
			}
			
			// left := $1[0].Outputs[0]
			// //right := $3[0].Outputs[0]
			
			// if left.IsRest {
			// 	// then it can't be a module name
			// 	// and we propagate the property to the right expression
			// 	// right.IsRest = true
			// } else {
			// 	// then left is a first (e.g first.rest) and right is a rest
			// 	// right.IsRest = true
			// 	// let's check if left is a package
			// 	if _, err := prgrm.GetPackage(left.Name); err != nil {
			// 		// the external property will be propagated to the following arguments
			// 		// this way we avoid considering these arguments as module names
			// 		//right.Package = pkg
			// 		$$ = $1
			// 	} else {
			// 		// left is not a package, then it's a struct or a function call
			// 		// if right.Operator != nil {
			// 		// 	// then left is a function call that returns a struct instance
			// 		// 	// we just return an expression with the GetField operator
						
			// 		// } else {
			// 		// 	// then left is a struct instance

			// 		// 	// GetField(left, right)
			// 		// }
			// 	}
			// }
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
                {$$ = $2}
	|       DEC_OP unary_expression
                {$$ = $2}
	|       unary_operator unary_expression // check
                {$$ = $2}
                ;

unary_operator:
                '&'
	|       MUL_OP
	|       ADD_OP
	|       SUB_OP
	|       '~' // check
	|       '!'
                ;

multiplicative_expression:
                unary_expression
	|       multiplicative_expression MUL_OP unary_expression
	|       multiplicative_expression '/' unary_expression
	|       multiplicative_expression '%' unary_expression
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
                { $$ = nil }
                ;

logical_or_expression:
                logical_and_expression
	|       logical_or_expression OR_OP logical_and_expression
                ;

conditional_expression:
                logical_or_expression
	|       logical_or_expression '?' expression COLON conditional_expression
                ;

assignment_expression:
                conditional_expression
	|       unary_expression assignment_operator assignment_expression
                {
			idx := len($3) - 1

			if $3[idx].Operator == nil {
				$3[idx].Operator = Natives[OP_IDENTITY]
				$1[0].Outputs[0].Size = $3[idx].Outputs[0].Size
				$3[idx].Inputs = $3[idx].Outputs
				$3[idx].Outputs = $1[0].Outputs
				
				$$ = $3
			} else {
				if $3[idx].Operator.IsNative {
					for i, _ := range $3[idx].Operator.Outputs {
						$1[0].Outputs[i].Size = Natives[$3[idx].Operator.OpCode].Outputs[i].Size
					}
				} else {
					for i, out := range $3[idx].Operator.Outputs {
						$1[0].Outputs[i].Size = out.Size
					}
				}

				$3[idx].Outputs = $1[0].Outputs
				$$ = $3
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
                { $$ = append($1, $3...) }
                ;

constant_expression:
                conditional_expression
                ;






/* declaration: */
/*                 declaration_specifiers init_declarator_list SEMICOLON */
/*                 ; */

/* declaration_specifiers: */
/* 		type_specifier declaration_specifiers */
/* 	|       type_specifier */
/* 	|       type_qualifier declaration_specifiers */
/* 	|       type_qualifier */
/*                 ; */



/* init_declarator_list: */
/*                 init_declarator */
/* 	|       init_declarator_list COMMA init_declarator */
/*                 ; */

/* init_declarator: */
/*                 declarator '=' initializer */
/* 	|       declarator */
/*                 ; */










/* initializer: */
/*                 LBRACE initializer_list RBRACE */
/* 	| LBRACE   initializer_list COMMA RBRACE */
/* 	|       assignment_expression */
/*                 ; */

/* initializer_list: */
/*                 designation initializer */
/* 	|       initializer */
/* 	|       initializer_list COMMA designation initializer */
/* 	|       initializer_list COMMA initializer */
/*                 ; */

/* designation:    designator_list '=' */
/*                 ; */

/* designator_list: */
/*                 designator */
/* 	|       designator_list designator */
/*                 ; */

/* designator: */
/*                 '[' constant_expression ']' */
/* 	| PERIOD   IDENTIFIER */
/*                 ; */






// statements
statement:      /* labeled_statement */
	/* |        */compound_statement
	|       expression_statement
	/* |       selection_statement */
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
                LBRACE RBRACE
                { $$ = nil }
	|       LBRACE block_item_list RBRACE
                {
                    $$ = $2
                }
                ;

block_item_list:
                statement
	|       block_item_list statement
                {
			$$ = append($1, $2...)
                }
                ;

// block_item:     /* declaration */
// 	/* |        */statement
//                 ;

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
                /* IF LPAREN expression RPAREN statement */
	/* |    IF LPAREN expression RPAREN statement ELSE statement */
                IF LPAREN expression RPAREN compound_statement elseif_list else_statement
                { $$ = nil }
	|       SWITCH LPAREN expression RPAREN statement
                { $$ = nil }
                ;

elseif:         ELSE IF expression compound_statement
        ;

elseif_list:
        |       elseif_list elseif
        ;

else_statement:
        |       ELSE compound_statement
        ;



iteration_statement:
                // FOR expression_statement expression_statement statement
        //         { $$ = nil }
	// |       
		FOR LPAREN expression_statement expression_statement expression RPAREN statement
                {
			jmpFn := Natives[OP_JMP]
			
			upExpr := MakeExpression(jmpFn)
			trueArg := WritePrimary(TYPE_BOOL, encoder.Serialize(true))
			upLines := WritePrimary(TYPE_I32, encoder.Serialize(int32((len($7) + len($5) + len($4) + 2) * -1)))
			downLines := WritePrimary(TYPE_I32, encoder.SerializeAtomic(int32(0)))
			upExpr.AddInput(trueArg[0].Outputs[0])
			upExpr.AddInput(upLines[0].Outputs[0])
			upExpr.AddInput(downLines[0].Outputs[0])
			
			downExpr := MakeExpression(jmpFn)
			
			if len($4[len($4) - 1].Outputs) < 1 {
				predicate := MakeParameter(MakeGenSym(LOCAL_PREFIX), $4[len($4) - 1].Operator.Outputs[0].Type)
				$4[len($4) - 1].AddOutput(predicate)
				downExpr.AddInput(predicate)
			} else {
				predicate := $4[len($4) - 1].Outputs[0]
				downExpr.AddInput(predicate)
			}
			thenLines := WritePrimary(TYPE_I32, encoder.SerializeAtomic(int32(0)))
			elseLines := WritePrimary(TYPE_I32, encoder.SerializeAtomic(int32(len($5) + len($7) + 1)))
			
			downExpr.AddInput(thenLines[0].Outputs[0])
			downExpr.AddInput(elseLines[0].Outputs[0])
			
			exprs := $3
			exprs = append(exprs, $4...)
			exprs = append(exprs, downExpr)
			exprs = append(exprs, $7...)
			exprs = append(exprs, $5...)
			exprs = append(exprs, upExpr)
			
			$$ = exprs
                }
	/* |       FOR LPAREN declaration expression_statement RPAREN statement */
	/* |       FOR LPAREN declaration expression_statement expression RPAREN statement */
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
















// lines:
//                 /* empty */
//         |       lines line
//         |       lines SEMICOLON
//         ;

// line:
//                 definitionDeclaration
//         |       structDeclaration
//         |       packageDeclaration
//         |       importDeclaration
//         |       functionDeclaration
//         |       selector
//         |       stepping
//         |       debugging
//         |       affordance
//         |       remover
//         ;

// importDeclaration:
//                 IMPORT STRING
//         ;

// affordance:
//                 TAG
//                 /* Function Affordances */
//         |       AFF FUNC IDENT
//         |       AFF FUNC IDENT LBRACE INT RBRACE
//         |       AFF FUNC IDENT LBRACE STRING RBRACE
//         |       AFF FUNC IDENT LBRACE STRING INT RBRACE
//                 /* Module Affordances */
//         |       AFF PACKAGE IDENT
//         |       AFF PACKAGE IDENT LBRACE INT RBRACE
//         |       AFF PACKAGE IDENT LBRACE STRING RBRACE
//         |       AFF PACKAGE IDENT LBRACE STRING INT RBRACE
//                 /* Struct Affordances */
//         |       AFF STRUCT IDENT
//         |       AFF STRUCT IDENT LBRACE INT RBRACE
//         |       AFF STRUCT IDENT LBRACE STRING RBRACE
//         |       AFF STRUCT IDENT LBRACE STRING INT RBRACE
//                 /* Struct Affordances */
//         |       AFF EXPR IDENT
//         |       AFF EXPR IDENT LBRACE INT RBRACE
//         |       AFF EXPR IDENT LBRACE STRING RBRACE
//         |       AFF EXPR IDENT LBRACE STRING INT RBRACE
//         ;

// stepping:       TSTEP INT INT

//         |       DSTACK
//         |       DPROGRAM
//         ;

// remover:        REM FUNC IDENT
//         |       REM PACKAGE IDENT
//         |       REM DEF IDENT
//         |       REM STRUCT IDENT
//         |       REM IMPORT STRING
//         |       REM EXPR IDENT FUNC IDENT
//         |       REM FIELD IDENT STRUCT IDENT
//         |       REM INPUT IDENT FUNC IDENT
//         |       REM OUTPUT IDENT FUNC IDENT
//         ;

// selectorLines:
//                 /* empty */
//         |       LBRACE lines RBRACE
//         ;

// selectorExpressionsAndStatements:
//                 /* empty */
//         |       LBRACE expressionsAndStatements RBRACE
//         ;

// selectorFields:
//                 /* empty */
//         |       LBRACE fields RBRACE
//         ;

// selector:       SPACKAGE IDENT
//                 selectorLines
//         |       SFUNC IDENT
//                 selectorExpressionsAndStatements
//         |       SSTRUCT IDENT
//                 selectorFields
//         ;






// assignOperator:
//                 ASSIGN
//         |       CASSIGN
//         |       PLUSEQ
//         |       MINUSEQ
//         |       MULTEQ
//         |       DIVEQ
//         |       REMAINDEREQ
//         |       EXPEQ
//         |       LEFTSHIFTEQ
//         |       RIGHTSHIFTEQ
//         |       BITANDEQ
//         |       BITXOREQ
//         |       BITOREQ
//         ;

// packageDeclaration:
//                 PACKAGE IDENT
//                 ;

// definitionAssignment:
//                 /* empty */
//         |       assignOperator argument
//         |       assignOperator ADDR argument
//         |       assignOperator VALUE argument
//                 ;

// definitionDeclaration:
//                 VAR IDENT BASICTYPE definitionAssignment
//         |       VAR IDENT IDENT
//         ;

// fields:
//                 parameter
//         |       SEMICOLON
//         |       debugging
//         |       fields parameter
//         |       fields SEMICOLON
//         |       fields debugging
//                 ;

// structFields:
//                 LBRACE fields RBRACE
//         |       LBRACE RBRACE
//         ;

// structDeclaration:
//                 TYPSTRUCT IDENT
//                 STRUCT structFields
//         ;

// functionParameters:
//                 LPAREN parameters RPAREN
//         |       LPAREN RPAREN
//         ;

// functionDeclaration:
//                 /* Methods */
//                 FUNC functionParameters IDENT functionParameters functionParameters
//                 functionStatements
//         |       FUNC functionParameters IDENT functionParameters
//                 functionStatements
//                 /* Functions */
//         |       FUNC IDENT functionParameters
//                 functionStatements
//         |       FUNC IDENT functionParameters functionParameters
//                 functionStatements
//         ;

// parameter:
//                 IDENT BASICTYPE
//         |       IDENT IDENT
//         |       IDENT MULT IDENT
//         ;

// parameters:
//                 parameter
//         |       parameters COMMA parameter
//         ;

// functionStatements:
//                 LBRACE expressionsAndStatements RBRACE
//         |       LBRACE RBRACE
//         ;

// expressionsAndStatements:
//                 nonAssignExpression
//         |       assignExpression
//         |       statement
//         |       selector
//         |       stepping
//         |       debugging
//         |       affordance
//         |       remover
//         |       expressionsAndStatements nonAssignExpression
//         |       expressionsAndStatements assignExpression
//         |       expressionsAndStatements statement
//         |       expressionsAndStatements selector
//         |       expressionsAndStatements stepping
//         |       expressionsAndStatements debugging
//         |       expressionsAndStatements affordance
//         |       expressionsAndStatements remover
//         ;


// assignExpression:
//                 VAR IDENT BASICTYPE definitionAssignment
//         |       VAR IDENT LBRACK RBRACK IDENT
//         |       argumentsList assignOperator argumentsList
//         ;

// nonAssignExpression:
//                 IDENT arguments
//         |       argument PLUSPLUS
//         |       argument MINUSMINUS
//         ;

// beginFor:       FOR
//                 ;

// conditionControl:
//                 nonAssignExpression
//         |       argument
//         ;

// returnArg:
//                 SEMICOLON
//         |       argumentsList
//                 ;

// statement:      RETURN returnArg
//         |       GOTO IDENT
//         |       IF conditionControl
//                 LBRACE
//                 expressionsAndStatements RBRACE elseStatement
//         |       beginFor
//                 nonAssignExpression
//                 LBRACE expressionsAndStatements RBRACE
//         |       beginFor
//                 argument
//                 LBRACE expressionsAndStatements RBRACE
//         |       beginFor // $<i>1
//                 forLoopAssignExpression // $2
//                 SEMICOLON conditionControl
//                 SEMICOLON forLoopAssignExpression //$<bool>9
//                 LBRACE expressionsAndStatements RBRACE
//         |       VAR IDENT IDENT                
//         |       SEMICOLON
//         ;

// forLoopAssignExpression:                
//         |       assignExpression
//         |       nonAssignExpression
//         ;

// elseStatement:
//                 /* empty */
//         |       ELSE
//                 LBRACE
//                 expressionsAndStatements RBRACE
//                 ;

// expressions:
//                 nonAssignExpression
//         |       assignExpression
//         |       expressions nonAssignExpression
//         |       expressions assignExpression
//         ;











// inferPred:      inferObj
//         |       inferCond
//         |       inferPred COMMA inferObj
//         |       inferPred COMMA inferCond
//         ;

// inferCond:      IDENT LPAREN inferPred RPAREN
//         |       BOOLEAN
//         ;

// relationalOp:   EQUAL
//         |       GTHAN
//         |       LTHAN
//         |       UNEQUAL
//                 ;

// inferActionArg:
//                 inferObj
//         |       IDENT
//         |       AFFVAR relationalOp argument
//         |       AFFVAR relationalOp nonAssignExpression
//         |       AFFVAR relationalOp AFFVAR
//         ;

// inferAction:
// 		IDENT LPAREN inferActionArg RPAREN
//         ;

// inferActions:
//                 inferAction
//         |       inferActions inferAction
//                 ;

// inferRule:      IF inferCond LBRACE inferActions RBRACE
//         |       IF inferObj LBRACE inferActions RBRACE
//         ;

// inferRules:     inferRule
//         |       inferRules inferRule
//         ;

// inferWeight:    FLOAT
//         |       INT
//         |       IDENT
//         ;

// inferObj:
//         |       IDENT VALUE inferWeight
//         |       IDENT VALUE nonAssignExpression
// ;

// inferObjs:      inferObj
//         |       inferObjs COMMA inferObj
//         ;

// inferTarget:    IDENT LPAREN IDENT RPAREN
//         ;

// inferTargets:   inferTarget
//         |       inferTargets inferTarget
//         ;

// inferClauses:   inferObjs
//         |       inferRules
//         |       inferTargets
//         ;




// structLitDef:
//                 TAG argument
//         |       TAG nonAssignExpression
                    
// ;

// structLitDefs:  structLitDef
//         |       structLitDefs COMMA structLitDef
//         ;

// structLiteral:
//         |       
// 		LBRACE structLitDefs RBRACE
//         ;

// /*
//   Fix this, there has to be a way to compress these rules
// */
// argument:       argument PLUS argument
//         |       nonAssignExpression PLUS nonAssignExpression
//         |       argument PLUS nonAssignExpression
//         |       nonAssignExpression PLUS argument


//         |       argument MINUS argument
//         |       nonAssignExpression MINUS nonAssignExpression
//         |       argument MINUS nonAssignExpression
//         |       nonAssignExpression MINUS argument

//         |       argument MULT argument
//         |       nonAssignExpression MULT nonAssignExpression
//         |       argument MULT nonAssignExpression
//         |       nonAssignExpression MULT argument

//         |       argument DIV argument
//         |       nonAssignExpression DIV nonAssignExpression
//         |       argument DIV nonAssignExpression
//         |       nonAssignExpression DIV argument

//         |       argument REMAINDER argument
//         |       nonAssignExpression REMAINDER nonAssignExpression
//         |       argument REMAINDER nonAssignExpression
//         |       nonAssignExpression REMAINDER argument

//         |       argument LEFTSHIFT argument
//         |       nonAssignExpression LEFTSHIFT nonAssignExpression
//         |       argument LEFTSHIFT nonAssignExpression
//         |       nonAssignExpression LEFTSHIFT argument

//         |       argument RIGHTSHIFT argument
//         |       nonAssignExpression RIGHTSHIFT nonAssignExpression
//         |       argument RIGHTSHIFT nonAssignExpression
//         |       nonAssignExpression RIGHTSHIFT argument

//         |       argument EXP argument
//         |       nonAssignExpression EXP nonAssignExpression
//         |       argument EXP nonAssignExpression
//         |       nonAssignExpression EXP argument

//         |       argument EQUAL argument
//         |       nonAssignExpression EQUAL nonAssignExpression
//         |       argument EQUAL nonAssignExpression
//         |       nonAssignExpression EQUAL argument

//         |       argument UNEQUAL argument
//         |       nonAssignExpression UNEQUAL nonAssignExpression
//         |       argument UNEQUAL nonAssignExpression
//         |       nonAssignExpression UNEQUAL argument

//         |       argument GTHAN argument
//         |       nonAssignExpression GTHAN nonAssignExpression
//         |       argument GTHAN nonAssignExpression
//         |       nonAssignExpression GTHAN argument

//         |       argument GTHANEQ argument
//         |       nonAssignExpression GTHANEQ nonAssignExpression
//         |       argument GTHANEQ nonAssignExpression
//         |       nonAssignExpression GTHANEQ argument

//         |       argument LTHAN argument
//         |       nonAssignExpression LTHAN nonAssignExpression
//         |       argument LTHAN nonAssignExpression
//         |       nonAssignExpression LTHAN argument

//         |       argument LTHANEQ argument
//         |       nonAssignExpression LTHANEQ nonAssignExpression
//         |       argument LTHANEQ nonAssignExpression
//         |       nonAssignExpression LTHANEQ argument

//         |       argument OR argument
//         |       nonAssignExpression OR nonAssignExpression
//         |       argument OR nonAssignExpression
//         |       nonAssignExpression OR argument

//         |       argument AND argument
//         |       nonAssignExpression AND nonAssignExpression
//         |       argument AND nonAssignExpression
//         |       nonAssignExpression AND argument

//         |       argument BITAND argument
//         |       nonAssignExpression BITAND nonAssignExpression
//         |       argument BITAND nonAssignExpression
//         |       nonAssignExpression BITAND argument

//         |       argument BITOR argument
//         |       nonAssignExpression BITOR nonAssignExpression
//         |       argument BITOR nonAssignExpression
//         |       nonAssignExpression BITOR argument

//         |       argument BITXOR argument
//         |       nonAssignExpression BITXOR nonAssignExpression
//         |       argument BITXOR nonAssignExpression
//         |       nonAssignExpression BITXOR argument

//         |       argument BITCLEAR argument
//         |       nonAssignExpression BITCLEAR nonAssignExpression
//         |       argument BITCLEAR nonAssignExpression
//         |       nonAssignExpression BITCLEAR argument

//         |       NOT argument
//         |       NOT nonAssignExpression
//         |       LPAREN argument RPAREN
//         |       BYTENUM
//         |       INT
//         |       LONG
//         |       FLOAT
//         |       DOUBLE
//         |       BOOLEAN
//         |       STRING
//         |       IDENT
//         |       NEW IDENT LBRACE structLitDefs RBRACE
//         |       IDENT LBRACK INT RBRACK
//         |       INFER LBRACE inferClauses RBRACE
//         |       BASICTYPE LBRACE argumentsList RBRACE
//                 // empty arrays
//         |       BASICTYPE LBRACE RBRACE
//         ;

// arguments:
//                 LPAREN argumentsList RPAREN
//         |       LPAREN RPAREN
//         ;

// argumentsList:  argument
//         |       nonAssignExpression
//         |       ADDR argument
//         |       VALUE argument
//         |       VALUE nonAssignExpression
                
//         |       argumentsList COMMA argument
//         |       argumentsList COMMA nonAssignExpression
//         |       argumentsList COMMA ADDR argument
//         |       argumentsList COMMA VALUE argument
//         |       argumentsList COMMA VALUE nonAssignExpression
//         ;

%%
