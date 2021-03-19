TODO.md

# Todo List for CX

This is a list of notes for compiler issues to fix

## Local Variable Declarations can Shadow Global Variables

If we have global variable declared in scope, a local variable can overshadow the global variable.

We need to throw warning when this occurs.

We need to verify the existing behavior and how it should be handled in CX.

Generally, locally declared variables should not shadow globally declared variables in same package. So it should throw warning in compiler.

Golang Suports. CX also supports shadowing of global variable with local variable.

```
package main
import "fmt"

var i int = 55

func main() {
	var i int = 44
	fmt.Println(i)
}
```

Is valid in golang. CX seems to have same behavior, but should verify with unit tests.


## CX Support for single quoted strings

We need to support single quoted strings.

We also need to support escape characters.

We need to unit test the string library and escape character behavior.

## Exposing GC functions to user

We need to expose the GC functions to user. For triggering GC operations during waits.

We need to be able to cap the GC execution time.

We need to be able to run GC operations in a second thread.

We need to implement memory remapping and movement of objects during run time.

## Read() function on Windows does not work

Read() reads from terminal?

Should be "ReadInput" or "ReadKeyboardInput" or "ReadUserInput" and needs a better name

## func defined with no arguments, called WITH arguments causes PANIC

```
func main() {
cat(1)
}

func cat() {
}
```

Should have compile time error. Not panic at run time.

## When initializing slices/arrays (with values), can't have final closed brace on it's own line

```
package main

func main() {
var s [3]str = [3]str{
"a",
"b",
"c"
}
}
```

c.cx:7: syntax error: unexpected SEMICOLON

"Unexpected semi-colon" is clearly wrong error message.

```
In golang, for last element of array it is expected to be , (comma) or } (RBRACE), otherwise there would be syntax error.
Although, C++ supports this kind of syntax.
```

## "break" exits all(?) outer "for loops" 

https://github.com/skycoin/cx/issues/17

```
func main() {
for h := 0; h < 5; h++ {
for i := 0; i < 5; i++ {
for j := 0; j < 5; j++ {
break
printf("j: %d \n", j)
	}
  		str.print("out of j loop")
	}

	str.print("out of i loop")
	}
	str.print("out of h loop")
}

```

We need to do break, and break(1), break(2) to denote how many levels break should break out of.

Also for continue.