# Changelog

### v0.5.19 (NOT YET RELEASED)
* Functions as first-class objects
* Callbacks
* Refactor parser actions
  
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
  * Fix #120: CX was throwing redeclaration errors in multiple return expressions
  * Fix #98: CX now throws an error when trying to redeclare a variable
  * Fix #92: Conflict when calling multiple callbacks using GLFW
  * Fix #49: Trying to access fields of non-struct instances now throws an appropriate error
  * Fix #39: Short-variable declarations now work properly with arrays, e.g. `bar := foo[0]`
  * Fix #84: `++` suffix operator now working
  * Fix #112: `printf` now prints either a `MISSING` or `EXTRA` when there are fewer or more arguments than format directives, respectively
### v0.5.17
  * Fixed issues:
	* Fix #111: Trying to use global variables from other packages is no longer allowed without their owner package prefixed to them, i.e. `foo` was allowed, now it must be written as `pkg.foo`
	* Fix #90: `goto` now works properly on Windows
	* Fix #91: Methods with pointer receivers are working now
