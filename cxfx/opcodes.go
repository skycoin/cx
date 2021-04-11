// +build cxfx

// 24:58
// 16/02 : 20:20->23:50
package cxfx

import (
	. "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/packages/http"
	"github.com/skycoin/cx/cxos"
)

const (
	// gogl
	OP_GL_INIT = iota + cxos.END_OF_OS_OPS
	OP_GL_DESTROY
	OP_GL_STRS
	OP_GL_FREE
	OP_GL_NEW_TEXTURE
	OP_GL_NEW_TEXTURE_CUBE
	OP_GL_NEW_GIF
	OP_GL_FREE_GIF
	OP_GL_GIF_FRAME_TO_TEXTURE
	OP_GL_UPLOAD_IMAGE_TO_TEXTURE
	OP_CX_RELEASE_TEXTURE
	OP_CX_TEXTURE_GET_PIXEL
	OP_GL_APPEND_F32
	OP_GL_APPEND_UI32
	OP_GL_APPEND_UI16

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
	OP_GL_FRONT_FACE
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
	OP_GL_DRAW_ELEMENTS
	OP_GL_BIND_TEXTURE
	OP_GL_DELETE_TEXTURES
	OP_GL_GEN_TEXTURES

	// gl_1_3
	OP_GL_ACTIVE_TEXTURE

	// gl_1_4
	OP_GL_BLEND_FUNC_SEPARATE

	// gl_1_5
	OP_GL_BIND_BUFFER
	OP_GL_DELETE_BUFFERS
	OP_GL_GEN_BUFFERS
	OP_GL_BUFFER_DATA
	OP_GL_BUFFER_SUB_DATA

	// gl_2_0
	OP_GL_DRAW_BUFFERS
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
	OP_GL_GET_PROGRAMIV
	OP_GL_GET_PROGRAM_INFO_LOG
	OP_GL_GET_SHADERIV
	OP_GL_GET_SHADER_INFO_LOG
	OP_GL_GET_UNIFORM_LOCATION
	OP_GL_LINK_PROGRAM
	OP_GL_SHADER_SOURCE
	OP_GL_USE_PROGRAM
	OP_GL_UNIFORM_1F
	OP_GL_UNIFORM_2F
	OP_GL_UNIFORM_3F
	OP_GL_UNIFORM_4F
	OP_GL_UNIFORM_1I
	OP_GL_UNIFORM_2I
	OP_GL_UNIFORM_3I
	OP_GL_UNIFORM_4I
	OP_GL_UNIFORM_1FV
	OP_GL_UNIFORM_2FV
	OP_GL_UNIFORM_3FV
	OP_GL_UNIFORM_4FV
	OP_GL_UNIFORM_1IV
	OP_GL_UNIFORM_2IV
	OP_GL_UNIFORM_3IV
	OP_GL_UNIFORM_4IV
	OP_GL_UNIFORM_MATRIX_2FV
	OP_GL_UNIFORM_MATRIX_3FV
	OP_GL_UNIFORM_MATRIX_4FV
	OP_GL_UNIFORM_V4F
	OP_GL_UNIFORM_M44F
	OP_GL_UNIFORM_M44FV
	OP_GL_VERTEX_ATTRIB_POINTER
	OP_GL_VERTEX_ATTRIB_POINTER_I32

	// gl_3_0
	OP_GL_CLEAR_BUFFER_I
	OP_GL_CLEAR_BUFFER_UI
	OP_GL_CLEAR_BUFFER_F
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
	OP_GL_GENERATE_MIPMAP
	OP_GL_BIND_VERTEX_ARRAY
	OP_GL_BIND_VERTEX_ARRAY_CORE
	OP_GL_DELETE_VERTEX_ARRAYS
	OP_GL_DELETE_VERTEX_ARRAYS_CORE
	OP_GL_GEN_VERTEX_ARRAYS
	OP_GL_GEN_VERTEX_ARRAYS_CORE

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
	OP_GLFW_GET_TIME
	OP_GLFW_SWAP_INTERVAL
	OP_GLFW_SET_START_CALLBACK
	OP_GLFW_SET_STOP_CALLBACK
	OP_GLFW_SET_KEYBOARD_CALLBACK
	OP_GLFW_SET_MOUSE_CALLBACK
	OP_GLFW_SET_FRAMEBUFFER_SIZE_CALLBACK
	OP_GLFW_SET_WINDOW_POS_CALLBACK
	OP_GLFW_SET_WINDOW_SIZE_CALLBACK
	OP_GLFW_SET_KEY_CALLBACK
	OP_GLFW_SET_MOUSE_BUTTON_CALLBACK
	OP_GLFW_SET_CURSOR_POS_CALLBACK
	OP_GLFW_GET_CURSOR_POS
	OP_GLFW_SET_INPUT_MODE
	OP_GLFW_SET_WINDOW_POS
	OP_GLFW_GET_KEY
	OP_GLFW_FUNC_I32_I32
	OP_GLFW_CALL_I32_I32
	OP_GLFW_GET_WINDOW_CONTENT_SCALE
	OP_GLFW_GET_MONITOR_CONTENT_SCALE

	// gltext
	OP_GLTEXT_LOAD_TRUE_TYPE
	OP_GLTEXT_LOAD_TRUE_TYPE_CORE
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
	RegisterOpCode(OP_GL_INIT, "gl.Init", opGlInit, nil, nil)
	RegisterOpCode(OP_GL_DESTROY, "gl.Destroy", opGlDestroy, nil, nil)
	RegisterOpCode(OP_GL_STRS, "gl.Strs", opGlStrs, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GL_FREE, "gl.Free", opGlFree, http.In(ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GL_NEW_TEXTURE, "gl.NewTexture", opGlNewTexture, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_NEW_TEXTURE_CUBE, "gl.NewTextureCube", opGlNewTextureCube, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_NEW_GIF, "gl.NewGIF", opGlNewGIF, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_FREE_GIF, "gl.FreeGIF", opGlFreeGIF, http.In(ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GL_GIF_FRAME_TO_TEXTURE, "gl.GIFFrameToTexture", opGlGIFFrameToTexture, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_UPLOAD_IMAGE_TO_TEXTURE, "gl.UploadImageToTexture", opGlUploadImageToTexture, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_BOOL), nil)
	RegisterOpCode(OP_CX_RELEASE_TEXTURE, "gl.CxReleaseTexture", opCxReleaseTexture, http.In(ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_CX_TEXTURE_GET_PIXEL, "gl.CxTextureGetPixel", opCxTextureGetPixel, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32))

	RegisterOpCode(OP_GL_APPEND_F32, "gl.AppendF32", opGlAppend, http.In(ast.Slice(constants.TYPE_UI8), ast.ConstCxArg_F32), http.Out(ast.Slice(constants.TYPE_UI8)))
	RegisterOpCode(OP_GL_APPEND_UI16, "gl.AppendUI16", opGlAppend, http.In(ast.Slice(constants.TYPE_UI8), ast.ConstCxArg_UI16), http.Out(ast.Slice(constants.TYPE_UI8)))
	RegisterOpCode(OP_GL_APPEND_UI32, "gl.AppendUI32", opGlAppend, http.In(ast.Slice(constants.TYPE_UI8), ast.ConstCxArg_UI32), http.Out(ast.Slice(constants.TYPE_UI8)))

	// gl_0.0
	RegisterOpCode(OP_GL_MATRIX_MODE, "gl.MatrixMode", opGlMatrixMode, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_ROTATEF, "gl.Rotatef", opGlRotatef, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_TRANSLATEF, "gl.Translatef", opGlTranslatef, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_LOAD_IDENTITY, "gl.LoadIdentity", opGlLoadIdentity, nil, nil)
	RegisterOpCode(OP_GL_PUSH_MATRIX, "gl.PushMatrix", opGlPushMatrix, nil, nil)
	RegisterOpCode(OP_GL_POP_MATRIX, "gl.PopMatrix", opGlPopMatrix, nil, nil)
	RegisterOpCode(OP_GL_ENABLE_CLIENT_STATE, "gl.EnableClientState", opGlEnableClientState, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_COLOR3F, "gl.Color3f", opGlColor3f, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_COLOR4F, "gl.Color4f", opGlColor4f, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_BEGIN, "gl.Begin", opGlBegin, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_END, "gl.End", opGlEnd, nil, nil)
	RegisterOpCode(OP_GL_NORMAL3F, "gl.Normal3f", opGlNormal3f, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_VERTEX_2F, "gl.Vertex2f", opGlVertex2f, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_VERTEX_3F, "gl.Vertex3f", opGlVertex3f, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_LIGHTFV, "gl.Lightfv", opGlLightfv, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_FRUSTUM, "gl.Frustum", opGlFrustum, http.In(ast.ConstCxArg_F64, ast.ConstCxArg_F64, ast.ConstCxArg_F64, ast.ConstCxArg_F64, ast.ConstCxArg_F64, ast.ConstCxArg_F64), nil)
	RegisterOpCode(OP_GL_TEX_ENVI, "gl.TexEnvi", opGlTexEnvi, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_ORTHO, "gl.Ortho", opGlOrtho, http.In(ast.ConstCxArg_F64, ast.ConstCxArg_F64, ast.ConstCxArg_F64, ast.ConstCxArg_F64, ast.ConstCxArg_F64, ast.ConstCxArg_F64), nil)
	RegisterOpCode(OP_GL_SCALEF, "gl.Scalef", opGlScalef, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_TEX_COORD_2D, "gl.TexCoord2d", opGlTexCoord2d, http.In(ast.ConstCxArg_F64, ast.ConstCxArg_F64), nil)
	RegisterOpCode(OP_GL_TEX_COORD_2F, "gl.TexCoord2f", opGlTexCoord2f, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)

	// gl_1_0
	RegisterOpCode(OP_GL_CULL_FACE, "gl.CullFace", opGlCullFace, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_FRONT_FACE, "gl.FrontFace", opGlFrontFace, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_HINT, "gl.Hint", opGlHint, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_SCISSOR, "gl.Scissor", opGlScissor, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_TEX_PARAMETERI, "gl.TexParameteri", opGlTexParameteri, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_TEX_IMAGE_2D, "gl.TexImage2D", opGlTexImage2D, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI32)), nil)
	RegisterOpCode(OP_GL_CLEAR, "gl.Clear", opGlClear, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_CLEAR_COLOR, "gl.ClearColor", opGlClearColor, http.In(ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_CLEAR_STENCIL, "gl.ClearStencil", opGlClearStencil, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_CLEAR_DEPTH, "gl.ClearDepth", opGlClearDepth, http.In(ast.ConstCxArg_F64), nil)
	RegisterOpCode(OP_GL_STENCIL_MASK, "gl.StencilMask", opGlStencilMask, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_COLOR_MASK, "gl.ColorMask", opGlColorMask, http.In(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL), nil)
	RegisterOpCode(OP_GL_DEPTH_MASK, "gl.DepthMask", opGlDepthMask, http.In(ast.ConstCxArg_BOOL), nil)
	RegisterOpCode(OP_GL_DISABLE, "gl.Disable", opGlDisable, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_ENABLE, "gl.Enable", opGlEnable, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_BLEND_FUNC, "gl.BlendFunc", opGlBlendFunc, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_STENCIL_FUNC, "gl.StencilFunc", opGlStencilFunc, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_STENCIL_OP, "gl.StencilOp", opGlStencilOp, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DEPTH_FUNC, "gl.DepthFunc", opGlDepthFunc, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_GET_ERROR, "gl.GetError", opGlGetError, nil, http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_GET_TEX_LEVEL_PARAMETERIV, "gl.GetTexLevelParameteriv", opGlGetTexLevelParameteriv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_DEPTH_RANGE, "gl.DepthRange", opGlDepthRange, http.In(ast.ConstCxArg_F64, ast.ConstCxArg_F64), nil)
	RegisterOpCode(OP_GL_VIEWPORT, "gl.Viewport", opGlViewport, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)

	// gl_1_1
	RegisterOpCode(OP_GL_DRAW_ARRAYS, "gl.DrawArrays", opGlDrawArrays, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DRAW_ELEMENTS, "gl.DrawElements", opGlDrawElements, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_BIND_TEXTURE, "gl.BindTexture", opGlBindTexture, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DELETE_TEXTURES, "gl.DeleteTextures", opGlDeleteTextures, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_GEN_TEXTURES, "gl.GenTextures", opGlGenTextures, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))

	// gl_1_3
	RegisterOpCode(OP_GL_ACTIVE_TEXTURE, "gl.ActiveTexture", opGlActiveTexture, http.In(ast.ConstCxArg_I32), nil)

	// gl_1_4
	RegisterOpCode(OP_GL_BLEND_FUNC_SEPARATE, "gl.BlendFuncSeparate", opGlBlendFuncSeparate, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)

	// gl_1_5
	RegisterOpCode(OP_GL_BIND_BUFFER, "gl.BindBuffer", opGlBindBuffer, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DELETE_BUFFERS, "gl.DeleteBuffers", opGlDeleteBuffers, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_GEN_BUFFERS, "gl.GenBuffers", opGlGenBuffers, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_BUFFER_DATA, "gl.BufferData", opGlBufferData, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI8), ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_BUFFER_SUB_DATA, "gl.BufferSubData", opGlBufferSubData, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI8)), nil)

	//gl_2_0
	RegisterOpCode(OP_GL_DRAW_BUFFERS, "gl.DrawBuffers", opGlDrawBuffers, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI32)), nil)
	RegisterOpCode(OP_GL_STENCIL_OP_SEPARATE, "gl.StencilOpSeparate", opGlStencilOpSeparate, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_STENCIL_FUNC_SEPARATE, "gl.StencilFuncSeparate", opGlStencilFuncSeparate, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_STENCIL_MASK_SEPARATE, "gl.StencilMaskSeparate", opGlStencilMaskSeparate, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_ATTACH_SHADER, "gl.AttachShader", opGlAttachShader, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_BIND_ATTRIB_LOCATION, "gl.BindAttribLocation", opGlBindAttribLocation, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GL_COMPILE_SHADER, "gl.CompileShader", opGlCompileShader, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_CREATE_PROGRAM, "gl.CreateProgram", opGlCreateProgram, nil, http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_CREATE_SHADER, "gl.CreateShader", opGlCreateShader, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_DELETE_PROGRAM, "gl.DeleteProgram", opGlDeleteProgram, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DELETE_SHADER, "gl.DeleteShader", opGlDeleteShader, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DETACH_SHADER, "gl.DetachShader", opGlDetachShader, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_ENABLE_VERTEX_ATTRIB_ARRAY, "gl.EnableVertexAttribArray", opGlEnableVertexAttribArray, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_GET_ATTRIB_LOCATION, "gl.GetAttribLocation", opGlGetAttribLocation, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_GET_PROGRAMIV, "gl.GetProgramiv", opGlGetProgramiv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_GET_PROGRAM_INFO_LOG, "gl.GetProgramInfoLog", opGlGetProgramInfoLog, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_STR))
	RegisterOpCode(OP_GL_GET_SHADERIV, "gl.GetShaderiv", opGlGetShaderiv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_GET_SHADER_INFO_LOG, "gl.GetShaderInfoLog", opGlGetShaderInfoLog, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_STR))
	RegisterOpCode(OP_GL_GET_UNIFORM_LOCATION, "gl.GetUniformLocation", opGlGetUniformLocation, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_LINK_PROGRAM, "gl.LinkProgram", opGlLinkProgram, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_SHADER_SOURCE, "gl.ShaderSource", opGlShaderSource, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GL_USE_PROGRAM, "gl.UseProgram", opGlUseProgram, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_UNIFORM_1F, "gl.Uniform1f", opGlUniform1f, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_UNIFORM_2F, "gl.Uniform2f", opGlUniform2f, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_UNIFORM_3F, "gl.Uniform3f", opGlUniform3f, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_UNIFORM_4F, "gl.Uniform4f", opGlUniform4f, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_UNIFORM_1I, "gl.Uniform1i", opGlUniform1i, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_UNIFORM_2I, "gl.Uniform2i", opGlUniform2i, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_UNIFORM_3I, "gl.Uniform3i", opGlUniform3i, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_UNIFORM_4I, "gl.Uniform4i", opGlUniform4i, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_UNIFORM_1FV, "gl.Uniform1fv", opGlUniform1fv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32, ast.Slice(constants.TYPE_F32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_2FV, "gl.Uniform2fv", opGlUniform2fv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32, ast.Slice(constants.TYPE_F32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_3FV, "gl.Uniform3fv", opGlUniform3fv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_4FV, "gl.Uniform4fv", opGlUniform4fv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_1IV, "gl.Uniform1iv", opGlUniform1iv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_2IV, "gl.Uniform2iv", opGlUniform2iv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_3IV, "gl.Uniform3iv", opGlUniform3iv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_4IV, "gl.Uniform4iv", opGlUniform4iv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_MATRIX_2FV, "gl.UniformMatrix2fv", opGlUniformMatrix2fv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_BOOL, ast.Slice(constants.TYPE_F32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_MATRIX_3FV, "gl.UniformMatrix3fv", opGlUniformMatrix3fv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_BOOL, ast.Slice(constants.TYPE_F32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_MATRIX_4FV, "gl.UniformMatrix4fv", opGlUniformMatrix4fv, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_BOOL, ast.Slice(constants.TYPE_F32)), nil)
	RegisterOpCode(OP_GL_UNIFORM_V4F, "gl.UniformV4F", opGlUniformV4F, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_UND_TYPE), nil)
	RegisterOpCode(OP_GL_UNIFORM_M44F, "gl.UniformM44F", opGlUniformM44F, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL, ast.ConstCxArg_UND_TYPE), nil)
	RegisterOpCode(OP_GL_UNIFORM_M44FV, "gl.UniformM44FV", opGlUniformM44FV, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_BOOL, ast.Slice(constants.TYPE_UNDEFINED)), nil)
	RegisterOpCode(OP_GL_VERTEX_ATTRIB_POINTER, "gl.VertexAttribPointer", opGlVertexAttribPointer, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_BOOL, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_VERTEX_ATTRIB_POINTER_I32, "gl.VertexAttribPointerI32", opGlVertexAttribPointerI32, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_BOOL, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)

	// gl_3_0
	RegisterOpCode(OP_GL_CLEAR_BUFFER_I, "gl.ClearBufferI", opGlClearBufferI, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_CLEAR_BUFFER_UI, "gl.ClearBufferUI", opGlClearBufferUI, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_UI32, ast.ConstCxArg_UI32, ast.ConstCxArg_UI32, ast.ConstCxArg_UI32), nil)
	RegisterOpCode(OP_GL_CLEAR_BUFFER_F, "gl.ClearBufferF", opGlClearBufferF, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_F32), nil)
	RegisterOpCode(OP_GL_BIND_RENDERBUFFER, "gl.BindRenderbuffer", opGlBindRenderbuffer, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DELETE_RENDERBUFFERS, "gl.DeleteRenderbuffers", opGlDeleteRenderbuffers, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_GEN_RENDERBUFFERS, "gl.GenRenderbuffers", opGlGenRenderbuffers, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_RENDERBUFFER_STORAGE, "gl.RenderbufferStorage", opGlRenderbufferStorage, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_BIND_FRAMEBUFFER, "gl.BindFramebuffer", opGlBindFramebuffer, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DELETE_FRAMEBUFFERS, "gl.DeleteFramebuffers", opGlDeleteFramebuffers, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_GEN_FRAMEBUFFERS, "gl.GenFramebuffers", opGlGenFramebuffers, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_CHECK_FRAMEBUFFER_STATUS, "gl.CheckFramebufferStatus", opGlCheckFramebufferStatus, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_FRAMEBUFFER_TEXTURE_2D, "gl.FramebufferTexture2D", opGlFramebufferTexture2D, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_FRAMEBUFFER_RENDERBUFFER, "gl.FramebufferRenderbuffer", opGlFramebufferRenderbuffer, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_GENERATE_MIPMAP, "gl.GenerateMipmap", opGlGenerateMipmap, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_BIND_VERTEX_ARRAY, "gl.BindVertexArray", opGlBindVertexArray, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_BIND_VERTEX_ARRAY_CORE, "gl.BindVertexArrayCore", opGlBindVertexArrayCore, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DELETE_VERTEX_ARRAYS, "gl.DeleteVertexArrays", opGlDeleteVertexArrays, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_DELETE_VERTEX_ARRAYS_CORE, "gl.DeleteVertexArraysCore", opGlDeleteVertexArraysCore, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GL_GEN_VERTEX_ARRAYS, "gl.GenVertexArrays", opGlGenVertexArrays, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GL_GEN_VERTEX_ARRAYS_CORE, "gl.GenVertexArraysCore", opGlGenVertexArraysCore, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))

	// goglfw
	RegisterOpCode(OP_GLFW_FULLSCREEN, "glfw.Fullscreen", opGlfwFullscreen, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)

	// glfw
	RegisterOpCode(OP_GLFW_INIT, "glfw.Init", opGlfwInit, nil, nil)
	RegisterOpCode(OP_GLFW_WINDOW_HINT, "glfw.WindowHint", opGlfwWindowHint, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GLFW_CREATE_WINDOW, "glfw.CreateWindow", opGlfwCreateWindow, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_MAKE_CONTEXT_CURRENT, "glfw.MakeContextCurrent", opGlfwMakeContextCurrent, http.In(ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_SHOULD_CLOSE, "glfw.ShouldClose", opGlfwShouldClose, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_GLFW_SET_SHOULD_CLOSE, "glfw.SetShouldClose", opGlfwSetShouldClose, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL), nil)
	RegisterOpCode(OP_GLFW_POLL_EVENTS, "glfw.PollEvents", opGlfwPollEvents, nil, nil)
	RegisterOpCode(OP_GLFW_SWAP_BUFFERS, "glfw.SwapBuffers", opGlfwSwapBuffers, http.In(ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_GET_FRAMEBUFFER_SIZE, "glfw.GetFramebufferSize", opGlfwGetFramebufferSize, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLFW_GET_WINDOW_POS, "glfw.GetWindowPos", opGlfwGetWindowPos, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLFW_GET_WINDOW_SIZE, "glfw.GetWindowSize", opGlfwGetWindowSize, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLFW_GET_TIME, "glfw.GetTime", opGlfwGetTime, nil, http.Out(ast.ConstCxArg_F64))
	RegisterOpCode(OP_GLFW_SWAP_INTERVAL, "glfw.SwapInterval", opGlfwSwapInterval, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GLFW_SET_START_CALLBACK, "glfw.SetStartCallback", opGlfwSetStartCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_SET_STOP_CALLBACK, "glfw.SetStopCallback", opGlfwSetStopCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_SET_KEYBOARD_CALLBACK, "glfw.SetKeyboardCallback", opGlfwSetKeyboardCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_SET_MOUSE_CALLBACK, "glfw.SetMouseCallback", opGlfwSetMouseCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_SET_FRAMEBUFFER_SIZE_CALLBACK, "glfw.SetFramebufferSizeCallback", opGlfwSetFramebufferSizeCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_SET_WINDOW_POS_CALLBACK, "glfw.SetWindowPosCallback", opGlfwSetWindowPosCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_SET_WINDOW_SIZE_CALLBACK, "glfw.SetWindowSizeCallback", opGlfwSetWindowSizeCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLFW_SET_KEY_CALLBACK, "glfw.SetKeyCallback", opGlfwSetKeyCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)                          // TODO : to deprecate
	RegisterOpCode(OP_GLFW_SET_MOUSE_BUTTON_CALLBACK, "glfw.SetMouseButtonCallback", opGlfwSetMouseButtonCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil) // TODO : to deprecate
	RegisterOpCode(OP_GLFW_SET_CURSOR_POS_CALLBACK, "glfw.SetCursorPosCallback", opGlfwSetCursorPosCallback, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), nil)       // TODO : to deprecate
	RegisterOpCode(OP_GLFW_GET_CURSOR_POS, "glfw.GetCursorPos", opGlfwGetCursorPos, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_F64, ast.ConstCxArg_F64))
	RegisterOpCode(OP_GLFW_SET_INPUT_MODE, "glfw.SetInputMode", opGlfwSetInputMode, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GLFW_SET_WINDOW_POS, "glfw.SetWindowPos", opGlfwSetWindowPos, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GLFW_GET_KEY, "glfw.GetKey", opGlfwGetKey, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLFW_FUNC_I32_I32, "glfw.func_i32_i32", opGlfwFuncI32I32, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLFW_CALL_I32_I32, "glfw.call_i32_i32", opGlfwCallI32I32, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GLFW_GET_WINDOW_CONTENT_SCALE, "glfw.GetWindowContentScale", opGlfwGetWindowContentScale, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLFW_GET_MONITOR_CONTENT_SCALE, "glfw.GetMonitorContentScale", opGlfwGetMonitorContentScale, nil, http.Out(ast.ConstCxArg_F32, ast.ConstCxArg_F32))

	// gltext
	RegisterOpCode(OP_GLTEXT_LOAD_TRUE_TYPE, "gltext.LoadTrueType", opGltextLoadTrueType, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GLTEXT_LOAD_TRUE_TYPE_CORE, "gltext.LoadTrueTypeCore", opGltextLoadTrueTypeCore, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_GLTEXT_PRINTF, "gltext.Printf", opGltextPrintf, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_F32, ast.ConstCxArg_F32, ast.ConstCxArg_STR), nil)
	RegisterOpCode(OP_GLTEXT_METRICS, "gltext.Metrics", opGltextMetrics, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLTEXT_TEXTURE, "gltext.Texture", opGltextTexture, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLTEXT_NEXT_GLYPH, "gltext.NextGlyph", opGltextNextGlyph, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLTEXT_GLYPH_BOUNDS, "gltext.GlyphBounds", opGltextGlyphBounds, nil, http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLTEXT_GLYPH_METRICS, "gltext.GlyphMetrics", opGltextGlyphMetrics, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32))
	RegisterOpCode(OP_GLTEXT_GLYPH_INFO, "gltext.GlyphInfo", opGltextGlyphInfo, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32))

	// goal
	RegisterOpCode(OP_AL_LOAD_WAV, "al.LoadWav", opAlLoadWav, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_I64, ast.Slice(constants.TYPE_UI8)))

	// openal
	RegisterOpCode(OP_AL_CLOSE_DEVICE, "al.CloseDevice", opAlCloseDevice, nil, nil)
	RegisterOpCode(OP_AL_DELETE_BUFFERS, "al.DeleteBuffers", opAlDeleteBuffers, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_AL_DELETE_SOURCES, "al.DeleteSources", opAlDeleteSources, http.In(ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_AL_DEVICE_ERROR, "al.DeviceError", opAlDeviceError, nil, http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_AL_ERROR, "al.Error", opAlError, nil, http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_AL_EXTENSIONS, "al.Extensions", opAlExtensions, nil, http.Out(ast.ConstCxArg_STR))
	RegisterOpCode(OP_AL_OPEN_DEVICE, "al.OpenDevice", opAlOpenDevice, nil, nil)
	RegisterOpCode(OP_AL_PAUSE_SOURCES, "al.PauseSources", opAlPauseSources, http.In(ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_AL_PLAY_SOURCES, "al.PlaySources", opAlPlaySources, http.In(ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_AL_RENDERER, "al.Renderer", opAlRenderer, nil, http.Out(ast.ConstCxArg_STR))
	RegisterOpCode(OP_AL_REWIND_SOURCES, "al.RewindSources", opAlRewindSources, http.In(ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_AL_STOP_SOURCES, "al.StopSources", opAlStopSources, http.In(ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_AL_VENDOR, "al.Vendor", opAlVendor, nil, http.Out(ast.ConstCxArg_STR))
	RegisterOpCode(OP_AL_VERSION, "al.Version", opAlVersion, nil, http.Out(ast.ConstCxArg_STR))
	RegisterOpCode(OP_AL_GEN_BUFFERS, "al.GenBuffers", opAlGenBuffers, http.In(ast.ConstCxArg_I32), http.Out(ast.Slice(constants.TYPE_I32)))
	RegisterOpCode(OP_AL_BUFFER_DATA, "al.BufferData", opAlBufferData, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_UND_TYPE, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_AL_GEN_SOURCES, "al.GenSources", opAlGenSources, http.In(ast.ConstCxArg_I32), http.Out(ast.Slice(constants.TYPE_I32)))
	RegisterOpCode(OP_AL_SOURCE_BUFFERS_PROCESSED, "al.SourceBuffersProcessed", opAlSourceBuffersProcessed, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_AL_SOURCE_BUFFERS_QUEUED, "al.SourceBuffersQueued", opAlSourceBuffersQueued, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_AL_SOURCE_QUEUE_BUFFERS, "al.SourceQueueBuffers", opAlSourceQueueBuffers, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32)), nil)
	RegisterOpCode(OP_AL_SOURCE_STATE, "al.SourceState", opAlSourceState, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_AL_SOURCE_UNQUEUE_BUFFERS, "al.SourceUnqueueBuffers", opAlSourceUnqueueBuffers, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32)), nil)
}
