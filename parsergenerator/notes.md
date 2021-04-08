notes

---

---

Seperate out pipelines
- lexer -> parser -> ast -> VM execution of AST
- lexer -> parser -> ast -> Compile to Golang -> Compiler Golang to Binary

1> Where is the stage where we compile to golang?
- can this be moved to its own module?

2> Where is the code for the AST inpreter (the VM)?

3> Can the AST construction be moved out from module that has the compiler?

---

The compiler needs to be seperated into stages
- stage1
- stage2
- stage3