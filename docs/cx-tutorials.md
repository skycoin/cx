# Getting started with CX
## Hello World

We will first create a Hello World CX application. 

Create a `hello-world.cx` file with the following content. 

```
package main

func main () {
    str.print("Hello, World!")
}
```

Running `cx hello-world.cx` prints out `Hello, world!`. More basic CX patterns and examples can be found in the `examples` directory. 

## Basic Options

`cx --help` will print available commands and options. `cx --version` prints the currently installed version of CX. 

Frequently used options are:

* `--repl` loads program and runs it in REPL mode useful for debugging
* `--web` exposes endpoint for code evaluation at: http://127.0.0.1:5336/eval)


## REPL tutorial

Once CX has been successfully installed, running `cx` should print
this in your terminal:

```
CX 0.5.13
More information about CX is available at http://cx.skycoin.com/ and https://github.com/skycoin/cx/
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

# Learning CX

To get your hand dirty with the CX language, you can read the following reference:

## 1. CX Book

You can find the book source code and its releases in its [CX book Github
repository](https://github.com/Skycoin/cx-book).

## 2. Language Guide

You can find the language guide which is more up to date than the book in [language guide](LanguageGuide.md) section.