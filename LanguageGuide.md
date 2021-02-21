# Table of Contents
* [Syntax](#syntax)
  * [Comments](#comments)
  * [Declarations](#declarations)
    * [Allowed Names](#allowed-names)
    * [Strict Type System](#strict-type-system-1)
    * [Primitive Types](#primitive-types)
    * [Global variables](#global-variables)
    * [Local variables](#local-variables)
    * [Arrays](#arrays)
    * [Slices](#slices)
    * [Literals](#literals)
    * [Functions](#functions)
    * [Custom Types](#custom-types)
    * [Methods](#methods)
    * [Packages](#packages)
  * [Statements](#statements)
    * [If and if/else](#if-and-ifelse)
    * [Else if](#else-if)
    * [For loop](#for-loop)
    * [Goto](#goto)
  * [Expressions](#expressions)
  * [Assignments and Initializations](#assignments-and-initializations)
  * [Affordances](#affordances)
* [Runtime](#runtime)
  * [Packages](#packages-1)
  * [Data Structures](#data-structures)
    * [Literals](#literals-1)
    * [Variables](#variables)
    * [Primitive types](#primitive-types-1)
    * [Arrays](#arrays-1)
    * [Slices](#slices-1)
    * [Structures](#structures)
    * [Pointers](#pointers)
    * [Escape Analysis](#escape-analysis)
  * [Control Flow](#control-flow)
    * [Functions](#functions-1)
    * [Methods](#methods-1)
    * [If and if/else](#if-and-ifelse-1)
    * [For loop](#for-loop-1)
    * [Go-to](#go-to)
  * [Affordances](#affordances-1)
* [Native Functions](#native-functions)
  * [Type-inferenced Functions](#type-inferenced-functions)
  * [Slice Functions](#slice-functions)
  * [Input/Output Functions](#inputoutput-functions)
  * [Parse Functions](#parse-functions)
  * [Unit Testing](#unit-testing)
  * [`bool` Type Functions](#bool-type-functions)
  * [`str` Type Functions](#str-type-functions)
  * [`i32` Type Functions](#i32-type-functions)
  * [`i64` Type Functions](#i64-type-functions)
  * [`f32` Type Functions](#f32-type-functions)
  * [`f64` Type Functions](#f64-type-functions)
  * [`time` Package Functions](#time-package-functions)
  * [`os` Package Functions](#os-package-functions)
  * [`gl` Package Functions](#gl-package-functions)
  * [`glfw` Package Functions](#glfw-package-functions)
  * [`gltext` Package Functions](#gltext-package-functions)

# Syntax
[[Back to the Table of Contents] ↑](#table-of-contents)

In this section, we're going to have a look at how a CX program looks
like. Basically, the following sections are not going to discuss about
the logic behind the various CX constructs, i.e. how they
behave; we're only going to see how they look like.

## Comments
[[Back to the Table of Contents] ↑](#table-of-contents)

Some of the code snippets that follow have comments in them, i.e.,
blocks of text that are not actually "run" by the CX compiler or
interpreter. Just like in C, Golang and many other programming
languages, single line comments are created by placing double slashes
(//) before the text being commented. For example:

```
// Example of adding two 32 bit integers in CX

i32.add(3, 4)       // This will be ignored

// End of the program
```

Mult-line comments are opened by writing slash-asterisk (/\*), and are
closed by writing asterisk-slash (\*/).

```
/* This code won't be executed
str.print("Hello world!")
*/
```

## Declarations
[[Back to the Table of Contents] ↑](#table-of-contents)

A declaration refers to a *named* element in a program's
structure, which are described using other constructs, such as
expressions and other statements. For example: a function can be
referred by its name and it's constructed by expressions and local
variable declarations.

### Allowed Names
[[Back to the Table of Contents] ↑](#table-of-contents)

Any name that satisfies the PCRE regular expression
`[_a-zA-Z][_a-zA-Z0-9]*` is allowed as an identifier for a declared
element. In other words, an identifier can start with an underscore
(*_*) or any lowercase or uppercase letter, and can be followed by 0
or more underscores or lowercase or uppercase letters, and any number
from 0 to 9.

### Strict Type System
[[Back to the Table of Contents] ↑](#table-of-contents)

One of CX's goals is to provide a very strict type system. The purpose
of this is to reduce runtime errors as much as possible. In order to
achieve this goal, many of CX's native functions are constrained to a
single type signature. For example, if you want to add two 32-bit
integers, you'd need to use `i32.add`. In contrast, if you want to add
two 64-bit integers, you'd use `i64.add`. These functions can help the
programmer to ensure that a particular data type is being received or
sent during a process.

If the programmer doesn't want to use those type-specific functions,
CX still provides type inference in some cases. For example, instead
of writing `i32.add(5, 5)` you can just write `5 + 5`. In this case,
CX is going to see that you're using 32-bit integers, and the parser
is going to transform that expression to `i32.add(5, 5)`. However, if
you try to do `5 + 5L`, i.e. if you try to add a 32-bit integer to a
64-bit integer, CX will throw a compile-time error because you're
mixing types.

The proper way to handle types in CX is to explicitly parse
everything. This way one can be sure that you're always going to be
handling the desired type. So, retaking the previous example, you'd
need to parse one of them to match the other's type, either
`i32.i64(5) + 5L` or `5 + i64.i32(5L)`.

### Primitive Types
[[Back to the Table of Contents] ↑](#table-of-contents)

There are seven primitive types in CX: *bool*, *str*, *byte*, *i32*,
*i64*, *f32*, and *f64*. Those represent Booleans (*true*
or *false*), character strings, bytes, 32-bit integers, 64-bit
integers, single precision and double precision floating-point
numbers, respectively.

### Global variables
[[Back to the Table of Contents] ↑](#table-of-contents)

Global variables are different from local variables regarding
scope. Global variables are available to any function defined in a
package, and to any package that is importing the package that
contains that global declaration. An example of some global variables
is shown below.

```
package main

var global1 i32
var global2 i64

func foo () {
    i32.print(global1)
    i64.print(global2)
}

func main () {
    global1 = 5
    global2 = 5L

    i32.print(global1)
    i64.print(global2)
}
```

In the example above we can see that both the `main` and `foo`
functions are printing the values of the two global variables
defined. They are going to print the same values, as they are
referring to the same variables.

### Local variables
[[Back to the Table of Contents] ↑](#table-of-contents)

In contrast to global variables, local variables are constrained to
the function where they are declared. This means that is not possible
for another function to call a variable defined in another function.

```
package main

func foo () {
	i32.print(local) // this expression will throw a compile-time error
}

func main () {
	var local i32
	local = 5

	i32.print(5)
	foo()
}
```

If you try to run the example above, CX will throw an error similar to
this: `error: examples/testing.cx:4 identifier 'local' does not
exist`, so CX will not even try to run that program. If we could
de-activate CX's compile-time type checking, and the program above
could make it to the runtime, CX would not print 5 when running
`foo()`, as that function is unaware of that variable.

### Arrays
[[Back to the Table of Contents] ↑](#table-of-contents)

Arrays (or vectors) and multi-dimensional arrays (or matrices) can be
declared using a syntax similar to C's.

```
package main

type Point struct {
    x i32
    y i32
}

func main () {
    var arr1 [5]i32
    var arr2 [5]Point
    var arr3 [2][2]f32

    arr1[0] = 10
    arr2[1] = 20
}
```

In the example above we see the declaration of an array of 5 elements
of type *i32*, followed by an array of the same cardinality but of
type *Point*, which is a custom type. Custom types are discussed in a
later section. Lastly, we see an example of a 2x2 matrix of type
*f32*.

Lastly, we can see how we can initialize an array using the bracket
notation, e.g. `arr1[0] = 10`.

### Slices
[[Back to the Table of Contents] ↑](#table-of-contents)

Golang-like slices exist in CX (dynamic arrays). Slices are declared
similarly to arrays, with the only difference that the size is
omitted.

```
package main

type Point struct {
    x i32
    y i32
}

func main () {
    var slc1 []i32
    var slc2 []Point
    var slc3 [][]f32
}
```

Slices, unlike arrays, cannot be directly initialized using the
bracket notation (unless you use the native function `make`
first). You can use the bracket notation to reassign values to a
slice, once an element associated to the index that you want to use
already exists, as shown in the example below.

```
package main

func main () {
    var slc []i32

    slc = append(slc, 1)
    slc = append(slc, 2)

    slc[0] = 10
    slc[2] = 30 // This is not allowed, as len(slc) == 2, not 3
}
```

As this behavior is more related to the logic behind slices, it is
further explained in the *Runtime->Data Structures->Slices* section.

### Literals
[[Back to the Table of Contents] ↑](#table-of-contents)

A literal is any data structure that is not being referenced by any
variable yet. For example: `1`, `true`, `[]i32{1, 2, 3}`, `Point{x:
10, y: 20}`.

Particularly, it is worth noting the cases of array, slice and struct
literals.

```
package main

type Point struct {
    x i32
    y i32
}

func main () {
    var a Point
    var b [5]i32
    var c []i32

    a = Point{x: 10, y: 20}
    b = [5]i32{1, 2, 3, 4, 5}
    c = []i32{100, 200, 300}
}
```

In the example above we can see examples of struct (`Point{x: 10,
y: 20}`), array (`[5]i32{1, 2, 3, 4, 5}`), and slice (`[]i32{100, 200, 300}`)
literals, in that order. These literals exist to simplify the creation
of such data structures.

### Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

Functions in CX are similar in syntax to functions in Go. The only
exception is that named outputs are enforced at the moment (this will
most likely change in the future).

```
package main

func foo () {

}

func main () {
	foo()
}
```

The example above doesn't do anything, but it illustrates the anatomy
of a function. In the case of `foo`, we have an empty function
declaration, and then we have `main`, which is defined by a single
call to `foo`. Functions can also receive inputs and return outputs,
as in the example below.

```
package main

func foo (in i32) {
    i32.print(in) // this will print 5
}

func bar () (out i32){
    out = 10
}

func main () {
    foo(5)

    var local i32
    local = bar()
    i32.print(local) // this will print 10
}
```

In this case, `foo` is declared to receive one input parameter, and
`bar` is declared to return one output parameter.

### Custom Types
[[Back to the Table of Contents] ↑](#table-of-contents)

If primitive types are not enough, you can define your own custom
types by combining the primitive types and other constructs like
slices, arrays, and even other custom types.

```
package main

type Point struct {
    x i32
    y i32
}

func main () {
    var p Point
    p.x = 10
    p.y = 20

    printf("Point coordinates: (%d, %d)\n", p.x, p.y)
}
```

In the example above, we can see a custom type that defines a
*Point* as the combination of two 32-bit integers (*i32*). After
declaring the custom type, you can start declaring variables of that
type anywhere in the package where it was declared in. The code in
`foo` shows how you can create and use an instance of that structure.

### Methods
[[Back to the Table of Contents] ↑](#table-of-contents)

A variation of functions that are associated to custom types are
*methods*.

```
package main

type Point struct {
	x i32
	y i32
}

type Line struct {
	a Point
	b Point
}

func (p Point) print () {
	printf("Point coordinates: (%d, %d)\n", p.x, p.y)
}

func (l Line) print () {
	printf("Line point A: (%d, %d), Line point B: (%d, %d)\n", l.a.x, l.a.y, l.b.x, l.b.y)
}

func main () {
	var l Line
	var p1 Point
	var p2 Point

	p1.x = 10
	p1.y = 20
	p2.x = 11
	p2.y = 21

	l.a = p1
	l.b = p2

	p1.print()
	p2.print()
	l.print()
}
```

In the example above, we define two custom types: *Point* and
*Line*. The type line is defined by two fields of type
*Point*, and the type *Point* is defined as coordinate defined by two
fields of type *i32*.

As a simple example, we create two methods called `print`, one for the
type *Point* and another for the type *Line*. In the case of
`Point.print`, we just print the two coordinates, and in
the case of `Line.print` we print the coordinates of the two points
that define the *Line* instance.

### Packages
[[Back to the Table of Contents] ↑](#table-of-contents)

In the previous examples we have always been using a single package:
`main`. If your program grows too large it's convenient to divide your
code into different packages.


```
package foo

func fn (in i32) {
    i32.print(in)
}

package bar

func fn () (out i32) {
    out = 5
}

package main
import "foo"
import "bar"

func main () {
    foo.fn(10) // prints 10

    var num i32
    num = bar.fn()

    i32.print(num) // prints
}
```

In the example above, we can see how two functions with the same name
(`fn`) are declared, each in a separate package. Both of these
functions have different signatures, as `foo.fn` accepts a single
input parameter and `bar.fn` doesn't accept any inputs but returns a
single output parameter.

We can then see how the `main` package `import`s both the `foo` and
`bar` packages, to later call each of these functions.

## Statements
[[Back to the Table of Contents] ↑](#table-of-contents)

Statements are different to declarations, as they don't create any
named elements in a program. They are used to control the flow of a
program.

### If and if/else
[[Back to the Table of Contents] ↑](#table-of-contents)

The most basic statement is the *if* statement, which is going to
execute a block of code only if a condition is true.

```
package main

func main () {
    if false {
        str.print("This will never print")
    }
}
```

The example above won't do anything, as the condition for the *if*
statement is always going to evaluate to *false*.

```
package main

func main () {
    if true {
        str.print("This will always print")
    }
}
```

In contrast, the example above will always print.

```
package main

func main () {
    if true {
        str.print("This will always print")
    } else {
        str.print("This will never print")
    }
}
```

Lastly, the example above shows how to write an *if/else* statement in CX.

As a note about its syntax, the predicates or conditions don't need to
be enclosed in parentheses, just like in Go.

### Else if
[[Back to the Table of Contents] ↑](#table-of-contents)

Instead of simply adding one alternative path, you can string together a series
of *else if* blocks, which check for as many different conditions as you like.
Giving you similar functionality as Go's *switch*/*select* blocks (containing
various conditions/cases).

```
package main

func main () {
   var i i32
   i = 0

   if i == 0 {
     str.print("i is 0")
   } else if i == 1 {
     str.print("i is 1")
   } else if i == 2 {
     str.print("i is 2")
   } else {
     str.print("i is NOT 0, 1 or 2")
   }
}
```

### For loop
[[Back to the Table of Contents] ↑](#table-of-contents)

CX's only looping statement is the *for* loop. Similar to Go, the
*for* loop in CX can be used as the *while* statement in other
programming languages, and as a traditional *for* statement.

```
package main

func main () {
	for true {
		str.print("forever")
	}
}
```

As the simplest example of a loop, we have the infinite loop shown in
the example above. In this case, the loop will print the character
string "forever" indefinitely. If you try this code, remember that you
can cancel the program's execution by hitting `Ctrl-C`.

```
package main

func main () {
	for str.eq("hi", "hi") {
		str.print("hi")
	}
}
```

The code above shows another example, one where we use an expression
as its predicate, rather than a literal `true` or `false`. It is worth
mentioning that you can replace `str.eq("hi", "hi")` by `"hi" == "hi"`.


```
package main

func main () {
	var c i32
	for c = 0 ; c < 10; c++ {
		i32.print(c)
	}
}
```

The traditional *for* loop shown in the example above. In languages
like C, you need to first declare your counter variable, and then you
have the option to initialize or reassign the counter in the first
part of the *for* loop. The second part of the *for* loop is reserved
for the predicate, and the last part is usually used to increment the
counter. Nevertheless, just like in C, there's nothing stopping you
from doing whatever you want in the first and last parts. However, the
predicate part needs to be an expression that evaluates to a Boolean.

```
package main

func main () {
	for c := 0; c < 10; c++ {
		i32.print(c)
	}
}
```

A more Go-ish way of declaring and initializing the counter is to use
an inline declaration, as seen in the example above.

```
package main

func main () {
	var c i32
	c = 0
	for ; c < 10; c++ {
		i32.print(c)
	}
}
```

Lastly, the for loop can also completely omit the initialization part,
as seen above.

### Goto
[[Back to the Table of Contents] ↑](#table-of-contents)

`goto` can be used to immediately jump the execution of a program to
the corresponding labeled expression.

```
package main

func main () (out i32) {
	goto label2
label1:
	str.print("this should never be reached")
label2:
	str.print("this should be printed")
}
```

In the example above, we see how a `goto` statement forces CX to
ignore executing the expression labeled as `label1`, and instead jumps
to the `label2` expression.

## Expressions
[[Back to the Table of Contents] ↑](#table-of-contents)

Expressions are basically function calls. But the term expression also
takes into consideration the variables that are receiving the
function's output arguments, the input arguments, and any dereference
operations.

```
package main

func foo () (arr [2]i32) {
    arr = [2]i32{10, 20}
}

func main () {
    i32.print(foo()[0])
}
```

For example, the expression `i32.print(foo()[0])` in the code above
consists of two function calls, and the array returned by the call
to `foo` is "dereferenced" to its *0th* element.

## Assignments and Initializations
[[Back to the Table of Contents] ↑](#table-of-contents)

As in many other C-like languages, assignments are done using the
equal (`=`) sign.

```
package main

func main () {
    var foo i32
    foo = 5
    foo = 50
}
```

In the case of the code above, the variable `foo` is declared and then
initialized to `5` using the `=` sign. Then we reassign the `foo`
variable to the value `50`.

```
package main

func main () {
    foo := 5
    foo = 50
}
```

As in other programming languages, *short variable declarations* exist
in CX. The `:=` token can be used to tell CX to infer a variable's
type. This way, CX declares and initializes at the same time, as seen
in the example above.

## Affordances
[[Back to the Table of Contents] ↑](#table-of-contents)

The affordance system in CX uses a special operator: `->`.  This
operator takes a series of statements that have the form of function
calls, and transforms them to a series of instructions that can be
internally interpreted by the affordance system.


```
package main

func exprPredicate (expr Expression) (res bool) {
	if expr.Operator == "i32.add" {
		res = true
	}
}

func prgrmPredicate (prgrm Program) () {
	if prgrm.FreeHeap > 50 {
		res = true
	}
}

func main () {
	num1 := 5
	num2 := 10

targetExpr:
	sum := i32.add(0, 0)

	tgt := ->{
		pkg(main)
		fn(main)
		expr(targetExpr)
	}

	fltrs := ->{
		filter(exprPredicate)
		filter(prgrmPredicate)
	}

	aff.print(tgt)
	aff.print(fltrs)

	affs := aff.query(fltrs, tgt)

	aff.on(affs, tgt)
	aff.of(affs, tgt)

	aff.inform(affs, 0, tgt)
	aff.request(affs, 0, tgt)
}
```

# Runtime
[[Back to the Table of Contents] ↑](#table-of-contents)

The previous section presents the language features from a syntax
perspective. In this section we'll cover what's the logic behind these
features: how they interact with other elements in your program, and
what are the intrinsic capabilities of each of these features.

## Packages
[[Back to the Table of Contents] ↑](#table-of-contents)

Packages are CX's mechanism for better organizing your code. Although
it is theoretically possible to store a big project in a single
package, the code will most likely become very hard to understand. In
CX the programmer is encouraged to place the files that define the
code of a package in separate directory. Any subdirectory in a
package's directory should also contain only source code files that
define elements of the same package. Nevertheless, CX will not
throw any error if you don't follow this way of laying out your source
files. In fact, you can declare different packages in a single source
code file.

## Data Structures
[[Back to the Table of Contents] ↑](#table-of-contents)

Data structures are particular arrangements of bytes that the language
interprets and stores in special ways. The most basic data structures
represent basic data, such as numbers and character strings, but these
basic types can be used to construct more complex data types.

### Literals
[[Back to the Table of Contents] ↑](#table-of-contents)

A literal is any data structure that is not being referenced by any
variable yet. For example: `1`, `true`, `[]i32{1, 2, 3}`, `Point{x:
10, y: 20}`.

It's important to make a distinction, particularly with *arrays*,
*slices* and *struct instances*.

```
package main

type Point struct {
	x i32
	y i32
}

func main () {
	var p1 Point
	p1.x = 10
	p1.y = 20

	p2 := Point{x: 11, y: 21}

	i32.print(p2.x)
	i32.print(p2.y)
}
```

In the example above we are creating two instances of the `Point`
type. The first method we use does not involve struct literals, as a
variable of that type is first created and then initialized.

In the second case (`p2`), the full struct instance is first
created. CX creates an anonymous struct instance as soon as it
encounters `Points{x: 11, y: 21}`, and then it proceeds to assign that
literal to the `p2` variable, using *short variable declarations*.

```
package main

func main () {
	var arr1 [3]i32
	arr1[0] = 1
	arr1[1] = 2
	arr1[2] = 3

	arr2 := [3]i32{10, 20, 30}
}
```

```
package main

func main () {
	var slc1 []i32
	slc1 = append(slc1, 1)
	slc1 = append(slc1, 2)
	slc1 = append(slc1, 3)

	slc2 := []i32{10, 20, 30}
}
```
Similarly, in the two examples above we can see how we can declare
array and slice variables and then we initialize them. In the case of
arrays, we use the bracket notation, and for slices we have to use
`append`, as `slc1` starts with a size and capacity of 0. In the cases
of `arr2` and `slc2`, we use literals to initialize them more
conveniently.

Regarding numbers, you need to be aware that implicit casting does not
exist in CX. This means that the number `34` cannot be assigned to a
variable of type *i64*. In order to assign it, you need to either
parse it using the native function `i32.i64` or you can create a
64-bit integer literal. To create a number literal of a type other
than *i32*, you can use different suffixes: `B`, `L` and `D`, for
*byte*, *i64* (long) and *f64*, respectively. So, assuming `foo` is of
type *i64*, you can do this assignment: `foo = 34L`.

### Variables
[[Back to the Table of Contents] ↑](#table-of-contents)

When CX compiles a program, it knows how many bytes need to be
reserved in the stack for each of the functions. CX can know this
thanks to variable declarations.


```
package main

type Point struct {
    x i32
    y i32
}

func foo (inp Point) {
    var test1 i64
    var test2 bool
}

func main () {
    var test3 i32
    var test4 f32
}
```

The two functions declared in the example above are going to reserve 17
and 8 bytes in the stack, respectively. In the case of the first
function, `foo` needs to reserve space for an input parameter of type
`Point`, which requires 8 bytes (because of the two *i32* fields), and
two local variables: one 64-bit integer that requires 8 bytes and a
Boolean that requires a single byte. In the case of `main`, CX needs
to reserve bytes for two local variables: a 32-bit integer and a
single-precision floating point number, where each of them require 4
bytes.

```
package main

var global1 i32

func main () {
    var local i32
}
```

Local variables are different than global variables. In order for
globals to have a global scope they need to be allocated in a
different memory segment than local variables. This different memory
segment does not shrink or get bigger like the stack. This means that
any global variable is going to be kept "alive" as long as the program
keeps being executed.

A global scope means that variables of this type are accessible to any
function declared in the same package where the variable is declared,
and to any function of other packages that are importing this package.

```
package main

func main () {
    var foo i32
    i32.print(foo) // prints 0
}
```

In CX every variable is going to initially point to a *nil*
value. This *nil* value is basically a series of one or more zeroes,
depending on the size of the data type of a given variable. For
example, in the code above we see that we have declared a variable of
type *i32* and we immediately print its value without initializing
it. This CX program will print 0, as the value of `foo` is `[0 0 0 0]`
in the stack (4 zeroes, as a 32-bit integer is represented by 4
bytes).

In the case of data types that point to variable-sized
structures, such as slices or character strings, these are initialized
to a nil pointer, which is represented by 4 zeroed bytes. This nil
pointer is located in the heap memory segment, instead of the stack.

### Primitive types
[[Back to the Table of Contents] ↑](#table-of-contents)

There are seven primitive types in CX: *bool*, *str*, *byte*, *i32*,
*i64*, *f32*, and *f64*. These types can be used to construct other
more complex types, as will be seen in the next sections.

*bool* and *byte* both require a single byte to represent their
 values. In the case of *bool*, there are only two possible values:
 `true` or `false`. In the case of *byte* you can represent up to 256
 values, which range from 0 to 255. Next in size, we have *i32* and
 *f32* , where both of them require 4 bytes, and then we have *i64*
 and *f64*, which require 8 bytes each.

Now, strings are special as they are static and dynamic sized at the
same time. If you have a look at how a variable of type *str* reserves
memory in the stack, you'll see that it requires 4 bytes, regardless
of what text it's pointing to. The explanation behind this is that any
*str* in CX actually behaves like a pointer behind the scenes, and the
actual string gets stored in the heap memory segment.

```
package main

func main () {
	var foo str

	foo = str.concat("Hello, ", "World!")
	foo = "Hi"
}
```

When CX compiles the example above, three strings are first stored in
the data memory segment (just like global variables, as these strings
are constants, memory-wise): `"Hello, "`, `"World"` and `"Hi"`. When
the program is executed, `str.concat` is called, which creates a new
string by concatenating `"Hello, "` and `"World!"`, and this new
character string is allocated in the heap memory segment. Then `foo`
is assigned only the address of this new character string. Then we
immediately re-assign `foo` with the address of `"Hi"`. This means
that `foo` was first assigned a memory address located in the data
memory segment, and then it was assigned an address located in the
heap.

### Arrays
[[Back to the Table of Contents] ↑](#table-of-contents)

Arrays, as in other programming languages, are used to create
collections of data structures. These data structures can be primitive
types, custom types or even arrays or slices.


```
package main

type Point struct {
    x i32
    y i32
}

func main () {
    var [5]i32
    var [5]Point
}
```

In the example above, we're creating two arrays, one of a primitive
type and the other one of a custom type. CX reserves memory for these
arrays in the stack as soon as the function that contains them is
called. In this case, 60 bytes are going to be reserved for `main` as
soon as the program starts its execution, as `main` acts as the
program's entry point. You need to be careful with arrays, as those
can easily fill up your memory, especially with multi-dimensional
arrays (or matrices).

Also, another point to consider is performance. While accessing arrays
is almost as fast as accessing an atomic variable, arrays can be
troublesome when being sent/received as to/from functions. The reason
behind this is that an array needs to be copied whenever it is sent to
another function. If you're working with arrays of millions of
elements and you need to be sending that arrays millions of times to
another function, it's going to impact your program's performance a
lot. A way to avoid this is to either use pointers to arrays or slices.

### Slices
[[Back to the Table of Contents] ↑](#table-of-contents)

Dynamic arrays don't exist in CX. This means that the following
code is not a valid CX program:

```
package main

func main () {
    var size i32
    size = 13
    var arr [size] // this is not valid
}
```

If you need an array that can grow in size as required, you need to
use slices. Behind the scenes, slices are just arrays with some extra
features. First of all, any slice in CX goes directly to the heap, as
it's a data structure that is going to be changing in size. In
contrast, arrays are always going to be stored in the stack, unless
we're handling pointers to arrays. However, this behavior may change
in the future, when CX's escape analysis mechanism improves (for
example, the compiler can determine if an array is never going to change
its size, and decide to keep it in the stack).

The second characteristic of slices in CX is how they change their
size. Any slice, when it's first declared, starts with a size and
capacity of 0. The size represents how many elements are in a given
slice, while the capacity represents how many elements can be
allocated in that slice without having to be relocated in the heap.

```
package main

func main () {
	var slc []i32

	slc = append(slc, 1)
	slc = append(slc, 2)
	slc = append(slc, 3)
	slc = append(slc, 4)
}
```

In the code above we can see how we declare a slice and then we
initialize it using the `append`function. After all the `append`s,
we'll end up with a slice of size 4 and capacity 4, and this
`append`ing process will create the following objects in the heap:

```
[0 0 0 0 0 12 0 0 0 1 0 0 0 1 0 0 0 1 0 0 0 0 0 0 0 0 16 0 0 0 2 0 0 0 2 0 0 0 1 0 0 0 2 0 0 0 0 0 0 0 0 24 0 0 0 4 0 0 0 4 0 0 0 1 0 0 0 2 0 0 0 3 0 0 0 4 0 0 0]
```

First, the slice `slc` starts with 0 objects in it; it is pointing to
*nil*. Then, after the first `append`, the object
`[0 0 0 0 0 12 0 0 0 1 0 0 0 1 0 0 0 1 0 0 0]` is allocated to the
heap. The first five bytes are used by CX's garbage collector. The
next 4 bytes indicate the size of the object, and the remaining bytes
are the actual slice `slc`. The first four bytes of `slc` tell us its
current size, while the next four tell us its capacity. The remaining
bytes of this object are the elements of the slice.

The following object,
`[0 0 0 0 0 16 0 0 0 2 0 0 0 2 0 0 0 1 0 0 0 2 0 0 0]`, shows now a
size of 2 and a capacity of 2, with the 32-bit integers `1` and `2` as
its elements. The last object, `0 0 0 0 0 24 0 0 0 4 0 0 0 4 0 0 0 1 0
0 0 2 0 0 0 3 0 0 0 4 0 0 0`, needs careful attention. We can see that
our objects jumped from size 1 to 2 and finally 4. The same happened
to its capacity, and the containing elements are now `1`, `2`, `3` and
`4`. What happened to the slice of size 3 and capacity 3? First of
all, capacities are increased by getting doubled each time the size of
an object is greater than its capacity, so we would never get a slice
of capacity 4 by following this method. Next, we need to think on what
is capacity used for.

Slices are just arrays, which means that they can't be resized. The
dynamic nature of slices is emulated by copying the *full* slice to
somewhere else in memory, but with a greater capacity. However, this
will only happen if adding a new element to the existing slice would
overflow it. This is why slices keep track of two metrics: *size* and
*capacity*, i.e. how many actual elements are in the slice, and how
many elements the currently allocated slice can hold, respectively.

```
package main

func main () {
    var arr1 [1]i32
    arr1[0] = 1 // add the first value

    var arr2 [2]i32 // double the size
    arr2[0] = arr1[0] // copy previous array
    arr2[1] = 2 // add the second value

    var arr3 [4]i32 // double the size
    arr3[0] = arr2[0] // copy previous array
    arr3[1] = arr2[1] // copy previous array

    arr3[2] = 3 // add the third value
    arr3[3] = 4 // add the fourth value
}
```

The example above shows the behavior of the slice in the previous
example, but using arrays.

### Structures
[[Back to the Table of Contents] ↑](#table-of-contents)

Structures are CX's mechanism for creating custom types, as in many other
C-like languages. Structures are basically a grouping of other
primitive or custom types (called *fields*) that together create
another type of data structure. For example, a *point* can be defined
by its coordinates in a two-dimensional space. In order to create a
type `Point`, you can use a structure that contains two fields of type
*i32*, one for `x` and another for `y`, as in the example below.


```
package main

type Point struct {
	x i32
	y i32
}

func main () {
	var p Point
	p.x = 10
	p.y = 20
}
```

Whenever an instance of a structure is created by either declaring a
variable of that type or by creating a literal of that type, CX
reserves memory to hold space for all the fields defined in the
structure declaration. Like in C, the bytes are reserved depending on
the order of the fields in the structure declaration.

```
package main

type struct1 struct {
	field1 bool
	field2 i32
	field3 i64
}

type struct2 struct {
	field1 i64
	field2 bool
	field3 i32
}

func main () {
	var s1 struct1
	var s2 struct2
}
```

For example, in the code above a call to `main` will reserve a total
of 26 bytes in the stack. In the case of the first struct instance,
the first byte is going to represent `field1` of type *bool*, the next
four bytes are going to represent `field2` of type *i32*, and the
final 8 bytes are going to represent `field3` of type *i64*. In the
case of the next struct instance, the first eight bytes represent an
*i64* field so, although both struct instances contain the same number
of fields and of the same type, the byte layout changes.

### Pointers
[[Back to the Table of Contents] ↑](#table-of-contents)

Sometimes it's useful to pass variables to functions by reference
instead of by value.

```
package main
import "time"

func foo (nums [100][100]i32) {
	// do something with nums
}

func main () {
	var start i64
	var end i64
	var nums [100][100]i32

	start = time.UnixMilli()

	for c := 0; c < 10000; c++ {
		foo(nums)
	}

	end = time.UnixMilli()

	printf("elapsed time: \t%d milliseconds\n", end - start)
}
```

The example above is very inefficient, as CX is going to be sending a
10,000 element matrix to `foo` 10,000 times. Every time `foo` is
called, every byte of that matrix needs to be copied for `foo`. In my
computer the example above takes around 638 milliseconds to run.

```
package main
import "time"

func foo (nums *[100][100]i32) {
	// do something with nums
}

func main () {
	var start i64
	var end i64
	var nums [100][100]i32

	start = time.UnixMilli()

	for c := 0; c < 10000; c++ {
		foo(&nums)
	}

	end = time.UnixMilli()

	printf("elapsed time: \t%d milliseconds\n", end - start)
}
```

A new version of the last program is shown above. In contrast to the
last program, the code above sends a pointer to the matrix to
`foo`. A pointer in CX uses only 4 bytes (in the future, pointers will
use 8 bytes in 64-bit systems and 4 bytes in 32-bit systems), so
instead of copying 10,000 bytes, we only copy 4 bytes to `foo` every
time we call it. This version of the program takes only 3 milliseconds
to run in my computer.

```
package main

func foo (inp i32) {
	inp = 10
}

func main () {
	var num i32
	num = 15
	i32.print(num) // prints 15
	foo(num)
	i32.print(num) // prints 15
}
```

In the example above, we send`num` to `foo`, and then we re-assign
the input's value to `10`. If we print the value of `num` before and
after calling `foo`, we can see that in both instances `15` will be
printed to the console.

```
package main

func foo (num *i32) {
	*num = 10
}

func main () {
	var num i32
	num = 15
	i32.print(num) // prints 15
	foo(&num)
	i32.print(num) // prints 10
}
```

The code above is a pointer-version of the previous example. In this
case, instead of sending `num` by value, we send it by reference,
using the `&` operator. `foo` also changed, and it now accepts a
pointer to a 32-bit integer, i.e. `*i32`. After running the example,
you'll notice that, this time, `foo` is now changing `num`'s value.

### Escape Analysis
[[Back to the Table of Contents] ↑](#table-of-contents)

Consider the following example:

```
package main

func foo () (pNum *i32) {
	var num i32
	num = 5 // this is in the stack

	pNum = &num
}

func stackDestroyer () {
	var arr [5]i32
}

func main () {
	var pNum *i32
	pNum = foo()

	stackDestroyer()

	i32.print(*pNum)
}
```

If we store `foo`'s `num`'s value (`5`) in the stack, and then we call
`stackDestroyer`, isn't `arr` going to overwrite the bytes storing the
`5`? This doesn't happen, because that `5` is now in the heap. But
this doesn't mean that any value being pointed to is going to be moved
to the heap. For example, let's re-examine one of the examples
presented in the *Pointers* section:

```
package main

func foo (num *i32) {
	*num = 10
}

func main () {
	var num i32
	num = 15
	i32.print(num) // prints 15
	foo(&num)
	i32.print(num) // prints 10
}
```

If any value being pointed to by a pointer was sent to the heap, we
wouldn't be able to change `num`s value, which is stored in the stack;
we would be changing the heap's copied value.

```
package main

func foo () (pNum *i32) {
	var num i32
	var pNum *i32

	num = 5

	pNum = &num
}

func main () {
	var pNum *i32
	pNum = foo()

	i32.print(*pNum) // prints 5, which is stored in the heap
}
```

Basically, in order to fix this problem, whenever a pointer needs to
be returned from a function, the value it is pointing to "escapes" to
the heap. In the example above, we can see that `num`'s value is going
to be preserved by escaping to the heap, as we are returning a pointer
to it from `foo`.

```
package main

func foo () (pNum *i32) {
	var num i32
	var pNum *i32

	num = 5 // this is in the stack

	pNum = &num // the pointer will be returned, so the value is sent to the heap
}

func stackDestroyer () {
	var arr [5]i32
}

func main () {
	var pNum *i32
	pNum = foo()
	stackDestroyer() // if 5 does not escape, it would be destroyed by this function

	i32.print(*pNum) // prints 5, which is stored in the heap
}
```

We can check this behavior even further in the example above. After
calling `foo`, we call `stackDestroyer`, which overwrites the
following 20 bytes after `main`'s stack frame. Yet, when we call
`i32.print(*pNum)`, we'll see that we still have access to a `5`. This
`5` is not the one created in `foo`, though, but a copy of it that was
allocated in the heap.

## Control Flow
[[Back to the Table of Contents] ↑](#table-of-contents)

Once we have the appropriate data structures for our program, we'll
now need to process them. In order to do so, we need to have access to
some control flow structures.

### Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

Functions are used to encapsulate routines that we plan to be
frequently calling. In addition to encapsulating a series of
expressions and statements, we can also receive input parameters and
return output parameters, just like mathematical functions.

```
package main

func main () {
	var players []str
	players = []str{"Richard", "Mario", "Edward"}

	str.print("=======================")
	str.print(str.concat("Name: \t", players[0]))
	str.print("=======================")

	str.print("=======================")
	str.print(str.concat("Name: \t", players[1]))
	str.print("=======================")

	str.print("=======================")
	str.print(str.concat("Name: \t", players[2]))
	str.print("=======================")
}
```

For example, if we see the code above we'll notice that it seems
repetitive. We can fix this by creating a function, as seen in the
example below.


```
package main

func drawBox (player str) {
	str.print("=======================")
	str.print(str.concat("Name: \t", player))
	str.print("=======================")
}

func main () {
	var players []str
	players = []str{"Richard", "Mario", "Edward"}

	drawBox(players[0])
	drawBox(players[1])
	drawBox(players[2])
}
```

### Methods
[[Back to the Table of Contents] ↑](#table-of-contents)

Methods are useful when we want to associate a particular function to
a particular custom type (associating functions to primitive types is
not allowed). This allows us to create more readable code.

```
package main

type Player struct {
	Name str
	HP i32
	Mana i32
	Lives i32
}

type Monster struct {
	Name str
	HP i32
	Mana i32
}

func (player Player) draw () {
	str.print(sprintf("\n\tName: \t%s\n\tHP: \t%d\n\tMana: \t%d\n\tLives: \t%d\n\n%s",
		player.Name,
		player.HP,
		player.Mana,
		player.Lives,
		`
─▄████▄▄░
▄▀█▀▐└─┐░░
█▄▐▌▄█▄┘██
└▄▄▄▄▄┘███
██▒█▒███▀`))
}

func (monster Monster) draw () {
	str.print(sprintf("\n\tName: \t%s\n\tHP: \t%d\n\tMana: \t%d\n\n%s",
		monster.Name,
		monster.HP,
		monster.Mana,
		`
╲╲╭━━━━╮╲╲
╭╮┃▆┈┈▆┃╭╮
┃╰┫▽▽▽▽┣╯┃
╰━┫△△△△┣━╯
╲╲┃┈┈┈┈┃╲╲
╲╲┃┈┏┓┈┃╲╲
▔▔╰━╯╰━╯▔▔`))
}

func main () {
	var player Player
	player.Name = "Mario"
	player.HP = 10
	player.Mana = 10
	player.Lives = 3

	player.draw()

	var monster Monster
	monster.Name = "Domo-kun"
	monster.HP = 7
	monster.Mana = 4

	monster.draw()
}
```

The example above shows us how we can create two versions of the
function `draw`, and the behavior of each depends on the custom type
that we're using to call it.

### If and if/else
[[Back to the Table of Contents] ↑](#table-of-contents)

*if* and *if/else* statements are used to execute a block of
 instructions only if certain condition is true or false. Behind the
 scenes, *if* and *if/else* statements are parsed to a series of
 `jmp` instructions internally. For example, in the case of an *if*
 statement, we will jump *0* instructions if certain predicate is
 true, and it will jump *n* instructions if the predicate is false,
 where *n* is the number of instructions in the *if* block of
 instructions.

```
package main

func main () {
	if true {
		str.print("hi")
	}
	str.print("bye")
}
```

```
Program
0.- Package: main
	Functions
		0.- Function: main () ()
			0.- Expression: jmp(true bool)
			1.- Expression: str.print("" str)
			2.- Expression: jmp(true bool)
			3.- Expression: str.print("" str)
		1.- Function: *init () ()
```

In the two code snippets above we can see how an if statement is
translated by the parser to a set of two `jmp` instructions. These
`jmp` instructions have some meta data in them that is not shown in
the second snippet: how many lines to jump if its predicate is true
and how many lines to jump if the predicate is false. `jmp` is
not meant to be used by CX programmers (it's only part of the CX base
language), so you don't need to worry about it.

```
package main

type Player struct {
	Name str
	HP i32
	Mana i32
	Lives i32
}

type Monster struct {
	Name str
	HP i32
	Mana i32
}

func main () {
	var player Player
	player.Name = "Mario"
	player.HP = 10
	player.Mana = 10
	player.Lives = 3

	var monster Monster
	monster.Name = "Domo-kun"
	monster.HP = 7
	monster.Mana = 4

	if player.HP < 5 {
		str.print("===DANGER!===")
	} else {
		str.print("===YOU CAN DO IT!===")
	}

	if monster.HP < 10 {
		str.print(sprintf("===%s is bleeding!===", monster.Name))
	}

	if monster.HP < 5 {
		str.print(sprintf("===%s is dying!===", monster.Name))
	}

	if monster.HP == 0 {
		str.print(sprintf("===%s is dead!===", monster.Name))
	}
}
```

Continuing with the example from the previous section (to some
extent), let's use *if* and *if/else* statements to determine what
messages are going to be displayed to the user. These messages
represent the state of the player or the monster, depending on their
hit points (HP).

### For loop
[[Back to the Table of Contents] ↑](#table-of-contents)

The *for* loop is the only looping mechanism in CX. Just like *if* and
*if/else* statements are constructed using `jmp` statements, *for*
loop statements are also constructed the same way.

```
package main

func main () {
	for c := 0; c < 10; c++ {
		i32.print(c)
	}
}
```

```
Program
0.- Package: main
	Functions
		0.- Function: main () ()
			0.- Declaration: c i32
			1.- Expression: c i32 = identity(0 i32)
			2.- Expression: *lcl_0 bool = lt(c i32, 10 i32)
			3.- Expression: jmp(*lcl_0 bool)
			4.- Expression: i32.print(c i32)
			5.- Declaration: c i32
			6.- Expression: c i32 = i32.add(c i32, 1 i32)
			7.- Expression: jmp(true bool)
		1.- Function: *init () ()
```

The code snippets above illustrate how a *for* loop that counts from 0
to 9 is translated to a set of of `jmp` instructions.

```
package main

type Player struct {
	Name str
	HP i32
	Mana i32
	Lives i32
}

type Monster struct {
	Name str
	HP i32
	Mana i32
}

func (player Player) draw () {
	str.print(sprintf("\n\tName: \t%s\n\tHP: \t%d\n\tMana: \t%d\n\tLives: \t%d\n\n%s",
		player.Name,
		player.HP,
		player.Mana,
		player.Lives,
		`
─▄████▄▄░
▄▀█▀▐└─┐░░
█▄▐▌▄█▄┘██
└▄▄▄▄▄┘███
██▒█▒███▀`))
}

func (monster Monster) draw () {
	str.print(sprintf("\n\tName: \t%s\n\tHP: \t%d\n\tMana: \t%d\n\n%s",
		monster.Name,
		monster.HP,
		monster.Mana,
		`
╲╲╭━━━━╮╲╲
╭╮┃▆┈┈▆┃╭╮
┃╰┫▽▽▽▽┣╯┃
╰━┫△△△△┣━╯
╲╲┃┈┈┈┈┃╲╲
╲╲┃┈┏┓┈┃╲╲
▔▔╰━╯╰━╯▔▔`))
}

func (player Player) attack (cmd str, monster *Monster) {
	if bool.or(cmd == "M", cmd == "m") {
		var dmg i32
		dmg = i32.rand(1, 4)
		(*monster).HP = (*monster).HP - dmg
		printf("'%s' suffered a magic attack. Lost %d HP. New HP is %d\n", (*monster).Name, dmg, (*monster).HP)
	} else {
		var dmg i32
		dmg = i32.rand(1, 2)
		(*monster).HP = (*monster).HP - dmg
		printf("'%s' suffered a physical attack. Lost %d HP. New HP is %d\n", (*monster).Name, dmg, (*monster).HP)
	}
}

func (monster Monster) attack (cmd str, player *Player) {
	var dmg i32
	dmg = i32.rand(1, 5)
	(*player).HP = (*player).HP - dmg

	printf("'%s' suffered a physical attack. Lost %d HP. New HP is %d\n", (*player).Name, dmg, (*player).HP)
}

func battleStatus (player Player, monster Monster) {
	if player.HP < 5 {
		str.print("===DANGER!===")
	} else {
		str.print("===YOU CAN DO IT!===")
	}

	if player.HP == 0 {
		str.print("===YOU DIED===")
	}

	if monster.HP < 10 && monster.HP >= 5 {
		str.print(sprintf("===%s is bleeding!===", monster.Name))
	}

	if monster.HP < 5 && monster.HP > 0 {
		str.print(sprintf("===%s is dying!===", monster.Name))
	}

	if monster.HP <= 0 {
		str.print(sprintf("===%s is dead!===", monster.Name))
	}
}

func main () {
	var player Player
	player.Name = "Mario"
	player.HP = 10
	player.Mana = 10
	player.Lives = 3

	var monster Monster
	monster.Name = "Domo-kun"
	monster.HP = 7
	monster.Mana = 4

	player.draw()
	monster.draw()

	for true {
		if player.HP < 1 || monster.HP < 1 {
			return
		}

		printf("Command? (M)agic; (P)hysical; (E)xit\t")
		var cmd str
		cmd = read()

		if cmd == "E" || cmd == "e" {
			return
		}

		player.draw()
		monster.draw()

		player.attack(cmd, &monster)
		monster.attack(cmd, &player)
		battleStatus(player, monster)
	}
}
```

Lastly, we can see how we use a *for* loop to create something similar
to a REPL for the program that we have been building in the last few sections.

### Go-to
[[Back to the Table of Contents] ↑](#table-of-contents)

The last control flow mechanism is *go-to*, which is achieved through
the `goto` statement.


```
package main

func main () (out i32) {
beginning:
	printf("What animal do you like the most: (C)at; (D)og; (P)igeon\n")

	var cmd str
	cmd = read()

	if cmd == "C" || cmd == "c" {
		goto cat
	}

	if cmd == "D" || cmd == "d" {
		goto dog
	}

	if cmd == "P" || cmd == "p" {
		goto pigeon
	}

cat:
	str.print("meow")
	goto beginning
dog:
	str.print("woof")
	goto beginning
pigeon:
	str.print("tweet")
	goto beginning
}
```

The program above creates an infinite loop by using `goto`s. The loop
will keep asking the user to input commands, and will jump to certain
expression depending on the command.

## Affordances
[[Back to the Table of Contents] ↑](#table-of-contents)

![CX Affordances](documentation/images/affordances.png)

# Native Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

## Type-inferenced Functions

CX has a small set of functions that are not associated to a single
type signature. For example, instead of using `i32.add` to add two
32-bit integers, you can use the generalized `add`
function. Furthermore, whenever you use arithmetic operators, such as
`+`, `-` or `%`, these are translated to their corresponding
"type-inferenced" function, e.g. `num = 5 + 5` is translated to `num =
i32.add(5, 5)`. These native functions still follow CX's philosophy of
having a strict typing system, as the types of the arguments sent to
these native functions must be the same.

Note that after listing a group of similar "type-inferenced" functions
below, we list the compatible types for the corresponding functions.

### `eq`
### `uneq`

Note: the preceding functions only work with arguments of type
*bool*, *byte*, *str*, *i32*, *i64*, *f32* or *f64*.

#### Example

```
package main

func main () {
	bool.print(eq(5, 5))
	bool.print(5 == 5) // alternative

	bool.print(uneq("hihi", "byebye"))
	bool.print("hihi" != "byebye") // alternative
}
```

### `lt`
### `gt`
### `lteq`
### `gteq`

Note: the preceding function only works with arguments of type *byte*,
*bool*, *str*, *i32*, *i64*, *f32* or *f64*.

#### Example

```
package main

func main () {
	bool.print(lt(3B, 4B))
	bool.print(3B < 4B) // alternative

	bool.print(gt("hello", "hi!"))
	bool.print("hello" > "hi!") // alternative

	bool.print(lteq(5.3D, 5.3D))
	bool.print(5.3D <= 5.3D) // alternative

	bool.print(gteq(10L, 3L))
	bool.print(10L >= 3L) // alternative
}
```

### `bitand`
### `bitor`
### `bitxor`
### `bitclear`
### `bitshl`
### `bitshr`

Note: the preceding functions only work with arguments of type *i32*
or *i64*.

#### Example

```
package main

func main () {
	i32.print(bitand(5, 1))
	i32.print(5 & 1) // alternative

	i64.print(bitor(3L, 2L))
	i64.print(3L | 2L) // alternative

	i32.print(bitxor(10, 2))
	i32.print(10 ^ 2) // alternative

	i64.print(bitclear(5L, 2L))
	i64.print(5L &^ 2L) // alternative

	i32.print(bitshl(2, 3))
	i32.print(2 << 3) // alternative

	i32.print(bitshr(16, 3))
	i32.print(16 >> 3) // alternative
}
```

### `add`
### `sub`
### `mul`
### `div`

Note: the preceding functions only work with arguments of type
*byte*, *i32*, *i64*, *f32* or *f64*.

#### Example

```
package main

func main () {
	byte.print(add(5B, 10B))
	byte.print(5B + 10B) // alternative

	i32.print(sub(3, 7))
	i32.print(3 - 7) // alternative

	i64.print(mul(4L, 5L))
	i64.print(4L * 5L) // alternative

	f32.print(div(4.3, 2.1))
	f32.print(4.3 / 2.1) // alternative
}
```

### `mod`
[[Back to the Table of Contents] ↑](#table-of-contents)

Note: the preceding function only works with arguments of type *byte*,
*i32* or *i64*.

#### Example

```
package main

func main () {
	byte.print(mod(5B, 3B))
	byte.print(5B % 3B) // alternative
}
```

### `len`
[[Back to the Table of Contents] ↑](#table-of-contents)

Note: the preceding function only works with arguments of type *str*,
*arrays* or *slices*.

#### Example

```
package main

func main () {
	var string str
	var array [5]i32
	var slice []i32

	string = "this should print 20"
	array = [5]i32{1, 2, 3, 4, 5}
	slice = []i32{10, 20, 30}

	i32.print(len(string)) // prints 20
	i32.print(len(array)) // prints 5
	i32.print(len(slice)) // prints 3
}
```

### `printf`
[[Back to the Table of Contents] ↑](#table-of-contents)

Note: the preceding function requires a format *str* as its first
argument, followed by any number of arguments of type *str*, *i32*,
*i64*, *f32* or *f64*. The format string recognizes the following
directives: `%s` for strings, `%d` for integers and `%f` for floating
point numbers.

#### Example

```
package main

func main () {
	var name str
	var age i32
	var wrongPI f32
	var error f64

	name = "Richard"
	age = 14
	wrongPI = 3.16
	error = 0.0000000000000001D

	printf("Hello, %s. My name is %s. I see that you calculated the value of PI wrong (%f). I think this is not so bad, considering your young age of %d. When I was %d years old, I remember I miscalculated it, too (I got %f as a result, using a numerical method). If you are using a numerical method, please consider reaching an error lower than %f to get an acceptable result, and not a ridiculous value such as %f. \n\nBest regards!\n", name, "Edward", wrongPI, age, 25, 3.1417, error, 0.1)
}
```

### `sprintf`
[[Back to the Table of Contents] ↑](#table-of-contents)

Note: the preceding function requires a format *str* as its first
argument, followed by any number of arguments of type *str*, *i32*,
*i64*, *f32* or *f64*. The format string recognizes the following
directives: `%s` for strings, `%d` for integers and `%f` for floating
point numbers.

```
package main

func main () {
	var reply str
	var name str
	var title str

	name = "Edward"
	title = "Richard 8 PI"

	reply = sprintf("Thank you for contacting our technical support, %s. We see that you are having trouble with our video game titled '%s', targetted to kids under age %d. If you provide us with your parents e-mail address, we'll be glad to help you!", name, title, 14)

	str.print(reply)
}
```

## Slice Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

### `append`

#### Example

```
package main

func main () {
	var slc1 []i32
	slc1 = append(slc1, 1)
	slc1 = append(slc1, 2)

	var slc2 []i32
	slc2 = append(slc1, 3)
	slc2 = append(slc2, 4)

	i32.print(len(slc1)) // prints 2
	i32.print(len(slc2)) // prints 4
}
```

## Input/Output Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

The following functions are used to handle input from the user and to
print output to a terminal.

### `read`

#### Example

```
package main

func main () {
	var password str

	for true {
		printf("What's the password, kid? ")
		password = read()

		if password == "123" {
			str.print("Welcome back.")
			return
		} else {
			str.print("Wrong, but you'll get another chance.")
		}
	}
}
```

### `byte.print`
### `bool.print`
### `str.print`
### `i32.print`
### `i64.print`
### `f32.print`
### `f64.print`
### `printf`

#### Example

```
package main

func main () {
    byte.print(5B)
    bool.print(true)
    str.print("Hello!")
    i32.print(5)
    i64.print(5L)
    f32.print(5.0)
    f64.print(5.0D)
    printf("For a better example, check section Type-inferenced Functions'")
}
```

## Parse Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

All parse functions follow the same pattern: `XXX.YYY` where *XXX* is
the receiving type and *YYY* is the target type. You can read these
functions as "parse XXX to YYY".

### `byte.str`
### `byte.i32`
### `byte.i64`
### `byte.f32`
### `byte.f64`


#### Example

```
package main

func main () {
	var b byte
	b = 30B

	str.print(str.concat("Hello, ", byte.str(b)))
	i32.print(5 + byte.i32(b))
	i64.print(10L + byte.i64(b))
	f32.print(33.3 + byte.f32(b))
	f64.print(50.111D + byte.f64(b))
}
```

### `i32.byte`
### `i32.str`
### `i32.i64`
### `i32.f32`
### `i32.f64`

#### Example

```
package main

func main () {
	var num i32
	num = 43

	str.print(str.concat("Hello, ", i32.str(num)))
	byte.print(5B + i32.byte(num))
	i64.print(10L + i32.i64(num))
	f32.print(33.3 + i32.f32(num))
	f64.print(50.111D + i32.f64(num))
}
```

### `i64.byte`
### `i64.str`
### `i64.i32`
### `i64.f32`
### `i64.f64`

#### Example

```
package main

func main () {
	var num i64
	num = 43L

	str.print(str.concat("Hello, ", i64.str(num)))
	byte.print(5B + i64.byte(num))
	i64.print(10L + i64.i64(num))
	f32.print(33.3 + i64.f32(num))
	f64.print(50.111D + i64.f64(num))
}
```

### `f32.byte`
### `f32.str`
### `f32.i32`
### `f32.i64`
### `f32.f64`

#### Example

```
package main

func main () {
	var num f32
	num = 43.33

	str.print(str.concat("Hello, ", f32.str(num)))
	byte.print(5B + f32.byte(num))
	i32.print(33 + f32.f32(num))
	i64.print(10L + f32.i64(num))
	f64.print(50.111D + f32.f64(num))
}
```

### `f64.byte`
### `f64.str`
### `f64.i32`
### `f64.i64`
### `f64.f32`

#### Example

```
package main

func main () {
	var num f64
	num = 43.33D

	str.print(str.concat("Hello, ", f64.str(num)))
	byte.print(5B + f64.byte(num))
	i32.print(33 + f64.f32(num))
	i64.print(10L + f64.i64(num))
	f32.print(50.111 + f64.f32(num))
}
```

### `str.byte`
### `str.i32`
### `str.i64`
### `str.f32`
### `str.f64`

#### Example

```
package main

func main () {
	var num str
	num = "33"

	byte.print(5B + str.byte(num))
	i32.print(33 + str.f32(num))
	i64.print(10L + str.i64(num))
	f32.print(50.111 + str.f32(num))
	f64.print(50.111D + str.f32(num))
}
```

## Unit Testing
[[Back to the Table of Contents] ↑](#table-of-contents)

The `assert` function is used to test the value of an expression
against another value. This function is useful to test that a
package is working as intended.

### assert

#### Example

```
package main

func foo() (res str) {
    res = "Working well"
}

func main () {
    var results []bool

    results = append(results, assert(5 + 5, 10, "Something went wrong with 5 + 5"))
    results = append(results, assert(foo(), "Working well", "Something went wrong with foo()"))

    var successfulTests i32
	for c := 0; c < len(results); c++ {
		if results[c] {
			successfulTests = successfulTests + 1
		}
	}

    printf("%d tests were performed\n", len(results))
    printf("%d were successful\n", successfulTests)
    printf("%d failed\n", len(results) - successfulTests)
}
```

## `bool` Type Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

### `bool.print`
### `bool.eq`
### `bool.uneq`
### `bool.not`
### `bool.or`
### `bool.and`

#### Example

```
package main

func main () {
	bool.print(bool.eq(true, true))
	bool.print(bool.uneq(false, true))
	bool.print(bool.not(false))
	bool.print(bool.or(true, false))
	bool.print(bool.and(true, true))
}
```

## `str` Type Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

### `str.print`
### `str.concat`

#### Example

```
package main

func main () {
	str.print(str.concat("Hello, ", "World!"))
}
```

## `i32` Type Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

The following functions are of general purpose and are restricted to
work with data structures of type *i32* where it makes sense.

### `i32.print`
### `i32.add`
### `i32.sub`
### `i32.mul`
### `i32.div`
### `i32.mod`
### `i32.abs`

#### Example

```
package main

func main () {
	i32.print(i32.add(5, 7))
	i32.print(i32.sub(6, 3))
	i32.print(i32.mul(4, 8))
	i32.print(i32.div(15, 3))
	i32.print(i32.mod(5, 3))
	i32.print(i32.abs(-13))
}
```

### `i32.log`
### `i32.log2`
### `i32.log10`
### `i32.pow`
### `i32.sqrt`

#### Example

```
package main

func main () {
	i32.print(i32.log(13))
	i32.print(i32.log2(3))
	i32.print(i32.log10(12))
	i32.print(i32.pow(4, 4))
	i32.print(i32.sqrt(2))
}
```

### `i32.gt`
### `i32.gteq`
### `i32.lt`
### `i32.lteq`
### `i32.eq`
### `i32.uneq`

#### Example

```
package main

func main () {
	bool.print(i32.gt(5, 3))
	bool.print(i32.gteq(3, 8))
	bool.print(i32.lt(4, 3))
	bool.print(i32.lteq(8, 6))
	bool.print(i32.eq(-9, -9))
	bool.print(i32.uneq(3, 3))
}
```

### `i32.bitand`
### `i32.bitor`
### `i32.bitxor`
### `i32.bitclear`
### `i32.bitshl`
### `i32.bitshr`

#### Example

```
package main

func main () {
	i32.print(i32.bitand(2, 5))
	i32.print(i32.bitor(8, 3))
	i32.print(i32.bitxor(3, 9))
	i32.print(i32.bitclear(4, 4))
	i32.print(i32.bitshl(5, 9))
	i32.print(i32.bitshr(1, 6))
}
```

### `i32.max`
### `i32.min`

#### Example

```
package main

func main () {
	i32.print(i32.max(2, 5))
	i32.print(i32.min(10, 3))
}
```

### `i32.rand`

#### Example

```
package main

func main () {
	i32.print(i32.rand(0, 100))
}
```

## `i64` Type Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

The following functions are of general purpose and are restricted to
work with data structures of type *i64* where it makes sense.

### `i64.print`
### `i64.add`
### `i64.sub`
### `i64.mul`
### `i64.div`
### `i64.mod`
### `i64.abs`

#### Example

```
package main

func main () {
	i64.print(i64.add(5L, 7L))
	i64.print(i64.sub(6L, 3L))
	i64.print(i64.mul(4L, 8L))
	i64.print(i64.div(15L, 3L))
	i64.print(i64.mod(5L, 3L))
	i64.print(i64.abs(-13L))
}
```

### `i64.log`
### `i64.log2`
### `i64.log10`
### `i64.pow`
### `i64.sqrt`

#### Example

```
package main

func main () {
	i64.print(i64.log(13L))
	i64.print(i64.log2(3L))
	i64.print(i64.log10(12L))
	i64.print(i64.pow(4L, 4L))
	i64.print(i64.sqrt(2L))
}
```

### `i64.gt`
### `i64.gteq`
### `i64.lt`
### `i64.lteq`
### `i64.eq`
### `i64.uneq`

#### Example

```
package main

func main () {
	bool.print(i64.gt(5L, 3L))
	bool.print(i64.gteq(3L, 8L))
	bool.print(i64.lt(4L, 3L))
	bool.print(i64.lteq(8L, 6L))
	bool.print(i64.eq(-9L, -9L))
	bool.print(i64.uneq(3L, 3L))
}
```

### `i64.bitand`
### `i64.bitor`
### `i64.bitxor`
### `i64.bitclear`
### `i64.bitshl`
### `i64.bitshr`

#### Example

```
package main

func main () {
	i64.print(i64.bitand(2L, 5L))
	i64.print(i64.bitor(8L, 3L))
	i64.print(i64.bitxor(3L, 9L))
	i64.print(i64.bitclear(4L, 4L))
	i64.print(i64.bitshl(5L, 9L))
	i64.print(i64.bitshr(1L, 6L))
}
```

### `i64.max`
### `i64.min`

#### Example

```
package main

func main () {
	i64.print(i64.max(2L, 5L))
	i64.print(i64.min(10L, 3L))
}
```

### `i64.rand`

#### Example

```
package main

func main () {
	i64.print(i64.rand(0L, 100L))
}
```

## `f32` Type Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

The following functions are of general purpose and are restricted to
work with data structures of type *f32* where it makes sense.

### `f32.print`
### `f32.add`
### `f32.sub`
### `f32.mul`
### `f32.div`
### `f32.abs`

#### Example

```
package main

func main () {
	f32.print(f32.add(5.3, 10.5))
	f32.print(f32.sub(3.2, 6.7))
	f32.print(f32.mul(-7.9, -7.1))
	f32.print(f32.div(10.3, 2.4))
	f32.print(f32.abs(-3.14159))
}
```

### `f32.log`
### `f32.log2`
### `f32.log10`
### `f32.pow`
### `f32.sqrt`

#### Example

```
package main

func main () {
	f32.print(f32.log(2.3))
	f32.print(f32.log2(3.4))
	f32.print(f32.log10(3.0))
	f32.print(f32.pow(-5.3, 2.0))
	f32.print(f32.sqrt(4.0))
}
```

### `f32.sin`
### `f32.cos`

#### Example

```
package main

func main () {
	f32.print(f32.sin(1.0))
	f32.print(f32.cos(2.0))
}
```

### `f32.gt`
### `f32.gteq`
### `f32.lt`
### `f32.lteq`
### `f32.eq`
### `f32.uneq`

#### Example

```
package main

func main () {
	bool.print(f32.gt(5.3, 3.1))
	bool.print(f32.gteq(3.7, 1.9))
	bool.print(f32.lt(2.4, 5.5))
	bool.print(f32.lteq(8.4, 3.2))
	bool.print(f32.eq(10.3, 10.3))
	bool.print(f32.uneq(8.9, 3.3))
}
```

### `f32.max`
### `f32.min`

#### Example

```
package main

func main () {
	f32.print(f32.max(3.3, 4.2))
	f32.print(f32.min(5.8, 9.9))
}
```

## `f64` Type Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

The following functions are of general purpose and are restricted to
work with data structures of type *f64* where it makes sense.

### `f64.print`
### `f64.add`
### `f64.sub`
### `f64.mul`
### `f64.div`
### `f64.abs`

#### Example

```
package main

func main () {
	f64.print(f64.add(5.3D, 10.5D))
	f64.print(f64.sub(3.2D, 6.7D))
	f64.print(f64.mul(-7.9D, -7.1D))
	f64.print(f64.div(10.3D, 2.4D))
	f64.print(f64.abs(-3.14159D))
}
```

### `f64.log`
### `f64.log2`
### `f64.log10`
### `f64.pow`
### `f64.sqrt`

#### Example

```
package main

func main () {
	f64.print(f64.log(2.3D))
	f64.print(f64.log2(3.4D))
	f64.print(f64.log10(3.0D))
	f64.print(f64.pow(-5.3D, 2.0D))
	f64.print(f64.sqrt(4.0D))
}
```

### `f64.sin`
### `f64.cos`

#### Example

```
package main

func main () {
	f64.print(f64.sin(1.0D))
	f64.print(f64.cos(2.0D))
}
```

### `f64.gt`
### `f64.gteq`
### `f64.lt`
### `f64.lteq`
### `f64.eq`
### `f64.uneq`

#### Example

```
package main

func main () {
	bool.print(f64.gt(5.3D, 3.1D))
	bool.print(f64.gteq(3.7D, 1.9D))
	bool.print(f64.lt(2.4D, 5.5D))
	bool.print(f64.lteq(8.4D, 3.2D))
	bool.print(f64.eq(10.3D, 10.3D))
	bool.print(f64.uneq(8.9D, 3.3D))
}
```

### `f64.max`
### `f64.min`

#### Example

```
package main

func main () {
	f64.print(f64.max(3.3D, 4.2D))
	f64.print(f64.min(5.8D, 9.9D))
}
```

## `time` Package Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

The functions in the `time` package deal with real-time in your
programs. They are used to measure and stop time. Note that in order
to use these functions you need to import the `time` package.

### `time.Sleep`
### `time.UnixMilli`
### `time.UnixNano`

#### Example

```
package main
import "time"

func main () {
	var start i64
	var end i64

	start = time.UnixMilli()
	time.Sleep(1000)
	end = time.UnixMilli()

	printf("elapsed time in milliseconds: \t%d\n", end - start)

	start = time.UnixNano()
	time.Sleep(1000)
	end = time.UnixNano()

	printf("elapsed time in nanoseconds: \t%d\n", end - start)
}
```

## `os` Package Functions

The `os` package provides functions that serve as an interface to CX's
underlaying operating system.

### `os.GetWorkingDirectory`
### `os.Open`
### `os.Close`

#### Example

```
package main
import "os"

func main () {
	var wd str
	wd = os.GetWorkingDirectory()

	var fileName str
	fileName = str.concat(wd, "testing.cx")

	os.Open(fileName)
	os.Close(fileName)
}
```

## `gl` Package Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

"OpenGL is the premier environment for developing portable,
interactive 2D and 3D graphics applications. Since its introduction in
1992, OpenGL has become the industry's most widely used and supported
2D and 3D graphics application programming interface (API), bringing
thousands of applications to a wide variety of computer
platforms. OpenGL fosters innovation and speeds application
development by incorporating a broad set of rendering, texture
mapping, special effects, and other powerful visualization
functions. Developers can leverage the power of OpenGL across all
popular desktop and workstation platforms, ensuring wide application
deployment." This description was extracted from OpenGL's website
(https://www.opengl.org/).


### `gl.ActiveTexture`
### `gl.AttachShader`
### `gl.Begin`
### `gl.BindAttribLocation`
### `gl.BindBuffer`
### `gl.BindFramebuffer`
### `gl.BindTexture`
### `gl.BindVertexArray`
### `gl.BlendFunc`
### `gl.BufferData`
### `gl.BufferSubData`
### `gl.CheckFramebufferStatus`
### `gl.ClearColor`
### `gl.ClearDepth`
### `gl.Clear`
### `gl.Color3f`
### `gl.Color4f`
### `gl.CompileShader`
### `gl.CreateProgram`
### `gl.CreateShader`
### `gl.CullFace`
### `gl.DeleteBuffers`
### `gl.DeleteFramebuffers`
### `gl.DeleteProgram`
### `gl.DeleteShader`
### `gl.DeleteTextures`
### `gl.DeleteVertexArrays`
### `gl.DepthFunc`
### `gl.DepthMask`
### `gl.DetachShader`
### `gl.Disable`
### `gl.DrawArrays`
### `gl.EnableClientState`
### `gl.EnableVertexAttribArray`
### `gl.Enable`
### `gl.End`
### `gl.FramebufferTexture2D`
### `gl.Free`
### `gl.Frustum`
### `gl.GenBuffers`
### `gl.GenFramebuffers`
### `gl.GenTextures`
### `gl.GenVertexArrays`
### `gl.GetAttribLocation`
### `gl.GetError`
### `gl.GetShaderiv`
### `gl.GetTexLevelParameteriv`
### `gl.Hint`
### `gl.Init`
### `gl.Lightfv`
### `gl.LinkProgram`
### `gl.LoadIdentity`
### `gl.MatrixMode`
### `gl.NewTexture`
### `gl.Normal3f`
### `gl.Ortho`
### `gl.PopMatrix`
### `gl.PushMatrix`
### `gl.Rotatef`
### `gl.Scalef`
### `gl.ShaderSource`
### `gl.Strs`
### `gl.TexCoord2d`
### `gl.TexCoord2f`
### `gl.TexEnvi`
### `gl.TexImage2D`
### `gl.TexParameteri`
### `gl.Translatef`
### `gl.Uniform1f`
### `gl.Uniform1i`
### `gl.UseProgram`
### `gl.Vertex2f`
### `gl.Vertex3f`
### `gl.VertexAttribPointerI32`
### `gl.VertexAttribPointer`
### `gl.Viewport`
### `gl.getUniformLocation`

## `glfw` Package Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

"GLFW is an Open Source, multi-platform library for OpenGL, OpenGL ES
and Vulkan development on the desktop. It provides a simple API for
creating windows, contexts and surfaces, receiving input and events."
This description was extracted from GLFW's website (https://www.glfw.org/).


### `glfw.CreateWindow`
### `glfw.GetCursorPos`
### `glfw.GetFramebufferSize`
### `glfw.GetTime`
### `glfw.Init`
### `glfw.MakeContextCurrent`
### `glfw.PollEvents`
### `glfw.SetCursorPosCallback`
### `glfw.SetInputMode`
### `glfw.SetKeyCallback`
### `glfw.SetMouseButtonCallback`
### `glfw.SetShouldClose`
### `glfw.SetWindowPos`
### `glfw.ShouldClose`
### `glfw.SwapBuffers`
### `glfw.SwapInterval`
### `glfw.WindowHint`

## `gltext` Package Functions
[[Back to the Table of Contents] ↑](#table-of-contents)

"The gltext package offers a simple set of text rendering utilities
for OpenGL programs. It deals with TrueType and Bitmap (raster)
fonts. Text can be rendered in various directions (Left-to-right,
right-to-left and top-to-bottom). This allows for correct display of
text for various languages." This description was extracted from
gltext's website (https://github.com/go-gl/gltext).


The `gltext` functions can be used to display character strings on
windows. Different fonts can be used by loading font files using `gltext.LoadTrueType`.

### `gltext.GlyphBounds () (i32, i32)`
### `gltext.LoadTrueType (str, str, i32, i32, i32, i32) ()`
### `gltext.Metrics (str, str) (i32, i32)`
### `gltext.NextRune (str, str, i32) (i32, i32, i32, i32, i32, i32, i32)`
### `gltext.Printf`
### `gltext.Texture`