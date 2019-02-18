# CX version 0.6.1

2019-02-15

Today the Skycoin development team releases the CX programming language
version 0.6.1. This is the first bugfix release of the CX 0.6 series.

The focus of this release is to improve the quality of the language compiler
and interpreter.


## New in This Release

### Language Improvements

There is two major improvement over 0.6.0:

 * Support for lexical scoping.

   This means that a construct like this:

   ```
   if (...) {
      var x i32 = 0
      i32.print(x)
   }
   x = 10
   ```
   will fail because the scope of `x` ends at the right brace that ends the if
   clause. From CX 0.6.1 the CX language will enforce strict lexical scoping.

 * `if/elseif` and `if/elseif/else` constructs now work correctly.
 * `+` can now be used to concatenate strings:
 
 ```
 var conc str
 conc = "Hello " + "world!"
 ```
 
 * Expressions can now be used when initializing slice/array literals:
 
 ```
 var slc []i32
 slc = []i32{i32.add(1, 2), foo()}
 ```

### Library Improvements

CX 0.6.1 does not contain any library improvements.

 * Test functions (`assert`, `panic`, `test`) now check that their first two input arguments are of the same type. This way comparing values of different types is prohibited. In a future release a function that checks the equality of their arguments regardless of types will be introduced.

### Fixed issues

  * \#54: No compilation error when defining a struct with duplicate fields.
  * \#67: No compilation error when var is accessed outside of its declaring scope.
  * \#69: glfw.GetCursorPos() throws error
  * \#82: Empty code blocks (even if they contain commented-out lines) crash like this.
  * \#99: Short variable declarations are not working with calls to methods or functions.
  * \#102: String concatenation using the + operator doesn't work.
  * \#135: No compilation error when using arithmetic operators on struct instances.
  * \#153: Panic in when assigning an empty initializer list to a []i32 variable.
  * \#169: No compilation error when assigning a i32 value to a []i32 variable.
  * \#170: No compilation error when comparing value of different types.
  * \#247: No compilation error when variables are inline initialized.
  * \#244: Crash when using a constant expression in a slice literal expression.
	* The problem actually involved the incapability of using expressions as
	values in slice literals

### Documentation

 * Updated CX Book.

   The CX book is being updated, but this is not part of this release.  The
   book version 0.6 will have its own release at a later time.  If you want to
   look at the work in progress, you can find a snapshot of it in the `book/`
   subdirectory. 

## About CX

CX is the programming language for smart contracts on the
[Skycoin](https://www.skycoin.net/) blockchain. CX is a general purpose,
interpreted and compiled programming language, with a very strict type system
and a syntax similar to Golang's. CX provides a new programming paradigm based
on the concept of affordances.
