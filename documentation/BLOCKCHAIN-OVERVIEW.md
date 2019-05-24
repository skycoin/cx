Table of Contents
=================

   * [CX Chains Overview](#cx-chains-overview)
      * [Introduction](#introduction)
      * [Public and Secret Keys Generation](#public-and-secret-keys-generation)
      * [In order to initialize a new CX chain, secret and public keys need to be generated.](#in-order-to-initialize-a-new-cx-chain-secret-and-public-keys-need-to-be-generated)
      * [newcoin and cxcoin](#newcoin-and-cxcoin)
      * [Wallet](#wallet)
         * [Seed](#seed)
      * [fiber.toml](#fibertoml)
      * [Blockchain Code](#blockchain-code)
         * [Program State](#program-state)
      * [Transaction Code](#transaction-code)
         * [Transaction and Block Sizes](#transaction-and-block-sizes)
         * [Genesis Address](#genesis-address)
         * [Genesis Signature](#genesis-signature)
         * [Public and Secret Keys](#public-and-secret-keys)
      * [Program State Structure](#program-state-structure)
         * [Code Segment](#code-segment)
         * [Stack Segment](#stack-segment)
         * [Data Segment](#data-segment)
         * [Heap Segment](#heap-segment)
      * [Testing and Broadcasting a Transaction](#testing-and-broadcasting-a-transaction)
      * [Publisher and Peer Nodes](#publisher-and-peer-nodes)

# CX Chains Overview

## Introduction

This document has the purpose of dissecting CX chains into its different parts and to describe the processes that involve these parts. As the CX chains feature is still in a prototype/alpha stage, the descriptions contained in this file are subject to change. Changes can be expected after a new release of CX is created and after Github pull requests are merged into the `develop` branch. 

## Public and Secret Keys Generation

In order to initialize a new CX chain, secret and public keys need to be generated.
- 


## `newcoin` and `cxcoin`

Currently, CX has the `newcoin` command as a dependency (located in `cmd/newcoin`). In order to initialize a new CX chain, `newcoin` needs to create the `cxcoin` command (located in `cmd/cxcoin`) using the parameters defined in `./fiber.toml`. This process should be simplified by mimicking the process defined by `cmd/cxcoin/cxcoin.go` in `cxgo/main.go`. In other words, instead of having to call the process defined by `cxcoin`, we can embed the process in CX to eliminate the `newcoin` and `cxcoin` dependencies.

The workflow -- which occurs when running, for example, `cx --blockchain --secret-key $SECRET_KEY --public-key $PUBLIC_KEY examples/blockchain/counter-bc.cx` -- is as follows:

1. `newcoin` is installed by running `go install ./cmd/newcoin/...`
2. `newcoin` is run in order to create `cxcoin`
> It is worthy to note that the name `cxcoin` could be changed to something else by using the `--program-name` flag, but this behavior has not been tested yet.
3. `cxcoin` is installed by running `go install ./cmd/cxcoin/...`
4. `cxcoin` is run to initialize the CX chain. This process involves:
    1. Running `go run ./cmd/cxcoin/cxcoin.go --block-publisher=true --blockchain-secret-key=$SECRET_KEY`.

    > The data directory for the publisher node is stored in `$HOME/.cxcoin/`. Every time a new CX chain is initialized, its data directory is deleted first. The name of this directory can change, depending on the value of `--program-name`.
    2. As this is a new blockchain, the genesis block will be created and a genesis signature is generated. This genesis signature will be different for every blockchain, even if the blockchain private and public keys are the same.
    3. Using the genesis signature and the secret key, a CX chain creates the first transaction in the genesis block, which is a transaction without a transaction code and with an unspent output storing the initial program state, defined by running the blockchain code.
    4. `fiber.toml`'s field `genesis_signature_str` is updated automatically with the new genesis signature.
    5. The genesis signature is printed to standard output, so the user can take note of it.


## Wallet

The transactions that are going to be run against the program state that is being stored on the blockchain need to be signed in order to meet the constraints imposed by the Skycoin blockchain platform.

[//]: # "UPDATE POINT"
Although, in theory, a secret key should be enough to sign a transaction, CX chains require at the moment to generate a wallet to be used to sign transactions. This wallet can be generated using the `cx` command, for example:

```
$ cx --create-wallet --wallet-seed "museum nothing practice weird wheel dignity economy attend mask recipe minor dress"
```

### Seed

[//]: # "UPDATE POINT"
Any transaction that occurs in a CX chain can be seen as a transaction between two accounts, which are represented by two addresses.  At the moment, these addresses involved in the transactions are hardcoded in CX. As a consequent, in order to run any transaction in a CX chain, a wallet created from the seed `"museum nothing practice weird wheel dignity economy attend mask recipe minor dress"` needs to be created.

The two addresses involved in any CX chain transaction are `TkyD4wD64UE6M5BkNQA17zaf7Xcg4AufwX` and `2PBcLADETphmqWV7sujRZdh3UcabssgKAEB`. If this was a transaction involving the transfer of SKY from one address to another, the former would be the address that is sending SKY to the latter.

## `fiber.toml`

## Blockchain Code

### Program State

At the moment, only the data segment is updated.

## Transaction Code

The program state that is stored on the blockchain can either be mutated or queried by running a program that "imports" the program state. Actually, in order to have any sort of access to the program state, a program needs to import packages that are stored in the program state. As a minimalist example, consider the following code:

```go
package number

var Num i32

func main() {
  Num = 10
}
```

In order to modify the value of `Num`, the following transaction code can be used:

```
package main
import "number"

func main() {
  number.Num = 11
}
```

As can be seen, this resembles exactly what would happen in a CX program that is importing a package, either located in the same file or in an external file. The difference between a CX chain and the aforementioned situation is that in a CX chain, the program state will be preserved. For example, consider the following transaction code:

```
package main
import "number"

func main() {
  i32.print(number.Num)
}
```

In the case above, `i32.print(number.Num)` will print `11`, because the previous transaction code modified the value of `number.Num`.

### Transaction and Block Sizes



### Genesis Address

### Genesis Signature

### Public and Secret Keys

## Program State Structure

The program state of a CX chain is represented by the serialization of a CX program.

### Code Segment

### Stack Segment

### Data Segment

### Heap Segment

## Testing and Broadcasting a Transaction

## Publisher and Peer Nodes

