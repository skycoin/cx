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


		"github.com/mndrix/golog"
		"github.com/mndrix/golog/read"
		"github.com/mndrix/golog/term"
	)

	var program bytes.Buffer
	
	var cxt = MakeContext()
	var m = golog.NewInteractiveMachine()
	
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
	//var dProgram bool = false
	var tag string = ""
	var asmNL = "\n"

	func warnf(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format, args...)
		os.Stderr.Sync()
	}
%}

%union {
	i int
	i32 int32
	f32 float32
	i64 int64
	f64 float64
	tok string
	bool bool
	string string

	line int

	parameter *CXParameter
	parameters []*CXParameter

	argument *CXArgument
	arguments []*CXArgument

	expression *CXExpression
	expressions []*CXExpression

	field *CXField
	fields []*CXField

	name string
	names []string
}

%token  <i32>           INT BOOLEAN
%token  <f32>           FLOAT
%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK IDENT
                        VAR COMMA COMMENT STRING PACKAGE IF ELSE FOR TYPSTRUCT STRUCT
                        ASSIGN CASSIGN GTHAN LTHAN LTEQ GTEQ IMPORT RETURN GOTO
                        /* Types */
                        BOOL STR I32 I64 F32 F64 BYTE BOOLA BYTEA I32A I64A F32A F64A STRA
                        /* Selectors */
                        SPACKAGE SSTRUCT SFUNC
                        /* Removers */
                        REM DEF EXPR FIELD INPUT OUTPUT CLAUSES OBJECT OBJECTS
                        /* Stepping */
                        STEP PSTEP TSTEP
                        /* Debugging */
                        DSTACK DPROGRAM DQUERY DSTATE
                        /* Affordances */
                        AFF TAG
                        /* Prolog */
                        CCLAUSES QUERY COBJECT COBJECTS

%type   <tok>           typeSpecifier
                        
%type   <parameter>     parameter
%type   <parameters>    parameters functionParameters
%type   <argument>      argument definitionAssignment
%type   <arguments>     arguments argumentsList nonAssignExpression
%type   <fields>        fields structFields
%type   <expression>    assignExpression
%type   <expressions>   elseStatement
//%type   <names>         nonAssignExpression
%type   <bool>          selectorLines selectorExpressionsAndStatements selectorFields

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
        |       prolog
        ;

prolog:
                COBJECTS
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				for i, object := range mod.Objects {
					fmt.Printf("%d.- %s\n", i, object.Name)
				}
			}
                }
        |       CCLAUSES STRING
                {
			clauses := strings.TrimPrefix($2, "\"")
			clauses = strings.TrimSuffix(clauses, "\"")

			b := bytes.NewBufferString(clauses)
			m = m.Consult(b)
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.AddClauses(clauses)
			}
                }
        |       COBJECT IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.AddObject(MakeObject($2));
			}
                }
        |       QUERY STRING
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				query := strings.TrimPrefix($2, "\"")
				query = strings.TrimSuffix(query, "\"")
				mod.AddQuery(query)
			}
                }
        |       DQUERY STRING
                {
			query := strings.TrimPrefix($2, "\"")
			query = strings.TrimSuffix(query, "\"")

			goal, err := read.Term(query)
			if err == nil {
				variables := term.Variables(goal)
				answers := m.ProveAll(goal)

				yesNoAnswer := false
				if len(answers) == 0 {
					fmt.Println("no.")
					yesNoAnswer = true
				} else if variables.Size() == 0 {
					fmt.Println("yes.")
					yesNoAnswer = true
				}

				if !yesNoAnswer {
					for i, answer := range answers {
						lines := make([]string, 0)
						variables.ForEach(func(name string, variable interface{}) {
							v := variable.(*term.Variable)
							val := answer.Resolve_(v)
							line := fmt.Sprintf("%s = %s", name, val)
							lines = append(lines, line)
						})

						warnf("%s", strings.Join(lines, "\n"))
						if i == len(answers)-1 {
							fmt.Printf("\t.\n\n")
						} else {
							warnf("\t;\n")
						}
					}
				}
			} else {
				fmt.Println("Problem parsing the query.")
			}
                }
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
							PrintAffordances(expr.GetAffordances())
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
							affs := expr.GetAffordances()
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
							affs := expr.GetAffordances()
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
							affs := expr.GetAffordances()
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
					fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Value, def.Typ))
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
        |       REM CLAUSES
                {
			m = golog.NewInteractiveMachine()
                }
        |       REM OBJECT IDENT
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.RemoveObject($3)
			}
                }
        |       REM OBJECTS
                {
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.RemoveObjects()
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
        ;

typeSpecifier:
                I32 {$$ = $1}
        |       I64 {$$ = $1}
        |       F32 {$$ = $1}
        |       F64 {$$ = $1}
        |       BOOL {$$ = $1}
        |       BYTE {$$ = $1}
        |       BOOLA {$$ = $1}
        |       BYTEA {$$ = $1}
        |       I32A {$$ = $1}
        |       I64A {$$ = $1}
        |       F32A {$$ = $1}
        |       F64A {$$ = $1}
        |       STR {$$ = $1}
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
				}
				// if strct, err := cxt.GetStruct($3, mod.Name); err == nil {
				// 	if flds, err := strct.GetFields(); err == nil {
				// 		for _, fld := range flds {
				// 			zeroVal := MakeDefaultValue(fld.Typ)
				// 			defName := fmt.Sprintf("%s.%s", $2, fld.Name)
				// 			mod.AddDefinition(MakeDefinition(defName, zeroVal, fld.Typ))
				// 		}
				// 	}
				// } else {
				// 	fmt.Printf("Type '%s' not defined\n", $3)
				// }
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
				fn := MakeFunction($2)
				mod.AddFunction(fn)
				if fn, err := mod.GetCurrentFunction(); err == nil {
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
                /* { */
		/* 	$$ = $2 */
                /* } */
        |       LBRACE RBRACE
                /* { */
		/* 	$$ = nil */
                /* } */
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
        |       prolog
        |       expressionsAndStatements nonAssignExpression
        |       expressionsAndStatements assignExpression
        |       expressionsAndStatements statement
        |       expressionsAndStatements selector
        |       expressionsAndStatements stepping
        |       expressionsAndStatements debugging
        |       expressionsAndStatements affordance
        |       expressionsAndStatements remover
        |       expressionsAndStatements prolog
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
							
							typ := []byte($3)
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
        |       argumentsList assignOperator argumentsList
                {
			argsL := $1
			argsR := $3
			
			if fn, err := cxt.GetCurrentFunction(); err == nil {

				for i, argL := range argsL {
					// argL is going to be the output name
					typeParts := strings.Split(argsR[i].Typ, ".")

					var typ string
					var secondTyp string
					var idFn string

					if len(typeParts) > 1 {
						//typ = typeParts[0] // ident
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
						//idFn = MakeIdentityOpName(secondTyp)
						idFn = "identity"
					}

					// if argsR[i].Typ == "ident" {
					// 	idFn = "identity"
					// 	typ = "str"
					// } else {
					// 	idFn = MakeIdentityOpName(argsR[i].Typ)
					// 	fmt.Println(argsR[i].Typ)
					// 	typ = argsR[i].Typ
					// }

					//idFn = MakeIdentityOpName(typ)

					if op, err := cxt.GetFunction(idFn, CORE_MODULE); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = yyS[yypt-0].line + 1
						}

						fn.AddExpression(expr)
						expr.AddTag(tag)
						tag = ""

						expr.AddOutputName(string(*argL.Value))

						arg := MakeArgument(argsR[i].Value, typ)
						expr.AddArgument(arg)
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
						// var outNames []string
						// var args []*CXArgument
						
						for i, out := range op.Outputs {
							outNames[i] = MakeGenSym(NON_ASSIGN_PREFIX)
							byteName := []byte(outNames[i])
							args[i] = MakeArgument(&byteName, fmt.Sprintf("ident.%s", out.Typ))
							//args[i] = MakeArgument(&byteName, "ident")
							expr.AddOutputName(outNames[i])
						}
						
						$$ = args
					} else {
						fmt.Printf("Function '%s' not defined\n", $1)
					}
				}
			}
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

						label := []byte($2)
						expr.AddArgument(MakeArgument(&label, "str"))
					}
				}
			}
                }
        |       IF nonAssignExpression
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

					//predVal := []byte($2)
					predVal := $2[0].Value

					goToExpr.AddArgument(MakeArgument(predVal, "ident"))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
				}
			}
                }
        |       IF IDENT
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
					//elseLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>5 - 3))
					thenLines := encoder.Serialize(int32(1))

					predVal := []byte($2)

					goToExpr.AddArgument(MakeArgument(&predVal, "ident"))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
				}
			}
                }
        |       IF BOOLEAN
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

					predVal := encoder.Serialize($2)
					
					goToExpr.AddArgument(MakeArgument(&predVal, "bool"))
					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
				}
			}
                }
        |       FOR nonAssignExpression
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
					
					//predVal := []byte($2)
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
						thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>5 - 1))

						alwaysTrue := encoder.Serialize(int32(1))

						goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
						goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
						goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
					}
					
				}
			}
                }
        |       FOR IDENT
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
					
					predVal := []byte($2)
					
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
						thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>5))

						alwaysTrue := encoder.Serialize(int32(1))

						goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
						goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
						goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))
					}
					
				}
			}
                }
        |       FOR BOOLEAN
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
        |       FOR forLoopAssignExpression ';' nonAssignExpression
                {//$<int>5
			if fn, err := cxt.GetCurrentFunction(); err == nil {

				//beforeGoTo := len(fn.Expressions)
				
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

				//afterGoTo := len(fn.Expressions)
				//$<i>$ = afterGoTo - beforeGoTo
				$<i>$ = len(fn.Expressions)
			}
                }
                ';' forLoopAssignExpression
                {//$<int>8
			//increment goTo
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				goToExpr := fn.Expressions[$<i>5 - 1]
				if $<bool>7 {
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

					thenLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>5 + 1))
					// elseLines := encoder.Serialize(int32(0)) // this is added later in $12

					predVal := $4[0].Value
					
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
					goToExpr := fn.Expressions[$<i>8 - 1]

					if $<bool>7 {
						predVal := $4[0].Value

						thenLines := encoder.Serialize(int32(-($<i>8 - $<i>5 + 1)))
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

							thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>8 - 2))
							elseLines := encoder.Serialize(int32(0))

							goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
							goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
							goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))

							condGoToExpr := fn.Expressions[$<i>5 - 1]

							condThenLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>5 + 1))
							
							condGoToExpr.AddArgument(MakeArgument(&condThenLines, "i32"))
						}
					} else {
						predVal := $4[0].Value

						thenLines := encoder.Serialize(int32(1))
						elseLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>5 + 2))
						
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

							thenLines := encoder.Serialize(int32(-len(fn.Expressions) + $<i>5 - 1))
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
						
						typ := []byte($3)
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

argument:
                INT
                {
			val := encoder.Serialize($1)
                        $$ = MakeArgument(&val, "i32")
                }
        |       FLOAT
                {
			
			val := encoder.Serialize($1)
			$$ = MakeArgument(&val, "f32")
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
                        
                        val := []byte(str)
                        $$ = MakeArgument(&val, "str")
                }
        |       IDENT
                {
			val := []byte($1)
                        $$ = MakeArgument(&val, "ident")
                }
        |       IDENT LBRACK INT RBRACK
                {
			val := []byte(fmt.Sprintf("%s[%d", $1, $3))
			$$ = MakeArgument(&val, "ident")
                }
        |       typeSpecifier LBRACE argumentsList RBRACE
                {
			switch $1 {
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
				$$ = MakeArgument(&vals, "[]byte")
			case "[]i32":
				vals := make([]int32, len($3))
				for i, arg := range $3 {
					if arg.Typ == "ident" {
						// if fn, err := cxt.GetCurrentFunction(); err == nil {
						// 	if op, err := cxt.GetFunction("[]i32.write", CORE_MODULE); err == nil {
						// 		expr := MakeExpression(op)
						// 		fn.AddExpression(expr)
						// 		expr.AddArgument()
						// 	}
							
						// } else {
						// 	fmt.Println(err)
						// }
					} else {
						var val int32
						encoder.DeserializeRaw(*arg.Value, &val)
						vals[i] = val
					}
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

