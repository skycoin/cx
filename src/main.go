package main

import (
	"fmt"
	"log"
	"os"
	
	"bytes"
	"errors"
	"strconv"
	//"encoding/binary"
	"strings"
	"crypto/sha256"
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

// type Datum struct {
// 	dtype DataType
// 	value [DATUM_LENGTH]byte
// }

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
type State []Datum

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
	// Applies a series of operators to state
	lm_name := stat.lambda.name.name
	
	inputs := make([]*Datum, len(stat.inputs))
	outputs := make([]*Datum, len(stat.outputs))

	//Datum inputs
	for i := 0; i < len(stat.inputs); i++ {
		//inputs[i] = &(*state)[stat.inputs[i]]
		inputs[i] = &(*state)[stat.inputs[i]]
	}
	//Datum outputs
	for i := 0; i < len(stat.outputs); i++ {
		//outputs[i] = &(*state)[stat.outputs[i]]
		outputs[i] = &(*state)[stat.outputs[i]]
	}

	// checking for native functions
	switch lm_name {
	case "addInt32":
		// arg1, err1 := (*state)[stat.inputs[0]].ReadInt32()
		// arg2, err2 := (*state)[stat.inputs[1]].ReadInt32()

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
		
		//(*state)[stat.outputs[0]].WriteInt32(addInt32(args...))
	}




	
	// var lm_times2Int32 Lambda = Lambda{name: s_times2Int32.Intern(),
	// 	sigin: []DataType{_int32},
	// 	sigout: []DataType{_int32},
	// 	statements: make([]Statement, 1)}
	
	// statement2 := Statement{lambda: &lm_addInt32,
	// 	inputs: []int{0, 0},
	// 	outputs: []int{1}}
	// lm_times2Int32.statements[0] = statement2


	
	// if not native, we need to execute current lambda statements


	//*[1][2][3][4][5][6]
	//substate := State(append(inputs, outputs...))
	//stat.lambda.Execute(&substate)
	stat.lambda.Execute(state)
}

	// var s_times2Int32 Symbol = Symbol{name: "times2"}
	// var lm_times2Int32 Lambda = Lambda{name: s_times2Int32.Intern(),
	// 	sigin: []DataType{_int32},
	// 	sigout: []DataType{_int32}}
	
	// statement2 := Statement{lambda: &lm_addInt32,
	// 	inputs: []int{0, 0},
	// 	outputs: []int{1}}
	// lm_times2Int32.statements[0] = statement2

func (lm *Lambda) Execute(state *State) {
	for i := 0; i < len(lm.statements); i++ {
		lm.statements[i].Execute(state)
	}
	
	// we don't need to return the state, we simply examine its contents in the
	// context where we defined it
	//return state
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




	
	
	


	// In a signature space, for example (int32, int32, int32) (int32)


	
	// &vtable{space: make([]byte, 0),
	// 	sigin: lm.sigin,
	// 	sigout: lm.sigout}

	//lm_vt := Lambdas[sig]


	return nil
}

func dbg(elt string) {
	fmt.Println(elt)
}

func main() {

	var state State
	var dat1 Datum
	var dat2 Datum
	var dat3 Datum
	dat1.WriteInt32(int32(50))
	dat2.WriteInt32(int32(20))
	dat3.WriteInt32(int32(3300))
	
	state = append(state, dat1)
	state = append(state, dat2)
	state = append(state, dat3)
	
	foo := Symbol{name: "foo"}
	bar := Symbol{name: "bar"}
	tar := Symbol{name: "tar"}
	mar := Symbol{name: "mar"}
	foo.Intern()
	bar.Intern()
	tar.Intern()
	fmt.Println(mar.Intern())

	fmt.Println(Symbol_vt.Delegate())
	

	fun := Lambda{name: &foo,
		sigin: []DataType{_uint32, _float32},
		sigout: []DataType{_uint32}}
	
	fmt.Println(fun)

	fun.Defun()


	// crypto
	sum := sha256.Sum256([]byte("123"))
	fmt.Println(len(sum))
	//fmt.Printf("%x", h.Sum(nil))

	// type Datum struct {
	// 	dtype DataType
	// 	value [DATUM_LENGTH]byte
	// }
	// type State []Datum

	//bs := make([]byte, 100)
	//value := int64(10)
	//binary.LittleEndian.PutUint64(bs, value)
	//binary.Write(bs, binary.LittleEndian, value)

	value := strconv.FormatInt(10, 10)
	
	fmt.Println([]byte(value))
	fmt.Println(strconv.Atoi(string([]byte(value))))
	
	//func (datum *Datum) WriteInt32 (num int32) (*Datum) {

	var dat *Datum = new(Datum)
	dat.WriteInt32(int32(1234567890))

	fmt.Println(dat.dtype)
	fmt.Println(dat.value)

	val32, err := dat.ReadInt32()
	if err == nil {
		fmt.Println(val32 + 1)
	}

	// type Lambda struct {
	// 	//name *Symbol
	// 	//change to Symbol later
	// 	name string
	// 	//offset int //Necessary?
	// 	sigin []DataType
	// 	sigout []DataType
	// 	statements []Statement
	// }

	var s_addInt32 Symbol = Symbol{name: "addInt32"}	

	// This IS required in order to initialize the system with native functions
	var lm_addInt32 Lambda = Lambda{name: s_addInt32.Intern(),
		sigin: []DataType{_int32, _int32},
		sigout: []DataType{_int32}}
	
	statement1 := Statement{lambda: &lm_addInt32,
		inputs: []int{0, 1},
		outputs: []int{2}}

	//statements are linked to their current context (space) in a lambda
	//in this case, the statement below ONLY works for times2 lambda OR
	//for another lambda where the state has a similar structure
	var s_times2Int32 Symbol = Symbol{name: "times2"}
	var lm_times2Int32 Lambda = Lambda{name: s_times2Int32.Intern(),
		sigin: []DataType{_int32},
		sigout: []DataType{_int32},
		statements: make([]Statement, 1)}
	
	statement2 := Statement{lambda: &lm_addInt32,
		inputs: []int{0, 0},
		outputs: []int{1}}
	lm_times2Int32.statements[0] = statement2

	dbg(".....")
	lm_times2Int32.Execute(&state)
	fmt.Println(state)

	//Problem with space: calling a lambda inside another lambda
	//Solution: we need to create another array of *Datums (subspace)

	
	

	

	fmt.Println(statement2)
	
	fmt.Println(statement1.lambda.name.name)
		
	fmt.Println(lm_addInt32)
	
	// type Statement struct {
	// 	lambda *Lambda
	// 	inputs []int //Positions in the state
	// 	outputs []int //Positions in the state
	// }

	fmt.Println(state)
	statement1.Execute(&state)
	fmt.Println(state)
	
	// //fmt.Println(Symbol_vt.buffer.Bytes())
	// //fmt.Fprintf(Symbol_vt.buffer, "hello")
	// //Symbol_vt.buffer.WriteTo(os.Stdout)

	
	// // sym, err := Symbol_vt.buffer.ReadString(delimr)
	// // if err == nil {
	// // 	fmt.Println(sym[0:len(sym)-1])
	// // }
	
	// // Symbol_vt.buffer.WriteTo(os.Stdout)
	// //binary.Write(Symbol_vt.buffer)


	// var symb *Symbol
	// symb.Intern("bar")
}
