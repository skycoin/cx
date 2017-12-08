%{
	package main
	import (
		"strings"
		"bytes"
		"fmt"
		"os"
		"time"
		
		"github.com/skycoin/skycoin/src/cipher/encoder"
		. "github.com/skycoin/cx/src/base"
                )

	var program bytes.Buffer
	
	var cxt = MakeContext()

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
	//var dProgram bool = false
	var tag string = ""
	var asmNL = "\n"

	func warnf(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format, args...)
		os.Stderr.Sync()
	}

	func binaryOp (op string, arg1, arg2 *CXArgument, line int) *CXArgument {
		var opName string
		var typArg1 string
		// var typArg2 string
		// _ = typArg2

		if arg1.Typ == "ident" {
			var identName string
			encoder.DeserializeRaw(*arg1.Value, &identName)

			if typ, err := GetIdentType(identName, cxt); err == nil {
				typArg1 = typ
			} else {
				fmt.Println(err)
				//typArg1 = "ident"
			}
		} else {
			typArg1 = arg1.Typ
		}

		// if arg2.Typ == "ident" {
		// 	var identName string
		// 	encoder.DeserializeRaw(*arg2.Value, &identName)

		// 	if typ, err := GetIdentType(identName, cxt); err == nil {
		// 		typArg2 = typ
		// 	} else {
		// 		fmt.Println(err)
		// 	}
		// } else {
		// 	typArg2 = arg1.Typ
		// }
		
		switch op {
		case "+":
			opName = fmt.Sprintf("%s.add", typArg1)
		case "-":
			opName = fmt.Sprintf("%s.sub", typArg1)
		case "*":
			opName = fmt.Sprintf("%s.mul", typArg1)
		case "/":
			opName = fmt.Sprintf("%s.div", typArg1)
		case "%":
			opName = fmt.Sprintf("%s.mod", typArg1)
		case ">":
			opName = fmt.Sprintf("%s.gt", typArg1)
		case "<":
			opName = fmt.Sprintf("%s.lt", typArg1)
		case "<=":
			opName = fmt.Sprintf("%s.lteq", typArg1)
		case ">=":
			opName = fmt.Sprintf("%s.gteq", typArg1)
		case "<<":
			opName = fmt.Sprintf("%s.bitshl", typArg1)
		case ">>":
			opName = fmt.Sprintf("%s.bitshr", typArg1)
		case "**":
			opName = fmt.Sprintf("%s.pow", typArg1)
		case "&&":
			opName = "and"
		case "||":
			opName = "or"
		case "==":
			opName = fmt.Sprintf("%s.eq", typArg1)
		case "!=":
			opName = fmt.Sprintf("%s.uneq", typArg1)
		}

		if fn, err := cxt.GetCurrentFunction(); err == nil {
			if op, err := cxt.GetFunction(opName, CORE_MODULE); err == nil {
				expr := MakeExpression(op)
				if !replMode {
					expr.FileLine = line
				}
				fn.AddExpression(expr)
				expr.AddTag(tag)
				tag = ""
				expr.AddArgument(arg1)
				expr.AddArgument(arg2)

				outName := MakeGenSym(NON_ASSIGN_PREFIX)
				byteName := encoder.Serialize(outName)
				
				expr.AddOutputName(outName)
				return MakeArgument(&byteName, "ident")
			}
		}
		return nil
	}

        func unaryOp (op string, arg1 *CXArgument, line int) *CXArgument {
		var opName string
		var typArg1 string

		if arg1.Typ == "ident" {
			var identName string
			encoder.DeserializeRaw(*arg1.Value, &identName)

			if typ, err := GetIdentType(identName, cxt); err == nil {
				typArg1 = typ
			} else {
				fmt.Println(err)
			}
		} else {
			typArg1 = arg1.Typ
		}

		switch op {
		case "++":
			opName = fmt.Sprintf("%s.add", typArg1)
		case "--":
			opName = fmt.Sprintf("%s.sub", typArg1)
		}
		
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			if op, err := cxt.GetFunction(opName, CORE_MODULE); err == nil {
				expr := MakeExpression(op)
				if !replMode {
					expr.FileLine = line
				}
				fn.AddExpression(expr)
				expr.AddTag(tag)
				tag = ""

				
				expr.AddArgument(arg1)

				// var one *CXArgument

				switch typArg1 {
				case "i32":
					sOne := encoder.Serialize(int32(1))
					expr.AddArgument(MakeArgument(&sOne, "i32"))
				case "i64":
					sOne := encoder.Serialize(int64(1))
					expr.AddArgument(MakeArgument(&sOne, "i64"))
				case "f32":
					sOne := encoder.Serialize(float32(1))
					expr.AddArgument(MakeArgument(&sOne, "f32"))
				case "f64":
					sOne := encoder.Serialize(float64(1))
					expr.AddArgument(MakeArgument(&sOne, "f64"))
				}

				var outName string
				if arg1.Typ == "ident" {
					encoder.DeserializeRaw(*arg1.Value, &outName)
				} else {
					outName = MakeGenSym(NON_ASSIGN_PREFIX)
				}
				
				byteName := encoder.Serialize(outName)
				
				expr.AddOutputName(outName)
				return MakeArgument(&byteName, "ident")
			}
		}
		return nil
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

	parameter *CXParameter
	parameters []*CXParameter

	argument *CXArgument
	arguments []*CXArgument

        definition *CXDefinition
	definitions []*CXDefinition

	expression *CXExpression
	expressions []*CXExpression

	field *CXField
	fields []*CXField

	name string
	names []string
}

%token  <byt>           BYTENUM
%token  <i32>           INT BOOLEAN
%token  <i64>           LONG
%token  <f32>           FLOAT
%token  <f64>           DOUBLE
%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK IDENT
                        VAR COMMA COMMENT STRING PACKAGE IF ELSE FOR TYPSTRUCT STRUCT
                        ASSIGN CASSIGN IMPORT RETURN GOTO GTHAN LTHAN EQUAL COLON NEW
                        EQUALWORD GTHANWORD LTHANWORD
                        GTHANEQ LTHANEQ UNEQUAL AND OR
                        PLUS MINUS MULT DIV AFFVAR
                        PLUSPLUS MINUSMINUS REMAINDER LEFTSHIFT RIGHTSHIFT EXP
                        NOT
                        BITAND BITXOR BITOR BITCLEAR
                        PLUSEQ MINUSEQ MULTEQ DIVEQ REMAINDEREQ EXPEQ
                        LEFTSHIFTEQ RIGHTSHIFTEQ BITANDEQ BITXOREQ BITOREQ
                        /* Types */
                        BOOL STR I32 I64 F32 F64 BYTE BOOLA BYTEA I32A I64A F32A F64A STRA
                        /* Selectors */
                        SPACKAGE SSTRUCT SFUNC
                        /* Removers */
                        REM DEF EXPR FIELD INPUT OUTPUT CLAUSES OBJECT OBJECTS
                        /* Stepping */
                        STEP PSTEP TSTEP
                        /* Debugging */
                        DSTACK DPROGRAM DSTATE
                        /* Affordances */
                        AFF TAG INFER WEIGHT

%type   <tok>           typeSpecifier assignOperator relationalOp
                        
%type   <parameter>     parameter
%type   <parameters>    parameters functionParameters
%type   <argument>      argument definitionAssignment
%type   <arguments>     arguments argumentsList nonAssignExpression conditionControl
%type   <definition>    structLitDef
%type   <definitions>   structLitDefs structLiteral
%type   <fields>        fields structFields
%type   <expression>    assignExpression
%type   <expressions>   elseStatement

%type   <bool>          selectorLines selectorExpressionsAndStatements selectorFields
%type   <stringA>       inferObj inferObjs inferRule inferRules inferClauses inferPred inferCond inferAction inferActions inferActionArg inferTarget inferTargets
%type   <string>        inferWeight

%left                   OR
%left                   AND
%left                   BITOR
%left                   BITXOR
%left                   BITAND
%left                   EQUAL UNEQUAL
%left                   GTHAN LTHAN GTHANEQ LTHANEQ
%left                   LEFTSHIFT RIGHTSHIFT
%left                   PLUS MINUS
%left                   REMAINDER MULT DIV
%left                   EXP
%left                   PLUSPLUS MINUSMINUS
%left                   NOT
%right                  LPAREN IDENT

                        
%%

lines:
                /* empty */
        |       lines line
                {
			// if replMode && dProgram {
			// 	cxt.PrintProgram(false)
			// }
                }
        |       lines ';'
                {
                    
                }
        ;

line:
                definitionDeclaration
        |       structDeclaration
        |       packageDeclaration
        |       importDeclaration
        |       functionDeclaration
        |       selector
        |       stepping
        |       debugging
        |       affordance
        |       remover
        ;

importDeclaration:
                IMPORT STRING
                {
			impName := strings.TrimPrefix($2, "\"")
			impName = strings.TrimSuffix(impName, "\"")
			if imp, err := cxt.GetModule(impName); err == nil {
				if mod, err := cxt.GetCurrentModule(); err == nil {
					mod.AddImport(imp)
				}
			}
                }
        ;

affordance:
                TAG
                {
			tag = strings.TrimSuffix($1, ":")
                }
                /* Function Affordances */
        |       AFF FUNC IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetFunction($3, mod.Name); err == nil {
					affs := fn.GetAffordances()
					for i, aff := range affs {
						fmt.Printf("(%d)\t%s\n", i, aff.Description)
					}
				}
			}
                }
        |       AFF FUNC IDENT LBRACE INT RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetFunction($3, mod.Name); err == nil {
					affs := fn.GetAffordances()
					affs[$5].ApplyAffordance()
				}
			}
                }
        |       AFF FUNC IDENT LBRACE STRING RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetFunction($3, mod.Name); err == nil {
					affs := fn.GetAffordances()
					filter := strings.TrimPrefix($5, "\"")
					filter = strings.TrimSuffix(filter, "\"")
					affs = FilterAffordances(affs, filter)
					for i, aff := range affs {
						fmt.Printf("(%d)\t%s\n", i, aff.Description)
					}
				}
			}
                }
        |       AFF FUNC IDENT LBRACE STRING INT RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetFunction($3, mod.Name); err == nil {
					affs := fn.GetAffordances()
					filter := strings.TrimPrefix($5, "\"")
					filter = strings.TrimSuffix(filter, "\"")
					affs = FilterAffordances(affs, filter)
					affs[$6].ApplyAffordance()
				}
			}
                }
                /* Module Affordances */
        |       AFF PACKAGE IDENT
                {
			if mod, err := cxt.GetModule($3); err == nil {
				affs := mod.GetAffordances()
				for i, aff := range affs {
					fmt.Printf("(%d)\t%s\n", i, aff.Description)
				}
			}
                }
        |       AFF PACKAGE IDENT LBRACE INT RBRACE
                {
			if mod, err := cxt.GetModule($3); err == nil {
				affs := mod.GetAffordances()
				affs[$5].ApplyAffordance()
			}
                }
        |       AFF PACKAGE IDENT LBRACE STRING RBRACE
                {
			if mod, err := cxt.GetModule($3); err == nil {
				affs := mod.GetAffordances()
				filter := strings.TrimPrefix($5, "\"")
				filter = strings.TrimSuffix(filter, "\"")
				affs = FilterAffordances(affs, filter)
				for i, aff := range affs {
					fmt.Printf("(%d)\t%s\n", i, aff.Description)
				}
			}
                }
        |       AFF PACKAGE IDENT LBRACE STRING INT RBRACE
                {
			if mod, err := cxt.GetModule($3); err == nil {
				affs := mod.GetAffordances()
				filter := strings.TrimPrefix($5, "\"")
				filter = strings.TrimSuffix(filter, "\"")
				affs = FilterAffordances(affs, filter)
				affs[$6].ApplyAffordance()
			}
                }
                /* Struct Affordances */
        |       AFF STRUCT IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if strct, err := cxt.GetStruct($3, mod.Name); err == nil {
					affs := strct.GetAffordances()
					for i, aff := range affs {
						fmt.Printf("(%d)\t%s\n", i, aff.Description)
					}
				}
			}
			
                }
        |       AFF STRUCT IDENT LBRACE INT RBRACE
                {

			if mod, err := cxt.GetCurrentModule(); err == nil {
				if strct, err := cxt.GetStruct($3, mod.Name); err == nil {
					affs := strct.GetAffordances()
					affs[$5].ApplyAffordance()
				}
			}
                }
        |       AFF STRUCT IDENT LBRACE STRING RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if strct, err := cxt.GetStruct($3, mod.Name); err == nil {
					affs := strct.GetAffordances()
					filter := strings.TrimPrefix($5, "\"")
					filter = strings.TrimSuffix(filter, "\"")
					affs = FilterAffordances(affs, filter)
					for i, aff := range affs {
						fmt.Printf("(%d)\t%s\n", i, aff.Description)
					}
				}
			}
                }
        |       AFF STRUCT IDENT LBRACE STRING INT RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if strct, err := cxt.GetStruct($3, mod.Name); err == nil {
					affs := strct.GetAffordances()
					filter := strings.TrimPrefix($5, "\"")
					filter = strings.TrimSuffix(filter, "\"")
					affs = FilterAffordances(affs, filter)
					affs[$6].ApplyAffordance()
				}
			}
                }
                /* Struct Affordances */
        |       AFF EXPR IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, expr := range fn.Expressions {
						if expr.Tag == $3 {
							PrintAffordances(expr.GetAffordances(nil))
							break
						}
					}
					// if expr, err := fn.GetCurrentExpression(); err == nil {
					// 	affs := expr.GetAffordances()
					// 	for i, aff := range affs {
					// 		fmt.Printf("(%d)\t%s\n", i, aff.Description)
					// 	}
					// }
				}
			}
                }
        |       AFF EXPR IDENT LBRACE INT RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, expr := range fn.Expressions {
						if expr.Tag == $3 {
							affs := expr.GetAffordances(nil)
							affs[$5].ApplyAffordance()
							break
						}
					}
				}
			}
                }
        |       AFF EXPR IDENT LBRACE STRING RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, expr := range fn.Expressions {
						if expr.Tag == $3 {
							affs := expr.GetAffordances(nil)
							filter := strings.TrimPrefix($5, "\"")
							filter = strings.TrimSuffix(filter, "\"")
							PrintAffordances(FilterAffordances(affs, filter))
							break
						}
					}
				}
			}
                }
        |       AFF EXPR IDENT LBRACE STRING INT RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, expr := range fn.Expressions {
						if expr.Tag == $3 {
							affs := expr.GetAffordances(nil)
							filter := strings.TrimPrefix($5, "\"")
							filter = strings.TrimSuffix(filter, "\"")
							affs = FilterAffordances(affs, filter)
							affs[$6].ApplyAffordance()
							break
						}
					}
				}
			}
                }
        ;

stepping:       TSTEP INT INT
                {
			if $2 == 0 {
				// Maybe nothing for now
			} else {
				if $2 < 0 {
					nCalls := $2 * -1
					for i := int32(0); i < nCalls; i++ {
						time.Sleep(time.Duration(int32($3)) * time.Millisecond)
						cxt.UnRun(1)
					}
				} else {

					for i := int32(0); i < $2; i++ {
						time.Sleep(time.Duration(int32($3)) * time.Millisecond)
						err := cxt.Run(dStack, 1)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
                }
        |       STEP INT
                {
			if $2 == 0 {
				// we run until halt or end of program;
				if err := cxt.Run(dStack, -1); err != nil {
					fmt.Println(err)
				}
			} else {
				if $2 < 0 {
					nCalls := $2 * -1
					cxt.UnRun(int(nCalls))
				} else {
					//fmt.Println(cxt.Run(dStack, int($2)))

					err := cxt.Run(dStack, int($2))
					if err != nil {
						fmt.Println(err)
					}
					
					// if err := cxt.Run(dStack, int($2)); err != nil {
					// 	fmt.Println(err)
					// }
				}
			}
                }
        // |       PSTEP INT
        //         {
			
        //         }
        ;

debugging:      DSTATE
                {
			if len(cxt.CallStack.Calls) > 0 {
				for _, def := range cxt.CallStack.Calls[len(cxt.CallStack.Calls) - 1].State {
					fmt.Printf("%s(%s):\t\t%s\n", def.Name, def.Typ, PrintValue(def.Name, def.Value, def.Typ, cxt))
				}
			}
                }
        |       DSTACK BOOLEAN
                {
			if $2 > 0 {
				dStack = true
                        } else {
				dStack = false
			}
                }
        |       DPROGRAM
                {
			cxt.PrintProgram(false)
			// if $2 > 0 {
			// 	dProgram = true
                        // } else {
			// 	dProgram = false
			// }
                }
        ;

remover:        REM FUNC IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.RemoveFunction($3)
			}
                }
        |       REM PACKAGE IDENT
                {
			cxt.RemoveModule($3)
                }
        |       REM DEF IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.RemoveDefinition($3)
			}
                }
        |       REM STRUCT IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.RemoveStruct($3)
			}
                }
        |       REM IMPORT STRING
                {
			impName := strings.TrimPrefix($3, "\"")
			impName = strings.TrimSuffix(impName, "\"")
			
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.RemoveImport(impName)
			}
                }
        |       REM EXPR IDENT FUNC IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.Context.GetFunction($5, mod.Name); err == nil {
					for i, expr := range fn.Expressions {
						if expr.Tag == $3 {
							fn.RemoveExpression(i)
						}
					}
				}
			}
                }
        |       REM FIELD IDENT STRUCT IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if strct, err := cxt.GetStruct($5, mod.Name); err == nil {
					strct.RemoveField($3)
				}
				
			}
                }
        |       REM INPUT IDENT FUNC IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.Context.GetFunction($5, mod.Name); err == nil {
					fn.RemoveInput($3)
				}
			}
                }
        |       REM OUTPUT IDENT FUNC IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.Context.GetFunction($5, mod.Name); err == nil {
					fn.RemoveOutput($3)
				}
			}
                }
                // no, too complex. just wipe out entire expression
        // |       REM ARG INT EXPR INT FUNC IDENT
        //         {
	// 		if mod, err := cxt.GetCurrentModule(); err == nil {
	// 			if fn, err := mod.Context.GetFunction($5, mod.Name); err == nil {
	// 				fn.RemoveExpression(int($3))
	// 			}
	// 		}
        //         }
        // |       REM OUTNAME
        ;

selectorLines:
                /* empty */
                {
			$$ = false
                }
        |       LBRACE lines RBRACE
                {
			$$ = true
                }
        ;

selectorExpressionsAndStatements:
                /* empty */
                {
			$$ = false
                }
        |       LBRACE expressionsAndStatements RBRACE
                {
			$$ = true
                }
        ;

selectorFields:
                /* empty */
                {
			$$ = false
                }
        |       LBRACE fields RBRACE
                {
			if strct, err := cxt.GetCurrentStruct(); err == nil {
				for _, fld := range $2 {
					fldFromParam := MakeField(fld.Name, fld.Typ)
					strct.AddField(fldFromParam)
				}
			}
			$$ = true
                }
        ;

selector:       SPACKAGE IDENT
                {
			var previousModule *CXModule
			if mod, err := cxt.GetCurrentModule(); err == nil {
				previousModule = mod
			} else {
				fmt.Println("A current module does not exist")
			}
			if _, err := cxt.SelectModule($2); err == nil {
				//fmt.Println(fmt.Sprintf("== Changed to package '%s' ==", mod.Name))
			} else {
				fmt.Println(err)
			}

			replTargetMod = $2
			replTargetStrct = ""
			replTargetFn = ""
			
			$<string>$ = previousModule.Name
                }
                selectorLines
                {
			if $<bool>4 {
				if _, err := cxt.SelectModule($<string>3); err == nil {
					//fmt.Println(fmt.Sprintf("== Changed to package '%s' ==", mod.Name))
				}
			}
                }
        |       SFUNC IDENT
                {
			var previousFunction *CXFunction
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				previousFunction = fn
			} else {
				fmt.Println("A current function does not exist")
			}
			if _, err := cxt.SelectFunction($2); err == nil {
				//fmt.Println(fmt.Sprintf("== Changed to function '%s' ==", fn.Name))
			} else {
				fmt.Println(err)
			}

			replTargetMod = ""
			replTargetStrct = ""
			replTargetFn = $2
			
			$<string>$ = previousFunction.Name
                }
                selectorExpressionsAndStatements
                {
			if $<bool>4 {
				if _, err := cxt.SelectFunction($<string>3); err == nil {
					//fmt.Println(fmt.Sprintf("== Changed to function '%s' ==", fn.Name))
				}
			}
                }
        |       SSTRUCT IDENT
                {
			var previousStruct *CXStruct
			if fn, err := cxt.GetCurrentStruct(); err == nil {
				previousStruct = fn
			} else {
				fmt.Println("A current struct does not exist")
			}
			if _, err := cxt.SelectStruct($2); err == nil {
				//fmt.Println(fmt.Sprintf("== Changed to struct '%s' ==", fn.Name))
			} else {
				fmt.Println(err)
			}

			replTargetStrct = $2
			replTargetMod = ""
			replTargetFn = ""
			
			$<string>$ = previousStruct.Name
                }
                selectorFields
                {
			if $<bool>4 {
				if _, err := cxt.SelectStruct($<string>3); err == nil {
					//fmt.Println(fmt.Sprintf("== Changed to struct '%s' ==", strct.Name))
				}
			}
                }
        ;

assignOperator:
                ASSIGN
        |       CASSIGN
        |       PLUSEQ
        |       MINUSEQ
        |       MULTEQ
        |       DIVEQ
        |       REMAINDEREQ
        |       EXPEQ
        |       LEFTSHIFTEQ
        |       RIGHTSHIFTEQ
        |       BITANDEQ
        |       BITXOREQ
        |       BITOREQ
        ;

typeSpecifier:
                I32 {$$ = $1}
        |       I64 {$$ = $1}
        |       F32 {$$ = $1}
        |       F64 {$$ = $1}
        |       BOOL {$$ = $1}
        |       BYTE {$$ = $1}
        |       BOOLA {$$ = $1}
        |       STRA {$$ = $1}
        |       BYTEA {$$ = $1}
        |       I32A {$$ = $1}
        |       I64A {$$ = $1}
        |       F32A {$$ = $1}
        |       F64A {$$ = $1}
        |       STR {$$ = $1}
        |       NEW IDENT {$$ = $2}
        ;

packageDeclaration:
                PACKAGE IDENT
                {
			mod := MakeModule($2)
			cxt.AddModule(mod)
                }
                ;

definitionAssignment:
                /* empty */
                {
			$$ = nil
                }
        |       assignOperator argument
                {
			$$ = $2
                }
                ;

definitionDeclaration:
                VAR IDENT typeSpecifier definitionAssignment
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				var val *CXArgument;
				if $4 == nil {
					val = MakeArgument(MakeDefaultValue($3), $3)
				} else {
					switch $3 {
					case "byte":
						var ds int32
						encoder.DeserializeRaw(*$4.Value, &ds)
						//new := encoder.Serialize(byte(ds))
						new := []byte{byte(ds)}
						val = MakeArgument(&new, "byte")
					case "i64":
						var ds int32
						encoder.DeserializeRaw(*$4.Value, &ds)
						new := encoder.Serialize(int64(ds))
						val = MakeArgument(&new, "i64")
					case "f64":
						var ds float32
						encoder.DeserializeRaw(*$4.Value, &ds)
						new := encoder.Serialize(float64(ds))
						val = MakeArgument(&new, "f64")
					default:
						val = $4
					}
				}

				mod.AddDefinition(MakeDefinition($2, val.Value, $3))
			}
                }
        |       VAR IDENT IDENT
                {
			// we have to initialize all the fields
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if zeroVal, err := ResolveStruct($3, cxt); err == nil {
					mod.AddDefinition(MakeDefinition($2, &zeroVal, $3))
				} else {
					fmt.Println(fmt.Sprintf("%d: definition declaration: %s", yyS[yypt-0].line + 1, err))
				}
			}
                }
        ;

fields:
                parameter
                {
			var flds []*CXField
                        flds = append(flds, MakeFieldFromParameter($1))
			$$ = flds
                }
        |       ';'
                {
			var flds []*CXField
			$$ = flds
                }
        |       debugging
                {
			var flds []*CXField
			$$ = flds
                }
        |       fields parameter
                {
			$1 = append($1, MakeFieldFromParameter($2))
			$$ = $1
                }
        |       fields ';'
                {
			$$ = $1
                }
        |       fields debugging
                {
			$$ = $1
                }
                ;

structFields:
                LBRACE fields RBRACE
                {
			$$ = $2
                }
        |       LBRACE RBRACE
                {
			$$ = nil
                }
        ;

structDeclaration:
                TYPSTRUCT IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				strct := MakeStruct($2)
				mod.AddStruct(strct)

				// creating manipulation functions for this type a la common lisp
				// append
				fn := MakeFunction(fmt.Sprintf("[]%s.append", $2))
				fn.AddInput(MakeParameter("arr", fmt.Sprintf("[]%s", $2)))
				fn.AddInput(MakeParameter("strctInst", $2))
				fn.AddOutput(MakeParameter("_arr", fmt.Sprintf("[]%s", $2)))
				mod.AddFunction(fn)

				if op, err := cxt.GetFunction("cstm.append", CORE_MODULE); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = yyS[yypt-0].line + 1
					}
					expr.AddOutputName("_arr")
					sArr := encoder.Serialize("arr")
					arrArg := MakeArgument(&sArr, "str")
					sStrctInst := encoder.Serialize("strctInst")
					strctInstArg := MakeArgument(&sStrctInst, "str")
					expr.AddArgument(arrArg)
					expr.AddArgument(strctInstArg)
					fn.AddExpression(expr)
				} else {
					fmt.Println(err)
				}
				// read
				fn = MakeFunction(fmt.Sprintf("[]%s.read", $2))
				fn.AddInput(MakeParameter("arr", fmt.Sprintf("[]%s", $2)))
				fn.AddInput(MakeParameter("index", "i32"))
				fn.AddOutput(MakeParameter("strctInst", $2))
				mod.AddFunction(fn)

				if op, err := cxt.GetFunction("cstm.read", CORE_MODULE); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = yyS[yypt-0].line + 1
					}
					expr.AddOutputName("strctInst")
					sArr := encoder.Serialize("arr")
					arrArg := MakeArgument(&sArr, "str")
					sIndex := encoder.Serialize("index")
					indexArg := MakeArgument(&sIndex, "ident")
					expr.AddArgument(arrArg)
					expr.AddArgument(indexArg)
					fn.AddExpression(expr)
				} else {
					fmt.Println(err)
				}
				// write
				fn = MakeFunction(fmt.Sprintf("[]%s.write", $2))
				fn.AddInput(MakeParameter("arr", fmt.Sprintf("[]%s", $2)))
				fn.AddInput(MakeParameter("index", "i32"))
				fn.AddInput(MakeParameter("inst", $2))
				fn.AddOutput(MakeParameter("_arr", fmt.Sprintf("[]%s", $2)))
				mod.AddFunction(fn)

				if op, err := cxt.GetFunction("cstm.write", CORE_MODULE); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = yyS[yypt-0].line + 1
					}
					sArr := encoder.Serialize("arr")
					arrArg := MakeArgument(&sArr, "str")
					sIndex := encoder.Serialize("index")
					indexArg := MakeArgument(&sIndex, "ident")
					sInst := encoder.Serialize("inst")
					instArg := MakeArgument(&sInst, "str")
					expr.AddArgument(arrArg)
					expr.AddArgument(indexArg)
					expr.AddArgument(instArg)

					expr.AddOutputName("_arr")
					fn.AddExpression(expr)
				} else {
					fmt.Println(err)
				}
				// len
				fn = MakeFunction(fmt.Sprintf("[]%s.len", $2))
				fn.AddInput(MakeParameter("arr", fmt.Sprintf("[]%s", $2)))
				fn.AddOutput(MakeParameter("len", "i32"))
				mod.AddFunction(fn)

				if op, err := cxt.GetFunction("cstm.len", CORE_MODULE); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = yyS[yypt-0].line + 1
					}
					expr.AddOutputName("len")
					sArr := encoder.Serialize("arr")
					arrArg := MakeArgument(&sArr, "str")
					expr.AddArgument(arrArg)
					fn.AddExpression(expr)
				} else {
					fmt.Println(err)
				}
				// make
				fn = MakeFunction(fmt.Sprintf("[]%s.make", $2))
				fn.AddInput(MakeParameter("len", "i32"))
				//fn.AddInput(MakeParameter("typ", "str"))
				fn.AddOutput(MakeParameter("arr", fmt.Sprintf("[]%s", $2)))
				mod.AddFunction(fn)

				if op, err := cxt.GetFunction("cstm.make", CORE_MODULE); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = yyS[yypt-0].line + 1
					}
					expr.AddOutputName("arr")
					sLen := encoder.Serialize("len")
					sTyp := encoder.Serialize(fmt.Sprintf("[]%s", $2))
					lenArg := MakeArgument(&sLen, "ident")
					typArg := MakeArgument(&sTyp, "str")
					expr.AddArgument(lenArg)
					expr.AddArgument(typArg)
					fn.AddExpression(expr)
				} else {
					fmt.Println(err)
				}
			}
                }
                STRUCT structFields
                {
			if strct, err := cxt.GetCurrentStruct(); err == nil {
				for _, fld := range $5 {
					fldFromParam := MakeField(fld.Name, fld.Typ)
					strct.AddField(fldFromParam)
				}
			}
                }
        ;

functionParameters:
                LPAREN parameters RPAREN
                {
                    $$ = $2
                }
        |       LPAREN RPAREN
                {
                    $$ = nil
                }
        ;

functionDeclaration:
                FUNC IDENT functionParameters functionParameters
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				inFn = true
				fn := MakeFunction($2)
				mod.AddFunction(fn)
				if fn, err := mod.GetCurrentFunction(); err == nil {

					//checking if there are duplicate parameters
					dups := append($3, $4...)
					for _, param := range dups {
						for _, dup := range dups {
							if param.Name == dup.Name && param != dup {
								panic(fmt.Sprintf("%d: duplicate input and/or output parameters in function '%s'", yyS[yypt-0].line+1, $2))
							}
						}
					}
					
					for _, inp := range $3 {
						fn.AddInput(inp)
					}
					for _, out := range $4 {
						fn.AddOutput(out)
					}
				}
			}
                }
                functionStatements
        ;

parameter:
                IDENT typeSpecifier
                {
			$$ = MakeParameter($1, $2)
                }
        |       IDENT IDENT
                {
			$$ = MakeParameter($1, $2)
                }
        ;

parameters:
                parameter
                {
			var params []*CXParameter
                        params = append(params, $1)
                        $$ = params
                }
        |       parameters COMMA parameter
                {
			$1 = append($1, $3)
                        $$ = $1
                }
        ;

functionStatements:
                LBRACE expressionsAndStatements RBRACE
                {
			inFn = false
                }
        |       LBRACE RBRACE
                {
			inFn = false
                }
        ;

expressionsAndStatements:
                nonAssignExpression
        |       assignExpression
        |       statement
        |       selector
        |       stepping
        |       debugging
        |       affordance
        |       remover
        |       expressionsAndStatements nonAssignExpression
        |       expressionsAndStatements assignExpression
        |       expressionsAndStatements statement
        |       expressionsAndStatements selector
        |       expressionsAndStatements stepping
        |       expressionsAndStatements debugging
        |       expressionsAndStatements affordance
        |       expressionsAndStatements remover
        ;


assignExpression:
                VAR IDENT typeSpecifier definitionAssignment
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					if $4 == nil {
						if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
							expr := MakeExpression(op)
							if !replMode {
								expr.FileLine = yyS[yypt-0].line + 1
							}

							fn.AddExpression(expr)
							expr.AddOutputName($2)
							
							typ := encoder.Serialize($3)
							arg := MakeArgument(&typ, "str")
							expr.AddArgument(arg)

							if strct, err := cxt.GetStruct($3, mod.Name); err == nil {
								for _, fld := range strct.Fields {
									expr := MakeExpression(op)
									if !replMode {
										expr.FileLine = yyS[yypt-0].line + 1
									}
									fn.AddExpression(expr)
									expr.AddOutputName(fmt.Sprintf("%s.%s", $2, fld.Name))
									typ := []byte(fld.Typ)
									arg := MakeArgument(&typ, "str")
									expr.AddArgument(arg)
								}
							}
						}
					} else {
						switch $3 {
						case "bool":
							var ds int32
							encoder.DeserializeRaw(*$4.Value, &ds)
							new := encoder.SerializeAtomic(ds)
							val := MakeArgument(&new, "bool")
							
							if op, err := cxt.GetFunction("bool.id", mod.Name); err == nil {
								expr := MakeExpression(op)
								if !replMode {
									expr.FileLine = yyS[yypt-0].line + 1
								}
								fn.AddExpression(expr)
								expr.AddOutputName($2)
								expr.AddArgument(val)
							}
						case "byte":
							var ds int32
							encoder.DeserializeRaw(*$4.Value, &ds)
							new := []byte{byte(ds)}
							val := MakeArgument(&new, "byte")
							
							if op, err := cxt.GetFunction("byte.id", mod.Name); err == nil {
								expr := MakeExpression(op)
								if !replMode {
									expr.FileLine = yyS[yypt-0].line + 1
								}
								fn.AddExpression(expr)
								expr.AddOutputName($2)
								expr.AddArgument(val)
							}
						case "i64":
							var ds int32
							encoder.DeserializeRaw(*$4.Value, &ds)
							new := encoder.Serialize(int64(ds))
							val := MakeArgument(&new, "i64")

							if op, err := cxt.GetFunction("i64.id", mod.Name); err == nil {
								expr := MakeExpression(op)
								if !replMode {
									expr.FileLine = yyS[yypt-0].line + 1
								}
								fn.AddExpression(expr)
								expr.AddOutputName($2)
								expr.AddArgument(val)
							}
						case "f64":
							var ds float32
							encoder.DeserializeRaw(*$4.Value, &ds)
							new := encoder.Serialize(float64(ds))
							val := MakeArgument(&new, "f64")

							if op, err := cxt.GetFunction("f64.id", mod.Name); err == nil {
								expr := MakeExpression(op)
								if !replMode {
									expr.FileLine = yyS[yypt-0].line + 1
								}
								fn.AddExpression(expr)
								expr.AddOutputName($2)
								expr.AddArgument(val)
							}
						default:
							val := $4
							var getFn string
							switch $3 {
							case "i32": getFn = "i32.id"
							case "f32": getFn = "f32.id"
							case "[]bool": getFn = "[]bool.id"
							case "[]byte": getFn = "[]byte.id"
							case "[]str": getFn = "[]str.id"
							case "[]i32": getFn = "[]i32.id"
							case "[]i64": getFn = "[]i64.id"
							case "[]f32": getFn = "[]f32.id"
							case "[]f64": getFn = "[]f64.id"
							}

							if op, err := cxt.GetFunction(getFn, mod.Name); err == nil {
								expr := MakeExpression(op)
								if !replMode {
									expr.FileLine = yyS[yypt-0].line + 1
								}
								fn.AddExpression(expr)
								expr.AddOutputName($2)
								expr.AddArgument(val)
							}
						}
					}
				}
			}
                }
        |       VAR IDENT LBRACK RBRACK IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
						expr := MakeExpression(op)

						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}
						fn.AddExpression(expr)
						expr.AddOutputName($2)
						typ := encoder.Serialize(fmt.Sprintf("[]%s", $5))
						arg := MakeArgument(&typ, "str")
						expr.AddArgument(arg)
					}
				}
			}
                }
        |       argumentsList assignOperator argumentsList
                {
			argsL := $1
			argsR := $3

			if len(argsL) > len(argsR) {
				panic(fmt.Sprintf("%d: trying to assign values to variables using a function with no output parameters", yyS[yypt-0].line + 1))
			}

			if fn, err := cxt.GetCurrentFunction(); err == nil {
				for i, argL := range argsL {
					if argsR[i] == nil {
						continue
					}
					// argL is going to be the output name
					typeParts := strings.Split(argsR[i].Typ, ".")

					var typ string
					var secondTyp string
					var idFn string

					if len(typeParts) > 1 {
						typ = "str"
						secondTyp = typeParts[1] // i32, f32, etc
					} else if typeParts[0] == "ident" {
						typ = "str"
						secondTyp = "ident"
					} else {
						typ = typeParts[0] // i32, f32, etc
					}

					if $2 == ":=" || $2 == "=" {
						if secondTyp == "" {
							idFn = MakeIdentityOpName(typ)
						} else {
							idFn = "identity"
						}

						if op, err := cxt.GetFunction(idFn, CORE_MODULE); err == nil {
							expr := MakeExpression(op)
							if !replMode {
								expr.FileLine = yyS[yypt-0].line + 1
							}

							fn.AddExpression(expr)
							expr.AddTag(tag)
							tag = ""

							var outName string
							encoder.DeserializeRaw(*argL.Value, &outName)
							
							expr.AddOutputName(outName)

							arg := MakeArgument(argsR[i].Value, typ)
							expr.AddArgument(arg)
						}
					} else {
						// +=, -=, *=, etc.
						var opName string
						var typName string

						if secondTyp == "ident" {
							var identName string
							encoder.DeserializeRaw(*argsR[i].Value, &identName)

							if argTyp, err := GetIdentType(identName, cxt); err == nil {
								typName = argTyp
							} else {
								panic(err)
							}
						} else if secondTyp == "" {
							typName = typ
						} else {
							typName = secondTyp
						}

						// here here

						switch $2 {
						case "+=":
							opName = "add"
						case "-=":
							opName = "sub"
						case "*=":
							opName = "mul"
						case "/=":
							opName = "div"
						case "%=":
							opName = "mod"
						case "**=":
							opName = "pow"
						case "<<=":
							opName = "bitshl"
						case ">>=":
							opName = "bitshr"
						case "&=":
							opName = "bitand"
						case "^=":
							opName = "bitxor"
						case "|=":
							opName = "bitor"
						}

						if op, err := cxt.GetFunction(fmt.Sprintf("%s.%s", typName, opName), CORE_MODULE); err == nil {
							expr := MakeExpression(op)
							if !replMode {
								expr.FileLine = yyS[yypt-0].line + 1
							}

							fn.AddExpression(expr)
							expr.AddTag(tag)
							tag = ""

							var outName string
							encoder.DeserializeRaw(*argL.Value, &outName)
							
							expr.AddOutputName(outName)
							expr.AddArgument(argL)
							
							if len(typeParts) > 1 {
								expr.AddArgument(MakeArgument(argsR[i].Value, "ident"))
							} else {
								expr.AddArgument(MakeArgument(argsR[i].Value, typeParts[0]))
							}
						}
					}
				}
			}
                }
        ;

nonAssignExpression:
                IDENT arguments
                {
			var modName string
			var fnName string
			var err error
			identParts := strings.Split($1, ".")
			
			if len(identParts) == 2 {
				modName = identParts[0]
				fnName = identParts[1]
			} else {
				fnName = identParts[0]
				mod, e := cxt.GetCurrentModule()
				modName = mod.Name
				err = e
			}

			found := false
			currModName := ""
			if mod, err := cxt.GetCurrentModule(); err == nil {
				currModName = mod.Name
				for _, imp := range mod.Imports {
					if modName == imp.Name {
						found = true
						break
					}
				}
			}

			isModule := false
			if _, err := cxt.GetModule(modName); err == nil {
				isModule = true
			}
			
			if !found && !IsNative(modName + "." + fnName) && modName != currModName && isModule {
				fmt.Printf("%d: module '%s' was not imported or does not exist\n", yyS[yypt-0].line + 1, modName)
			} else {
				if err == nil {
					if fn, err := cxt.GetCurrentFunction(); err == nil {
						if op, err := cxt.GetFunction(fnName, modName); err == nil {
							expr := MakeExpression(op)
							if !replMode {
								expr.FileLine = yyS[yypt-0].line + 1
							}
							fn.AddExpression(expr)
							expr.AddTag(tag)
							tag = ""
							for _, arg := range $2 {
								typeParts := strings.Split(arg.Typ, ".")

								arg.Typ = typeParts[0]
								expr.AddArgument(arg)
							}

							lenOut := len(op.Outputs)
							outNames := make([]string, lenOut)
							args := make([]*CXArgument, lenOut)
							
							for i, out := range op.Outputs {
								outNames[i] = MakeGenSym(NON_ASSIGN_PREFIX)
								byteName := encoder.Serialize(outNames[i])
								args[i] = MakeArgument(&byteName, fmt.Sprintf("ident.%s", out.Typ))

								expr.AddOutputName(outNames[i])
							}
							
							$$ = args
						} else {
							fmt.Printf("Function '%s' not defined\n", $1)
						}
					}
				}
			}
			
                }
        |       argument PLUSPLUS
                {
			$$ = []*CXArgument{unaryOp($2, $1, yyS[yypt-0].line + 1)}
                }
        |       argument MINUSMINUS
                {
			$$ = []*CXArgument{unaryOp($2, $1, yyS[yypt-0].line + 1)}
                }
        ;

beginFor:       FOR
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                ;

conditionControl:
                nonAssignExpression
                {
			$$ = $1
                }
        |       argument
                {
			$$ = []*CXArgument{$1}
                }
        ;

statement:      RETURN
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}
						fn.AddExpression(expr)
						val := MakeDefaultValue("bool")
						expr.AddArgument(MakeArgument(val, "bool"))
						lines := encoder.SerializeAtomic(int32(-len(fn.Expressions)))
						expr.AddArgument(MakeArgument(&lines, "i32"))
						expr.AddArgument(MakeArgument(&lines, "i32"))
					}
				}
			}
                }
        |       GOTO IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					// this one is goTo, not baseGoTo
					if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}
						fn.AddExpression(expr)

						//label := []byte($2)
						label := encoder.Serialize($2)
						expr.AddArgument(MakeArgument(&label, "str"))
					}
				}
			}
                }
        |       IF conditionControl
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}
						fn.AddExpression(expr)
					}
				}
			}
                }
                LBRACE
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                expressionsAndStatements RBRACE elseStatement
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
                                        goToExpr := fn.Expressions[$<i>5 - 1]

					var elseLines []byte
					if $<i>8 > 0 {
						elseLines = encoder.Serialize(int32(len(fn.Expressions) - $<i>5 - $<i>8 + 1))
					} else {
						elseLines = encoder.Serialize(int32(len(fn.Expressions) - $<i>5 + 1))
					}
					
					thenLines := encoder.Serialize(int32(1))

					var typ string
					if len($2[0].Typ) > len("ident.") && $2[0].Typ[:len("ident.")] == "ident." {
						typ = "ident"
					} else {
						typ = $2[0].Typ
					}

					//goToExpr.AddArgument(MakeArgument(predVal, "ident"))
					goToExpr.AddArgument(MakeArgument($2[0].Value, typ))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
				}
			}
                }
        |       beginFor
                nonAssignExpression
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}
						fn.AddExpression(expr)
					}
				}
			}
                }
                LBRACE expressionsAndStatements RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					goToExpr := fn.Expressions[$<i>3]
					elseLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>3 + 1))
					thenLines := encoder.Serialize(int32(1))

					//if multiple value return, take first one for condition
					predVal := $2[0].Value
					
					goToExpr.AddArgument(MakeArgument(predVal, "ident"))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
					
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						goToExpr := MakeExpression(goToFn)
						if !replMode {
							goToExpr.FileLine = lineNo
						}
						fn.AddExpression(goToExpr)

						elseLines := encoder.Serialize(int32(0))
						thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>1 + 1))

						alwaysTrue := encoder.Serialize(int32(1))

						goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
						goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
						goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
					}
					
				}
			}
                }
        |       beginFor IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}
						fn.AddExpression(expr)
					}
				}
			}
                }
                LBRACE expressionsAndStatements RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					goToExpr := fn.Expressions[$<i>1]

					elseLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>1 + 1))
					thenLines := encoder.Serialize(int32(1))
					
					predVal := encoder.Serialize($2)
					
					goToExpr.AddArgument(MakeArgument(&predVal, "ident"))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
					
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						goToExpr := MakeExpression(goToFn)
						if !replMode {
							goToExpr.FileLine = lineNo
						}
						fn.AddExpression(goToExpr)

						elseLines := encoder.Serialize(int32(0))
						thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>1 + 1))

						alwaysTrue := encoder.Serialize(int32(1))

						goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
						goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
						goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
					}
				}
			}
                }
        |       beginFor BOOLEAN
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}
						fn.AddExpression(expr)
					}
				}
			}
                }
                LBRACE
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                expressionsAndStatements RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					goToExpr := fn.Expressions[$<i>5 - 1]
					
					elseLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>5 + 2))
					thenLines := encoder.Serialize(int32(1))
					
					var predVal []byte
					if $2 == int32(1) {
						predVal = encoder.Serialize(int32(1))
					} else {
						predVal = encoder.Serialize(int32(0))
					}

					goToExpr.AddArgument(MakeArgument(&predVal, "bool"))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
					
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						goToExpr := MakeExpression(goToFn)
						if !replMode {
							goToExpr.FileLine = lineNo
						}
						fn.AddExpression(goToExpr)

						elseLines := encoder.Serialize(int32(0))
						thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>5))

						alwaysTrue := encoder.Serialize(int32(1))

						goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
						goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
						goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
					}
					
				}
			}
                }
        |       beginFor // $<i>1
                forLoopAssignExpression // $2
                {//$<i>3
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
//                              ';' nonAssignExpression
                ';' conditionControl
                {//$<i>6
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                {//$<i>7
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				if mod, err := cxt.GetCurrentModule(); err == nil {
					if fn, err := mod.GetCurrentFunction(); err == nil {
						if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
							expr := MakeExpression(goToFn)
							if !replMode {
								expr.FileLine = yyS[yypt-0].line + 1
							}
							fn.AddExpression(expr)
						}
					}
				} else {
					fmt.Println(err)
				}

				$<i>$ = len(fn.Expressions)
			}
                }
                ';' forLoopAssignExpression //$<bool>9
                {//$<int>10
			//increment goTo
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				goToExpr := fn.Expressions[$<i>7 - 1]
				if $<bool>9 {
					if mod, err := cxt.GetCurrentModule(); err == nil {
						if fn, err := mod.GetCurrentFunction(); err == nil {
							if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
								expr := MakeExpression(goToFn)
								if !replMode {
									expr.FileLine = yyS[yypt-0].line + 1
								}
								fn.AddExpression(expr)
							}
						}
					}

					thenLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>7 + 1))
					// elseLines := encoder.Serialize(int32(0)) // this is added later in $12

					predVal := $5[0].Value
					
					goToExpr.AddArgument(MakeArgument(predVal, "ident"))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					// goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
				}
				$<i>$ = len(fn.Expressions)
			}
                }
                LBRACE expressionsAndStatements RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					goToExpr := fn.Expressions[$<i>10 - 1]

					if $<bool>9 {
						predVal := $5[0].Value

						thenLines := encoder.Serialize(int32(-($<i>10 - $<i>3) + 1))
						elseLines := encoder.Serialize(int32(0))

						goToExpr.AddArgument(MakeArgument(predVal, "bool"))
						goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
						goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))

						if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
							goToExpr := MakeExpression(goToFn)
							if !replMode {
								goToExpr.FileLine = lineNo
							}
							fn.AddExpression(goToExpr)

							alwaysTrue := encoder.Serialize(int32(1))

							thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>7) + 1)
							elseLines := encoder.Serialize(int32(0))

							goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
							goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
							goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))

							condGoToExpr := fn.Expressions[$<i>7 - 1]

							condThenLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>7 + 1))
							
							condGoToExpr.AddArgument(MakeArgument(&condThenLines, "i32"))
						}
					} else {
						predVal := $5[0].Value

						thenLines := encoder.Serialize(int32(1))
						elseLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>7 + 2))
						
						goToExpr.AddArgument(MakeArgument(predVal, "ident"))
						goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
						goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))

						if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
							goToExpr := MakeExpression(goToFn)
							if !replMode {
								goToExpr.FileLine = lineNo
							}
							fn.AddExpression(goToExpr)
							
							alwaysTrue := encoder.Serialize(int32(1))

							thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>3 + 1))
							elseLines := encoder.Serialize(int32(0))

							goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
							goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
							goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
						}
					}
				}
			}
                }
        |       VAR IDENT IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}

						fn.AddExpression(expr)
						expr.AddOutputName($2)
						
						typ := encoder.Serialize($3)
						arg := MakeArgument(&typ, "str")
						expr.AddArgument(arg)
					}
				}
			}
                }
        |       ';'
        ;

forLoopAssignExpression:
                {
			$<bool>$ = false
                }
        |       assignExpression
                {
			$<bool>$ = true
                }
        |       nonAssignExpression
                {
			$<bool>$ = true
                }
        ;

elseStatement:
                /* empty */
                {
                    $<i>$ = 0
                }
        |       ELSE
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := mod.GetCurrentFunction(); err == nil {
                            if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
				    expr := MakeExpression(goToFn)
				    if !replMode {
					    expr.FileLine = yyS[yypt-0].line + 1
				    }
				    fn.AddExpression(expr)
                            }
                        }
                    }
                }
                LBRACE
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                expressionsAndStatements RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					goToExpr := fn.Expressions[$<i>4 - 1]
					
					elseLines := encoder.Serialize(int32(0))
					thenLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>4 + 1))

					alwaysTrue := encoder.Serialize(int32(1))

					goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))

					$<i>$ = len(fn.Expressions) - $<i>4
				}
			}
                }
                ;

expressions:
                nonAssignExpression
        |       assignExpression
        |       expressions nonAssignExpression
        |       expressions assignExpression
        ;











inferPred:      inferObj
                {
			$$ = $1
                }
        |       inferCond
                {
			$$ = $1
                }
        |       inferPred COMMA inferObj
                {
			$1 = append($1, $3...)
			$$ = $1
                }
        |       inferPred COMMA inferCond
                {
			$1 = append($1, $3...)
			$$ = $1
                }
        ;

inferCond:      IDENT LPAREN inferPred RPAREN
                {
			$$ = append($3, $1)
                }
        |       BOOLEAN
                {
			var obj string
			if $1 == 1 {
				obj = "true"
			} else {
				obj = "false"
			}
			$$ = []string{obj, "1.0", "weight"}
                }
        ;

// inferArg:       IDENT
//                 {
// 			$$ = fmt.Sprintf("%s", $1)
//                 }
//         |       STRING
//                 {
// 			str := strings.TrimPrefix($1, "\"")
//                         str = strings.TrimSuffix(str, "\"")
// 			$$ = str
//                 }
//         |       FLOAT
//                 {
// 			$$ = fmt.Sprintf("%f", $1)
//                 }
//         |       INT
//                 {
// 			$$ = fmt.Sprintf("%d", $1)
//                 }
//         |       nonAssignExpression
//                 {
// 			var ident string
// 			encoder.DeserializeRaw(*$1[0].Value, &ident)
// 			$$ = ident
//                 }
//         ;

// inferOp:        EQUAL
//         |       GTHAN
//         |       LTHAN
// ;

relationalOp:   EQUAL
        |       GTHAN
        |       LTHAN
        |       UNEQUAL
                ;

inferActionArg:
                inferObj
                {
			$$ = $1
                }
        |       IDENT
                {
                    $$ = []string{$1}
                }
        |       AFFVAR relationalOp argument
                {
			argStr := ArgToString($3)
			$$ = []string{$1, argStr, $2}
                }
        |       MULT relationalOp argument
                {
			argStr := ArgToString($3)
			$$ = []string{$1, argStr, $2}
                }
        |       MULT relationalOp nonAssignExpression
                {
			var identName string
			encoder.DeserializeRaw(*$3[0].Value, &identName)
			$$ = []string{$1, identName, $2}
                }
        |       MULT relationalOp MULT
                {
			$$ = []string{$1, $1, $2}
                }
        ;

inferAction:
		IDENT LPAREN inferActionArg RPAREN
                {
			$$ = append($3, $1)
                }
        ;

inferActions:
                inferAction
                {
			$$ = $1
                }
        |       inferActions inferAction
                {
			$1 = append($1, $2...)
			$$ = $1
                }
                ;

inferRule:      IF inferCond LBRACE inferActions RBRACE
                {
			if $4 == nil {
				$$ = $2
			} else {
				block := append($2, "if")
				block = append(block, $4...)
				block = append(block, "endif")
				$$ = block
			}
                }
        |       IF inferObj LBRACE inferActions RBRACE
                {
			if $4 == nil {
				$$ = $2
			} else {
				block := append($2, "single")
				block = append(block, "if")
				block = append(block, $4...)
				block = append(block, "endif")
				$$ = block
			}
                }
        ;

inferRules:     inferRule
                {
			$$ = $1
                }
        |       inferRules inferRule
                {
			$1 = append($1, $2...)
			$$ = $1
                }
        ;

inferWeight:    FLOAT
                {
			$$ = fmt.Sprintf("%f", $1)
                }
        |       INT
                {
			$$ = fmt.Sprintf("%d", $1)
                }
        |       IDENT
                {
			$$ = fmt.Sprintf("%s", $1)
                }
        ;

inferObj:
                {
			$$ = []string{}
                }
        |       IDENT WEIGHT inferWeight
                {
			$$ = []string{$1, $3, "weight"}
                }
        |       IDENT WEIGHT nonAssignExpression
                {
			var ident string
			encoder.DeserializeRaw(*$3[0].Value, &ident)
			$$ = []string{$1, ident, "weight"}
                }
;

inferObjs:      inferObj
                {
			$$ = $1
                }
        |       inferObjs COMMA inferObj
                {
			$1 = append($1, $3...)
			$$ = $1
                }
        ;

inferTarget:    IDENT LPAREN IDENT RPAREN
                {
			$$ = []string{$3, $1}
                }
        ;

inferTargets:   inferTarget
                {
			$$ = $1
                }
        |       inferTargets inferTarget
                {
			$1 = append($1, $2...)
			$$ = $1
                }
        ;

inferClauses:   inferObjs
        |       inferRules
        |       inferTargets
        ;




structLitDef:
                TAG argument
                {
			name := strings.TrimSuffix($1, ":")
			$$ = MakeDefinition(name, $2.Value, $2.Typ)
                }
        |       TAG nonAssignExpression
                {
			name := strings.TrimSuffix($1, ":")
			$$ = MakeDefinition(name, $2[0].Value, $2[0].Typ)
                }
                    
;

structLitDefs:  structLitDef
                {
			var defs []*CXDefinition
			$$ = append(defs, $1)
                }
        |       structLitDefs COMMA structLitDef
                {
			$1 = append($1, $3)
			$$ = $1
                }
        ;

structLiteral:
                {
			$$ = nil
                }
        |       
		LBRACE structLitDefs RBRACE
                {
			$$ = $2
                }
        ;

argument:
                argument PLUS argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument MINUS argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument MULT argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument DIV argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument REMAINDER argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument LEFTSHIFT argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument RIGHTSHIFT argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument EXP argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument EQUAL argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument UNEQUAL argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument GTHAN argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument GTHANEQ argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument LTHAN argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument LTHANEQ argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument OR argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       argument AND argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       NOT argument
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				if op, err := cxt.GetFunction("not", CORE_MODULE); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = yyS[yypt-0].line + 1
					}
					fn.AddExpression(expr)
					expr.AddTag(tag)
					tag = ""
					expr.AddArgument($2)

					outName := MakeGenSym(NON_ASSIGN_PREFIX)
					byteName := encoder.Serialize(outName)
					
					expr.AddOutputName(outName)
					$$ = MakeArgument(&byteName, "ident")
				}
			}
		}
        |       LPAREN argument RPAREN
                {
			$$ = $2
                }
        |       BYTENUM
                {
			val := encoder.Serialize($1)
                        $$ = MakeArgument(&val, "byte")
                }
        |       INT
                {
			val := encoder.Serialize($1)
                        $$ = MakeArgument(&val, "i32")
                }
        |       LONG
                {
			val := encoder.Serialize($1)
                        $$ = MakeArgument(&val, "i64")
                }
        |       FLOAT
                {
			val := encoder.Serialize($1)
			$$ = MakeArgument(&val, "f32")
                }
        |       DOUBLE
                {
			val := encoder.Serialize($1)
			$$ = MakeArgument(&val, "f64")
                }
        |       BOOLEAN
                {
			val := encoder.Serialize($1)
			$$ = MakeArgument(&val, "bool")
                }
        |       STRING
                {
			str := strings.TrimPrefix($1, "\"")
                        str = strings.TrimSuffix(str, "\"")

			val := encoder.Serialize(str)
			
                        $$ = MakeArgument(&val, "str")
                }
        |       IDENT
                {
			val := encoder.Serialize($1)
			$$ = MakeArgument(&val, "ident")
                }
        |       typeSpecifier LBRACE structLitDefs RBRACE
                {
			val := encoder.Serialize($2)
		
			if len($4) < 1 {
				$$ = MakeArgument(&val, "ident")
			} else {
				// then it's a struct literal
				if mod, err := cxt.GetCurrentModule(); err == nil {
					if fn, err := cxt.GetCurrentFunction(); err == nil {
						if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
							expr := MakeExpression(op)
							if !replMode {
								expr.FileLine = yyS[yypt-0].line + 1
							}
							fn.AddExpression(expr)

							outName := MakeGenSym(NON_ASSIGN_PREFIX)
							sOutName := encoder.Serialize(outName)

							expr.AddOutputName(outName)
							typ := encoder.Serialize($1)
							expr.AddArgument(MakeArgument(&typ, "str"))

							$$ = MakeArgument(&sOutName, fmt.Sprintf("ident.%s", $1))
							for _, def := range $3 {
								typeParts := strings.Split(def.Typ, ".")

								var typ string
								var secondTyp string
								var idFn string

								if len(typeParts) > 1 {
									typ = "str"
									secondTyp = typeParts[1] // i32, f32, etc
								} else if typeParts[0] == "ident" {
									typ = "str"
									secondTyp = "ident"
								} else {
									typ = typeParts[0] // i32, f32, etc
								}

								if secondTyp == "" {
									idFn = MakeIdentityOpName(typ)
								} else {
									idFn = "identity"
								}
								
								if op, err := cxt.GetFunction(idFn, CORE_MODULE); err == nil {
									expr := MakeExpression(op)
									if !replMode {
										expr.FileLine = yyS[yypt-0].line + 1
									}
									fn.AddExpression(expr)
									expr.AddTag(tag)
									tag = ""

									outName := fmt.Sprintf("%s.%s", outName, def.Name)
									expr.AddOutputName(outName)
									arg := MakeArgument(def.Value, typ)
									expr.AddArgument(arg)
								}
							}
						}
					}
				}
			}
                }
        |       IDENT LBRACK INT RBRACK
                {
			val := encoder.Serialize(fmt.Sprintf("%s[%d", $1, $3))
			$$ = MakeArgument(&val, "ident")
                }
        |       INFER LBRACE inferClauses RBRACE
                {
			val := encoder.Serialize($3)
			$$ = MakeArgument(&val, "[]str")
                }
        |       typeSpecifier LBRACE argumentsList RBRACE
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil && inFn {
					if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}

						outName := MakeGenSym(NON_ASSIGN_PREFIX)
						sOutName := encoder.Serialize(outName)

						fn.AddExpression(expr)
						expr.AddOutputName(outName)

						typ := encoder.Serialize($1)
						arg := MakeArgument(&typ, "str")
						expr.AddArgument(arg)
						
						if op, err := cxt.GetFunction(fmt.Sprintf("%s.append", $1), mod.Name); err == nil {
							for _, arg := range $3 {
								typeParts := strings.Split(arg.Typ, ".")
								arg.Typ = typeParts[0]
								expr := MakeExpression(op)
								fn.AddExpression(expr)
								expr.AddOutputName(outName)
								expr.AddArgument(MakeArgument(&sOutName, "ident"))
								expr.AddArgument(CastArgumentForArray($1, arg))
							}
						}
						$$ = MakeArgument(&sOutName, "ident")
					}
				} else {
					// then it's for a global definition
					switch $1 {
					case "[]str":
						vals := make([]string, len($3))
						for i, arg := range $3 {
							var val string
							encoder.DeserializeRaw(*arg.Value, &val)
							vals[i] = val
						}
						sVal := encoder.Serialize(vals)
						$$ = MakeArgument(&sVal, "[]str")
                                        case "[]bool":
						vals := make([]int32, len($3))
						for i, arg := range $3 {
							var val int32
							encoder.DeserializeRaw(*arg.Value, &val)
							vals[i] = val
						}
						sVal := encoder.Serialize(vals)
						$$ = MakeArgument(&sVal, "[]bool")
                                        case "[]byte":
						vals := make([]byte, len($3))
						for i, arg := range $3 {
							var val int32
							encoder.DeserializeRaw(*arg.Value, &val)
							vals[i] = byte(val)
						}
						sVal := encoder.Serialize(vals)
						$$ = MakeArgument(&sVal, "[]byte")
                                        case "[]i32":
						vals := make([]int32, len($3))
						for i, arg := range $3 {
							var val int32
							encoder.DeserializeRaw(*arg.Value, &val)
							vals[i] = val
						}
						sVal := encoder.Serialize(vals)
						$$ = MakeArgument(&sVal, "[]i32")
                                        case "[]i64":
						vals := make([]int64, len($3))
						for i, arg := range $3 {
							var val int32
							encoder.DeserializeRaw(*arg.Value, &val)
							vals[i] = int64(val)
						}
						sVal := encoder.Serialize(vals)
						$$ = MakeArgument(&sVal, "[]i64")
                                        case "[]f32":
						vals := make([]float32, len($3))
						for i, arg := range $3 {
							var val float32
							encoder.DeserializeRaw(*arg.Value, &val)
							vals[i] = val
						}
						sVal := encoder.Serialize(vals)
						$$ = MakeArgument(&sVal, "[]f32")
                                        case "[]f64":
						vals := make([]float64, len($3))
						for i, arg := range $3 {
							var val float32
							encoder.DeserializeRaw(*arg.Value, &val)
							vals[i] = float64(val)
						}
						sVal := encoder.Serialize(vals)
						$$ = MakeArgument(&sVal, "[]f64")
					}
                                }
			}
                }
                // empty arrays
        |       typeSpecifier LBRACE RBRACE
                {
			switch $1 {
			case "[]bool":
                                vals := make([]int32, 0)
				sVal := encoder.Serialize(vals)

				$$ = MakeArgument(&sVal, "[]bool")
			case "[]byte":
                                vals := make([]byte, 0)
				sVal := encoder.Serialize(vals)
				$$ = MakeArgument(&sVal, "[]byte")
			case "[]str":
                                vals := make([]string, 0)
				sVal := encoder.Serialize(vals)
				$$ = MakeArgument(&sVal, "[]str")
			case "[]i32":
				vals := make([]int32, 0)
				sVal := encoder.Serialize(vals)
				$$ = MakeArgument(&sVal, "[]i32")
			case "[]i64":
                                vals := make([]int64, 0)
				sVal := encoder.Serialize(vals)
				$$ = MakeArgument(&sVal, "[]i64")
			case "[]f32":
                                vals := make([]float32, 0)
				sVal := encoder.Serialize(vals)
				$$ = MakeArgument(&sVal, "[]f32")
			case "[]f64":
                                vals := make([]float64, 0)
				sVal := encoder.Serialize(vals)
				$$ = MakeArgument(&sVal, "[]f64")
			}
                }
        ;

arguments:
                LPAREN argumentsList RPAREN
                {
                    $$ = $2
                }
        |       LPAREN RPAREN
                {
                    $$ = nil
                }
        ;

argumentsList:
                argument
                {
			var args []*CXArgument
			args = append(args, $1)
			$$ = args
                }
        |       nonAssignExpression
                {
			args := $1
			$$ = args
                }
        |       argumentsList COMMA argument
                {
			$1 = append($1, $3)
			$$ = $1
                }
        |       argumentsList COMMA nonAssignExpression
                {
			args := $3

			$1 = append($1, args...)
			$$ = $1
                }
        ;

%%

