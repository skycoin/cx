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

## Randomly Generated Program

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
