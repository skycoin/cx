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

## Behavior of Break and Continue

Also
- "continue" exits the entire current loop, AND an OUTER loop #15
- "break" exits all(?) outer "for loops" 

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

AND/OR introduced "labeled statements"
```
I think “labeled statements” would be ideal here, if CX doesn’t already support them. Go supports them, but I haven’t coded in Go. Here’s what they look like in Swift:
```

Labels are just integer identifiers for stack depth and is same as break(1), break(2), etc

## Maps for CX

CX needs to have maps implemented

## Using BREAK outside intended structures (loops)

CX should warn users that `break` is used incorrectly (for examples: if break is used in IF statement)

## Better error message for slice/array index missuse (when using commas inside square brackets)

```
package main

type Shape struct {
   Cells [][]bool
}

var Curr Shape

func main() {
   Curr.Cells[0] = append(Curr.Cells[0], true)
   Curr.Cells[0,0] = append(Curr.Cells[0,0], true)


   i32.print(len(Curr.Cells[0]))
   i32.print(len(Curr.Cells[0,0]))
}
```

https://github.com/skycoin/cx/issues/22

## Compilation error when using logical NOT operator on the return value of a function call

```
package main

func foo()(out bool) {
	out = true
}

func main()() {
	var b bool = !foo()
	test(b, false, "")
}
```

Expected behaviour: No compilation error

## No compilation error when redeclaring a global variable in the same package

```
package main

var i i32 = 5
var i i32 = 4

func main()() {
	panic(false, true, "must not compile")
}
```

Should give compile error for redeclaration of `i` variable

## No compilation error when appending a value of type B in a slice of type []A

```
package main

type fooA struct {
	f f32
}

type fooB struct {
	d f64
	i i32
}

func main()() {
	var sa []fooA
	var b fooB
	sa = append(sa, b)
	panic(true, false, "must not compile")
}
```

Should give compile error.

## Compilation error when using NEG operator on struct fields

```
package main

type too struct {
	x i32
}

func main()() {
	var t too
	t.x = -12

	t.x = -t.x // gives compile error
}
```

Expected behaviour: No compilation error

## Runtime error when calling str.i32 with a string not representing a valid i32 number.

To Reproduce:
```
package main
func main()() {
	i := str.i32("a") // gives runtime error
	k := str.i32("-2147483649") // gives runtime error
}
```

Parsing methods (str.*) should return an error on failure.
Should return error instead of crashing. 
If string is literal, then should be converted at compile time and not run time?

## Compilation error when using return value of len function in short hand expression

```
var s[]str
s = append(s, "33")
i := len(s)
```

Should not give compilation error.

## Panic when trying to print a negative literal as a string with printf

```
package main

func main()() {
	printf("%s\n", -1)
}
```

This should only give a runtime error and continue, but should not panic.
Ideally this should be done at compile time when the format string is a literal.

## No compilation error when appending to a pointer to a slice

```
var si []i32
var psi *[]i32 = &si
psi = append(psi, 4)
```

Should give compile error.
Error in golang: `first argument to append must be slice; have *[]int32`

## Compilation error when using multiple return values in short hand expression

```
package main
func foo()(i i32, f f32) {
	i = 5
	f = 1.0
}

func main(){	
	i, f := foo()
}
```

Expected behavior: No compilation error

##Compilation error when initializing multidimensiona arrays/slices with values

```
package main
func main() {
	var i [][]i32 = [][]i32{ {1, 2}, {3, 4} }
}
```
NOTE: Adding `[][]i32` before `{1, 2}` and `{3, 4}` compiles program without errors.
Expected behavior: No compilation error

## Compilation error when using inline initalization of a slice var as function argument

```
package main

func fooi(slice []i32) {}

func foos(slice []str) {}

func main()() {
	fooi([]i32{1, 2, 3})
	foos([]str {"foo", "slice", "str"})
}
```

Current error when running above code: `error: test1.cx:19 identifier '*tmp_4' does not exist`
Expected behavior: No compilation error

## No compilation error when using empty argument list after function call

```
package main
func foo()() {
	printf("foo\n")
}
func main()() {
	foo() () // should give compile error here
}
```

It appeares parenthesis after foo() are being ignored. Following doesn't give error either:
`foo() ()()()()()()()()`

## Compile error when using boolean variable in for loop

```
package main
func main()() {
	var b bool = true
	for b {
		printf("true\n")
	}
}
```
Expected behavior: No compilation error

## No compilation error when using void return value of a function in a for loop expression

```
package main
func foo()() {
}

func main() {
	for foo() { // panic, runtime error
		printf("true\n")
	}
}
```
Expected behavior: Should give compilation error

## Compilation error when using unary negative operator on a function call

```
package main

func foo()(out i32) {
    out = 4
}

func main()() {
    var a i32 = -foo()
}
```
Expected behavior: No compilation error

## No Compilation Error when function is called without parentheses

```package pkg

func FuncInPkg () {
str.print("FuncInPkg called with NO parentheses")
}

package main
import "pkg"

func main () {
pkg.FuncInPkg
}
```
Expected behavior: Some kind of error or at least a warning that you didn't use parentheses.

## Compilation error when using i32.add() in/as array's index (inside brackets)

```
package main

func main (){
        var nums [5]i32
        var temp i32
        for i := 0; i < 5; i++{
        	for j := 0; j < 5; j++{
        		if nums[j] > nums[i32.add(j, 1)]{ 
		                temp = nums[j]
			        nums[j] = nums[i32.add(j, 1)] // works here
			        nums[i32.add(j, 1)] = temp // here it produces compiler error
		        }
	        }
        }
}
```
Expected behavior: No compilation error

## Compiler Error Improvement for indexing non-indexable types 

```
package main

func main()() {
	var b bool = true
	var bb bool = b[0]
	panic(true, false, "must not compile")
}
```
Expected behavior: No panic.
Behavior should throw compiler error at compile time, 
"variable is not indexable type" or error stating that type cannot be indexed.
Should be compile time error, not run time error.

## Compilation error when declaring array of struct defined in other package

```
package pack

type Too struct {
    i i32
}

package main

type Moo struct {
    i i32
}

func main()() {
    var moo [3]Moo
    var too [3]pack.Too
}
```
Expected behavior: No compilation error

## CX_RUNTIME_ERROR when using a negative i32 variable to initialize a f32 variable

```
package main

func main()() {
    var f1 f32 = -1
}
```
Expected behavior: Compilation Error: `invalid implicit type conversion from i32 to f32`

## Compilation error: Wrong type inferred for variable declared with short declaration operator :=

```
package main

import "os"

func foo(i i32) (s str) {
s = "test"
}

func main()() {
    file0 := foo("test")
    file := os.Open("test")
    fileC := os.Close("test")
}
```
Expected behavior: No compilation error

## Compilation error when running declaring Pointer to Array

```
package main

type Point struct {
	x i32
	y i32
}

func testArrayPointerPnts (pnums [2]*Point) {}
func testArrayPointer (pnums [2]*i32) {}

func main () {
	var pnums [2]*i32
	n1 := 1
	n2 := 2

	pnums[0] = &n1
	pnums[1] = &n2

	testArrayPointer(pnums)

	i32.print(*pnums[0])

	var apPnts [2]*Point
	testArrayPointerPnts(apPnts)
}
```
Expected behavior: No compilation error

## Compilation error when redeclaring variables from different scopes/blocks

```
package main

type test_s struct {
    f f32
}
func test_()(out test_s) {
    out.f = 0.0
}

func main()() {
    {
        var s test_s = test_()
        var i i32 = 0
        var t test_s
    }
    {
        var s test_s = test_()
        var i i32 = 0
        var t test_s
    }
}
```
Expected behavior: No compilation error