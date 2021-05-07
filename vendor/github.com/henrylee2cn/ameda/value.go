package ameda

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

// Value go underlying type data
type Value struct {
	typPtr uintptr
	kind   reflect.Kind
	ptr    unsafe.Pointer
	_iPtr  unsafe.Pointer // avoid being GC
}

// ValueOf unpacks i to go underlying type data.
func ValueOf(i interface{}) Value {
	checkValueUsable()
	return newT(unsafe.Pointer(&i))
}

// ValueFrom gets go underlying type data from reflect.Value.
func ValueFrom(v reflect.Value) Value {
	checkValueUsable()
	return newT(unsafe.Pointer(&v))
}

func newT(iPtr unsafe.Pointer) Value {
	typPtr := *(*uintptr)(iPtr)
	return Value{
		typPtr: typPtr,
		kind:   kind(typPtr),
		ptr:    pointerElem(unsafe.Pointer(uintptr(iPtr) + ptrOffset)),
		_iPtr:  iPtr,
	}
}

// RuntimeTypeID returns the underlying type ID in current runtime from reflect.Type.
// NOTE:
//  *A and A returns the same runtime type ID;
//  It is 10 times performance of t.String().
func RuntimeTypeID(t reflect.Type) int32 {
	checkValueUsable()
	typPtr := uintptrElem(uintptr(unsafe.Pointer(&t)) + ptrOffset)
	return *(*int32)(unsafe.Pointer(typPtr + rtypeStrOffset))
}

// RuntimeTypeID gets the underlying type ID in current runtime.
// NOTE:
//  *A and A gets the same runtime type ID;
//  It is 10 times performance of reflect.TypeOf(i).String().
func (v Value) RuntimeTypeID() int32 {
	return *(*int32)(unsafe.Pointer(v.typPtr + rtypeStrOffset))
}

// Kind gets the reflect.Kind fastly.
func (v Value) Kind() reflect.Kind {
	return v.kind
}

// Elem returns the Value that the interface i contains
// or that the pointer i points to.
func (v Value) Elem() Value {
	k := v.kind
	switch k {
	default:
		return v
	case reflect.Interface:
		return newT(v.ptr)
	case reflect.Ptr:
		var has bool
		v.kind, v.typPtr, has = typeUnderlying(k, v.typPtr)
		if !has {
			return v
		}
		if v.kind == reflect.Ptr {
			v.ptr = pointerElem(v.ptr)
		}
		return v
	}
}

// UnderlyingElem returns the underlying Value that the interface i contains
// or that the pointer i points to.
func (v Value) UnderlyingElem() Value {
	for v.kind == reflect.Ptr || v.kind == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// Pointer gets the pointer of i.
// NOTE:
//  *T and T, gets diffrent pointer
func (v Value) Pointer() uintptr {
	switch v.Kind() {
	case reflect.Invalid:
		return 0
	case reflect.Slice:
		return uintptrElem(uintptr(v.ptr)) + sliceDataOffset
	default:
		return uintptr(v.ptr)
	}
}

// IsNil reports whether its argument i is nil.
func (v Value) IsNil() bool {
	return unsafe.Pointer(v.Pointer()) == nil
}

// FuncForPC returns a *Func describing the function that contains the
// given program counter address, or else nil.
//
// If pc represents multiple functions because of inlining, it returns
// the a *Func describing the innermost function, but with an entry
// of the outermost function.
//
// NOTE: Its kind must be a reflect.Func, otherwise it returns nil
func (v Value) FuncForPC() *runtime.Func {
	return runtime.FuncForPC(*(*uintptr)(v.ptr))
}

func typeUnderlying(k reflect.Kind, typPtr uintptr) (reflect.Kind, uintptr, bool) {
	typPtr2 := uintptrElem(typPtr + elemOffset)
	k2 := kind(typPtr2)
	if k2 == reflect.Invalid {
		return k, typPtr, false
	}
	return k2, typPtr2, true
}

func kind(typPtr uintptr) reflect.Kind {
	if unsafe.Pointer(typPtr) == nil {
		return reflect.Invalid
	}
	k := *(*uint8)(unsafe.Pointer(typPtr + kindOffset))
	return reflect.Kind(k & kindMask)
}

func uintptrElem(ptr uintptr) uintptr {
	return *(*uintptr)(unsafe.Pointer(ptr))
}

func pointerElem(p unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(p)
}

var errValueUsable error

func init() {
	goVersion := strings.TrimPrefix(runtime.Version(), "go")
	a, err := StringsToInts(strings.Split(goVersion, "."))
	if err != nil {
		errValueUsable = err
		return
	}
	if a[0] != 1 || a[1] < 9 {
		errValueUsable = fmt.Errorf("required go>=1.9, but current version is go" + goVersion)
	}
}

func checkValueUsable() {
	if errValueUsable != nil {
		panic(errValueUsable)
	}
}

var (
	e         = emptyInterface{typ: new(rtype)}
	ptrOffset = func() uintptr {
		return unsafe.Offsetof(e.word)
	}()
	rtypeStrOffset = func() uintptr {
		return unsafe.Offsetof(e.typ.str)
	}()
	kindOffset = func() uintptr {
		return unsafe.Offsetof(e.typ.kind)
	}()
	elemOffset = func() uintptr {
		return unsafe.Offsetof(new(ptrType).elem)
	}()
	sliceLenOffset = func() uintptr {
		return unsafe.Offsetof(new(reflect.SliceHeader).Len)
	}()
	sliceDataOffset = func() uintptr {
		return unsafe.Offsetof(new(reflect.SliceHeader).Data)
	}()
)

// NOTE: The following definitions must be consistent with those in the standard package!!!

const (
	kindMask = (1 << 5) - 1
)

type (
	// reflectValue struct {
	// 	typ *rtype
	// 	ptr unsafe.Pointer
	// 	flag
	// }
	emptyInterface struct {
		typ  *rtype
		word unsafe.Pointer
	}
	rtype struct {
		size       uintptr
		ptrdata    uintptr  // number of bytes in the type that can contain pointers
		hash       uint32   // hash of type; avoids computation in hash tables
		tflag      tflag    // extra type information flags
		align      uint8    // alignment of variable with this type
		fieldAlign uint8    // alignment of struct field with this type
		kind       uint8    // enumeration for C
		alg        *typeAlg // algorithm table
		gcdata     *byte    // garbage collection data
		str        nameOff  // string form
		ptrToThis  typeOff  // type for pointer to this type, may be zero
	}
	ptrType struct {
		rtype
		elem *rtype // pointer element (pointed at) type
	}
	typeAlg struct {
		hash  func(unsafe.Pointer, uintptr) uintptr
		equal func(unsafe.Pointer, unsafe.Pointer) bool
	}
	nameOff int32 // offset to a name
	typeOff int32 // offset to an *rtype
	flag    uintptr
	tflag   uint8
)
