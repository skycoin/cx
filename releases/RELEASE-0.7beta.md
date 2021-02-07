# CX version 0.7beta

2019-05-19

Today the Skycoin development team releases the CX programming language
version 0.7beta. This is a public beta release of CX 0.7, which contains
improvements in all areas of CX: new features, language improvements,
libraries and bugfixes.

The focus of this release is the **integration of CX into the Skycoin Fiber
blockchains**, called CX chains.  Fiber is the way to create specialized
blockchains with CX code adapted to it, using Skycoin blockchain
technology. It is important to note that it is not the Skycoin blockchain
itself that is affected by CX chain, but a new blockchain using the same
technology.

The description on how to do the integration of CX and the blockchain can be
read here: https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial.

This release also contains a modified version of the Skycoin blockchain
repository with adaptions to run CX programs as transactions. These changes
will be moved back into the main Skycoin repository after extensive testing.

## New in This Release

### Blockchain integration

This is the first release where blockchain integration is actively supported.
This section will see many additions in upcoming versions.

### Language Improvements

There are no language-specific improvements in CX 0.7beta.

### IDE
* The embryo for the IDE was removed because it lacked maturity. We will move
  to a textmate-based solution in the future.

### Library Improvements

* Add new math bindings for f32 and f64 types: isnan, rand, acos, asin.
* Add new glfw bindings: Fullscreen, GetWindowPos, GetWindowSize, SetFramebufferSizeCallback, SetWindowPosCallback, SetWindowSizeCallback.

### Fixed issues

* \#292: Compilation error when left hand side of an assignment expression is a struct field.
* \#309: Serialization is not taking into account non-default stack sizes.
* \#310: Splitting a serialized program into its blockchain and transaction parts.
* \#311: `CurrentFunction` and `CurrentStruct` are causing errors in programs with more than 1 package.
* \#312: Deserialization is not setting correctly the sizes for the CallStack, HeapStartsAt and StackSize fields of the CXProgram structure.

### Documentation

* The documentation for the CX chains can be found on the CX wiki here:
  https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial.

## About CX

CX is the blockchain programming language on the
[Skycoin](https://www.skycoin.com/) blockchain. CX is a general purpose,
interpreted and compiled programming language, with a very strict type system
and a syntax similar to Golang's. CX provides a new programming paradigm based
on the concept of affordances.
