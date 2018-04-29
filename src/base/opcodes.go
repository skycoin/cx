package base

// op codes
const (
	OP_IDENTITY = iota
	OP_READ_ARRAY
	OP_JMP
	OP_DEBUG

	OP_BOOL_PRINT
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
	
	OP_STR_PRINT
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
)

func execNative (prgrm *CXProgram) {
	call := &prgrm.CallStack[prgrm.CallCounter]
	stack := &prgrm.Stacks[0]
	expr := call.Operator.Expressions[call.Line]
	opCode := expr.Operator.OpCode
	fp := call.FramePointer

	switch opCode {
	case OP_IDENTITY: identity(expr, stack, fp)
		// case OP_READ_ARRAY: read_array(expr, stack, fp)
	case OP_JMP: jmp(expr, stack, fp, call)
	case OP_DEBUG: prgrm.PrintStack()

	case OP_BYTE_PRINT: byte_print(expr, stack, fp)

	case OP_BOOL_PRINT: bool_print(expr, stack, fp)
	case OP_BOOL_NOT: bool_not(expr, stack, fp)
	case OP_BOOL_OR: bool_or(expr, stack, fp)
	case OP_BOOL_AND: bool_and(expr, stack, fp)

	case OP_I32_BYTE: i32_i32(expr, stack, fp)
	case OP_I32_STR: i32_i32(expr, stack, fp)
	case OP_I32_I32: i32_i32(expr, stack, fp)
	case OP_I32_I64: i32_i32(expr, stack, fp)
	case OP_I32_F32: i32_i32(expr, stack, fp)
	case OP_I32_F64: i32_i32(expr, stack, fp)
	case OP_I32_PRINT: i32_print(expr, stack, fp)
	case OP_I32_ADD: i32_add(expr, stack, fp)
	case OP_I32_SUB: i32_sub(expr, stack, fp)
	case OP_I32_MUL: i32_mul(expr, stack, fp)
	case OP_I32_DIV: i32_div(expr, stack, fp)
	case OP_I32_ABS: i32_abs(expr, stack, fp)
	case OP_I32_POW: i32_pow(expr, stack, fp)
	case OP_I32_GT: i32_gt(expr, stack, fp)
	case OP_I32_GTEQ: i32_gteq(expr, stack, fp)
	case OP_I32_LT: i32_lt(expr, stack, fp)
	case OP_I32_LTEQ: i32_lteq(expr, stack, fp)
	case OP_I32_EQ: i32_eq(expr, stack, fp)
	case OP_I32_UNEQ: i32_uneq(expr, stack, fp)
	case OP_I32_MOD: i32_mod(expr, stack, fp)
	case OP_I32_RAND: i32_rand(expr, stack, fp)
	case OP_I32_BITAND: i32_bitand(expr, stack, fp)
	case OP_I32_BITOR: i32_bitor(expr, stack, fp)
	case OP_I32_BITXOR: i32_bitxor(expr, stack, fp)
	case OP_I32_BITCLEAR: i32_bitclear(expr, stack, fp)
	case OP_I32_BITSHL: i32_bitshl(expr, stack, fp)
	case OP_I32_BITSHR: i32_bitshr(expr, stack, fp)
		
	case OP_I64_PRINT: i64_print(expr, stack, fp)
	case OP_I64_ADD: i64_add(expr, stack, fp)
	case OP_I64_SUB: i64_sub(expr, stack, fp)
	case OP_I64_MUL: i64_mul(expr, stack, fp)
	case OP_I64_DIV: i64_div(expr, stack, fp)
	case OP_I64_ABS: i64_abs(expr, stack, fp)
	case OP_I64_POW: i64_pow(expr, stack, fp)
	case OP_I64_GT: i64_gt(expr, stack, fp)
	case OP_I64_GTEQ: i64_gteq(expr, stack, fp)
	case OP_I64_LT: i64_lt(expr, stack, fp)
	case OP_I64_LTEQ: i64_lteq(expr, stack, fp)
	case OP_I64_EQ: i64_eq(expr, stack, fp)
	case OP_I64_UNEQ: i64_uneq(expr, stack, fp)
	case OP_I64_MOD: i64_mod(expr, stack, fp)
	case OP_I64_RAND: i64_rand(expr, stack, fp)
	case OP_I64_BITAND: i64_bitand(expr, stack, fp)
	case OP_I64_BITOR: i64_bitor(expr, stack, fp)
	case OP_I64_BITXOR: i64_bitxor(expr, stack, fp)
	case OP_I64_BITCLEAR: i64_bitclear(expr, stack, fp)
	case OP_I64_BITSHL: i64_bitshl(expr, stack, fp)
	case OP_I64_BITSHR: i64_bitshr(expr, stack, fp)

	case OP_F32_PRINT: f32_print(expr, stack, fp)
	case OP_F32_ADD: f32_add(expr, stack, fp)
	case OP_F32_SUB: f32_sub(expr, stack, fp)
	case OP_F32_MUL: f32_mul(expr, stack, fp)
	case OP_F32_DIV: f32_div(expr, stack, fp)
	case OP_F32_ABS: f32_abs(expr, stack, fp)
	case OP_F32_POW: f32_pow(expr, stack, fp)
	case OP_F32_GT: f32_gt(expr, stack, fp)
	case OP_F32_GTEQ: f32_gteq(expr, stack, fp)
	case OP_F32_LT: f32_lt(expr, stack, fp)
	case OP_F32_LTEQ: f32_lteq(expr, stack, fp)
	case OP_F32_EQ: f32_eq(expr, stack, fp)
	case OP_F32_UNEQ: f32_uneq(expr, stack, fp)
	case OP_F32_COS: f32_cos(expr, stack, fp)
	case OP_F32_SIN: f32_sin(expr, stack, fp)

	case OP_F64_PRINT: f64_print(expr, stack, fp)
	case OP_F64_ADD: f64_add(expr, stack, fp)
	case OP_F64_SUB: f64_sub(expr, stack, fp)
	case OP_F64_MUL: f64_mul(expr, stack, fp)
	case OP_F64_DIV: f64_div(expr, stack, fp)
	case OP_F64_ABS: f64_abs(expr, stack, fp)
	case OP_F64_POW: f64_pow(expr, stack, fp)
	case OP_F64_GT: f64_gt(expr, stack, fp)
	case OP_F64_GTEQ: f64_gteq(expr, stack, fp)
	case OP_F64_LT: f64_lt(expr, stack, fp)
	case OP_F64_LTEQ: f64_lteq(expr, stack, fp)
	case OP_F64_EQ: f64_eq(expr, stack, fp)
	case OP_F64_UNEQ: f64_uneq(expr, stack, fp)
	case OP_F64_COS: f64_cos(expr, stack, fp)
	case OP_F64_SIN: f64_sin(expr, stack, fp)

		
	case OP_STR_PRINT: str_print(expr, stack, fp)
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
	case OP_GOTO: goTo(expr, call)
	case OP_REMCX:
	case OP_ADDCX:
	case OP_QUERY:
	case OP_EXECUTE:
	case OP_INDEX:
	case OP_NAME:
	case OP_EVOLVE:
	case OP_TEST_START:
	case OP_TEST_STOP:
	case OP_TEST_ERROR:
	case OP_TEST:

	case OP_TIME_SLEEP: time_Sleep(expr, stack, fp)
	case OP_TIME_UNIX:
	case OP_TIME_UNIX_MILLI: time_UnixMilli(expr, stack, fp)
	case OP_TIME_UNIX_NANO:



		// opengl
	case OP_GL_INIT: gl_Init()
	case OP_GL_CREATE_PROGRAM: gl_CreateProgram(expr, stack, fp)
	case OP_GL_LINK_PROGRAM: gl_LinkProgram(expr, stack, fp)
	case OP_GL_CLEAR: gl_Clear(expr, stack, fp)
	case OP_GL_USE_PROGRAM: gl_UseProgram(expr, stack, fp)
	case OP_GL_BIND_BUFFER: gl_BindBuffer(expr, stack, fp)
	case OP_GL_BIND_VERTEX_ARRAY: gl_BindVertexArray(expr, stack, fp)
	case OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY: gl_EnableVertexAttribArray(expr, stack, fp)
	case OP_GL_VERTEX_ATTRIB_POINTER: gl_VertexAttribPointer(expr, stack, fp)
	case OP_GL_DRAW_ARRAYS: gl_DrawArrays(expr, stack, fp)
	case OP_GL_GEN_BUFFERS: gl_GenBuffers(expr, stack, fp)
	case OP_GL_BUFFER_DATA: gl_BufferData(expr, stack, fp)
	case OP_GL_GEN_VERTEX_ARRAYS: gl_GenVertexArrays(expr, stack, fp)
	case OP_GL_CREATE_SHADER: gl_CreateShader(expr, stack, fp)
	case OP_GL_STRS: gl_Strs(expr, stack, fp)
	case OP_GL_FREE: gl_Free(expr, stack, fp)
	case OP_GL_SHADER_SOURCE: gl_ShaderSource(expr, stack, fp)
	case OP_GL_COMPILE_SHADER: gl_CompileShader(expr, stack, fp)
	case OP_GL_GET_SHADERIV: gl_GetShaderiv(expr, stack, fp)
	case OP_GL_ATTACH_SHADER: gl_AttachShader(expr, stack, fp)
	case OP_GL_MATRIX_MODE: gl_MatrixMode(expr, stack, fp)
	case OP_GL_ROTATEF: gl_Rotatef(expr, stack, fp)
	case OP_GL_TRANSLATEF: gl_Translatef(expr, stack, fp)
	case OP_GL_LOAD_IDENTITY: gl_LoadIdentity()
	case OP_GL_PUSH_MATRIX: gl_PushMatrix()
	case OP_GL_POP_MATRIX: gl_PopMatrix()
	case OP_GL_ENABLE_CLIENT_STATE: gl_EnableClientState(expr, stack, fp)
	case OP_GL_BIND_TEXTURE: gl_BindTexture(expr, stack, fp)
	case OP_GL_COLOR3F: gl_Color3f(expr, stack, fp)
	case OP_GL_COLOR4F: gl_Color4f(expr, stack, fp)
	case OP_GL_BEGIN: gl_Begin(expr, stack, fp)
	case OP_GL_END: gl_End()
	case OP_GL_NORMAL3F: gl_Normal3f(expr, stack, fp)
	case OP_GL_VERTEX_2F: gl_Vertex2f(expr, stack, fp)
	case OP_GL_VERTEX_3F: gl_Vertex3f(expr, stack, fp)
	case OP_GL_ENABLE: gl_Enable(expr, stack, fp)
	case OP_GL_CLEAR_COLOR: gl_ClearColor(expr, stack, fp)
	case OP_GL_CLEAR_DEPTH: gl_ClearDepth(expr, stack, fp)
	case OP_GL_DEPTH_FUNC: gl_DepthFunc(expr, stack, fp)
	case OP_GL_LIGHTFV: gl_Lightfv(expr, stack, fp)
	case OP_GL_FRUSTUM: gl_Frustum(expr, stack, fp)
	case OP_GL_DISABLE: gl_Disable(expr, stack, fp)
	case OP_GL_HINT: gl_Hint(expr, stack, fp)
	case OP_GL_NEW_TEXTURE: gl_NewTexture(expr, stack, fp)
	case OP_GL_DEPTH_MASK: gl_DepthMask(expr, stack, fp)
	case OP_GL_TEX_ENVI: gl_TexEnvi(expr, stack, fp)
	case OP_GL_BLEND_FUNC: gl_BlendFunc(expr, stack, fp)
	case OP_GL_ORTHO: gl_Ortho(expr, stack, fp)
	case OP_GL_VIEWPORT: gl_Viewport(expr, stack, fp)
	case OP_GL_SCALEF: gl_Scalef(expr, stack, fp)
	case OP_GL_TEX_COORD_2D: gl_TexCoord2d(expr, stack, fp)

		// glfw
	case OP_GLFW_INIT: glfw_Init(expr, stack, fp)
	case OP_GLFW_WINDOW_HINT: glfw_WindowHint(expr, stack, fp)
	case OP_GLFW_CREATE_WINDOW: glfw_CreateWindow(expr, stack, fp)
	case OP_GLFW_MAKE_CONTEXT_CURRENT: glfw_MakeContextCurrent(expr, stack, fp)
	case OP_GLFW_SHOULD_CLOSE: glfw_ShouldClose(expr, stack, fp)
	case OP_GLFW_SET_SHOULD_CLOSE: glfw_SetShouldClose(expr, stack, fp)
	case OP_GLFW_POLL_EVENTS: glfw_PollEvents()
	case OP_GLFW_SWAP_BUFFERS: glfw_SwapBuffers(expr, stack, fp)
	case OP_GLFW_GET_FRAMEBUFFER_SIZE: glfw_GetFramebufferSize(expr, stack, fp)
	case OP_GLFW_SET_KEY_CALLBACK: glfw_SetKeyCallback(expr, stack, fp)
	case OP_GLFW_GET_TIME: glfw_GetTime(expr, stack, fp)
	case OP_GLFW_SET_MOUSE_BUTTON_CALLBACK: glfw_SetMouseButtonCallback(expr, stack, fp)
	case OP_GLFW_SET_CURSOR_POS_CALLBACK: glfw_SetCursorPosCallback(expr, stack, fp)
	case OP_GLFW_GET_CURSOR_POS: glfw_GetCursorPos(expr, stack, fp)
	case OP_GLFW_SET_INPUT_MODE: glfw_SetInputMode(expr, stack, fp)
	}
}

// For the parser. These shouldn't be used in the runtime for performance reasons
var OpNames map[int]string = map[int]string{
	OP_IDENTITY: "identity",
	OP_READ_ARRAY: "read",
	OP_JMP: "jmp",
	OP_DEBUG: "debug",

	OP_BYTE_PRINT: "byte.print",

	OP_BOOL_PRINT: "bool.print",
	OP_BOOL_NOT: "bool.not",
	OP_BOOL_OR: "bool.or",
	OP_BOOL_AND: "bool.and",

	OP_I32_BYTE: "i32.byte",
	OP_I32_STR: "i32.str",
	OP_I32_I32: "i32.i32",
	OP_I32_I64: "i32.i64",
	OP_I32_F32: "i32.f32",
	OP_I32_F64: "i32.f64",
	OP_I32_PRINT: "i32.print",
	OP_I32_ADD: "i32.add",
	OP_I32_SUB: "i32.sub",
	OP_I32_MUL: "i32.mul",
	OP_I32_DIV: "i32.div",
	OP_I32_ABS: "i32.abs",
	OP_I32_POW: "i32.pow",
	OP_I32_GT: "i32.gt",
	OP_I32_GTEQ: "i32.gteq",
	OP_I32_LT: "i32.lt",
	OP_I32_LTEQ: "i32.lteq",
	OP_I32_EQ: "i32.eq",
	OP_I32_UNEQ: "i32.uneq",
	OP_I32_MOD: "i32.mod",
	OP_I32_RAND: "i32.rand",
	OP_I32_BITAND: "i32.bitand",
	OP_I32_BITOR: "i32.bitor",
	OP_I32_BITXOR: "i32.bitxor",
	OP_I32_BITCLEAR: "i32.bitclear",
	OP_I32_BITSHL: "i32.bitshl",
	OP_I32_BITSHR: "i32.bitshr",

	OP_I64_PRINT: "i64.print",
	OP_I64_ADD: "i64.add",
	OP_I64_SUB: "i64.sub",
	OP_I64_MUL: "i64.mul",
	OP_I64_DIV: "i64.div",
	OP_I64_ABS: "i64.abs",
	OP_I64_POW: "i64.pow",
	OP_I64_GT: "i64.gt",
	OP_I64_GTEQ: "i64.gteq",
	OP_I64_LT: "i64.lt",
	OP_I64_LTEQ: "i64.lteq",
	OP_I64_EQ: "i64.eq",
	OP_I64_UNEQ: "i64.uneq",
	OP_I64_MOD: "i64.mod",
	OP_I64_RAND: "i64.rand",
	OP_I64_BITAND: "i64.bitand",
	OP_I64_BITOR: "i64.bitor",
	OP_I64_BITXOR: "i64.bitxor",
	OP_I64_BITCLEAR: "i64.bitclear",
	OP_I64_BITSHL: "i64.bitshl",
	OP_I64_BITSHR: "i64.bitshr",

	OP_F32_PRINT: "f32.print",
	OP_F32_ADD: "f32.add",
	OP_F32_SUB: "f32.sub",
	OP_F32_MUL: "f32.mul",
	OP_F32_DIV: "f32.div",
	OP_F32_ABS: "f32.abs",
	OP_F32_POW: "f32.pow",
	OP_F32_GT: "f32.gt",
	OP_F32_GTEQ: "f32.gteq",
	OP_F32_LT: "f32.lt",
	OP_F32_LTEQ: "f32.lteq",
	OP_F32_EQ: "f32.eq",
	OP_F32_UNEQ: "f32.uneq",
	OP_F32_COS: "f32.cos",
	OP_F32_SIN: "f32.sin",

	OP_F64_PRINT: "f64.print",
	OP_F64_ADD: "f64.add",
	OP_F64_SUB: "f64.sub",
	OP_F64_MUL: "f64.mul",
	OP_F64_DIV: "f64.div",
	OP_F64_ABS: "f64.abs",
	OP_F64_POW: "f64.pow",
	OP_F64_GT: "f64.gt",
	OP_F64_GTEQ: "f64.gteq",
	OP_F64_LT: "f64.lt",
	OP_F64_LTEQ: "f64.lteq",
	OP_F64_EQ: "f64.eq",
	OP_F64_UNEQ: "f64.uneq",
	OP_F64_COS: "f64.cos",
	OP_F64_SIN: "f64.sin",

	OP_STR_PRINT: "str.print",

	OP_TIME_SLEEP: "time.Sleep",
	OP_TIME_UNIX_MILLI: "time.UnixMilli",



	// opengl
	OP_GL_INIT: "gl.Init",
	OP_GL_CREATE_PROGRAM: "gl.CreateProgram",
	OP_GL_LINK_PROGRAM: "gl.LinkProgram",
	OP_GL_CLEAR: "gl.Clear",
	OP_GL_USE_PROGRAM: "gl.UseProgram",
	OP_GL_BIND_BUFFER: "gl.BindBuffer",
	OP_GL_BIND_VERTEX_ARRAY: "gl.BindVertexArray",
	OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY: "gl.EnableVertexAttribArray",
	OP_GL_VERTEX_ATTRIB_POINTER: "gl.VertexAttribPointer",
	OP_GL_DRAW_ARRAYS: "gl.DrawArrays",
	OP_GL_GEN_BUFFERS: "gl.GenBuffers",
	OP_GL_BUFFER_DATA: "gl.BufferData",
	OP_GL_GEN_VERTEX_ARRAYS: "gl.GenVertexArrays",
	OP_GL_CREATE_SHADER: "gl.CreateShader",
	OP_GL_STRS: "gl.Strs",
	OP_GL_FREE: "gl.Free",
	OP_GL_SHADER_SOURCE: "gl.ShaderSource",
	OP_GL_COMPILE_SHADER: "gl.CompileShader",
	OP_GL_GET_SHADERIV: "gl.GetShaderiv",
	OP_GL_ATTACH_SHADER: "gl.AttachShader",
	OP_GL_MATRIX_MODE: "gl.MatrixMode",
	OP_GL_ROTATEF: "gl.Rotatef",
	OP_GL_TRANSLATEF: "gl.Translatef",
	OP_GL_LOAD_IDENTITY: "gl.LoadIdentity",
	OP_GL_PUSH_MATRIX: "gl.PushMatrix",
	OP_GL_POP_MATRIX: "gl.PopMatrix",
	OP_GL_ENABLE_CLIENT_STATE: "gl.EnableClientState",
	OP_GL_BIND_TEXTURE: "gl.BindTexture",
	OP_GL_COLOR3F: "gl.Color3f",
	OP_GL_COLOR4F: "gl.Color4f",
	OP_GL_BEGIN: "gl.Begin",
	OP_GL_END: "gl.End",
	OP_GL_NORMAL3F: "gl.Normal3f",
	OP_GL_VERTEX_2F: "gl.Vertex2f",
	OP_GL_VERTEX_3F: "gl.Vertex3f",
	OP_GL_ENABLE: "gl.Enable",
	OP_GL_CLEAR_COLOR: "gl.ClearColor",
	OP_GL_CLEAR_DEPTH: "gl.ClearDepth",
	OP_GL_DEPTH_FUNC: "gl.DepthFunc",
	OP_GL_LIGHTFV: "gl.Lightfv",
	OP_GL_FRUSTUM: "gl.Frustum",
	OP_GL_DISABLE: "gl.Disable",
	OP_GL_HINT: "gl.Hint",
	OP_GL_NEW_TEXTURE: "gl.NewTexture",
	OP_GL_DEPTH_MASK: "gl.DepthMask",
	OP_GL_TEX_ENVI: "gl.TexEnvi",
	OP_GL_BLEND_FUNC: "gl.BlendFunc",
	OP_GL_ORTHO: "gl.Ortho",
	OP_GL_VIEWPORT: "gl.Viewport",
	OP_GL_SCALEF: "gl.Scalef",
	OP_GL_TEX_COORD_2D: "gl.TexCoord2d",

	// glfw
	OP_GLFW_INIT: "glfw.Init",
	OP_GLFW_WINDOW_HINT: "glfw.WindowHint",
	OP_GLFW_CREATE_WINDOW: "glfw.CreateWindow",
	OP_GLFW_MAKE_CONTEXT_CURRENT: "glfw.MakeContextCurrent",
	OP_GLFW_SHOULD_CLOSE: "glfw.ShouldClose",
	OP_GLFW_SET_SHOULD_CLOSE: "glfw.SetShouldClose",
	OP_GLFW_POLL_EVENTS: "glfw.PollEvents",
	OP_GLFW_SWAP_BUFFERS: "glfw.SwapBuffers",
	OP_GLFW_GET_FRAMEBUFFER_SIZE: "glfw.GetFramebufferSize",
	OP_GLFW_SET_KEY_CALLBACK: "glfw.SetKeyCallback",
	OP_GLFW_GET_TIME: "glfw.GetTime",
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK: "glfw.SetMouseButtonCallback",
	OP_GLFW_SET_CURSOR_POS_CALLBACK: "glfw.SetCursorPosCallback",
	OP_GLFW_GET_CURSOR_POS: "glfw.GetCursorPos",
	OP_GLFW_SET_INPUT_MODE: "glfw.SetInputMode",
}

// For the parser. These shouldn't be used in the runtime for performance reasons
var OpCodes map[string]int = map[string]int{
	"identity": OP_IDENTITY,
	"read": OP_READ_ARRAY,
	"jmp": OP_JMP,
	"debug": OP_DEBUG,

	"byte.print": OP_BYTE_PRINT,

	"bool.print": OP_BOOL_PRINT,
	"bool.not": OP_BOOL_NOT,
	"bool.or": OP_BOOL_OR,
	"bool.and": OP_BOOL_AND,

	"i32.byte" : OP_I32_BYTE,
	"i32.str" : OP_I32_STR,
	"i32.i32" : OP_I32_I32,
	"i32.i64" : OP_I32_I64,
	"i32.f32" : OP_I32_F32,
	"i32.f64" : OP_I32_F64,
	"i32.print": OP_I32_PRINT,
	"i32.add": OP_I32_ADD,
	"i32.sub": OP_I32_SUB,
	"i32.mul": OP_I32_MUL,
	"i32.div": OP_I32_DIV,
	"i32.abs": OP_I32_ABS,
	"i32.pow": OP_I32_POW,
	"i32.gt": OP_I32_GT,
	"i32.gteq": OP_I32_GTEQ,
	"i32.lt": OP_I32_LT,
	"i32.lteq": OP_I32_LTEQ,
	"i32.eq": OP_I32_EQ,
	"i32.uneq": OP_I32_UNEQ,
	"i32.mod": OP_I32_MOD,
	"i32.rand": OP_I32_RAND,
	"i32.bitand": OP_I32_BITAND,
	"i32.bitor": OP_I32_BITOR,
	"i32.bitxor": OP_I32_BITXOR,
	"i32.bitclear": OP_I32_BITCLEAR,
	"i32.bitshl": OP_I32_BITSHL,
	"i32.bitshr": OP_I32_BITSHR,

	"i64.print": OP_I64_PRINT,
	"i64.add": OP_I64_ADD,
	"i64.sub": OP_I64_SUB,
	"i64.mul": OP_I64_MUL,
	"i64.div": OP_I64_DIV,
	"i64.abs": OP_I64_ABS,
	"i64.pow": OP_I64_POW,
	"i64.gt": OP_I64_GT,
	"i64.gteq": OP_I64_GTEQ,
	"i64.lt": OP_I64_LT,
	"i64.lteq": OP_I64_LTEQ,
	"i64.eq": OP_I64_EQ,
	"i64.uneq": OP_I64_UNEQ,
	"i64.mod": OP_I64_MOD,
	"i64.rand": OP_I64_RAND,
	"i64.bitand": OP_I64_BITAND,
	"i64.bitor": OP_I64_BITOR,
	"i64.bitxor": OP_I64_BITXOR,
	"i64.bitclear": OP_I64_BITCLEAR,
	"i64.bitshl": OP_I64_BITSHL,
	"i64.bitshr": OP_I64_BITSHR,

	"f32.print": OP_F32_PRINT,
	"f32.add": OP_F32_ADD,
	"f32.sub": OP_F32_SUB,
	"f32.mul": OP_F32_MUL,
	"f32.div": OP_F32_DIV,
	"f32.abs": OP_F32_ABS,
	"f32.pow": OP_F32_POW,
	"f32.gt": OP_F32_GT,
	"f32.gteq": OP_F32_GTEQ,
	"f32.lt": OP_F32_LT,
	"f32.lteq": OP_F32_LTEQ,
	"f32.eq": OP_F32_EQ,
	"f32.uneq": OP_F32_UNEQ,
	"f32.cos": OP_F32_COS,
	"f32.sin": OP_F32_SIN,

	"f64.print": OP_F64_PRINT,
	"f64.add": OP_F64_ADD,
	"f64.sub": OP_F64_SUB,
	"f64.mul": OP_F64_MUL,
	"f64.div": OP_F64_DIV,
	"f64.abs": OP_F64_ABS,
	"f64.pow": OP_F64_POW,
	"f64.gt": OP_F64_GT,
	"f64.gteq": OP_F64_GTEQ,
	"f64.lt": OP_F64_LT,
	"f64.lteq": OP_F64_LTEQ,
	"f64.eq": OP_F64_EQ,
	"f64.uneq": OP_F64_UNEQ,
	"f64.cos": OP_F64_COS,
	"f64.sin": OP_F64_SIN,






	"str.print": OP_STR_PRINT,

	"time.Sleep": OP_TIME_SLEEP,
	"time.UnixMilli": OP_TIME_UNIX_MILLI,




	// opengl
	"gl.Init": OP_GL_INIT,
	"gl.CreateProgram": OP_GL_CREATE_PROGRAM,
	"gl.LinkProgram": OP_GL_LINK_PROGRAM,
	"gl.Clear": OP_GL_CLEAR,
	"gl.UseProgram": OP_GL_USE_PROGRAM,
	"gl.BindBuffer": OP_GL_BIND_BUFFER,
	"gl.BindVertexArray": OP_GL_BIND_VERTEX_ARRAY,
	"gl.EnableVertexAttribArray": OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY,
	"gl.VertexAttribPointer": OP_GL_VERTEX_ATTRIB_POINTER,
	"gl.DrawArrays": OP_GL_DRAW_ARRAYS,
	"gl.GenBuffers": OP_GL_GEN_BUFFERS,
	"gl.BufferData": OP_GL_BUFFER_DATA,
	"gl.GenVertexArrays": OP_GL_GEN_VERTEX_ARRAYS,
	"gl.CreateShader": OP_GL_CREATE_SHADER,
	"gl.Strs": OP_GL_STRS,
	"gl.Free": OP_GL_FREE,
	"gl.ShaderSource": OP_GL_SHADER_SOURCE,
	"gl.CompileShader": OP_GL_COMPILE_SHADER,
	"gl.GetShaderiv": OP_GL_GET_SHADERIV,
	"gl.AttachShader": OP_GL_ATTACH_SHADER,
	"gl.MatrixMode": OP_GL_MATRIX_MODE,
	"gl.Rotatef": OP_GL_ROTATEF,
	"gl.Translatef": OP_GL_TRANSLATEF,
	"gl.LoadIdentity": OP_GL_LOAD_IDENTITY,
	"gl.PushMatrix": OP_GL_PUSH_MATRIX,
	"gl.PopMatrix": OP_GL_POP_MATRIX,
	"gl.EnableClientState": OP_GL_ENABLE_CLIENT_STATE,
	"gl.BindTexture": OP_GL_BIND_TEXTURE,
	"gl.Color3f": OP_GL_COLOR3F,
	"gl.Color4f": OP_GL_COLOR4F,
	"gl.Begin": OP_GL_BEGIN,
	"gl.End": OP_GL_END,
	"gl.Normal3f": OP_GL_NORMAL3F,
	"gl.Vertex2f": OP_GL_VERTEX_2F,
	"gl.Vertex3f": OP_GL_VERTEX_3F,
	"gl.Enable": OP_GL_ENABLE,
	"gl.ClearColor": OP_GL_CLEAR_COLOR,
	"gl.ClearDepth": OP_GL_CLEAR_DEPTH,
	"gl.DepthFunc": OP_GL_DEPTH_FUNC,
	"gl.Lightfv": OP_GL_LIGHTFV,
	"gl.Frustum": OP_GL_FRUSTUM,
	"gl.Disable": OP_GL_DISABLE,
	"gl.Hint": OP_GL_HINT,
	"gl.NewTexture": OP_GL_NEW_TEXTURE,
	"gl.DepthMask": OP_GL_DEPTH_MASK,
	"gl.TexEnvi": OP_GL_TEX_ENVI,
	"gl.BlendFunc": OP_GL_BLEND_FUNC,
	"gl.Ortho": OP_GL_ORTHO,
	"gl.Viewport": OP_GL_VIEWPORT,
	"gl.Scalef": OP_GL_SCALEF,
	"gl.TexCoord2d": OP_GL_TEX_COORD_2D,
	
	// glfw
	"glfw.Init": OP_GLFW_INIT,
	"glfw.WindowHint": OP_GLFW_WINDOW_HINT,
	"glfw.CreateWindow": OP_GLFW_CREATE_WINDOW,
	"glfw.MakeContextCurrent": OP_GLFW_MAKE_CONTEXT_CURRENT,
	"glfw.ShouldClose": OP_GLFW_SHOULD_CLOSE,
	"glfw.SetShouldClose": OP_GLFW_SET_SHOULD_CLOSE,
	"glfw.PollEvents": OP_GLFW_POLL_EVENTS,
	"glfw.SwapBuffers": OP_GLFW_SWAP_BUFFERS,
	"glfw.GetFramebufferSize": OP_GLFW_GET_FRAMEBUFFER_SIZE,
	"glfw.SetKeyCallback": OP_GLFW_SET_KEY_CALLBACK,
	"glfw.GetTime": OP_GLFW_GET_TIME,
	"glfw.SetMouseButtonCallback": OP_GLFW_SET_MOUSE_BUTTON_CALLBACK,
	"glfw.SetCursorPosCallback": OP_GLFW_SET_CURSOR_POS_CALLBACK,
	"glfw.GetCursorPos": OP_GLFW_GET_CURSOR_POS,
	"glfw.SetInputMode": OP_GLFW_SET_INPUT_MODE,
}

var Natives map[int]*CXFunction = map[int]*CXFunction{
	OP_IDENTITY: MakeNative(OP_IDENTITY, []int{TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_READ_ARRAY: MakeNative(OP_READ_ARRAY, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_JMP: MakeNative(OP_JMP, []int{TYPE_BOOL, TYPE_I32, TYPE_I32}, []int{}),
	OP_DEBUG: MakeNative(OP_DEBUG, []int{}, []int{}),

	OP_BYTE_PRINT: MakeNative(OP_BYTE_PRINT, []int{TYPE_BYTE}, []int{}),

	OP_BOOL_PRINT: MakeNative(OP_BOOL_PRINT, []int{TYPE_BOOL}, []int{}),
	OP_BOOL_NOT: MakeNative(OP_BOOL_NOT, []int{TYPE_BOOL}, []int{TYPE_BOOL}),
	OP_BOOL_OR: MakeNative(OP_BOOL_OR, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),
	OP_BOOL_AND: MakeNative(OP_BOOL_AND, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),

	OP_I32_BYTE: MakeNative(OP_I32_BYTE, []int{TYPE_I32}, []int{TYPE_BYTE}),
	OP_I32_STR: MakeNative(OP_I32_STR, []int{TYPE_I32}, []int{TYPE_STR}),
	OP_I32_I32: MakeNative(OP_I32_I32, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_I64: MakeNative(OP_I32_I64, []int{TYPE_I32}, []int{TYPE_I64}),
	OP_I32_F32: MakeNative(OP_I32_F32, []int{TYPE_I32}, []int{TYPE_F32}),
	OP_I32_F64: MakeNative(OP_I32_F64, []int{TYPE_I32}, []int{TYPE_F64}),
	
	OP_I32_PRINT: MakeNative(OP_I32_PRINT, []int{TYPE_I32}, []int{}),
	OP_I32_ADD: MakeNative(OP_I32_ADD, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_SUB: MakeNative(OP_I32_SUB, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_MUL: MakeNative(OP_I32_MUL, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_DIV: MakeNative(OP_I32_DIV, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_ABS: MakeNative(OP_I32_ABS, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_POW: MakeNative(OP_I32_POW, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_I32_GT: MakeNative(OP_I32_GT, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_GTEQ: MakeNative(OP_I32_GTEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_LT: MakeNative(OP_I32_LT, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_LTEQ: MakeNative(OP_I32_LTEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_EQ: MakeNative(OP_I32_EQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_UNEQ: MakeNative(OP_I32_UNEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_MOD: MakeNative(OP_I32_MOD, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_RAND: MakeNative(OP_I32_RAND, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITAND: MakeNative(OP_I32_BITAND, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITOR: MakeNative(OP_I32_BITOR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITXOR: MakeNative(OP_I32_BITXOR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITCLEAR: MakeNative(OP_I32_BITCLEAR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITSHL: MakeNative(OP_I32_BITSHL, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_BITSHR: MakeNative(OP_I32_BITSHR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),

	OP_I64_PRINT: MakeNative(OP_I64_PRINT, []int{TYPE_I64}, []int{}),
	OP_I64_ADD: MakeNative(OP_I64_ADD, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_SUB: MakeNative(OP_I64_SUB, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_MUL: MakeNative(OP_I64_MUL, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_DIV: MakeNative(OP_I64_DIV, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_ABS: MakeNative(OP_I64_ABS, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_POW: MakeNative(OP_I64_POW, []int{TYPE_I64}, []int{TYPE_I64}),
	OP_I64_GT: MakeNative(OP_I64_GT, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_GTEQ: MakeNative(OP_I64_GTEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_LT: MakeNative(OP_I64_LT, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_LTEQ: MakeNative(OP_I64_LTEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_EQ: MakeNative(OP_I64_EQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_UNEQ: MakeNative(OP_I64_UNEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
	OP_I64_MOD: MakeNative(OP_I64_MOD, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_RAND: MakeNative(OP_I64_RAND, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITAND: MakeNative(OP_I64_BITAND, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITOR: MakeNative(OP_I64_BITOR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITXOR: MakeNative(OP_I64_BITXOR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITCLEAR: MakeNative(OP_I64_BITCLEAR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITSHL: MakeNative(OP_I64_BITSHL, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_BITSHR: MakeNative(OP_I64_BITSHR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	
	OP_F32_PRINT: MakeNative(OP_F32_PRINT, []int{TYPE_F32}, []int{}),
	OP_F32_ADD: MakeNative(OP_F32_ADD, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_SUB: MakeNative(OP_F32_SUB, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_MUL: MakeNative(OP_F32_MUL, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_DIV: MakeNative(OP_F32_DIV, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
	OP_F32_ABS: MakeNative(OP_F32_ABS, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_POW: MakeNative(OP_F32_POW, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_GT: MakeNative(OP_F32_GT, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_GTEQ: MakeNative(OP_F32_GTEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_LT: MakeNative(OP_F32_LT, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_LTEQ: MakeNative(OP_F32_LTEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_EQ: MakeNative(OP_F32_EQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_UNEQ: MakeNative(OP_F32_UNEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
	OP_F32_COS: MakeNative(OP_F32_COS, []int{TYPE_F32}, []int{TYPE_F32}),
	OP_F32_SIN: MakeNative(OP_F32_SIN, []int{TYPE_F32}, []int{TYPE_F32}),

	OP_F64_PRINT: MakeNative(OP_F64_PRINT, []int{TYPE_F64}, []int{}),
	OP_F64_ADD: MakeNative(OP_F64_ADD, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_SUB: MakeNative(OP_F64_SUB, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_MUL: MakeNative(OP_F64_MUL, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_DIV: MakeNative(OP_F64_DIV, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
	OP_F64_ABS: MakeNative(OP_F64_ABS, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_POW: MakeNative(OP_F64_POW, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_GT: MakeNative(OP_F64_GT, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_GTEQ: MakeNative(OP_F64_GTEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_LT: MakeNative(OP_F64_LT, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_LTEQ: MakeNative(OP_F64_LTEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_EQ: MakeNative(OP_F64_EQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_UNEQ: MakeNative(OP_F64_UNEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
	OP_F64_COS: MakeNative(OP_F64_COS, []int{TYPE_F64}, []int{TYPE_F64}),
	OP_F64_SIN: MakeNative(OP_F64_SIN, []int{TYPE_F64}, []int{TYPE_F64}),
	
	OP_STR_PRINT: MakeNative(OP_STR_PRINT, []int{TYPE_STR}, []int{}),

	OP_TIME_SLEEP: MakeNative(OP_TIME_SLEEP, []int{TYPE_I32}, []int{}),
	OP_TIME_UNIX_MILLI: MakeNative(OP_TIME_UNIX_MILLI, []int{}, []int{TYPE_I64}),

	// opengl
	OP_GL_INIT: MakeNative(OP_GL_INIT, []int{}, []int{}),
	OP_GL_CREATE_PROGRAM: MakeNative(OP_GL_CREATE_PROGRAM, []int{}, []int{TYPE_I32}),
	OP_GL_LINK_PROGRAM: MakeNative(OP_GL_LINK_PROGRAM, []int{TYPE_I32}, []int{}),
	OP_GL_CLEAR: MakeNative(OP_GL_CLEAR, []int{TYPE_I32}, []int{}),
	OP_GL_USE_PROGRAM: MakeNative(OP_GL_USE_PROGRAM, []int{TYPE_I32}, []int{}),
	OP_GL_BIND_BUFFER: MakeNative(OP_GL_BIND_BUFFER, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_BIND_VERTEX_ARRAY: MakeNative(OP_GL_BIND_VERTEX_ARRAY, []int{TYPE_I32}, []int{}),
	OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY: MakeNative(OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY, []int{TYPE_I32}, []int{}),
	OP_GL_VERTEX_ATTRIB_POINTER: MakeNative(OP_GL_VERTEX_ATTRIB_POINTER, []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_BOOL, TYPE_I32}, []int{}),
	OP_GL_DRAW_ARRAYS: MakeNative(OP_GL_DRAW_ARRAYS, []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_GEN_BUFFERS: MakeNative(OP_GL_GEN_BUFFERS, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_GL_BUFFER_DATA: MakeNative(OP_GL_BUFFER_DATA, []int{TYPE_I32, TYPE_I32, TYPE_F32, TYPE_I32}, []int{}),
	OP_GL_GEN_VERTEX_ARRAYS: MakeNative(OP_GL_GEN_VERTEX_ARRAYS, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_GL_CREATE_SHADER: MakeNative(OP_GL_CREATE_SHADER, []int{TYPE_I32}, []int{TYPE_I32}),
	OP_GL_STRS: MakeNative(OP_GL_STRS, []int{TYPE_STR, TYPE_STR}, []int{}),
	OP_GL_FREE: MakeNative(OP_GL_FREE, []int{TYPE_STR}, []int{}),
	OP_GL_SHADER_SOURCE: MakeNative(OP_GL_SHADER_SOURCE, []int{TYPE_I32, TYPE_I32, TYPE_STR}, []int{}),
	OP_GL_COMPILE_SHADER: MakeNative(OP_GL_COMPILE_SHADER, []int{TYPE_I32}, []int{}),
	OP_GL_GET_SHADERIV: MakeNative(OP_GL_GET_SHADERIV, []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_ATTACH_SHADER: MakeNative(OP_GL_ATTACH_SHADER, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_MATRIX_MODE: MakeNative(OP_GL_MATRIX_MODE, []int{TYPE_I32}, []int{}),
	OP_GL_ROTATEF: MakeNative(OP_GL_ROTATEF, []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_TRANSLATEF: MakeNative(OP_GL_TRANSLATEF, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_LOAD_IDENTITY: MakeNative(OP_GL_LOAD_IDENTITY, []int{}, []int{}),
	OP_GL_PUSH_MATRIX: MakeNative(OP_GL_PUSH_MATRIX, []int{}, []int{}),
	OP_GL_POP_MATRIX: MakeNative(OP_GL_POP_MATRIX, []int{}, []int{}),
	OP_GL_ENABLE_CLIENT_STATE: MakeNative(OP_GL_ENABLE_CLIENT_STATE, []int{TYPE_I32}, []int{}),
	OP_GL_BIND_TEXTURE: MakeNative(OP_GL_BIND_TEXTURE, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_COLOR3F: MakeNative(OP_GL_COLOR3F, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_COLOR4F: MakeNative(OP_GL_COLOR4F, []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_BEGIN: MakeNative(OP_GL_BEGIN, []int{TYPE_I32}, []int{}),
	OP_GL_END: MakeNative(OP_GL_END, []int{}, []int{}),
	OP_GL_NORMAL3F: MakeNative(OP_GL_NORMAL3F, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),

	OP_GL_VERTEX_2F: MakeNative(OP_GL_VERTEX_2F, []int{TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_VERTEX_3F: MakeNative(OP_GL_VERTEX_3F, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	
	OP_GL_ENABLE: MakeNative(OP_GL_ENABLE, []int{TYPE_I32}, []int{}),
	OP_GL_CLEAR_COLOR: MakeNative(OP_GL_CLEAR_COLOR, []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_CLEAR_DEPTH: MakeNative(OP_GL_CLEAR_DEPTH, []int{TYPE_F64}, []int{}),
	OP_GL_DEPTH_FUNC: MakeNative(OP_GL_DEPTH_FUNC, []int{TYPE_I32}, []int{}),
	OP_GL_LIGHTFV: MakeNative(OP_GL_LIGHTFV, []int{TYPE_I32, TYPE_I32, TYPE_F32}, []int{}),
	OP_GL_FRUSTUM: MakeNative(OP_GL_FRUSTUM, []int{TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64}, []int{}),
	OP_GL_DISABLE: MakeNative(OP_GL_DISABLE, []int{TYPE_I32}, []int{}),
	OP_GL_HINT: MakeNative(OP_GL_HINT, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_NEW_TEXTURE: MakeNative(OP_GL_NEW_TEXTURE, []int{TYPE_STR}, []int{TYPE_I32}),
	OP_GL_DEPTH_MASK: MakeNative(OP_GL_DEPTH_MASK, []int{TYPE_BOOL}, []int{}),
	OP_GL_TEX_ENVI: MakeNative(OP_GL_TEX_ENVI, []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_BLEND_FUNC: MakeNative(OP_GL_BLEND_FUNC, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_ORTHO: MakeNative(OP_GL_ORTHO, []int{TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64}, []int{}),
	OP_GL_VIEWPORT: MakeNative(OP_GL_VIEWPORT, []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{}),
	OP_GL_SCALEF: MakeNative(OP_GL_SCALEF, []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{}),
	OP_GL_TEX_COORD_2D: MakeNative(OP_GL_TEX_COORD_2D, []int{TYPE_F32, TYPE_F32}, []int{}),
	
	// glfw
	OP_GLFW_INIT: MakeNative(OP_GLFW_INIT, []int{}, []int{}),
	OP_GLFW_WINDOW_HINT: MakeNative(OP_GLFW_WINDOW_HINT, []int{TYPE_I32, TYPE_I32}, []int{}),
	OP_GLFW_CREATE_WINDOW: MakeNative(OP_GLFW_CREATE_WINDOW, []int{TYPE_STR, TYPE_I32, TYPE_I32, TYPE_STR}, []int{}),
	OP_GLFW_MAKE_CONTEXT_CURRENT: MakeNative(OP_GLFW_MAKE_CONTEXT_CURRENT, []int{TYPE_STR}, []int{}),
	OP_GLFW_SHOULD_CLOSE: MakeNative(OP_GLFW_SHOULD_CLOSE, []int{TYPE_STR}, []int{TYPE_BOOL}),
	OP_GLFW_SET_SHOULD_CLOSE: MakeNative(OP_GLFW_SET_SHOULD_CLOSE, []int{TYPE_STR, TYPE_BOOL}, []int{TYPE_BOOL}),
	OP_GLFW_POLL_EVENTS: MakeNative(OP_GLFW_POLL_EVENTS, []int{}, []int{}),
	OP_GLFW_SWAP_BUFFERS: MakeNative(OP_GLFW_SWAP_BUFFERS, []int{TYPE_STR}, []int{}),
	OP_GLFW_GET_FRAMEBUFFER_SIZE: MakeNative(OP_GLFW_GET_FRAMEBUFFER_SIZE, []int{TYPE_STR}, []int{TYPE_I32, TYPE_I32}),
	OP_GLFW_SET_KEY_CALLBACK: MakeNative(OP_GLFW_SET_KEY_CALLBACK, []int{TYPE_STR, TYPE_STR}, []int{}),
	OP_GLFW_GET_TIME: MakeNative(OP_GLFW_GET_TIME, []int{}, []int{TYPE_F64}),
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK: MakeNative(OP_GLFW_SET_MOUSE_BUTTON_CALLBACK, []int{TYPE_STR, TYPE_STR}, []int{}),
	OP_GLFW_SET_CURSOR_POS_CALLBACK: MakeNative(OP_GLFW_SET_CURSOR_POS_CALLBACK, []int{TYPE_STR, TYPE_STR}, []int{}),
	OP_GLFW_GET_CURSOR_POS: MakeNative(OP_GLFW_GET_CURSOR_POS, []int{TYPE_STR}, []int{TYPE_F64, TYPE_F64}),
	OP_GLFW_SET_INPUT_MODE: MakeNative(OP_GLFW_SET_INPUT_MODE, []int{TYPE_STR, TYPE_I32, TYPE_I32}, []int{}),
}
