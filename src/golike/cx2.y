%{
package main
import (
  "strings"
  "fmt"
  "github.com/skycoin/skycoin/src/cipher/encoder"
  . "github.com/skycoin/cx/src/base"
)

var cxt = MakeContext()
var lineNo int = 0

%}

%union {
    i64 int64
    f64 float64
    
    tok string

    param *CXParameter
    params []*CXParameter

    arg *CXArgument
    args []*CXArgument

    expr *CXExpression
    exprs []*CXExpression

    fld *CXField
    flds []*CXField

    names []string
}

%token  <i64>           INT
%token  <f64>           FLOAT BOOLEAN
%token  <tok>           FUNC OP LPAREN RPAREN LBRACE RBRACE IDENT
                        /* Types */
                        STR I32 I64 F32 F64 BYTE BYTEA I32A I64A F32A F64A
                        ASSIGN CASSIGN GTHAN LTHAN LTEQ GTEQ
                        VAR COMMA COMMENT STRING PACKAGE IF ELSE WHILE TYPSTRUCT STRUCT UNEXPECTED

%type   <tok>           typeSpec
                        
%type   <param>         param
%type   <params>        params
%type   <arg>           arg defAssign
%type   <args>          args
/* %type   <fld>           fld */
/* %type   <flds>          flds */
%type   <expr>          expr arglessExpr argsExpr opExpr ifExpr whileExpr
%type   <exprs>         exprs elseExpr
%type   <names>         outNames

%%

lines:
        |       lines line
        ;

line:           defDecl
        |       structDecl
        |       packageDecl
        |       funcDecl
        |       UNEXPECTED
                {
                    yylex.Error(fmt.Sprintf("error in line #%d: unexpected character '%s'", lineNo, $1))
                }
        ;

assignOp:       ASSIGN
        |       CASSIGN
        ;

typeSpec:       I32 {$$ = $1}
        |       I64 {$$ = $1}
        |       F32 {$$ = $1}
        |       F64 {$$ = $1}
        |       BYTE {$$ = $1}
        |       BYTEA {$$ = $1}
        |       I32A {$$ = $1}
        |       I64A {$$ = $1}
        |       F32A {$$ = $1}
        |       F64A {$$ = $1}
        |       STR {$$ = $1}
        |       STRUCT IDENT {$$ = $2}
        ;

packageDecl:    PACKAGE IDENT
                {
                    cxt.AddModule(MakeModule($2))
                }
                ;

defAssign:
                {
                    $$ = nil
                }
        |       assignOp arg
                {
                    $$ = $2
                }
                ;

// second ident typ here
defDecl:        VAR IDENT typeSpec defAssign
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        var val *CXArgument;
                        if $4 == nil {
                                var zeroVal []byte
                                switch $3 {
                                        case "byte": zeroVal = []byte{byte(0)}
                                        case "i32": zeroVal = encoder.Serialize(int32(0))
                                        case "i64": zeroVal = encoder.Serialize(int64(0))
                                        case "f32": zeroVal = encoder.Serialize(float32(0))
                                        case "f64": zeroVal = encoder.Serialize(float64(0))
                                        case "[]byte": zeroVal = []byte{byte(0)}
                                        case "[]i32": zeroVal = encoder.Serialize([]int32{0})
                                        case "[]i64": zeroVal = encoder.Serialize([]int64{0})
                                        case "[]f32": zeroVal = encoder.Serialize([]float32{0})
                                        case "[]f64": zeroVal = encoder.Serialize([]float64{0})
                                        }
                                val = MakeArgument(&zeroVal, MakeType($3))
                            } else {
                            switch $3 {
                                case "byte":
                                var ds int64
                                encoder.DeserializeRaw(*$4.Value, &ds)
                                //new := encoder.Serialize(byte(ds))
                                new := []byte{byte(ds)}
                                val = MakeArgument(&new, MakeType("byte"))
                                case "i32":
                                var ds int64
                                encoder.DeserializeRaw(*$4.Value, &ds)
                                new := encoder.Serialize(int32(ds))
                                val = MakeArgument(&new, MakeType("i32"))
                                case "i64": /* stays the same */ 
                                case "f32":
                                var ds float64
                                encoder.DeserializeRaw(*$4.Value, &ds)
                                new := encoder.Serialize(float32(ds))
                                val = MakeArgument(&new, MakeType("f32"))
                                case "f64":
                                val = $4
                                case "[]byte":
                                /* var ds []int64 */
                                /* encoder.DeserializeRaw(*$4.Value, &ds) */
                                /* new := make([]byte, len(ds)) */
                                /* for i, val := range ds { */
                                /*         new[i] = byte(val) */
                                /*     } */
                                /* val = MakeArgument(&new, MakeType("[]byte")) */
                                val = $4
                                case "[]i32":
                                /* var ds []int64 */
                                /* encoder.DeserializeRaw(*$4.Value, &ds) */
                                /* new := make([]int32, len(ds)) */
                                /* for i, val := range ds { */
                                /*         new[i] = int32(val) */
                                /*     } */
                                /* sNew := encoder.Serialize(new) */
                                /* val = MakeArgument(&sNew, MakeType("[]i32")) */
                                val = $4
                                case "[]i64":
                                val = $4
                                case "[]f32":
                                /* var ds []float64 */
                                /* encoder.DeserializeRaw(*$4.Value, &ds) */
                                /* new := make([]float32, len(ds)) */
                                /* for i, val := range ds { */
                                /*         new[i] = float32(val) */
                                /*     } */
                                /* sNew := encoder.Serialize(new) */
                                /* val = MakeArgument(&sNew, MakeType("[]f32")) */
                                val = $4
                                case "[]f64": /* stays the same */
                                val = $4
                            }
                            //val = $4
                        }
                
                        mod.AddDefinition(MakeDefinition($2, val.Value, MakeType($3)))
                    }
                }
        ;

structDecl:     TYPSTRUCT IDENT
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        strct := MakeStruct($2)
                        mod.AddStruct(strct)
                    }
                }
                STRUCT LBRACE params RBRACE
                {
                    if strct, err := cxt.GetCurrentStruct(); err == nil {
                        for _, fld := range $6 {
                            fldFromParam := MakeField(fld.Name, fld.Typ)
                            strct.AddField(fldFromParam)
                            }
                    }
                }
        ;

// second ident typ here
/* fld:          IDENT TYP */
/*                 { */
/*                     $$ = MakeField($1, MakeType($2)) */
/*                 } */
/*         ; */

/* flds: */
/*                 { */
/*                     $$ = nil */
/*                 } */
/*         |       fld */
/*                 { */
/*                     var flds []*CXField */
/*                         flds = append(flds, $1) */
/*                         $$ = flds */
/*                 } */
/*         |       flds fld */
/*                 { */
/*                     $1 = append($1, $2) */
/*                         $$ = $1 */
/*                 } */
/*         ; */

funcDecl:       FUNC IDENT LPAREN params RPAREN LPAREN params RPAREN
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
                LBRACE exprs RBRACE
        ;



/* IF THEN */
/* IF THEN ELSE */

/* ELSECLAUSE */
/*         |       ELSE */


// second ident typ here
param:          IDENT typeSpec
                {
                    $$ = MakeParameter($1, MakeType($2))
                }
        ;

params:
                {
                    $$ = nil
                }
        |       param
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
        ;



/* statement */
/*   : matched */
/*   | unmatched */
/* ; */

/* matched */
/*   : IF '(' expression ')' matched ELSE matched */
/*   | other_statement */
/* ; */

/* unmatched */
/*   : IF '(' expression ')' statement */
/*   | IF '(' expression ')' matched ELSE unmatched */
/* ; */

/* other_statement */
/*   : iteration_statement */
/*   | compound_statement */
/*   | ... other alternatives ... */
/* ; */


argsExpr:       outNames assignOp IDENT LPAREN args RPAREN
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

arglessExpr:    IDENT LPAREN args RPAREN
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
        |       VAR IDENT typeSpec defAssign
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := cxt.GetCurrentFunction(); err == nil {
                            if $4 == nil {
                                    if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
                                    expr := MakeExpression(op)
                                    fn.AddExpression(expr)
                                    expr.AddOutputName($2)
                                    typ := []byte($3)
                                    arg := MakeArgument(&typ, MakeType("str"))
                                    expr.AddArgument(arg)

                                    if strct, err := cxt.GetStruct($3, mod.Name); err == nil {
                                            for _, fld := range strct.Fields {
                                                expr := MakeExpression(op)
                                                fn.AddExpression(expr)
                                                expr.AddOutputName(fmt.Sprintf("%s.%s", $2, fld.Name))
                                                typ := []byte(fld.Typ.Name)
                                                arg := MakeArgument(&typ, MakeType("str"))
                                                expr.AddArgument(arg)
                                                }
                                        }
                                    }
                                } else {
                                switch $3 {
                                    case "byte":
                                    var ds int64
                                    encoder.DeserializeRaw(*$4.Value, &ds)
                                    new := []byte{byte(ds)}
                                    val := MakeArgument(&new, MakeType("byte"))
                
                                    if op, err := cxt.GetFunction("idByte", mod.Name); err == nil {
                                    expr := MakeExpression(op)
                                    fn.AddExpression(expr)
                                    expr.AddOutputName($2)
                                    expr.AddArgument(val)
                                    }
                                    case "i32":
                                    var ds int64
                                    encoder.DeserializeRaw(*$4.Value, &ds)
                                    new := encoder.Serialize(int32(ds))
                                    val := MakeArgument(&new, MakeType("i32"))

                                    if op, err := cxt.GetFunction("idI32", mod.Name); err == nil {
                                    expr := MakeExpression(op)
                                    fn.AddExpression(expr)
                                    expr.AddOutputName($2)
                                    expr.AddArgument(val)
                                    }
                                    case "f32":
                                    var ds float64
                                    encoder.DeserializeRaw(*$4.Value, &ds)
                                    new := encoder.Serialize(float32(ds))
                                    val := MakeArgument(&new, MakeType("f32"))

                                    if op, err := cxt.GetFunction("idF32", mod.Name); err == nil {
                                    expr := MakeExpression(op)
                                    fn.AddExpression(expr)
                                    expr.AddOutputName($2)
                                    expr.AddArgument(val)
                                    }
                                    default:
                                    val := $4
                                    var getFn string
                                    switch $3 {
                                        case "i64": getFn = "idI64"
                                        case "f64": getFn = "idF64"
                                        case "[]byte": getFn = "idByteA"
                                        case "[]i32": getFn = "idI32A"
                                        case "[]i64": getFn = "idI64A"
                                        case "[]f32": getFn = "idF32A"
                                        case "[]f64": getFn = "idF64A"
                                    }

                                    if op, err := cxt.GetFunction(getFn, mod.Name); err == nil {
                                    expr := MakeExpression(op)
                                    fn.AddExpression(expr)
                                    expr.AddOutputName($2)
                                    expr.AddArgument(val)
                                    }
                                }
                            
                            }
                        }
                    }
                }
        ;

whileExpr:      WHILE
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := mod.GetCurrentFunction(); err == nil {
                            if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
                            expr := MakeExpression(goToFn)
                            fn.AddExpression(expr)
                            }
                        }
                    }
                }
                arg LBRACE exprs RBRACE
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := mod.GetCurrentFunction(); err == nil {
                            
                        goToExpr := fn.Expressions[len(fn.Expressions) - len($5) - 1]
                            
                        elseLines := encoder.Serialize(int32(len($5) + 2))
                        thenLines := encoder.Serialize(int32(1))
                
                        predVal := *$3.Value
                        
                        goToExpr.AddArgument(MakeArgument(&predVal, $3.Typ))
                        goToExpr.AddArgument(MakeArgument(&thenLines, MakeType("i32")))
                        goToExpr.AddArgument(MakeArgument(&elseLines, MakeType("i32")))

                        if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
                            goToExpr := MakeExpression(goToFn)
                            fn.AddExpression(goToExpr)

                            elseLines := encoder.Serialize(int32(1))
                            thenLines := encoder.Serialize(int32(-len($5)))

                            goToExpr.AddArgument(MakeArgument(&predVal, $3.Typ))
                            goToExpr.AddArgument(MakeArgument(&thenLines, MakeType("i32")))
                            goToExpr.AddArgument(MakeArgument(&elseLines, MakeType("i32")))
                            }
                        
                        }
                    }
                }
                ;

ifExpr:         IF
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := mod.GetCurrentFunction(); err == nil {
                            if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
                            expr := MakeExpression(goToFn)
                            fn.AddExpression(expr)
                            }
                        }
                    }
                }
                arg LBRACE exprs RBRACE elseExpr
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := mod.GetCurrentFunction(); err == nil {
                            
                            var goToExpr *CXExpression
                            if len($7) > 0 {
                                    goToExpr = fn.Expressions[len(fn.Expressions) - 2 - len($5) - len($7)]
                                } else {
                                goToExpr = fn.Expressions[len(fn.Expressions) - 2 - len($5)]
                            }
                
                        elseLines := encoder.Serialize(int32(len($5) + 2))
                            thenLines := encoder.Serialize(int32(1))
                
                            predVal := *$3.Value

                            goToExpr.AddArgument(MakeArgument(&predVal, $3.Typ))
                            goToExpr.AddArgument(MakeArgument(&thenLines, MakeType("i32")))
                            goToExpr.AddArgument(MakeArgument(&elseLines, MakeType("i32")))
                        }
                    }
                }
                ;

elseExpr:
                {
                    exprs := make([]*CXExpression, 0)
                        $$ = exprs
                }
        |       ELSE
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := mod.GetCurrentFunction(); err == nil {
                            if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
                            expr := MakeExpression(goToFn)
                            fn.AddExpression(expr)
                            }
                        }
                    }
                }
                LBRACE exprs RBRACE
                {
                    if mod, err := cxt.GetCurrentModule(); err == nil {
                        if fn, err := mod.GetCurrentFunction(); err == nil {
                            
                            goToExpr := fn.Expressions[len(fn.Expressions) - 1 - len($4)]
                            
                            elseLines := encoder.Serialize(int32(0))
                            thenLines := encoder.Serialize(int32(len($4) + 1))

                            predVal := []byte{1}

                            goToExpr.AddArgument(MakeArgument(&predVal, MakeType("byte")))
                            goToExpr.AddArgument(MakeArgument(&thenLines, MakeType("i32")))
                            goToExpr.AddArgument(MakeArgument(&elseLines, MakeType("i32")))
                        }
                    }
                    
                    $$ = $4
                }
                ;

opExpr:         arglessExpr
        |       argsExpr
        ;

expr:           ifExpr
        |       whileExpr
        |       opExpr
        ;

exprs: 
                {
                exprs := make([]*CXExpression, 0)
                        $$ = exprs
                }
        |       expr
                {
                    exprs := make([]*CXExpression, 1)
                        exprs[0] = $1
                        $$ = exprs
                        
                }
        |       exprs expr
                {
                    $1 = append($1, $2)
                        $$ = $1
                }
        ;

arg:            INT
                {
                val := encoder.SerializeAtomic($1)
                    $$ = MakeArgument(&val, MakeType("i64"))
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
                // first ident typ here
        |       typeSpec LBRACE args RBRACE
                {
                    switch $1 {
                            case "[]byte":
                                vals := make([]byte, len($3))
                                    for i, arg := range $3 {
                                            var val int64
                                            encoder.DeserializeRaw(*arg.Value, &val)
                                            vals[i] = byte(val)
                                        }
                sVal := encoder.Serialize(vals)
                                    $$ = MakeArgument(&sVal, MakeType("[]i32"))
                                    case "[]i32":
                                vals := make([]int32, len($3))
                                    for i, arg := range $3 {
                                            var val int64
                                            encoder.DeserializeRaw(*arg.Value, &val)
                                            vals[i] = int32(val)
                                        }
                sVal := encoder.Serialize(vals)
                                    $$ = MakeArgument(&sVal, MakeType("[]i32"))
                                    case "[]i64":
                                vals := make([]int64, len($3))
                                    for i, arg := range $3 {
                                            var val int64
                                            encoder.DeserializeRaw(*arg.Value, &val)
                                            vals[i] = val
                                        }
                        sVal := encoder.Serialize(vals)
                                    $$ = MakeArgument(&sVal, MakeType("[]i64"))
                                    case "[]f32":
                                vals := make([]float32, len($3))
                                    for i, arg := range $3 {
                                            var val float64
                                            encoder.DeserializeRaw(*arg.Value, &val)
                                            vals[i] = float32(val)
                                        }
                        sVal := encoder.Serialize(vals)
                                    $$ = MakeArgument(&sVal, MakeType("[]f32"))
                                    case "[]f64":
                                vals := make([]float64, len($3))
                                    for i, arg := range $3 {
                                            var val float64
                                            encoder.DeserializeRaw(*arg.Value, &val)
                                            vals[i] = val
                                        }
                        sVal := encoder.Serialize(vals)
                                    $$ = MakeArgument(&sVal, MakeType("[]f64"))
                                    }
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

/*
  The goTo is now ready, now we need to build the expressions
*/

%%
