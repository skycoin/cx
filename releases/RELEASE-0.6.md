# CX version 0.6

2019-01-25

Today the Skycoin development team releases the CX programming language
version 0.6.

The focus of this release is two-fold:
 * Prepare CX for the integration into the blockchain
 * Fix many of the issues in the language in previous releases

## New in This Release

### Serialization

A CX program has a 'binary' representation that is intended to go on the
blockchain. This binary representation is often called a bytecode, and is a
common feature in many language.  The difference between CX and most other
languages is that in CX you have a bytecode representation of both the program
itself, and a running state. CX 0.6 has serialization of any program and
running state as well as deserialization of both.

So in the next step of the Sky development, you will be able to store a CX
program on a blockchain (a chain on Fiber) and store the state of it also on
the chain.  The next time the program is started, the state will be read from
the chain, and continued from where it left off, run for a while, providing a
service of some kind, and then the new state will be put back into a new
block. This is similar to what is called a smart contract on other
blockchains.

### Language Improvements

There are a number of new language features that were present in late release
of CX 0.5, but are now officially supported.  These are:

 * Functions as first-class objects
 * Callbacks (used in some libraries, like OpenGL)
 * Much improved handling of slices, esp. resize/copy/insert/remove functions, and boundary checks.
 * New control flow keywords: `break` and `continue`
 * New formatting for `printf`: `%v` for values of any type
 * Labels can now appear anywhere in a function.
 * Improved error reporting
 * ...and many improvements that were introduced during late stages of CX 0.5

### Library Improvements

 * Added GIF support to OpenGL

### Many fixed issues
 * \#32: Panic if return value is used in an expression
 * \#40: Slice keeps growing though it's cleared inside the loop
 * \#41: Scope not working in loops
 * \#50: No compilation error when using an invalid identifier
 * \#51: Silent name clash between packages
 * \#52: Some implicit casts were not being caught at compile time
 * \#53: CX was not catching an error involving invalid indirections
 * \#55: Single character declarations are now allowed
 * \#59: Fields of a struct returned by a function call can now be accessed
 * \#61: No compilation error when passing *i32 as an i32 arg and conversely
 * \#62: No compilation error when dereferencing an i32 var
 * \#63: Fixed a problem where inline initializations didn't work with dereferences
 * \#65: Return statements now work in CX, with and without return arguments
 * \#77: Fixed errors related to sending references of structs to functions and assigning references to struct literals
 * \#101: Using different types in shorthands now throws an error
 * \#104: Dubious error message when indexing an array with a substraction expression
 * \#105: Dubious error message when inline initializing a slice
 * \#108: Solved a bug that occured when two functions were named the same in different packages
 * \#131: Problem with struct literals in short variable declarations
 * \#132: Short declarations can now be assigned values coming from function calls
 * \#154: Sending pointers to slices to functions is now possible
 * \#167: Passing the address of a slice element is now possible
 * \#199: Trying to call an undefined function no longer throws a segfault
 * \#214: Fixed an error related to type deduction in references to struct fields
 * \#218: Type checking now works with receiving variables of unexpected types


### Documentation

 * Updated CX Book.

   The CX book for CX programmers is updated during the development phase of
   CX 0.6.

   However, most of the contents is still relevant to CX 0.5.x. It still does
   not have all the information about new features in CX 0.6 and will be
   released in a new CX 0.6 edition as soon as possible.

 * Improved Documentation About the CX Internals

   It has become much easier to contribute to the development of CX because
   there is now documentation on what to do. This documentation contains
   suggestions for many contributing roles, not only software development.  It
   also contains documentation on internal data structures and control flow,
   thereby shortening the learning period for new developers.

<!-- ## Roadmap -->

<!-- The next step after the release of version 0.6 is to integrate it into the -->
<!-- Skycoin blockchain.  The team will also update the CX book to contain all the -->
<!-- new features in CX 0.6 as well as improve the documentation of the libraries. -->

<!-- The next version of CX, version 0.7, will focus on adding Affordances, a -->
<!-- mechanism to secure and change programs in running state. It will also contain -->
<!-- more debug features and an IDE, that already exists as work in progress. -->

## About CX

CX is the programming language for smart contracts on the [Skycoin](https://www.skycoin.net/) blockchain. CX is a general purpose, interpreted and compiled programming language, with a very strict type system and a syntax similar to Golang's. CX provides a new programming paradigm based on the concept of affordances.


