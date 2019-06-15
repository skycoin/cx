# CX version 0.7beta

2019-06-15

Today the Skycoin development team releases the CX programming language
version 0.7.0. This is the first public release of CX 0.7, and it contains
improvements in several areas of CX: new features, libraries and bugfixes.

The focus of this release is the **integration of CX into the Skycoin Fiber
blockchains**, called CX chains.  Fiber is the way to create specialized
blockchains with CX code adapted to it, using Skycoin blockchain
technology. It is important to note that it is not the Skycoin blockchain
itself that is affected by CX chain, but a new blockchain using the same
technology.

Due to the importance of getting the blockchain programming language right,
there was a previous public beta of CX 0.7. This public release fixes all
known issues with the CX integration with the blockchain from the beta. The
most important one is that we increased the size of the biggest possible CX
chain program from 32 KB to 64 MB. Note that there are still a number of
issues with CX, but they are generic rather than related to the blockchain
integration.

The description on how to do the integration of CX and the blockchain can be
read here: https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial.

This release also contains a modified version of the Skycoin blockchain
repository with adaptations to run CX programs as transactions. These changes
will be moved back into the main Skycoin repository after extensive testing.

## New in This Release

### Blockchain integration

This is the first public release where blockchain integration is actively supported.
This section will see many additions in upcoming versions.

New in version 0.7.0 is that CX itself is used to generate addresses,
something that previously used the `cmd/cli` command.

### New Debug Option for CX Developers

Developers of CX itself can now debug CX by making the lexical analyzer output
the tokens that are being returned.

### Language Improvements

There are no language-specific improvements in CX 0.7.0.

### IDE

There is no longer an IDE in the CX repository.

### Library Improvements

There are no library improvements in CX 0.7.0

### Fixed issues

* \#357: Error running cx in blockchain broadcast mode.
* \#360: Panic when package keyword is misspelled
* \#373: Error in address used to generate a CSRF token. Port was 6001 instead of 6421.
* \#388: Array fail on cx --blockchain.
* \#389: CX chains errors with big programs.

### Documentation

* The documentation for the CX chains has been improved. It can be found on
  the CX wiki here: https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial.

## About CX

CX is the blockchain programming language for use on the
[Skycoin](https://www.skycoin.net/) blockchain and associated Fiber
blockchains. CX is a general purpose, interpreted and compiled programming
language, with a very strict type system and a syntax similar to Golang's. CX
provides a new programming paradigm based on the concept of affordances.
