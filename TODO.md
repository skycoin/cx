TODO.md

# Todo List for CX

This is a list of notes for compiler issues to fix

## Local Variable Declarations can Shadow Global Variables

If we have global variable declared in scope, a local variable can overshadow the global variable.

We need to throw warning when this occurs.

We need to verify the existing behavior and how it should be handled in CX.

Generally, locally declared variables should not shadow globally declared variables in same package. So it should throw warning in compiler.

## Exposing GC functions to user

We need to expose the GC functions to user. For triggering GC operations during waits.

We need to be able to cap the GC execution time.

We need to be able to run GC operations in a second thread.

We need to implement memory remapping and movement of objects during run time.