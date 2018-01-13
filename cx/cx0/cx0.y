%{
	package cx0
	import (
		"fmt"
		"bytes"
		. "github.com/skycoin/cx/src/base"
	)

	var PRGRM *CXProgram

	var replMode bool = false
	var inREPL bool = false
	var inFn bool = false
	var fileName string

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

%type   <argument>      parameter
%type   <arguments>     fields structFields parameters functionParameters

%left                   OR
%left                   AND
%left                   BITCLEAR
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
        |       lines ';'
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
        ;

affordance:
                TAG
                /* Function Affordances */
        |       AFF FUNC IDENT
        |       AFF FUNC IDENT LBRACE INT RBRACE
        |       AFF FUNC IDENT LBRACE STRING RBRACE
        |       AFF FUNC IDENT LBRACE STRING INT RBRACE
                /* Module Affordances */
        |       AFF PACKAGE IDENT
        |       AFF PACKAGE IDENT LBRACE INT RBRACE
        |       AFF PACKAGE IDENT LBRACE STRING RBRACE
        |       AFF PACKAGE IDENT LBRACE STRING INT RBRACE
                /* Struct Affordances */
        |       AFF STRUCT IDENT
        |       AFF STRUCT IDENT LBRACE INT RBRACE
        |       AFF STRUCT IDENT LBRACE STRING RBRACE
        |       AFF STRUCT IDENT LBRACE STRING INT RBRACE
                /* Struct Affordances */
        |       AFF EXPR IDENT
        |       AFF EXPR IDENT LBRACE INT RBRACE
        |       AFF EXPR IDENT LBRACE STRING RBRACE
        |       AFF EXPR IDENT LBRACE STRING INT RBRACE
        ;

stepping:       TSTEP INT INT
        |       STEP INT
        ;

debugging:      DSTATE
        |       DSTACK
        |       DPROGRAM
        ;

remover:        REM FUNC IDENT
        |       REM PACKAGE IDENT
        |       REM DEF IDENT
        |       REM STRUCT IDENT
        |       REM IMPORT STRING
        |       REM EXPR IDENT FUNC IDENT
        |       REM FIELD IDENT STRUCT IDENT
        |       REM INPUT IDENT FUNC IDENT
        |       REM OUTPUT IDENT FUNC IDENT
        ;

selectorLines:
                /* empty */
        |       LBRACE lines RBRACE
        ;

selectorExpressionsAndStatements:
                /* empty */
        |       LBRACE expressionsAndStatements RBRACE
        ;

selectorFields:
                /* empty */
        |       LBRACE fields RBRACE
        ;

selector:       SPACKAGE IDENT
                selectorLines
        |       SFUNC IDENT
                selectorExpressionsAndStatements
        |       SSTRUCT IDENT
                selectorFields
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

packageDeclaration:
                PACKAGE IDENT
                {
			mod := MakeModule($2)
			PRGRM.AddModule(mod)
                }
                ;

definitionAssignment:
                /* empty */
        |       assignOperator argument
        |       assignOperator ADDR argument
        |       assignOperator VALUE argument
                ;

definitionDeclaration:
                VAR IDENT BASICTYPE definitionAssignment
                {
			if mod, err := PRGRM.GetCurrentModule(); err == nil {
				def := MakeGlobal($2, TypeNameToInt($3))
				mod.AddGlobal(def)
			} else {
				panic(err)
			}
			
                }
        |       VAR IDENT IDENT
                {
			if mod, err := PRGRM.GetCurrentModule(); err == nil {
				def := MakeGlobal($2, TypeNameToInt($3))
				mod.AddGlobal(def)
			} else {
				panic(err)
			}
                }
        ;

fields:
                parameter
                {
			var flds []*CXArgument
                        flds = append(flds, $1)
			$$ = flds
                }
        |       ';'
                {
			var flds []*CXArgument
			$$ = flds
                }
        |       debugging
                {
			var flds []*CXArgument
			$$ = flds
                }
        |       fields parameter
                {
			$1 = append($1, $2)
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
			if mod, err := PRGRM.GetCurrentModule(); err == nil {
				strct := MakeStruct($2)
				mod.AddStruct(strct)
			}
                }
                STRUCT structFields
                {
			if strct, err := PRGRM.GetCurrentStruct(); err == nil {
				for _, fld := range $5 {
					fldFromParam := MakeField(fld.Name, fld.Type)
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
                /* Methods */
                FUNC functionParameters IDENT functionParameters functionParameters functionStatements
                {
			if len($2) > 1 {
				panic(fmt.Sprintf("%s: %d: method '%s' has multiple receivers", fileName, yyS[yypt-0].line+1, $3))
			}

			if mod, err := PRGRM.GetCurrentModule(); err == nil {
				if IsBasicType($2[0].Type) {
					panic(fmt.Sprintf("%s: %d: cannot define methods on basic type %s", fileName, yyS[yypt-0].line+1, $2[0].Type))
				}
				
				inFn = true
				fn := MakeFunction(fmt.Sprintf("%s.%s", $2[0].Type, $3))
				mod.AddFunction(fn)
				if fn, err := mod.GetCurrentFunction(); err == nil {

					//checking if there are duplicate parameters
					dups := append($4, $5...)
					dups = append(dups, $2...)
					for _, param := range dups {
						for _, dup := range dups {
							if param.Name == dup.Name && param != dup {
								panic(fmt.Sprintf("%s: %d: duplicate receiver, input and/or output parameters in method '%s'", fileName, yyS[yypt-0].line+1, $3))
							}
						}
					}

					for _, rec := range $2 {
						fn.AddInput(rec)
					}
					for _, inp := range $4 {
						fn.AddInput(inp)
					}
					for _, out := range $5 {
						fn.AddOutput(out)
					}
				}
			}
                }
        |       FUNC functionParameters IDENT functionParameters functionStatements
                {
			if len($2) > 1 {
				panic(fmt.Sprintf("%s: %d: method '%s' has multiple receivers", fileName, yyS[yypt-0].line+1, $3))
			}
			
			if mod, err := PRGRM.GetCurrentModule(); err == nil {
				if IsBasicType($2[0].Type) {
					panic(fmt.Sprintf("%s: %d: cannot define methods on basic type %s", fileName, yyS[yypt-0].line+1, $2[0].Type))
				}
				
				inFn = true
				fn := MakeFunction(fmt.Sprintf("%s.%s", $2[0].Type, $3))
				mod.AddFunction(fn)
				if fn, err := mod.GetCurrentFunction(); err == nil {

					//checking if there are duplicate parameters
					dups := append($2, $4...)
					for _, param := range dups {
						for _, dup := range dups {
							if param.Name == dup.Name && param != dup {
								panic(fmt.Sprintf("%s: %d: duplicate receiver, input and/or output parameters in method '%s'", fileName, yyS[yypt-0].line+1, $3))
							}
						}
					}

					for _, rec := range $2 {
						fn.AddInput(rec)
					}
					for _, inp := range $4 {
						fn.AddInput(inp)
					}
				}
			}
                }
                /* Functions */
        |       FUNC IDENT functionParameters functionStatements
                {
			if mod, err := PRGRM.GetCurrentModule(); err == nil {
				inFn = true
				fn := MakeFunction($2)
				mod.AddFunction(fn)
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, inp := range $3 {
						fn.AddInput(inp)
					}
				}
			}
                }
        |       FUNC IDENT functionParameters functionParameters functionStatements
                {
			if mod, err := PRGRM.GetCurrentModule(); err == nil {
				inFn = true
				fn := MakeFunction($2)
				mod.AddFunction(fn)
				if fn, err := mod.GetCurrentFunction(); err == nil {

					//checking if there are duplicate parameters
					dups := append($3, $4...)
					for _, param := range dups {
						for _, dup := range dups {
							if param.Name == dup.Name && param != dup {
								panic(fmt.Sprintf("%s: %d: duplicate input and/or output parameters in function '%s'", fileName, yyS[yypt-0].line+1, $2))
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
        ;

parameter:
                IDENT BASICTYPE
                {
			$$ = MakeParameter($1, TypeNameToInt($2))
                }
        |       IDENT IDENT
                {
			$$ = MakeParameter($1, TypeNameToInt($2))
                }
        |       IDENT MULT IDENT
                {
			typ := TypeNameToInt($3)
			param := MakeParameter($1, typ)
			param.IsPointer = true
			$$ = param
                }
        ;

parameters:
                parameter
                {
			var params []*CXArgument
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
        |       LBRACE RBRACE
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
                VAR IDENT BASICTYPE definitionAssignment
        |       VAR IDENT LBRACK RBRACK IDENT
        |       argumentsList assignOperator argumentsList
        ;

nonAssignExpression:
                IDENT arguments
        |       argument PLUSPLUS
        |       argument MINUSMINUS
        ;

beginFor:       FOR
                ;

conditionControl:
                nonAssignExpression
        |       argument
        ;

returnArg:
                ';'
        |       argumentsList
                ;

statement:      RETURN returnArg
        |       GOTO IDENT
        |       IF conditionControl
                LBRACE
                expressionsAndStatements RBRACE elseStatement
        |       beginFor
                nonAssignExpression
                LBRACE expressionsAndStatements RBRACE
        |       beginFor
                argument
                LBRACE expressionsAndStatements RBRACE
        |       beginFor // $<i>1
                forLoopAssignExpression // $2
                ';' conditionControl
                ';' forLoopAssignExpression //$<bool>9
                LBRACE expressionsAndStatements RBRACE
        |       VAR IDENT IDENT
        |       ';'
        ;

forLoopAssignExpression:
        |       assignExpression
        |       nonAssignExpression
        ;

elseStatement:
                /* empty */
        |       ELSE
                LBRACE
                expressionsAndStatements RBRACE
                ;

expressions:
                nonAssignExpression
        |       assignExpression
        |       expressions nonAssignExpression
        |       expressions assignExpression
        ;

inferPred:      inferObj
        |       inferCond
        |       inferPred COMMA inferObj
        |       inferPred COMMA inferCond
        ;

inferCond:      IDENT LPAREN inferPred RPAREN
        |       BOOLEAN
        ;

relationalOp:   EQUAL
        |       GTHAN
        |       LTHAN
        |       UNEQUAL
                ;

inferActionArg:
                inferObj
        |       IDENT
        |       AFFVAR relationalOp argument
        |       AFFVAR relationalOp nonAssignExpression
        |       AFFVAR relationalOp AFFVAR
        ;

inferAction:
		IDENT LPAREN inferActionArg RPAREN
        ;

inferActions:
                inferAction
        |       inferActions inferAction
                ;

inferRule:      IF inferCond LBRACE inferActions RBRACE
        |       IF inferObj LBRACE inferActions RBRACE
        ;

inferRules:     inferRule
        |       inferRules inferRule
        ;

inferWeight:    FLOAT
        |       INT
        |       IDENT
        ;

inferObj:
        |       IDENT VALUE inferWeight
        |       IDENT VALUE nonAssignExpression
;

inferObjs:      inferObj
        |       inferObjs COMMA inferObj
        ;

inferTarget:    IDENT LPAREN IDENT RPAREN
        ;

inferTargets:   inferTarget
        |       inferTargets inferTarget
        ;

inferClauses:   inferObjs
        |       inferRules
        |       inferTargets
        ;




structLitDef:
                TAG argument
        |       TAG nonAssignExpression
                    
;

structLitDefs:  structLitDef
        |       structLitDefs COMMA structLitDef
        ;

/*
  Fix this, there has to be a way to compress these rules
*/
argument:       argument PLUS argument
        |       nonAssignExpression PLUS nonAssignExpression
        |       argument PLUS nonAssignExpression
        |       nonAssignExpression PLUS argument


        |       argument MINUS argument
        |       nonAssignExpression MINUS nonAssignExpression
        |       argument MINUS nonAssignExpression
        |       nonAssignExpression MINUS argument

        |       argument MULT argument
        |       nonAssignExpression MULT nonAssignExpression
        |       argument MULT nonAssignExpression
        |       nonAssignExpression MULT argument

        |       argument DIV argument
        |       nonAssignExpression DIV nonAssignExpression
        |       argument DIV nonAssignExpression
        |       nonAssignExpression DIV argument

        |       argument REMAINDER argument
        |       nonAssignExpression REMAINDER nonAssignExpression
        |       argument REMAINDER nonAssignExpression
        |       nonAssignExpression REMAINDER argument

        |       argument LEFTSHIFT argument
        |       nonAssignExpression LEFTSHIFT nonAssignExpression
        |       argument LEFTSHIFT nonAssignExpression
        |       nonAssignExpression LEFTSHIFT argument

        |       argument RIGHTSHIFT argument
        |       nonAssignExpression RIGHTSHIFT nonAssignExpression
        |       argument RIGHTSHIFT nonAssignExpression
        |       nonAssignExpression RIGHTSHIFT argument

        |       argument EXP argument
        |       nonAssignExpression EXP nonAssignExpression
        |       argument EXP nonAssignExpression
        |       nonAssignExpression EXP argument

        |       argument EQUAL argument
        |       nonAssignExpression EQUAL nonAssignExpression
        |       argument EQUAL nonAssignExpression
        |       nonAssignExpression EQUAL argument

        |       argument UNEQUAL argument
        |       nonAssignExpression UNEQUAL nonAssignExpression
        |       argument UNEQUAL nonAssignExpression
        |       nonAssignExpression UNEQUAL argument

        |       argument GTHAN argument
        |       nonAssignExpression GTHAN nonAssignExpression
        |       argument GTHAN nonAssignExpression
        |       nonAssignExpression GTHAN argument

        |       argument GTHANEQ argument
        |       nonAssignExpression GTHANEQ nonAssignExpression
        |       argument GTHANEQ nonAssignExpression
        |       nonAssignExpression GTHANEQ argument

        |       argument LTHAN argument
        |       nonAssignExpression LTHAN nonAssignExpression
        |       argument LTHAN nonAssignExpression
        |       nonAssignExpression LTHAN argument

        |       argument LTHANEQ argument
        |       nonAssignExpression LTHANEQ nonAssignExpression
        |       argument LTHANEQ nonAssignExpression
        |       nonAssignExpression LTHANEQ argument

        |       argument OR argument
        |       nonAssignExpression OR nonAssignExpression
        |       argument OR nonAssignExpression
        |       nonAssignExpression OR argument

        |       argument AND argument
        |       nonAssignExpression AND nonAssignExpression
        |       argument AND nonAssignExpression
        |       nonAssignExpression AND argument

        |       argument BITAND argument
        |       nonAssignExpression BITAND nonAssignExpression
        |       argument BITAND nonAssignExpression
        |       nonAssignExpression BITAND argument

        |       argument BITOR argument
        |       nonAssignExpression BITOR nonAssignExpression
        |       argument BITOR nonAssignExpression
        |       nonAssignExpression BITOR argument

        |       argument BITXOR argument
        |       nonAssignExpression BITXOR nonAssignExpression
        |       argument BITXOR nonAssignExpression
        |       nonAssignExpression BITXOR argument

        |       argument BITCLEAR argument
        |       nonAssignExpression BITCLEAR nonAssignExpression
        |       argument BITCLEAR nonAssignExpression
        |       nonAssignExpression BITCLEAR argument

        |       NOT argument
        |       NOT nonAssignExpression
        |       LPAREN argument RPAREN
        |       BYTENUM
        |       INT
        |       LONG
        |       FLOAT
        |       DOUBLE
        |       BOOLEAN
        |       STRING
        |       IDENT
        |       NEW IDENT LBRACE structLitDefs RBRACE
        |       IDENT LBRACK INT RBRACK
        |       INFER LBRACE inferClauses RBRACE
        |       BASICTYPE LBRACE argumentsList RBRACE
                // empty arrays
        |       BASICTYPE LBRACE RBRACE
        ;

arguments:
                LPAREN argumentsList RPAREN
        |       LPAREN RPAREN
        ;

argumentsList:  argument
        |       nonAssignExpression
        |       ADDR argument
        |       VALUE argument
        |       VALUE nonAssignExpression
        |       argumentsList COMMA argument
        |       argumentsList COMMA nonAssignExpression
        |       argumentsList COMMA ADDR argument
        |       argumentsList COMMA VALUE argument
        |       argumentsList COMMA VALUE nonAssignExpression
        ;

%%
