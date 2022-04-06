# 2022 CX Roadmap

## CX Refactoring
1. Give complete details to all functions and their arguments.
2. Implement new CXStruct definition.
3. Add unit tests to all functions and more acceptance tests to achieve 100% coverage.

## Integrate CX and Golang
Integrating CX and Golang. Golang can call functions in CX program that are exposed and CX program can also call the game API. Example in a game created with Golang, a feature/support/API to program an object with CX. 
1. CXStruct definition to Golang struct definition conversion and vice versa. This is needed for the function input arguments and returns. CX program takes a golang struct and read out the data from the golang struct to CX struct with a pointer.

## Taskbar Launcher for CX
This is going to launch the CX Playground. It’s going to open up a web browser on the user’s computer. It’s going to let them type in a CX program and run it. This is going to be expanded into what’s called a CX-IDE or CX-Playground. That will be github.com/skycoin/cx-ide or github.com/skycoin/cx-playground. 

Use same tasbar library that we're using for skywire(https://github.com/getlantern/systray).

## CX Module format or Package format
1. Package format - This is how we’re representing the file. Every file has a length, a name, and a hash. Every module is a list of files – file structs. And then, we have a list of the package structs. And then, we can serialize it and hash that to get the ID for the whole program.
2. CX App Store - a web server for CX or a CX Repository. To be an added feature in CX-Playground/CX-IDE. This will use CX Package Format in compiling cx program.

## Web Interface for CX game objects
To change the behavior live while the game is running, change the behavior, change the image, sprite, sound effects, etc. Example is, there's a list of agents/objects, 35 ships, while the game is running, I can change the ship to attack or run away through the web interface while the game is running. To be added as a feature in CX-Playground/CX-IDE

## CX Memory - a requirement for CX to run on an embedded system
1. Every variable and everything that needs to be defined has to be in that linear array as embedded system only have one memory stack.
2. the stack layout or the frame layout - A function that gives us, for a function body, (1) the list of all the variables and their name, (2) the size of each variable, and (3) their byte offset. If I have this list, I can take the list and I can say, “Give me all the I32s that are defined in this function.” 

## Function Scope
Push everything to functional scope, meaning that there is no local scopes within the functions. We may enforce local scoping rules so that, if you have a loop and there is a variable in the loop, when the loop finishes, that variable will still be defined, but you will not be able to access the variable outside the loop from the CX program.

## CX Program Encapsulation - Global Removal
Restrict the compiler to only do one compilation at a time when it’s producing the AST and then force the encapsulation of the AST inside of a struct so that we can execute multiple concurrent CX programs at the same time.

