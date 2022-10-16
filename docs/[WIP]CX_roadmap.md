# 2022 CX Roadmap

## CX Refactoring
- [ ] Give complete details to all functions and their arguments. (in progress)
- [ ] Add unit tests to all functions and more acceptance tests to achieve 100% coverage.

## New CX AST Format
- [x] Add arrays for CXFunctions then reference by id
- [x] Add arrays for CXPackages then reference by id
- [x] Add arrays for CXlines then reference by id
- [x] Add arrays for CXArgs then reference by id
- [x] Add arrays for CXAtomicOps then reference by id
- [ ] Implement CXOPERATION_TYPE for CXExpression
    - [x] CX_ATOMIC_OPERATOR
    - [x] CX_LINE
    - [ ] CX_ARGUMENT
    - [ ] CX_STRUCT_DEF
    - [ ] CX_FUNCTION_DEF
    - [ ] CX_PACKAGE_DEF
    - [ ] CX_GLOBAL_DEF

## Struct Definition
- [ ] Complete new CXStruct definition specifications
- [ ] Complete new CXTypeSignature definition specifications.
- [ ] CXStruct Implementation
    - [x] CXStruct for function Inputs
    - [x] CXStruct for function Outputs
    - [ ] CXStruct for function variable layout
    - [ ] CXStruct for struct definitions
- [ ] CXTypeSignature Implementation
    - [x] atomic
    - [ ] PointerAtomic
    - [ ] ArrayAtomic
    - [ ] ArrayPointerAtomic
    - [ ] SliceAtomic
    - [ ] SlicePointerAtomic
    - [ ] Struct
    - [ ] PointerStruct
    - [ ] ArrayStruct
    - [ ] ArrayPointerStruct
    - [ ] SliceStruct
    - [ ] SlicePointerStruct
    - [ ] Complex
    - [ ] PointerComplex
    - [ ] ArrayComplex
    - [ ] ArrayPointerComplex
    - [ ] SliceComplex
    - [ ] SlicePointerComplex

Objectives:
- To eliminate CxArguments and replace with CXTypeSignature, we need int id codes for types (enum)
    - type (with a sizeof function for type)
    - offset for each variable
    Simple types are like int, float, int64, etc.

## Struct-def Library
For interfacing CX to Golang. Golang can call functions in CX program that are exposed and CX program can also call the golang game API. Example, in a game created with Golang, there will be a feature/support/API to program an object with CX. 

- [ ] CXStruct definition to Golang struct definition conversion and vice versa - this is needed for the function input arguments and returns. CX program takes a golang struct and read out the data from the golang struct to CX struct with a pointer.
- [ ] Unit tests of struct-def library.

Objectives:
- Allow CX to read/write from golang structs directly using "unsafe".

Notes:
- Support only for simple types, no pointers, only atomics.
- Library reads in a golang struct definition and gets the 
    - variable type
    - offset of variable from start of struct
    - size of variable

## Web Interface for CX game objects
To change the behavior while the game is running, change the behavior, change the image, sprite, sound effects, etc. Example is, there's a list of agents/objects, 35 ships, while the game is running, I can change the ship to attack or run away through the web interface while the game is running. To be added as a feature in CX-Playground/CX-IDE.

Objectives:
- To have the option to change the behavior of an object in a running game by modifying the CX Program in a friendly web interface.

## CX Memory - a requirement for CX to run on an embedded system
Every variable and everything that needs to be defined has to be in that linear array as embedded system only have one memory stack.

## Stack layout or the frame layout 
An API that gets "All variables in scope", the list of all the variables and their name, the type of variable, the length of variable (size) and offset
for a struct or for function, etc. 
Every variable used inside a function, every var that receives an assignment has to have a name,offset and size and including temp vars

## Function Scope
Push everything to functional scope, meaning that there is no local scopes within the functions. We may enforce local scoping rules so that, if you have a loop and there is a variable in the loop, when the loop finishes, that variable will still be defined, but you will not be able to access the variable outside the loop from the CX program.

## CX Program Encapsulation - Global Removal
A struct encapsulating the state of a single CX program, instead of a global variable containing the program.

- [x] All execution functions will receive and pass CXProgram pointer so they won't need to access global CXProgram variable.
- [ ] Restrict the compiler to only do one compilation at a time when it’s producing the AST and then force the encapsulation of the AST inside of a struct so that we can execute multiple concurrent CX programs at the same time.

Objectives:
- allow multiple CX programs to be loaded/run at same time by one golang process
- cleanup code base

Notes:
- Compiler will still use global when outputing AST from source file, but lock compilation so that only one program can be compiled at a time. (Maybe use a channel to pass in data with compiled AST returned)

## Abstract Binary Interface (ABI)
- [x] Define the ABI
- [ ] then we can take a CX function, then compiled it to x86 assembly or x64 with LLVM
- [ ] then we can "Execute" the function directly on the CPU
- [ ] it takes in byte array, cx program, etc then runs the assembly language instructions, that modify the byte array directly; no interpreter

Objectives:
- so functions in AST can be reduced completely to assembly language and executed natively eventually.

## CX Compiler Frontend
- [x] Stage 1: Package format - This is how we’re representing the file. Every file has a length, a name, and a hash. Every package is a list of files – file structs. And then, we have a list of the package structs. And then, we can serialize it and hash that to get the ID for the whole program.
- [x] Stage 2: Declaration Extraction
- [ ] Stage 3: Type Checks
- [ ] Stage 4: Function Body Extraction
- [ ] Integration of Stage 1-4 to the compiling process

## CX App Store
CX App Store will be a downloadable app for windows and mac. Once it is downloaded and installed, when ran, the app will start its own local web app server and automatically open up a web browser that redirects to that web app. It will show lists of CX programs/apps/games that the user can run. For example, the user clicks "CXPacman", it will then automatically open up an OpenGL window that runs the CXPacman game.
- [ ] Taskbar launcher 
- [ ] Windows and Mac installer
- [ ] Web App

Objectives:
- To easily download a CX program from CX App Store website and run the CX program.

## CX Apps
- [ ] CXS - CXShell

## CX Playground
- [x] Running server
- [ ] More Functionailities for CXPlayground

## Misc
- [ ] Memory Explorer
- [ ] App/library that shows all CX objects, layout values, location, etc
- [ ] Count of object types, total memory, used, unused, etc
- [ ] Stack trace library
- [ ] AST manipulation/inspection API
- [ ] CX Evolves task API
