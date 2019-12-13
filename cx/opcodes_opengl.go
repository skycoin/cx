// +build opengl

package cxcore

const (
	// gogl
	OP_GL_INIT = iota + END_OF_BASE_OPS
	OP_GL_STRS
	OP_GL_FREE
	OP_GL_NEW_TEXTURE
	OP_GL_NEW_GIF
	OP_GL_FREE_GIF
	OP_GL_GIF_FRAME_TO_TEXTURE

	// gl_0_0
	OP_GL_MATRIX_MODE
	OP_GL_ROTATEF
	OP_GL_TRANSLATEF
	OP_GL_LOAD_IDENTITY
	OP_GL_PUSH_MATRIX
	OP_GL_POP_MATRIX
	OP_GL_ENABLE_CLIENT_STATE
	OP_GL_COLOR3F
	OP_GL_COLOR4F
	OP_GL_BEGIN
	OP_GL_END
	OP_GL_NORMAL3F
	OP_GL_VERTEX_2F
	OP_GL_VERTEX_3F
	OP_GL_LIGHTFV
	OP_GL_FRUSTUM
	OP_GL_TEX_ENVI
	OP_GL_ORTHO
	OP_GL_SCALEF
	OP_GL_TEX_COORD_2D
	OP_GL_TEX_COORD_2F

	// gl_1_0
	OP_GL_CULL_FACE
	OP_GL_HINT
	OP_GL_SCISSOR
	OP_GL_TEX_PARAMETERI
	OP_GL_TEX_IMAGE_2D
	OP_GL_CLEAR
	OP_GL_CLEAR_COLOR
	OP_GL_CLEAR_STENCIL
	OP_GL_CLEAR_DEPTH
	OP_GL_STENCIL_MASK
	OP_GL_COLOR_MASK
	OP_GL_DEPTH_MASK
	OP_GL_DISABLE
	OP_GL_ENABLE
	OP_GL_BLEND_FUNC
	OP_GL_STENCIL_FUNC
	OP_GL_STENCIL_OP
	OP_GL_DEPTH_FUNC
	OP_GL_GET_ERROR
	OP_GL_GET_TEX_LEVEL_PARAMETERIV
	OP_GL_DEPTH_RANGE
	OP_GL_VIEWPORT

	// gl_1_1
	OP_GL_DRAW_ARRAYS
	OP_GL_BIND_TEXTURE
	OP_GL_DELETE_TEXTURES
	OP_GL_GEN_TEXTURES

	// gl_1_3
	OP_GL_ACTIVE_TEXTURE

	// gl_1_5
	OP_GL_BIND_BUFFER
	OP_GL_DELETE_BUFFERS
	OP_GL_GEN_BUFFERS
	OP_GL_BUFFER_DATA
	OP_GL_BUFFER_SUB_DATA

	// gl_2_0
	OP_GL_STENCIL_OP_SEPARATE
	OP_GL_STENCIL_FUNC_SEPARATE
	OP_GL_STENCIL_MASK_SEPARATE
	OP_GL_ATTACH_SHADER
	OP_GL_BIND_ATTRIB_LOCATION
	OP_GL_COMPILE_SHADER
	OP_GL_CREATE_PROGRAM
	OP_GL_CREATE_SHADER
	OP_GL_DELETE_PROGRAM
	OP_GL_DELETE_SHADER
	OP_GL_DETACH_SHADER
	OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY
	OP_GL_GET_ATTRIB_LOCATION
	OP_GL_GET_SHADERIV
	OP_GL_GET_UNIFORM_LOCATION
	OP_GL_LINK_PROGRAM
	OP_GL_SHADER_SOURCE
	OP_GL_USE_PROGRAM
	OP_GL_UNIFORM_1F
	OP_GL_UNIFORM_1I
	OP_GL_VERTEX_ATTRIB_POINTER
	OP_GL_VERTEX_ATTRIB_POINTER_I32

	// gl_3_0
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
	OP_GL_BIND_VERTEX_ARRAY
	OP_GL_DELETE_VERTEX_ARRAYS
	OP_GL_GEN_VERTEX_ARRAYS

	// goglfw
	OP_GLFW_FULLSCREEN

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
	OP_GLFW_GET_WINDOW_POS
	OP_GLFW_GET_WINDOW_SIZE
	OP_GLFW_SWAP_INTERVAL
	OP_GLFW_SET_KEY_CALLBACK
	OP_GLFW_SET_KEY_CALLBACK_EX
	OP_GLFW_GET_TIME
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK_EX
	OP_GLFW_SET_CURSOR_POS_CALLBACK
	OP_GLFW_SET_CURSOR_POS_CALLBACK_EX
	OP_GLFW_SET_FRAMEBUFFER_SIZE_CALLBACK
	OP_GLFW_SET_WINDOW_POS_CALLBACK
	OP_GLFW_SET_WINDOW_SIZE_CALLBACK
	OP_GLFW_GET_CURSOR_POS
	OP_GLFW_SET_INPUT_MODE
	OP_GLFW_SET_WINDOW_POS
	OP_GLFW_GET_KEY
	OP_GLFW_FUNC_I32_I32
	OP_GLFW_CALL_I32_I32

	// gltext
	OP_GLTEXT_LOAD_TRUE_TYPE
	OP_GLTEXT_LOAD_TRUE_TYPE_EX
	OP_GLTEXT_PRINTF
	OP_GLTEXT_METRICS
	OP_GLTEXT_TEXTURE
	OP_GLTEXT_NEXT_GLYPH
	OP_GLTEXT_GLYPH_BOUNDS
	OP_GLTEXT_GLYPH_METRICS
	OP_GLTEXT_GLYPH_INFO

	// goal
	OP_AL_LOAD_WAV

	// openal
	OP_AL_CLOSE_DEVICE
	OP_AL_DELETE_BUFFERS
	OP_AL_DELETE_SOURCES
	OP_AL_DEVICE_ERROR
	OP_AL_ERROR
	OP_AL_EXTENSIONS
	OP_AL_OPEN_DEVICE
	OP_AL_PAUSE_SOURCES
	OP_AL_PLAY_SOURCES
	OP_AL_RENDERER
	OP_AL_REWIND_SOURCES
	OP_AL_STOP_SOURCES
	OP_AL_VENDOR
	OP_AL_VERSION
	OP_AL_GEN_BUFFERS
	OP_AL_BUFFER_DATA
	OP_AL_GEN_SOURCES
	OP_AL_SOURCE_BUFFERS_PROCESSED
	OP_AL_SOURCE_BUFFERS_QUEUED
	OP_AL_SOURCE_QUEUE_BUFFERS
	OP_AL_SOURCE_STATE
	OP_AL_SOURCE_UNQUEUE_BUFFERS
)

func init() {
	// gogl
	AddOpCode(OP_GL_INIT, "gl.Init",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_GL_STRS, "gl.Strs",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_FREE, "gl.Free",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_NEW_TEXTURE, "gl.NewTexture",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_NEW_GIF, "gl.NewGIF",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_FREE_GIF, "gl.FreeGIF",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GIF_FRAME_TO_TEXTURE, "gl.GIFFrameToTexture",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})

	// gl_0.0
	AddOpCode(OP_GL_MATRIX_MODE, "gl.MatrixMode",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_ROTATEF, "gl.Rotatef",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_TRANSLATEF, "gl.Translatef",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_LOAD_IDENTITY, "gl.LoadIdentity",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_GL_PUSH_MATRIX, "gl.PushMatrix",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_GL_POP_MATRIX, "gl.PopMatrix",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_GL_ENABLE_CLIENT_STATE, "gl.EnableClientState",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_COLOR3F, "gl.Color3f",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_COLOR4F, "gl.Color4f",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_BEGIN, "gl.Begin",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_END, "gl.End",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_GL_NORMAL3F, "gl.Normal3f",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_VERTEX_2F, "gl.Vertex2f",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_VERTEX_3F, "gl.Vertex3f",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_LIGHTFV, "gl.Lightfv",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_FRUSTUM, "gl.Frustum",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_TEX_ENVI, "gl.TexEnvi",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_ORTHO, "gl.Ortho",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_SCALEF, "gl.Scalef",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_TEX_COORD_2D, "gl.TexCoord2d",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_TEX_COORD_2F, "gl.TexCoord2f",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})

	// gl_1_0
	AddOpCode(OP_GL_CULL_FACE, "gl.CullFace",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_HINT, "gl.Hint",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_SCISSOR, "gl.Scissor",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_TEX_PARAMETERI, "gl.TexParameteri",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_TEX_IMAGE_2D, "gl.TexImage2D",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_CLEAR, "gl.Clear",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_CLEAR_COLOR, "gl.ClearColor",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_CLEAR_STENCIL, "gl.ClearStencil",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_CLEAR_DEPTH, "gl.ClearDepth",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_STENCIL_MASK, "gl.StencilMask",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_COLOR_MASK, "gl.ColorMask",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DEPTH_MASK, "gl.DepthMask",
		[]*CXArgument{newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DISABLE, "gl.Disable",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_ENABLE, "gl.Enable",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_BLEND_FUNC, "gl.BlendFunc",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_STENCIL_FUNC, "gl.StencilFunc",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_STENCIL_OP, "gl.StencilOp",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DEPTH_FUNC, "gl.DepthFunc",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GET_ERROR, "gl.GetError",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_GET_TEX_LEVEL_PARAMETERIV, "gl.GetTexLevelParameteriv",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_DEPTH_RANGE, "gl.DepthRange",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_VIEWPORT, "gl.Viewport",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})

	// gl_1_1
	AddOpCode(OP_GL_DRAW_ARRAYS, "gl.DrawArrays",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_BIND_TEXTURE, "gl.BindTexture",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DELETE_TEXTURES, "gl.DeleteTextures",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GEN_TEXTURES, "gl.GenTextures",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})

	// gl_1_3
	AddOpCode(OP_GL_ACTIVE_TEXTURE, "gl.ActiveTexture",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})

	// gl_1_5
	AddOpCode(OP_GL_BIND_BUFFER, "gl.BindBuffer",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DELETE_BUFFERS, "gl.DeleteBuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GEN_BUFFERS, "gl.GenBuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_BUFFER_DATA, "gl.BufferData",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_F32, true), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_BUFFER_SUB_DATA, "gl.BufferSubData",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})

	//gl_2_0
	AddOpCode(OP_GL_STENCIL_OP_SEPARATE, "gl.StencilOpSeparate",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_STENCIL_FUNC_SEPARATE, "gl.StencilFuncSeparate",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_STENCIL_MASK_SEPARATE, "gl.StencilMaskSeparate",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_ATTACH_SHADER, "gl.AttachShader",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_BIND_ATTRIB_LOCATION, "gl.BindAttribLocation",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_COMPILE_SHADER, "gl.CompileShader",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_CREATE_PROGRAM, "gl.CreateProgram",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_CREATE_SHADER, "gl.CreateShader",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_DELETE_PROGRAM, "gl.DeleteProgram",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DELETE_SHADER, "gl.DeleteShader",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DETACH_SHADER, "gl.DetachShader",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY, "gl.EnableVertexAttribArray",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GET_ATTRIB_LOCATION, "gl.GetAttribLocation",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_GET_SHADERIV, "gl.GetShaderiv",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GET_UNIFORM_LOCATION, "gl.GetUniformLocation",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_LINK_PROGRAM, "gl.LinkProgram",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_SHADER_SOURCE, "gl.ShaderSource",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_USE_PROGRAM, "gl.UseProgram",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_UNIFORM_1F, "gl.Uniform1f",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_UNIFORM_1I, "gl.Uniform1i",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_VERTEX_ATTRIB_POINTER, "gl.VertexAttribPointer",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_BOOL, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_VERTEX_ATTRIB_POINTER_I32, "gl.VertexAttribPointerI32",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_BOOL, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})

	// gl_3_0
	AddOpCode(OP_GL_BIND_RENDERBUFFER, "gl.BindRenderbuffer",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DELETE_RENDERBUFFERS, "gl.DeleteRenderbuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GEN_RENDERBUFFERS, "gl.GenRenderbuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_RENDERBUFFER_STORAGE, "gl.RenderbufferStorage",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_BIND_FRAMEBUFFER, "gl.BindFramebuffer",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DELETE_FRAMEBUFFERS, "gl.DeleteFramebuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GEN_FRAMEBUFFERS, "gl.GenFramebuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_CHECK_FRAMEBUFFER_STATUS, "gl.CheckFramebufferStatus",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GL_FRAMEBUFFER_TEXTURE_2D, "gl.FramebufferTexture2D",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_FRAMEBUFFER_RENDERBUFFER, "gl.FramebufferRenderbuffer",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_BIND_VERTEX_ARRAY, "gl.BindVertexArray",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_DELETE_VERTEX_ARRAYS, "gl.DeleteVertexArrays",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GL_GEN_VERTEX_ARRAYS, "gl.GenVertexArrays",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})

	// goglfw
	AddOpCode(OP_GLFW_FULLSCREEN, "glfw.Fullscreen",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_BOOL, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})

	// glfw
	AddOpCode(OP_GLFW_INIT, "glfw.Init",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_WINDOW_HINT, "glfw.WindowHint",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_CREATE_WINDOW, "glfw.CreateWindow",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_MAKE_CONTEXT_CURRENT, "glfw.MakeContextCurrent",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SHOULD_CLOSE, "glfw.ShouldClose",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_GLFW_SET_SHOULD_CLOSE, "glfw.SetShouldClose",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_POLL_EVENTS, "glfw.PollEvents",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SWAP_BUFFERS, "glfw.SwapBuffers",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_GET_FRAMEBUFFER_SIZE, "glfw.GetFramebufferSize",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLFW_GET_WINDOW_POS, "glfw.GetWindowPos",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLFW_GET_WINDOW_SIZE, "glfw.GetWindowSize",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLFW_SWAP_INTERVAL, "glfw.SwapInterval",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_KEY_CALLBACK, "glfw.SetKeyCallback",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_KEY_CALLBACK_EX, "glfw.SetKeyCallbackEx",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_GET_TIME, "glfw.GetTime",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_GLFW_SET_MOUSE_BUTTON_CALLBACK, "glfw.SetMouseButtonCallback",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_MOUSE_BUTTON_CALLBACK_EX, "glfw.SetMouseButtonCallbackEx",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_CURSOR_POS_CALLBACK, "glfw.SetCursorPosCallback",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_CURSOR_POS_CALLBACK_EX, "glfw.SetCursorPosCallbackEx",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_FRAMEBUFFER_SIZE_CALLBACK, "glfw.SetFramebufferSizeCallback",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_WINDOW_POS_CALLBACK, "glfw.SetWindowPosCallback",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_WINDOW_SIZE_CALLBACK, "glfw.SetWindowSizeCallback",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_GET_CURSOR_POS, "glfw.GetCursorPos",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)})
	AddOpCode(OP_GLFW_SET_INPUT_MODE, "glfw.SetInputMode",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_SET_WINDOW_POS, "glfw.SetWindowPos",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLFW_GET_KEY, "glfw.GetKey",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLFW_FUNC_I32_I32, "glfw.func_i32_i32",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLFW_CALL_I32_I32, "glfw.call_i32_i32",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})

	// gltext
	AddOpCode(OP_GLTEXT_LOAD_TRUE_TYPE, "gltext.LoadTrueType",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLTEXT_PRINTF, "gltext.Printf",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_GLTEXT_METRICS, "gltext.Metrics",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLTEXT_TEXTURE, "gltext.Texture",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLTEXT_NEXT_GLYPH, "gltext.NextGlyph",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLTEXT_GLYPH_BOUNDS, "gltext.GlyphBounds",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLTEXT_GLYPH_METRICS, "gltext.GlyphMetrics",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})
	AddOpCode(OP_GLTEXT_GLYPH_INFO, "gltext.GlyphInfo",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)})

	// goal
	AddOpCode(OP_AL_LOAD_WAV, "al.LoadWav",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false),
			newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false),
			newOpPar(TYPE_I32, false), newOpPar(TYPE_I64, false), newOpPar(TYPE_UI8, true)})

	// openal
	AddOpCode(OP_AL_CLOSE_DEVICE, "al.CloseDevice",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_AL_DELETE_BUFFERS, "al.DeleteBuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_AL_DELETE_SOURCES, "al.DeleteSources",
		[]*CXArgument{newOpPar(TYPE_I32, true)},
		[]*CXArgument{})
	AddOpCode(OP_AL_DEVICE_ERROR, "al.DeviceError",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_AL_ERROR, "al.Error",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_AL_EXTENSIONS, "al.Extensions",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_AL_OPEN_DEVICE, "al.OpenDevice",
		[]*CXArgument{},
		[]*CXArgument{})
	AddOpCode(OP_AL_PAUSE_SOURCES, "al.PauseSources",
		[]*CXArgument{newOpPar(TYPE_I32, true)},
		[]*CXArgument{})
	AddOpCode(OP_AL_PLAY_SOURCES, "al.PlaySources",
		[]*CXArgument{newOpPar(TYPE_I32, true)},
		[]*CXArgument{})
	AddOpCode(OP_AL_RENDERER, "al.Renderer",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_AL_REWIND_SOURCES, "al.RewindSources",
		[]*CXArgument{newOpPar(TYPE_I32, true)},
		[]*CXArgument{})
	AddOpCode(OP_AL_STOP_SOURCES, "al.StopSources",
		[]*CXArgument{newOpPar(TYPE_I32, true)},
		[]*CXArgument{})
	AddOpCode(OP_AL_VENDOR, "al.Vendor",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_AL_VERSION, "al.Version",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_AL_GEN_BUFFERS, "al.GenBuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, true)})
	AddOpCode(OP_AL_BUFFER_DATA, "al.BufferData",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_AL_GEN_SOURCES, "al.GenSources",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, true)})
	AddOpCode(OP_AL_SOURCE_BUFFERS_PROCESSED, "al.SourceBuffersProcessed",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_AL_SOURCE_BUFFERS_QUEUED, "al.SourceBuffersQueued",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_AL_SOURCE_QUEUE_BUFFERS, "al.SourceQueueBuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, true)},
		[]*CXArgument{})
	AddOpCode(OP_AL_SOURCE_STATE, "al.SourceState",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_AL_SOURCE_UNQUEUE_BUFFERS, "al.SourceUnueueBuffers",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, true)},
		[]*CXArgument{})

	// exec
	handleOpcode := func(opCode int) opcodeHandler {
		switch opCode {
		// gogl
		case OP_GL_INIT:
			return op_gl_Init
		case OP_GL_STRS:
			return op_gl_Strs
		case OP_GL_FREE:
			return op_gl_Free
		case OP_GL_NEW_TEXTURE:
			return op_gl_NewTexture
		case OP_GL_NEW_GIF:
			return op_gl_NewGIF
		case OP_GL_FREE_GIF:
			return op_gl_FreeGIF
		case OP_GL_GIF_FRAME_TO_TEXTURE:
			return op_gl_GIFFrameToTexture

		// gl_0_0
		case OP_GL_MATRIX_MODE:
			return op_gl_MatrixMode
		case OP_GL_ROTATEF:
			return op_gl_Rotatef
		case OP_GL_TRANSLATEF:
			return op_gl_Translatef
		case OP_GL_LOAD_IDENTITY:
			return op_gl_LoadIdentity
		case OP_GL_PUSH_MATRIX:
			return op_gl_PushMatrix
		case OP_GL_POP_MATRIX:
			return op_gl_PopMatrix
		case OP_GL_ENABLE_CLIENT_STATE:
			return op_gl_EnableClientState
		case OP_GL_COLOR3F:
			return op_gl_Color3f
		case OP_GL_COLOR4F:
			return op_gl_Color4f
		case OP_GL_BEGIN:
			return op_gl_Begin
		case OP_GL_END:
			return op_gl_End
		case OP_GL_NORMAL3F:
			return op_gl_Normal3f
		case OP_GL_VERTEX_2F:
			return op_gl_Vertex2f
		case OP_GL_VERTEX_3F:
			return op_gl_Vertex3f
		case OP_GL_LIGHTFV:
			return op_gl_Lightfv
		case OP_GL_FRUSTUM:
			return op_gl_Frustum
		case OP_GL_TEX_ENVI:
			return op_gl_TexEnvi
		case OP_GL_ORTHO:
			return op_gl_Ortho
		case OP_GL_SCALEF:
			return op_gl_Scalef
		case OP_GL_TEX_COORD_2D:
			return op_gl_TexCoord2d
		case OP_GL_TEX_COORD_2F:
			return op_gl_TexCoord2f

		// gl_1_0
		case OP_GL_CULL_FACE:
			return op_gl_CullFace
		case OP_GL_HINT:
			return op_gl_Hint
		case OP_GL_SCISSOR:
			return op_gl_Scissor
		case OP_GL_TEX_PARAMETERI:
			return op_gl_TexParameteri
		case OP_GL_TEX_IMAGE_2D:
			return op_gl_TexImage2D
		case OP_GL_CLEAR:
			return op_gl_Clear
		case OP_GL_CLEAR_COLOR:
			return op_gl_ClearColor
		case OP_GL_CLEAR_STENCIL:
			return op_gl_ClearStencil
		case OP_GL_CLEAR_DEPTH:
			return op_gl_ClearDepth
		case OP_GL_STENCIL_MASK:
			return op_gl_StencilMask
		case OP_GL_COLOR_MASK:
			return op_gl_ColorMask
		case OP_GL_DEPTH_MASK:
			return op_gl_DepthMask
		case OP_GL_DISABLE:
			return op_gl_Disable
		case OP_GL_ENABLE:
			return op_gl_Enable
		case OP_GL_BLEND_FUNC:
			return op_gl_BlendFunc
		case OP_GL_STENCIL_FUNC:
			return op_gl_StencilFunc
		case OP_GL_STENCIL_OP:
			return op_gl_StencilOp
		case OP_GL_DEPTH_FUNC:
			return op_gl_DepthFunc
		case OP_GL_GET_ERROR:
			return op_gl_GetError
		case OP_GL_GET_TEX_LEVEL_PARAMETERIV:
			return op_gl_GetTexLevelParameteriv
		case OP_GL_DEPTH_RANGE:
			return op_gl_DepthRange
		case OP_GL_VIEWPORT:
			return op_gl_Viewport

		// gl_1_1
		case OP_GL_DRAW_ARRAYS:
			return op_gl_DrawArrays
		case OP_GL_BIND_TEXTURE:
			return op_gl_BindTexture
		case OP_GL_DELETE_TEXTURES:
			return op_gl_DeleteTextures
		case OP_GL_GEN_TEXTURES:
			return op_gl_GenTextures

		// gl_1_3
		case OP_GL_ACTIVE_TEXTURE:
			return op_gl_ActiveTexture

		// gl_1_5
		case OP_GL_BIND_BUFFER:
			return op_gl_BindBuffer
		case OP_GL_DELETE_BUFFERS:
			return op_gl_DeleteBuffers
		case OP_GL_GEN_BUFFERS:
			return op_gl_GenBuffers
		case OP_GL_BUFFER_DATA:
			return op_gl_BufferData
		case OP_GL_BUFFER_SUB_DATA:
			return op_gl_BufferSubData

		// gl_2_0
		case OP_GL_STENCIL_OP_SEPARATE:
			return op_gl_StencilOpSeparate
		case OP_GL_STENCIL_FUNC_SEPARATE:
			return op_gl_StencilFuncSeparate
		case OP_GL_STENCIL_MASK_SEPARATE:
			return op_gl_StencilMaskSeparate
		case OP_GL_ATTACH_SHADER:
			return op_gl_AttachShader
		case OP_GL_BIND_ATTRIB_LOCATION:
			return op_gl_BindAttribLocation
		case OP_GL_COMPILE_SHADER:
			return op_gl_CompileShader
		case OP_GL_CREATE_PROGRAM:
			return op_gl_CreateProgram
		case OP_GL_CREATE_SHADER:
			return op_gl_CreateShader
		case OP_GL_DELETE_PROGRAM:
			return op_gl_DeleteProgram
		case OP_GL_DELETE_SHADER:
			return op_gl_DeleteShader
		case OP_GL_DETACH_SHADER:
			return op_gl_DetachShader
		case OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY:
			return op_gl_EnableVertexAttribArray
		case OP_GL_GET_ATTRIB_LOCATION:
			return op_gl_GetAttribLocation
		case OP_GL_GET_SHADERIV:
			return op_gl_GetShaderiv
		case OP_GL_GET_UNIFORM_LOCATION:
			return op_gl_GetUniformLocation
		case OP_GL_LINK_PROGRAM:
			return op_gl_LinkProgram
		case OP_GL_SHADER_SOURCE:
			return op_gl_ShaderSource
		case OP_GL_USE_PROGRAM:
			return op_gl_UseProgram
		case OP_GL_UNIFORM_1F:
			return op_gl_Uniform1f
		case OP_GL_UNIFORM_1I:
			return op_gl_Uniform1i
		case OP_GL_VERTEX_ATTRIB_POINTER:
			return op_gl_VertexAttribPointer
		case OP_GL_VERTEX_ATTRIB_POINTER_I32:
			return op_gl_VertexAttribPointerI32

		// gl_3_0
		case OP_GL_BIND_RENDERBUFFER:
			return op_gl_BindRenderbuffer
		case OP_GL_DELETE_RENDERBUFFERS:
			return op_gl_DeleteRenderbuffers
		case OP_GL_GEN_RENDERBUFFERS:
			return op_gl_GenRenderbuffers
		case OP_GL_RENDERBUFFER_STORAGE:
			return op_gl_RenderbufferStorage
		case OP_GL_BIND_FRAMEBUFFER:
			return op_gl_BindFramebuffer
		case OP_GL_DELETE_FRAMEBUFFERS:
			return op_gl_DeleteFramebuffers
		case OP_GL_GEN_FRAMEBUFFERS:
			return op_gl_GenFramebuffers
		case OP_GL_CHECK_FRAMEBUFFER_STATUS:
			return op_gl_CheckFramebufferStatus
		case OP_GL_FRAMEBUFFER_TEXTURE_2D:
			return op_gl_FramebufferTexture2D
		case OP_GL_FRAMEBUFFER_RENDERBUFFER:
			return op_gl_FramebufferRenderbuffer
		case OP_GL_BIND_VERTEX_ARRAY:
			return op_gl_BindVertexArray
		case OP_GL_DELETE_VERTEX_ARRAYS:
			return op_gl_DeleteVertexArrays
		case OP_GL_GEN_VERTEX_ARRAYS:
			return op_gl_GenVertexArrays

		// goglfw
		case OP_GLFW_FULLSCREEN:
			return op_glfw_Fullscreen

		// glfw
		case OP_GLFW_INIT:
			return op_glfw_Init
		case OP_GLFW_WINDOW_HINT:
			return op_glfw_WindowHint
		case OP_GLFW_CREATE_WINDOW:
			return op_glfw_CreateWindow
		case OP_GLFW_MAKE_CONTEXT_CURRENT:
			return op_glfw_MakeContextCurrent
		case OP_GLFW_SHOULD_CLOSE:
			return op_glfw_ShouldClose
		case OP_GLFW_SET_SHOULD_CLOSE:
			return op_glfw_SetShouldClose
		case OP_GLFW_POLL_EVENTS:
			return op_glfw_PollEvents
		case OP_GLFW_SWAP_BUFFERS:
			return op_glfw_SwapBuffers
		case OP_GLFW_GET_FRAMEBUFFER_SIZE:
			return op_glfw_GetFramebufferSize
		case OP_GLFW_GET_WINDOW_POS:
			return op_glfw_GetWindowPos
		case OP_GLFW_GET_WINDOW_SIZE:
			return op_glfw_GetWindowSize
		case OP_GLFW_SWAP_INTERVAL:
			return op_glfw_SwapInterval
		case OP_GLFW_SET_KEY_CALLBACK:
			return op_glfw_SetKeyCallback
		case OP_GLFW_SET_KEY_CALLBACK_EX:
			return op_glfw_SetKeyCallbackEx
		case OP_GLFW_GET_TIME:
			return op_glfw_GetTime
		case OP_GLFW_SET_MOUSE_BUTTON_CALLBACK:
			return op_glfw_SetMouseButtonCallback
		case OP_GLFW_SET_MOUSE_BUTTON_CALLBACK_EX:
			return op_glfw_SetMouseButtonCallbackEx
		case OP_GLFW_SET_CURSOR_POS_CALLBACK:
			return op_glfw_SetCursorPosCallback
		case OP_GLFW_SET_CURSOR_POS_CALLBACK_EX:
			return op_glfw_SetCursorPosCallbackEx
		case OP_GLFW_SET_FRAMEBUFFER_SIZE_CALLBACK:
			return op_glfw_SetFramebufferSizeCallback
		case OP_GLFW_SET_WINDOW_POS_CALLBACK:
			return op_glfw_SetWindowPosCallback
		case OP_GLFW_SET_WINDOW_SIZE_CALLBACK:
			return op_glfw_SetWindowSizeCallback
		case OP_GLFW_GET_CURSOR_POS:
			return op_glfw_GetCursorPos
		case OP_GLFW_SET_INPUT_MODE:
			return op_glfw_SetInputMode
		case OP_GLFW_SET_WINDOW_POS:
			return op_glfw_SetWindowPos
		case OP_GLFW_GET_KEY:
			return op_glfw_GetKey
		case OP_GLFW_FUNC_I32_I32:
			return op_glfw_func_i32_i32
		case OP_GLFW_CALL_I32_I32:
			return op_glfw_call_i32_i32

		// gltext
		case OP_GLTEXT_LOAD_TRUE_TYPE:
			return op_gltext_LoadTrueType
		case OP_GLTEXT_PRINTF:
			return op_gltext_Printf
		case OP_GLTEXT_METRICS:
			return op_gltext_Metrics
		case OP_GLTEXT_TEXTURE:
			return op_gltext_Texture
		case OP_GLTEXT_NEXT_GLYPH:
			return op_gltext_NextGlyph
		case OP_GLTEXT_GLYPH_BOUNDS:
			return op_gltext_GlyphBounds
		case OP_GLTEXT_GLYPH_METRICS:
			return op_gltext_GlyphMetrics
		case OP_GLTEXT_GLYPH_INFO:
			return op_gltext_GlyphInfo

		// goal
		case OP_AL_LOAD_WAV:
			return opAlLoadWav

		// openal
		case OP_AL_CLOSE_DEVICE:
			return opAlCloseDevice
		case OP_AL_DELETE_BUFFERS:
			return opAlDeleteBuffers
		case OP_AL_DELETE_SOURCES:
			return opAlDeleteSources
		case OP_AL_DEVICE_ERROR:
			return opAlDeviceError
		case OP_AL_ERROR:
			return opAlError
		case OP_AL_EXTENSIONS:
			return opAlExtensions
		case OP_AL_OPEN_DEVICE:
			return opAlOpenDevice
		case OP_AL_PAUSE_SOURCES:
			return opAlPauseSources
		case OP_AL_PLAY_SOURCES:
			return opAlPlaySources
		case OP_AL_RENDERER:
			return opAlRenderer
		case OP_AL_REWIND_SOURCES:
			return opAlRewindSources
		case OP_AL_STOP_SOURCES:
			return opAlStopSources
		case OP_AL_VENDOR:
			return opAlVendor
		case OP_AL_VERSION:
			return opAlVersion
		case OP_AL_GEN_BUFFERS:
			return opAlGenBuffers
		case OP_AL_BUFFER_DATA:
			return opAlBufferData
		case OP_AL_GEN_SOURCES:
			return opAlGenSources
		case OP_AL_SOURCE_BUFFERS_PROCESSED:
			return opAlSourceBuffersProcessed
		case OP_AL_SOURCE_BUFFERS_QUEUED:
			return opAlSourceBuffersQueued
		case OP_AL_SOURCE_QUEUE_BUFFERS:
			return opAlSourceQueueBuffers
		case OP_AL_SOURCE_STATE:
			return opAlSourceState
		case OP_AL_SOURCE_UNQUEUE_BUFFERS:
			return opAlSourceUnqueueBuffers
		}

		return nil
	}

	opcodeHandlerFinders = append(opcodeHandlerFinders, handleOpcode)
}
