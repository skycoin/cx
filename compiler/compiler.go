package compiler

type compiler struct {
	constants   []object.Object
	symbolTable *symbolTable
	scopes      []CompilationScope
	scopeIndex  int
}
