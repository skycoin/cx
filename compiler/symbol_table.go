package compiler

type symbolScope string

const (
	LocalScope    symbolScope = "Local"
	GlobalScope   symbolScope = "GLOBAL"
	BuildinScope  symbolScope = "BUILTIN"
	FreeScope     symbolScope = "FREE"
	FunctionScope symbolScope = "FUNCTION"
)

type Symbol struct {
	Name  string
	Scope symbolScope
	INDEX int
}

type symbolTable struct {
	Outer          *symbolTable
	store          map[string]Symbol
	numDefinations int
	FreeSymbols    []Symbol
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {

	s := NewSymbolTable()

	s.Outer = outer

	return s
}

func NewSymbolTable() *symbolTable {

	store := make(map[string]Symbol)
	free := []Symbol{}

	symbolTable := symbolTable{store: store, FreeSymbols: free}

	return &symbolTable
}

func (s *symbolTable) Define(name string) Symbol {

	symbol := Symbol{Name: name, Index: s.numDefinations}

	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {

		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinations++
	return symbol
}

func (s *symbolTable) Resolve(name string) (Symbol, bool) {

	obj, ok := s.store[name]

	if !ok && s.Outer != nil {

		obj, ok = s.Outer.Resolve(name)

		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuildinScope {
			return obj, ok
		}

		free := s.defineFree(obj)

		return free, ok
	}

	return obj, ok
}

func (s *symbolScope) defineFree(original Symbol) Symbol {

	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}

	symbol.Scope = FreeScope

	s.store[original.Name] = symbol

	return symbol
}

func (s *symbolTable) DefineFunctionName(name string) Symbol {

	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope}

	s.store[name] = symbol

	return symbol
}

func (s *symbolTable) DefineBuiltin(index int, name string) Symbol {

	symbol := Symbol{Name: name, Index: index, Scope: BuildinScope}
	s.store[name] = symbol
	return symbol
}
