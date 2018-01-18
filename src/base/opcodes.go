package base

// op codes
const (
	OP_IDENTITY = iota
	OP_I32_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_ABS
	OP_MOD
	OP_POW
	OP_COS
	OP_SIN
	OP_BITAND
	OP_BITOR
	OP_BITXOR
	OP_BITCLEAR
	OP_BITSHL
	OP_BITSHR
	OP_I32_PRINT
	OP_MAKE
	OP_READ
	OP_WRITE
	OP_LEN
	OP_CONCAT
	OP_APPEND
	OP_COPY
	OP_CAST
	OP_EQ
	OP_UNEQ
	OP_LT
	OP_GT
	OP_LTEQ
	OP_GTEQ
	OP_RAND
	OP_AND
	OP_OR
	OP_NOT
	OP_SLEEP
	OP_HALT
	OP_GOTO
	OP_REMCX
	OP_ADDCX
	OP_QUERY
	OP_EXECUTE
	OP_INDEX
	OP_NAME
	OP_EVOLVE
	OP_TEST_START
	OP_TEST_STOP
	OP_TEST_ERROR
	OP_TEST
	OP_TIME_UNIX
	OP_TIME_UNIXMILLI
	OP_TIME_UNIXNANO
)


// For the parser. These shouldn't be used in the runtime for performance reasons
var OpNames map[int]string = map[int]string{
	OP_IDENTITY: "identity",
	OP_I32_ADD: "i32.add",
}

// For the parser. These shouldn't be used in the runtime for performance reasons
var OpCodes map[string]int = map[string]int{
	"identity": OP_IDENTITY,
	"i32.add": OP_I32_ADD,
}

// // inputs, then outputs
// var OpSignature map[int][2][]int = map[int][2][]int{
// 	OP_IDENTITY: [2][]int{[]int{TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}},
// 	OP_I32_ADD: [2][]int{[]int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}},
// }

var Natives map[int]*CXFunction = map[int]*CXFunction{
	OP_IDENTITY: MakeNative(OP_IDENTITY, []int{TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_I32_ADD: MakeNative(OP_I32_ADD, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
}
