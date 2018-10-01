![cx logo](https://user-images.githubusercontent.com/26845312/32426758-2a4bbb00-c282-11e7-858e-a1eaf3ea92f3.png)

Table of Contents
=================

   * [Table of Contents](#table-of-contents)
   * [CX Programming Language](#cx-programming-language)
      * [Strict Type System](#strict-type-system)
   * [CX Roadmap](#cx-roadmap)
   * [Video Games in CX](#video-games-in-cx)
   * [CX Playground](#cx-playground)
   * [Changelog](#changelog)
   * [Installation](#installation)
      * [Binary Releases](#binary-releases)
      * [Compiling from Source](#compiling-from-source)
      * [Installing Go](#installing-go)
   * [Additional Notes Before the Actual Installation](#additional-notes-before-the-actual-installation)
      * [Linux: Installing OpenGL and GLFW Dependencies](#linux-installing-opengl-and-glfw-dependencies)
         * [Debian-based Linux Distributions](#debian-based-linux-distributions)
      * [Windows: Installing GCC](#windows-installing-gcc)
      * [Installing CX](#installing-cx)
         * [Windows](#windows)
      * [Updating CX](#updating-cx)
   * [Running CX](#running-cx)
      * [CX REPL](#cx-repl)
         * [Running CX Programs](#running-cx-programs)
         * [Other Options](#other-options)
         * [Hello World](#hello-world)
   * [Syntax](#syntax)
      * [Comments](#comments)
      * [Declarations](#declarations)
         * [Allowed Names](#allowed-names)
         * [Strict Type System](#strict-type-system-1)
         * [Primitive Types](#primitive-types)
         * [Global variables](#global-variables)
         * [Local variables](#local-variables)
         * [Functions](#functions)
         * [Custom Types](#custom-types)
         * [Methods](#methods)
         * [Packages](#packages)
      * [Statements](#statements)
         * [If and if/else](#if-and-ifelse)
         * [For loop](#for-loop)
         * [Goto](#goto)
      * [Expressions](#expressions)
      * [Assignments and Initializations](#assignments-and-initializations)
   * [Runtime](#runtime)
      * [Packages](#packages-1)
      * [Data Structures](#data-structures)
         * [Variables](#variables)
         * [Primitive types](#primitive-types-1)
         * [Arrays](#arrays)
         * [Slices](#slices)
         * [Structures](#structures)
      * [Control Flow](#control-flow)
         * [Functions](#functions-1)
         * [Methods](#methods-1)
         * [If and if/else](#if-and-ifelse-1)
         * [For loop](#for-loop-1)
         * [Go-to](#go-to)
   * [Native Functions](#native-functions)
      * [Parse functions](#parse-functions)
      * [Unit testing](#unit-testing)
      * [OpenGL](#opengl)
      * [GLFW](#glfw)
      * [gltext](#gltext)

# CX Programming Language

CX is a general purpose, interpreted and compiled programming
language, with a very strict type system and a syntax
similar to Golang's. CX provides a new programming paradigm based on
the concept of affordances, where the user can ask the programming
language at runtime what can be done with a CX object (functions,
expressions, packages, etc.), and interactively or automatically choose
one of the affordances to be applied. This paradigm has the main objective
of providing an additional security layer for decentralized,
blockchain-based applications, but can also be used for general
purpose programming.

In the following sections, the reader can find a short tutorial on how
to use the main features of the language. In previous versions of this
README file, the tutorial was written in a *book-ish* style and was
targetted to a beginner audience. We're going to
be making a transition from that style to a more technical style,
without falling into a pure documentation style. The reason behind
this is that a CX book is now going to be published, which is going to
be targetted to a more beginner audience. Thus, this README now has
the purpose of quickly demonstrating the capabilities of the language,
how to install it, etc.

This tutorial/documentation is divided into four parts, which can
broadly be described as *introduction*, *syntax*, *runtime* and
*native functions*. The first section presents general information
about the language, such as how to install it, the development
roadmap, etc. The second and third sections are more *tutorial-ish*,
where the reader can find information about the core language,
i.e. how your program needs to look so it's considered a valid CX
program (*syntax*), and how your CX program is going to be executing
(*runtime*). The last section follows a more documentation style,
where each of the native functions of the language is presented, along
with an example.



Feel free to [create an issue](https://github.com/skycoin/cx/issues)
requesting a better explanation of a feature.

## Strict Type System

As mentioned in the description of the language, CX has a strict
type system. Most of the native functions in CX are associated to a
single type signature. For example, `str.print` seen in the "Hello,
World!" program only accepts strings as its input argument. There are
versions of `print` for each of the primitive types, such as
`i32.print`, `f32.print`, etc. The purpose of this is to

There are native functions in CX (the functions in the
core language) associated to

# CX Roadmap

![CX Roadmap](https://raw.githubusercontent.com/skycoin/cx/master/cx-roadmap.png)

# Video Games in CX

In order to test out the language, programmers from the Skycoin
community have already developed a number of video games, and more are
in-coming. Below is a list of the current video games that we are
aware of. If you are developing a video game using CX, please let us
know in our official telegram for video game development
https://t.me/skycoin_game_dev, or in the general CX group https://t.me/skycoin_cx.

- https://github.com/Lunier/Snake
- https://github.com/galah4d/casino-cx
- https://github.com/redcurse/SynthCat-Brick-Breaker
- https://github.com/atang152/crappyBall-cx
- https://github.com/galah4d/pacman-cx

# CX Playground

If you want to test some CX examples, you can do it in the CX
[Playground (http://cx.skycoin.net/)](http://cx.skycoin.net/).

# Changelog

Check out the latest additions and bug fixes in the [changelog](https://github.com/skycoin/cx/blob/master/CHANGELOG.md).






# Installation

## Binary Releases

This repository provides new binary releases of the language every
week. Check this link and download the appropriate binary release for
your platfrom: https://github.com/skycoin/cx/releases

More platforms will be added in the future.

CX has been successfully installed and tested in recent versions of
Linux (Ubuntu), MacOS X and Windows. Nevertheless, if you run into any
problems, please create an issue and we'll try to solve the problem as
soon as possible.

Once you have downloaded and de-compressed the binary release file,
you should place it somewhere in your operating system's $PATH
environment variable (or similar). The purpose of this is to have cx
globally accessible when using the terminal.

If you don't want to have it globally accessible, you can always try
out CX locally, inside the directory where you have the binary file.

## Compiling from Source

If a binary release is not currently available for your platfrom or if
you want to have a nightly build of CX, you'll have to compile from
source. If you're not familiarized with Go, Git, your OS's terminal or
your OS's package manager (to name a few), we *strongly* recommend you
to try out a binary release. If you find any bugs or problems with the
binary release, submit an issue here:
https://github.com/skycoin/cx/issues, and we'll fix it for the next
week's release.

## Installing Go

In order to compile CX from source, first make sure that you have Go
installed by running `go version`. It should output something similar to:

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
CX 0.5.13
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
0.- Package: main
	Functions
		0.- Function: main () ()
			0.- Expression: str.print("" str)
		1.- Function: *init () ()
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

### Running CX Programs

To run a CX program, you have to type, for example, `cx
the-program.cx`. Let's try to run some examples from the `examples`
directory in this repository. In a terminal, type this:

```
cd $GOPATH/src/github.com/skycoin/cx/
cx examples/hello-world.cx
```

This should print `Hello World!` in the terminal. Now try running `cx
examples/game.cx`.

### Other Options

If you write `cx --help` or `cx -h`, you should see a text describing
CX's usage, options and more.

Some interesting options are:

* `--base` which generates a CX program's assembly code (in Go)
* `--compile` which generates an executable file
* `--repl` which loads the program and makes CX run in REPL mode
(useful for debugging a program)
* `--web` which starts CX as a RESTful web service (you can send code
  to be evaluated to this endpoint: http://127.0.0.1:5336/eval)

### Hello World

Do you want to know how CX looks? This is how you print "Hello, World!"
in a terminal:

```
package main

func main () {
    str.print("Hello, World!")
}
```

Every CX program must have at least a *main* package, and a *main*
function. As mentioned before, CX has a strict type system,
where functions can only be associated with a single type
signature. As a consequence,
if we want to print a string, as in the example above, we have to call
*str*'s print function, where *str* is a package containing string
related functions.

However, there are some exceptions, mainly to functions where it makes
sense to have a generalized type signature. For example, the `len`
function accepts arrays of any type as its argument, instead of having
`[]i32.len()` or `[][]str.len()`. Another example is `sprintf`, which
is used to construct a string using a format string and a series of
arguments of any type.


# Syntax

In this section, we're going to have a look at how a CX program looks
like. Basically, the following sections are not going to discuss about
the logic behind the various CX constructs, i.e. how they
behave; we're only going to see how they look like.

## Comments

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

A declaration refers to a *named* element in a program's
structure, which are described using other constructs, such as
expressions and other statements. For example: a function can be
referred by its name and it's constructed by expressions and local
variable declarations.

### Allowed Names

Any name that satisfies the PCRE regular expression
`[_a-zA-Z][_a-zA-Z0-9]*` is allowed as an identifier for a declared
element. In other words, an identifier can start with an underscore
(*_*) or any lowercase or uppercase letter, and can be followed by 0
or more underscores or lowercase or uppercase letters, and any number
from 0 to 9.

### Strict Type System

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

There are seven primitive types in CX: *bool*, *str*, *byte*, *i32*,
*i64*, *f32*, and *f64*. Those represent Booleans (*true*
or *false*), character strings, bytes, 32-bit integers, 64-bit
integers, single precision and double precision floating-point
numbers, respectively.

### Global variables

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

Statements are different to declarations, as they don't create any
named elements in a program. They are used to control the flow of a
program.

### If and if/else

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

### For loop

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

# Runtime

The previous section presents the language features from a syntax
perspective. In this section we'll cover what's the logic behind these
features: how they interact with other elements in your program, and
what are the intrinsic capabilities of each of these features.

## Packages

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

Data structures are particular arrangements of bytes that the language
interprets and stores in special ways. The most basic data structures
represent basic data, such as numbers and character strings, but these
basic types can be used to construct more complex data types.

### Literals

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

### Variables

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
bytes). In the case of data types that point to variable-sized
structures, such as slices or character strings, these are initialized
to a nil pointer, which is represented by 4 zeroed bytes.

### Primitive types
### Arrays
### Slices
### Structures

## Control Flow
### Functions
<!-- talk about Expressions in here -->
### Methods
### If and if/else
### For loop
### Go-to

# Native Functions
## Parse functions
## Unit testing
## OpenGL
## GLFW
## gltext
