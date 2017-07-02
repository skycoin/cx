package base

import (
	
)

// It should be cxt.Execute()

func (expr *cxExpression) Execute(state map[string]*cxDefinition) *cxArgument {
	fn := expr.Operator
	args := expr.Arguments
	fnName := fn.Name

	// checking if arguments are identifiers to extract and replace them
	for i := 0; i < len(args); i++ {
		if args[i].Typ.Name == "ident" {
			tmpDef := state[string(*args[i].Value)]
			args[i] = &cxArgument{Value: tmpDef.Value, Typ: tmpDef.Typ}
		}
	}

	// checking for native functions
	switch fnName {
	case "addI32":
		return addI32(args[0], args[1])
	default: // not native function
		
		// making a copy of the current state
		// to add/replace definitions for the new scope
		scope := make(map[string]*cxDefinition)

		for k, v := range state {
			scope[k] = v
		}

		for i, input := range fn.Inputs {
			def := &cxDefinition{
				Name: input.Name,
				Typ: input.Typ,
				Value: args[i].Value,
			}
			scope[input.Name] = def
		}

		stop := 0
		if fn.Output.Name != "" {
			stop = len(fn.Expressions)
		} else {
			stop = len(fn.Expressions) - 1
		}
		
		for i := 0; i < stop; i++ {
			fn.Expressions[i].Execute(scope)
		}

		if fn.Output.Name != "" {
			// if end-user didn't assign a value to that named output, we should raise an error
			return &cxArgument{
				Value: scope[fn.Output.Name].Value,
				Typ: scope[fn.Output.Name].Typ}
		} else {
			return fn.Expressions[len(fn.Expressions) - 1].Execute(scope)
		}
	}
}
