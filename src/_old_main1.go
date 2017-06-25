package main

// Core ideas of the Object Model in a COLA
//
// The implementation in C pays too much attention to memory allocation in order to provide a far more efficient process.
// For the implementation in Go, we need to focus on the core ideas of the object model in a COLA, which are:
// Instead of making a distinction between state and behaviour, this object model decides to use behaviour to represent both.
// If one wants to access state, it can be done through accessors (methods)
// The kernel of the object model is a virtual table (vtable) which has the most primitive methods for the root object
// These primitive methods are: lookup, addMethod, allocate, and delegate
// lookup: gets a method according to a message name
// addMethod: adds a new method to the object's vtable
// allocate: stores the state of an object in memory, returns the objects address, casted to a object type. This is important and was hard to understand; this way of storing state ensures that vtables store only objects adresses, without the need to be aware of the payload structure or size
// delegate: used to create a subtype (inheritance) of an object. When an object's vtable doesn't contain a method (behaviour), its parent is searched for that method.
//
// The most primitive entities required for the object model are: object, symbol, and vtable.
// Vtables are explained above.
// An object is simply a link between the stored state in memory and the behaviours stored in vtables.
// A symbol is an object used to lookup a method in a vtable. Its only absolutely necessary method or behaviour is "intern"
//
// Implementation in Go:
//
// As Go doesn't provide much power for manually allocating and manipulating memory as C (as far as I know), the vtables shall be implemented as slices, where the memory is automatically allocated and freed by Go's garbage collector
//
// Objects need to be referenced by pointers


// Update:
// Data is going to be stored in byte slices. Each Object from that type holds the offset for the slice of bytes that is storing the data

import (
	"fmt"
	"os"
	"log"
	"errors"
	"reflect"
	
	"bytes"
	"encoding/binary"
	
	_ "object"
)

func init() {
	log.SetOutput(os.Stdout)
}

type Vtable struct {
	kv map[interface{}]interface{}
	parent *Vtable
}

type Object interface {}

// type Closure struct {
// 	// We need a way to store functions
// }

type Symbol struct {
	//name string
}

// for testing access to object with state
type Number struct {
	num int
}

var SymbolList *Vtable

var Vtable_vt *Vtable
var Object_vt *Vtable
var Symbol_vt *Vtable
var Closure_vt *Vtable

var S_addMethod interface{}
//var S_allocate Object
var S_delegated interface{}
var S_lookup interface{}

//
//func (obj Symbol)

func Vtable_delegated(self *Vtable) (*Vtable) {
	var child Vtable
	
	child.kv = make(map[interface{}]interface{})
	child.parent = self
	
	return &child
}

//fmt.Println(reflect.TypeOf(fun1).Kind() == reflect.TypeOf(fun2))

func Vtable_addMethod(self *Vtable, key interface{}, method interface{}) (interface{}) {
	if m, ok := self.kv[key]; ok {
		return m
	} else {
		self.kv[key] = method
		return method
	}
}
// here
func Vtable_lookup(self *Vtable, key interface{}) (interface{}, error) {
	if method, ok := self.kv[key]; ok {
		return method, nil
	} else if self.parent != nil {
		return Vtable_lookup(self.parent, key)
	} else {
		return nil, errors.New("Method not found")
	}
}

// func Symbol_new(name string) (Object) {
// 	var sym Symbol
// 	sym.name = name
// 	return &sym
// }

func Symbol_intern(name string) (interface{}) {
	if sym, ok := SymbolList.kv[name]; ok {
		return &sym
	} else {
		var sym Symbol
		//sym.name = name
		SymbolList.kv[name] = &sym
		return &sym
	}
}

func Send(receiver *Vtable, message interface{}, arguments ...interface{}) (interface{}, error) {
	method, err := Vtable_lookup(receiver, message)	
	if err == nil {
		var numIn int = reflect.TypeOf(method).NumIn()
		var numOut int = reflect.TypeOf(method).NumOut()
		
		if len(arguments) != numIn {
			return nil, errors.New("Wrong number of input parameters")
		}

		for i := 0; i < numIn; i++ {
			if reflect.TypeOf(arguments[i]) != reflect.TypeOf(method).In(i) &&
				reflect.TypeOf(method).In(i).String() != "interface {}" {
				return nil, errors.New("Input type mismatch")
			}
		}

		for i := 0; i < numOut; i++ {

			reflect.TypeOf(method).Out(i)
			
			// if reflect.TypeOf(arguments[i]) != reflect.TypeOf(method).Out(i) &&
			// 	reflect.TypeOf(method).Out(i).String() != "interface {}" {
			// 	return nil, errors.New("Output type mismatch")
			// }
		}

		v := reflect.ValueOf(method)
		rargs := make([]reflect.Value, len(arguments))
		for i, a := range arguments {
			rargs[i] = reflect.ValueOf(a)
		}
		
		// if numOut > 0 {
		// 	outType = reflect.TypeOf(method).Out(0)
		// }
		
		return v.Call(rargs), nil
		//return v.Call(rargs), reflect.TypeOf(method).Out(0), nil
		
	} else {
		return nil, errors.New("Method not found")
	}
}

func foobar() {
	fmt.Println("hello from foobar")
}
func barfoo(name string) {
	fmt.Println("hello" + name)
}



/////////

type _Vtable struct {
	buffer *bytes.Buffer
	parent *Vtable
}

type Symbol struct {
	//name string
}

Symbol_vt

func _Symbol_intern(name string) (interface{}) {
	if sym, ok := SymbolList.kv[name]; ok {
		return &sym
	} else {
		var sym Symbol
		//sym.name = name
		SymbolList.kv[name] = &sym
		return &sym
	}
}

func main() {
	//initializing
	Vtable_vt = Vtable_delegated(nil)
	
	Object_vt = Vtable_delegated(nil)
	Vtable_vt.parent = Object_vt

	Symbol_vt = Vtable_delegated(Object_vt)
	Closure_vt = Vtable_delegated(Object_vt) //

	SymbolList = Vtable_delegated(nil)

	S_lookup = Symbol_intern("lookup")
	S_addMethod = Symbol_intern("addMethod")
	S_delegated = Symbol_intern("delegated")
	
	Vtable_addMethod(Vtable_vt, &S_lookup, Vtable_lookup)
	Vtable_addMethod(Vtable_vt, &S_addMethod, Vtable_addMethod)
	Vtable_addMethod(Vtable_vt, &S_delegated, Vtable_delegated)

	//method, err := Send(Vtable_vt, &S_addMethod, Vtable_vt, &S_delegated, Vtable_delegated)

	//fmt.Println(err)
	
	//fmt.Println(Send(Vtable_vt, &S_, Vtable_vt, &S_delegated))
	//method.(func(*Vtable)(*Vtable))(Vtable_vt)

	//foo1, typ, err := Send(Vtable_vt, &S_lookup, Vtable_vt, &S_delegated)
	//foo1, err := Send(Vtable_vt, &S_delegated, Vtable_vt)

	//fmt.Println(foo1.(Vtable))

	// buff := new(bytes.Buffer)
	
	// binary.Write(buff, binary.LittleEndian, num2)

	var vbuff Vtable2
	var num1 uint16 = 130
	var num2 uint16 = 132
	vbuff.buffer = new(bytes.Buffer)
	binary.Write(vbuff.buffer, binary.LittleEndian, num1)
	fmt.Println(vbuff.buffer.Next(2))
	binary.Write(vbuff.buffer, binary.LittleEndian, num2)
	fmt.Println(vbuff.buffer.Bytes())
	
	buf := new(bytes.Buffer)
	var num uint16 = 1234
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes())

	// var S_foo interface{} = Symbol_intern("foo")
	// var S_bar interface{} = Symbol_intern("bar")
	// Vtable_addMethod(Vtable_vt, &S_foo, foobar)
	// Vtable_addMethod(Vtable_vt, &S_bar, barfoo)
	
	// foo, err1 := Vtable_lookup(Vtable_vt, &S_foo)
	// bar, err2 := Vtable_lookup(Vtable_vt, &S_bar)

	// Important
	//foo2 := foo1.([]reflect.Value)[0]
	//foo2 := reflect.ValueOf(reflect.Indirect(foo1.([]reflect.Value)[0]))
	// fmt.Println(reflect.Indirect(foo2).FieldByName("kv"))
	// fmt.Println(reflect.TypeOf(foo2))

	//foo, err1 := Vtable_lookup(foo2, &S_foo)

	//fmt.Println(reflect.TypeOf(reflect.ValueOf(foo1)).NumField())
	//fmt.Println(reflect.ValueOf(foo1).Index(0).(Vtable))
	//fmt.Println(reflect.ValueOf(foo1).Index(0).)

	//fmt.Println(err)
	

	// testing
	//fmt.Println(&S_lookup)
	//fmt.Println(&S_addMethod)
	//fmt.Println("...")
	//fmt.Println(Send(Vtable_vt, &S_lookup, Vtable_vt, &S_delegated))


	
	// fmt.Println(Send(Vtable_vt, &S_delegated, Vtable_vt))
	
	// fmt.Printf("%v\n", SymbolList)
	// fmt.Printf("%v\n", Vtable_vt)

	// var S_foo interface{} = Symbol_intern("foo")
	// var S_bar interface{} = Symbol_intern("bar")
	// Vtable_addMethod(Vtable_vt, &S_foo, foobar)
	// Vtable_addMethod(Vtable_vt, &S_bar, barfoo)
	
	// foo, err1 := Vtable_lookup(Vtable_vt, &S_foo)
	// bar, err2 := Vtable_lookup(Vtable_vt, &S_bar)

	// if err1 == nil {
	// 	foo.(func())()
	// }
	// if err2 == nil {
	// 	bar.(func(string))(" cat")
	// }

	// // testing inheritance
	// var Inh_vt = Vtable_delegated(Vtable_vt)
	// //var S_inh interface{} = Symbol_intern("inheritance")
	// inh, err3 := Vtable_lookup(Inh_vt, &S_bar)
	// if err3 == nil {
	// 	inh.(func(string))(" cat")
	// } else {
	// 	fmt.Println(err3)
	// }


	// //func Send(receiver *Vtable, message interface{}, arguments ...interface{}) {

	// Send(Vtable_vt, &S_bar, " little cat")
	// Send(Vtable_vt, &S_foo)



	
	//fmt.Println(reflect.TypeOf(1))
	

	//fmt.Println(reflect.TypeOf(reflect.TypeOf(inh).In(0)))
	//fmt.Println(reflect.TypeOf(inh).Out(0))


	//fmt.Println(reflect.TypeOf(inh).Kind())			
	
	//fmt.Println(Vtable_lookup(Vtable_vt, S_foo))
}
