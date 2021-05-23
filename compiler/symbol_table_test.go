package compiler

import "testing"

func TestDefine(t *testing.T) {

	expected := map[string]Symbol{

		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 0},
		"c": Symbol{Name: "c", Scope: LocalScope, Index: 0},
		"d": Symbol{Name: "d", Scope: LocalScope, Index: 0},
		"e": Symbol{Name: "e", Scope: LocalScope, Index: 0},
		"f": Symbol{Name: "f", Scope: LocalScope, Index: 0},
	}

	global := NewSymbolTable()

	a := global.Define("a")

	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b")

	if a != expected["b"] {
		t.Errorf("expected a=%+v, got=%+v", expected["b"], a)
	}

	firstLocal := NewEnclosedSymbolTable(global)

	c := firstLocal.define("c")

	if c != expected["c"] {
		t.Errorf("expected c=%+v", expected["c"], c)
	}

}
