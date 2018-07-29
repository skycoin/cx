![cx logo](https://user-images.githubusercontent.com/26845312/32426758-2a4bbb00-c282-11e7-858e-a1eaf3ea92f3.png)

Table of Contents
=================

   * [Table of Contents](#table-of-contents)
   * [CX Programming Language](#cx-programming-language)
   * [CX Playground](#cx-playground)
   * [Changelog](#changelog)
   * [Installation](#installation)
      * [Installing Go](#installing-go)
   * [Additional Notes Before the Actual Installation](#additional-notes-before-the-actual-installation)
      * [Linux: Installing OpenGL and GLFW Dependencies](#linux-installing-opengl-and-glfw-dependencies)
         * [Debian-based Linux Distributions](#debian-based-linux-distributions)
      * [Windows: Installing GCC](#windows-installing-gcc)
      * [Installing CX - Method 1: The "so easy it might not work" Solution](#installing-cx---method-1-the-so-easy-it-might-not-work-solution)
         * [Windows](#windows)
      * [Installing CX - Method 2: The "not so easy, but still easy" Solution](#installing-cx---method-2-the-not-so-easy-but-still-easy-solution)
      * [Updating CX](#updating-cx)
   * [Running CX](#running-cx)
      * [CX REPL](#cx-repl)
      * [Running CX Programs](#running-cx-programs)
      * [Other Options](#other-options)
   * [CX Tutorial](#cx-tutorial)
   * [Hello World](#hello-world)
   * [Shorthands to Native Functions](#shorthands-to-native-functions)
      * [Relational Operators](#relational-operators)
      * [Logical Operators](#logical-operators)
      * [Arithmetic Operators](#arithmetic-operators)
      * [Arithmetic Assignment Operators](#arithmetic-assignment-operators)
   * [Comments](#comments)
   * [Data](#data)
      * [Atomic Data](#atomic-data)
         * [Booleans](#booleans)
         * [Integers](#integers)
         * [Floats](#floats)
         * [Bytes](#bytes)
         * [Strings](#strings)
      * [Arrays](#arrays)
         * [Multidimensional Arrays](#multidimensional-arrays)
      * [Variables](#variables)
         * [Local Variables](#local-variables)
         * [Global Variables](#global-variables)
      * [Structs](#structs)
   * [Expressions](#expressions)
   * [Flow Control](#flow-control)
      * [If and if/else](#if-and-ifelse)
      * [For Loop](#for-loop)
      * [Go-to](#go-to)
      * [Return](#return)
   * [Functions](#functions)
   * [Methods](#methods)
   * [Packages](#packages)
      * [Creating and Importing Packages](#creating-and-importing-packages)
      * [Working with Different Files](#working-with-different-files)
   * [Debugging](#debugging)
      * [Halt](#halt)
      * [Unit Testing](#unit-testing)
   * [Affordances](#affordances)
      * [Limiting the Affordance System's Search Space](#limiting-the-affordance-systems-search-space)
   * [Experimental Features](#experimental-features)
      * [Evolutionary Algorithm](#evolutionary-algorithm)
      * [Serialization](#serialization)
      * [OpenGL 1.2 API](#opengl-12-api)

# CX Programming Language

CX is a general purpose, interpreted and compiled programming
language, with a very strict type system and a syntax
similar to Golang's. CX provides a new programming paradigm based on
the concept of affordances, where the user can ask the programming
language at runtime what can be done with a CX object (functions,
expressions, packages, etc.), and interactively or automatically choose
one of the affordances to be applied.

# CX Playground

If you want to test some CX examples, you can do it in the CX
[Playground (http://cx.skycoin.net/)](http://cx.skycoin.net/).

# Changelog

Check out the latest additions and bug fixes in the [changelog](https://github.com/skycoin/cx/blob/master/CHANGELOG.md).

# Installation

CX has been successfully installed and tested in recent versions of
Linux (Ubuntu), MacOS X and Windows. Nevertheless, if you run into any
problems, please create an issue and we'll try to solve the problem as
soon as possible.

## Installing Go

First, make sure that you have Go installed by running `go
version`. It should output something similar to:

```
go version go1.8.3 darwin/amd64
```

**You need a version greater than 1.8, and >1.10 is recommended**

Some linux distros' package managers install very old versions of
Go. You can try first with a binary from your favorite package
manager, but if the installation starts showing errors, try with the
latest version before creating an issue.

Go should also be properly configured (you can read the installation
instructions by clicking [here](https://golang.org/doc/install). Particularly:

* Make sure that you have added the Go binary to your `$PATH`.
  * If you installed Go using a package manager, the Go binary is most
    likely already in your `$PATH` variable.
  * If you already installed Go, but running "go" in a terminal throws
    a "command not found" error, this is most likely the problem.
* Make sure that you have configured your `$GOPATH` environment
variable.
* Make sure that you have added `$GOPATH/bin` to you `$PATH`.
  * If you have binaries installed in `$GOPATH/bin` but you can't use
    them by just typing their name wherever you are in the file system
    in a terminal, then this will solve the problem.

As an example configuration, considering you're using *bash* in
*Ubuntu*, you would append to your `~/.bashrc` file this:

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin/
```

Don't just copy/paste that; think on what you're doing!

# Additional Notes Before the Actual Installation

## Linux: Installing OpenGL and GLFW Dependencies

### Debian-based Linux Distributions

\* Based on instructions from [Viscript](https://github.com/skycoin/viscript)'s repository.

CX comes with OpenGL and GLFW APIs. In order to use them, you need to
install some dependencies. If you're using a Debian based Linux
distro, such as Ubuntu, you can run these commands:

```
sudo apt-get install libxi-dev
sudo apt-get install libgl1-mesa-dev
sudo apt-get install libxrandr-dev
sudo apt-get install libxcursor-dev
sudo apt-get install libxinerama-dev
```

and you should be ready to go.

## Windows: Installing GCC

You might need to install GCC. Try installing everything first
without installing GCC, and if an error similar to "gcc: command not
found" is shown, you can fix this by installing MinGW.

Don't get GCC through Cygwin; apparently, [Cygwin has compatibility
issues with Go](https://github.com/golang/go/issues/7265#issuecomment-66091041).

Users have reported that using either [MingW](http://www.mingw.org/)
or [tdm-gcc](tdm-gcc.tdragon.net(), where tdm-gcc seems to be the
easiest way.

## Installing CX

Make sure that you have `curl` and `git` installed. Run this command in a terminal:

```
sh <(curl -s https://raw.githubusercontent.com/skycoin/cx/master/cx.sh)
```

If you're skeptical about what this command does, you can check the
source code in this project. Basically, this script checks if you have
all the necessary Golang packages and tries to install them for
you. The script even downloads this repository and installs CX for
you. This means that you can run `cx` after running the script and see
the REPL right away (if the script worked). To exit the REPL, you can press `Ctrl-D`.

You should test your installation by running `cx
$GOPATH/src/github.com/skycoin/cx/tests`.

As an alternative, you could clone into this repository and run cx.sh
in a terminal.

### Windows

An installation script is also provided for Windows named cx.bat.
The Windows version of this method would be to manually
download the provided [batch script](https://github.com/skycoin/cx/blob/master/cx.bat) (which is similar to the bash
script for *nix systems described above), and run it in a terminal.

You should test your installation by running `cx
%GOPATH%\src\github.com\skycoin\cx\tests`.

## Updating CX

Now you can update CX by simply running the installation script
again:

```
./cx.sh
```

or, in Windows:

```
cx.bat
```

# Running CX
## CX REPL

Once CX has been successfully installed, running `cx` should print
this in your terminal:

```
CX 0.5.6
More information about CX is available at http://cx.skycoin.net/ and https://github.com/skycoin/cx/

:func main {...
	*
```

This is the CX REPL
([read-eval-print loop](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop)),
where you can debug and modify CX programs. The CX REPL starts with a
barebones CX structure (a `main` package and a `main` function) that
you can use to start building a program.

Let's create a small program to test the REPL. First, write
`str.print("Testing the REPL")` after the `*`, and press enter. After
pressing enter you'll see the message "Testing the REPL" on the screen. If you
then write `:dp` (short for `:dProgram` or *debug program*), you
should get the current program AST printed:

```
Program
0.- Module: main
	Functions
		0.- Function: main () ()
			0.- Expression: str.print("Testing the REPL" str)
```

As we can see, we have a `main` package, a `main` function, and we
have a single expression: `str.print("Testing the REPL")`.

Let's now create a new function. In order to do this, we first need to
leave the `main` function. At this moment, any expression (or function
call) that we add to our program is going to be added to `main`. To
exit a function declaration, press `Ctrl+D`. The prompt (`*`) should
have changed indentation, and the REPL now shouldn't print `:func main
{...` above the prompt:

```
:func main {...
	* 

* 
```

Now, let's enter a function prototype (an empty function which only
specifies the name, the inputs and the outputs):

```
* func sum (num1 i32, num2 i32) (num3 i32) {}

* 
```

You can check that the function was indeed added by issuing a `:dp`
command. If we want to add expressions to `sum`, we have to select it:

```
* :func sum

:func sum {...
	* 
```

Notice that there's a semicolon before `func sum`. Now we can add an expression to it:

```
:func sum {...
	* num3 = num1 + num2
```

Now, exit `sum` and select `main` with the command `:func main`. Let's
add a call to `sum` and print the value that it returns when giving
the arguments 10 and 20:

```
:func main {...
	* i32.print(sum(10, 20))
30
```

## Running CX Programs

To run a CX program, you have to type, for example, `cx
the-program.cx`. Let's try to run some examples from the `examples`
directory in this repository. In a terminal, type this:

```
cd $GOPATH/src/github.com/skycoin/cx/
cx examples/hello-world.cx
```

This should print `Hello World!` in the terminal. Now try running `cx
examples/game.cx`.

## Other Options

If you write `cx --help` or `cx -h`, you should see a text describing
CX's usage, options and more.

Some interesting options are:

* `--base` which generates a CX program's assembly code (in Go)
* `--compile` which generates an executable file
* `--repl` which loads the program and makes CX run in REPL mode
(useful for debugging a program)
* `--web` which starts CX as a RESTful web service (you can send code
  to be evaluated to this endpoint: http://127.0.0.1:5336/eval)

# CX Tutorial

In the following sections, the reader can find a short tutorial on how
to use the main features of the language. It can be used as introductory
material for people with no experience in programming as some sections
are explained as if this was the main audience. Nevertheless, the true
purpose of the tutorial is to demonstrate all the features that the
language currently supports.

Feel free to [create an issue](https://github.com/skycoin/cx/issues)
requesting a better explanation of a feature.

# Hello World

Do you want to know how CX looks? This is how you print "Hello World!"
in a terminal:

```
package main

func main () () {
	str.print("Hello World!")
}
```

Every CX program must have at least a *main* package, and a *main*
function. As mentioned before, CX has a very stricty type system,
where functions can only be associated with a single type
signature. As a consequence,
if we want to print a string, as in the example above, we have to call
*str*'s print function, where *str* is a package containing string
related functions.

# Shorthands to Native Functions

In the following sections you might see function calls like:

```
i32.add(5, 10)
i64.sub(i32.i64(10), i32.i64(5)) 
```

In previous versions of CX, you couldn't write
*infix* operations (e.g. `5 + 10`), but this is not the case
anymore. In other words, the two examples above can now be written as:

```
5 + 10
10 - 5
```

Also, now you don't need to use the cast functions `i32.i64`,
`f32.f64`, etc. If you want to tell the compiler that you want a
number to be interpreted as a *byte*, an *i64* or as an *f64*, you can use the
suffixes `B`, `L` and `D`, respectively:

```
byte.print(55B)
i64.add(31L, 62L)
f64.mul(25D, 25D)
```

or

```
31L + 62L
25D * 25D
```

A bunch of shorthands for are now implemented
Many shorthands to native functions are now implemented. These
shorthands are similar to those present in many other programming
languages, like `&&` (the *and* logical operator) or `>` (the
*greather than* relational operator). The following subsections
present these shorthands.

## Relational Operators

\* These shorthands can also be used with other
data types, such as *i64*, *f32*, *str*, etc., and they can also be
used with identifiers (variables).

| Shorthand 	| Native Function 	|
|:---------:	|:---------------:	|
|   5 == 5  	|   i32.eq(5, 5)  	|
|   5 != 5  	|  i32.uneq(5, 5) 	|
|   5 > 5   	|   i32.gt(5, 5)  	|
|   5 >= 5  	|  i32.gteq(5, 5) 	|
|   5 < 5   	|   i32.lt(5, 5)  	|
|   5 <= 5  	|  i32.lteq(5, 5) 	|

## Logical Operators

\* These shorthands can also be used with identifiers (variables).

|   Shorthand   	| Native Function 	|
|:-------------:	|:---------------:	|
|     !true     	|    bool.not(true)    	|
|  true && true 	| bool.and(true, true) 	|
| true \|\| false 	| bool.or(true, false) 	|

## Arithmetic Operators

\* These shorthands can also be used with other
data types, such as *i64*, *f32*, etc., and they can also be
used with identifiers (variables).

| Shorthand 	|   Native Function  	|
|:---------:	|:------------------:	|
|    5++    	|    i32.add(5, 1)   	|
|    5--    	|    i32.sub(5, 1)   	|
|   5 + 5   	|    i32.add(5, 5)   	|
|   5 - 5   	|    i32.sub(5, 5)   	|
|   5 * 5   	|    i32.mul(5, 5)   	|
|   5 / 5   	|    i32.div(5, 5)   	|
|   5 % 5   	|    i32.mod(5, 5)   	|
|   5 << 5  	|  i32.bitshl(5, 5)  	|
|   5 >> 5  	|  i32.bitshr(5, 5)  	|
|   5 & 5   	|  i32.bitand(5, 5)  	|
|   5 \| 5   	|   i32.bitor(5, 5)  	|
|   5 ^ 5   	|  i32.bitxor(5, 5)  	|
|   5 &^ 5  	| i32.bitclear(5, 5) 	|

## Arithmetic Assignment Operators

\* These shorthands can also be used with other
data types that make sense (for example, there's no `f64.mod` native function), such as *i64*, *f32*, etc., and they can also be
used with identifiers (variables).

| Shorthand 	|      Native Function     	|
|:---------:	|:------------------------:	|
|  foo += 5 	|   foo = i32.add(foo, 5)  	|
|  foo -= 5 	|   foo = i32.sub(foo, 5)  	|
|  foo *= 5 	|   foo = i32.mul(foo, 5)  	|
|  foo /= 5 	|   foo = i32.div(foo, 5)  	|
|  foo %= 5 	|   foo = i32.mod(foo, 5)  	|
| foo <<= 5 	| foo = i32.bitshl(foo, 5) 	|
| foo >>= 5 	| foo = i32.bitshr(foo, 5) 	|
|  foo &= 5 	| foo = i32.bitand(foo, 5) 	|
|  foo \|= 5 	|  foo = i32.bitor(foo, 5) 	|
|  foo ^= 5 	| foo = i32.bitxor(foo, 5) 	|

# Comments

Some of the code snippets that follow have comments in them, i.e.,
blocks of text that are not actually "run" by the CX compiler or
interpreter. Just like in C, Golang and many other programming
languages, single line comments are created by placing double slashes
(//) before the text being commented. For example:

```
// Example of summing two 32 bit integers in CX

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

# Data

Every programming language is designed to manipulate some kind of data
using some kind of process. Let's first have a look at the simplest
kind of data that we can create in CX, and then move on to slightly more
complex data structures.

## Atomic Data
### Booleans

Booleans can be either *true* or *false*, and they are mainly used to
control the flow of a program. As an example, let's print both
possible values to the terminal:

```
bool.print(true)
bool.print(false)
```

### Integers

CX can work with either 32 or 64 bit integers. The types themselves
are called *i32* and *i64*  respectively. Any number without decimal
points are considered to be *i32*, e.g. *5* or *12*.

Unlike i32 numbers, the programmer needs to explicitly tell CX when an
i64 number is required. For example, to print the *i64* number "15",
you'd need to write:

```
i64.print(i32.i64(15))
```

### Floats

Floating-point numbers come in two sizes, just like integers: 32 and
64 bits. 32-bit floats are named *f32*s, while 64-bit floats are named
*f64*s. Similar to *i32*, the programmer does not need any explicit
casting; CX simply regards any number that has a decimal point as an
*f32*. In the case of *f64* numbers, the programmer needs to cast an
*f32* number to *f64* before being passed as an argument to a
function, for example:

```
i64.add(i32.i64(30.0), i32.i64(20.0))
```

### Bytes

Bytes can hold any number from 0 to 255. The programmer can create a
byte number by casting an *i32* to *byte*, or by appending the suffix
B to an integer:

```
i32.byte(100)
byte.print(55B)
```

### Strings

Strings are internally represented as array of bytes. CX's parser
recognizes any chain of characters enclosed by a pair of double quotes
(") as a string:

```
"I'm a string"
"I'm
    also a
            string"
```

A character string in CX is said to be of type *str*.

## Arrays

Until this point, all data types that we have mentioned have been
"atomic," which means that they hold only one piece of information
(except strings, that are actually byte arrays, but they can be seen
as a single unit of information too).

An array is a collection of atomics, where every element
contained in an array must be of the same type.

A programmer can create arrays of each of the atomic types by writing
"[N]", where N is the number of elements for that array, followed by
the desired type and a list of elements that initialize the array,
enclosed in curly braces. For example:


```
[3]bool{true, false, true}
[2]byte{3, 2}
[3]i32{0, 1, 2}
[3]i64{7, 7, 7}
[3]f32{3.5, 1.2, 8.9}
[1]f64{0.3}
[2]str{"hello", "world"}
```

As can be noted, numbers in *[N]i64* and *[N]f64* do not need to be
cast explicitly. The reason behind this is that they are already being
explicitly cast: *[N]i64* and *[N]f64* are telling CX that every element
is of that type.

There are a number of native functions associated to array types. For
example, to obtain the number of elements in an *[]i32*, we can use
*len*:

```
i32.print(len([3]i32{0, 1, 2})
```

We can read and write specific elements from arrays by using the
bracket notation:

```
foo[2] = 33
i32.print(foo[2])
```

### Multidimensional Arrays

If you want to group together arrays, you can create *multidimensional
arrays*. As an example, imagine that we want to hold a tilemap for a
videogame. We could represent the map as a 2-dimensional array which
contains the index of the sprite that needs to be drawn in that tile:


```
1 1 1 0 1 1 1
1 0 0 0 0 0 1
1 2 0 0 2 0 1
1 0 2 0 0 0 1
1 0 0 0 0 0 1
1 1 1 1 1 1 1
```

We can create this map using the following code:

```
var tilemap [6][7]
tilemap[0] = [7]i32{1, 1, 1, 0, 1, 1, 1}
tilemap[0] = [7]i32{1, 0, 0, 0, 0, 0, 1}
tilemap[0] = [7]i32{1, 2, 0, 0, 2, 0, 1}
tilemap[0] = [7]i32{1, 0, 2, 0, 0, 0, 1}
tilemap[0] = [7]i32{1, 0, 0, 0, 0, 0, 1}
tilemap[0] = [7]i32{1, 1, 1, 1, 1, 1, 1}
```

## Variables

Passing literal values to functions give us already a lot of power,
but imagine that we need to pass a 5 string array to three different
functions. In order to avoid writing the same array three times, we
can make a variable hold the value for us, and then pass the variable
to the three different functions:

```
var names [5]str = [5]str{"Edward", "Daniel", "Melissa", "Roger", "Ron"}

notify(names)
saveToDatabase(names)
```

Finally, a variable can be declared and not
*explicitly* initialized. Unlike in languages like C, where a variable
can end up pointing to garbage in memory, every variable in CX is
implicitly initialized to its zero value (unless explicitly
initialized to something else, of course). Numerical variables are
initialized to 0 or 0.0, booleans are initialized to false, strings
are initialized to an empty string (""), and arrays are initialized to
an empty array.

### Local Variables

Variables that are declared inside a function are said to be "local."
This means that other functions do not have access to this variable:

```
func outside () () {
   var greeting str = "hello"
}
func main () () {
  str.print(greeting)
}
```

In this example, the *main* function does not have access to the
*greeting* variable declared in the *outside* function. *greeting* is
a local variable in *outside* function's scope.

Another thing to have in mind is that local variables are not shared
in different calls to the same function. For example:

```
func recur (num i32) () {
  var state i32
  state = i32.add(num, 1)
  if i32.eq(state, 3) {
    return
  } else {
    recur(state)
  }
}

func main () () {
   recur(1)
}
```

In the example above we can see the example of a recursive function: a
function that calls itself. The *main* function calls *recur* and
sends 1 as its argument. *recur* declares a variable called state,
which is defined as the sum of its only argument and 1. If *state* is
equal to 3, the function returns, but if it isn't, it calls itself,
sending its *state* variable as an argument. In this next call to
recur, what's *state* value going to be? Is it going to be 2, the
value obtained in the previous call? No, each function call in CX is
associated with its own "state." The correct answer is 3, which is
obtained by adding 1 to 2, the argument sent to this new call to *recur*.

### Global Variables

What can we do if we want state to be shared among several function
calls? We can use global variables. Global variables are declared
outside of any function declaration, and they must be declared before
any function that plans to use it. If the programmer wants a function
to use a global variable that is declared after a function
declaration, CX should either raise an error or declare a local
variable with that name, depending on the context in which the
variable is trying to be used (being sent as an argument to another
function or assigning a new value to the variable).

```
var counter i32

func inc () () {
  counter = i32.add(counter, 1)
}

func main () () {
  i32.print(counter)
  inc()
  i32.print(counter)
}
```

In the example above, *counter* is defined as a global variable, which
means that any function declared after it, has access to its
value. The *main* function starts its block of expressions by printing
*counter*'s value, which is 0. Then a call to *inc* is performed,
which increases *counter*'s value by 1. *main* then prints again the
value of counter, which should now be 2.

## Structs

We can create groups of variables by using *struct*s. *struct*s are
useful for representing more complex data abstractions, where
different data types and data encapsulations are needed. For example,
if the programmer needs to represent a *Point*, it can be defined in
the following way:

```
type Point struct {
  x i32
  y i32
}
```

The name of the *struct* needs to be surrounded by the keywords *type*
and *struct*. The *struct*'s *fields* are enclosed by curly braces,
and each field is defined by a name or identifier and a type assigned
to that identifier.

Any basic type can be used as the type of a *struct* *field*:

```
type Student struct {
  name str
  age i32
  height f32
}
```

*struct*s can also serve as field types in other *struct*s, but we
need to remember to declare the *struct*s in the correct order:

```
type Color struct {
  r i32
  g i32
  b i32
}

type Point struct {
  x i32
  y i32
}

type Shape struct {
  color Color
  vertices []Point
}
```

If we had declared *Shape* before *Color* or *Point*, CX would raise
an error telling us that type "Color" or type "Point" is not defined.

<!-- As can be noted, as soon as we declare a new type using a *struct*, we -->
<!-- automatically have access to another type: arrays of that type of -->
<!-- *struct*s. CX not only creates this additional type for us, but a set -->
<!-- of functions to manipulate this new array type. Remember, CX is very -->
<!-- strict regarding its type system, so if we want to know what's the -->
<!-- length of an array of *Point*s, we'd need to call *[]Point.len* to -->
<!-- find out: -->

<!-- ``` -->
<!-- var color Color -->
<!-- color.r = 31 -->
<!-- color.g = 23 -->
<!-- color.b = 131 -->

<!-- points := []Point.make(3) -->

<!-- myShape := new Shape{ -->
<!--   color: color, -->
<!--   vertices: points -->
<!-- } -->

<!-- []Point.write(myShape.vertices, 0, new Point{x: 1, y: 2}) -->
<!-- []Point.write(myShape.vertices, 1, new Point{x: 3, y: 5}) -->
<!-- []Point.write(myShape.vertices, 2, new Point{x: 2, y: 7}) -->
<!-- ``` -->

<!-- Woa woa! A lot is happening in the example above. Let's analyze this -->
<!-- step by step. First, we can notice that we can now declare variables -->
<!-- of custom types (*struct*s): we create a *Color* variable named -->
<!-- *color*. Similarly to many other programming languages that are -->
<!-- capable of declaring C structs, we can access the *struct* fields by -->
<!-- using a dot following the variable name, and we tell CX to assign -->
<!-- different values to the r, g, b fields of the *Color* type. -->

<!-- Then we can see how we declare and initialize the *points* variable -->
<!-- by assigning the result of the function call *[]Point.make(3)*. CX -->
<!-- also created a *make* function to initialize arrays of the *Point* -->
<!-- type. Each of the *Point*s in this newly created array are initialized -->
<!-- to its 0 form, i.e., they are *Point* *struct* instances with *x = 0* -->
<!-- and *y = 0*. -->

<!-- Now, in order to create a *Shape* instance, we use the keyword "new" -->
<!-- followed by the name of the *struct* we want to create. This is an -->
<!-- alternative to using the "var name type" form. The advantage to this -->
<!-- new form (pun intended) is that we can use the created literal as an -->
<!-- argument to a function call. Anyway, by using the *new Struct* form, -->
<!-- we can also directly specify what are going to be the values for the -->
<!-- *struct* instance fields. In this case, the instance's field color is -->
<!-- set to the color variable defined above, and the vertices field is set -->
<!-- to the points array, also defined above. -->

<!-- We can see how the empty points are re-assigned by using the -->
<!-- *[]Point.write* function. *myShape.vertices* is sent to *[]Point.write* as the -->
<!-- first argument, which means that we want to write a new value of type -->
<!-- *Point* in the *vertices* array. Each of the three calls writes a new -->
<!-- *Point* literal to each of the available indexes (0, 1 and 3). -->

<!-- As a final note, you can create functions associated to structs, which -->
<!-- are called `methods` (following Go's convention). You can learn more -->
<!-- about CX methods by reading their -->
<!-- [own section here, in this document](#methods). -->

# Expressions

An expression consists of a function to be
called, a set of arguments that are sent to the function, and a set of
receiving variables which will hold the outputs of the function being
called. Everything inside a CX function is an expression or is converted to
expressions.

```
var foo i32
foo = i32.mul(3, 5)
```

<!-- In the example above, we are telling CX to do 3 times 5 and to store -->
<!-- the result in *foo*. *i32.mul* is the function to be called, 3 and 5 -->
<!-- are its arguments, and *foo* is the variable that will store the -->
<!-- output. As in Golang, functions in CX can return more than one value: -->

<!-- ``` -->
<!-- seconds, minutes, hours := getTime() -->
<!-- ``` -->

<!-- Arguments to function calls can be other expressions: -->

<!-- ``` -->
<!-- result := i32.mul(i32.add(5, 2), i32.sub(3, 1)) -->
<!-- ``` -->

<!-- When stating multiple receiving variables, you can provide different -->
<!-- expressions to be assigned to each variable, for example: -->

<!-- ``` -->
<!-- res1, res2 := i32.add(5, 3), i32.sub(10, 7) -->
<!-- i32.print(res1) -->
<!-- i32.print(res2) -->
<!-- ``` -->

<!-- After executing the example above, 8 will be printed followed by 3. -->

# Flow Control

Like in most programming languages (if not all, excluding perhaps some
[esoteric languages](https://en.wikipedia.org/wiki/Esoteric_programming_language)),
programs are executed from top to bottom.

We can control CX's normal flow by using *flow control statements*,
which are discussed below.

## If and if/else

The first *flow control statement* is *if*, which has the capability
of ignoring execution of a number of expressions.

```
if false {
  str.print("This will never print")
}
```

In the example above, we are asking "if x", where x is a predicate
that is evaluated to determine if it's either true or false. If the
predicate is true, then the expressions enclosed in curly brackets are
executed, and are ignored if otherwise. In this case, *x = false*, so
the expression won't be executed. Predicates can be either booleans,
or variables or expressions that evaluate to booleans:


```
if true {
	str.print("This will print")
}

if true || false && false {
	str.print("This will also print")
}
```

If statements can be nested:

```
if true {
  if false {
    str.print("This won't print")
  }
  if true {
    str.print("This will print")
  }
}
```

As in other programming languages, you can also create an *else*
clause. This *else* clause will execute if the predicate evaluates to
false.

```
if false {
  str.print("This won't print")
} else {
  str.print("But this will!")  
}
```

These if/else statements can also be nested:

```
if false {
  str.print("This won't print")
  } else {
  if true {
    if false {
        str.print("This won't print")
      } else {
        str.print("This will print")
      }
    }
}
```

## For Loop

Many programming languages provide several looping constructs, like
while, do-while, for-each, for, etc. CX followed the same strategy as
Golang, and only provides a *for* loop. First, let's have a look at an
infinite loop:

```
for true {
    str.print("This will print forever")
}
```

In this form, the *for* loop behaves similarly to an if statement: it
can receive a single predicate, which can be either a boolean, or a
variable or expression that evaluates to a boolean. For example:

```
for true || false && false {
		str.print("This will also print forever")
}
```

Using this *for* form, we can create a loop that prints the numbers
from 0 to 10 like this:

```
for c := 0; i32.lteq(c, 10); c = i32.add(c, 1) {
  i32.print(c)
}
```

This form is the same as the form provided by many other
programming languages like C. The first part declares and initializes
a variable that will usually serve as a counter, the second part has a
predicate, and the last part is usually used as a counter updater. As
in these other similar programming languages, the first and the third
parts can actually have any expression you want, but they are usually
used initialize counters and update counters, respectively.

## Go-to

The last flow control structure is *go-to*. *go-to*s are used to make
a CX program jump directly to a labelled expression.

```
	goto label4
label3:
	str.print("This should never be reached")
label4:
	str.print("This should be printed")
```

In the example above, the statement "goto label4" will make CX
directly jump to the expression labelled as "label4." This flow
control statement can be combined with other flow control statements
to create complex programs. The downside with *go-to*s is usually that
programs become harder to read, but they can become very powerful in the
correct hands.

## Return

If you want to make a function stop its execution and return to its
calling function, you can use `return`, followed by the arguments that
you want to return as outputs, separated by commas if more than one.

```
func safeDiv (num i32, den i32) (res i32) {
	if den == 0 {
		return
	}
	res = num / den
}
```

`return` can also be used without any arguments, but in this case, you
need to place a semicolon (;) in front of it. If a semicolon is not
present, the CX parser will consider the code in the following lines
to be the function's output:

```
func foo () (num i32) {
	num = 5
	return
	num = i32.add(5, 5)
}
```

The function above will return 10, instead of 5.

# Functions

We have already seen some examples of function calls, and an example
of a function declaration: the *main* function. We can define other
custom functions, though.

```
package main

var PI f32 = 3.14159

func circleArea (radius f32) (area f32) {
	f32.mul(f32.mul(radius, radius), PI)
}
```

In the example above, we are creating a new function that will be
contained in the *main* package. This function is called *circleArea*,
and has one input and one output parameter. If we analyze this
function, we can realize that it calculates the area of a circle,
given its radius. Once defined, it can be called in other functions
that are defined below it, such as the *main* function:

```
func main () {
  var area f32
  area = circleArea(i32.f32(2))
  f32.print(area)
}
```

# Methods

Functions can be associated with structs (custom types) by creating
them as methods. The advantages of methods are that the programmers
are restricted to calling the function only on a variable of the
specified type, and that you can have several functions named the same
way, but that behave differently depending on the variable type that
they are being called on. For example, let's consider the following
two structures:

```
type Point struct {
	x i32
	y i32
}

type Player struct {
	name str
	x i32
	y i32
}
```

Both of these structures could have methods called `Position`, which
print the x and y coordinates of a particular instance:

```
func (pl Player) Position () {
	str.print("This is a player")
	i32.print(pl.x)
    i32.print(pl.y)
}

func (pnt Point) Position () {
	str.print("This is a point")
	i32.print(pnt.x)
    i32.print(pnt.y)
}
```

To call these methods, you just have to write the name of a variable
of the desired type, followed by a period and the name of the method:

```
func main () {
    var myPoint Point
	var myPlayer Player
	myPoint.x = 10
	myPlayer.x = 30

	myPoint.Position()
	myPlayer.Position()
}
```

# Packages

Packages are a useful feature to encapsulate functions, structs and global
variables. We have been explicitly working with a single package since the
beginning of this tutorial: *main*. But we have also used functions in
other packages like *str*, *i32* and *f32*.

Functions from packages like *str* and *i32* are all implicitly imported, which
means that you can start using them right away. The reason behind this
is that these functions are so common that they could be
considered as part of the language's core library of functions.

Big projects usually need to be divided into several
packages. Depending on the programming language, these "packages"
can be called modules, namespaces, vocabularies, etc., but the idea
is always the same: to encapsulate. Encapsulation is very important,
because it allows the programmer to create groups of related
identifiers (variables, types, functions) that won't create conflicts
with other identifiers named the same, but that have a different
purpose. For example, in CX we have several versions of *print*, but
they have somewhat different purposes. Indeed, all of them will print
stuff in the terminal, but *str.print* will print a string,
*i32.print* will print a 32 bit integer, and so on. We could even have
a situation where the functions need to be named the same, and they
even receive the same type of argument, but do entirely different
things. For example, *printer.send* and *fax.send* could both receive
strings, but the former sends the string to a printer and the second
to a fax (if you don't know what a fax is, you can read about them in
[here](https://en.wikipedia.org/wiki/Fax)).

## Creating and Importing Packages

"How can I create a package in CX?" you might be wondering by now.

```
package myPackage
/*
  Functions, structs and globals are placed here
*/
```

You just need to write the keyword *package* followed by the name you
want to give to your package. Unlike Golang, you can have multiple
packages in a single source file and CX won't complain about it
(although you are encouraged to place different packages [in different
files](#working-with-different-files)). Whenever CX reads "package something", every CX statement that
follows will be attached to the "something" package, and this
behaviour will continue until "package somethingElse" is
encountered. Let's create a Math package:

```
package Math

var PI f32 = 3.14159

func double (num f32) (double f32) {
    double = f32.mul(num, 2)
}
```

This starts looking promising, doesn't it? But how can we start using
this package in other packages? You have to *import* it.

```
package main
import "Math"

func main () () {
    i32.print(Math.PI)
    i32.print(Math.double(Math.PI))
}
```

If we don't *import Math*, CX will raise an error telling us that the
module *Math* is not being imported or does not exist.

## Working with Different Files

As mentioned in the previous section, you can place different chunks
of code in different CX files. These files could share the same
package (i.e., they all start with the same `package` declaration),
they could work with different packages, or a combination of these.

If you decide to work with different files, remember to give all of
them as input to the `cx` command. For example:

```
cx file1.cx file2.cx file3.cx
```

In previous versions of CX, the order in which you gave these files as
input was important: if a function needed by `file2.cx` was present in
`file1.cx`, you needed to give `file1.cx` as input before
`file2.cx`. Now this is no longer needed, and you can use whichever
order you want.

# Debugging

Whenever an error is raised in CX, a read-eval-print loop (REPL) will
start, where the programmer can try to fix any errors with the
program. In the REPL, the programmer can modify functions, print the
call stack, print the program structure, and many other useful
things.

```
package main

func main () () {
	var foo i32 = 10
	var bar f32 = 5.5
	
	i32.div(1, 0)
}
```

In this example, we start adding some variables and then we try to
divide two numbers: 1 by 0. If we run this program, this will be
printed:

```
Call's State:
foo:		10
bar:		5.5

i32.div() Arguments:
0: 1
1: 0

7: i32.div: Division by 0
CX REPL
More information about CX is available at http://cx.skycoin.net/

* 
```

The "Call's State" section tells us what are the values of the
variables in the current call. In this case, we have two variables:
foo and bar, with 10 and 5.5 as their values, respectively. Next, CX
tells us the arguments that were sent to the expression raising the
error, which are 1 (argument #0) and 0 (argument #1). A description of
the error is also provided: CX tells us that a "Division by 0" was
raised, and it was caused by *i32.div* at line #7. Lastly, we can see
that the CX REPL starts, where the programmer can enter some commands
to try to fix the error.

## Halt

Sometimes we want to *halt* a program's execution: maybe because we
want to check what's the current call's state, maybe because we want
to see if the program's flow is reaching somewhere, maybe simply
because it's fun. Whatever the reason might be, CX provides us with a
very helpful function, which is curiously named "halt."

```
for true {
	str.print("Enter a number greater than 0: ")
	value := i32.read()

	if i32.gt(value, 0) {
		str.print("Good, good.")
	} else {
		halt("The number was not greater than 0.")
	}
}
```

In this example, CX enters an infinite loop using *for true*. The user
is asked to enter a number greater than 0, and this value is *read*
using the *i32.read* function. If this value is greater than 0, the
program congratulates the user, and if it is not greater than 0, the
program halts. As you can see, halt receives a string as it's
argument, which serves as an error message for the user.

## Unit Testing

It is a good idea to create *test* files, which contain code that
makes sure that your defined functions are behaving correctly. For
example, let's assume that we have a function that converts an integer
from 0 to 3 to it's corresponding "word name" (for example, 3 =>
"three").

```
func toWord (num i32) (name str) {
	if i32.eq(0, num) {
		name = "zero"
		return
	}
	if i32.eq(1, num) {
		name = "one"
		return
	}
	if i32.eq(2, num) {
		name = "two"
		return
	}
	if i32.eq(3, num) {
		name = "three"
		return
	}
	name = "error"
	return
}
```

We could obviously create some function calls and make sure that the
results are as expected:

```
str.print(toWord(0))
str.print(toWord(1))
str.print(toWord(2))
str.print(toWord(3))
str.print(toWord(4))
```

This solution can work for some time, but it's going to get too
complicated when your program reaches dozens of functions. A better
solution is to use CX's *test* package:

```
test.start()
test.str(toWord(0), "zero", "0 failed")
test.str(toWord(1), "one", "0 failed")
test.str(toWord(2), "two", "0 failed")
test.str(toWord(3), "three", "0 failed")
test.str(toWord(4), "error", "0 failed")
test.stop()
```

*test.start* tells CX that unit testing will start, and that
 errors should only be printed to the user instead of halting the
 program. *test.str* receives three arguments: the first two can be
 anything that returns a string, and the third one is an error message
 that will be displayed to the user in case the test fails. In order
 for the test to succeed, the evaluations of the two first arguments
 must be the same.

Now let's imagine that we have a function that needs to return an
error with a particular set of arguments. How can we test for such a
case? The solution is to use *test.error*.

```
i32.div(10, 0)
test.error("i32.div did not raise a division by 0 error")
```

Another feature that *test.start* provides is that CX becomes aware of
the raised errors. In the example above, i32.div(10, 0) *must* raise
an error (if it doesn't, that's an error). After evaluating the bugged
expression, we call *test.error*. *test.error* will raise an error if
an error was not raised by the preceding expression, and the string
provided as its first argument will be shown to the user.

The *test* package has test functions for each of CX's basic types. If
you want to see a complete example, you can see CX's unit tests
[here](https://github.com/skycoin/cx/blob/master/tests/test.cx) (that
is, you're seeing how we use CX to test that CX is working correctly).

# Affordances

If we create a CX function, what can we do with it? We can call it, we
can add more expressions to it, we can add more input and output
parameters, we can remove them, we can change its name, we can remove
the function entirely... These are called affordances in CX, and they
help us achieve meta-programming: programs that can get themselves
modified.

CX applies the affordance paradigm by using its affordance system and
inference engine. The affordance system can determine everything that
can be done to an object, and everything that that object can do to
its surroundings. The inference system filters these affordances
according to certain criteria.

As this is a complex subject in CX, let's go step by step. First, we
need to tell CX somehow what element is our target, i.e., what element
we want to get affordances of. To do this, we need to create a
*target*:

```
target := ->{
  pkg(main) fn(double) exp(multiplication)
}
```

Everything contained in *->{...}* is practically a different
mini-language, but it tries to resemble what we have seen in CX until
now, as much as possible. *pkg* is used to tell CX what package we
want to target, *fn* to target a function, and *exp* to target an
expression. In this case, *main*, *double* and *multiplication* are
not CX variables; they are simply identifiers for the inference
engine.

Now, we need to create something similar to a *knowledge base*,
containing *facts* or *objects*, as are called in CX. Objects simply
tell CX something that is true or, more correctly, something that
exists in the current environment.

```
objects := ->{
    cloudy $0.7,
    hot $0.2
}
```

As you can see, we keep using the *->{...}* syntax. We are stating
that the objects *cloudy* and *hot* exist. Notice that objects are
separated by commas, and that they have numbers next to them, preceded
by a dollar sign ($). These numbers are called *weights*, and they
help us assign a grade of truthiness to the object. For instance, we
are perceiving the weather to be 0.7 cloudy, or *very* cloudy perhaps,
while we are perceiving it to be only 0.2 hot, or *not so* hot.

Lastly, we need to define a set of *rules* that describe how the
stated objects are going to filter the affordances determined by the
affordance system.

```
if cloudy $0.8 {
  allow(*.lightSensitive == true)
  obj(drones $1.0)
}
if and(cloudy $0.5, hot $0.1) {
  allow(*.numberWheels > 2)
  reject(*.solarPowered == true)
  obj(rovers $1.0)
}
if true {
  allow(*.class == "bipedal")
}  
if or(drones $1.0, rovers $1.0) {
  reject(*.class == "bipedal")
}
```

Let's imagine that we want to program some logic that determines what
kind of robots can or should be deployed to a particular
environment. We need to keep in mind that
at the beginning, the affordance system should have thrown *every*
robot that exists in the system and that is accessible to the targeted
CX expression in this case.

The first rule is telling the inference system to allow any robot that
is light sensitive, and then we dynamically add another object to the
object set: *drones $1.0* (actually, this is an actual set, as in the
mathematical sense, i.e., objects can't be repeated). Rules can add
new objects to affect the execution of the following rules. Adding the
*drones* object can be interpreted as saying "we can/will deploy
drones."

The second rule tells the inference system to allow any robot that has
more than 2 wheels, but to reject any robot that is solar powered
(because it's too cloudy). The robots meeting these criteria are
rovers, so we add an object *rovers*.

The third rule is interesting because, regardless of what objects we
have and regardless of their weights, bipedal robots are going to be
deployed... *unless* the fourth rule is true. The fourth rule tells
the inference engine to reject bipedal robots if we are deploying
either drones or rovers already.

Now we have all the requirements to create a query: a target, objects
and rules. But first, we need to assign a label to the expression that
we want to query.

```
multiplication:
	i32.mul(5, 0)
```

The expression above is a simple multiplication. Notice how we are
multiplying 5 by 0. CX's affordance system, when querying expressions,
is always going find out what arguments we can send as the last output
that we have given. In this case, we want to replace the 0. To label
the expression, we use the same kind of labels that we use for *go-to*
statements.

```
affs := aff.query(target, objs, rules)
```

The code above will query the affordance system and perform the
filtering according to the provided objects and rules. Notice that
*aff.query* returns something and is stored in *affs*. *aff.query*
returns all the affordances for the queried expression. We can print
the obtained affordances by using *aff.print*:

```
aff.print(affs)
```

And should return something similar to:

```
(0)	Operator: AddArgument	Name: foo1
(1)	Operator: AddArgument	Name: foo2
(2)	Operator: AddArgument	Name: foo3
(3)	Operator: AddArgument	Name: foo4
```

The numbers enclosed in parentheses are the affordance indexes, and
they are used to know what affordance we want to execute or apply in
particular. For example:

```
aff.execute(target, affs, 0)
```

*aff.execute* executes the nth affordance from a list of affordances
 to a target. In this case, we want to apply the 0th affordance from
 *affs* to *target* (which is the multiplication
 expression). Following the robots example, we actually want to send
 all the robots that meet the criteria defined by the rules. We can
 loop over all the affordances by using *aff.len*:

```
for c := 0; i32.lt(c, aff.len(affs)); c = i32.add(c, 1) {
	aff.execute(target, affs, c)
deployRobot:
	deployRobot(new Robot{class: "default"})
}
```

You should always send a default value, in case the affordance system
returns 0 results. Also, CX will complain that *deployRobot* requires
one argument and won't compile or interpret the program. If you don't
want a default value, just send a dummy value like 0, an empty
struct instance, etc., and then have an *if* statement check if the
resulting affordances array is empty before continuing with *aff.execute*.

## Limiting the Affordance System's Search Space

As you may have noticed, the affordance system takes into
consideration *all* your program as its search space. If you're
looking for *i32*s and you have some *i32* fields in some global
struct instances in an imported package, CX will consider it. This
behaviour might be useful for some type of applications: for example,
creating an IDE that lists you all the variables that you can send to an
expression as an argument. However, for many other applications, you
might want to limit this search space.

Let's imagine that you want to allow all the variables with values that are
greater than 2, you could write this rule:

```
rules := ->{
    if true {
        allow(x > 2)
    }
}
```

But the inference process is taking too long, even if you only wanted
to search in a pair of arrays that were sent as arguments to the
currently called function:

```
func doSomething (array1 []i32, array2 []i32) () {
 /* ... */
}
```
How can we limit the search space to *locals* and *arrays*? With the
`search` function:

```
rules := ->{
    if true {
        search(locals)
        search(arrays)
        allow(x > 2)
    }
}
```

You can create combinations for limiting the search space using the
following keywords: *nonArrays*, *arrays*, *structs*, *locals*,
*globals*, *allScopes*, and *allTypes*. As you can see, these keywords
limit the *type* (arrays or non-arrays) and *scope* (locals and
globals). Additionally, you can tell CX if you want to also search in
struct instances' fields.

Not limiting the search space is equivalent to the following:

```
rules := ->{
    if true {
        search(structs)
        search(allScopes)
        search(allTypes)
    }
}
```

# Experimental Features

CX continues growing every day. As a consequence, some features are
still in its infancy stages. This also means that they are very prone
to change and be improved in the future. Particularly, in this
tutorial we'll mention CX's evolutionary algorithm, CX's capability to
serialize itself, and GLFW and OpenGL API.

## Evolutionary Algorithm

After creating the first prototype of CX's affordance system, one of
the first ideas that came to mind was to use it to create an
evolutionary algorithm; particularly, a genetic programming (GP)
algorithm. Programming languages that can manipulate themselves or
that can programmatically create programs are usual targets to
experiment with new GP ideas, such as lisps.

In the case of CX, the affordance system is ideal for creating such
algorithms. CX can start with a null object that represents a CX
program, and then query that object to know what we can do with it. As
it is empty, the first and only action would be to add a *main*
package. The second option only option would be to create a *main*
function. From this point, the possibilities start to grow
exponentially: add new definitions, new functions, new input and
output parameters, add expressions to these functions, etc. Actually,
one of the first experiments using the affordance system was to create
a random program.

With some restrictions, it was easy to create a basic GP algorithm for
CX: we need to target a single function which will act as a solution
to a problem, this function needs to have only one input and one
output (although this will probably change later on), we'll be adding
only a limited number of expressions to the function, the operators
for these expressions are limited to a certain set and, lastly, the
arguments to these expressions are limited to local variables only. We
could provide a deeper explanation on how it is constructed, but
that's out of the scope of this tutorial, but you can [read the
source code](https://github.com/skycoin/cx/blob/master/src/base/evolution.go).

Let's suppose that we have an function that we don't know how it
works, but we know what are its outputs when sent certain inputs (for
example, a stock market time-series).

```
[-10 -9 -8 -7 -6 -5 -4 -3 -2 -1 0 1 2 3 4 5 6 7 8 9 10]
```

We decide to test the function with the input values above, and the
unknown function responded with the following outputs:

```
[110 90 72 56 42 30 20 12 6 2 0 0 2 6 12 20 30 42 56 72 90]
```

How can we know what function responds that way? Well, one option is
to just think a little but and after some minutes you should find out
that a possible function is *n\*n+n*, but let's pretend that this is a
harder problem. How can we find a solution to this problem? One way is
to use a [curve-fitting](https://en.wikipedia.org/wiki/Curve_fitting)
algorithm, such as neural networks or, you guessed it, genetic
programming.

First, let's program the *real function*:

```
func realFn (n f64) (out f64) {
	out = f64.add(f64.mul(n, n), n)
}
```

And now, let's create the function that *simulates* the real function:

```
func simFn (n f64) (out f64) {}
```

Woa! Wait, it's empty! Yeah, we're going to ask CX to fill it for us
using its evolutionary algorithm. If we had a rough idea of how the
real function is composed, we could help CX by writing some
expressions that approximate the solution, like:

```
func simFn (n f64) (out f64) {
  someHelp := f64.mul(n, n)
}
```

And the evolutionary algorithm could converge to a solution faster (or
maybe not. Maybe your guess is not as good as you think). But let's
make CX do all the work and leave *simFn* empty. Now, we only need to
call *evolve*:

```
evolve("simFn", "f64.add|f64.mul|f64.sub", inps, outs, 5, 100, f32.f64(0.1))
```

*evolve*'s first parameter is used to indicate CX what's the target
 function to be evolved, which is *simFn* in this case. The second
 parameter is a string representing a bag of functions to be used to
 find the solution. If you have no idea of what could be part of the
 solution, you could just write "." and CX will use every function it
 knows to create a solution. The third and fourth parameters are the
 inputs and outputs that represent the real function's behaviour. The
 sixth parameter represents how many expressions you want the solution
 to have.

The last two parameters are known as *stop criteria*, which
 are: for how many iterations do you want CX to run the evolutionary
 algorithm, and what is a "good enough" error to reach before
 stopping. But why do we need these parameters? Why not run the
 evolutionary algorithm until it finds *THE* solution? Well, for this
 example we have chosen a very easy problem to solve, but most
 problems in the real world are very hard to solve. *Bio-inspired*
 search algorithms, such as GP algorithms, are considered as
 *heuristics*, which is a fancy word to say that they won't
 necessarily reach the optimal solution. There are indeed algorithms
 that are guaranteed to find an optimal solution, but these algorithms
 can take a lot of time to find it (and by a lot of time, we mean
 weeks or even months, depending on the problem and how rich you are
 te get a suitable server, of course). Anyway, *stop criteria* are
 used to tell the evolutionary algorithm when it should stop
 searching.

Anyway, after either reaching 100 iterations or an error lower or
equal to 0.1, our evolutionary algorithm will stop, and now we can
test the solution:

```
str.print("Testing evolved solution")
for c := 0; i32.lt(c, []f64.len(inps)); c = i32.add(c, 1) {
	f64.print(simFn([]f64.read(inps, c)))
}
```

As the problem is very easy to solve, the code above should print the
same numbers that are present in the outputs array. If you are curious
on how the evolved function looks like, you can add a call to *halt*
before the program finishes, and type *:dProgram;* in the *REPL*. You
should see something like this:

```
1.- Function: simFn (n f64) (out f64)
			0.- Expression: var_1037 = f64.mul(n f64, n f64)
			1.- Expression: var_9874 = f64.sub(var_1037 , var_1037 )
			2.- Expression: var_9905 = f64.mul(var_9874 , n f64)
			3.- Expression: var_9936 = f64.add(var_1037 , var_9874 )
			4.- Expression: out = f64.add(var_9936 , n f64)
```

That function is an equivalent function to *n\*n+n*. Awesome, right?

## Serialization

CX programs can fully or partially serialize themselves. At the
moment, CX's serialization feature only has functions to completely
serialize the calling program, and to deserialize it back and print
its structure to the user

*serialize* is a parameter-less function, which returns a byte array
 that represents the current program.

```
sPrgrm := serialize()
```

And *deserialize* receives a byte array which should represent a
serialized program.

```
deserialize(sPrgrm)
```

After a call to *deserialize*, the terminal should print the program's
abstract syntax tree.

## OpenGL 1.2 API

CX, at the moment, provides at least the necessary functions to run a
basic game where you can move a character through a tile map. If you
are interested on what are the currently implemented functions, here's
a list for OpenGL 1.2:

```
func gl.Init () (error str) {}
func gl.CreateProgram () (progId i32) {}
func gl.LinkProgram (progId i32) () {}
func gl.Clear (mask i32) () {}
func gl.UseProgram (progId i32) () {}

func gl.Viewport (x i32, y i32, width i32, height i32) () {}
func gl.BindBuffer (target i32, buffer i32) () {}
func gl.BindVertexArray (target i32) () {}
func gl.EnableVertexAttribArray (index i32) () {}
func gl.VertexAttribPointer (index i32, size i32, xtype i32, normalized bool, stride i32) () {}
func gl.DrawArrays (mode i32, first i32, count i32) () {}
func gl.GenBuffers (n i32, buffers i32) () {}
func gl.BufferData (target i32, size i32, data []f32, usage i32) () {}
func gl.GenVertexArrays (n i32, arrays i32) () {}
func gl.CreateShader (xtype i32) (shader i32) {}

func gl.Strs (source str, freeFn str) () {}
func gl.Free (freeFn str) () {}
func gl.ShaderSource (shader i32, count i32, xstring str) () {}
func gl.CompileShader (shader i32) () {}
func gl.GetShaderiv (shader i32, pname i32, params i32) () {}
func gl.AttachShader (program i32, shader i32) () {}

func gl.MatrixMode (mode i32) () {}
func gl.LoadIdentity () () {}
func gl.Rotatef (angle f32, x f32, y f32, z f32) () {}
func gl.Translatef (x f32, y f32, z f32) () {}
func gl.Scalef (x f32, y f32, z f32) () {}
func gl.TexCoord2d (s f32, t f32) () {}
func gl.PushMatrix () () {}
func gl.PopMatrix () () {}
func gl.EnableClientState (array i32) () {}

func gl.BindTexture (target i32, texture i32) () {}
func gl.Color3f (red f32, green f32, blue f32) () {}
func gl.Color4f (red f32, green f32, blue f32, alpha f32) () {}
func gl.Begin (mode i32) () {}
func gl.End () () {}
func gl.Normal3f (nx f32, ny f32, nz f32) () {}
func gl.TexCoord2f (s f32, t f32) () {}
func gl.Vertex2f (nx f32, ny f32) () {}
func gl.Vertex3f (nx f32, ny f32, nz f32) () {}

func gl.Enable (cap i32) () {}
func gl.Disable (cap i32) () {}
func gl.ClearColor (red f32, green f32, blue f32, alpha f32) () {}
func gl.ClearDepth (depth f64) () {}
func gl.DepthFunc (xfunc i32) () {}
func gl.Lightfv (light i32, pname i32, params f32) () {}
func gl.Frustum (left f64, right f64, bottom f64, top f64, zNear f64, zFar f64) () {}

func gl.NewTexture (file str) (texture i32) {}
func gl.DepthMask (flag bool) () {}
func gl.TexEnvi (target i32, pname i32, param i32) () {}
func gl.BlendFunc (sfactor i32, dfactor i32) () {}
func gl.Hint (target i32, mode i32) () {}

func gl.Ortho (left f32, right f32, bottom f32, top f32, zNear f32, zFar f32) () {}
```

OpenGL constants:

```
var FALSE i32 = 0
var TRUE i32 = 1
var QUADS i32 = 7
var COLOR_BUFFER_BIT i32 = 16384
var DEPTH_BUFFER_BIT i32 = 256
var ARRAY_BUFFER i32 = 34962
var FLOAT i32 = 5126
var TRIANGLES i32 = 4
var POLYGON i32 = 9
var VERTEX_SHADER i32 = 35633
var FRAGMENT_SHADER i32 = 35632
var MODELVIEW i32 = 5888

var TEXTURE_2D i32 = 3553

var PROJECTION i32 = 5889
var TEXTURE i32 = 5890
var COLOR i32 = 6144

var MODELVIEW_MATRIX i32 = 2982
var VERTEX_ARRAY i32 = 32884

var STREAM_DRAW i32 = 35040
var STREAM_READ i32 = 35041
var STREAM_COPY i32 = 35042

var STATIC_DRAW i32 = 35044
var STATIC_READ i32 = 35045
var STATIC_COPY i32 = 35046

var DYNAMIC_DRAW i32 = 35048
var DYNAMIC_READ i32 = 35049
var DYNAMIC_COPY i32 = 35050

var BLEND i32 = 3042
var DEPTH_TEST i32 = 2929
var LIGHTING i32 = 2896
var LEQUAL i32 = 515
var LIGHT0 i32 = 16384
var AMBIENT i32 = 4608
var DIFFUSE i32 = 4609
var POSITION i32 = 4611

var TEXTURE_ENV i32 = 8960
var TEXTURE_ENV_MODE i32 = 8704
var MODULATE i32 = 8448
var DECAL i32 = 8449
var BLEND i32 = 3042
var REPLACE i32 = 7681

var SRC_ALPHA i32 = 770
var ONE_MINUS_SRC_ALPHA i32 = 771

var DITHER i32 = 3024
var POINT_SMOOTH i32 = 2832
var LINE_SMOOTH i32 = 2848
var POLYGON_SMOOTH i32 = 2881
var DONT_CARE i32 = 4352
var POLYGON_SMOOTH_HINT i32 = 3155
var MULTISAMPLE_ARB i32 = 32925
```

And here's a list for the GLFW functions:

```
func glfw.Init () () {}
func glfw.WindowHint (target i32, hint i32) () {}
func glfw.CreateWindow (window str, width i32, height i32, title str) () {}
func glfw.MakeContextCurrent (window str) () {}
func glfw.ShouldClose (window str) (flag bool) {}
func glfw.PollEvents () () {}
func glfw.SwapBuffers (window str) () {}
func glfw.GetFramebufferSize (window str) (width i32, height i32) {}
func glfw.SetKeyCallback (window str, fnName str) () {}
func glfw.SetMouseButtonCallback (window str, fnName str) () {}
func glfw.SetCursorPosCallback (window str, fnName str) () {}
func glfw.GetCursorPos (window str) (x f64, y f64) {}
func glfw.SetInputMode (window str, mode i32, value i32) () {}
func glfw.GetTime () (time f64) {}
```

GLFW constants:

```
var False i32 = 0
var True i32 = 1
var Press i32 = 1

var Cursor i32 = 208897
var StickyKeys i32 = 208898
var StickyMouseButtons i32 = 208899
var CursorNormal i32 = 212993
var CursorHidden i32 = 212994
var CursorDisabled i32 = 212995

var Resizable i32 = 131075
var ContextVersionMajor i32 = 139266
var ContextVersionMinor i32 = 139267
var OpenGLProfile i32 = 139272
var OpenGLCoreProfile i32 = 204801
var OpenGLForwardCompatible i32 = 139270

var MouseButtonLast i32 = 7
var MouseButtonLeft i32 = 0
var MouseButtonRight i32 = 1
var MouseButtonMiddle i32 = 2
```

If you're interested on having a look at the applications that have
been created using these APIs, check the [opengl examples folder](https://github.com/skycoin/cx/tree/master/examples/opengl).
