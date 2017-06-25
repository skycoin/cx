package main

import (
	"fmt"
	"log"
	"os"

	//"generator"
	"parser"
	"lexer"
	"go/token"
	"go/ast"
	"go/printer"
	
	"bytes"
	"errors"
	"strconv"
	"strings"
	//"crypto/sha256"
)

const SYMBOL_LENGTH int = 32
const DATUM_LENGTH int = 32

func init() {
	log.SetOutput(os.Stdout)
}

type DataType int
const (
	_bool DataType = iota
	_int8
	_int16
	_int32
	_int64
	_uint8
	_uint16
	_uint32
	_uint64
	_float32
	_float64
)

type vtable struct {
	space []byte
	sigin []DataType
	sigout []DataType
	//children []*vtable // Children are going to be used to register a function according to its signature
	// Nope, I think we don't need it
}

type Symbol struct {
	name string
	offset int
}

// LGP
// Modules
// Signatures
// []byte
// we can have a hash table representing a directory of the functions. lookup by name, instructions are returned

/*
--Symbols
--Delegate vtable (create new)
Statements
Functions
Bindings
+Basic Types
*/

func (vt *vtable) Delegate() (*vtable) {
	// We wouldn't use the vtable exactly as in a COLA
	// Specifically, we won't be using inheritance in the same way
	// What we will use is their byte space idea for encapsulation
	var child vtable
	
	child.space = make([]byte, 0)
	//child.parent = vt
	
	return &child
}


// NOPE, we're using the Datum struct for a generalized representation
// Datum is the CX struct Brandon mentioned
// type Int32 struct {
// 	// These types are going to be stored in a State structure
// 	// They are going to store values that represent this type
// 	// Each statement knows how to handle the data contained in these types
// 	val int32
// 	// For basic types it might not make much sense
// }

// We would create several vtables to encapsulate types of objects with the same offset (fixed structs)
var Object_vt *vtable = &vtable{space: make([]byte, 0)}
var Symbol_vt *vtable = &vtable{space: make([]byte, 0)}

func (datum *Datum) ReadBool () (bool, error) {
	if datum.dtype == _bool {
		if datum.value[0] == byte(1) {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, errors.New("Datum is not of type Bool")
}

func (datum *Datum) WriteBool (b bool) (*Datum) {
	var bs [DATUM_LENGTH]byte
	
	datum.dtype = _bool
	
	if b {
		bs[0] = byte(1)
		
	} else {
		bs[0] = byte(0)
	}
	
	datum.value = bs
	
	return datum
}

// Extremely inefficient way, but enough for testing. Change to bit math maybe
func (datum *Datum) ReadInt32 () (int32, error) {
	if datum.dtype == _int32 {
		var i int
		for i = 0; i < DATUM_LENGTH; i++ {
			if datum.value[i] == byte(0) {
				break
			}
		}
		value, err := strconv.ParseInt(string(datum.value[:i]), 10, 32)
		return int32(value), err
	}
	return 0, errors.New("Datum is not of type Int32")
}

// Extremely inefficient way, but enough for testing. Change to bit math maybe
func (datum *Datum) WriteInt32 (num int32) (*Datum) {
	var bs [DATUM_LENGTH]byte
	value := []byte(strconv.FormatInt(int64(num), 10))

	for i := 0; i < len(value); i++ {
		bs[i] = value[i]
	}

	datum.dtype = _int32
	datum.value = bs
	
	return datum
}

// Symbols are used in order to determine if a Symbol is already defined
func (sym *Symbol) Intern() (*Symbol) {
	var i int
	for i = 0; i < len(Symbol_vt.space) ; i += SYMBOL_LENGTH {
		next_symbol := Symbol_vt.space[i:i+SYMBOL_LENGTH]

		
		if bytes.Equal([]byte(sym.name), next_symbol[0:len(sym.name)]) {
			sym.offset = i
			return sym
		}
	}

	sym_bytes := append([]byte(sym.name), make([]byte, SYMBOL_LENGTH - len(sym.name))...)
	Symbol_vt.space = append(Symbol_vt.space, sym_bytes...)
	sym.offset = i
	return sym
}

func strSig(sigin []DataType, sigout []DataType) string {
	var buffer bytes.Buffer
	
	buffer.WriteString(strings.Trim(strings.Replace(fmt.Sprint(sigin), " ", "", -1), "[]"))
	buffer.WriteString(",")
	buffer.WriteString(strings.Trim(strings.Replace(fmt.Sprint(sigout), " ", "", -1), "[]"))
	
	return buffer.String()
}


type Signature string
// In doubt. We need to think more about this
type SHA256 [32]byte


// Holds inputs, outputs and variables inside a lambda
type Datum struct {
	dtype DataType
	value [DATUM_LENGTH]byte
}
type State []*Datum

type Lambda struct {
	name *Symbol
	//offset int //Necessary?
	sigin []DataType
	sigout []DataType
	statements []Statement
}

type Statement struct {
	lambda *Lambda
	inputs []int //Positions in the state
	outputs []int //Positions in the state
}

func addInt32(nums ...int32) (int32) {
	total := int32(0)
	for _, num := range nums {
		total += num
	}
	return total
}

func allClear(e []error) bool {
	for i := 1; i < len(e); i++ {
		if e[i] != nil {
			return false
		}
	}
	return true
}

func (stat *Statement) Execute(state *State) {
	if stat.lambda == nil {
		return // A statement ALWAYS needs a lambda
	}
	lm_name := stat.lambda.name.name

	if stat.outputs == nil {
		return
	}
	
	inputs := make([]*Datum, len(stat.inputs))
	outputs := make([]*Datum, len(stat.outputs))
	
	//Datum inputs
	for i := 0; i < len(stat.inputs); i++ {
		inputs[i] = (*state)[stat.inputs[i]]
	}
	//Datum outputs
	for i := 0; i < len(stat.outputs); i++ {
		outputs[i] = (*state)[stat.outputs[i]]
	}

	// checking for native functions
	// this part is getting very bloated with repetitive code. refactor
	switch lm_name {
	case "addInt32":
		args := make([]int32, len(inputs))
		errs := make([]error, len(inputs))
		
		for i := 0; i < len(inputs); i++ {
			args[i], errs[i] = inputs[i].ReadInt32()
		}

		if allClear(errs) {
			for i := 0; i < len(outputs); i++ {
				outputs[i].WriteInt32(addInt32(args...))
			}
		}
	case "timesInt32":
		args := make([]int32, len(inputs))
		errs := make([]error, len(inputs))
		
		for i := 0; i < len(inputs); i++ {
			args[i], errs[i] = inputs[i].ReadInt32()
		}

		if allClear(errs) {
			for i := 0; i < len(outputs); i++ {
				outputs[i].WriteInt32(args[0] * args[1])
			}
		}
	case ">":
		args := make([]int32, len(inputs))
		errs := make([]error, len(inputs))

		for i := 0; i < len(inputs); i++ {
			args[i], errs[i] = inputs[i].ReadInt32()
		}

		if allClear(errs) {
			for i := 0; i < len(outputs); i++ {
				outputs[i].WriteBool(args[0] > args[1])
			}
		}
	case "==":
		args := make([]int32, len(inputs))
		errs := make([]error, len(inputs))

		for i := 0; i < len(inputs); i++ {
			args[i], errs[i] = inputs[i].ReadInt32()
		}

		if allClear(errs) {
			for i := 0; i < len(outputs); i++ {
				outputs[i].WriteBool(args[0] == args[1])
			}
		}
	default:
		// not native function. executing lambda with substate
		substate := State(append(inputs, outputs...))
		stat.lambda.Execute(&substate)
	}
}

func (lm *Lambda) Execute(state *State) {
	lm_name := ""
	if lm.name != nil {
		lm_name = lm.name.name
	}
	
	switch lm_name {
		// If statement could never be treated as a single statement, as it is actually defined by a group of statements (thus Lambda.Execute needs to handle it)
	case "if":
		stat := lm.statements[0] //It NEEDS to be the first statement of an if
		inputs := make([]*Datum, len(stat.inputs))
		outputs := make([]*Datum, len(stat.outputs))

		for i := 0; i < len(stat.inputs); i++ {
			inputs[i] = (*state)[stat.inputs[i]]
		}
		for i := 0; i < len(stat.outputs); i++ {
			outputs[i] = (*state)[stat.outputs[i]]
		}
		
		var predicate_result *Datum = new(Datum)
		predicate_result.dtype = _bool
		var tmp_state State = append(inputs, predicate_result)

		stat.outputs = []int{2}
		stat.Execute(&tmp_state)

		result, err := tmp_state[2].ReadBool()

		if err == nil {
			if result {
				if lm.statements[1].lambda != nil {
					lm.statements[1].Execute(state)
				}
			} else if lm.statements[2].lambda != nil {
				lm.statements[2].Execute(state)
			}
		}
	default:
		for i := 0; i < len(lm.statements); i++ {
			lm.statements[i].Execute(state)
		}
	}
	
	
	// we don't need to return the state, we simply examine its contents in the
	// context where we defined it
}

func (lm *Lambda) Defun() (*Lambda) {
	// Interns a lambda
	sig := Signature(strSig(lm.sigin, lm.sigout))
	fmt.Println(sig)

	// Create vtable if it doesn't exist
	// if _, ok := Lambdas[sig]; !ok {
	// 	Lambdas[sig] = make(map[SHA256][]byte, 0)
	// }

	// ToDo: We hash the native names, e.g. "ADD" => [166 ...]
	// The value is "ADD" (or a Symbol? The symbol could could)
	// The key is the hash of "ADD"
	// This needs to be initialized by us
	//
	// A program is going to be an array of strings in Lambda struct
	///for example: ["ADD", "SUB"]
	// The key of this program is the hash of these strings
	// The value is an array of the hashes of the strings
	//
	// This way, in order to construct the final program, we only need
	///to look for all the hashes until we have only native functions
	//Lambdas[sig][/**/0]

	// ToDo: How to apply to program state (data)
	// and:  how to implement program flow (if, while)
	///for if: execute condition to determine what nodes to compute
	///for while: the statements remain the same, only state (data) changes

	// Registers (?): Inputs and Outputs
	//
	// Registers need to be part of the generated program
	// (int32, int32, int32, int32, int32) /* Number of features */ (int32) /* Output */
	/// The program needs to communicate with the register space to allocate
	// Registers are going to be in []byte spaces
	// We initialize the []byte space (inputs)
	// The first (or last?) [n]byte is the output(s)
	//
	// We need to save the result somewhere too (and we need to indicate this when saving the lambda)




	// randomly generated
	// or user generated
	// r[5] = ADD(r[0], r[3])
	// r[0] = SUB(r[7], r[5])
	// r[5] = SUB(r[2], 72) // This lambda is represented as 5,SUB,2,Uint32
	// No! We just send this function an input space
	/// But now


	// As we noted before, a lambda needs to be represented differently than a statement
	// Just like in regular programming, a lambda is going to be a generalized representation
	/// e.g. ADD(Uint32, Uint32)
	// Now, a statement is, for example R[0] = ADD(R[5], R[2])
	// In here, a statement is always be assigned to a register.
	// To construct a statement, we can have a tuple in the _____ space, where the first entry is the function
	/// The second entry is the assignment variable, and the rest are the arguments to the function




	// What are the necessary elements until now:
	// Drop symbols; these are just going to be
	//////// In LGP we wouldn't need these, but what about general programming?
	//////// In general programming, we are going to parse the symbols into register indexes (R[1...3])
	/////// This means that we need to be able to manage a variable amount of registers (i.e. len(registers) != len(inputs), as in classic LGP)

	return nil
}

func dbg(elt string) {
	fmt.Println(elt)
}

func dbgState(s *State, desc string) {
	dbg(desc)

	// determining bytes used in byte slice
	// can and should be used later in other parts of the code
	// (or maybe not? seems slow)
	
	for i:= 0; i < len((*s)); i++ {
		var j int
		for j = 0; j < len((*s)[i].value); j++ {
			if (*s)[i].value[j] == byte(0) {
				break
			}
		}
		fmt.Println(string((*s)[i].value[:j]))
	}
}

// func MakeSymbol(name string) *Symbol {
// 	&Symbol{name: name}.Intern()
// }

// func MakeLambda(name string, sigin []DataType, sigout []DataType, statements ...Statement) {
// 	return &Lambda{name: makeSymbol(name),
// 		sigin: sigin, //[]DataType{_int32, _int32}
// 		sigout: sigout, //[]DataType{_int32}
// 		statements: statements}
// }

func main() {
	foo := Symbol{name: "foo"}
	bar := Symbol{name: "bar"}
	tar := Symbol{name: "tar"}
	mar := Symbol{name: "mar"}
	
	foo.Intern()
	bar.Intern()
	tar.Intern()
	mar.Intern()

	// crypto
	//sum := sha256.Sum256([]byte("123"))
	//fmt.Println(len(sum))
	//fmt.Printf("%x", h.Sum(nil))

	var dat *Datum = new(Datum)
	dat.WriteInt32(int32(1234567890))

	fmt.Println(dat.dtype)
	fmt.Println(dat.value)

	val32, err := dat.ReadInt32()
	if err == nil {
		fmt.Println(val32 + 1)
	}

	// This is required in order to initialize the system with native functions
	/* Start Native Function addInt32 */
	var s_addInt32 Symbol = Symbol{name: "addInt32"}
	var lm_addInt32 Lambda = Lambda{name: s_addInt32.Intern(),
		sigin: []DataType{_int32, _int32},
		sigout: []DataType{_int32}}
	/* End Native Function addInt32 */
	

	//statements are linked to their current context (space) in a lambda
	//in this case, the statement below ONLY works for times2 lambda OR
	//for another lambda where the state has a similar structure


	/* Start Function Times2Int32 */
	var s_times2Int32 Symbol = Symbol{name: "times2"}
	var lm_times2Int32 Lambda = Lambda{name: s_times2Int32.Intern(),
		sigin: []DataType{_int32},
		sigout: []DataType{_int32},
		statements: make([]Statement, 1)}
	
	statement2 := Statement{lambda: &lm_addInt32,
		inputs: []int{0, 0},
		outputs: []int{1}}
	lm_times2Int32.statements[0] = statement2
	/* End Function Times2Int32 */

	/* Start Native Function >, greater than (gt) */
	var s_gt Symbol = Symbol{name: ">"}
	var lm_gt Lambda = Lambda{name: s_gt.Intern(),
		sigin: []DataType{_int32, _int32},
		sigout: []DataType{_bool}}
	/* End Native Function >, greater than (gt) */

	dbg(".....")

	var state State =
		[]*Datum {
		new(Datum).WriteInt32(int32(50)),
		new(Datum).WriteInt32(int32(20)),
		new(Datum).WriteInt32(int32(3300))}

	//Testing calling a single statement
	statement1 := Statement{lambda: &lm_addInt32,
		inputs: []int{0, 1},
		outputs: []int{2}}
	
	dbgState(&state, "Initial state")
	
	statement1.Execute(&state)

	dbgState(&state, "Sum of first and second inputs")

	//Testing calling a lambda
	
	lm_times2Int32.Execute(&state)

	dbgState(&state, "Double of first input and writes to second slot in state")

	//Lambda which doubles each of two inputs and then sums them

	var state_sumDoubles State =
		[]*Datum {
		new(Datum).WriteInt32(int32(10)),
		new(Datum).WriteInt32(int32(30)),
		new(Datum).WriteInt32(int32(0))}
	
	var s_sumDoublesInt32 Symbol = Symbol{name: "sumDoubles"}
	var lm_sumDoublesInt32 Lambda = Lambda{name: s_sumDoublesInt32.Intern(),
		sigin: []DataType{_int32, _int32},
		sigout: []DataType{_int32},
		statements: make([]Statement, 3)}
	
	stat_times2 := Statement{lambda: &lm_times2Int32}
	stat_sum := Statement{lambda: &lm_addInt32}

	// Doubling first input
	stat_times2.inputs = []int{0}
	stat_times2.outputs = []int{0}
	lm_sumDoublesInt32.statements[0] = stat_times2
	// Doubling second input
	stat_times2.inputs = []int{1}
	stat_times2.outputs = []int{1}
	lm_sumDoublesInt32.statements[1] = stat_times2

	// Summing input1 and input2
	stat_sum.inputs = []int{0, 1}
	stat_sum.outputs = []int{2} // Writing result to third datum
	lm_sumDoublesInt32.statements[2] = stat_sum

	dbgState(&state_sumDoubles, "Created another state")
	lm_sumDoublesInt32.Execute(&state_sumDoubles)
	dbgState(&state_sumDoubles, "Doubles each input, then sums both inputs and writes the sum to third slot")

	// testing lambda without a name (duh, an actual lambda)
	var lambda Lambda = Lambda{//name: symbol.Intern(),
		sigin: []DataType{_int32},
		sigout: []DataType{_int32},
		statements: make([]Statement, 10)}

	// Doubling input
	stat_times2.inputs = []int{0}
	stat_times2.outputs = []int{0}
	// doubling first input 10 times
	for i := 0; i < len(lambda.statements); i++ {
		lambda.statements[i] = stat_times2
	}

	// sending same state as previous example for convenience
	lambda.Execute(&state_sumDoubles)
	dbgState(&state_sumDoubles, "Testing nameless function (actual lambda), doubles first slot 10 times and writes each result to first slot")

	fmt.Println("Testing Boolean datums")
	true_dat := new(Datum).WriteBool(true)
	fmt.Println(true_dat.ReadBool())
	false_dat := new(Datum).WriteBool(false)
	fmt.Println(false_dat.ReadBool())


	fmt.Println("Testing if statement")

	// Creating another state to test if statement
	var state_if State =
		[]*Datum {
		new(Datum).WriteInt32(int32(100)),
		new(Datum).WriteInt32(int32(30))} // Notice that we don't need a third slot to hold the result of predicate. This result is stored in a new, temporary state
	
	// if temporary prototype
	var s_if Symbol = Symbol{name: "if"}
	var lm_if Lambda = Lambda{name: s_if.Intern(),
		sigin: []DataType{}, // bool input, or maybe no input. if won't change main state
		sigout: []DataType{}, // no output, as "if" itself won't change state
		statements: make([]Statement, 3)} // This could actually work: we are making room for 3 statements that every if lambda should have, but we haven't initialized them




	// (lm [int32 int32] [int32]
	// 	(+ %0 %1))
	


	
	// predicate
	stat_boolean := Statement{lambda: &lm_gt,
		inputs: []int{0, 1},
		//outputs: []int{2} // We won't use this, as it's going to be used for an if statement
		// If we wanted to store the value, we need to provide an output slot
		// If it's not for an if statement and we don't provide an output, the result should simply be discarded
	}
	// if slot0 > slot1, then we double slot0
	stat_ifthen := Statement{lambda: &lm_times2Int32,
		inputs: []int{0},
		outputs: []int{0}}
	// if not, then we double slot1
	stat_else := Statement{lambda: &lm_times2Int32,
		inputs: []int{1},
		outputs: []int{1}}

	lm_if.statements[0] = stat_boolean
	lm_if.statements[1] = stat_ifthen
	lm_if.statements[2] = stat_else

	dbgState(&state_if, "New state for testing if statement")
	lm_if.Execute(&state_if)
	dbgState(&state_if, "If slot1 is greater than slot 2, double slot1, if not, double slot2")



	
	var state_factorial State =
		[]*Datum {
		new(Datum).WriteInt32(int32(5)),
		new(Datum).WriteInt32(int32(4)),
		
		new(Datum).WriteInt32(int32(-1)),
		new(Datum).WriteInt32(int32(0))} // We should only need 2 slots, but we can't use non-slot values at the moment. For example, if we need the value of PI = 3.14159, this value needs to be in a slot in a state

	// Move these to initialization
	var lm_timesInt32 Lambda = Lambda{name: (&Symbol{name: "*"}).Intern(),
		sigin: []DataType{_int32, _int32},
		sigout: []DataType{_int32}}
	
	var lm_equalInt32 Lambda = Lambda{name: (&Symbol{name: "=="}).Intern(),
		sigin: []DataType{_int32, _int32},
		sigout: []DataType{_int32}}

	stat_equal := Statement{lambda: &lm_equalInt32,
		inputs: []int{1, 3},
		//outputs: []int{1} //for if statement
	}
	
	stat_times := Statement{lambda: &lm_timesInt32,
		inputs: []int{0, 1},
		outputs: []int{0}}

	stat_decrement := Statement{lambda: &lm_addInt32,
		inputs: []int{1, 2},
		outputs: []int{1}}

	var lm_factorial Lambda = Lambda{name: (&Symbol{name: "factorial"}).Intern(),
		sigin: []DataType{_int32, _int32, _int32, _int32}, // we should only need 2, but we need extra features (see above)
		sigout: []DataType{_int32},
		statements: make([]Statement, 3)}

	stat_factorial := Statement{lambda: &lm_factorial,
		inputs: []int{0, 1, 2, 3},
		outputs: []int{0}}
	
	// re-using if lambda from above
	lm_if.statements[0] = stat_equal
	lm_if.statements[1] = Statement{} // clearing ifthen statement
	lm_if.statements[2] = stat_factorial

	lm_factorial.statements[0] = stat_times
	lm_factorial.statements[1] = stat_decrement
	lm_factorial.statements[2] = stat_factorial

	dbgState(&state_factorial, "New state for testing factorial lambda")
	//lm_if.Execute(&state_factorial)
	dbgState(&state_factorial, "Slot0 should be 120")



	l := lexer.Lex("hello", "(ns main \"fmt\") (def foo (fn [[bar]] (let [[n 40]] (+ 10 10))))")
	p := parser.Parse(l)
	// a := generator.GenerateAST(p)
	// fset := token.NewFileSet()
	// ast.Print(fset, a)

	// var buf bytes.Buffer
	// printer.Fprint(&buf, fset, a)

	fmt.Println(p)

	// Creating an AST manually

	anyType := &ast.SelectorExpr{
		X:   ast.NewIdent("core"),
		Sel: ast.NewIdent("Any"),
	}
	
	fmt.Println("...")
	a := &ast.File{Name: ast.NewIdent("testing")}
	decls := make([]ast.Decl, 1)

	decls[0] = &ast.FuncDecl{
		Name: ast.NewIdent("foo"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("bar"),
						},
						Type: anyType,
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("bar"),
						},
						Type: anyType,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BinaryExpr{
						X: ast.Expr(&ast.BasicLit{
							Kind: token.INT,
							Value: "10",
						}),
						Y: ast.Expr(&ast.BasicLit{
							Kind: token.INT,
							Value: "20",
						}),
						Op: token.ADD}},
			},
		},
	}

	fset := token.NewFileSet()
	ast.Print(fset, a)

	var buf bytes.Buffer
	printer.Fprint(&buf, fset, a)

	fmt.Println("...")
	fmt.Printf("%s\n", buf.String())
	
	// func factorial (slot0 int32, slot1 int32) (slot0 int32) {
	// 	slot0 := slot0 * slot1
	// 	dec(slot1) // This will be slot1 := slot1 - slot2 for now
	// 	if x == 0 {
	
	// 	} else {
	// 		factorial(slot0, slot1)
	// 	}
	// }

	// // We allocate a state or scope with two slots: 5 and 4, and send this state to factorial
	// factorial(5, 4) 
	
	// print(slot0)	

}
