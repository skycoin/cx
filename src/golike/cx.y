%{
package main
import (
  "strings"
  "fmt"
  "github.com/skycoin/skycoin/src/cipher/encoder"
  . "github.com/skycoin/cx/src/base"
)

var cxt = MakeContext()
var lineNo int = 1

%}

%union {
    i32 int32
    f64 float64
    tok string
    fun *CXFunction
    params []*CXParameter
    param *CXParameter
    args []*CXArgument
    arg *CXArgument
    outNames []string
    fnAdder []*CXExpression
    expr *CXExpression

    cxt *CXContext
    mod *CXModule
}

%token  <i32>           INT
%token  <f64>           FLOAT

%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE IDENT TYP VAR COMMA COMMENT STRING PACKAGE IF ELSE

%type   <cxt>           cxtAdder
%type   <mod>           modAdder
%type   <fnAdder>       fnAdder
%type   <expr>          expr ifStat
%type   <param>         param
%type   <params>        params
%type   <arg>           arg
%type   <args>          args
%type   <outNames>      outNames

%right '='
%left '+' '-'
%left '*' '/'
%nonassoc ','

%%

input:          
        |       input line
;

line:           term
        |       cxtAdder
        |       modAdder
                ;

cxtAdder:           PACKAGE IDENT
                {
                    fmt.Printf("")
                    cxt.AddModule(MakeModule($2))
                        $$ = cxt
                        }
                ;

param:          IDENT TYP
                {
                    $$ = MakeParameter($1, MakeType($2))
                }
        ;

params:         param
                {
                    var params []*CXParameter
                        params = append(params, $1)
                        $$ = params
                        }
        |       params COMMA param
                {
                    $1 = append($1, $3)
                        $$ = $1
                        }
        |
                {
                    $$ = nil
                }
        ;

arg:            INT
                {
                val := encoder.SerializeAtomic($1)
                    $$ = MakeArgument(&val, MakeType("i32"))
                }
        |       FLOAT
                {
                    val := encoder.Serialize($1)
                    $$ = MakeArgument(&val, MakeType("f64"))
                }
        |       STRING
                {
                    var str string
                        str = strings.TrimPrefix($1, "\"")
                        str = strings.TrimSuffix(str, "\"")
                        
                        val := []byte(str)
                        $$ = MakeArgument(&val, MakeType("str"))
                        }
        |       IDENT
                {
                val := []byte($1)
                        $$ = MakeArgument(&val, MakeType("ident"))
                }
        |       TYP LBRACE args RBRACE
                {
                    var val []byte
                        for _, arg := range $3 {
                                val = append(val, *arg.Value...)
                            }
                    $$ = MakeArgument(&val, MakeType($1))
                }
        ;

args:
                {
                    
                }
        |       arg
                {
                    var args []*CXArgument
                        args = append(args, $1)
                        $$ = args
                }
        |       args COMMA arg
                {
                    $1 = append($1, $3)
                        $$ = $1
                }
        ;

outNames:       IDENT
                {
                outNames := make([]string, 1)
                        outNames[0] = $1
                        $$ = outNames
                }
        |       outNames COMMA IDENT
                {
                    $1 = append($1, $3)
                        $$ = $1
                }
        ;

expr:           IDENT LPAREN args RPAREN
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := cxt.GetCurrentFunction(); err == nil {
                            if op, err := cxt.GetFunction($1, mod.Name); err == nil {
                            expr := MakeExpression(op)

                            fn.AddExpression(expr)

                            if expr, err := fn.GetCurrentExpression(); err == nil {
                                    for _, arg := range $3 {
                                            expr.AddArgument(arg)
                                        }
                                    $$ = expr
                                }
                
                            } else {
                                panic(err)
                            }
                        }
                
                    } else {
                        panic(err)
                            }
                }
        |       outNames OP IDENT LPAREN args RPAREN
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := cxt.GetCurrentFunction(); err == nil {
                            if op, err := cxt.GetFunction($3, mod.Name); err == nil {
                            expr := MakeExpression(op)

                            for _, outName := range $1 {
                                    expr.AddOutputName(outName)
                                }
                
                            fn.AddExpression(expr)

                            if expr, err := fn.GetCurrentExpression(); err == nil {
                                for _, arg := range $5 {
                                        expr.AddArgument(arg)
                                    }
                                $$ = expr
                            }
                
                            } else {
                                panic(err)
                            }
                        }
                
                    } else {
                        panic(err)
                    }
                }
        ;

fnAdder:        
                {
                exprs := make([]*CXExpression, 0)
                        $$ = exprs
                }
        |
                ifStat
                {
                
                }
        |       expr
                {
                    exprs := make([]*CXExpression, 1)
                        exprs[0] = $1
                        $$ = exprs
                        
                }
        |       fnAdder expr
                {
                    $1 = append($1, $2)
                        $$ = $1
                }
        ;

elseStat:       ELSE LBRACE fnAdder RBRACE
                {
                fmt.Println("else")
                fmt.Println($3)

                /* if mod, err := cxt.GetCurrentModule(); err == nil { */
                /*     if fn, err := mod.GetCurrentFunction(); err == nil { */
                /* fn.AddExpression() */
                /* } */
                /* } */
                
                }
                ;
/*
  The goTo is now ready, now we need to build the expressions
*/
ifStat:         IF arg LBRACE fnAdder RBRACE
                {
                fmt.Println("if/then")
                fmt.Println($4)
                }
        |       IF arg LBRACE fnAdder RBRACE elseStat
                {
                fmt.Println("if/then with else")
                fmt.Println($4)
                }
        ;

def:           VAR IDENT TYP OP
                arg
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                    mod.AddDefinition(MakeDefinition($2, $5.Value, MakeType($3)))
                    }
                }
        ;

fun:            FUNC IDENT LPAREN params RPAREN LPAREN params RPAREN
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        mod.AddFunction(MakeFunction($2))
                        if fn, err := mod.GetCurrentFunction(); err == nil {
                            for _, inp := range $4 {
                                    fn.AddInput(inp)
                                }
                            for _, out := range $7 {
                                    fn.AddOutput(out)
                                }
                        }
                    }
                }
                LBRACE fnAdder RBRACE
        ;

modAdder:       def
                {
                    
                }
        |       fun
                {
                    
                }
        |       modAdder fun
                {
                
                }
        |       modAdder def
                {
                
                }
                
        ;

// cxtAdder should be in line

term : ';'
	{
	}
	| '\n'
	{
	}
	;

%%
