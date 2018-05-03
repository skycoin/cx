%{
	package main
	import (
		"fmt"
		"strings"
		//"github.com/skycoin/cx/cx/cx0"
		"github.com/skycoin/skycoin/src/cipher/encoder"
		. "github.com/skycoin/cx/cx-master/src/base"
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

	parameter *CXArgument
	parameters []*CXArgument

	argument *CXArgument
	arguments []*CXArgument

        definition *CXArgument
	definitions []*CXArgument

	expression *CXExpression
	expressions []*CXExpression

	field *CXArgument
	fields []*CXArgument

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

%type   <tok>           assignOperator relationalOp
                        
%type   <parameter>     parameter
%type   <parameters>    parameters functionParameters
%type   <argument>      argument definitionAssignment
%type   <arguments>     arguments argumentsList nonAssignExpression conditionControl returnArg
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
                globalDeclaration
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
			Import($2)
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
			Affordance(AFF_FUNC, AFF_TYP1, $3, "", int32(0))
                }
        |       AFF FUNC IDENT LBRACE INT RBRACE
                {
			Affordance(AFF_FUNC, AFF_TYP2, $3, "", $5)
                }
        |       AFF FUNC IDENT LBRACE STRING RBRACE
                {
			Affordance(AFF_FUNC, AFF_TYP3, $3, $5, int32(0))
                }
        |       AFF FUNC IDENT LBRACE STRING INT RBRACE
                {
			Affordance(AFF_FUNC, AFF_TYP4, $3, $5, $6)
                }
                /* Module Affordances */
        |       AFF PACKAGE IDENT
                {
			Affordance(AFF_PKG, AFF_TYP1, $3, "", int32(0))
                }
        |       AFF PACKAGE IDENT LBRACE INT RBRACE
                {
			Affordance(AFF_PKG, AFF_TYP2, $3, "", $5)
                }
        |       AFF PACKAGE IDENT LBRACE STRING RBRACE
                {
			Affordance(AFF_PKG, AFF_TYP3, $3, $5, int32(0))
                }
        |       AFF PACKAGE IDENT LBRACE STRING INT RBRACE
                {
			Affordance(AFF_PKG, AFF_TYP4, $3, $5, $6)
                }
                /* Struct Affordances */
        |       AFF STRUCT IDENT
                {
			Affordance(AFF_STRCT, AFF_TYP1, $3, "", int32(0))
                }
        |       AFF STRUCT IDENT LBRACE INT RBRACE
                {
			Affordance(AFF_STRCT, AFF_TYP2, $3, "", $5)
                }
        |       AFF STRUCT IDENT LBRACE STRING RBRACE
                {
			Affordance(AFF_STRCT, AFF_TYP3, $3, $5, int32(0))
                }
        |       AFF STRUCT IDENT LBRACE STRING INT RBRACE
                {
			Affordance(AFF_STRCT, AFF_TYP4, $3, $5, $6)
                }
                /* Struct Affordances */
        |       AFF EXPR IDENT
                {
			Affordance(AFF_EXPR, AFF_TYP1, $3, "", int32(0))
                }
        |       AFF EXPR IDENT LBRACE INT RBRACE
                {
			Affordance(AFF_EXPR, AFF_TYP2, $3, "", $5)
                }
        |       AFF EXPR IDENT LBRACE STRING RBRACE
                {
			Affordance(AFF_EXPR, AFF_TYP3, $3, $5, int32(0))
                }
        |       AFF EXPR IDENT LBRACE STRING INT RBRACE
                {
			Affordance(AFF_EXPR, AFF_TYP4, $3, $5, $6)
                }
        ;

stepping:       TSTEP INT INT
                {
			Stepping(int($2), int($3), true)
                }
        |       STEP INT
                {
			Stepping(int($2), 0, false)
                }
        ;

debugging:      DSTATE
                {
			DebugState()
                }
        |       DSTACK
                {
			DebugStack()
                }
        |       DPROGRAM
                {
			cxt.PrintProgram(false)
                }
        ;

remover:        REM FUNC IDENT
                {
			Remover(REM_TYP_FUNC, $3, "")
                }
        |       REM PACKAGE IDENT
                {
			Remover(REM_TYP_PKG, $3, "")
                }
        |       REM DEF IDENT
                {
			Remover(REM_TYP_GLBL, $3, "")
                }
        |       REM STRUCT IDENT
                {
			Remover(REM_TYP_STRCT, $3, "")
                }
        |       REM IMPORT STRING
                {
			Remover(REM_TYP_STRCT, $3, "")
                }
        |       REM EXPR IDENT FUNC IDENT
                {
			Remover(REM_TYP_EXPR, $3, $5)
                }
        |       REM FIELD IDENT STRUCT IDENT
                {
			Remover(REM_TYP_FLD, $3, $5)
                }
        |       REM INPUT IDENT FUNC IDENT
                {
			Remover(REM_TYP_INPUT, $3, $5)
                }
        |       REM OUTPUT IDENT FUNC IDENT
                {
			Remover(REM_TYP_OUTPUT, $3, $5)
                }
                // no, too complex. just wipe out entire expression
        // |       REM ARG INT EXPR INT FUNC IDENT
        //         {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			if fn, err := mod.Program.GetFunction($5, mod.Name); err == nil {
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
			SelectorFields($2)
                }
        ;

selector:       SPACKAGE IDENT
                {
			$<string>$ = Selector($2, SELECT_TYP_PKG)
                }
                selectorLines
                {
			if $<bool>4 {
				if _, err := cxt.SelectModule($<string>3); err == nil {
				}
			}
                }
        |       SFUNC IDENT
                {
			$<string>$ = Selector($2, SELECT_TYP_FUNC)
                }
                selectorExpressionsAndStatements
                {
			if $<bool>4 {
				if _, err := cxt.SelectFunction($<string>3); err == nil {
				}
			}
                }
        |       SSTRUCT IDENT
                {
			$<string>$ = Selector($2, SELECT_TYP_STRCT)
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

/* typeSpecifier: */
/*                 I32 {$$ = $1} */
/*         |       I64 {$$ = $1} */
/*         |       F32 {$$ = $1} */
/*         |       F64 {$$ = $1} */
/*         |       BOOL {$$ = $1} */
/*         |       BYTE {$$ = $1} */
/*         |       BOOLA {$$ = $1} */
/*         |       STRA {$$ = $1} */
/*         |       BYTEA {$$ = $1} */
/*         |       I32A {$$ = $1} */
/*         |       I64A {$$ = $1} */
/*         |       F32A {$$ = $1} */
/*         |       F64A {$$ = $1} */
/*         |       STR {$$ = $1} */
/*         |       NEW IDENT {$$ = $2} */
/*         ; */

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

globalDeclaration:
                VAR IDENT BASICTYPE definitionAssignment
                {
			GlobalDeclaration(true, $2, $3, $4, yyS[yypt-0].line + 1)
                }
        |       VAR IDENT IDENT
                {
			GlobalDeclaration(false, $2, $3, nil, yyS[yypt-0].line + 1)
                }
        ;

fields:
                parameter
                {
			var flds []*CXArgument
                        flds = append(flds, MakeFieldFromParameter($1))
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
			StructDeclaration($2, yyS[yypt-0].line + 1)
                }
                STRUCT structFields
                {
			StructDeclarationFields($5)
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
                FUNC functionParameters IDENT functionParameters functionParameters
                {
			FunctionDeclarationHeader(METHOD_INP_OUT, $3, $2, $4, $5, yyS[yypt-0].line + 1)
                }
                functionStatements
        |       FUNC functionParameters IDENT functionParameters
                {
			FunctionDeclarationHeader(METHOD_INP, $3, $2, $4, nil, yyS[yypt-0].line + 1)
                }
                functionStatements
                /* Functions */
        |       FUNC IDENT functionParameters
                {
			FunctionDeclarationHeader(FUNC_INP, $2, nil, $3, nil, yyS[yypt-0].line + 1)
                }
                functionStatements
        |       FUNC IDENT functionParameters functionParameters
                {
			FunctionDeclarationHeader(FUNC_INP_OUT, $2, nil, $3, $4, yyS[yypt-0].line + 1)
                }
                functionStatements
        ;

parameter:
                IDENT BASICTYPE
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
                VAR IDENT BASICTYPE definitionAssignment
                {
			AssignBasicVar($2, $3, $4, yyS[yypt-0].line + 1)
                }
        |       VAR IDENT LBRACK RBRACK IDENT
                {
			AssignCustomVar($2, $5, yyS[yypt-0].line + 1)
                }
        |       argumentsList assignOperator argumentsList
                {
			AssignExpression($1, $2, $3, yyS[yypt-0].line + 1)
                }
        ;

nonAssignExpression:
                IDENT arguments
                {
			$$ = NonAssignFunctionCall($1, $2, yyS[yypt-0].line + 1)
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

returnArg:
                ';'
                {
                    $$ = nil
                }
        |       argumentsList
                {
			$$ = $1
                }
                ;

statement:      RETURN returnArg
                {
			StatementReturn($2, yyS[yypt-0].line + 1)
                }
        |       GOTO IDENT
                {
			StatementGoTo($2, yyS[yypt-0].line + 1)
                }
        |       IF conditionControl
                {
			StatementIfCondition(yyS[yypt-0].line + 1)
                }
                LBRACE
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                expressionsAndStatements RBRACE elseStatement
                {
			StatementIfElse($<i>5, $2, $<i>8)
                }
        |       beginFor
                nonAssignExpression
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                {
			StatementForCondExpression(yyS[yypt-0].line + 1)
                }
                LBRACE expressionsAndStatements RBRACE
                {
			StatementForFinalizer($<i>1, $2, $<i>3, true, yyS[yypt-0].line + 1)
                }
        |       beginFor
                argument
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                {
			StatementForCondExpression(yyS[yypt-0].line + 1)
                }
                LBRACE expressionsAndStatements RBRACE
                {
			StatementForFinalizer($<i>1, []*CXArgument{$2}, $<i>3, false, yyS[yypt-0].line + 1)
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
			$<i>$ = StatementForCondLenExpressions(yyS[yypt-0].line + 1)
                }
                ';' forLoopAssignExpression //$<bool>9
                {//$<int>10
			//increment goTo
			$<i>$ = StatementForLoopAssignLenExpressions($5, $<i>7, $<bool>9, yyS[yypt-0].line + 1)
                }
                LBRACE expressionsAndStatements RBRACE
                {
			StatementForThreePartsFinalizer($5, $<i>7, $<i>3, $<i>10, $<bool>9, yyS[yypt-0].line + 1)
                }
        |       VAR IDENT IDENT
                {
			VariableDeclaration($2, $3, yyS[yypt-0].line + 1)
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
			ElseStatementInitializer(yyS[yypt-0].line + 1)
                }
                LBRACE
                {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				$<i>$ = len(fn.Expressions)
			}
                }
                expressionsAndStatements RBRACE
                {
			$<i>$ = ElseStatementFinalizer($<i>4)
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
        // |       MULT relationalOp argument
        //         {
	// 		argStr := ArgToString($3)
	// 		$$ = []string{$1, argStr, $2}
        //         }
        |       AFFVAR relationalOp nonAssignExpression
                {
			var identName string
			encoder.DeserializeRaw(*$3[0].Value, &identName)
			$$ = []string{$1, identName, $2}
                }
        |       AFFVAR relationalOp AFFVAR
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
        |       IDENT VALUE inferWeight
                {
			$$ = []string{$1, $3, "weight"}
                }
        |       IDENT VALUE nonAssignExpression
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
			var defs []*CXArgument
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

/*
  Fix this, there has to be a way to compress these rules
*/
argument:       argument PLUS argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression PLUS nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument PLUS nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression PLUS argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}


        |       argument MINUS argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression MINUS nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument MINUS nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression MINUS argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument MULT argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression MULT nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument MULT nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression MULT argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument DIV argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression DIV nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument DIV nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression DIV argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument REMAINDER argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression REMAINDER nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument REMAINDER nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression REMAINDER argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument LEFTSHIFT argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression LEFTSHIFT nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument LEFTSHIFT nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression LEFTSHIFT argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument RIGHTSHIFT argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression RIGHTSHIFT nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument RIGHTSHIFT nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression RIGHTSHIFT argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument EXP argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression EXP nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument EXP nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression EXP argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument EQUAL argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression EQUAL nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument EQUAL nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression EQUAL argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument UNEQUAL argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression UNEQUAL nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument UNEQUAL nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression UNEQUAL argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument GTHAN argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression GTHAN nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument GTHAN nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression GTHAN argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument GTHANEQ argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression GTHANEQ nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument GTHANEQ nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression GTHANEQ argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument LTHAN argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression LTHAN nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument LTHAN nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression LTHAN argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument LTHANEQ argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression LTHANEQ nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument LTHANEQ nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression LTHANEQ argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument OR argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression OR nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument OR nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression OR argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument AND argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression AND nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument AND nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression AND argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument BITAND argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression BITAND nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument BITAND nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression BITAND argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument BITOR argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression BITOR nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument BITOR nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression BITOR argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument BITXOR argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression BITXOR nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument BITXOR nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression BITXOR argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       argument BITCLEAR argument
                {
			$$ = binaryOp($2, $1, $3, yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression BITCLEAR nonAssignExpression
                {
			$$ = binaryOp($2, $1[0], $3[0], yyS[yypt-0].line + 1)
		}
        |       argument BITCLEAR nonAssignExpression
                {
			$$ = binaryOp($2, $1, $3[0], yyS[yypt-0].line + 1)
		}
        |       nonAssignExpression BITCLEAR argument
                {
			$$ = binaryOp($2, $1[0], $3, yyS[yypt-0].line + 1)
		}

        |       NOT argument
                {
			$$ = UnaryPrefixOp($2, nil, true, yyS[yypt-0].line + 1)
		}
        |       NOT nonAssignExpression
                {
			$$ = UnaryPrefixOp(nil, $2, false, yyS[yypt-0].line + 1)
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
        |       NEW IDENT LBRACE structLitDefs RBRACE
                {
			$$ = StructLiteralDeclaration($2, $4, yyS[yypt-0].line + 1)
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
        |       BASICTYPE LBRACE argumentsList RBRACE
                {
			$$ = BasicArrayLiteralDeclaration($1, $3, yyS[yypt-0].line + 1, false)
                }
                // empty arrays
        |       BASICTYPE LBRACE RBRACE
                {
			$$ = BasicArrayLiteralDeclaration($1, nil, yyS[yypt-0].line + 1, true)
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

argumentsList:  argument
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
