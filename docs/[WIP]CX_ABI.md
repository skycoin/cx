# Abstract Binary Interface (ABI) and its Memory Layout Standard

## Questions to be answered
1. What is a CXTypeSignature?
2. Memory layout for function inputs, function outputs, function body
3. How returns should be rewritten?
4. Naming convention for input variable structs, output variable structs, body variable structs, temp variable struct, and global variable structs.
5. How stack is extended or truncated?
6. How functions are called and returned?

---
## What is a CXTypeSignature?
A CXTypeSignature is a representation of a type. It contains the name, offset, type, and meta. 

    type CXTypeSignature struct {
        NameStringID int
        Offset       int
        Type         CXTypeSignature_TYPE
        Meta         CXTypeSignature_META
    }

Examples:

The type signature of `TestVariables int32` would be:

    CXTypeSignature{
      Name: "TestVariable"
      Offset: OffsetInStack
      Type: int32
      Meta: unused
    }

The type signature of `TestArrayVariables []int32` would be:

    CXTypeSignature{
      Name: "TestArrayVariable"
      Offset: OffsetInStack
      Type: array of int32
      Meta: int32
    }

An array of CXTypeSignature can be written as a CXStruct.
For example, we have 

    []CXTypeSignature{
      TestVariables int32,
      TestArrayVariables []int32
      }
We can write this as

    CXStruct{
      Name: "TestStruct"
      Fields:  []CXTypeSignature{
                  TestVariables int32,
                  TestArrayVariables []int32
                },
    }

---
## Memory Layout for Function Inputs, Outputs, and Body

Suppose we have this function:
```
func TestFunction(x int, y int, z int) (out int) {}
```

In the example function above we have
- function name - "TestFunction"
- function inputs - (x int, y int, z int) which we can write as an array of TypeSignature or as a CXStruct like below: 
  ```
  FunctionInputs CXStruct { 
    x int
    y int
    z int
  }
  ```
- function outputs - (out int) which we can write as an array of TypeSignature or as a CXStruct like below:
  ```
  FunctionOutputs CXStruct { 
    out int
  }
  ```

So we will now have,
```
func FunctionName(FunctionInputs CXStruct) (FunctionOutputs CXStruct) {
    // body
}
```

When you call TestFunction() in assembly
1. x int is pushed down to stack
2. y int is pushed down to stack
3. z int is pushed down to stack
4. The Instruction pointer (IP) is set to location of code for function

And pushing x,y,z down to stack is same as pushing FunctionInputs CXStruct down to stack.

Now suppose we have this function:
```
func TestFunction(int x, int y, int z) (a int, b int) {
  var1 int
  var2 int

  var1=10
  var2=20
  ...
}
```

You could write this as

```
FnStructIn struct{
  x int
  y int
  z int
}

FnStructOut struct{
  a int
  b int
}

FnVariablesStruct struct{
  var1 int
  var2 int
}

func TestFunction(FnStructIn) (FnStructOut) {
  FnVariablesStruct.var1=10
  FnVariablesStruct.var2=20
}
```
Now we get to function body, on the stack we will have:

1. FnStructInput - all the input variables.
2. FnStructOutput - all the output variables that might be assigned.
3. FnVariablesStruct - a struct with all the variables defined in the function body.
4. FnTempVariablesStruct - all temp variables used inside the function body.
5. PackageGlobalsStruct - global variables accessible within the package.

In reality we want to put FnStructOutputs Above FnStructInputs because after function is done executing we throw out FnStructInputs (we dont need anymore), but the return still needs the output data from the function.

So the stack is layed out like this. 

First on stack memory is,

    FnStructOutputs - all the output variables for function that are returned or can be assigned
    - if we have return (2,3) we rewrite it as "Out.A = 2, Out.B  = 3, return". 
    - all returns are assigning the output variables, then calling a "return" that takes no arguements

Next is,

    FnStructInputs - This is all the variables we feed into function as parameters

Then we have the variables in the function body that are assigned

    FnVariablesStruct

Then we have the temp variables used inside the function

    FnTempVariablesStruct

Then we have the global variables that are accessible within the package

    PackageGlobalsStruct


---
## Naming Convention And How Returns Are Written
Now pretend we have a package, like CxMath and we have function TestFunction

```
Package CxMath

func TestFunction(int x, int y, int z) (a int, b int) {
  var1 int = 10
  var2 int = 20

  a = var1
  b = var2

  return a,b
}
```

So we now have package CxMath with a function TestFunction

We say the function name is 

    CxMath.TestFunction

And then 
- we have the input variable struct
- we have the output variable struct
- we have the variables defined in the function body

We write a struct name with an _ in front if its "internal" 

    CxMath._IN_TestFunction
    CxMath._OUT_TestFunction
    CxMath._BODY_TestFunction

And _ is OK, because no user defined function can start with _

    so return (a int = var1, b int = var2)

is a return expression but we rewrite into 

    CxMath._OUT_TestFunction.a = CxMath._BODY_TestFunction.var1;
    CxMath._OUT_TestFunction.b = CxMath._BODY_TestFunction.var2;
    return;

However these structs dont need to actually exist because 

* InStruct
* OutStruct
* BodyStruct

On the stack

    InStruct is 
    - x int
    - y int
    - z int

    OutStruct is
    - a int
    - b int

    BodyStruct is
    - var1 int
    - var2 int

The three structs on stack, one after another is same as

    FunctionStruct
    - x int //start of InStruct
    - y int
    - z int
    - a int //start of OutStruct
    - b int
    - var1 int //start of BodyStruct
    - var2 int

Its just appended

And you do not even need the "structs" its actually just
- list of variable signatures
- list of offsets (byte) from stack
- size of each type (Computed from type signature)

You see?

You can actually put the offset in the signature type or have a type signature struct with an offset parameter. Then a function body's variables or a struct is really just a []CxTypeSignature

---
## How Functions Are Called and Returned
When we call a function
- we take size of output variables, and size of input variables
- we expand the stack by that much, the size of those two together
- we zero all the bytes from current stack size to new stack size
- then we push or write the inputs to the function
- then we JMP or jump to the code for the function


Also on stack, we write down the JMP position or instruction pointer location before function was called so when function call is done, we can just put that value back and continue executing

So we 
- push the instruction pointer (current execution position in AST) to stack
- we allocate stack data/bytes we need 
- we zero those variables/bytes (or we make sure that people cannot use unitialized variables); but are same
- we write the input variables
- then we JMP to AST portion where function execution starts

if we had an assembly language, there is no such thing as a "function", that is what it would do but we actually call "ExecuteFunction" or something . Then when function is done
- we write outputs
- then we reduce stack size
- then we pop the execution pointer back

So when we call a function in a function body

- we create a struct in function body to hold the return variables
- then when function returns we might copy the values from the output struct in function body to the one in the calling function body

might be easiest?

- Or we can pass an integer in, with offset of the place the return struct is supposed to be written back to; which will normally be on the stack
