package base

// op codes
const (
	OP_IDENTITY = iota
	OP_READ_ARRAY
	OP_JMP
	OP_DEBUG

	OP_UND_EQUAL
	OP_UND_UNEQUAL
	OP_UND_BITAND
	OP_UND_BITXOR
	OP_UND_BITOR
	OP_UND_BITCLEAR
	OP_UND_MUL
	OP_UND_DIV
	OP_UND_MOD
	OP_UND_ADD
	OP_UND_SUB
	OP_UND_BITSHL
	OP_UND_BITSHR
	OP_UND_LT
	OP_UND_GT
	OP_UND_LTEQ
	OP_UND_GTEQ
	OP_UND_LEN
	OP_UND_PRINTF
	OP_UND_SPRINTF

	OP_BOOL_PRINT

	OP_BOOL_EQUAL
	OP_BOOL_UNEQUAL
	OP_BOOL_NOT
	OP_BOOL_OR
	OP_BOOL_AND

	OP_BYTE_PRINT

	OP_I32_BYTE
	OP_I32_STR
	OP_I32_I32
	OP_I32_I64
	OP_I32_F32
	OP_I32_F64
	OP_I32_PRINT
	OP_I32_ADD
	OP_I32_SUB
	OP_I32_MUL
	OP_I32_DIV
	OP_I32_ABS
	OP_I32_POW
	OP_I32_GT
	OP_I32_GTEQ
	OP_I32_LT
	OP_I32_LTEQ
	OP_I32_EQ
	OP_I32_UNEQ
	OP_I32_MOD
	OP_I32_RAND
	OP_I32_BITAND
	OP_I32_BITOR
	OP_I32_BITXOR
	OP_I32_BITCLEAR
	OP_I32_BITSHL
	OP_I32_BITSHR
	OP_I32_SQRT
	OP_I32_LOG
	OP_I32_LOG2
	OP_I32_LOG10

	OP_I32_MAX
	OP_I32_MIN
	OP_I32_SIN
	OP_I32_COS

	OP_I64_PRINT
	OP_I64_ADD
	OP_I64_SUB
	OP_I64_MUL
	OP_I64_DIV
	OP_I64_ABS
	OP_I64_POW
	OP_I64_GT
	OP_I64_GTEQ
	OP_I64_LT
	OP_I64_LTEQ
	OP_I64_EQ
	OP_I64_UNEQ
	OP_I64_MOD
	OP_I64_RAND
	OP_I64_BITAND
	OP_I64_BITOR
	OP_I64_BITXOR
	OP_I64_BITCLEAR
	OP_I64_BITSHL
	OP_I64_BITSHR
	OP_I64_SQRT
	OP_I64_LOG
	OP_I64_LOG10
	OP_I64_LOG2
	OP_I64_MAX
	OP_I64_MIN
	OP_I64_SIN
	OP_I64_COS

	OP_F32_PRINT
	OP_F32_ADD
	OP_F32_SUB
	OP_F32_MUL
	OP_F32_DIV
	OP_F32_ABS
	OP_F32_POW
	OP_F32_GT
	OP_F32_GTEQ
	OP_F32_LT
	OP_F32_LTEQ
	OP_F32_EQ
	OP_F32_UNEQ
	OP_F32_COS
	OP_F32_SIN
	OP_F32_SQRT
	OP_F32_LOG
	OP_F32_LOG2
	OP_F32_LOG10
	OP_F32_MAX
	OP_F32_MIN

	OP_F64_PRINT
	OP_F64_ADD
	OP_F64_SUB
	OP_F64_MUL
	OP_F64_DIV
	OP_F64_ABS
	OP_F64_POW
	OP_F64_GT
	OP_F64_GTEQ
	OP_F64_LT
	OP_F64_LTEQ
	OP_F64_EQ
	OP_F64_UNEQ
	OP_F64_COS
	OP_F64_SIN

	OP_F64_SQRT
	OP_F64_LOG
	OP_F64_LOG2
	OP_F64_LOG10
	OP_F64_MAX
	OP_F64_MIN

	OP_STR_PRINT
	OP_STR_EQ

	
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

	OP_TIME_SLEEP
	OP_TIME_UNIX
	OP_TIME_UNIX_MILLI
	OP_TIME_UNIX_NANO

	// opengl
	OP_GL_INIT
	OP_GL_CREATE_PROGRAM
	OP_GL_LINK_PROGRAM
	OP_GL_CLEAR
	OP_GL_USE_PROGRAM
	OP_GL_BIND_BUFFER
	OP_GL_BIND_VERTEX_ARRAY
	OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY
	OP_GL_VERTEX_ATTRIB_POINTER
	OP_GL_DRAW_ARRAYS
	OP_GL_GEN_BUFFERS
	OP_GL_BUFFER_DATA
	OP_GL_GEN_VERTEX_ARRAYS
	OP_GL_CREATE_SHADER
	OP_GL_STRS
	OP_GL_FREE
	OP_GL_SHADER_SOURCE
	OP_GL_COMPILE_SHADER
	OP_GL_GET_SHADERIV
	OP_GL_ATTACH_SHADER
	OP_GL_MATRIX_MODE
	OP_GL_ROTATEF
	OP_GL_TRANSLATEF
	OP_GL_LOAD_IDENTITY
	OP_GL_PUSH_MATRIX
	OP_GL_POP_MATRIX
	OP_GL_ENABLE_CLIENT_STATE
	OP_GL_BIND_TEXTURE
	OP_GL_COLOR3F
	OP_GL_COLOR4F
	OP_GL_BEGIN
	OP_GL_END
	OP_GL_NORMAL3F
	OP_GL_VERTEX_2F
	OP_GL_VERTEX_3F
	OP_GL_ENABLE
	OP_GL_CLEAR_COLOR
	OP_GL_CLEAR_DEPTH
	OP_GL_DEPTH_FUNC
	OP_GL_LIGHTFV
	OP_GL_FRUSTUM
	OP_GL_DISABLE
	OP_GL_HINT
	OP_GL_NEW_TEXTURE
	OP_GL_DEPTH_MASK
	OP_GL_TEX_ENVI
	OP_GL_BLEND_FUNC
	OP_GL_ORTHO
	OP_GL_VIEWPORT
	OP_GL_SCALEF
	OP_GL_TEX_COORD_2D

	// glfw
	OP_GLFW_INIT
	OP_GLFW_WINDOW_HINT
	OP_GLFW_CREATE_WINDOW
	OP_GLFW_MAKE_CONTEXT_CURRENT
	OP_GLFW_SHOULD_CLOSE
	OP_GLFW_SET_SHOULD_CLOSE
	OP_GLFW_POLL_EVENTS
	OP_GLFW_SWAP_BUFFERS
	OP_GLFW_GET_FRAMEBUFFER_SIZE
	OP_GLFW_SET_KEY_CALLBACK
	OP_GLFW_GET_TIME
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK
	OP_GLFW_SET_CURSOR_POS_CALLBACK
	OP_GLFW_GET_CURSOR_POS
	OP_GLFW_SET_INPUT_MODE

	// http
	OP_HTTP_GET
)

func execNative(prgrm *CXProgram) {
	call := &prgrm.CallStack[prgrm.CallCounter]
	stack := &prgrm.Stacks[0]
	expr := call.Operator.Expressions[call.Line]
	opCode := expr.Operator.OpCode
	fp := call.FramePointer

	switch opCode {
	case OP_IDENTITY:
		op_identity(expr, stack, fp)
		// case OP_READ_ARRAY: op_read_array(expr, stack, fp)
	case OP_JMP:
		op_jmp(expr, stack, fp, call)
	case OP_DEBUG:
		prgrm.PrintStack()

	case OP_UND_EQUAL:
		op_equal(expr, stack, fp)
	case OP_UND_UNEQUAL:
		op_unequal(expr, stack, fp)
	case OP_UND_BITAND:
		op_bitand(expr, stack, fp)
	case OP_UND_BITXOR:
		op_bitxor(expr, stack, fp)
	case OP_UND_BITOR:
		op_bitor(expr, stack, fp)
	case OP_UND_BITCLEAR:
		op_bitclear(expr, stack, fp)
	case OP_UND_MUL:
		op_mul(expr, stack, fp)
	case OP_UND_DIV:
		op_div(expr, stack, fp)
	case OP_UND_MOD:
		op_mod(expr, stack, fp)
	case OP_UND_ADD:
		op_add(expr, stack, fp)
	case OP_UND_SUB:
		op_sub(expr, stack, fp)
	case OP_UND_BITSHL:
		op_bitshl(expr, stack, fp)
	case OP_UND_BITSHR:
		op_bitshr(expr, stack, fp)
	case OP_UND_LT:
		op_lt(expr, stack, fp)
	case OP_UND_GT:
		op_gt(expr, stack, fp)
	case OP_UND_LTEQ:
		op_lteq(expr, stack, fp)
	case OP_UND_GTEQ:
		op_gteq(expr, stack, fp)
	case OP_UND_LEN:
		op_len(expr, stack, fp)
	case OP_UND_PRINTF:
		op_printf(expr, stack, fp)
	case OP_UND_SPRINTF:
		op_sprintf(expr, stack, fp)

	case OP_BYTE_PRINT:
		op_byte_print(expr, stack, fp)

	case OP_BOOL_PRINT:
		op_bool_print(expr, stack, fp)
	case OP_BOOL_EQUAL:
		op_bool_equal(expr, stack, fp)
	case OP_BOOL_UNEQUAL:
		op_bool_unequal(expr, stack, fp)
	case OP_BOOL_NOT:
		op_bool_not(expr, stack, fp)
	case OP_BOOL_OR:
		op_bool_or(expr, stack, fp)
	case OP_BOOL_AND:
		op_bool_and(expr, stack, fp)

	case OP_I32_BYTE:
		op_i32_i32(expr, stack, fp)
	case OP_I32_STR:
		op_i32_i32(expr, stack, fp)
	case OP_I32_I32:
		op_i32_i32(expr, stack, fp)
	case OP_I32_I64:
		op_i32_i32(expr, stack, fp)
	case OP_I32_F32:
		op_i32_i32(expr, stack, fp)
	case OP_I32_F64:
		op_i32_i32(expr, stack, fp)
	case OP_I32_PRINT:
		op_i32_print(expr, stack, fp)
	case OP_I32_ADD:
		op_i32_add(expr, stack, fp)
	case OP_I32_SUB:
		op_i32_sub(expr, stack, fp)
	case OP_I32_MUL:
		op_i32_mul(expr, stack, fp)
	case OP_I32_DIV:
		op_i32_div(expr, stack, fp)
	case OP_I32_ABS:
		op_i32_abs(expr, stack, fp)
	case OP_I32_POW:
		op_i32_pow(expr, stack, fp)
	case OP_I32_GT:
		op_i32_gt(expr, stack, fp)
	case OP_I32_GTEQ:
		op_i32_gteq(expr, stack, fp)
	case OP_I32_LT:
		op_i32_lt(expr, stack, fp)
	case OP_I32_LTEQ:
		op_i32_lteq(expr, stack, fp)
	case OP_I32_EQ:
		op_i32_eq(expr, stack, fp)
	case OP_I32_UNEQ:
		op_i32_uneq(expr, stack, fp)
	case OP_I32_MOD:
		op_i32_mod(expr, stack, fp)
	case OP_I32_RAND:
		op_i32_rand(expr, stack, fp)
	case OP_I32_BITAND:
		op_i32_bitand(expr, stack, fp)
	case OP_I32_BITOR:
		op_i32_bitor(expr, stack, fp)
	case OP_I32_BITXOR:
		op_i32_bitxor(expr, stack, fp)
	case OP_I32_BITCLEAR:
		op_i32_bitclear(expr, stack, fp)
	case OP_I32_BITSHL:
		op_i32_bitshl(expr, stack, fp)
	case OP_I32_BITSHR:
		op_i32_bitshr(expr, stack, fp)
	case OP_I32_SQRT:
		op_i32_sqrt(expr, stack, fp)
	case OP_I32_LOG:
		op_i32_log(expr, stack, fp)
	case OP_I32_LOG2:
		op_i32_log2(expr, stack, fp)
	case OP_I32_LOG10:
		op_i32_log10(expr, stack, fp)

	case OP_I32_MAX:
		op_i32_max(expr, stack, fp)
	case OP_I32_MIN:
		op_i32_min(expr, stack, fp)
	case OP_I32_SIN:
		op_i32_sin(expr, stack, fp)
	case OP_I32_COS:
		op_i32_cos(expr, stack, fp)

	case OP_I64_PRINT:
		op_i64_print(expr, stack, fp)
	case OP_I64_ADD:
		op_i64_add(expr, stack, fp)
	case OP_I64_SUB:
		op_i64_sub(expr, stack, fp)
	case OP_I64_MUL:
		op_i64_mul(expr, stack, fp)
	case OP_I64_DIV:
		op_i64_div(expr, stack, fp)
	case OP_I64_ABS:
		op_i64_abs(expr, stack, fp)
	case OP_I64_POW:
		op_i64_pow(expr, stack, fp)
	case OP_I64_GT:
		op_i64_gt(expr, stack, fp)
	case OP_I64_GTEQ:
		op_i64_gteq(expr, stack, fp)
	case OP_I64_LT:
		op_i64_lt(expr, stack, fp)
	case OP_I64_LTEQ:
		op_i64_lteq(expr, stack, fp)
	case OP_I64_EQ:
		op_i64_eq(expr, stack, fp)
	case OP_I64_UNEQ:
		op_i64_uneq(expr, stack, fp)
	case OP_I64_MOD:
		op_i64_mod(expr, stack, fp)
	case OP_I64_RAND:
		op_i64_rand(expr, stack, fp)
	case OP_I64_BITAND:
		op_i64_bitand(expr, stack, fp)
	case OP_I64_BITOR:
		op_i64_bitor(expr, stack, fp)
	case OP_I64_BITXOR:
		op_i64_bitxor(expr, stack, fp)
	case OP_I64_BITCLEAR:
		op_i64_bitclear(expr, stack, fp)
	case OP_I64_BITSHL:
		op_i64_bitshl(expr, stack, fp)
	case OP_I64_BITSHR:
		op_i64_bitshr(expr, stack, fp)
	case OP_I64_SQRT:
		op_i64_sqrt(expr, stack, fp)
	case OP_I64_LOG:
		op_i64_log(expr, stack, fp)
	case OP_I64_LOG2:
		op_i64_log2(expr, stack, fp)
	case OP_I64_LOG10:
		op_i64_log10(expr, stack, fp)
	case OP_I64_MAX:
		op_i64_max(expr, stack, fp)
	case OP_I64_MIN:
		op_i64_min(expr, stack, fp)
	case OP_I64_SIN:
		op_i64_sin(expr, stack, fp)
	case OP_I64_COS:
		op_i64_cos(expr, stack, fp)

	case OP_F32_PRINT:
		op_f32_print(expr, stack, fp)
	case OP_F32_ADD:
		op_f32_add(expr, stack, fp)
	case OP_F32_SUB:
		op_f32_sub(expr, stack, fp)
	case OP_F32_MUL:
		op_f32_mul(expr, stack, fp)
	case OP_F32_DIV:
		op_f32_div(expr, stack, fp)
	case OP_F32_ABS:
		op_f32_abs(expr, stack, fp)
	case OP_F32_POW:
		op_f32_pow(expr, stack, fp)
	case OP_F32_GT:
		op_f32_gt(expr, stack, fp)
	case OP_F32_GTEQ:
		op_f32_gteq(expr, stack, fp)
	case OP_F32_LT:
		op_f32_lt(expr, stack, fp)
	case OP_F32_LTEQ:
		op_f32_lteq(expr, stack, fp)
	case OP_F32_EQ:
		op_f32_eq(expr, stack, fp)
	case OP_F32_UNEQ:
		op_f32_uneq(expr, stack, fp)
	case OP_F32_COS:
		op_f32_cos(expr, stack, fp)
	case OP_F32_SIN:
		op_f32_sin(expr, stack, fp)
	case OP_F32_SQRT:
		op_f32_sqrt(expr, stack, fp)
	case OP_F32_LOG:
		op_f32_log(expr, stack, fp)
	case OP_F32_LOG2:
		op_f32_log2(expr, stack, fp)
	case OP_F32_LOG10:
		op_f32_log10(expr, stack, fp)
	case OP_F32_MAX:
		op_f32_max(expr, stack, fp)
	case OP_F32_MIN:
		op_f32_min(expr, stack, fp)

	case OP_F64_PRINT:
		op_f64_print(expr, stack, fp)
	case OP_F64_ADD:
		op_f64_add(expr, stack, fp)
	case OP_F64_SUB:
		op_f64_sub(expr, stack, fp)
	case OP_F64_MUL:
		op_f64_mul(expr, stack, fp)
	case OP_F64_DIV:
		op_f64_div(expr, stack, fp)
	case OP_F64_ABS:
		op_f64_abs(expr, stack, fp)
	case OP_F64_POW:
		op_f64_pow(expr, stack, fp)
	case OP_F64_GT:
		op_f64_gt(expr, stack, fp)
	case OP_F64_GTEQ:
		op_f64_gteq(expr, stack, fp)
	case OP_F64_LT:
		op_f64_lt(expr, stack, fp)
	case OP_F64_LTEQ:
		op_f64_lteq(expr, stack, fp)
	case OP_F64_EQ:
		op_f64_eq(expr, stack, fp)
	case OP_F64_UNEQ:
		op_f64_uneq(expr, stack, fp)
	case OP_F64_COS:
		op_f64_cos(expr, stack, fp)
	case OP_F64_SIN:
		op_f64_sin(expr, stack, fp)
	case OP_F64_SQRT:
		op_f64_sqrt(expr, stack, fp)
	case OP_F64_LOG:
		op_f64_log(expr, stack, fp)
	case OP_F64_LOG2:
		op_f64_log2(expr, stack, fp)
	case OP_F64_LOG10:
		op_f64_log10(expr, stack, fp)
	case OP_F64_MAX:
		op_f64_max(expr, stack, fp)
	case OP_F64_MIN:
		op_f64_min(expr, stack, fp)
	case OP_STR_PRINT:
		op_str_print(expr, stack, fp)
	case OP_STR_EQ:
		op_str_eq(expr, stack, fp)
	case OP_MAKE:
	case OP_READ:
	case OP_WRITE:
	case OP_LEN:
	case OP_CONCAT:
	case OP_APPEND:
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
	case OP_TEST_START:
		isTesting = true
	case OP_TEST_STOP:
		isTesting = false
	case OP_TEST_ERROR:
		isErrorPresent = false
	case OP_TEST:
		op_test_value(expr, stack, fp)
	case OP_TIME_SLEEP:
		op_time_Sleep(expr, stack, fp)
	case OP_TIME_UNIX:
	case OP_TIME_UNIX_MILLI:
		op_time_UnixMilli(expr, stack, fp)
	case OP_TIME_UNIX_NANO:
		op_time_UnixNano(expr, stack, fp)

		// opengl
	case OP_GL_INIT:
		op_gl_Init()
	case OP_GL_CREATE_PROGRAM:
		op_gl_CreateProgram(expr, stack, fp)
	case OP_GL_LINK_PROGRAM:
		op_gl_LinkProgram(expr, stack, fp)
	case OP_GL_CLEAR:
		op_gl_Clear(expr, stack, fp)
	case OP_GL_USE_PROGRAM:
		op_gl_UseProgram(expr, stack, fp)
	case OP_GL_BIND_BUFFER:
		op_gl_BindBuffer(expr, stack, fp)
	case OP_GL_BIND_VERTEX_ARRAY:
		op_gl_BindVertexArray(expr, stack, fp)
	case OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY:
		op_gl_EnableVertexAttribArray(expr, stack, fp)
	case OP_GL_VERTEX_ATTRIB_POINTER:
		op_gl_VertexAttribPointer(expr, stack, fp)
	case OP_GL_DRAW_ARRAYS:
		op_gl_DrawArrays(expr, stack, fp)
	case OP_GL_GEN_BUFFERS:
		op_gl_GenBuffers(expr, stack, fp)
	case OP_GL_BUFFER_DATA:
		op_gl_BufferData(expr, stack, fp)
	case OP_GL_GEN_VERTEX_ARRAYS:
		op_gl_GenVertexArrays(expr, stack, fp)
	case OP_GL_CREATE_SHADER:
		op_gl_CreateShader(expr, stack, fp)
	case OP_GL_STRS:
		op_gl_Strs(expr, stack, fp)
	case OP_GL_FREE:
		op_gl_Free(expr, stack, fp)
	case OP_GL_SHADER_SOURCE:
		op_gl_ShaderSource(expr, stack, fp)
	case OP_GL_COMPILE_SHADER:
		op_gl_CompileShader(expr, stack, fp)
	case OP_GL_GET_SHADERIV:
		op_gl_GetShaderiv(expr, stack, fp)
	case OP_GL_ATTACH_SHADER:
		op_gl_AttachShader(expr, stack, fp)
	case OP_GL_MATRIX_MODE:
		op_gl_MatrixMode(expr, stack, fp)
	case OP_GL_ROTATEF:
		op_gl_Rotatef(expr, stack, fp)
	case OP_GL_TRANSLATEF:
		op_gl_Translatef(expr, stack, fp)
	case OP_GL_LOAD_IDENTITY:
		op_gl_LoadIdentity()
	case OP_GL_PUSH_MATRIX:
		op_gl_PushMatrix()
	case OP_GL_POP_MATRIX:
		op_gl_PopMatrix()
	case OP_GL_ENABLE_CLIENT_STATE:
		op_gl_EnableClientState(expr, stack, fp)
	case OP_GL_BIND_TEXTURE:
		op_gl_BindTexture(expr, stack, fp)
	case OP_GL_COLOR3F:
		op_gl_Color3f(expr, stack, fp)
	case OP_GL_COLOR4F:
		op_gl_Color4f(expr, stack, fp)
	case OP_GL_BEGIN:
		op_gl_Begin(expr, stack, fp)
	case OP_GL_END:
		op_gl_End()
	case OP_GL_NORMAL3F:
		op_gl_Normal3f(expr, stack, fp)
	case OP_GL_VERTEX_2F:
		op_gl_Vertex2f(expr, stack, fp)
	case OP_GL_VERTEX_3F:
		op_gl_Vertex3f(expr, stack, fp)
	case OP_GL_ENABLE:
		op_gl_Enable(expr, stack, fp)
	case OP_GL_CLEAR_COLOR:
		op_gl_ClearColor(expr, stack, fp)
	case OP_GL_CLEAR_DEPTH:
		op_gl_ClearDepth(expr, stack, fp)
	case OP_GL_DEPTH_FUNC:
		op_gl_DepthFunc(expr, stack, fp)
	case OP_GL_LIGHTFV:
		op_gl_Lightfv(expr, stack, fp)
	case OP_GL_FRUSTUM:
		op_gl_Frustum(expr, stack, fp)
	case OP_GL_DISABLE:
		op_gl_Disable(expr, stack, fp)
	case OP_GL_HINT:
		op_gl_Hint(expr, stack, fp)
	case OP_GL_NEW_TEXTURE:
		op_gl_NewTexture(expr, stack, fp)
	case OP_GL_DEPTH_MASK:
		op_gl_DepthMask(expr, stack, fp)
	case OP_GL_TEX_ENVI:
		op_gl_TexEnvi(expr, stack, fp)
	case OP_GL_BLEND_FUNC:
		op_gl_BlendFunc(expr, stack, fp)
	case OP_GL_ORTHO:
		op_gl_Ortho(expr, stack, fp)
	case OP_GL_VIEWPORT:
		op_gl_Viewport(expr, stack, fp)
	case OP_GL_SCALEF:
		op_gl_Scalef(expr, stack, fp)
	case OP_GL_TEX_COORD_2D:
		op_gl_TexCoord2d(expr, stack, fp)

		// glfw
	case OP_GLFW_INIT:
		op_glfw_Init(expr, stack, fp)
	case OP_GLFW_WINDOW_HINT:
		op_glfw_WindowHint(expr, stack, fp)
	case OP_GLFW_CREATE_WINDOW:
		op_glfw_CreateWindow(expr, stack, fp)
	case OP_GLFW_MAKE_CONTEXT_CURRENT:
		op_glfw_MakeContextCurrent(expr, stack, fp)
	case OP_GLFW_SHOULD_CLOSE:
		op_glfw_ShouldClose(expr, stack, fp)
	case OP_GLFW_SET_SHOULD_CLOSE:
		op_glfw_SetShouldClose(expr, stack, fp)
	case OP_GLFW_POLL_EVENTS:
		op_glfw_PollEvents()
	case OP_GLFW_SWAP_BUFFERS:
		op_glfw_SwapBuffers(expr, stack, fp)
	case OP_GLFW_GET_FRAMEBUFFER_SIZE:
		op_glfw_GetFramebufferSize(expr, stack, fp)
	case OP_GLFW_SET_KEY_CALLBACK:
		op_glfw_SetKeyCallback(expr, stack, fp)
	case OP_GLFW_GET_TIME:
		op_glfw_GetTime(expr, stack, fp)
	case OP_GLFW_SET_MOUSE_BUTTON_CALLBACK:
		op_glfw_SetMouseButtonCallback(expr, stack, fp)
	case OP_GLFW_SET_CURSOR_POS_CALLBACK:
		op_glfw_SetCursorPosCallback(expr, stack, fp)
	case OP_GLFW_GET_CURSOR_POS:
		op_glfw_GetCursorPos(expr, stack, fp)
	case OP_GLFW_SET_INPUT_MODE:
		op_glfw_SetInputMode(expr, stack, fp)
	case OP_HTTP_GET:
		op_http_get(expr, stack, fp)
	}
}

// For the parser. These shouldn't be used in the runtime for performance reasons
var OpNames map[int]string = map[int]string{
	OP_IDENTITY:   "identity",
	OP_READ_ARRAY: "read",
	OP_JMP:        "jmp",
	OP_DEBUG:      "debug",

	OP_UND_EQUAL:    "eq",
	OP_UND_UNEQUAL:  "uneq",
	OP_UND_BITAND:   "bitand",
	OP_UND_BITXOR:   "bitxor",
	OP_UND_BITOR:    "bitor",
	OP_UND_BITCLEAR: "bitclear",
	OP_UND_MUL:      "mul",
	OP_UND_DIV:      "div",
	OP_UND_MOD:      "mod",
	OP_UND_ADD:      "add",
	OP_UND_SUB:      "sub",
	OP_UND_BITSHL:   "bitshl",
	OP_UND_LT:       "lt",
	OP_UND_GT:       "gt",
	OP_UND_LTEQ:     "lteq",
	OP_UND_GTEQ:     "gteq",
	OP_UND_LEN:      "len",
	OP_UND_PRINTF:   "printf",
	OP_UND_SPRINTF:  "sprintf",

	OP_BYTE_PRINT: "byte.print",

	OP_BOOL_PRINT:   "bool.print",
	OP_BOOL_EQUAL:   "bool.eq",
	OP_BOOL_UNEQUAL: "bool.uneq",
	OP_BOOL_NOT:     "bool.not",
	OP_BOOL_OR:      "bool.or",
	OP_BOOL_AND:     "bool.and",

	OP_I32_BYTE:     "i32.byte",
	OP_I32_STR:      "i32.str",
	OP_I32_I32:      "i32.i32",
	OP_I32_I64:      "i32.i64",
	OP_I32_F32:      "i32.f32",
	OP_I32_F64:      "i32.f64",
	OP_I32_PRINT:    "i32.print",
	OP_I32_ADD:      "i32.add",
	OP_I32_SUB:      "i32.sub",
	OP_I32_MUL:      "i32.mul",
	OP_I32_DIV:      "i32.div",
	OP_I32_ABS:      "i32.abs",
	OP_I32_POW:      "i32.pow",
	OP_I32_GT:       "i32.gt",
	OP_I32_GTEQ:     "i32.gteq",
	OP_I32_LT:       "i32.lt",
	OP_I32_LTEQ:     "i32.lteq",
	OP_I32_EQ:       "i32.eq",
	OP_I32_UNEQ:     "i32.uneq",
	OP_I32_MOD:      "i32.mod",
	OP_I32_RAND:     "i32.rand",
	OP_I32_BITAND:   "i32.bitand",
	OP_I32_BITOR:    "i32.bitor",
	OP_I32_BITXOR:   "i32.bitxor",
	OP_I32_BITCLEAR: "i32.bitclear",
	OP_I32_BITSHL:   "i32.bitshl",
	OP_I32_BITSHR:   "i32.bitshr",
	OP_I32_SQRT:     "i32.sqrt",
	OP_I32_LOG:      "i32.log",
	OP_I32_LOG2:     "i32.log2",
	OP_I32_LOG10:    "i32.log10",
	OP_I32_MAX:      "i32.max",
	OP_I32_MIN:      "i32.min",
	OP_I32_SIN:      "i32.sin",
	OP_I32_COS:      "i32.cos",
	OP_I64_PRINT:    "i64.print",
	OP_I64_ADD:      "i64.add",
	OP_I64_SUB:      "i64.sub",
	OP_I64_MUL:      "i64.mul",
	OP_I64_DIV:      "i64.div",
	OP_I64_ABS:      "i64.abs",
	OP_I64_POW:      "i64.pow",
	OP_I64_GT:       "i64.gt",
	OP_I64_GTEQ:     "i64.gteq",
	OP_I64_LT:       "i64.lt",
	OP_I64_LTEQ:     "i64.lteq",
	OP_I64_EQ:       "i64.eq",
	OP_I64_UNEQ:     "i64.uneq",
	OP_I64_MOD:      "i64.mod",
	OP_I64_RAND:     "i64.rand",
	OP_I64_BITAND:   "i64.bitand",
	OP_I64_BITOR:    "i64.bitor",
	OP_I64_BITXOR:   "i64.bitxor",
	OP_I64_BITCLEAR: "i64.bitclear",
	OP_I64_BITSHL:   "i64.bitshl",
	OP_I64_BITSHR:   "i64.bitshr",
	OP_I64_SQRT:     "i64.sqrt",
	OP_I64_LOG:      "i64.log",
	OP_I64_LOG2:     "i64.log2",
	OP_I64_LOG10:    "i64.log10",
	OP_I64_MAX:      "i64.max",
	OP_I64_MIN:      "i64.min",
	OP_I64_COS:      "i64.cos",
	OP_I64_SIN:      "i64.sin",
	OP_F32_PRINT:    "f32.print",
	OP_F32_ADD:      "f32.add",
	OP_F32_SUB:      "f32.sub",
	OP_F32_MUL:      "f32.mul",
	OP_F32_DIV:      "f32.div",
	OP_F32_ABS:      "f32.abs",
	OP_F32_POW:      "f32.pow",
	OP_F32_GT:       "f32.gt",
	OP_F32_GTEQ:     "f32.gteq",
	OP_F32_LT:       "f32.lt",
	OP_F32_LTEQ:     "f32.lteq",
	OP_F32_EQ:       "f32.eq",
	OP_F32_UNEQ:     "f32.uneq",
	OP_F32_COS:      "f32.cos",
	OP_F32_SIN:      "f32.sin",
	OP_F32_SQRT:     "f32.sqrt",
	OP_F32_LOG:      "f32.log",
	OP_F32_LOG2:     "f32.log2",
	OP_F32_LOG10:    "f32.log10",
	OP_F32_MAX:      "f32.max",
	OP_F32_MIN:      "f32.min",
	OP_F64_PRINT:    "f64.print",
	OP_F64_ADD:      "f64.add",
	OP_F64_SUB:      "f64.sub",
	OP_F64_MUL:      "f64.mul",
	OP_F64_DIV:      "f64.div",
	OP_F64_ABS:      "f64.abs",
	OP_F64_POW:      "f64.pow",
	OP_F64_GT:       "f64.gt",
	OP_F64_GTEQ:     "f64.gteq",
	OP_F64_LT:       "f64.lt",
	OP_F64_LTEQ:     "f64.lteq",
	OP_F64_EQ:       "f64.eq",
	OP_F64_UNEQ:     "f64.uneq",
	OP_F64_COS:      "f64.cos",
	OP_F64_SIN:      "f64.sin",
	OP_F64_SQRT:     "f64.sqrt",
	OP_F64_LOG:      "f64.log",
	OP_F64_LOG2:     "f64.log2",
	OP_F64_LOG10:    "f64.log10",
	OP_F64_MAX:      "f64.max",
	OP_F64_MIN:      "f64.min",

	OP_STR_PRINT: "str.print",
	OP_STR_EQ: "str.eq",

	OP_TIME_SLEEP:      "time.Sleep",
	OP_TIME_UNIX_MILLI: "time.UnixMilli",
	OP_TIME_UNIX_NANO:  "time.UnixNano",

	OP_TEST_START: "test.start",
	OP_TEST_STOP:  "test.stop",
	OP_TEST_ERROR: "test.error",

	OP_TEST: "test",

	// opengl
	OP_GL_INIT:                       "gl.Init",
	OP_GL_CREATE_PROGRAM:             "gl.CreateProgram",
	OP_GL_LINK_PROGRAM:               "gl.LinkProgram",
	OP_GL_CLEAR:                      "gl.Clear",
	OP_GL_USE_PROGRAM:                "gl.UseProgram",
	OP_GL_BIND_BUFFER:                "gl.BindBuffer",
	OP_GL_BIND_VERTEX_ARRAY:          "gl.BindVertexArray",
	OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY: "gl.EnableVertexAttribArray",
	OP_GL_VERTEX_ATTRIB_POINTER:      "gl.VertexAttribPointer",
	OP_GL_DRAW_ARRAYS:                "gl.DrawArrays",
	OP_GL_GEN_BUFFERS:                "gl.GenBuffers",
	OP_GL_BUFFER_DATA:                "gl.BufferData",
	OP_GL_GEN_VERTEX_ARRAYS:          "gl.GenVertexArrays",
	OP_GL_CREATE_SHADER:              "gl.CreateShader",
	OP_GL_STRS:                       "gl.Strs",
	OP_GL_FREE:                       "gl.Free",
	OP_GL_SHADER_SOURCE:              "gl.ShaderSource",
	OP_GL_COMPILE_SHADER:             "gl.CompileShader",
	OP_GL_GET_SHADERIV:               "gl.GetShaderiv",
	OP_GL_ATTACH_SHADER:              "gl.AttachShader",
	OP_GL_MATRIX_MODE:                "gl.MatrixMode",
	OP_GL_ROTATEF:                    "gl.Rotatef",
	OP_GL_TRANSLATEF:                 "gl.Translatef",
	OP_GL_LOAD_IDENTITY:              "gl.LoadIdentity",
	OP_GL_PUSH_MATRIX:                "gl.PushMatrix",
	OP_GL_POP_MATRIX:                 "gl.PopMatrix",
	OP_GL_ENABLE_CLIENT_STATE:        "gl.EnableClientState",
	OP_GL_BIND_TEXTURE:               "gl.BindTexture",
	OP_GL_COLOR3F:                    "gl.Color3f",
	OP_GL_COLOR4F:                    "gl.Color4f",
	OP_GL_BEGIN:                      "gl.Begin",
	OP_GL_END:                        "gl.End",
	OP_GL_NORMAL3F:                   "gl.Normal3f",
	OP_GL_VERTEX_2F:                  "gl.Vertex2f",
	OP_GL_VERTEX_3F:                  "gl.Vertex3f",
	OP_GL_ENABLE:                     "gl.Enable",
	OP_GL_CLEAR_COLOR:                "gl.ClearColor",
	OP_GL_CLEAR_DEPTH:                "gl.ClearDepth",
	OP_GL_DEPTH_FUNC:                 "gl.DepthFunc",
	OP_GL_LIGHTFV:                    "gl.Lightfv",
	OP_GL_FRUSTUM:                    "gl.Frustum",
	OP_GL_DISABLE:                    "gl.Disable",
	OP_GL_HINT:                       "gl.Hint",
	OP_GL_NEW_TEXTURE:                "gl.NewTexture",
	OP_GL_DEPTH_MASK:                 "gl.DepthMask",
	OP_GL_TEX_ENVI:                   "gl.TexEnvi",
	OP_GL_BLEND_FUNC:                 "gl.BlendFunc",
	OP_GL_ORTHO:                      "gl.Ortho",
	OP_GL_VIEWPORT:                   "gl.Viewport",
	OP_GL_SCALEF:                     "gl.Scalef",
	OP_GL_TEX_COORD_2D:               "gl.TexCoord2d",

	// glfw
	OP_GLFW_INIT:                      "glfw.Init",
	OP_GLFW_WINDOW_HINT:               "glfw.WindowHint",
	OP_GLFW_CREATE_WINDOW:             "glfw.CreateWindow",
	OP_GLFW_MAKE_CONTEXT_CURRENT:      "glfw.MakeContextCurrent",
	OP_GLFW_SHOULD_CLOSE:              "glfw.ShouldClose",
	OP_GLFW_SET_SHOULD_CLOSE:          "glfw.SetShouldClose",
	OP_GLFW_POLL_EVENTS:               "glfw.PollEvents",
	OP_GLFW_SWAP_BUFFERS:              "glfw.SwapBuffers",
	OP_GLFW_GET_FRAMEBUFFER_SIZE:      "glfw.GetFramebufferSize",
	OP_GLFW_SET_KEY_CALLBACK:          "glfw.SetKeyCallback",
	OP_GLFW_GET_TIME:                  "glfw.GetTime",
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK: "glfw.SetMouseButtonCallback",
	OP_GLFW_SET_CURSOR_POS_CALLBACK:   "glfw.SetCursorPosCallback",
	OP_GLFW_GET_CURSOR_POS:            "glfw.GetCursorPos",
	OP_GLFW_SET_INPUT_MODE:            "glfw.SetInputMode",
	// http
	OP_HTTP_GET:                       "http.Get",
}

// For the parser. These shouldn't be used in the runtime for performance reasons
var OpCodes map[string]int = map[string]int{
	"identity": OP_IDENTITY,
	"read":     OP_READ_ARRAY,
	"jmp":      OP_JMP,
	"debug":    OP_DEBUG,

	"equal":    OP_UND_EQUAL,
	"unequal":  OP_UND_UNEQUAL,
	"bitand":   OP_UND_BITAND,
	"bitxor":   OP_UND_BITXOR,
	"bitor":    OP_UND_BITOR,
	"bitclear": OP_UND_BITCLEAR,
	"mul":      OP_UND_MUL,
	"div":      OP_UND_DIV,
	"mod":      OP_UND_MOD,
	"add":      OP_UND_ADD,
	"sub":      OP_UND_SUB,
	"bitshl":   OP_UND_BITSHL,
	"lt":       OP_UND_LT,
	"gt":       OP_UND_GT,
	"lteq":     OP_UND_LTEQ,
	"gteq":     OP_UND_GTEQ,
	"len":      OP_UND_LEN,
	"printf":   OP_UND_PRINTF,
	"sprintf":  OP_UND_SPRINTF,

	"byte.print": OP_BYTE_PRINT,

	"bool.print": OP_BOOL_PRINT,
	"bool.eq":    OP_BOOL_EQUAL,
	"bool.uneq":  OP_BOOL_UNEQUAL,
	"bool.not":   OP_BOOL_NOT,
	"bool.or":    OP_BOOL_OR,
	"bool.and":   OP_BOOL_AND,

	"i32.byte":     OP_I32_BYTE,
	"i32.str":      OP_I32_STR,
	"i32.i32":      OP_I32_I32,
	"i32.i64":      OP_I32_I64,
	"i32.f32":      OP_I32_F32,
	"i32.f64":      OP_I32_F64,
	"i32.print":    OP_I32_PRINT,
	"i32.add":      OP_I32_ADD,
	"i32.sub":      OP_I32_SUB,
	"i32.mul":      OP_I32_MUL,
	"i32.div":      OP_I32_DIV,
	"i32.abs":      OP_I32_ABS,
	"i32.pow":      OP_I32_POW,
	"i32.gt":       OP_I32_GT,
	"i32.gteq":     OP_I32_GTEQ,
	"i32.lt":       OP_I32_LT,
	"i32.lteq":     OP_I32_LTEQ,
	"i32.eq":       OP_I32_EQ,
	"i32.uneq":     OP_I32_UNEQ,
	"i32.mod":      OP_I32_MOD,
	"i32.rand":     OP_I32_RAND,
	"i32.bitand":   OP_I32_BITAND,
	"i32.bitor":    OP_I32_BITOR,
	"i32.bitxor":   OP_I32_BITXOR,
	"i32.bitclear": OP_I32_BITCLEAR,
	"i32.bitshl":   OP_I32_BITSHL,
	"i32.bitshr":   OP_I32_BITSHR,
	"i32.sqrt":     OP_I32_SQRT,
	"i32.log":      OP_I32_LOG,
	"i32.log2":     OP_I32_LOG2,
	"i32.log10":    OP_I32_LOG10,
	"i32.max":      OP_I32_MAX,
	"i32.min":      OP_I32_MIN,
	"i32.sin":      OP_I32_SIN,
	"i32.cos":      OP_I32_COS,
	"i64.print":    OP_I64_PRINT,
	"i64.add":      OP_I64_ADD,
	"i64.sub":      OP_I64_SUB,
	"i64.mul":      OP_I64_MUL,
	"i64.div":      OP_I64_DIV,
	"i64.abs":      OP_I64_ABS,
	"i64.pow":      OP_I64_POW,
	"i64.gt":       OP_I64_GT,
	"i64.gteq":     OP_I64_GTEQ,
	"i64.lt":       OP_I64_LT,
	"i64.lteq":     OP_I64_LTEQ,
	"i64.eq":       OP_I64_EQ,
	"i64.uneq":     OP_I64_UNEQ,
	"i64.mod":      OP_I64_MOD,
	"i64.rand":     OP_I64_RAND,
	"i64.bitand":   OP_I64_BITAND,
	"i64.bitor":    OP_I64_BITOR,
	"i64.bitxor":   OP_I64_BITXOR,
	"i64.bitclear": OP_I64_BITCLEAR,
	"i64.bitshl":   OP_I64_BITSHL,
	"i64.bitshr":   OP_I64_BITSHR,
	"i64.sqrt":     OP_I64_SQRT,
	"i64.log":      OP_I64_LOG,
	"i64.log2":     OP_I64_LOG2,
	"i64.log10":    OP_I64_LOG10,
	"i64.max":      OP_I64_MAX,
	"i64.min":      OP_I64_MIN,
	"i64.cos":      OP_I64_COS,
	"i64.sin":      OP_I64_SIN,
	"f32.print":    OP_F32_PRINT,
	"f32.add":      OP_F32_ADD,
	"f32.sub":      OP_F32_SUB,
	"f32.mul":      OP_F32_MUL,
	"f32.div":      OP_F32_DIV,
	"f32.abs":      OP_F32_ABS,
	"f32.pow":      OP_F32_POW,
	"f32.gt":       OP_F32_GT,
	"f32.gteq":     OP_F32_GTEQ,
	"f32.lt":       OP_F32_LT,
	"f32.lteq":     OP_F32_LTEQ,
	"f32.eq":       OP_F32_EQ,
	"f32.uneq":     OP_F32_UNEQ,
	"f32.cos":      OP_F32_COS,
	"f32.sin":      OP_F32_SIN,
	"f32.sqrt":     OP_F32_SQRT,
	"f32.log":      OP_F32_LOG,
	"f32.log2":     OP_F32_LOG2,
	"f32.log10":    OP_F32_LOG10,
	"f32.max":      OP_F32_MAX,
	"f32.min":      OP_F32_MIN,
	"f64.print":    OP_F64_PRINT,
	"f64.add":      OP_F64_ADD,
	"f64.sub":      OP_F64_SUB,
	"f64.mul":      OP_F64_MUL,
	"f64.div":      OP_F64_DIV,
	"f64.abs":      OP_F64_ABS,
	"f64.pow":      OP_F64_POW,
	"f64.gt":       OP_F64_GT,
	"f64.gteq":     OP_F64_GTEQ,
	"f64.lt":       OP_F64_LT,
	"f64.lteq":     OP_F64_LTEQ,
	"f64.eq":       OP_F64_EQ,
	"f64.uneq":     OP_F64_UNEQ,
	"f64.cos":      OP_F64_COS,
	"f64.sin":      OP_F64_SIN,
	"f64.sqrt":     OP_F64_SQRT,
	"f64.log":      OP_F64_LOG,
	"f64.log2":     OP_F64_LOG2,
	"f64.log10":    OP_F64_LOG10,
	"f64.max":      OP_F64_MAX,
	"f64.min":      OP_F64_MIN,

	"str.print": OP_STR_PRINT,
	"str.eq": OP_STR_EQ,

	"time.Sleep":     OP_TIME_SLEEP,
	"time.UnixMilli": OP_TIME_UNIX_MILLI,
	"time.UnixNano":  OP_TIME_UNIX_NANO,

	"test.start": OP_TEST_START,
	"test.stop":  OP_TEST_STOP,
	"test":       OP_TEST,

	// opengl
	"gl.Init":                    OP_GL_INIT,
	"gl.CreateProgram":           OP_GL_CREATE_PROGRAM,
	"gl.LinkProgram":             OP_GL_LINK_PROGRAM,
	"gl.Clear":                   OP_GL_CLEAR,
	"gl.UseProgram":              OP_GL_USE_PROGRAM,
	"gl.BindBuffer":              OP_GL_BIND_BUFFER,
	"gl.BindVertexArray":         OP_GL_BIND_VERTEX_ARRAY,
	"gl.EnableVertexAttribArray": OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY,
	"gl.VertexAttribPointer":     OP_GL_VERTEX_ATTRIB_POINTER,
	"gl.DrawArrays":              OP_GL_DRAW_ARRAYS,
	"gl.GenBuffers":              OP_GL_GEN_BUFFERS,
	"gl.BufferData":              OP_GL_BUFFER_DATA,
	"gl.GenVertexArrays":         OP_GL_GEN_VERTEX_ARRAYS,
	"gl.CreateShader":            OP_GL_CREATE_SHADER,
	"gl.Strs":                    OP_GL_STRS,
	"gl.Free":                    OP_GL_FREE,
	"gl.ShaderSource":            OP_GL_SHADER_SOURCE,
	"gl.CompileShader":           OP_GL_COMPILE_SHADER,
	"gl.GetShaderiv":             OP_GL_GET_SHADERIV,
	"gl.AttachShader":            OP_GL_ATTACH_SHADER,
	"gl.MatrixMode":              OP_GL_MATRIX_MODE,
	"gl.Rotatef":                 OP_GL_ROTATEF,
	"gl.Translatef":              OP_GL_TRANSLATEF,
	"gl.LoadIdentity":            OP_GL_LOAD_IDENTITY,
	"gl.PushMatrix":              OP_GL_PUSH_MATRIX,
	"gl.PopMatrix":               OP_GL_POP_MATRIX,
	"gl.EnableClientState":       OP_GL_ENABLE_CLIENT_STATE,
	"gl.BindTexture":             OP_GL_BIND_TEXTURE,
	"gl.Color3f":                 OP_GL_COLOR3F,
	"gl.Color4f":                 OP_GL_COLOR4F,
	"gl.Begin":                   OP_GL_BEGIN,
	"gl.End":                     OP_GL_END,
	"gl.Normal3f":                OP_GL_NORMAL3F,
	"gl.Vertex2f":                OP_GL_VERTEX_2F,
	"gl.Vertex3f":                OP_GL_VERTEX_3F,
	"gl.Enable":                  OP_GL_ENABLE,
	"gl.ClearColor":              OP_GL_CLEAR_COLOR,
	"gl.ClearDepth":              OP_GL_CLEAR_DEPTH,
	"gl.DepthFunc":               OP_GL_DEPTH_FUNC,
	"gl.Lightfv":                 OP_GL_LIGHTFV,
	"gl.Frustum":                 OP_GL_FRUSTUM,
	"gl.Disable":                 OP_GL_DISABLE,
	"gl.Hint":                    OP_GL_HINT,
	"gl.NewTexture":              OP_GL_NEW_TEXTURE,
	"gl.DepthMask":               OP_GL_DEPTH_MASK,
	"gl.TexEnvi":                 OP_GL_TEX_ENVI,
	"gl.BlendFunc":               OP_GL_BLEND_FUNC,
	"gl.Ortho":                   OP_GL_ORTHO,
	"gl.Viewport":                OP_GL_VIEWPORT,
	"gl.Scalef":                  OP_GL_SCALEF,
	"gl.TexCoord2d":              OP_GL_TEX_COORD_2D,

	// glfw
	"glfw.Init":                   OP_GLFW_INIT,
	"glfw.WindowHint":             OP_GLFW_WINDOW_HINT,
	"glfw.CreateWindow":           OP_GLFW_CREATE_WINDOW,
	"glfw.MakeContextCurrent":     OP_GLFW_MAKE_CONTEXT_CURRENT,
	"glfw.ShouldClose":            OP_GLFW_SHOULD_CLOSE,
	"glfw.SetShouldClose":         OP_GLFW_SET_SHOULD_CLOSE,
	"glfw.PollEvents":             OP_GLFW_POLL_EVENTS,
	"glfw.SwapBuffers":            OP_GLFW_SWAP_BUFFERS,
	"glfw.GetFramebufferSize":     OP_GLFW_GET_FRAMEBUFFER_SIZE,
	"glfw.SetKeyCallback":         OP_GLFW_SET_KEY_CALLBACK,
	"glfw.GetTime":                OP_GLFW_GET_TIME,
	"glfw.SetMouseButtonCallback": OP_GLFW_SET_MOUSE_BUTTON_CALLBACK,
	"glfw.SetCursorPosCallback":   OP_GLFW_SET_CURSOR_POS_CALLBACK,
	"glfw.GetCursorPos":           OP_GLFW_GET_CURSOR_POS,
	"glfw.SetInputMode":           OP_GLFW_SET_INPUT_MODE,
	// http
	"http.Get":                    OP_HTTP_GET,
}

var Natives map[int]*CXFunction = map[int]*CXFunction{
	OP_IDENTITY:     MakeNative(OP_IDENTITY, []int{TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_READ_ARRAY:   MakeNative(OP_READ_ARRAY, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_JMP:          MakeNative(OP_JMP, []int{TYPE_BOOL}, []int{}),
	OP_DEBUG:        MakeNative(OP_DEBUG, []int{}, []int{}),

	OP_UND_EQUAL:    MakeNative(OP_UND_EQUAL, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
	OP_UND_UNEQUAL:  MakeNative(OP_UND_UNEQUAL, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
	OP_UND_BITAND:   MakeNative(OP_UND_BITAND, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_BITXOR:   MakeNative(OP_UND_BITXOR, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_BITOR:    MakeNative(OP_UND_BITOR, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_BITCLEAR: MakeNative(OP_UND_BITCLEAR, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_MUL:      MakeNative(OP_UND_MUL, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_DIV:      MakeNative(OP_UND_DIV, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_MOD:      MakeNative(OP_UND_MOD, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_ADD:      MakeNative(OP_UND_ADD, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_SUB:      MakeNative(OP_UND_SUB, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_BITSHL:   MakeNative(OP_UND_BITSHL, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_BITSHR:   MakeNative(OP_UND_BITSHR, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_UND_LT:       MakeNative(OP_UND_LT, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
	OP_UND_GT:       MakeNative(OP_UND_GT, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
	OP_UND_LTEQ:     MakeNative(OP_UND_LTEQ, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
	OP_UND_GTEQ:     MakeNative(OP_UND_GTEQ, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
	OP_UND_LEN:      MakeNative(OP_UND_LEN, []int{TYPE_UNDEFINED}, []int{TYPE_I32}),
	OP_UND_PRINTF:   MakeNative(OP_UND_PRINTF, []int{TYPE_UNDEFINED}, []int{}),
	OP_UND_SPRINTF:  MakeNative(OP_UND_SPRINTF, []int{TYPE_UNDEFINED}, []int{TYPE_STR}),

	OP_BYTE_PRINT:   MakeNative(OP_BYTE_PRINT, []int{TYPE_BYTE}, []int{}),

	OP_BOOL_PRINT:   MakeNative(OP_BOOL_PRINT, []int{TYPE_BOOL}, []int{}),
	OP_BOOL_EQUAL:   MakeNative(OP_BOOL_EQUAL, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),
	OP_BOOL_UNEQUAL: MakeNative(OP_BOOL_UNEQUAL, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),
	OP_BOOL_NOT:     MakeNative(OP_BOOL_NOT, []int{TYPE_BOOL}, []int{TYPE_BOOL}),
	OP_BOOL_OR:      MakeNative(OP_BOOL_OR, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),
	OP_BOOL_AND:     MakeNative(OP_BOOL_AND, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),

	OP_I32_BYTE:     MakeNative(OP_I32_BYTE, []int{TYPE_I32}, []int{TYPE_BYTE}),
	OP_I32_STR:      MakeNative(OP_I32_STR, []int{TYPE_I32}, []int{TYPE_STR}),
	OP_I32_I32:      MakeNative(OP_I32_I32, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_I64:      MakeNative(OP_I32_I64, []int{TYPE_I32}, []int{TYPE_I64}),
	OP_I32_F32:      MakeNative(OP_I32_F32, []int{TYPE_I32}, []int{TYPE_F32}),
	OP_I32_F64:      MakeNative(OP_I32_F64, []int{TYPE_I32}, []int{TYPE_F64}),

	OP_I32_PRINT:    MakeNative(OP_I32_PRINT, []int{TYPE_I32}, []int{}),
	OP_I32_ADD:      MakeNative(OP_I32_ADD, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_SUB:      MakeNative(OP_I32_SUB, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_MUL:      MakeNative(OP_I32_MUL, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_DIV:      MakeNative(OP_I32_DIV, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_ABS:      MakeNative(OP_I32_ABS, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_POW:      MakeNative(OP_I32_POW, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_GT:       MakeNative(OP_I32_GT, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_GTEQ:     MakeNative(OP_I32_GTEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_LT:       MakeNative(OP_I32_LT, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_LTEQ:     MakeNative(OP_I32_LTEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_EQ:       MakeNative(OP_I32_EQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_UNEQ:     MakeNative(OP_I32_UNEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_MOD:      MakeNative(OP_I32_MOD, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_RAND:     MakeNative(OP_I32_RAND, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITAND:   MakeNative(OP_I32_BITAND, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITOR:    MakeNative(OP_I32_BITOR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITXOR:   MakeNative(OP_I32_BITXOR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITCLEAR: MakeNative(OP_I32_BITCLEAR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITSHL:   MakeNative(OP_I32_BITSHL, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITSHR:   MakeNative(OP_I32_BITSHR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_SQRT:     MakeNative(OP_I32_SQRT, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_LOG:      MakeNative(OP_I32_LOG, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_LOG2:     MakeNative(OP_I32_LOG2, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_LOG10:    MakeNative(OP_I32_LOG10, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_MAX:      MakeNative(OP_I32_MAX, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_MIN:      MakeNative(OP_I32_MIN, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_SIN:      MakeNative(OP_I32_SIN, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_COS:      MakeNative(OP_I32_COS, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I64_PRINT:    MakeNative(OP_I64_PRINT, []int{TYPE_I64}, []int{}),
	OP_I64_ADD:      MakeNative(OP_I64_ADD, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_SUB:      MakeNative(OP_I64_SUB, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_MUL:      MakeNative(OP_I64_MUL, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_DIV:      MakeNative(OP_I64_DIV, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_ABS:      MakeNative(OP_I64_ABS, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_POW:      MakeNative(OP_I64_POW, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_GT:       MakeNative(OP_I64_GT, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_GTEQ:     MakeNative(OP_I64_GTEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_LT:       MakeNative(OP_I64_LT, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_LTEQ:     MakeNative(OP_I64_LTEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_EQ:       MakeNative(OP_I64_EQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_UNEQ:     MakeNative(OP_I64_UNEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_MOD:      MakeNative(OP_I64_MOD, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_RAND:     MakeNative(OP_I64_RAND, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITAND:   MakeNative(OP_I64_BITAND, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITOR:    MakeNative(OP_I64_BITOR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITXOR:   MakeNative(OP_I64_BITXOR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITCLEAR: MakeNative(OP_I64_BITCLEAR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITSHL:   MakeNative(OP_I64_BITSHL, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITSHR:   MakeNative(OP_I64_BITSHR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_SQRT:     MakeNative(OP_I64_SQRT, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_LOG:      MakeNative(OP_I64_LOG, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_LOG2:     MakeNative(OP_I64_LOG2, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_LOG10:    MakeNative(OP_I64_LOG10, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_MAX:      MakeNative(OP_I64_MAX, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_MIN:      MakeNative(OP_I64_MIN, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_COS:      MakeNative(OP_I64_COS, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_SIN:      MakeNative(OP_I64_SIN, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_F32_PRINT:    MakeNative(OP_F32_PRINT, []int{TYPE_F32}, []int{}),
	OP_F32_ADD:      MakeNative(OP_F32_ADD, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_SUB:      MakeNative(OP_F32_SUB, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_MUL:      MakeNative(OP_F32_MUL, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_DIV:      MakeNative(OP_F32_DIV, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_ABS:      MakeNative(OP_F32_ABS, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_POW:      MakeNative(OP_F32_POW, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_GT:       MakeNative(OP_F32_GT, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_GTEQ:     MakeNative(OP_F32_GTEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_LT:       MakeNative(OP_F32_LT, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_LTEQ:     MakeNative(OP_F32_LTEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_EQ:       MakeNative(OP_F32_EQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_UNEQ:     MakeNative(OP_F32_UNEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_COS:      MakeNative(OP_F32_COS, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_SIN:      MakeNative(OP_F32_SIN, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_SQRT:     MakeNative(OP_F32_SQRT, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_LOG:      MakeNative(OP_F32_LOG, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_LOG2:     MakeNative(OP_F32_LOG2, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_LOG10:    MakeNative(OP_F32_LOG10, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_MAX:      MakeNative(OP_F32_MAX, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_MIN:      MakeNative(OP_F32_MIN, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F64_PRINT:    MakeNative(OP_F64_PRINT, []int{TYPE_F64}, []int{}),
	OP_F64_ADD:      MakeNative(OP_F64_ADD, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_SUB:      MakeNative(OP_F64_SUB, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_MUL:      MakeNative(OP_F64_MUL, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_DIV:      MakeNative(OP_F64_DIV, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_ABS:      MakeNative(OP_F64_ABS, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_POW:      MakeNative(OP_F64_POW, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_GT:       MakeNative(OP_F64_GT, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_GTEQ:     MakeNative(OP_F64_GTEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_LT:       MakeNative(OP_F64_LT, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_LTEQ:     MakeNative(OP_F64_LTEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_EQ:       MakeNative(OP_F64_EQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_UNEQ:     MakeNative(OP_F64_UNEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_COS:      MakeNative(OP_F64_COS, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_SIN:      MakeNative(OP_F64_SIN, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_SQRT:     MakeNative(OP_F64_SQRT, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_LOG:      MakeNative(OP_F64_LOG, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_LOG2:     MakeNative(OP_F64_LOG2, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_LOG10:    MakeNative(OP_F64_LOG10, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_MIN:      MakeNative(OP_F64_MIN, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_MAX:      MakeNative(OP_F64_MAX, []int{TYPE_F64}, []int{TYPE_F64}),
	
	OP_STR_PRINT:    MakeNative(OP_STR_PRINT, []int{TYPE_STR}, []int{}),
	OP_STR_EQ:       MakeNative(OP_STR_EQ, []int{TYPE_STR, TYPE_STR}, []int{TYPE_BOOL}),

	OP_TIME_SLEEP:      MakeNative(OP_TIME_SLEEP, []int{TYPE_I32}, []int{}),
	OP_TIME_UNIX_MILLI: MakeNative(OP_TIME_UNIX_MILLI, []int{}, []int{TYPE_I64}),
	OP_TIME_UNIX_NANO:  MakeNative(OP_TIME_UNIX_NANO, []int{}, []int{TYPE_I64}),

	OP_TEST_START: MakeNative(OP_TEST_START, []int{}, []int{}),
	OP_TEST_STOP:  MakeNative(OP_TEST_START, []int{}, []int{}),
	OP_TEST:       MakeNative(OP_TEST, []int{TYPE_UNDEFINED, TYPE_UNDEFINED, TYPE_STR}, []int{}),

	// opengl
	OP_GL_INIT:                       MakeNative(OP_GL_INIT, []int{}, []int{}),
	OP_GL_CREATE_PROGRAM:             MakeNative(OP_GL_CREATE_PROGRAM, []int{}, []int{TYPE_I32}),
	OP_GL_LINK_PROGRAM:               MakeNative(OP_GL_LINK_PROGRAM, []int{TYPE_I32}, []int{}),
	OP_GL_CLEAR:                      MakeNative(OP_GL_CLEAR, []int{TYPE_I32}, []int{}),
	OP_GL_USE_PROGRAM:                MakeNative(OP_GL_USE_PROGRAM, []int{TYPE_I32}, []int{}),
	OP_GL_BIND_BUFFER:                MakeNative(OP_GL_BIND_BUFFER, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_BIND_VERTEX_ARRAY:          MakeNative(OP_GL_BIND_VERTEX_ARRAY, []int{TYPE_I32}, []int{}),
	OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY: MakeNative(OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY, []int{TYPE_I32}, []int{}),
	OP_GL_VERTEX_ATTRIB_POINTER:      MakeNative(OP_GL_VERTEX_ATTRIB_POINTER, []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_BOOL, TYPE_I32}, []int{}),
	OP_GL_DRAW_ARRAYS:                MakeNative(OP_GL_DRAW_ARRAYS, []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_GEN_BUFFERS:                MakeNative(OP_GL_GEN_BUFFERS, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_GL_BUFFER_DATA:                MakeNative(OP_GL_BUFFER_DATA, []int{TYPE_I32, TYPE_I32, TYPE_F32, TYPE_I32}, []int{}),
	OP_GL_GEN_VERTEX_ARRAYS:          MakeNative(OP_GL_GEN_VERTEX_ARRAYS, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_GL_CREATE_SHADER:              MakeNative(OP_GL_CREATE_SHADER, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_GL_STRS:                       MakeNative(OP_GL_STRS, []int{TYPE_STR, TYPE_STR}, []int{}),
	OP_GL_FREE:                       MakeNative(OP_GL_FREE, []int{TYPE_STR}, []int{}),
	OP_GL_SHADER_SOURCE:              MakeNative(OP_GL_SHADER_SOURCE, []int{TYPE_I32, TYPE_I32, TYPE_STR}, []int{}),
	OP_GL_COMPILE_SHADER:             MakeNative(OP_GL_COMPILE_SHADER, []int{TYPE_I32}, []int{}),
	OP_GL_GET_SHADERIV:               MakeNative(OP_GL_GET_SHADERIV, []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_ATTACH_SHADER:              MakeNative(OP_GL_ATTACH_SHADER, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_MATRIX_MODE:                MakeNative(OP_GL_MATRIX_MODE, []int{TYPE_I32}, []int{}),
	OP_GL_ROTATEF:                    MakeNative(OP_GL_ROTATEF, []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_TRANSLATEF:                 MakeNative(OP_GL_TRANSLATEF, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_LOAD_IDENTITY:              MakeNative(OP_GL_LOAD_IDENTITY, []int{}, []int{}),
	OP_GL_PUSH_MATRIX:                MakeNative(OP_GL_PUSH_MATRIX, []int{}, []int{}),
	OP_GL_POP_MATRIX:                 MakeNative(OP_GL_POP_MATRIX, []int{}, []int{}),
	OP_GL_ENABLE_CLIENT_STATE:        MakeNative(OP_GL_ENABLE_CLIENT_STATE, []int{TYPE_I32}, []int{}),
	OP_GL_BIND_TEXTURE:               MakeNative(OP_GL_BIND_TEXTURE, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_COLOR3F:                    MakeNative(OP_GL_COLOR3F, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_COLOR4F:                    MakeNative(OP_GL_COLOR4F, []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_BEGIN:                      MakeNative(OP_GL_BEGIN, []int{TYPE_I32}, []int{}),
	OP_GL_END:                        MakeNative(OP_GL_END, []int{}, []int{}),
	OP_GL_NORMAL3F:                   MakeNative(OP_GL_NORMAL3F, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),

	OP_GL_VERTEX_2F: MakeNative(OP_GL_VERTEX_2F, []int{TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_VERTEX_3F: MakeNative(OP_GL_VERTEX_3F, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),

	OP_GL_ENABLE:       MakeNative(OP_GL_ENABLE, []int{TYPE_I32}, []int{}),
	OP_GL_CLEAR_COLOR:  MakeNative(OP_GL_CLEAR_COLOR, []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_CLEAR_DEPTH:  MakeNative(OP_GL_CLEAR_DEPTH, []int{TYPE_F64}, []int{}),
	OP_GL_DEPTH_FUNC:   MakeNative(OP_GL_DEPTH_FUNC, []int{TYPE_I32}, []int{}),
	OP_GL_LIGHTFV:      MakeNative(OP_GL_LIGHTFV, []int{TYPE_I32, TYPE_I32, TYPE_F32}, []int{}),
	OP_GL_FRUSTUM:      MakeNative(OP_GL_FRUSTUM, []int{TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64}, []int{}),
	OP_GL_DISABLE:      MakeNative(OP_GL_DISABLE, []int{TYPE_I32}, []int{}),
	OP_GL_HINT:         MakeNative(OP_GL_HINT, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_NEW_TEXTURE:  MakeNative(OP_GL_NEW_TEXTURE, []int{TYPE_STR}, []int{TYPE_I32}),
	OP_GL_DEPTH_MASK:   MakeNative(OP_GL_DEPTH_MASK, []int{TYPE_BOOL}, []int{}),
	OP_GL_TEX_ENVI:     MakeNative(OP_GL_TEX_ENVI, []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_BLEND_FUNC:   MakeNative(OP_GL_BLEND_FUNC, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_ORTHO:        MakeNative(OP_GL_ORTHO, []int{TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64}, []int{}),
	OP_GL_VIEWPORT:     MakeNative(OP_GL_VIEWPORT, []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_SCALEF:       MakeNative(OP_GL_SCALEF, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_TEX_COORD_2D: MakeNative(OP_GL_TEX_COORD_2D, []int{TYPE_F32, TYPE_F32}, []int{}),

	// glfw
	OP_GLFW_INIT:                      MakeNative(OP_GLFW_INIT, []int{}, []int{}),
	OP_GLFW_WINDOW_HINT:               MakeNative(OP_GLFW_WINDOW_HINT, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GLFW_CREATE_WINDOW:             MakeNative(OP_GLFW_CREATE_WINDOW, []int{TYPE_STR, TYPE_I32, TYPE_I32, TYPE_STR}, []int{}),
	OP_GLFW_MAKE_CONTEXT_CURRENT:      MakeNative(OP_GLFW_MAKE_CONTEXT_CURRENT, []int{TYPE_STR}, []int{}),
	OP_GLFW_SHOULD_CLOSE:              MakeNative(OP_GLFW_SHOULD_CLOSE, []int{TYPE_STR}, []int{TYPE_BOOL}),
	OP_GLFW_SET_SHOULD_CLOSE:          MakeNative(OP_GLFW_SET_SHOULD_CLOSE, []int{TYPE_STR, TYPE_BOOL}, []int{TYPE_BOOL}),
	OP_GLFW_POLL_EVENTS:               MakeNative(OP_GLFW_POLL_EVENTS, []int{}, []int{}),
	OP_GLFW_SWAP_BUFFERS:              MakeNative(OP_GLFW_SWAP_BUFFERS, []int{TYPE_STR}, []int{}),
	OP_GLFW_GET_FRAMEBUFFER_SIZE:      MakeNative(OP_GLFW_GET_FRAMEBUFFER_SIZE, []int{TYPE_STR}, []int{TYPE_I32, TYPE_I32}),
	OP_GLFW_SET_KEY_CALLBACK:          MakeNative(OP_GLFW_SET_KEY_CALLBACK, []int{TYPE_STR, TYPE_STR}, []int{}),
	OP_GLFW_GET_TIME:                  MakeNative(OP_GLFW_GET_TIME, []int{}, []int{TYPE_F64}),
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK: MakeNative(OP_GLFW_SET_MOUSE_BUTTON_CALLBACK, []int{TYPE_STR, TYPE_STR}, []int{}),
	OP_GLFW_SET_CURSOR_POS_CALLBACK:   MakeNative(OP_GLFW_SET_CURSOR_POS_CALLBACK, []int{TYPE_STR, TYPE_STR}, []int{}),
	OP_GLFW_GET_CURSOR_POS:            MakeNative(OP_GLFW_GET_CURSOR_POS, []int{TYPE_STR}, []int{TYPE_F64, TYPE_F64}),
	OP_GLFW_SET_INPUT_MODE:            MakeNative(OP_GLFW_SET_INPUT_MODE, []int{TYPE_STR, TYPE_I32, TYPE_I32}, []int{}),

	// http
	OP_HTTP_GET:           MakeNative(OP_HTTP_GET, []int{TYPE_STR}, []int{TYPE_STR}),
}
