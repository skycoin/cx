![cx logo](https://user-images.githubusercontent.com/26845312/32426758-2a4bbb00-c282-11e7-858e-a1eaf3ea92f3.png)

# CX Programming Language and Blockchain

[![Build Status](https://travis-ci.com/skycoin/cx.svg?branch=develop)](https://travis-ci.com/skycoin/cx) [![Build status](https://ci.appveyor.com/api/projects/status/y04pofhhfmpw8vef/branch/master?svg=true)](https://ci.appveyor.com/project/skycoin/cx/branch/master)

CX is a general purpose, interpreted and compiled programming
language, with a very strict type system and a syntax
similar to Golang's. CX provides a new programming paradigm based on
the concept of affordances.

## CX Chain

Generate new chain spec.
```bash
$ cxchain-cli newchain -hi=100 -ss=100 ./examples/blockchain/counter-bc.cx
```

Run publisher node with generated chain spec.
* Obtain chain secret key from generated `{coin}.chain_keys.json` file.
```bash
$ CXCHAIN_SK={secret_key} cxchain -enable-all-api-sets
```

Run client node with generated chain spec (use different data dir, and ports to publisher node).
```bash
$ cxchain -enable-all-api-sets -data-dir "$HOME/.cxchain/skycoin2" -port 6002 -web-interface-port 6422 
```

Run transaction against publisher node.
```bash
$ cxchain-cli run ./examples/blockchain/counter-txn.cx
```

Run transaction against client node and inject.
```bash
$ CXCHAIN_GENESIS_SK={genesis_secret_key} cxchain-cli run -n "http://127.0.0.1:6422" -i ./examples/blockchain/counter-txn.cx
```