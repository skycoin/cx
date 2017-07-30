%{
package main
import (
  "strings"
  "github.com/skycoin/skycoin/src/cipher/encoder"
  . "github.com/skycoin/cx/src/base"
)

var cxt = MakeContext()

%}

%union {
    i32 int32
    f64 float64
    //    str string
    tok string
    fun *CXFunction
    params []*CXParameter
    param *CXParameter
    args []*CXArgument
    arg *CXArgument
    outNames []string
    
    cxt *CXContext
    mod *CXModule
}

%token  <i32>           INT
%token  <f64>           FLOAT

%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE IDENT KEYWORD TYP VAR COMMA COMMENT STRING

%type   <cxt>           cxtAdder
%type   <mod>           modAdder
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
        |       modAdder
        |       cxtAdder
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

expr:           outNames OP IDENT LPAREN args RPAREN
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
        |       expr
        |       fnAdder expr
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

modAdder:       def {
                    
                }
        |       fun {
                    
                }
        |       modAdder fun {}
        |       modAdder def {}
                
        ;

cxtAdder:           KEYWORD IDENT
                {
                    cxt.AddModule(MakeModule($2))
                        $$ = cxt
                        }
                ;

term : ';'
	{
	}
	| '\n'
	{
	}
	;%%
