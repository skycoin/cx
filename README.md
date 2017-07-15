# CX Programming Language

This is a prototype version of the CX Programming Language. At the
moment, in order to create a program, the programmer needs to use the base
language directly.

## Basic Usage

```
package main

import (
    "fmt"
    "github.com/skycoin/skycoin/src/cipher/encoder"
    . "github.com/skycoin/cx/src/base"
)

func main() {

}

```

Start by creating a reference to the only currently supported primitive
type *i32*.

```
    i32 := MakeType("i32")
```

We will use the variable *i32* when creating objects that require a
type, like *definitions* and *fields* in structures. When executing a
function which uses definitions of type i32, the base language will
automatically recognize this primitive type.

Arguments that represent identifiers need to be of type *ident*. These
identifiers will point to a value of certain type (like i32). This
behaviour is similar to how Lisp handles symbols, and how symbols
point to other data. When using languages constructed using the CX
base language, this process should be transparent to the user. For the
moment, let's create the type *ident*:

```
    ident := MakeType("ident")
```

The last reference to a primitive that we need to make is to the
*addI32* function:

```
    addI32 := MakeFunction("addI32")
```

which simply adds two arguments of type i32.

A function is created using a *maker*, which is, in this case,
*MakeFunction*:

```
    double := MakeFunction("double")
```

This function will take a number, will double it and then will return
the result. Let's start with adding the input and the output:

```
    double := MakeFunction("double").
        AddInput(MakeParameter("num", i32)).
        AddOutput(MakeParameter("", i32))
```

*Adders* are used to add objects to another object. In this case,
*AddInput* and *AddOutput* were used to add an input parameter and an
output parameter to the double function, respectively. Both of these
adders require a *parameter* object, which is created by using the maker
*MakeParameter*. The input parameter has the name "num" and is of type
i32. The output is unnamed and is of type i32. In the CX base
language, one can give the output a name or not. The output of a
function will correspond to the name of the output parameter in the former case,
or will be the value of the last expression in the function body in
the latter. As this function will only require a single expression, an
unnamed output is more convenient. Let's now add this expression:

```
    double := MakeFunction("double").
		  AddInput(MakeParameter("num", i32)).
		  AddOutput(MakeParameter("", i32)).
		
		  AddExpression(MakeExpression(addI32).
		                AddArgument(MakeArgument(MakeValue("num"), ident)).
	                    AddArgument(MakeArgument(MakeValue("num"), ident)))
```

The maker *MakeExpression* requires an operator to make an expression,
which is addI32 in this case. Without any arguments, this maker would
create this s-expression: `(+)`. After adding twice the *num*
identifier to the expression, the s-expression would look like this:
`(+ num num)`.

The double function could be represented by the following
s-expression:

```
(defun double (num i32)
    (+ num num))
```

If we want to call this function, and actually double a number with
it, we can make another expression which uses *double* as its
operator. But first, let's create a variable which holds the number
25:

```
num := encoder.SerializeAtomic(int32(25))
```

As values in CX are serialized to byte arrays, we used Skycoin's
`cipher/encoder` package to serialize the number. Now, to double this
number and print the value:

```
fmt.Println(
		MakeExpression(double).
			AddArgument(MakeArgument(&num, i32)).
			
			Execute(mod.Definitions).Value)
```

In the next versions, a program should be executed as this
`cxt.Execute()` where `cxt` is the context holding references to
everything in a program. The program would run by executing a `main()`
function or by continuing executing the function and expression
in which it was paused at. For now, we need to create an expression
which uses *double* and the number we created as its argument, and
then call its `Execute` method. The execute method requires a
definitions array to act as its state.

By compiling and running the code above, this output should be
printed:

```
&[50 0 0 0]
```

which is a seralized 50 int32.

## Affordances

An affordance of an object represents something that we can do to that
object. At the moment, CX can apply some basic affordances on its
objects, like *adders* and *selectors*. Adders were used in the
previous section, and these basically add objects to other objects
(like adding an input parameter to a function). Selectors change the
current object the program is acting on: for example,
`SelectModule` will change a context to work on a desired module,
which means that any call to `AddFunction`, `AddStruct` or
`AddDefinition` will add objects to this module.

Let's create a "root" object or context and call its `GetAffordances`
method:

```
    cxt = MakeContext()
	PrintAffordances(cxt.GetAffordances())
```

`PrintAffordances` is a function used for debugging affordances; it
just prints an affordance's description to the console. The code above
should give the following output:

```
AddModule mod26
```

This is telling us that the only affordance over our empty context is
to add a module. If we decide to apply this affordance, the program
will generate a unique, but not very descriptive, identifier (mod26 in
this case). Let's apply this affordance:

```
    cxt.GetAffordances()[0].ApplyAffordance()
```

Note that we had to state that we desire to apply affordance
`0`. Other functions should provide the logic to follow in order to
determine what affordance to apply. A `FilterAffordances` function is
provided, which receives a series of string keywords to filter the
descriptions of the provided affordances. For example:

```
FilterAffordances(cxt.GetAffordances(),
    "AddExpression", "write")
```

this filter will return all the affordances which involve adding an
expression with operators containing the "write" keyword.

After applying the affordance obtained above, if we make another call
to `GetAffordances`, we should get an output similar to this:

```
AddModule mod27
SelectModule mod26
```

It is now telling us that we can either add another module or select
the previously created module. Let's now run the following code:

```
    PrintAffordances(cxt.GetCurrentModule().GetAffordances())
```

We are using a *getter* called `GetCurrentModule` on our context to
obtain the active module, and then we get its affordances and print
them. The output should be similar to this:

```
AddDefinition def27 i32
AddFunction fn28
AddStruct strct29
```

This tells us that we can add a new definition of type i32, a function
or a struct to the module. If we had more types (structs and more
primitive types), the list of *AddDefinition*s would be longer (one
for each available type in the current module, and for each imported
module).

If we want to know the structure of a program (context), we can use
the function `PrintProgram`, which is a function that has been used
only for debugging purposes. By running `cxt.PrintProgram(false)`, we
should get something like this:

```
Context
0.- Module: mod0
1.- Module: Math
	Definitions
		0.- Definition: def2 i32
		1.- Definition: def3 Complex
		2.- Definition: hugs i32
		3.- Definition: items StdLib.Stream
		4.- Definition: views i64
	Structs
		0.- Struct: Complex
		1.- Struct: List
			0.- Field: fld18 i32
			1.- Field: fld23 List
	Functions
		0.- Function: fn5 (in16 StdLib.Stream, in17 StdLib.Stream, num1 i32, num2 i32) sum i32
		1.- Function: fn10 (items StdLib.Stream) 
			0.- Expression: fn5(items StdLib.Stream, items StdLib.Stream, hugs i32, hugs i32)
2.- Module: StdLib
	Structs
		0.- Struct: Stream
	Functions
		0.- Function: fn14 () 
```

The `false` argument tells the function to not print the available
affordances for each object. Here is an example of this case:

```
Context
 * 0.- AddModule mod26
 * 1.- SelectModule mod0
 * 2.- SelectModule Math
 * 3.- SelectModule StdLib
 * 4.- SelectFunction fn5
 * 5.- SelectFunction fn10
 * 6.- SelectStruct List
 * 7.- SelectStruct Complex
 * 8.- SelectExpression Line # 0
0.- Module: mod0
	 * 0.- AddDefinition def27 i32
	 * 1.- AddImport Math
	 * 2.- AddImport StdLib
	 * 3.- AddFunction fn28
	 * 4.- AddStruct strct29
1.- Module: Math
	 * 0.- AddDefinition def30 i32
	 * 1.- AddDefinition def31 List
	 * 2.- AddDefinition def32 Complex
	 * 3.- AddDefinition def33 StdLib.Stream
	 * 4.- AddImport mod0
	 * 5.- AddImport StdLib
	 * 6.- AddFunction fn34
	 * 7.- AddStruct strct35
	Definitions
		0.- Definition: def2 i32
		1.- Definition: def3 List
		2.- Definition: hugs i32
		3.- Definition: items StdLib.Stream
		4.- Definition: views i64
	Structs
		0.- Struct: List
		 * 0.- AddField fld36 i32
		 * 1.- AddField fld37 List
		 * 2.- AddField fld38 Complex
		 * 3.- AddField fld39 StdLib.Stream
			0.- Field: fld18 i32
			1.- Field: fld23 List
		1.- Struct: Complex
		 * 0.- AddField fld40 i32
		 * 1.- AddField fld41 List
		 * 2.- AddField fld42 Complex
		 * 3.- AddField fld43 StdLib.Stream
	Functions
		0.- Function: fn5 (in16 StdLib.Stream, in17 StdLib.Stream, num1 i32, num2 i32) sum i32
		 * 0.- AddInput i32
		 * 1.- AddInput List
		 * 2.- AddInput Complex
		 * 3.- AddInput StdLib.Stream
		 * 4.- AddExpression fn10
		 * 5.- AddExpression StdLib.fn14
		1.- Function: fn10 (items StdLib.Stream) 
		 * 0.- AddInput i32
		 * 1.- AddInput List
		 * 2.- AddInput Complex
		 * 3.- AddInput StdLib.Stream
		 * 4.- AddOutput i32
		 * 5.- AddOutput List
		 * 6.- AddOutput Complex
		 * 7.- AddOutput StdLib.Stream
		 * 8.- AddExpression fn5
		 * 9.- AddExpression StdLib.fn14
			0.- Expression: fn5(items StdLib.Stream, items StdLib.Stream, def2 i32, hugs i32)
2.- Module: StdLib
	 * 0.- AddDefinition def44 i32
	 * 1.- AddDefinition def45 Stream
	 * 2.- AddImport mod0
	 * 3.- AddImport Math
	 * 4.- AddFunction fn46
	 * 5.- AddStruct strct47
	Structs
		0.- Struct: Stream
		 * 0.- AddField fld48 i32
		 * 1.- AddField fld49 Stream
	Functions
		0.- Function: fn14 () 
		 * 0.- AddInput i32
		 * 1.- AddInput Stream
		 * 2.- AddOutput i32
		 * 3.- AddOutput Stream
```

## Randomly Generating a Program

A function which continually applies affordances to an empty program
is provided. This function, `RandomProgram`, receives a number that
represents how many affordances will be applied, and returns the
generated context. For example, by running

```
    cxt := RandomProgram(100)
    cxt.PrintProgram(false)
```

we get an output similar to this:

```
0.- Module: mod26
	Definitions
		0.- Definition: def27 i32
		1.- Definition: def33 i32
		2.- Definition: def39 i32
	Functions
		0.- Function: fn31 (in36 i32, in38 i32, in51 i32) 
1.- Module: mod42
	Definitions
		0.- Definition: def113 mod42.strct89
		1.- Definition: def119 strct89
		2.- Definition: def135 mod42.strct89
		3.- Definition: def142 mod42.strct48
		4.- Definition: def63 strct48
		5.- Definition: def75 i32
	Structs
		0.- Struct: strct48
			0.- Field: fld54 strct48
			1.- Field: fld57 strct48
			2.- Field: fld59 i32
		1.- Struct: strct89
			0.- Field: fld94 mod42.strct89
			1.- Field: fld99 mod42.strct89
			2.- Field: fld128 strct89
	Functions
		0.- Function: fn44 (in49 strct48, in55 strct48, in58 strct48, in61 strct48, in66 strct48, in67 strct48, in68 strct48, in69 strct48, in95 mod42.strct48) 
		1.- Function: fn106 (in108 mod42.strct48, in116 mod42.strct89, in117 mod42.strct89, in138 mod42.strct89) 
2.- Module: mod146
```

which sadly didn't add any expression to any function it created (so
it wouldn't do anything if executed).

## Evolutionary Algorithm

For this small evolutionary algorithm, affordances are used to add
expressions to a function which represents a solution to a curve
fitting problem. Mutation is defined by randomly removing one expression from
a program, and then applying an affordance to replace the removed
expression. At the moment, all the necessary code to perform the
evolutionary process is located in the *base/evolution.go* file.

For this walkthrough, the observed points are obtained by evaluating
`f(x) = x*x*x - (3*x)` with the integer set `{-10, -9, 8, ..., 8, 9,
10}`. In *main.go*, these points are obtained by the following code:

```
dataIn := []int32{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	dataOut := make([]int32, len(dataIn))
	for i, in := range dataIn {
		dataOut[i] = in * in * in - (3 * in)
	}
```

The context (or program) must have at least a module named `main`, a
function named `main`, a solution function named `solution`, and the
`main` function must call solution with the input data as
arguments. Also, one can add the desired functions to be used as
operators for the expressions in the solution function (it can be
noted that more functions means a larger search space for the
evolutionary algorithm). For this walkthrough, the initial program is
defined by the following code:

```
num1 := encoder.SerializeAtomic(int32(0))

	cxt := MakeContext().AddModule(MakeModule("main"))
	if mod, err := cxt.GetCurrentModule(); err == nil {
		mod.AddDefinition(MakeDefinition("num1", &num1, i32))

		mod.AddFunction(MakeFunction("addI32"))
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("n1", i32))
			fn.AddInput(MakeParameter("n2", i32))
			fn.AddOutput(MakeParameter("out", i32))
		}
		
		mod.AddFunction(MakeFunction("subI32"))
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("n1", i32))
			fn.AddInput(MakeParameter("n2", i32))
			fn.AddOutput(MakeParameter("out", i32))
		}

		mod.AddFunction(MakeFunction("mulI32"))
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("n1", i32))
			fn.AddInput(MakeParameter("n2", i32))
			fn.AddOutput(MakeParameter("out", i32))
		}
		
		mod.AddFunction(MakeFunction("main"))
		mod.AddFunction(MakeFunction("double"))

		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("n", i32))
			fn.AddOutput(MakeParameter("out", i32))
			if addI32, err := cxt.GetFunction("mulI32", mod.Name); err == nil {
				fn.AddExpression(MakeExpression("out", addI32))

				if expr, err := cxt.GetCurrentExpression(); err == nil {
					expr.AddArgument(MakeArgument(MakeValue("n"), ident))
					expr.AddArgument(MakeArgument(MakeValue("n"), ident))
				}
			}
		}

		mod.AddFunction(MakeFunction("solution"))

		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("x", i32))
			fn.AddOutput(MakeParameter("f(x)", i32))
		}
		
		mod.SelectFunction("main")

		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddOutput(MakeParameter("outMain", i32))
			if solution, err := cxt.GetFunction("solution", mod.Name); err == nil {
				fn.AddExpression(MakeExpression("outMain", solution))
				
				if expr, err := cxt.GetCurrentExpression(); err == nil {
					expr.AddArgument(MakeArgument(MakeValue("num2"), ident))
				}
			}
		}
	}
```

Each of the inputs is going to be assigned to the global definition
`num1` (the reason for this is the mere convenience of reusing the
code for previous testings). As can be seen in the code, the operators
which can be used for the expressions are `addI32`, `subI32`,
`mulI32`, and `double`.

The `solution` function receives a single parameter `x`, and has the
named output `f(x)`.

The evolutionary process can be started by calling the method
`EvolveSolution` on the previously defined context:

```
    evolvedProgram := cxt.EvolveSolution(dataIn, dataOut, 5, 10000)
```

`EvolveSolution` takes four parameters: the data inputs, the data
outputs, the number of expressions which the solution should have, and
the number of maximum iterations the algorithm can perform, in that
order. In this case, the solution should have exactly five expressions
and the maximum number of iterations is 10,000. This exercise was ran
several times, and the optimal solution was always found in less than
1,000 iterations, so a stop condition was added to the algorithm,
which prevents any more evaluations to occur if the error is equal
to 0. An additional parameter (usually called epsilon) can be added to indicate the
algorithm to stop if the error is lower than this threshold.

While the algorithm is running, one can see the lowest achieved error
at each iteration printed to console. Once the evolutiona process has
finished, the evolved program can be printed by running:

```
    evolvedProgram.PrintProgram(false)
```

which should print something like this:

```
Context
0.- Module: main
	Definitions
		0.- Definition: num1 i32
		1.- Definition: num2 i32
	Functions
		0.- Function: solution (x i32) f(x) i32
			0.- Expression: x = mulI32(num1 i32, num1 i32)
			1.- Expression: var7388 = mulI32(x i32, num1 i32)
			2.- Expression: x = subI32(var7388 ident, num1 i32)
			3.- Expression: var38465 = addI32(num1 i32, num1 i32)
			4.- Expression: f(x) = subI32(x i32, var38465 ident)
		1.- Function: addI32 (n1 i32, n2 i32) out i32
		2.- Function: subI32 (n1 i32, n2 i32) out i32
		3.- Function: mulI32 (n1 i32, n2 i32) out i32
		4.- Function: main () outMain i32
			0.- Expression: outMain = solution(num2 i32)
		5.- Function: double (n i32) out i32
			0.- Expression: out = mulI32(n i32, n i32)
```

The simulated vs observed points can be printed by running the
following code:

```
var error int32 = 0
	for i, inp := range dataIn {
		num1 := encoder.SerializeAtomic(inp)
		if def, err := evolvedProgram.GetDefinition("num1"); err == nil {
			def.Value = &num1
		} else {
			fmt.Println(err)
		}

		evolvedProgram.Reset()
		evolvedProgram.Run(false, -1)

		// getting the simulated output
		var result int32
		output := evolvedProgram.CallStack[0].State["outMain"].Value
		encoder.DeserializeAtomic(*output, &result)

		diff := result - dataOut[i]
		fmt.Printf("Simulated #%d: %d\n", i, result)
		fmt.Printf("Observed #%d: %d\n", i, dataOut[i])
		if diff >= 0 {
			error += diff
		} else {
			error += diff * -1
		}
	}
```

Which will print something similar to this:

```
Simulated #0: -970
Observed #0: -970
Simulated #1: -702
Observed #1: -702
Simulated #2: -488
Observed #2: -488
Simulated #3: -322
Observed #3: -322
Simulated #4: -198
Observed #4: -198
Simulated #5: -110
Observed #5: -110
Simulated #6: -52
Observed #6: -52
Simulated #7: -18
Observed #7: -18
Simulated #8: -2
Observed #8: -2
Simulated #9: 2
Observed #9: 2
Simulated #10: 0
Observed #10: 0
Simulated #11: -2
Observed #11: -2
Simulated #12: 2
Observed #12: 2
Simulated #13: 18
Observed #13: 18
Simulated #14: 52
Observed #14: 52
Simulated #15: 110
Observed #15: 110
Simulated #16: 198
Observed #16: 198
Simulated #17: 322
Observed #17: 322
Simulated #18: 488
Observed #18: 488
Simulated #19: 702
Observed #19: 702
Simulated #20: 970
Observed #20: 970
```

And finally, the obtained error can be calculated by dividing the
accumulated absolute errors by the number of data points:

```
fmt.Println(error / int32(len(dataIn)))
```

## Call Stack Stepping

A call to an expression is represented by a struct. This struct
stores a reference to the expression's operator which is being called,
the line number currently being executed in the expression's operator,
a state (a set of definitions), and a return address (to what call do we
need to return once the current call finishes its execution). A
context or a program stores an array of calls, which is defined as a
*call stack*.

Implementing the execution of a program (by a series of calls,
represented by structs) in this way allows a program to store the
call stacks a program is creating during its execution. The storing of
a call stack is defined as a *step*, and the process of storing call
stacks is defined as *stepping* (these names are subject to change and
I will use them for convenience for now).

Stepping in the current implementation happens everytime a new call is
performed. This means that a program has access to every past point of
execution. This behaviour can easily be changed to saving the program
steps every *N* calls, in order to save system resources.

`MakeContextCopy(cxt *cxContext, stepNumber int)` can be used to
create a copy of a program. The second parameter, *stepNumber*, is
used to indicate at which point of execution one wants to create the
copy. For example, if a program was executing a loop, we can create a
copy of this program and "rollback" the steps to the point where the
program hasn't entered the loop yet. The following code makes a copy
of the evolved program from the last section, rollbacks to step 3, and
then executes again:

```
    copy := MakeContextCopy(evolvedProgram, 8)
	copy.Run(true, -1)
```

The `Run` method above receives two parameters: `withDebug bool`,
which prints the call stack at each step of the program execution; and
`nCalls int`, which tells the program how many calls it must run
before pausing its execution. The program doesn't raises an error if
it finishes its execution before the *nCalls* threshold is reached. If
we don't want a program to run for a certain number of calls, we can
simply give it a negative number of calls.

The method `ResetTo(stepNumber int)` can be called on a program to
rollback to the given step number, without making any copies of the
program. The code below rollsback a program to step 3, and then runs
the same program only for 8 more steps. This process is looped 5 times.

```
    for i := 0; i < 5; i++ {
      evolvedProgram.ResetTo(3)
      evolvedProgram.Run(true, 8)
    }
```

It is important to note that each execution in the loop above is
independent from the others, state-wise, i. e., they won't share
the values of their variables among them.

Call stack stepping can be used later on for debugging a program while
it's being executed. For example, if we have an evolutionary that has
been running for the last 5 days and it encounters an error, we can
inspect the state in each of the calls in the current stack, encounter
the problem, rollback *N* steps, make the necessary changes to the
program structure, and then resume the execution.

Another use for call stack stepping is for benchmarking code blocks of a
running program. We can create two copies of a program, each using a
particular solution to a problem (for example, one that uses a
solution based on looping, and the other uses recursion). We set each
of these copies to the start of the desired step, and run each program
until it finishes the code block being benchmarked. The code block
with the better performance can then be inserted to the original
program.

## Program Structure Stepping

A program is also aware of its own structure and the ordered steps
that have to be executed to reach its current structure. This means
that a program has the capability of rolling back to a certain step in
its structure, and we can create copies of a program at different
points of its structure stepping.

With the current implementation of CX, we must create a new program
and use the program steps of another to duplicate its structure to the
new program. For example, if we have a program stored in the variable
`cxt`, we can create three partial copies of it in the following way:

```
    copy1 := MakeContext()
	copy2 := MakeContext()
	copy3 := MakeContext()

    for i := 0; i < 15; i++ {
		cxt.ProgramSteps[i].Action(copy1)
	}

	for i := 0; i < 5; i++ {
		cxt.ProgramSteps[i].Action(copy2)
	}
	
	for i := 0; i < 10; i++ {
		cxt.ProgramSteps[i].Action(copy3)
	}

    copy1.PrintProgram(false)
	copy2.PrintProgram(false)
	copy3.PrintProgram(false)
```

The first loop will create a copy of the program until step 15, the
second loop a copy until step 5, and the final loop a copy until step 10.

The code above will print something similar to this:

```
Context
0.- Module: main
	Definitions
		0.- Definition: num1 i32
		1.- Definition: num2 i32
	Functions
		0.- Function: addI32 (n1 i32, n2 i32) out i32
		1.- Function: subI32 (n1 i32, n2 i32) out i32
		2.- Function: mulI32 (n1 i32, n2 i32) out i32
Context
0.- Module: main
	Definitions
		0.- Definition: num2 i32
		1.- Definition: num1 i32
	Functions
		0.- Function: addI32 (n1 i32) 
Context
0.- Module: main
	Definitions
		0.- Definition: num1 i32
		1.- Definition: num2 i32
	Functions
		0.- Function: addI32 (n1 i32, n2 i32) out i32
		1.- Function: subI32 (n1 i32, n2 i32) 
```

With call stack stepping, program structure stepping, and evolutionary
algorithms we could create a program which could stop itself at
certain step, measure the execution time of one of its functions,
mutate them and if it creates something better, automatically
modify itself, and then resume its execution. The same could be done
to replace buggy parts of a program: if a function raises an exception with a
certain combination of arguments, we can mutate the function until it
gives us the same outputs as the previous version of the function, but
also doesn't raise the exception.
