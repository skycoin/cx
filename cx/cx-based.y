%{
	package main
	import (
		"strings"
		"fmt"
		"os"
		"time"

		//"github.com/skycoin/cx/cx/cx0"
		"github.com/skycoin/skycoin/src/cipher/encoder"
		. "github.com/skycoin/cx/src/base"
                )

	var cxt = MakeContext()
	//var cxt = cx0.CXT

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
	var fileName string

	func warnf(format string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, format, args...)
		os.Stderr.Sync()
	}

	func binaryOp (op string, arg1, arg2 *CXArgument, line int) *CXArgument {
		var opName string
		var typArg1 string
		// var typArg2 string
		// _ = typArg2

		if (len(arg1.Typ) > len("ident.") && arg1.Typ[:len("ident.")] == "ident.") {
			arg1.Typ = "ident"
		}

		if (len(arg2.Typ) > len("ident.") && arg2.Typ[:len("ident.")] == "ident.") {
			arg2.Typ = "ident"
		}
		
		if arg1.Typ == "ident" {
			var identName string
			encoder.DeserializeRaw(*arg1.Value, &identName)

			if typ, err := GetIdentType(identName, line, fileName, cxt); err == nil {
				typArg1 = typ
			} else {
				fmt.Println(err)
			}
		} else {
			typArg1 = arg1.Typ
		}

		// if arg2.Typ == "ident" {
		// 	var identName string
		// 	encoder.DeserializeRaw(*arg2.Value, &identName)

		// 	if typ, err := GetIdentType(identName, line, fileName, cxt); err == nil {
		// 		typArg2 = typ
		// 	} else {
		// 		fmt.Println(err)
		// 	}
		// } else {
		// 	typArg2 = arg1.Typ
		// }

		
		// fmt.Println(typArg1)
		// fmt.Println(arg1.Typ)

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
		case "&":
			opName = fmt.Sprintf("%s.bitand", typArg1)
		case "|":
			opName = fmt.Sprintf("%s.bitor", typArg1)
		case "^":
			opName = fmt.Sprintf("%s.bitxor", typArg1)
		case "&^":
			opName = fmt.Sprintf("%s.bitclear", typArg1)
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
					expr.FileName = fileName
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

			if typ, err := GetIdentType(identName, line, fileName, cxt); err == nil {
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
					expr.FileName = fileName
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

%start                  translation_unit
%%

translation_unit:
                external_declaration
        |       translation_unit external_declaration
        ;

external_declaration:
                package_definition
        |       global_declaration
        |       struct_declaration
        ;

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
                ;

definitionAssignment:
                /* empty */
        |       assignOperator argument
        |       assignOperator ADDR argument
        |       assignOperator VALUE argument
                ;

definitionDeclaration:
                VAR IDENT BASICTYPE definitionAssignment
        |       VAR IDENT IDENT
        ;

fields:
                parameter
        |       ';'
        |       debugging
        |       fields parameter
        |       fields ';'
        |       fields debugging
                ;

structFields:
                LBRACE fields RBRACE
        |       LBRACE RBRACE
        ;

structDeclaration:
                TYPSTRUCT IDENT
                STRUCT structFields
        ;

functionParameters:
                LPAREN parameters RPAREN
        |       LPAREN RPAREN
        ;

functionDeclaration:
                /* Methods */
                FUNC functionParameters IDENT functionParameters functionParameters
                functionStatements
        |       FUNC functionParameters IDENT functionParameters
                functionStatements
                /* Functions */
        |       FUNC IDENT functionParameters
                functionStatements
        |       FUNC IDENT functionParameters functionParameters
                functionStatements
        ;

parameter:
                IDENT BASICTYPE
        |       IDENT IDENT
        |       IDENT MULT IDENT
        ;

parameters:
                parameter
        |       parameters COMMA parameter
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

structLiteral:
        |       
		LBRACE structLitDefs RBRACE
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
