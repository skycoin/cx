# Changelog

### v0.8.0
* Additions
  * Preperation for CX 1.0 Milestone
  * Massive internal compiler refactoring. Too many changes to list. See commit log

### v0.7.3
* Additions
  * Preperation for CX 1.0 Milestone

### v0.7.2
* Additions
  * CX package manager:
    - Added a flag for setting a CX workspace. This flag overrides the environment variable .
	- Import statements are now aware of the possibility of importing libraries that are defined in a CX workspace.
	- If the user has not set a CXPATH through an environment variable or --cxpath flag, CX will use by default `~/cx` and create all the necessary subdirectories.
	- If the directory supplied by the user to be used as CXPATH does not contain the following subdirectories: `src/`, `pkg/` and `bin/`, CX will create these subdirectories.
	- Increased opengl version from 2.1 to 3.2.
	- Increased glfw version from 3.2 to 3.3.
    - Added opcodes for reading/writing from/to binary file.
    - Added .hdr file loader.
	- Added opengl bindings : glClearBufferI, glClearBufferUI, glClearBufferF, glBlendFuncSeparate, glDrawBuffers.
	- Added glfw bindings : glfw.GetWindowContentScale, glfw.GetMonitorContentScale.
	- Added glfw constants : glfw.CocoaRetinaFramebuffer, glfw.ScaleToMonitor.
	- Added opengl constants : GL_NONE, GL_RED, GL_RGBA16F, GL_HALF_FLOAT, GL_UNSIGNED_INT_24_8, GL_R8.
	- Added opcodes for reading binary files : os.Seek, os.ReadUI16, os.ReadUI32, os.ReadF32, os.ReadUI16Slice, os.ReadUI32Slice, os.ReadF32Slice.
    - Added TextureGetPixel function.
    - Added basic Regexp library.
	- Added basic Cipher library.
* Changes
  * Updated affordances to play well with newer language features.
  * Removed cx-games as a module. It was just confusing as users would be
    redirected to an outdated version of the repo and the games are already
    mentioned in the README.
  * Improvements to `PrintProgram` (the function associated to the meta-command `:dp` in the REPL); the types printed to input/output parameters now print indexes, dereferences to fields and if they come from an external package.
  * Improvements to memory management when dealing with CX chains.
  * Multidimensional array literals are now working properly.
  * Multidimensional slice literals are now working properly.
  * Added i8/16/ui8/ui16/ui32/ui64 types.
  * Added str.lastindex built-in function.
Libraries
  * Add json bindings for reading json files: Open, Close, More, Next, Delim, Bool, Float64, Int64, Str.
  * Add json cx library to ease json parsing in cx.
  * Add HTTP library.
* Fixed issues
  * #306: Can't print double quotes.
  * #321: Can't do math inside 2nd square brackets (indexer) of an expression.
  * #419: Initializing array or slice with values equivalent to nil.
  * #481: Array literals: Problem with temporary variables in multi-dimensional arrays.
  * #482: Slice literals: Problem with temporary variables in multi-dimensional slices.
  * #131: Panic when using arithmetic operations.

* Documentation
  * ...
* Miscellaneous
  * ...

### v0.7.1 released 2019-07-26
* Additions
  * Added capability of storing heap objects in a CX chain's program state.
  * Added a function that prints information about the heap for debugging purposes (`debugHeap` in cx/utilities.go).
* Changes
  * Redesign of CX's garbage collector.
  * Changes to several functions that relied on allocating objects on the heap.
  * In previous versions of CX the data segment was living inside the heap segment. Now the data segment is properly separated from the heap segment.
  * Moved CX book sources to github.com/skycoin/cx-book. The releases and any code updates will be pushed to that repository.
    * The README file notifies the users about this change.
* Libraries
  * ...
* Fixed issues
  * #286: Compilation error when struct field is named 'input' or 'output'
  * #323: Installation issues on Windows after merging #320
* Documentation
  * ...
* Miscellaneous
  * ...

### v0.7.0 released 2019-06-15
* Additions
  * CX can now generate addresses to be used by CX chains.
  * New debug option --debug-lexer or -Dl to see which tokens that are
	returned by the lexer.
  * New debug option --debug-profile or -Dp to print timings information and dump cpu/mem
	profiles during compilation.
* Changes
  * Removed `cmd/cli`. The CX executable should now be used to generate CX chain addresses.
  * Running a CX chain with the --broadcast flag no longer runs the transaction code locally; it simply
	broadcasts the transaction code and it is run in the peer node to update the program state. If the user
	wants to test transaction code, the --transaction flag must be used.
  * Updated the style of the CX roadmap.
  * Changed max transaction size to 128 Mb. CX chains are storing their program state in two different unspent outputs in their skycoin fork. This means that if a CX program to be stored on a CX chain needs 64 Mb, then the CX chain will need at least a max transaction size of 128 Mb. This behavior needs to be corrected immediately in the next version of CX and the user needs a way to set these parameters via flags (they're hardcoded at the moment).
* Libraries
  * Add cx arg parsing library
* Fixed issues
  * #357: Error running cx in blockchain broadcast mode.
  * #360: Panic when package keyword is misspelled
  * #373: Error in address used to generate a CSRF token. Port was 6001 instead of 6421.
  * #338: Issues with slice of type T where sizeof T is different than 4.
  * #388: Array fail on cx --blockchain.
  * #389: CX chains errors with big programs.
* Documentation
  * New file `documentation/BLOCKCHAIN-OVERVIEW.md` which describes the processes and modules involved in CX chains.
  * The blockchain tutorial
	[documentation/BLOCKCHAIN.md](https://github.com/skycoin/cx/blob/develop/documentation/BLOCKCHAIN.md)
	will be used to reflect the state in the CX source code (`develop` branch)
  * The blockchain tutorial in the
	[wiki](https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial) will
	be used as a tutorial for the latest CX release. 
* Miscellaneous
  * ...

### v0.7beta released 2019-05-19
* Additions
  * First release of CX chains, i.e. CX programs stored on a Skycoin fiber blockchain
  * Added/forked the newcoin and skycoin-cli commands to the CX
    repositoryand adapted it to CX chains.
  * CX can now create a wallet by running `cx --create-wallet --wallet-seed $WALLET_SEED`
  * Added --wallet-id flag. This parameter replaces having to set the WALLET environment variable for CX chains.
* Changes
  * Transaction and block default sizes for CX chains changed from 32 Kb to 5 Mb.
* Libraries
  * Add math bindings for f32 and f64 types: isnan, rand, acos, asin.
  * Add glfw bindings: Fullscreen, GetWindowPos, GetWindowSize, SetFramebufferSizeCallback, SetWindowPosCallback, SetWindowSizeCallback.
* Fixed issues
  * #292: Compilation error when left hand side of an assignment expression is a struct field.
  * #309: Serialization is not taking into account non-default stack sizes.
  * #310: Splitting a serialized program into its blockchain and transaction parts.
  * #311: `CurrentFunction` and `CurrentStruct` are causing errors in programs with more than 1 package.
  * #312: Deserialization is not setting correctly the sizes for the CallStack, HeapStartsAt and StackSize fields of the CXProgram structure.
* Documentation
* IDE
  * Removed the current version of the IDE. We'll move to a textmate-based
    solution in the future.
* Miscellaneous
  * Vendored all external packages

### v0.6.2
* Additions
  * Support for exponents in floating point numbers (f32 and f64). Also
	improved decimal point parsing
	examples: 1.0 1. .1 1e+2 .1e-3
  * Support for multi-dimensional slices.
* Libraries
* Fixed issues
  * #58: No compilation error when indexing an array with a non integral var.
  * #70: Inline field and index "dereferences" to function calls' outputs.
  * #72: Multi-dimensional slices don't work.
  * #76: Using int literal 0 where 0.0 was needed gave no error.
  * #133: Panic when inserting a newline into a string literal.
  * #134: Panic when declaring a variable of an unknown type.
  * #155: Panic when trying to assign return value of a function returning void.
  * #156: Panic when using a function declared in another package without importing the package.
  * #166: Panic when calling a function from another package where the package name alias a local variable name
  * #271: CX floats cannot handle exponents
  * #284: Concatenation of str variables with + operator doesn't work.
  * #285: Short declaration doesn't compile with opcode return value.
  * #288: No compilation error when using float value in place of boolean expression.
  * #289: Panic when package contains duplicate function signature.
* Documentation
* IDE (WiP)
* Miscellaneous
  * CXFX, the game development toolkit, has been moved to its own
	repository at https://github.com/skycoin/FIXME

### v0.6.1
* Additions
  * Support for lexical scoping
  * `if/elseif` and `if/elseif/else` statements now work properly.
* Libraries
* Fixed issues
  * #54: No compilation error when defining a struct with duplicate fields.
  * #67: No compilation error when var is accessed outside of its declaring scope.
  * #69: glfw.GetCursorPos() throws error
  * #82: Empty code blocks (even if they contain commented-out lines) crash like this.
  * #99: Short variable declarations are not working with calls to methods or functions.
  * #102: String concatenation using the + operator doesn't work.
  * #135: No compilation error when using arithmetic operators on struct instances.
  * #153: Panic in when assigning an empty initializer list to a []i32 variable.
  * #169: No compilation error when assigning a i32 value to a []i32 variable.
  * #170: No compilation error when comparing value of different types.
  * #247: No compilation error when variables are inline initialized.
  * #244: Crash when using a constant expression in a slice literal expression.
	* The problem actually involved the incapability of using expressions as
	values in slice literals
* Documentation
* IDE (WiP)
* Miscellaneous

### v0.6
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

### v0.5.18 [2018-11-27 Tue 21:33]
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
