# 2022 CX Roadmap

## CX Refactoring
1. Implement new CX AST format. (in progress)
2. Give complete details to all functions and their arguments. (in progress)
3. Implement new CXStruct definition.
4. Add unit tests to all functions and more acceptance tests to achieve 100% coverage.

## Struct Definition
Implement new CXStruct definition.

Change CX type signatures
To eliminate CxArguments, we need int id codes for types (enum)
- type (with a sizeof function for type)
- offset for each variable

Simple types are like int, float, int64, etc.

## Struct-def Library
For interfacing CX to Golang. Golang can call functions in CX program that are exposed and CX program can also call the golang game API. Example, in a game created with Golang, there will be a feature/support/API to program an object with CX. 

1. CXStruct definition to Golang struct definition conversion and vice versa - this is needed for the function input arguments and returns. CX program takes a golang struct and read out the data from the golang struct to CX struct with a pointer.

Library reads in a golang struct definition and gets the 
- variable type
- offset of variable from start of struct
- size of variable

Objectives:
- Allow CX to read/write from golang structs directly using "unsafe".

Notes:
Support only for simple types, no pointers, only atomics

## Taskbar Launcher for CX
This is going to launch the CX Playground. It’s going to open up a web browser on the user’s computer. It’s going to let them type in a CX program and run it. This is going to be expanded into what’s called a CX-IDE or CX-Playground. That will be github.com/skycoin/cx-ide or github.com/skycoin/cx-playground. 

Use same tasbar library that we're using for skywire(https://github.com/getlantern/systray).

## CX Module format or Package format
1. Package format - This is how we’re representing the file. Every file has a length, a name, and a hash. Every module is a list of files – file structs. And then, we have a list of the package structs. And then, we can serialize it and hash that to get the ID for the whole program.
2. CX App Store - For storing CX programs/packages in a key-value store for downloading from website with simple API. To be an added feature in CX-Playground/CX-IDE. This will use CX Package Format in compiling cx program.

Objectives:
- To easily download a CX program from CX App Store website and run the CX program.

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

Objectives:
- allow multiple CX programs to be loaded/run at same time by one golang process
- cleanup code base

Notes:
- Restrict the compiler to only do one compilation at a time when it’s producing the AST and then force the encapsulation of the AST inside of a struct so that we can execute multiple concurrent CX programs at the same time.
- Compiler will still use global when outputing AST from source file, but lock compilation so that only one program can be compiled at a time. (Maybe use a channel to pass in data with compiled AST returned)