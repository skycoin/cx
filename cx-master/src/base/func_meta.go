package base

import (
	"fmt"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func rem_arg (tag *CXArgument, caller *CXFunction) error {
	if err := checkType("rem.arg", "str", tag); err == nil {
		var tg string
		encoder.DeserializeRaw(*tag.Value, &tg)

		for _, expr := range caller.Expressions {
			if expr.Tag == tg {
				expr.RemoveArgument()
				return nil
			}
		}
		return errors.New(fmt.Sprintf("rem.arg: no expression with label '%s' was found", tg))
	} else {
		return nil
	}
}

func add_arg (tag *CXArgument, ident *CXArgument, caller *CXFunction) error {
	if err := checkTwoTypes("add.arg", "str", "str", tag, ident); err == nil {
		var tg string
		encoder.DeserializeRaw(*tag.Value, &tg)
		
		for _, expr := range caller.Expressions {
			if expr.Tag == tg {
				expr.AddArgument(MakeArgument(ident.Value, "ident"))
				return nil
			}
		}
		return errors.New(fmt.Sprintf("add.arg: no expression with tag '%s' was found", tg))
	} else {
		return err
	}
}

func add_expr (tag *CXArgument, fnName *CXArgument, call *CXCall) error {
	if err := checkType("add.expr", "str", fnName); err == nil {
		mod := call.Package

		var opName string
		var tg string
		encoder.DeserializeRaw(*fnName.Value, &opName)
		encoder.DeserializeRaw(*tag.Value, &tg)
		
		if fn, err := mod.Program.GetFunction(opName, mod.Name); err == nil {
			expr := MakeExpression(fn)
			expr.AddTag(tg)
			
			call.Operator.AddExpression(expr)
			return nil
		} else {
			return nil
		}
	} else {
		return err
	}
}

func rem_expr (tag *CXArgument, caller *CXFunction) error {
	var tg string
	encoder.DeserializeRaw(*tag.Value, &tg)
	if err := checkType("remExpr", "str", tag); err == nil {
		for i, expr := range caller.Expressions {
			if expr.Tag == tg {
				caller.RemoveExpression(i)
				return nil
			}
		}
	} else {
		return err
	}
	return errors.New(fmt.Sprintf("remExpr: no expression with tag '%s' was found", tg))
}
