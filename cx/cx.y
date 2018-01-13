%{
	package main
	import (
		"strings"
		"fmt"
		"github.com/skycoin/skycoin/src/cipher/encoder"
		. "github.com/skycoin/cx/src/base"
                )

	var prgrm = MakeProgram(1024, 1024, 1024)

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
	var tag string = ""
	var asmNL = "\n"
	var fileName string
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
%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK IDENTIFIER
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

                        STRING_LITERAL
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

//%type   <tok>           assignOperator
                        
//%type   <argument>      argument definitionAssignment parameter
//%type   <arguments>     arguments argumentsList nonAssignExpression conditionControl returnArg parameters functionParameters
// %type   <expression>    assignExpression
// %type   <expressions>   elseStatement

// %type   <bool>          selectorExpressionsAndStatements

%%

primary_expression:
                IDENTIFIER
                string
        | '('   expression ')'
                ;

enumeration_constant:
                IDENTIFIER
                ;

string:         STRING_LITERAL
                ;

postfix_expression:
                primary_expression
	|       postfix_expression '[' expression ']'
	|       postfix_expression '(' ')'
	|       postfix_expression '(' argument_expression_list ')'
	|       postfix_expression '.' IDENTIFIER
	|       postfix_expression PTR_OP IDENTIFIER // check
	|       postfix_expression INC_OP
	|       postfix_expression DEC_OP
	| '('   type_name ')' '{' initializer_list '}' // check
	| '('   type_name ')' '{' initializer_list ',' '}' // check
                ;

argument_expression_list:
                assignment_expression
	|       argument_expression_list ',' assignment_expression
                ;

unary_expression:
                postfix_expression
	|       INC_OP unary_expression
	|       DEC_OP unary_expression
	|       unary_operator cast_expression // check
                ;

unary_operator:
                '&'
	|       '*'
	|       '+'
	|       '-'
	|       '~' // check
	|       '!'
                ;

cast_expression: // check
                unary_expression
	| '('   type_name ')' cast_expression // check
                ;

multiplicative_expression:
                cast_expression
	|       multiplicative_expression '*' cast_expression
	|       multiplicative_expression '/' cast_expression
	|       multiplicative_expression '%' cast_expression
                ;

additive_expression:
                multiplicative_expression
	|       additive_expression '+' multiplicative_expression
	|       additive_expression '-' multiplicative_expression
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
                ;

logical_or_expression:
                logical_and_expression
	|       logical_or_expression OR_OP logical_and_expression
                ;

conditional_expression:
                logical_or_expression
	|       logical_or_expression '?' expression ':' conditional_expression
                ;

assignment_expression:
                conditional_expression
	|       unary_expression assignment_operator assignment_expression
                ;

assignment_operator:
                '='
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
	|       expression ',' assignment_expression
                ;

constant_expression:
                conditional_expression
                ;

declaration:    declaration_specifiers ';'
	|       declaration_specifiers init_declarator_list ';'
                ;

declaration_specifiers:
		type_specifier declaration_specifiers
	|       type_specifier
	|       type_qualifier declaration_specifiers
	|       type_qualifier
                ;

init_declarator_list:
                init_declarator
	|       init_declarator_list ',' init_declarator
                ;

init_declarator:
                declarator '=' initializer
	|       declarator
                ;

// storage_class_specifier:
//                 TYPEDEF	/* identifiers must be flagged as TYPEDEF_NAME */
// 	|       EXTERN
// 	|       STATIC
// 	|       THREAD_LOCAL
// 	|       AUTO
// 	|       REGISTER
//                 ;

type_specifier: BOOL
        |       BYTE
        |       F32
        |       F64
        |       I8
        |       I16
        |       I32
        |       I64
        |       STR
        |       UI8
        |       UI16
        |       UI32
        |       UI64
	|       struct_or_union_specifier
	|       enum_specifier
	/* |       TYPEDEF_NAME // check */
                ;

struct_or_union_specifier:
                struct_or_union '{' struct_declaration_list '}'
	|       struct_or_union IDENTIFIER '{' struct_declaration_list '}'
	|       struct_or_union IDENTIFIER
                ;

struct_or_union:
                STRUCT
	|       UNION // check
                ;

struct_declaration_list:
                struct_declaration
	|       struct_declaration_list struct_declaration
                ;

struct_declaration:
                specifier_qualifier_list ';'	/* for anonymous struct/union */
	|       specifier_qualifier_list struct_declarator_list ';'
                ;

specifier_qualifier_list:
                type_specifier specifier_qualifier_list
	|       type_specifier
	|       type_qualifier specifier_qualifier_list
	|       type_qualifier
                ;

struct_declarator_list:
                struct_declarator
	|       struct_declarator_list ',' struct_declarator
                ;

struct_declarator:
                ':' constant_expression
	|       declarator ':' constant_expression
	|       declarator
                ;

enum_specifier: ENUM '{' enumerator_list '}'
	|       ENUM '{' enumerator_list ',' '}'
	|       ENUM IDENTIFIER '{' enumerator_list '}'
	|       ENUM IDENTIFIER '{' enumerator_list ',' '}'
	|       ENUM IDENTIFIER
                ;

enumerator_list:enumerator
	|       enumerator_list ',' enumerator
                ;

enumerator:     enumeration_constant '=' constant_expression
	|       enumeration_constant
                ;

type_qualifier: CONST
	/* |       RESTRICT */
	/* |       VOLATILE */
	/* |       ATOMIC */
                ;

/* function_specifier: */
/*                 INLINE */
/* 	|       NORETURN */
/*                 ; */

/* alignment_specifier: */
/*                 ALIGNAS '(' type_name ')' */
/* 	|       ALIGNAS '(' constant_expression ')' */
/*                 ; */

declarator:     pointer direct_declarator
	|       direct_declarator
                ;

direct_declarator:
                IDENTIFIER
	| '('   declarator ')'
	|       direct_declarator '[' ']'
	|       direct_declarator '[' '*' ']'
	|       direct_declarator '[' type_qualifier_list '*' ']'
	|       direct_declarator '[' type_qualifier_list assignment_expression ']'
	|       direct_declarator '[' type_qualifier_list ']'
	|       direct_declarator '[' assignment_expression ']'
	|       direct_declarator '(' parameter_type_list ')'
	|       direct_declarator '(' ')'
	|       direct_declarator '(' identifier_list ')'
                ;

pointer: '*'    type_qualifier_list pointer
	| '*'   type_qualifier_list
	| '*'   pointer
	| '*'
	;

                type_qualifier_list: type_qualifier
	|       type_qualifier_list type_qualifier
                ;


parameter_type_list:
                // parameter_list ',' ELLIPSIS
	// |       
		parameter_list
                ;

parameter_list: parameter_declaration
	|       parameter_list ',' parameter_declaration
                ;

parameter_declaration:
                declaration_specifiers declarator
	|       declaration_specifiers abstract_declarator
	|       declaration_specifiers
                ;

identifier_list:IDENTIFIER
	|       identifier_list ',' IDENTIFIER
                ;

type_name:      specifier_qualifier_list abstract_declarator
	|       specifier_qualifier_list
                ;

abstract_declarator:
                pointer direct_abstract_declarator
	|       pointer
	|       direct_abstract_declarator
                ;

direct_abstract_declarator:
                '(' abstract_declarator ')'
	| '[' ']'
	| '[' '*' ']'
	| '['   type_qualifier_list assignment_expression ']'
	| '['   type_qualifier_list ']'
	| '['   assignment_expression ']'
	|       direct_abstract_declarator '[' ']'
	|       direct_abstract_declarator '[' '*' ']'
	|       direct_abstract_declarator '[' type_qualifier_list assignment_expression ']'
	|       direct_abstract_declarator '[' type_qualifier_list ']'
	|       direct_abstract_declarator '[' assignment_expression ']'
	| '(' ')'
	| '('   parameter_type_list ')'
	|       direct_abstract_declarator '(' ')'
	|       direct_abstract_declarator '(' parameter_type_list ')'
                ;

initializer: '{'initializer_list '}'
	| '{'   initializer_list ',' '}'
	|       assignment_expression
                ;

initializer_list:
                designation initializer
	|       initializer
	|       initializer_list ',' designation initializer
	|       initializer_list ',' initializer
                ;

designation:    designator_list '='
                ;

designator_list:designator
	|       designator_list designator
                ;

designator: '[' constant_expression ']'
	| '.'   IDENTIFIER
                ;

statement:      labeled_statement
	|       compound_statement
	|       expression_statement
	|       selection_statement
	|       iteration_statement
	|       jump_statement
                ;

labeled_statement:
                IDENTIFIER ':' statement
	|       CASE constant_expression ':' statement
	|       DEFAULT ':' statement
                ;

compound_statement:
                '{' '}'
	| '{'   block_item_list '}'
                ;

block_item_list:block_item
	|       block_item_list block_item
                ;

block_item:     declaration
	|       statement
                ;

expression_statement:
                ';'
	|       expression ';'
                ;

selection_statement:
                IF '(' expression ')' statement ELSE statement
	|       IF '(' expression ')' statement
	|       SWITCH '(' expression ')' statement
                ;

iteration_statement:
                FOR '(' expression_statement expression_statement ')' statement
	|       FOR '(' expression_statement expression_statement expression ')' statement
	|       FOR '(' declaration expression_statement ')' statement
	|       FOR '(' declaration expression_statement expression ')' statement
                ;

jump_statement: GOTO IDENTIFIER ';'
	|       CONTINUE ';'
	|       BREAK ';'
	|       RETURN ';'
	|       RETURN expression ';'
                ;

translation_unit:
                external_declaration
	|       translation_unit external_declaration
                ;

external_declaration:
                function_definition
	|       declaration
                ;

function_definition:
                declaration_specifiers declarator declaration_list compound_statement
	|       declaration_specifiers declarator compound_statement
                ;

declaration_list:
                declaration
	|       declaration_list declaration
                ;






// lines:
//                 /* empty */
//         |       lines line
//         |       lines ';'
//                 ;

// line:
//                 packageDeclaration
//         |       importDeclaration
//         |       functionDeclaration
//         ;

// importDeclaration:
//                 IMPORT STRING
//                 {
// 			impName := strings.TrimPrefix($2, "\"")
// 			impName = strings.TrimSuffix(impName, "\"")
// 			if imp, err := prgrm.GetModule(impName); err == nil {
// 				if mod, err := prgrm.GetCurrentModule(); err == nil {
// 					mod.AddImport(imp)
// 				}
// 			}
//                 }
//         ;

// selectorExpressionsAndStatements:
//                 /* empty */
//                 {
// 			$$ = false
//                 }
//         |       LBRACE expressionsAndStatements RBRACE
//                 {
// 			$$ = true
//                 }
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
//                 {
// 			mod := MakeModule($2)
// 			prgrm.AddModule(mod)
//                 }
//                 ;

// definitionAssignment:
//                 /* empty */
//                 {
// 			$$ = nil
//                 }
//         |       assignOperator argument
//                 {
// 			$$ = $2
//                 }
//                 ;

// functionParameters:
//                 LPAREN parameters RPAREN
//                 {
// 			$$ = $2
//                 }
//         |       LPAREN RPAREN
//                 {
// 			$$ = nil
//                 }
//         ;

// functionDeclaration:
//                 FUNC IDENT functionParameters
//                 {
// 			if mod, err := prgrm.GetCurrentModule(); err == nil {
// 				inFn = true
// 				fn := MakeFunction($2)
// 				mod.AddFunction(fn)
// 				if fn, err := mod.GetCurrentFunction(); err == nil {
// 					for _, inp := range $3 {
// 						fn.AddInput(inp)
// 					}
// 				}
// 			}
//                 }
//                 functionStatements
//         |       FUNC IDENT functionParameters functionParameters
//                 {
// 			if mod, err := prgrm.GetCurrentModule(); err == nil {
// 				inFn = true
// 				fn := MakeFunction($2)
// 				mod.AddFunction(fn)
// 				if fn, err := mod.GetCurrentFunction(); err == nil {

// 					//checking if there are duplicate parameters
// 					dups := append($3, $4...)
// 					for _, param := range dups {
// 						for _, dup := range dups {
// 							if param.Name == dup.Name && param != dup {
// 								panic(fmt.Sprintf("%s: %d: duplicate input and/or output parameters in function '%s'", fileName, yyS[yypt-0].line+1, $2))
// 							}
// 						}
// 					}
					
// 					for _, inp := range $3 {
// 						fn.AddInput(inp)
// 					}
// 					for _, out := range $4 {
// 						fn.AddOutput(out)
// 					}
// 				}
// 			}
//                 }
//                 functionStatements
//         ;

// parameter:
//                 IDENT BASICTYPE
//                 {
// 			$$ = MakeParameter($1, TypeNameToInt($2))
//                 }
//         |       IDENT IDENT
//                 {
// 			$$ = MakeParameter($1, TypeNameToInt($2))
//                 }
//         |       IDENT MULT IDENT
//                 {
// 			typ := TypeNameToInt($3)
// 			param := MakeParameter($1, typ)
// 			param.IsPointer = true
// 			$$ = param
//                 }
//         ;

// parameters:
//                 parameter
//                 {
// 			var params []*CXArgument
//                         params = append(params, $1)
//                         $$ = params
//                 }
//         |       parameters COMMA parameter
//                 {
// 			$1 = append($1, $3)
//                         $$ = $1
//                 }
//         ;

// functionStatements:
//                 LBRACE expressionsAndStatements RBRACE
//                 {
// 			inFn = false
//                 }
//         |       LBRACE RBRACE
//                 {
// 			inFn = false
//                 }
//         ;

// expressionsAndStatements:
//                 nonAssignExpression
//         |       assignExpression
//         |       statement
//         |       expressionsAndStatements nonAssignExpression
//         |       expressionsAndStatements assignExpression
//         |       expressionsAndStatements statement
//         ;


// assignExpression:
//                 VAR IDENT BASICTYPE definitionAssignment
//                 {
//                     $$ = nil
//                 }
//         |       VAR IDENT LBRACK RBRACK IDENT
//                 {
//                     $$ = nil
//                 }
//         |       argumentsList assignOperator argumentsList
//                 {
// 			argsL := $1
// 			argsR := $3

// 			if fn, err := prgrm.GetCurrentFunction(); err == nil {
// 				offset := 0
// 				for _, inp := range fn.Inputs {
// 					offset += inp.Size
// 				}
// 				for _, expr := range fn.Expressions {
// 					for _, inp := range expr.Inputs {
// 						offset += inp.Size
// 					}
// 					for _, out := range expr.Outputs {
// 						offset += out.Size
// 					}
// 				}

// 				expr := MakeExpression(op)
// 				if !replMode {
// 					expr.FileLine = yyS[yypt-0].line + 1
// 					expr.FileName = fileName
// 				}

				
				
// 				for i, argL := range argsL {
// 					if op, err := prgrm.GetFunction(idFn, CORE_MODULE); err == nil {
// 						expr := MakeExpression(op)
// 						if !replMode {
// 							expr.FileLine = yyS[yypt-0].line + 1
// 							expr.FileName = fileName
// 						}

// 						fn.AddExpression(expr)
// 						expr.AddTag(tag)
// 						tag = ""

// 						var outName string
// 						encoder.DeserializeRaw(*argL.Value, &outName)

// 						// checking if identifier was previously declared
// 						if outType, err := GetIdentType(outName, yyS[yypt-0].line + 1, fileName, prgrm); err == nil {
// 							if len(typeParts) > 1 {
// 								if outType != secondTyp {
// 									panic(fmt.Sprintf("%s: %d: identifier '%s' was previously declared as '%s'; cannot use type '%s' in assignment", fileName, yyS[yypt-0].line + 1, outName, outType, secondTyp))
// 								}
// 							} else if typeParts[0] == "ident" {
// 								var identName string
// 								encoder.DeserializeRaw(*argsR[i].Value, &identName)
// 								if rightTyp, err := GetIdentType(identName, yyS[yypt-0].line + 1, fileName, prgrm); err == nil {
// 									if outType != ptrs + rightTyp {
// 										panic(fmt.Sprintf("%s: %d: identifier '%s' was previously declared as '%s'; cannot use type '%s' in assignment", fileName, yyS[yypt-0].line + 1, outName, outType, ptrs + rightTyp))
// 									}
// 								}
// 							}
// 						}

// 						if len(typeParts) > 1 || typeParts[0] == "ident" {
// 							var identName string
// 							encoder.DeserializeRaw(*argsR[i].Value, &identName)
// 							identName = ptrs + identName
// 							sIdentName := encoder.Serialize(identName)
// 							arg := MakeArgument(&sIdentName, typ)
// 							expr.AddArgument(arg)
// 						} else {
// 							arg := MakeArgument(argsR[i].Value, typ)
// 							expr.AddArgument(arg)
// 						}
						
// 						expr.AddOutputName(outName)
// 					}
// 				}
// 			}
//                 }
//         ;

// nonAssignExpression:
//                 IDENT arguments
//                 {
// 			var modName string
// 			var fnName string
// 			var err error
// 			var isMethod bool
// 			//var receiverType string
// 			identParts := strings.Split($1, ".")
			
// 			if len(identParts) == 2 {
// 				mod, _ := prgrm.GetCurrentModule()
// 				if typ, err := GetIdentType(identParts[0], yyS[yypt-0].line + 1, fileName, prgrm); err == nil {
// 					// then it's a method call
// 					if IsStructInstance(typ, mod) {
// 						isMethod = true
// 						//receiverType = typ
// 						modName = mod.Name
// 						fnName = fmt.Sprintf("%s.%s", typ, identParts[1])
// 					}
// 				} else {
// 					// then it's a module
// 					modName = identParts[0]
// 					fnName = identParts[1]
// 				}
// 			} else {
// 				fnName = identParts[0]
// 				mod, e := prgrm.GetCurrentModule()
// 				modName = mod.Name
// 				err = e
// 			}

// 			found := false
// 			currModName := ""
// 			if mod, err := prgrm.GetCurrentModule(); err == nil {
// 				currModName = mod.Name
// 				for _, imp := range mod.Imports {
// 					if modName == imp.Name {
// 						found = true
// 						break
// 					}
// 				}
// 			}

// 			isModule := false
// 			if _, err := prgrm.GetModule(modName); err == nil {
// 				isModule = true
// 			}
			
// 			if !found && !IsNative(modName + "." + fnName) && modName != currModName && isModule {
// 				fmt.Printf("%s: %d: module '%s' was not imported or does not exist\n", fileName, yyS[yypt-0].line + 1, modName)
// 			} else {
// 				if err == nil {
// 					if fn, err := prgrm.GetCurrentFunction(); err == nil {
// 						if op, err := prgrm.GetFunction(fnName, modName); err == nil {
// 							expr := MakeExpression(op)
// 							if !replMode {
// 								expr.FileLine = yyS[yypt-0].line + 1
// 								expr.FileName = fileName
// 							}
// 							fn.AddExpression(expr)
// 							expr.AddTag(tag)
// 							tag = ""

// 							if isMethod {
// 								sIdent := encoder.Serialize(identParts[0])
// 								$2 = append([]*CXArgument{MakeArgument(&sIdent, "ident")}, $2...)
// 							}
							
// 							for _, arg := range $2 {
// 								typeParts := strings.Split(arg.Type, ".")

// 								arg.Type = typeParts[0]
// 								expr.AddArgument(arg)
// 							}

// 							lenOut := len(op.Outputs)
// 							outNames := make([]string, lenOut)
// 							args := make([]*CXArgument, lenOut)
							
// 							for i, out := range op.Outputs {
// 								outNames[i] = MakeGenSym(NON_ASSIGN_PREFIX)
// 								byteName := encoder.Serialize(outNames[i])
// 								args[i] = MakeArgument(&byteName, fmt.Sprintf("ident.%s", out.Type))

// 								expr.AddOutputName(outNames[i])
// 							}
							
// 							$$ = args
// 						} else {
// 							fmt.Printf("%s: %d: function '%s' not defined\n", fileName, yyS[yypt-0].line + 1, $1)
// 						}
// 					}
// 				}
// 			}
//                 }
//         ;

// beginFor:       FOR
//                 {
// 			if fn, err := prgrm.GetCurrentFunction(); err == nil {
// 				$<i>$ = len(fn.Expressions)
// 			}
//                 }
//                 ;

// conditionControl:
//                 nonAssignExpression
//                 {
// 			$$ = $1
//                 }
//         |       argument
//                 {
// 			$$ = []*CXArgument{$1}
//                 }
//         ;

// returnArg:
//                 ';'
//                 {
// 			$$ = nil
//                 }
//         |       argumentsList
//                 {
// 			$$ = $1
//                 }
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
//                 ';' conditionControl
//                 ';' forLoopAssignExpression //$<bool>9
//                 LBRACE expressionsAndStatements RBRACE
//         |       VAR IDENT IDENT
//         |       ';'
//         ;

// forLoopAssignExpression:
//                 {
// 			$<bool>$ = false
//                 }
//         |       assignExpression
//                 {
// 			$<bool>$ = true
//                 }
//         |       nonAssignExpression
//                 {
// 			$<bool>$ = true
//                 }
//         ;

// elseStatement:
//                 /* empty */
//                 {
//                     $<i>$ = 0
//                 }
//         |       ELSE
//                 {
//                     if mod, err := prgrm.GetCurrentModule(); err == nil {
//                         if fn, err := mod.GetCurrentFunction(); err == nil {
//                             if goToFn, err := prgrm.GetFunction("baseGoTo", mod.Name); err == nil {
// 				    expr := MakeExpression(goToFn)
// 				    if !replMode {
// 					    expr.FileLine = yyS[yypt-0].line + 1
// 					    expr.FileName = fileName
// 				    }
// 				    fn.AddExpression(expr)
//                             }
//                         }
//                     }
//                 }
//                 LBRACE
//                 {
// 			if fn, err := prgrm.GetCurrentFunction(); err == nil {
// 				$<i>$ = len(fn.Expressions)
// 			}
//                 }
//                 expressionsAndStatements RBRACE
//                 {
// 			if mod, err := prgrm.GetCurrentModule(); err == nil {
// 				if fn, err := mod.GetCurrentFunction(); err == nil {
// 					goToExpr := fn.Expressions[$<i>4 - 1]
					
// 					elseLines := encoder.Serialize(int32(0))
// 					thenLines := encoder.Serialize(int32(len(fn.Expressions) - $<i>4 + 1))

// 					alwaysTrue := encoder.Serialize(int32(1))

// 					goToExpr.AddArgument(MakeArgument(&alwaysTrue, "bool"))
// 					goToExpr.AddArgument(MakeArgument(&thenLines, "i32"))
// 					goToExpr.AddArgument(MakeArgument(&elseLines, "i32"))

// 					$<i>$ = len(fn.Expressions) - $<i>4
// 				}
// 			}
//                 }
//                 ;

// expressions:
//                 nonAssignExpression
//         |       assignExpression
//         |       expressions nonAssignExpression
//         |       expressions assignExpression
//         ;

// /*
//   Fix this, there has to be a way to compress these rules
// */
// argument:       LPAREN argument RPAREN
//                 {
// 			$$ = $2
//                 }
//         |       BYTENUM
//                 {
// 			val := encoder.Serialize($1)
//                         $$ = MakeArgument(&val, "byte")
//                 }
//         |       INT
//                 {
// 			val := encoder.Serialize($1)
//                         $$ = MakeArgument(&val, "i32")
//                 }
//         |       LONG
//                 {
// 			val := encoder.Serialize($1)
//                         $$ = MakeArgument(&val, "i64")
//                 }
//         |       FLOAT
//                 {
// 			val := encoder.Serialize($1)
// 			$$ = MakeArgument(&val, "f32")
//                 }
//         |       DOUBLE
//                 {
// 			val := encoder.Serialize($1)
// 			$$ = MakeArgument(&val, "f64")
//                 }
//         |       BOOLEAN
//                 {
// 			val := encoder.Serialize($1)
// 			$$ = MakeArgument(&val, "bool")
//                 }
//         |       STRING
//                 {
// 			str := strings.TrimPrefix($1, "\"")
//                         str = strings.TrimSuffix(str, "\"")

// 			val := encoder.Serialize(str)
			
//                         $$ = MakeArgument(&val, "str")
//                 }
//         |       IDENT
//                 {
// 			val := encoder.Serialize($1)
// 			$$ = MakeArgument(&val, "ident")
//                 }
//         |       IDENT LBRACK INT RBRACK
//                 {
// 			val := encoder.Serialize(fmt.Sprintf("%s[%d", $1, $3))
// 			$$ = MakeArgument(&val, "ident")
//                 }
//         ;

// arguments:
//                 LPAREN argumentsList RPAREN
//                 {
//                     $$ = $2
//                 }
//         |       LPAREN RPAREN
//                 {
//                     $$ = nil
//                 }
//         ;

// argumentsList:  argument
//                 {
// 			var args []*CXArgument
// 			args = append(args, $1)
// 			$$ = args
//                 }
//         |       nonAssignExpression
//                 {
// 			args := $1
// 			$$ = args
//                 }
//         |       argumentsList COMMA argument
//                 {
// 			$1 = append($1, $3)
// 			$$ = $1
//                 }
//         |       argumentsList COMMA nonAssignExpression
//                 {
// 			args := $3

// 			$1 = append($1, args...)
// 			$$ = $1
//                 }
//         ;

%%
