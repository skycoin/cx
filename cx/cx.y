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
			arg.TotalSize = size
			dataOffset += size
			prgrm.Data = append(prgrm.Data, Data(byts)...)
			expr := MakeExpression(nil)
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

		// getting offset to use by statements (excluding inputs, outputs and receiver)
		var offset int
		if len(fn.Outputs) > 0 {
			lastOutput := fn.Outputs[len(fn.Outputs) - 1]
			offset = lastOutput.Offset + lastOutput.TotalSize
		} else if len(fn.Inputs) > 0 {
			lastInput := fn.Inputs[len(fn.Inputs) - 1]
			offset = lastInput.Offset + lastInput.TotalSize
		}

		// for _, out := range outputs {
		// 	out.Offset += offset
		// 	offset += out.TotalSize
		// }

		for _, expr := range exprs {
			fn.AddExpression(expr)
		}

		fn.Length = len(fn.Expressions)

		var symbols map[string]*CXArgument = make(map[string]*CXArgument, 0)
		for _, inp := range fn.Inputs {
			if inp.Name != "" {
				symbols[inp.Name] = inp
			}
		}
		for _, out := range fn.Outputs {
			if out.Name != "" {
				symbols[out.Name] = out
			}
		}

		for _, expr := range fn.Expressions {
			for _, inp := range expr.Inputs {
				if inp.Name != "" {
					if arg, found := symbols[inp.Name]; !found {
						// it should exist. error
						panic("identifier '" + inp.Name + "' does not exist")
					} else {
						inp.Offset = arg.Offset

						if inp.IsStruct {
							// checking if it's accessing fields
							inp.Offset = arg.Offset
							// this will only work for one field atm
							if len(inp.Fields) > 0 {
								var found bool
								for _, fld := range arg.CustomType.Fields {
									if inp.Fields[0].Name == fld.Name {
										inp.Fields[0].Lengths = fld.Lengths
										inp.Fields[0].Size = fld.Size
										inp.Fields[0].TotalSize = fld.TotalSize
										found = true
										break
									}
									inp.Offset += fld.TotalSize
								}
								if !found {
									panic("field '" + inp.Fields[0].Name + "' not found")
								}
							}
						} else {
							inp.Offset = arg.Offset
						}
						
						inp.Lengths = arg.Lengths
						inp.Size = arg.Size
						inp.TotalSize = arg.TotalSize
					}
				}
			}
			for _, out := range expr.Outputs {
				if out.Name != "" {
					if arg, found := symbols[out.Name]; !found {
						out.Offset = offset
						symbols[out.Name] = out
						offset += out.TotalSize
					} else {
						if out.IsStruct {
							// checking if it's accessing fields
							out.Offset = arg.Offset
							// this will only work for one field atm
							if len(out.Fields) > 0 {
								var found bool
								for _, fld := range arg.CustomType.Fields {
									if out.Fields[0].Name == fld.Name {
										out.Fields[0].Lengths = fld.Lengths
										out.Fields[0].Size = fld.Size
										out.Fields[0].TotalSize = fld.TotalSize
										found = true
										break
									}
									out.Offset += fld.TotalSize
								}
								if !found {
									panic("field '" + out.Fields[0].Name + "' not found")
								}
							}
						} else {
							out.Offset = arg.Offset
						}

						out.Lengths = arg.Lengths
						out.Size = arg.Size
						out.TotalSize = arg.TotalSize
					}
				}
			}
		}
		fn.Size = offset
	}

	func FunctionCall (exprs []*CXExpression, args []*CXExpression) []*CXExpression {
		expr := exprs[0]
		if expr.Operator == nil {
			if op, err := prgrm.GetFunction(expr.Outputs[0].Name, expr.Outputs[0].Package.Name); err == nil {
				expr.Operator = op
			} else {
				panic(err)
			}
		}

		var nestedExprs []*CXExpression
		for _, inpExpr := range args {
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
		
		return append(nestedExprs, exprs...)
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
                        ADD_OP SUB_OP MUL_OP DIV_OP MOD_OP REF_OP AFFVAR
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

%type   <expressions>   declaration
//                      %type   <expressions>   init_declarator_list
//                      %type   <expressions>   init_declarator

%type   <expressions>   initializer
%type   <expressions>   initializer_list
%type   <expressions>   designation
%type   <expressions>   designator_list
%type   <expressions>   designator

%type   <expressions>   expression
%type   <expressions>   block_item
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
        |       struct_declaration
        ;

// parameter_declaration

struct_declaration:
                TYPE IDENTIFIER STRUCT struct_fields
                {
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
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
                LBRACE RBRACE
                { $$ = nil }
        |       LBRACE fields RBRACE
                { $$ = $2 }
        ;

fields:         parameter_declaration SEMICOLON
                {
			if $1.IsArray {
				$1.TotalSize = $1.Size * TotalLength($1.Lengths)
			} else {
				$1.TotalSize = $1.Size
			}
			$$ = []*CXArgument{$1}
                }
        |       fields parameter_declaration SEMICOLON
                {
			if $2.IsArray {
				$2.TotalSize = $2.Size * TotalLength($2.Lengths)
			} else {
				$2.TotalSize = $2.Size
			}
			$$ = append($1, $2)
                }
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
        |       FUNC LPAREN parameter_type_list RPAREN IDENTIFIER
                {
			if len($3) > 1 {
				panic("method has multiple receivers")
			}
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
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
			FunctionDeclaration($1, $2, nil, $3)
                }
        |       function_header function_parameters function_parameters compound_statement
                {
			FunctionDeclaration($1, $2, $3, $4)
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
			if $1.IsArray {
				$1.TotalSize = $1.Size * TotalLength($1.Lengths)
			} else {
				$1.TotalSize = $1.Size
			}
			$$ = []*CXArgument{$1}
                }
	|       parameter_list COMMA parameter_declaration
                {
			if $3.IsArray {
				$3.TotalSize = $3.Size * TotalLength($3.Lengths)
			} else {
				$3.TotalSize = $3.Size
			}
			lastPar := $1[len($1) - 1]
			$3.Offset = lastPar.Offset + lastPar.TotalSize
			$$ = append($1, $3)
                }
                ;

parameter_declaration:
                declarator declaration_specifiers
                {
			$2.Name = $1.Name
			// $2.IsArray = $1.IsArray
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

declarator:     // MUL_OP direct_declarator
        //         {
        //             $2.IsPointer = true
        //             $$ = $2
        //         }
        // |       
		direct_declarator
                ;

direct_declarator:
                IDENTIFIER
                {
			arg := MakeArgument(TYPE_UNDEFINED)
			arg.Name = $1
			$$ = arg
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








// declaration_specifiers
declaration_specifiers:
                unary_operator declaration_specifiers
                {
			arg := $2
			$2.IsPointer = true
			$$ = arg
                }
        |       LBRACK INT_LITERAL RBRACK declaration_specifiers
                {
			arg := $4
                        arg.IsArray = true
			arg.Lengths = append([]int{int($2)}, arg.Lengths...)
			arg.TotalSize = arg.Size * TotalLength(arg.Lengths)
			// arg.Size = GetArgSize($4.Type)
			$$ = arg
                }
        |       type_specifier
                {
			arg := MakeArgument($1)
			arg.Size = GetArgSize($1)
			arg.TotalSize = arg.Size
			$$ = arg
                }
        |       IDENTIFIER
                {
			// custom type in the current package
			if pkg, err := prgrm.GetCurrentPackage(); err == nil {
				if strct, err := prgrm.GetStruct($1, pkg.Name); err == nil {
					arg := MakeArgument(TYPE_CUSTOM)
					arg.CustomType = strct
					arg.Size = strct.Size
					arg.TotalSize = strct.Size

					// for _, fld := range strct.Fields {
					// 	arg.Sizes = append(arg.Sizes, fld.Size)
					// }

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
	|       postfix_expression LBRACK expression RBRACK
                {

			$1[0].Outputs[0].IsArray = true
			// $1[0].Outputs[0].NumIndexes += 1

			if len($1[0].Outputs[0].Fields) > 0 {
				fld := $1[0].Outputs[0].Fields[len($1[0].Outputs[0].Fields) - 1]
				fld.Indexes = append(fld.Indexes, $3[0].Outputs[0])
			} else {
				$1[0].Outputs[0].Indexes = append($1[0].Outputs[0].Indexes, $3[0].Outputs[0])
			}
			
			expr := $1[len($1) - 1]
			// expr.Operator = Natives[OP_READ_ARRAY]
			if len(expr.Inputs) < 1 {
				expr.Inputs = append(expr.Inputs, $1[0].Outputs[0])
			}
			expr.Inputs = append(expr.Inputs, $3[0].Outputs[0])

			$$ = $1
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
	|       postfix_expression LPAREN RPAREN
                {
			$$ = FunctionCall($1, nil)
                }
	|       postfix_expression LPAREN argument_expression_list RPAREN
                {
			$$ = FunctionCall($1, $3)
                }
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
			left := $1[0].Outputs[0]
			// right := $3[0].Outputs[0]

			if left.IsRest {
				// fmt.Println("first")
				// then it can't be a module name
				// and we propagate the property to the right expression
				// right.IsRest = true
			} else {
				// fmt.Println("second")
				// then left is a first (e.g first.rest) and right is a rest
				// right.IsRest = true
				// let's check if left is a package
				if _, err := prgrm.GetPackage(left.Name); err == nil {
					// fmt.Println("second first")
					// the external property will be propagated to the following arguments
					// this way we avoid considering these arguments as module names
					//right.Package = pkg
					$$ = $1
				} else {
					if opCode, ok := OpCodes[$1[0].Outputs[0].Name + "." + $3]; ok {
						$1[0].Operator = Natives[opCode]
					} else {
						left.IsStruct = true
						fld := MakeArgument(TYPE_IDENTIFIER)
						fld.Name = $3
						left.Fields = append(left.Fields, fld)
					}
					
					// if fld, err := left
					// left.Offset += 
					
					// left is not a package, then it's a struct or a function call
					// if right.Operator != nil {
					// 	// then left is a function call that returns a struct instance
					// 	// we just return an expression with the GetField operator
						
					// } else {
					// 	// then left is a struct instance

					// 	// GetField(left, right)
					// }
				}
			}
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
                REF_OP
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
	|       and_expression REF_OP equality_expression
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
				$1[0].Outputs[0].Lengths = $3[idx].Outputs[0].Lengths
				$3[idx].Inputs = $3[idx].Outputs
				$3[idx].Outputs = $1[0].Outputs
				
				$$ = $3
			} else {
				if $3[idx].Operator.IsNative {
					for i, out := range $3[idx].Operator.Outputs {
						$1[0].Outputs[i].Size = Natives[$3[idx].Operator.OpCode].Outputs[i].Size
						$1[0].Outputs[i].Lengths = out.Lengths
					}
				} else {
					for i, out := range $3[idx].Operator.Outputs {
						$1[0].Outputs[i].Size = out.Size
						$1[0].Outputs[i].Lengths = out.Lengths
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







declaration:
                VAR declarator declaration_specifiers SEMICOLON
                {
			// this will tell the runtime that it's just a declaration
			expr := MakeExpression(nil)

			$3.Name = $2.Name
			expr.AddOutput($3)

			$$ = []*CXExpression{expr}
                }
        |       VAR declarator declaration_specifiers ASSIGN initializer SEMICOLON
                {
			$$ = nil
                }
                ;

/* init_declarator: */
/*                 declarator '=' initializer */
/*                 { */
/*                     $$ = nil */
/*                 } */
/*         |       declarator */
/*                 { */
/*                     $$ = nil */
/*                 } */
/*                 ; */


/* init_declarator_list: */
/*                 init_declarator */
/*                 { */
/*                     $$ = nil */
/*                 } */
/* 	|       init_declarator_list COMMA init_declarator */
/*                 { */
/*                     $$ = nil */
/*                 } */
/*                 ; */

/* init_declarator: */
/*                 declarator '=' initializer */
/*                 { */
/*                     $$ = nil */
/*                 } */
/*         |       declarator */
/*                 { */
/*                     $$ = nil */
/*                 } */
/*                 ; */






initializer:
        /*         LBRACE initializer_list RBRACE */
	/* |       LBRACE   initializer_list COMMA RBRACE */
	/* |        */assignment_expression
                ;

initializer_list:
                designation initializer
                {
                    $$ = nil
                }
	|       initializer
                {
                    $$ = nil
                }
	|       initializer_list COMMA designation initializer
                {
                    $$ = nil
                }
	|       initializer_list COMMA initializer
                {
			$$ = nil
                }
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

%%
