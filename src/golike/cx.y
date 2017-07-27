%{
package main
import (
  "fmt"
  . "github.com/skycoin/cx/src/base"
)

var cxt = MakeContext()

%}

%union {
    i int32
    f float64
    k string
    w string
    cxt CXContext
}

%token  <i>             INT
%token  <f>             FLOAT
%token  <k>             KEYWORD
%token  <w>             WORD

%type  <cxt>           cxtExpr

%left '+' '-'
%left '*' '/'

%%

input:
        |       input line
;

line:     '\n'
        |       cxtExpr '\n'      { fmt.Println($1.n); }
;

cxtExpr:           KEYWORD WORD
                {
                    cxt.AddModule(MakeModule($2))
                        fmt.Println()
                        cxt.PrintProgram(false)
                        return cxt
                }
;

%%
