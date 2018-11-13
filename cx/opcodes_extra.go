// +build extra full

package base

import (
	
)

const (
	// opengl
	OP_GL_INIT = iota + END_OF_BASE_OPS
	OP_GL_GET_ERROR
	OP_GL_CULL_FACE
	OP_GL_CREATE_PROGRAM
	OP_GL_DELETE_PROGRAM
	OP_GL_LINK_PROGRAM
	OP_GL_CLEAR
	OP_GL_USE_PROGRAM
	OP_GL_BIND_BUFFER
	OP_GL_BIND_VERTEX_ARRAY
	OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY
	OP_GL_VERTEX_ATTRIB_POINTER
	OP_GL_VERTEX_ATTRIB_POINTER_I32
	OP_GL_DRAW_ARRAYS
	OP_GL_GEN_BUFFERS
	OP_GL_DELETE_BUFFERS
	OP_GL_BUFFER_DATA
	OP_GL_BUFFER_SUB_DATA
	OP_GL_GEN_VERTEX_ARRAYS
	OP_GL_DELETE_VERTEX_ARRAYS
	OP_GL_CREATE_SHADER
	OP_GL_DETACH_SHADER
	OP_GL_DELETE_SHADER
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
	OP_GL_ACTIVE_TEXTURE
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
	OP_GL_TEX_COORD_2F

	/* gl_1_0 */
	OP_GL_SCISSOR
	OP_GL_TEX_IMAGE_2D
	OP_GL_TEX_PARAMETERI
	OP_GL_GET_TEX_LEVEL_PARAMETERIV

	/* gl_1_1 */
	OP_GL_BIND_TEXTURE
	OP_GL_GEN_TEXTURES
	OP_GL_DELETE_TEXTURES

	/* gl_2_0 */
	OP_GL_BIND_ATTRIB_LOCATION
	OP_GL_GET_ATTRIB_LOCATION
	OP_GL_GET_UNIFORM_LOCATION
	OP_GL_UNIFORM_1F
	OP_GL_UNIFORM_1I

	/* gl_3_0 */
	OP_GL_BIND_RENDERBUFFER
	OP_GL_DELETE_RENDERBUFFERS
	OP_GL_GEN_RENDERBUFFERS
	OP_GL_RENDERBUFFER_STORAGE
	OP_GL_BIND_FRAMEBUFFER
	OP_GL_DELETE_FRAMEBUFFERS
	OP_GL_GEN_FRAMEBUFFERS
	OP_GL_CHECK_FRAMEBUFFER_STATUS
	OP_GL_FRAMEBUFFER_TEXTURE_2D
	OP_GL_FRAMEBUFFER_RENDERBUFFER

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
	OP_GLFW_SWAP_INTERVAL
	OP_GLFW_SET_KEY_CALLBACK
	OP_GLFW_GET_TIME
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK
	OP_GLFW_SET_CURSOR_POS_CALLBACK
	OP_GLFW_GET_CURSOR_POS
	OP_GLFW_SET_INPUT_MODE
	OP_GLFW_SET_WINDOW_POS
	OP_GLFW_GET_KEY

	// gltext
	OP_GLTEXT_LOAD_TRUE_TYPE
	OP_GLTEXT_PRINTF
	OP_GLTEXT_METRICS
	OP_GLTEXT_TEXTURE
	OP_GLTEXT_NEXT_RUNE
	OP_GLTEXT_GLYPH_BOUNDS
)

func init () {
	// gl_0.0
	AddOpCode(OP_GL_INIT, "gl.Init", []int{}, []int{})
	AddOpCode(OP_GL_GET_ERROR, "gl.GetError", []int{}, []int{TYPE_I32})
	AddOpCode(OP_GL_CULL_FACE, "gl.CullFace", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_CREATE_PROGRAM, "gl.CreateProgram", []int{}, []int{TYPE_I32})
	AddOpCode(OP_GL_DELETE_PROGRAM, "gl.DeleteProgram", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_LINK_PROGRAM, "gl.LinkProgram", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_CLEAR, "gl.Clear", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_USE_PROGRAM, "gl.UseProgram", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_BIND_BUFFER, "gl.BindBuffer", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_BIND_VERTEX_ARRAY, "gl.BindVertexArray", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY, "gl.EnableVertexAttribArray", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_VERTEX_ATTRIB_POINTER, "gl.VertexAttribPointer", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_BOOL, TYPE_I32}, []int{})
	AddOpCode(OP_GL_VERTEX_ATTRIB_POINTER_I32, "gl.VertexAttribPointerI32", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_BOOL, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_DRAW_ARRAYS, "gl.DrawArrays", []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_GEN_BUFFERS, "gl.GenBuffers", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_GL_DELETE_BUFFERS, "gl.DeleteBuffers", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_BUFFER_DATA, "gl.BufferData", []int{TYPE_I32, TYPE_I32, TYPE_F32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_BUFFER_SUB_DATA, "gl.BufferSubData", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_GEN_VERTEX_ARRAYS, "gl.GenVertexArrays", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_GL_DELETE_VERTEX_ARRAYS, "gl.DeleteVertexArrays", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_CREATE_SHADER, "gl.CreateShader", []int{TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_GL_DETACH_SHADER, "gl.DetachShader", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_DELETE_SHADER, "gl.DeleteShader", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_STRS, "gl.Strs", []int{TYPE_STR, TYPE_STR}, []int{})
	AddOpCode(OP_GL_FREE, "gl.Free", []int{TYPE_STR}, []int{})
	AddOpCode(OP_GL_SHADER_SOURCE, "gl.ShaderSource", []int{TYPE_I32, TYPE_I32, TYPE_STR}, []int{})
	AddOpCode(OP_GL_COMPILE_SHADER, "gl.CompileShader", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_GET_SHADERIV, "gl.GetShaderiv", []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_ATTACH_SHADER, "gl.AttachShader", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_MATRIX_MODE, "gl.MatrixMode", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_ROTATEF, "gl.Rotatef", []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_TRANSLATEF, "gl.Translatef", []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_LOAD_IDENTITY, "gl.LoadIdentity", []int{}, []int{})
	AddOpCode(OP_GL_PUSH_MATRIX, "gl.PushMatrix", []int{}, []int{})
	AddOpCode(OP_GL_POP_MATRIX, "gl.PopMatrix", []int{}, []int{})
	AddOpCode(OP_GL_ENABLE_CLIENT_STATE, "gl.EnableClientState", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_ACTIVE_TEXTURE, "gl.ActiveTexture", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_COLOR3F, "gl.Color3f", []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_COLOR4F, "gl.Color4f", []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_BEGIN, "gl.Begin", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_END, "gl.End", []int{}, []int{})
	AddOpCode(OP_GL_NORMAL3F, "gl.Normal3f", []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_VERTEX_2F, "gl.Vertex2f", []int{TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_VERTEX_3F, "gl.Vertex3f", []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_ENABLE, "gl.Enable", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_CLEAR_COLOR, "gl.ClearColor", []int{TYPE_F32, TYPE_F32, TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_CLEAR_DEPTH, "gl.ClearDepth", []int{TYPE_F64}, []int{})
	AddOpCode(OP_GL_DEPTH_FUNC, "gl.DepthFunc", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_LIGHTFV, "gl.Lightfv", []int{TYPE_I32, TYPE_I32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_FRUSTUM, "gl.Frustum", []int{TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64}, []int{})
	AddOpCode(OP_GL_DISABLE, "gl.Disable", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GL_HINT, "gl.Hint", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_NEW_TEXTURE, "gl.NewTexture", []int{TYPE_STR}, []int{TYPE_I32})
	AddOpCode(OP_GL_DEPTH_MASK, "gl.DepthMask", []int{TYPE_BOOL}, []int{})
	AddOpCode(OP_GL_TEX_ENVI, "gl.TexEnvi", []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_BLEND_FUNC, "gl.BlendFunc", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_ORTHO, "gl.Ortho", []int{TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64, TYPE_F64}, []int{})
	AddOpCode(OP_GL_VIEWPORT, "gl.Viewport", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_SCALEF, "gl.Scalef", []int{TYPE_F32, TYPE_F32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_TEX_COORD_2D, "gl.TexCoord2d", []int{TYPE_F64, TYPE_F64}, []int{})
	AddOpCode(OP_GL_TEX_COORD_2F, "gl.TexCoord2f", []int{TYPE_F32, TYPE_F32}, []int{})

	// gl_1_0
	AddOpCode(OP_GL_SCISSOR, "gl.Scissor", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_TEX_IMAGE_2D, "gl.TexImage2D", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_TEX_PARAMETERI, "gl.TexParameteri", []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_GET_TEX_LEVEL_PARAMETERIV, "gl.GetTexLevelParameteriv", []int{TYPE_I32, TYPE_I32, TYPE_I32}, []int{TYPE_I32})

	// gl_1_1
	AddOpCode(OP_GL_BIND_TEXTURE, "gl.BindTexture", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_GEN_TEXTURES, "gl.GenTextures", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_GL_DELETE_TEXTURES, "gl.DeleteTextures", []int{TYPE_I32, TYPE_I32}, []int{})

	//gl_2_0
	AddOpCode(OP_GL_BIND_ATTRIB_LOCATION, "gl.BindAttribLocation", []int{TYPE_I32, TYPE_I32, TYPE_STR}, []int{})
	AddOpCode(OP_GL_GET_ATTRIB_LOCATION, "gl.GetAttribLocation", []int{TYPE_I32, TYPE_STR}, []int{TYPE_I32})
	AddOpCode(OP_GL_GET_UNIFORM_LOCATION, "gl.GetUniformLocation", []int{TYPE_I32, TYPE_STR}, []int{TYPE_I32})
	AddOpCode(OP_GL_UNIFORM_1F, "gl.Uniform1f", []int{TYPE_I32, TYPE_F32}, []int{})
	AddOpCode(OP_GL_UNIFORM_1I, "gl.Uniform1i", []int{TYPE_I32, TYPE_I32}, []int{})

	// gl_3_0
	AddOpCode(OP_GL_BIND_RENDERBUFFER, "gl.BindRenderbuffer", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_DELETE_RENDERBUFFERS, "gl.DeleteRenderbuffers", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_GEN_RENDERBUFFERS, "gl.GenRenderbuffers", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_GL_RENDERBUFFER_STORAGE, "gl.RenderbufferStorage", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_BIND_FRAMEBUFFER, "gl.BindFramebuffer", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_DELETE_FRAMEBUFFERS, "gl.DeleteFramebuffers", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_GEN_FRAMEBUFFERS, "gl.GenFramebuffers", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_GL_CHECK_FRAMEBUFFER_STATUS, "gl.CheckFramebufferStatus", []int{TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_GL_FRAMEBUFFER_TEXTURE_2D, "gl.FramebufferTexture2D", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GL_FRAMEBUFFER_RENDERBUFFER, "gl.FramebufferRenderbuffer", []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{})

	// glfw
	AddOpCode(OP_GLFW_INIT, "glfw.Init", []int{}, []int{})
	AddOpCode(OP_GLFW_WINDOW_HINT, "glfw.WindowHint", []int{TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GLFW_CREATE_WINDOW, "glfw.CreateWindow", []int{TYPE_STR, TYPE_I32, TYPE_I32, TYPE_STR}, []int{})
	AddOpCode(OP_GLFW_MAKE_CONTEXT_CURRENT, "glfw.MakeContextCurrent", []int{TYPE_STR}, []int{})
	AddOpCode(OP_GLFW_SHOULD_CLOSE, "glfw.ShouldClose", []int{TYPE_STR}, []int{TYPE_BOOL})
	AddOpCode(OP_GLFW_SET_SHOULD_CLOSE, "glfw.SetShouldClose", []int{TYPE_STR, TYPE_BOOL}, []int{})
	AddOpCode(OP_GLFW_POLL_EVENTS, "glfw.PollEvents", []int{}, []int{})
	AddOpCode(OP_GLFW_SWAP_BUFFERS, "glfw.SwapBuffers", []int{TYPE_STR}, []int{})
	AddOpCode(OP_GLFW_GET_FRAMEBUFFER_SIZE, "glfw.GetFramebufferSize", []int{TYPE_STR}, []int{TYPE_I32, TYPE_I32})
	AddOpCode(OP_GLFW_SWAP_INTERVAL, "glfw.SwapInterval", []int{TYPE_I32}, []int{})
	AddOpCode(OP_GLFW_SET_KEY_CALLBACK, "glfw.SetKeyCallback", []int{TYPE_STR, TYPE_STR}, []int{})
	AddOpCode(OP_GLFW_GET_TIME, "glfw.GetTime", []int{}, []int{TYPE_F64})
	AddOpCode(OP_GLFW_SET_MOUSE_BUTTON_CALLBACK, "glfw.SetMouseButtonCallback", []int{TYPE_STR, TYPE_STR}, []int{})
	AddOpCode(OP_GLFW_SET_CURSOR_POS_CALLBACK, "glfw.SetCursorPosCallback", []int{TYPE_STR, TYPE_STR}, []int{})
	AddOpCode(OP_GLFW_GET_CURSOR_POS, "glfw.GetCursorPos", []int{TYPE_STR}, []int{TYPE_F64, TYPE_F64})
	AddOpCode(OP_GLFW_SET_INPUT_MODE, "glfw.SetInputMode", []int{TYPE_STR, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GLFW_SET_WINDOW_POS, "glfw.SetWindowPos", []int{TYPE_STR, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GLFW_GET_KEY, "glfw.GetKey", []int{TYPE_STR, TYPE_I32}, []int{TYPE_I32})

	// gltext
	AddOpCode(OP_GLTEXT_LOAD_TRUE_TYPE, "gltext.LoadTrueType", []int{TYPE_STR, TYPE_STR, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32}, []int{})
	AddOpCode(OP_GLTEXT_PRINTF, "gltext.Printf", []int{TYPE_STR, TYPE_F32, TYPE_F32, TYPE_STR}, []int{})
	AddOpCode(OP_GLTEXT_METRICS, "gltext.Metrics", []int{TYPE_STR, TYPE_STR}, []int{TYPE_I32, TYPE_I32})
	AddOpCode(OP_GLTEXT_TEXTURE, "gltext.Texture", []int{TYPE_STR}, []int{TYPE_I32})
	AddOpCode(OP_GLTEXT_NEXT_RUNE, "gltext.NextRune", []int{TYPE_STR, TYPE_STR, TYPE_I32}, []int{TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32, TYPE_I32})
	AddOpCode(OP_GLTEXT_GLYPH_BOUNDS, "gltext.GlyphBounds", []int{}, []int{TYPE_I32, TYPE_I32})

	// exec
	execNative = func (prgrm *CXProgram) {
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

			// opengl
		case OP_GL_INIT:
			op_gl_Init()
		case OP_GL_GET_ERROR:
			op_gl_GetError(expr, fp)
		case OP_GL_CULL_FACE:
			op_gl_CullFace(expr, fp)
		case OP_GL_CREATE_PROGRAM:
			op_gl_CreateProgram(expr, fp)
		case OP_GL_DELETE_PROGRAM:
			op_gl_DeleteProgram(expr, fp)
		case OP_GL_LINK_PROGRAM:
			op_gl_LinkProgram(expr, fp)
		case OP_GL_CLEAR:
			op_gl_Clear(expr, fp)
		case OP_GL_USE_PROGRAM:
			op_gl_UseProgram(expr, fp)
		case OP_GL_BIND_BUFFER:
			op_gl_BindBuffer(expr, fp)
		case OP_GL_BIND_VERTEX_ARRAY:
			op_gl_BindVertexArray(expr, fp)
		case OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY:
			op_gl_EnableVertexAttribArray(expr, fp)
		case OP_GL_VERTEX_ATTRIB_POINTER:
			op_gl_VertexAttribPointer(expr, fp)
		case OP_GL_VERTEX_ATTRIB_POINTER_I32:
			op_gl_VertexAttribPointerI32(expr, fp)
		case OP_GL_DRAW_ARRAYS:
			op_gl_DrawArrays(expr, fp)
		case OP_GL_GEN_BUFFERS:
			op_gl_GenBuffers(expr, fp)
		case OP_GL_DELETE_BUFFERS:
			op_gl_DeleteBuffers(expr, fp)
		case OP_GL_BUFFER_DATA:
			op_gl_BufferData(expr, fp)
		case OP_GL_BUFFER_SUB_DATA:
			op_gl_BufferSubData(expr, fp)
		case OP_GL_GEN_VERTEX_ARRAYS:
			op_gl_GenVertexArrays(expr, fp)
		case OP_GL_DELETE_VERTEX_ARRAYS:
			op_gl_DeleteVertexArrays(expr, fp)
		case OP_GL_CREATE_SHADER:
			op_gl_CreateShader(expr, fp)
		case OP_GL_DETACH_SHADER:
			op_gl_DetachShader(expr, fp)
		case OP_GL_DELETE_SHADER:
			op_gl_DeleteShader(expr, fp)
		case OP_GL_STRS:
			op_gl_Strs(expr, fp)
		case OP_GL_FREE:
			op_gl_Free(expr, fp)
		case OP_GL_SHADER_SOURCE:
			op_gl_ShaderSource(expr, fp)
		case OP_GL_COMPILE_SHADER:
			op_gl_CompileShader(expr, fp)
		case OP_GL_GET_SHADERIV:
			op_gl_GetShaderiv(expr, fp)
		case OP_GL_ATTACH_SHADER:
			op_gl_AttachShader(expr, fp)
		case OP_GL_MATRIX_MODE:
			op_gl_MatrixMode(expr, fp)
		case OP_GL_ROTATEF:
			op_gl_Rotatef(expr, fp)
		case OP_GL_TRANSLATEF:
			op_gl_Translatef(expr, fp)
		case OP_GL_LOAD_IDENTITY:
			op_gl_LoadIdentity()
		case OP_GL_PUSH_MATRIX:
			op_gl_PushMatrix()
		case OP_GL_POP_MATRIX:
			op_gl_PopMatrix()
		case OP_GL_ENABLE_CLIENT_STATE:
			op_gl_EnableClientState(expr, fp)
		case OP_GL_ACTIVE_TEXTURE:
			op_gl_ActiveTexture(expr, fp)
		case OP_GL_COLOR3F:
			op_gl_Color3f(expr, fp)
		case OP_GL_COLOR4F:
			op_gl_Color4f(expr, fp)
		case OP_GL_BEGIN:
			op_gl_Begin(expr, fp)
		case OP_GL_END:
			op_gl_End()
		case OP_GL_NORMAL3F:
			op_gl_Normal3f(expr, fp)
		case OP_GL_VERTEX_2F:
			op_gl_Vertex2f(expr, fp)
		case OP_GL_VERTEX_3F:
			op_gl_Vertex3f(expr, fp)
		case OP_GL_ENABLE:
			op_gl_Enable(expr, fp)
		case OP_GL_CLEAR_COLOR:
			op_gl_ClearColor(expr, fp)
		case OP_GL_CLEAR_DEPTH:
			op_gl_ClearDepth(expr, fp)
		case OP_GL_DEPTH_FUNC:
			op_gl_DepthFunc(expr, fp)
		case OP_GL_LIGHTFV:
			op_gl_Lightfv(expr, fp)
		case OP_GL_FRUSTUM:
			op_gl_Frustum(expr, fp)
		case OP_GL_DISABLE:
			op_gl_Disable(expr, fp)
		case OP_GL_HINT:
			op_gl_Hint(expr, fp)
		case OP_GL_NEW_TEXTURE:
			op_gl_NewTexture(expr, fp)
		case OP_GL_DEPTH_MASK:
			op_gl_DepthMask(expr, fp)
		case OP_GL_TEX_ENVI:
			op_gl_TexEnvi(expr, fp)
		case OP_GL_BLEND_FUNC:
			op_gl_BlendFunc(expr, fp)
		case OP_GL_ORTHO:
			op_gl_Ortho(expr, fp)
		case OP_GL_VIEWPORT:
			op_gl_Viewport(expr, fp)
		case OP_GL_SCALEF:
			op_gl_Scalef(expr, fp)
		case OP_GL_TEX_COORD_2D:
			op_gl_TexCoord2d(expr, fp)
		case OP_GL_TEX_COORD_2F:
			op_gl_TexCoord2f(expr, fp)

			/* gl_1_0 */
		case OP_GL_SCISSOR:
			op_gl_Scissor(expr, fp)
		case OP_GL_TEX_IMAGE_2D:
			op_gl_TexImage2D(expr, fp)
		case OP_GL_TEX_PARAMETERI:
			op_gl_TexParameteri(expr, fp)
		case OP_GL_GET_TEX_LEVEL_PARAMETERIV:
			op_gl_GetTexLevelParameteriv(expr, fp)

			/* gl_1_1 */
		case OP_GL_BIND_TEXTURE:
			op_gl_BindTexture(expr, fp)
		case OP_GL_GEN_TEXTURES:
			op_gl_GenTextures(expr, fp)
		case OP_GL_DELETE_TEXTURES:
			op_gl_DeleteTextures(expr, fp)

			/* gl_2_0 */
		case OP_GL_BIND_ATTRIB_LOCATION:
			op_gl_BindAttribLocation(expr, fp)
		case OP_GL_GET_ATTRIB_LOCATION:
			op_gl_GetAttribLocation(expr, fp)
		case OP_GL_GET_UNIFORM_LOCATION:
			op_gl_GetUniformLocation(expr, fp)
		case OP_GL_UNIFORM_1F:
			op_gl_Uniform1f(expr, fp)
		case OP_GL_UNIFORM_1I:
			op_gl_Uniform1i(expr, fp)

			/* gl_3_0 */
		case OP_GL_BIND_RENDERBUFFER:
			op_gl_BindRenderbuffer(expr, fp)
		case OP_GL_DELETE_RENDERBUFFERS:
			op_gl_DeleteRenderbuffers(expr, fp)
		case OP_GL_GEN_RENDERBUFFERS:
			op_gl_GenRenderbuffers(expr, fp)
		case OP_GL_RENDERBUFFER_STORAGE:
			op_gl_RenderbufferStorage(expr, fp)
		case OP_GL_BIND_FRAMEBUFFER:
			op_gl_BindFramebuffer(expr, fp)
		case OP_GL_DELETE_FRAMEBUFFERS:
			op_gl_DeleteFramebuffers(expr, fp)
		case OP_GL_GEN_FRAMEBUFFERS:
			op_gl_GenFramebuffers(expr, fp)
		case OP_GL_CHECK_FRAMEBUFFER_STATUS:
			op_gl_CheckFramebufferStatus(expr, fp)
		case OP_GL_FRAMEBUFFER_TEXTURE_2D:
			op_gl_FramebufferTexture2D(expr, fp)
		case OP_GL_FRAMEBUFFER_RENDERBUFFER:
			op_gl_FramebufferRenderbuffer(expr, fp)

			// glfw
		case OP_GLFW_INIT:
			op_glfw_Init(expr, fp)
		case OP_GLFW_WINDOW_HINT:
			op_glfw_WindowHint(expr, fp)
		case OP_GLFW_CREATE_WINDOW:
			op_glfw_CreateWindow(expr, fp)
		case OP_GLFW_MAKE_CONTEXT_CURRENT:
			op_glfw_MakeContextCurrent(expr, fp)
		case OP_GLFW_SHOULD_CLOSE:
			op_glfw_ShouldClose(expr, fp)
		case OP_GLFW_SET_SHOULD_CLOSE:
			op_glfw_SetShouldClose(expr, fp)
		case OP_GLFW_POLL_EVENTS:
			op_glfw_PollEvents()
		case OP_GLFW_SWAP_BUFFERS:
			op_glfw_SwapBuffers(expr, fp)
		case OP_GLFW_GET_FRAMEBUFFER_SIZE:
			op_glfw_GetFramebufferSize(expr, fp)
		case OP_GLFW_SWAP_INTERVAL:
			op_glfw_SwapInterval(expr, fp)
		case OP_GLFW_SET_KEY_CALLBACK:
			op_glfw_SetKeyCallback(expr, fp)
		case OP_GLFW_GET_TIME:
			op_glfw_GetTime(expr, fp)
		case OP_GLFW_SET_MOUSE_BUTTON_CALLBACK:
			op_glfw_SetMouseButtonCallback(expr, fp)
		case OP_GLFW_SET_CURSOR_POS_CALLBACK:
			op_glfw_SetCursorPosCallback(expr, fp)
		case OP_GLFW_GET_CURSOR_POS:
			op_glfw_GetCursorPos(expr, fp)
		case OP_GLFW_SET_INPUT_MODE:
			op_glfw_SetInputMode(expr, fp)
		case OP_GLFW_SET_WINDOW_POS:
			op_glfw_SetWindowPos(expr, fp)
		case OP_GLFW_GET_KEY:
			op_glfw_GetKey(expr, fp)

			// gltext
		case OP_GLTEXT_LOAD_TRUE_TYPE:
			op_gltext_LoadTrueType(expr, fp)
		case OP_GLTEXT_PRINTF:
			op_gltext_Printf(expr, fp)
		case OP_GLTEXT_METRICS:
			op_gltext_Metrics(expr, fp)
		case OP_GLTEXT_TEXTURE:
			op_gltext_Texture(expr, fp)
		case OP_GLTEXT_NEXT_RUNE:
			op_gltext_NextRune(expr, fp)
		case OP_GLTEXT_GLYPH_BOUNDS:
			op_gltext_GlyphBounds(expr, fp)
		}
	}
}
