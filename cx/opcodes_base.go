// +build base extra full

package base

import (
	
)

// op codes
const (
	// time
	OP_TIME_SLEEP = iota + END_OF_BARE_OPS
	OP_TIME_UNIX
	OP_TIME_UNIX_MILLI
	OP_TIME_UNIX_NANO

	// serialize
	OP_SERIAL_PROGRAM

	// os
	OP_OS_GET_WORKING_DIRECTORY
	OP_OS_OPEN
	OP_OS_CLOSE

	// http
	OP_HTTP_GET

	// object explorer
	OP_OBJ_QUERY

	END_OF_BASE_OPS
)

func init () {
	// time
	AddOpCode(OP_TIME_SLEEP, "time.Sleep", []int{TYPE_I32}, []int{})
	AddOpCode(OP_TIME_UNIX_MILLI, "time.UnixMilli", []int{}, []int{TYPE_I64})
	AddOpCode(OP_TIME_UNIX_NANO, "time.UnixNano", []int{}, []int{TYPE_I64})

	// http
	AddOpCode(OP_HTTP_GET, "http.Get", []int{TYPE_STR}, []int{TYPE_STR})

	// os
	AddOpCode(OP_OS_GET_WORKING_DIRECTORY, "os.GetWorkingDirectory", []int{}, []int{TYPE_STR})
	AddOpCode(OP_OS_OPEN, "os.Open", []int{TYPE_STR}, []int{})
	AddOpCode(OP_OS_CLOSE, "os.Close", []int{TYPE_STR}, []int{})

	// exec
	execNative = func(prgrm *CXProgram) {
		call := &prgrm.CallStack[prgrm.CallCounter]
		expr := call.Operator.Expressions[call.Line]
		opCode := expr.Operator.OpCode
		fp := call.FramePointer

		switch opCode {
		case OP_IDENTITY:
			op_identity(expr, fp)
		case OP_JMP:
			op_jmp(expr, fp, call)
		case OP_DEBUG:
			prgrm.PrintStack()

		case OP_UND_EQUAL:
			op_equal(expr, fp)
		case OP_UND_UNEQUAL:
			op_unequal(expr, fp)
		case OP_UND_BITAND:
			op_bitand(expr, fp)
		case OP_UND_BITXOR:
			op_bitxor(expr, fp)
		case OP_UND_BITOR:
			op_bitor(expr, fp)
		case OP_UND_BITCLEAR:
			op_bitclear(expr, fp)
		case OP_UND_MUL:
			op_mul(expr, fp)
		case OP_UND_DIV:
			op_div(expr, fp)
		case OP_UND_MOD:
			op_mod(expr, fp)
		case OP_UND_ADD:
			op_add(expr, fp)
		case OP_UND_SUB:
			op_sub(expr, fp)
		case OP_UND_BITSHL:
			op_bitshl(expr, fp)
		case OP_UND_BITSHR:
			op_bitshr(expr, fp)
		case OP_UND_LT:
			op_lt(expr, fp)
		case OP_UND_GT:
			op_gt(expr, fp)
		case OP_UND_LTEQ:
			op_lteq(expr, fp)
		case OP_UND_GTEQ:
			op_gteq(expr, fp)
		case OP_UND_LEN:
			op_len(expr, fp)
		case OP_UND_PRINTF:
			op_printf(expr, fp)
		case OP_UND_SPRINTF:
			op_sprintf(expr, fp)
		case OP_UND_READ:
			op_read(expr, fp)

		case OP_BYTE_BYTE:
			op_byte_byte(expr, fp)
		case OP_BYTE_STR:
			op_byte_byte(expr, fp)
		case OP_BYTE_I32:
			op_byte_byte(expr, fp)
		case OP_BYTE_I64:
			op_byte_byte(expr, fp)
		case OP_BYTE_F32:
			op_byte_byte(expr, fp)
		case OP_BYTE_F64:
			op_byte_byte(expr, fp)

		case OP_BYTE_PRINT:
			op_byte_print(expr, fp)

		case OP_BOOL_PRINT:
			op_bool_print(expr, fp)
		case OP_BOOL_EQUAL:
			op_bool_equal(expr, fp)
		case OP_BOOL_UNEQUAL:
			op_bool_unequal(expr, fp)
		case OP_BOOL_NOT:
			op_bool_not(expr, fp)
		case OP_BOOL_OR:
			op_bool_or(expr, fp)
		case OP_BOOL_AND:
			op_bool_and(expr, fp)

		case OP_I32_BYTE:
			op_i32_i32(expr, fp)
		case OP_I32_STR:
			op_i32_i32(expr, fp)
		case OP_I32_I32:
			op_i32_i32(expr, fp)
		case OP_I32_I64:
			op_i32_i32(expr, fp)
		case OP_I32_F32:
			op_i32_i32(expr, fp)
		case OP_I32_F64:
			op_i32_i32(expr, fp)
			
		case OP_I32_PRINT:
			op_i32_print(expr, fp)
		case OP_I32_ADD:
			op_i32_add(expr, fp)
		case OP_I32_SUB:
			op_i32_sub(expr, fp)
		case OP_I32_MUL:
			op_i32_mul(expr, fp)
		case OP_I32_DIV:
			op_i32_div(expr, fp)
		case OP_I32_ABS:
			op_i32_abs(expr, fp)
		case OP_I32_POW:
			op_i32_pow(expr, fp)
		case OP_I32_GT:
			op_i32_gt(expr, fp)
		case OP_I32_GTEQ:
			op_i32_gteq(expr, fp)
		case OP_I32_LT:
			op_i32_lt(expr, fp)
		case OP_I32_LTEQ:
			op_i32_lteq(expr, fp)
		case OP_I32_EQ:
			op_i32_eq(expr, fp)
		case OP_I32_UNEQ:
			op_i32_uneq(expr, fp)
		case OP_I32_MOD:
			op_i32_mod(expr, fp)
		case OP_I32_RAND:
			op_i32_rand(expr, fp)
		case OP_I32_BITAND:
			op_i32_bitand(expr, fp)
		case OP_I32_BITOR:
			op_i32_bitor(expr, fp)
		case OP_I32_BITXOR:
			op_i32_bitxor(expr, fp)
		case OP_I32_BITCLEAR:
			op_i32_bitclear(expr, fp)
		case OP_I32_BITSHL:
			op_i32_bitshl(expr, fp)
		case OP_I32_BITSHR:
			op_i32_bitshr(expr, fp)
		case OP_I32_SQRT:
			op_i32_sqrt(expr, fp)
		case OP_I32_LOG:
			op_i32_log(expr, fp)
		case OP_I32_LOG2:
			op_i32_log2(expr, fp)
		case OP_I32_LOG10:
			op_i32_log10(expr, fp)

		case OP_I32_MAX:
			op_i32_max(expr, fp)
		case OP_I32_MIN:
			op_i32_min(expr, fp)

		case OP_I64_BYTE:
			op_i64_i64(expr, fp)
		case OP_I64_STR:
			op_i64_i64(expr, fp)
		case OP_I64_I32:
			op_i64_i64(expr, fp)
		case OP_I64_I64:
			op_i64_i64(expr, fp)
		case OP_I64_F32:
			op_i64_i64(expr, fp)
		case OP_I64_F64:
			op_i64_i64(expr, fp)

		case OP_I64_PRINT:
			op_i64_print(expr, fp)
		case OP_I64_ADD:
			op_i64_add(expr, fp)
		case OP_I64_SUB:
			op_i64_sub(expr, fp)
		case OP_I64_MUL:
			op_i64_mul(expr, fp)
		case OP_I64_DIV:
			op_i64_div(expr, fp)
		case OP_I64_ABS:
			op_i64_abs(expr, fp)
		case OP_I64_POW:
			op_i64_pow(expr, fp)
		case OP_I64_GT:
			op_i64_gt(expr, fp)
		case OP_I64_GTEQ:
			op_i64_gteq(expr, fp)
		case OP_I64_LT:
			op_i64_lt(expr, fp)
		case OP_I64_LTEQ:
			op_i64_lteq(expr, fp)
		case OP_I64_EQ:
			op_i64_eq(expr, fp)
		case OP_I64_UNEQ:
			op_i64_uneq(expr, fp)
		case OP_I64_MOD:
			op_i64_mod(expr, fp)
		case OP_I64_RAND:
			op_i64_rand(expr, fp)
		case OP_I64_BITAND:
			op_i64_bitand(expr, fp)
		case OP_I64_BITOR:
			op_i64_bitor(expr, fp)
		case OP_I64_BITXOR:
			op_i64_bitxor(expr, fp)
		case OP_I64_BITCLEAR:
			op_i64_bitclear(expr, fp)
		case OP_I64_BITSHL:
			op_i64_bitshl(expr, fp)
		case OP_I64_BITSHR:
			op_i64_bitshr(expr, fp)
		case OP_I64_SQRT:
			op_i64_sqrt(expr, fp)
		case OP_I64_LOG:
			op_i64_log(expr, fp)
		case OP_I64_LOG2:
			op_i64_log2(expr, fp)
		case OP_I64_LOG10:
			op_i64_log10(expr, fp)
		case OP_I64_MAX:
			op_i64_max(expr, fp)
		case OP_I64_MIN:
			op_i64_min(expr, fp)

		case OP_F32_BYTE:
			op_f32_f32(expr, fp)
		case OP_F32_STR:
			op_f32_f32(expr, fp)
		case OP_F32_I32:
			op_f32_f32(expr, fp)
		case OP_F32_I64:
			op_f32_f32(expr, fp)
		case OP_F32_F32:
			op_f32_f32(expr, fp)
		case OP_F32_F64:
			op_f32_f32(expr, fp)
			
		case OP_F32_PRINT:
			op_f32_print(expr, fp)
		case OP_F32_ADD:
			op_f32_add(expr, fp)
		case OP_F32_SUB:
			op_f32_sub(expr, fp)
		case OP_F32_MUL:
			op_f32_mul(expr, fp)
		case OP_F32_DIV:
			op_f32_div(expr, fp)
		case OP_F32_ABS:
			op_f32_abs(expr, fp)
		case OP_F32_POW:
			op_f32_pow(expr, fp)
		case OP_F32_GT:
			op_f32_gt(expr, fp)
		case OP_F32_GTEQ:
			op_f32_gteq(expr, fp)
		case OP_F32_LT:
			op_f32_lt(expr, fp)
		case OP_F32_LTEQ:
			op_f32_lteq(expr, fp)
		case OP_F32_EQ:
			op_f32_eq(expr, fp)
		case OP_F32_UNEQ:
			op_f32_uneq(expr, fp)
		case OP_F32_COS:
			op_f32_cos(expr, fp)
		case OP_F32_SIN:
			op_f32_sin(expr, fp)
		case OP_F32_SQRT:
			op_f32_sqrt(expr, fp)
		case OP_F32_LOG:
			op_f32_log(expr, fp)
		case OP_F32_LOG2:
			op_f32_log2(expr, fp)
		case OP_F32_LOG10:
			op_f32_log10(expr, fp)
		case OP_F32_MAX:
			op_f32_max(expr, fp)
		case OP_F32_MIN:
			op_f32_min(expr, fp)

		case OP_F64_BYTE:
			op_f64_f64(expr, fp)
		case OP_F64_STR:
			op_f64_f64(expr, fp)
		case OP_F64_I32:
			op_f64_f64(expr, fp)
		case OP_F64_I64:
			op_f64_f64(expr, fp)
		case OP_F64_F32:
			op_f64_f64(expr, fp)
		case OP_F64_F64:
			op_f64_f64(expr, fp)

		case OP_F64_PRINT:
			op_f64_print(expr, fp)
		case OP_F64_ADD:
			op_f64_add(expr, fp)
		case OP_F64_SUB:
			op_f64_sub(expr, fp)
		case OP_F64_MUL:
			op_f64_mul(expr, fp)
		case OP_F64_DIV:
			op_f64_div(expr, fp)
		case OP_F64_ABS:
			op_f64_abs(expr, fp)
		case OP_F64_POW:
			op_f64_pow(expr, fp)
		case OP_F64_GT:
			op_f64_gt(expr, fp)
		case OP_F64_GTEQ:
			op_f64_gteq(expr, fp)
		case OP_F64_LT:
			op_f64_lt(expr, fp)
		case OP_F64_LTEQ:
			op_f64_lteq(expr, fp)
		case OP_F64_EQ:
			op_f64_eq(expr, fp)
		case OP_F64_UNEQ:
			op_f64_uneq(expr, fp)
		case OP_F64_COS:
			op_f64_cos(expr, fp)
		case OP_F64_SIN:
			op_f64_sin(expr, fp)
		case OP_F64_SQRT:
			op_f64_sqrt(expr, fp)
		case OP_F64_LOG:
			op_f64_log(expr, fp)
		case OP_F64_LOG2:
			op_f64_log2(expr, fp)
		case OP_F64_LOG10:
			op_f64_log10(expr, fp)
		case OP_F64_MAX:
			op_f64_max(expr, fp)
		case OP_F64_MIN:
			op_f64_min(expr, fp)
		case OP_STR_PRINT:
			op_str_print(expr, fp)
		case OP_STR_CONCAT:
			op_str_concat(expr, fp)
		case OP_STR_EQ:
			op_str_eq(expr, fp)
			
		case OP_STR_BYTE:
			op_str_str(expr, fp)
		case OP_STR_STR:
			op_str_str(expr, fp)
		case OP_STR_I32:
			op_str_str(expr, fp)
		case OP_STR_I64:
			op_str_str(expr, fp)
		case OP_STR_F32:
			op_str_str(expr, fp)
		case OP_STR_F64:
			op_str_str(expr, fp)

		case OP_MAKE:
		case OP_READ:
		case OP_WRITE:
		case OP_LEN:
		case OP_CONCAT:
		case OP_APPEND:
			op_append(expr, fp)
		case OP_COPY:
		case OP_CAST:
		case OP_EQ:
		case OP_UNEQ:
		case OP_AND:
		case OP_OR:
		case OP_NOT:
		case OP_SLEEP:
		case OP_HALT:
		case OP_GOTO:
		case OP_REMCX:
		case OP_ADDCX:
		case OP_QUERY:
		case OP_EXECUTE:
		case OP_INDEX:
		case OP_NAME:
		case OP_EVOLVE:
		case OP_ASSERT:
			op_assert_value(expr, fp)

			// time
		case OP_TIME_SLEEP:
			op_time_Sleep(expr, fp)
		case OP_TIME_UNIX:
		case OP_TIME_UNIX_MILLI:
			op_time_UnixMilli(expr, fp)
		case OP_TIME_UNIX_NANO:
			op_time_UnixNano(expr, fp)

			// affordances
		case OP_AFF_PRINT:
			op_aff_print(expr, fp)
		case OP_AFF_QUERY:
			op_aff_query(expr, fp)
		case OP_AFF_ON:
			op_aff_on(expr, fp)
		case OP_AFF_OF:
			op_aff_of(expr, fp)
		case OP_AFF_INFORM:
			op_aff_inform(expr, fp)
		case OP_AFF_REQUEST:
			op_aff_request(expr, fp)

			// http
		case OP_HTTP_GET:
			op_http_get(expr, fp)

			// os
		case OP_OS_GET_WORKING_DIRECTORY:
			op_os_GetWorkingDirectory(expr, fp)
		case OP_OS_OPEN:
			op_os_Open(expr, fp)
		case OP_OS_CLOSE:
			op_os_Close(expr, fp)
		}
	}
}
