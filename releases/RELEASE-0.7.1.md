# CX version 0.7.1

2019-07-26

Today the Skycoin development team releases the CX programming language
version 0.7.1. This is the second public release of CX 0.7, and it contains
improvements mainly in the area of blockchain integration.

The focus of this release is the **ability to store heap objects on CX
Chains**.  This means that the programmer now has the ability to store
strings, arrays and slices on the blockchain, which was not possible before.

To achieve this, the CX team had to redesign and improve the builtin garbage
collector. This is a major internal improvement even though it contains no
user-visible changes.

CX chains are also known as Fiber chains.  Fiber is the way to create
specialized blockchains with CX code adapted to it, using Skycoin blockchain
technology. It is important to note that it is not the Skycoin blockchain
itself that is affected by CX chain, but a new blockchain using the same
technology.

Note that there are still a number of issues with CX, but they are generic
rather than related to the blockchain integration.

The description on how to do the integration of CX and the blockchain can be
read here: https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial.

This release also contains a modified version of the Skycoin blockchain
repository with adaptations to run CX programs as transactions. These changes
will be moved back into the main Skycoin repository after extensive testing.

## Backward Compatibility

CX 0.7.1 should be able to run programs on CX Chains created using CX 0.7.0.
This means that you will be able to upgrade CX and still run your existing
Skycoin Fiber applications.

## New in This Release

### Ability to store strings, arrays and slices on the blockchain

The ability to store all types of data on the blockchain is important for the
development of advanced dApps (distributed applications on the
blockchain). Some of these types store part of their contents on the heap
rather than in the data segment. Hence, CX 0.7.1 now stores both the contents of
the heap that are referenced by global variables on the blockchain.

### Language Improvements

There are no language-specific improvements in CX 0.7.1.

### Library Improvements

There are no library improvements in CX 0.7.1.

### Fixed issues

* \#286: Compilation error when struct field is named 'input' or 'output'
* \#323: Installation issues on Windows after merging \#320
* Several fixes to builtin functions that create heap objects which had hidden
  bugs that were not registered in the issue tracker.

### Documentation

* The documentation for the CX chains has been slightly improved. It can be
  found on the CX wiki here: https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial.

## About CX

CX is the blockchain programming language for use on the
[Skycoin](https://www.skycoin.net/) blockchain and associated Fiber
blockchains. CX is a general purpose, interpreted and compiled programming
language, with a very strict type system and a syntax similar to Golang's. CX
provides a new programming paradigm based on the concept of affordances.
