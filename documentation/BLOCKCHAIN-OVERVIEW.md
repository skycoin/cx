CX Chains Overview
=================

This document has the purpose of dissecting CX chains into its different parts and to describe the processes that involve these parts. As the CX chains feature is still in a prototype/alpha stage, the descriptions contained in this file are subject to change. Changes can be expected after a new release of CX is created and after Github pull requests are merged into the `develop` branch.

Each of the following sections describe either a process or a module that is part of the composition of a CX chain. The document starts explaining what a [Blockchain Code](#blockchain-code) and a [Transaction Code](#transaction-code) are. After reviewing these concepts, it is then explained how these programs are necessary to create, modify and query a CX chain's [Program State](#program-state). The [Program State Structure](#program-state-structure) is then described, which helps the reader understand how a CX chain stores its state and what are the capabilities of a CX chain.

In order to create a CX chain, a number of parameters need to be set first in different files by using different methods. First, it is described how a [Coin Template](#coin-template) file is used to generate a program that can be used to run something similar to a cryptocurrency. Then, a process to generate a [Genesis Address and Genesis Private and Public Keys](#genesis-address-and-genesis-private-and-public-keys) is described, which outputs data that will be used to populate a [fiber.toml Configuration File](#fibertoml-configuration-file) that is used for [Initializing a CX Chain](#initializing-a-cx-chain). After this initialization process completes, it is explained how a [Wallet](#wallet) is created using a [Seed](#seed) and how we can start [Testing and Injecting Transactions](#testing-and-injecting-transactions).

For any encountered bug or feature request, we encourage the reader to create a Github issue with the inquiry.

Table of Contents
=================

   * [Blockchain Code](#blockchain-code)
      * [Program State](#program-state)
      * [Program State Structure](#program-state-structure)
         * [Code Segment](#code-segment)
         * [Stack Segment](#stack-segment)
         * [Data Segment](#data-segment)
         * [Heap Segment](#heap-segment)
   * [Transaction Code](#transaction-code)
   * [Coin Template](#coin-template)
   * [Genesis Address and Genesis Private and Public Keys](#genesis-address-and-genesis-private-and-public-keys)
   * [fiber.toml Configuration File](#fibertoml-configuration-file)
   * [Initializing a CX Chain](#initializing-a-cx-chain)
   * [Publisher and Peer Nodes](#publisher-and-peer-nodes)
   * [Wallet](#wallet)
      * [Seed](#seed)
   * [Testing and Injecting Transactions](#testing-and-injecting-transactions)

# Blockchain Code

## Program State

At the moment, only the data segment is updated.

## Program State Structure

The program state of a CX chain is represented by the serialization of a CX program.

### Code Segment

### Stack Segment

### Data Segment

### Heap Segment

# Transaction Code

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

# Coin Template



# Genesis Address and Genesis Private and Public Keys

In order to initialize a new CX chain, secret and public keys need to be generated to create the genesis transaction. Generating these keys is achieved by running the following command:

```
cx --generate-address
```

The output of this command will be similar to the one below:

```json
{
    "meta": {
        "coin": "skycoin",
        "cryptoType": "",
        "encrypted": "false",
        "filename": "2019_05_23_a737.wlt",
        "label": "cxcoin",
        "lastSeed": "64ceea88ba937fecab483ab6d2d9f51d4a02548cba71dbc494bab9550c0e6346",
        "secrets": "",
        "seed": "2a998dcce5470b87207a790db446219c046972b1f5bb618b0a5e851c972cc3e8",
        "tm": "1558675717",
        "type": "deterministic",
        "version": "0.2"
    },
    "entries": [
        {
            "address": "U84KDcpRbEK8ReHs7Z85MZd3KiFCCjUYPY",
            "public_key": "027ab554fef1fb125c5ec5317b830126cba5ba554f56ce08afb44eef8ead9cdfc1",
            "secret_key": "e2529cf862bd5a01c044966897e3ab4173e3df39cf2034f4c1c749e1ef0c3672"
        }
    ]
}
```

The bits of interest from this output are the values of the JSON keys `address`, `public_key` and `secret_key`.

These values are used for editing the file `fiber.toml`, with the exception of `secret_key`. At the moment, the modification of this file needs to be done manually, but this process should be performed automatically in later versions of CX. The value of the secret key must be kept secret, as the name implies, as this key could be used to sign transactions by anyone who posseses it.

# `fiber.toml` Configuration File

`fiber.toml` is used to set parameters that are used during the initialization and operation of a CX chain. The file already contains some values that can be considered as default, such as the `genesis_timestamp` or `max_block_size`, but other fields need to be set up with different values for every CX chain. Particularly, the fields:

- `blockchain_pubkey_str`
- `genesis_address_str`
- `genesis_signature_str`

need to be updated. The values of the first two fields are updated with the values obtained by following the instructions in the section [Genesis Transaction](#genesis-transaction), while the last one is automatically generated and added to `fiber.toml` by initializing a blockchain by running a command of the form `cx --blockchain --public-key $PUBLIC_KEY --private-key $PRIVATE_KEY blockchain-code.cx`.

Other fields that can be of interest in this file are:

- `create_block_max_transaction_size`
- `max_block_size`
- `unconfirmed_max_transaction_size`.

<!--- Ask for formal definitions of these parameters -->

Lastly, any field related to the configuration of a cryptocurrency should be ignored and left untouched for the functionality of a CX chain, specifically: 

- `genesis_coin_volume`
- `create_block_burn_factor`
- `unconfirmed_burn_factor`
- `max_coin_supply`
- `user_burn_factor`

<!--- Discuss remaining fields: blockchain_seckey_str, default_connections, genesis_timestamp, peer_list_url, port, web_interface_port, distribution_addresses, distribution_addresses_total, initial_unlocked_count -->

# Initializing a CX Chain

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

# Publisher and Peer Nodes

<!--- Explain that creating a new CX chain requires deleting the data directory -->

# Wallet

The transactions that are going to be run against the program state that is being stored on the blockchain need to be signed in order to meet the constraints imposed by the Skycoin blockchain platform.

Although, in theory, a secret key should be enough to sign a transaction, CX chains require at the moment to generate a wallet to be used to sign transactions. This wallet can be generated using the `cx` command, for example:

```
$ cx --create-wallet --wallet-seed "museum nothing practice weird wheel dignity economy attend mask recipe minor dress"
```

## Seed

Any transaction that occurs in a CX chain can be seen as a transaction between two accounts, which are represented by two addresses.  At the moment, these addresses involved in the transactions are hardcoded in CX. As a consequent, in order to run any transaction in a CX chain, a wallet created from the seed `"museum nothing practice weird wheel dignity economy attend mask recipe minor dress"` needs to be created.

The two addresses involved in any CX chain transaction are `TkyD4wD64UE6M5BkNQA17zaf7Xcg4AufwX` and `2PBcLADETphmqWV7sujRZdh3UcabssgKAEB`. If this was a transaction involving the transfer of SKY from one address to another, the former would be the address that is sending SKY to the latter.

# Testing and Injecting Transactions
