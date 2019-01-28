# Changelog

### v0.6 (NOT YET RELEASED)
* Additions
  * Serialization and deserialization
  * Functions as first-class objects
  * Callbacks
  * Improve error reporting system
  * Add slice resize/copy/insert/remove functions
  * Add ReadF32Slice (gl.BufferData and gl.BufferSubData accept slice arguments)
  * Add slice helpers
  * Switch len and size in slice header to avoid unecessary alloc/copy.
  * Refactor cx/op_und.go::op_append
  * Refactor cx/utilities.go::WriteToSlice
  * Add runtime bound checks on slices
  * Print exit code string when a runtime error is thrown
  * CX process returns CX_RUNTIME_* exit code when a runtime error is thrown
  * Add strerror opcode returning a string matching the CX_* error code passed as argument
* Libraries
  * Added GIF support to OpenGL
* Fixed issues
  * #32: Panic if return value is used in an expression
  * #40: Slice keeps growing though it's cleared inside the loop
  * #41: Scope not working in loops
  * #50: No compilation error when using an invalid identifier
  * #51: Silent name clash between packages
  * #52: Some implicit casts were not being caught at compile time
  * #53: CX was not catching an error involving invalid indirections
  * #55: Single character declarations are now allowed
  * #59: Fields of a struct returned by a function call can now be accessed
  * #61: No compilation error when passing *i32 as an i32 arg and conversely
  * #62: No compilation error when dereferencing an i32 var
  * #63: Fixed a problem where inline initializations didn't work with dereferences
  * #65: Return statements now work in CX, with and without return arguments
  * #77: Fixed errors related to sending references of structs to functions and
        assigning references to struct literals
  * #101: Using different types in shorthands now throws an error
  * #104: Dubious error message when indexing an array with a substraction expression
  * #105: Dubious error message when inline initializing a slice
  * #108: Solved a bug that occured when two functions were named the same in
    different packages
  * #131: Problem with struct literals in short variable declarations
  * #132: Short declarations can now be assigned values coming from function calls
  * #154: Sending pointers to slices to functions is now possible
  * #167: Passing the address of a slice element is now possible
  * #199: Trying to call an undefined function no longer throws a segfault
  * #214: Fixed an error related to type deduction in references to struct fields
  * #218: Type checking now works with receiving variables of unexpected types
* Documentation
  * CONTRIBUTING.md: Information about how to contribute to CX
* IDE (WiP)
  * Added a simple guide
* CX GUI helper moved to its own repository at https://github.com/skycoin/cx-gui-helper

### v0.5.18 (CURRENT VERSION) [2018-11-27 Tue 21:33]
* **Affordances**:
  * Support for `affordances-of`: argument -> argument
  * Support for `affordances-of`: argument -> program
  * Support for `affordances-of`: argument -> struct
  * Support for `affordances-on`: argument -> argument
* `break`: Implemented
* `continue`: Implemented
* `printf`: Added %v format directive which tries to deduce an argument's type and prints its value
* `printf()`: Tell the type and value of extra arguments
* **Labels**: can now be anywhere in a function's block of statements
* Fixed issues:
  * Fix #39: Short-variable declarations now work properly with arrays, e.g. `bar := foo[0]`
  * Fix #49: Trying to access fields of non-struct instances now throws an appropriate error
  * Fix #84: `++` suffix operator now working
  * Fix #92: Conflict when calling multiple callbacks using GLFW
  * Fix #98: CX now throws an error when trying to redeclare a variable
  * Fix #112: `printf` now prints either a `MISSING` or `EXTRA` when there are
              fewer or more arguments than format directives, respectively
  * Fix #120: CX was throwing redeclaration errors in multiple return expressions

### v0.5.17
* Fixed issues:
  * Fix #90: `goto` now works properly on Windows
  * Fix #91: Methods with pointer receivers are working now
  * Fix #111: Trying to use global variables from other packages is no longer
              allowed without their owner package prefixed to them, i.e. `foo` was
	      allowed, now it must be written as `pkg.foo`
