> This tutorial is valid for the CX `develop` branch
> When a release is done, it should be copied to the CX wiki
>   https://github.com/skycoin/cx/wiki/CX-Chains-Tutorial


Table of Contents
=================

   * [Introduction to CX chains](#introduction-to-cx-chains)
   * [Getting Started](#getting-started)
      * [Explain that we can initialize in the blockchain code](#explain-that-we-can-initialize-in-the-blockchain-code)
   * [Setup Process](#setup-process)
   * [Limitations, bugs and non-desirable behaviors](#limitations-bugs-and-non-desirable-behaviors)
   * [More examples](#more-examples)
      * [Hello, world!](#hello-world)
      * [Blockchain counter](#blockchain-counter)
      * [Saving factorial state to run later](#saving-factorial-state-to-run-later)
      * [Mario video game, save score](#mario-video-game-save-score)
      * [Smart contract](#smart-contract)
      * [Paying salaries](#paying-salaries)

# Introduction to CX chains

This text gives you an introduction to running CX program "on the blockchain". This means that you will write programs in the CX language that are stored on a Skycoin Fiber blockchain. It is assumed that the reader of this text has a basic understanding of CX itself, as well as a number of key components in the Skycoin universe such as Fiber blockchains.

# Getting Started

A blockchain is a distributed ledger secured by cryptographic methods. The Skycoin blockchain, and most other blockchains like e.g. Bitcoin, contains records that represent an amount of a digital asset. Such an asset is normally a crypto coin like Skycoin or Bitcoin.

However, a blockchain can also store many other types of data. A _CX chain_ is a blockchain that contains a CX program and everything that is necessary to run it in a series of updates. A CX chain works by separating a CX program into two parts: its _program code_ and the so called _transaction code_. The program state is stored on the blockchain, and it is updated every time a transaction is run. In this case, "update" means that the new state is stored in a new block in the blockchain.

Consider the code below:

```
package main

func main() {

}
```

One could think that this CX program represents an empty program state, but this is not the case. The program state consists of the following parts:
* a code segment
* a stack segment
* a data segment
* a heap segment

This is exactly like any CX program, but in this case the running state represented by the segments listed above will be stored on a blockchain instead of in RAM.

In the case of the example above, it would represent a code segment containing a single package (`main`) with a single function (`main`). As the `main` function is empty, nothing will run and the stack and heap segments will be empty. Furthermore, as we do not have any global variables, the data segment will be empty too.

Now, let's have a look at a slightly more complex example:


```
package main

var name str = "Richard"

func main() {

}
```

In this case, we have a code segment similar to the previous example and we have two pieces of data in the data segment: a global string variable and the string value "Richard". The variable points to the string value so that the value can be accessed by using the name of the variable in the code.

In order to store the CX program on the Skycoin blockchain, we need to serialize it. To serialize a program, all the packages, structs, functions, globals, etc. of the CX program are converted to a series of bytes, called a slice. The same happens to the stack, the data and the heap memory segments. These slices are stored in a transaction on the blockchain, which results in a "CX chain".  (You can look at `cx/serialize.go` in the CX source code if you are curious about the details.)

As a consequence, if a program with many packages, functions, etc. is used to initiate a CX chain, the program state stored on the blockchain will be big, even if the program is not creating any heap objects or using deep function calls that fill the stack. Just to make it clearer, the following rather simple code would result in a very big program state:

```
package main

var foo [1000][1000]i32

func main() {

}
```

It is important to bear this in mind, as there are hardcoded limits to how big the program state of a CX chain can be, as well as how big a transaction can be.

We now know how and where the program state of a CX chain is going to be stored, but we also need to know how to use it. The program state will never be executed or mutated by itself. In order to modify its state, we need to run a transaction against it. Consider the following pieces of code:

```
package bc

func PlusOne(indata i32) (result i32) {
	result = indata + 1
}

package main

func main() {

}
```

and

```
package main
import "bc"

func main() {
	var num i32
	num = bc.PlusOne(5)
	i32.print(num)
}
```

The first code snippet is our _blockchain code_ and the second snippet is our _transaction code_. As can be seen in the blockchain code, there are two packages: `bc` and `main`. Then, if we pay attention to the transaction code, we can see that its `main` package is importing `bc`. The blockchain code (which defines the CX chain's program state) can be seen as a code library that can store mutable state, and the transaction code can be seen as a program that is importing the blockchain code -- which is exactly what is happening.

"But that's not useful, I can just create another file with that code and import it" you might be thinking. However, remember that you can also store state and that state can be mutated. Consider this case:

```
package bc

var numTransactions i32

func LogTransaction() {
	numTransactions++
	i32.print(numTransactions)
}

package main

func main() {

}
```

and


```
package main
import "bc"

func main() {
	bc.LogTransaction()
}
```

In the example above, we can see that we have a global variable called `numTransactions`, which is increased by one if `LogTransaction` is called. This simple program can help us track how many times a transaction has been created and broadcasted.

As you may have noticed, the `main` function in the blockchain code is empty. The whole `main` package and any function that is in that package will be destroyed once the program state is created and stored on the blockchain. This does not mean that the `main` package is useless though, as it can be used to initialize the program state. For example, have a look at the following blockchain and transaction codes:

```
package bc

var matrix [10][10]i32

func initMatrix() {
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			matrix[x][y] = 1
		}
	}
}

func PrintMatrix() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			printf("%d", matrix[i][j])
		}
		printf("\n")
	}
}

func FlipBit(x i32, y i32) {
	if matrix[x][y] == 1 {
		matrix[x][y] = 0
	} else {
		matrix[x][y] = 1
	}
}

package main
import "bc"

func main() {
	bc.initMatrix()
}
```

and

```
package main
import "bc"

func main() {
	bc.FlipBit(3, 4)
	bc.PrintMatrix()
}
```

We use `bc.initMatrix` to initialize `matrix` to 1s, and the transaction code
uses `bc.FlipBit` to change one of the cells of the matrix to 0.

Now, this should be enough to give you a quick introduction to what you can
do. Now, let us go into the details of how to set up a development environment
and create your own CX blockchain.

# Setup Process

First, this introduction assumes that you have a working CX installation,
installed from the CX sources.  If you do not, then first, follow the
instructions in the sources on how to do that.

**NOTE: At the current early stage, you need to run this on Linux, since the
  CX blockchain has not been tested on Windows or MacOS X.**

After you have a working CX, perform the following steps:

* Update your CX repository:

```bash
$ cd $GOPATH/src/github.com/skycoin/cx
$ git pull origin develop
```

* Recompile with `make build`
* Run `cx --generate-address` and write down the result.
* Replace the fields `blockchain_pubkey_str` and `genesis_address_str` in `./fiber.toml` with the values generated by `cx --generate-address` (`public_key` and `address`, respectively)
* Run

```bash
cx --blockchain --heap-initial 100 --stack-size 100 \
    --secret-key $SECRET_KEY \
    --public-key $PUBLIC_KEY \
    ./examples/blockchain/counter-bc.cx
```

* Write down the blockchain's genesis signature that is shown at the bottom. CX will update fiber.toml's `genesis_signature_str` field with your new genesis signature.
* Start a publisher node (use the data obtained from running the previous commands; you can create bash variables that hold these values for convenience):

```bash
cx --publisher --genesis-address $GENESIS_ADDRESS \
   --genesis-signature $GENESIS_SIGNATURE \
   --secret-key $SECRET_KEY \
   --public-key $PUBLIC_KEY
```

* Run a peer node (**for now, it MUST use port 6001)**:

```bash
cx --peer --genesis-address $GENESIS_ADDRESS \
   --port 6001 \
   --genesis-signature $GENESIS_SIGNATURE \
   --public-key $PUBLIC_KEY
```

* Create a wallet (for now you need to use this seed. This will be changed ASAP):

```bash
cx --create-wallet --wallet-seed "museum nothing practice weird wheel dignity economy attend mask recipe minor dress"
```

* Restart your peer node by pressing `Ctrl-C` to stop the peer node process, and then run again this:

```bash
cx --peer --genesis-address $GENESIS_ADDRESS \
   --port 6001 \
   --genesis-signature $GENESIS_SIGNATURE \
   --public-key $PUBLIC_KEY
```

* You can test a transaction (without broadcasting the new program state) with:

```bash
cx --transaction examples/blockchain/counter-txn.cx
```

* And you can broadcast it with:

```bash
cx --secret-key $SECRET_KEY --broadcast examples/blockchain/counter-txn.cx
```

# Limitations, bugs and non-desirable behaviors

At the time of writing, the CXChain implementation suffers from a number of limitations and bugs. But remember that this is only a beta version and we will work to improve the state of CX and Skycoin Fiber.

* CXChain has only been tested on Linux, and the Debian distribution at that.
* You cannot send and receive SKY or any other cryptocurrency on Fiber (which is a good thing, as this CX chains are still in their experimental stage).
* You need to wait a few seconds before creating and broadcasting a new transaction.
* We don't have any security mechanism to prevent someone from calling or accessing certain parts of a CX chain's program state.  If you need security at this point, you should set up a firewall to your development workstation or work offline.
* The wallet's address that is sending transactions is hard coded at the moment.
* We need a way to set a seed for random numbers for the initial program state. This way we ensure determinism in subsequent transactions. Also, this seed should not be able to be changed by any transaction.

# More examples

In the CX source code you can find a number of increasingly complex examples that we have put together to increase your understanding:

## Hello, world!

```bash
cx --heap-initial 100 --stack-size 100 --blockchain examples/blockchain/hello-world-bc.cx
```

```bash
cx --heap-initial 100 --stack-size 100 --transaction examples/blockchain/hello-world-txn.cx
```

## Blockchain counter

```bash
cx --heap-initial 100 --stack-size 100 --blockchain examples/blockchain/counter-bc.cx
```

```bash
cx --heap-initial 100 --stack-size 100 --transaction examples/blockchain/counter-txn.cx
```

```bash
cx --heap-initial 100 --stack-size 100 --broadcast examples/blockchain/counter-txn.cx
```
