# CX version 0.6.2

2019-03-xx

Today the Skycoin development team releases the CX programming language
version 0.6.2. This is the second bugfix release of the CX 0.6 series and it
fixes all known language bugs for CX 0.6.

The focus of this release is to improve the quality of the language compiler
and interpreter, although a few language enhancements have also been added.

## New in This Release

### Language Improvements

There are a few language improvement over 0.6.1:

 * Support for exponents in floating point numbers (f32 and f64). There is
   also improved decimal point parsing.

   Here are a few examples: `1.0` `1.` `.1` `1e+2` `.1e-3`

 * Support for multi-dimensional slices.

   This means that the following construct will now work:

   ```
   func main () {
       var slc1 [][]i32
       var slc2 []i32
       slc2 = append(slc2, 1)
       slc2 = append(slc2, 2)
       slc2 = append(slc2, 3)

       slc1 = append(slc1, slc2)
   }
   ```

### Library Improvements

 There are no library improvements in CX 0.6.2.

### Fixed issues

  * \#70: Inline field and index "dereferences" to function calls' outputs.
  * \#72: Multi-dimensional slices don't work.
  * \#76: Using int literal 0 where 0.0 was needed gave no error.
  * \#134: Panic when declaring a variable of an unknown type.
  * \#156: Panic when using a function declared in another package without importing the package.
  * \#166: Panic when calling a function from another package where the package name alias a local variable name
  * \#271: CX floats cannot handle exponents

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
