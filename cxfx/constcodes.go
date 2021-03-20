// +build cxfx

package cxfx

import (
	"github.com/skycoin/cx/cx"
)

// constant codes
// nolint golint
const (
	// https://www.khronos.org/registry/OpenGL/api/GL/glcorearb.h
	// gl_1_0
	CONST_GL_DEPTH_BUFFER_BIT = iota + 0x1FFFF
	CONST_GL_STENCIL_BUFFER_BIT
	CONST_GL_COLOR_BUFFER_BIT
	CONST_GL_FALSE
	CONST_GL_TRUE
	CONST_GL_POINTS
	CONST_GL_LINES
	CONST_GL_LINE_LOOP
	CONST_GL_LINE_STRIP
	CONST_GL_TRIANGLES
	CONST_GL_TRIANGLE_STRIP
	CONST_GL_TRIANGLE_FAN
	CONST_GL_QUADS
	CONST_GL_NEVER
	CONST_GL_LESS
	CONST_GL_EQUAL
	CONST_GL_LEQUAL
	CONST_GL_GREATER
	CONST_GL_NOTEQUAL
	CONST_GL_GEQUAL
	CONST_GL_ALWAYS
	CONST_GL_ZERO
	CONST_GL_ONE
	CONST_GL_SRC_ALPHA
	CONST_GL_ONE_MINUS_SRC_ALPHA
	CONST_GL_NONE
	CONST_GL_FRONT
	CONST_GL_BACK
	CONST_GL_FRONT_AND_BACK
	CONST_GL_NO_ERROR
	CONST_GL_INVALID_ENUM
	CONST_GL_INVALID_VALUE
	CONST_GL_INVALID_OPERATION
	CONST_GL_STACK_OVERFLOW
	CONST_GL_STACK_UNDERFLOW
	CONST_GL_OUT_OF_MEMORY
	CONST_GL_CW
	CONST_GL_CCW
	CONST_GL_LINE_SMOOTH
	CONST_GL_POLYGON_SMOOTH
	CONST_GL_CULL_FACE
	CONST_GL_CULL_FACE_MODE
	CONST_GL_FRONT_FACE
	CONST_GL_DEPTH_RANGE
	CONST_GL_DEPTH_TEST
	CONST_GL_DEPTH_WRITEMASK
	CONST_GL_DEPTH_CLEAR_VALUE
	CONST_GL_DEPTH_FUNC
	CONST_GL_STENCIL_TEST
	CONST_GL_STENCIL_CLEAR_VALUE
	CONST_GL_STENCIL_FUNC
	CONST_GL_STENCIL_VALUE_MASK
	CONST_GL_STENCIL_FAIL
	CONST_GL_STENCIL_PASS_DEPTH_FAIL
	CONST_GL_STENCIL_PASS_DEPTH_PASS
	CONST_GL_STENCIL_REF
	CONST_GL_STENCIL_WRITEMASK
	CONST_GL_DITHER
	CONST_GL_BLEND
	CONST_GL_SCISSOR_TEST
	CONST_GL_POLYGON_SMOOTH_HINT
	CONST_GL_TEXTURE_2D
	CONST_GL_TEXTURE_WIDTH
	CONST_GL_TEXTURE_HEIGHT
	CONST_GL_DONT_CARE
	CONST_GL_BYTE
	CONST_GL_UNSIGNED_BYTE
	CONST_GL_SHORT
	CONST_GL_UNSIGNED_SHORT
	CONST_GL_INT
	CONST_GL_UNSIGNED_INT
	CONST_GL_FLOAT
	CONST_GL_INVERT
	CONST_GL_TEXTURE
	CONST_GL_COLOR
	CONST_GL_DEPTH
	CONST_GL_STENCIL
	CONST_GL_STENCIL_INDEX
	CONST_GL_DEPTH_COMPONENT
	CONST_GL_RED
	CONST_GL_RGB
	CONST_GL_RGBA
	CONST_GL_KEEP
	CONST_GL_REPLACE
	CONST_GL_INCR
	CONST_GL_DECR
	CONST_GL_NEAREST
	CONST_GL_LINEAR
	CONST_GL_NEAREST_MIPMAP_NEAREST
	CONST_GL_LINEAR_MIPMAP_NEAREST
	CONST_GL_NEAREST_MIPMAP_LINEAR
	CONST_GL_LINEAR_MIPMAP_LINEAR
	CONST_GL_TEXTURE_MAG_FILTER
	CONST_GL_TEXTURE_MIN_FILTER
	CONST_GL_TEXTURE_WRAP_S
	CONST_GL_TEXTURE_WRAP_T
	CONST_GL_REPEAT

	// gl_1_1
	CONST_GL_RGBA8
	CONST_GL_VERTEX_ARRAY

	// gl_1_2
	CONST_GL_TEXTURE_WRAP_R
	CONST_GL_CLAMP_TO_EDGE

	// gl_1_3
	CONST_GL_TEXTURE0
	CONST_GL_MULTISAMPLE_ARB // remove _ARB
	CONST_GL_TEXTURE_CUBE_MAP
	CONST_GL_TEXTURE_CUBE_MAP_POSITIVE_X
	CONST_GL_TEXTURE_CUBE_MAP_NEGATIVE_X
	CONST_GL_TEXTURE_CUBE_MAP_POSITIVE_Y
	CONST_GL_TEXTURE_CUBE_MAP_NEGATIVE_Y
	CONST_GL_TEXTURE_CUBE_MAP_POSITIVE_Z
	CONST_GL_TEXTURE_CUBE_MAP_NEGATIVE_Z
	CONST_GL_CLAMP_TO_BORDER

	// gl_1_4
	CONST_GL_DEPTH_COMPONENT16
	CONST_GL_DEPTH_COMPONENT24
	CONST_GL_DEPTH_COMPONENT32
	CONST_GL_MIRRORED_REPEAT
	CONST_GL_INCR_WRAP
	CONST_GL_DECR_WRAP

	// gl_1_5
	CONST_GL_ARRAY_BUFFER
	CONST_GL_ELEMENT_ARRAY_BUFFER
	CONST_GL_STREAM_DRAW
	CONST_GL_STREAM_READ
	CONST_GL_STREAM_COPY
	CONST_GL_STATIC_DRAW
	CONST_GL_STATIC_READ
	CONST_GL_STATIC_COPY
	CONST_GL_DYNAMIC_DRAW
	CONST_GL_DYNAMIC_READ
	CONST_GL_DYNAMIC_COPY

	// gl_2_0
	CONST_GL_STENCIL_BACK_FUNC
	CONST_GL_STENCIL_BACK_FAIL
	CONST_GL_STENCIL_BACK_PASS_DEPTH_FAIL
	CONST_GL_STENCIL_BACK_PASS_DEPTH_PASS
	CONST_GL_FRAGMENT_SHADER
	CONST_GL_VERTEX_SHADER
	CONST_GL_COMPILE_STATUS
	CONST_GL_LINK_STATUS
	CONST_GL_INFO_LOG_LENGTH
	CONST_GL_STENCIL_BACK_REF
	CONST_GL_STENCIL_BACK_VALUE_MASK
	CONST_GL_STENCIL_BACK_WRITEMASK

	// gl_3_0
	CONST_GL_RGBA16F
	CONST_GL_RGB16F
	CONST_GL_DEPTH_COMPONENT32F
	CONST_GL_DEPTH32F_STENCIL8
	CONST_GL_FRAMEBUFFER_UNDEFINED
	CONST_GL_DEPTH_STENCIL_ATTACHMENT
	CONST_GL_DEPTH_STENCIL
	CONST_GL_UNSIGNED_INT_24_8
	CONST_GL_DEPTH24_STENCIL8
	CONST_GL_FRAMEBUFFER_COMPLETE
	CONST_GL_FRAMEBUFFER_INCOMPLETE_ATTACHMENT
	CONST_GL_FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT
	CONST_GL_FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER
	CONST_GL_FRAMEBUFFER_INCOMPLETE_READ_BUFFER
	CONST_GL_FRAMEBUFFER_UNSUPPORTED
	CONST_GL_COLOR_ATTACHMENT0
	CONST_GL_DEPTH_ATTACHMENT
	CONST_GL_STENCIL_ATTACHMENT
	CONST_GL_FRAMEBUFFER
	CONST_GL_RENDERBUFFER
	CONST_GL_FRAMEBUFFER_INCOMPLETE_MULTISAMPLE
	CONST_GL_STENCIL_INDEX1
	CONST_GL_STENCIL_INDEX4
	CONST_GL_STENCIL_INDEX8
	CONST_GL_STENCIL_INDEX16
	CONST_GL_HALF_FLOAT
	CONST_GL_R8

	// gl_3_2
	CONST_GL_FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS

	// gl_4_4
	CONST_GL_MIRROR_CLAMP_TO_EDGE

	// Fixed pipeline. Deprecated ?
	CONST_GL_POLYGON
	CONST_GL_MODELVIEW
	CONST_GL_PROJECTION
	CONST_GL_MODELVIEW_MATRIX
	CONST_GL_LIGHTING
	CONST_GL_LIGHT0
	CONST_GL_AMBIENT
	CONST_GL_DIFFUSE
	CONST_GL_POSITION
	CONST_GL_TEXTURE_ENV
	CONST_GL_TEXTURE_ENV_MODE
	CONST_GL_MODULATE
	CONST_GL_DECAL
	CONST_GL_POINT_SMOOTH

	// glfw
	CONST_GLFW_FALSE
	CONST_GLFW_TRUE
	CONST_GLFW_PRESS
	CONST_GLFW_RELEASE
	CONST_GLFW_REPEAT
	CONST_GLFW_KEY_UNKNOWN
	CONST_GLFW_CURSOR
	CONST_GLFW_STICKY_KEYS
	CONST_GLFW_STICKY_MOUSE_BUTTONS
	CONST_GLFW_CURSOR_NORMAL
	CONST_GLFW_CURSOR_HIDDEN
	CONST_GLFW_CURSOR_DISABLED
	CONST_GLFW_RESIZABLE
	CONST_GLFW_CONTEXT_VERSION_MAJOR
	CONST_GLFW_CONTEXT_VERSION_MINOR
	CONST_GLFW_OPENGL_PROFILE
	CONST_GLFW_OPENGL_COREPROFILE
	CONST_GLFW_OPENGL_FORWARD_COMPATIBLE
	CONST_GLFW_MOUSE_BUTTON_LAST
	CONST_GLFW_MOUSE_BUTTON_LEFT
	CONST_GLFW_MOUSE_BUTTON_RIGHT
	CONST_GLFW_MOUSE_BUTTON_MIDDLE
	CONST_GLFW_COCOA_RETINA_FRAMEBUFFER
	CONST_GLFW_SCALE_TO_MONITOR

	// gltext
	CONST_GLTEXT_LEFT_TO_RIGHT
	CONST_GLTEXT_RIGHT_TO_LEFT
	CONST_GLTEXT_TOP_TO_BOTTOM
	// openal
	CONST_AL_INVALID_NAME
	CONST_AL_INVALID_ENUM
	CONST_AL_INVALID_VALUE
	CONST_AL_INVALID_OPERATION
	CONST_AL_OUT_OF_MEMORY
	CONST_AL_INVERSE_DISTANCE
	CONST_AL_INVERSE_DISTANCE_CLAMPED
	CONST_AL_LINEAR_DISTANCE
	CONST_AL_LINEAR_DISTANCE_CLAMPED
	CONST_AL_EXPONENT_DISTANCE
	CONST_AL_EXPONENT_DISTANCE_CLAMPED
	CONST_AL_INITIAL
	CONST_AL_PLAYING
	CONST_AL_PAUSED
	CONST_AL_STOPPED
	CONST_AL_FORMAT_MONO_8
	CONST_AL_FORMAT_MONO_16
	CONST_AL_FORMAT_STEREO_8
	CONST_AL_FORMAT_STEREO_16
	CONST_AL_WAV_FORMAT_PCM
	CONST_AL_WAV_FORMAT_IEEE_FLOAT
)

// nolint typecheck
func init() {
	/* gl_1_0 */
	cxcore.ConstI32(CONST_GL_DEPTH_BUFFER_BIT, "gl.DEPTH_BUFFER_BIT", 0x00000100)
	cxcore.ConstI32(CONST_GL_STENCIL_BUFFER_BIT, "gl.STENCIL_BUFFER_BIT", 0x00000400)
	cxcore.ConstI32(CONST_GL_COLOR_BUFFER_BIT, "gl.COLOR_BUFFER_BIT", 0x00004000)
	cxcore.ConstI32(CONST_GL_FALSE, "gl.FALSE", 0)
	cxcore.ConstI32(CONST_GL_TRUE, "gl.TRUE", 1)
	cxcore.ConstI32(CONST_GL_POINTS, "gl.POINTS", 0x0000)
	cxcore.ConstI32(CONST_GL_LINES, "gl.LINES", 0x0001)
	cxcore.ConstI32(CONST_GL_LINE_LOOP, "gl.LINE_LOOP", 0x0002)
	cxcore.ConstI32(CONST_GL_LINE_STRIP, "gl.LINE_STRIP", 0x0003)
	cxcore.ConstI32(CONST_GL_TRIANGLES, "gl.TRIANGLES", 0x0004)
	cxcore.ConstI32(CONST_GL_TRIANGLE_STRIP, "gl.TRIANGLE_STRIP", 0x0005)
	cxcore.ConstI32(CONST_GL_TRIANGLE_FAN, "gl.TRIANGLE_FAN", 0x0006)
	cxcore.ConstI32(CONST_GL_QUADS, "gl.QUADS", 0x0007)
	cxcore.ConstI32(CONST_GL_NEVER, "gl.NEVER", 0x0200)
	cxcore.ConstI32(CONST_GL_LESS, "gl.LESS", 0x0201)
	cxcore.ConstI32(CONST_GL_EQUAL, "gl.EQUAL", 0x0202)
	cxcore.ConstI32(CONST_GL_LEQUAL, "gl.LEQUAL", 0x0203)
	cxcore.ConstI32(CONST_GL_GREATER, "gl.GREATER", 0x0204)
	cxcore.ConstI32(CONST_GL_NOTEQUAL, "gl.NOTEQUAL", 0x0205)
	cxcore.ConstI32(CONST_GL_GEQUAL, "gl.GEQUAL", 0x0206)
	cxcore.ConstI32(CONST_GL_ALWAYS, "gl.ALWAYS", 0x0207)
	cxcore.ConstI32(CONST_GL_ZERO, "gl.ZERO", 0)
	cxcore.ConstI32(CONST_GL_ONE, "gl.ONE", 1)
	cxcore.ConstI32(CONST_GL_SRC_ALPHA, "gl.SRC_ALPHA", 0x302)
	cxcore.ConstI32(CONST_GL_ONE_MINUS_SRC_ALPHA, "gl.ONE_MINUS_SRC_ALPHA", 0x303)
	cxcore.ConstI32(CONST_GL_NONE, "gl.NONE", 0)
	cxcore.ConstI32(CONST_GL_FRONT, "gl.FRONT", 0x404)
	cxcore.ConstI32(CONST_GL_BACK, "gl.BACK", 0x405)
	cxcore.ConstI32(CONST_GL_FRONT_AND_BACK, "gl.FRONT_AND_BACK", 0x408)
	cxcore.ConstI32(CONST_GL_NO_ERROR, "gl.NO_ERROR", 0)
	cxcore.ConstI32(CONST_GL_INVALID_ENUM, "gl.INVALID_ENUM", 0x500)
	cxcore.ConstI32(CONST_GL_INVALID_VALUE, "gl.INVALID_VALUE", 0x501)
	cxcore.ConstI32(CONST_GL_INVALID_OPERATION, "gl.INVALID_OPERATION", 0x502)
	cxcore.ConstI32(CONST_GL_STACK_OVERFLOW, "gl.STACK_OVERFLOW", 0x503)
	cxcore.ConstI32(CONST_GL_STACK_UNDERFLOW, "gl.STACK_UNDERFLOW", 0x504)
	cxcore.ConstI32(CONST_GL_OUT_OF_MEMORY, "gl.OUT_OF_MEMORY", 0x505)
	cxcore.ConstI32(CONST_GL_CW, "gl.CW", 0x0900)
	cxcore.ConstI32(CONST_GL_CCW, "gl.CCW", 0x0901)
	cxcore.ConstI32(CONST_GL_LINE_SMOOTH, "gl.LINE_SMOOTH", 0x0B20)
	cxcore.ConstI32(CONST_GL_POLYGON_SMOOTH, "gl.POLYGON_SMOOTH", 0x0B41)
	cxcore.ConstI32(CONST_GL_CULL_FACE, "gl.CULL_FACE", 0x0B44)
	cxcore.ConstI32(CONST_GL_CULL_FACE_MODE, "gl.CULL_FACE_MODE", 0x0B45)
	cxcore.ConstI32(CONST_GL_FRONT_FACE, "gl.FRONT_FACE", 0x0B46)
	cxcore.ConstI32(CONST_GL_DEPTH_RANGE, "gl.DEPTH_RANGE", 0x0B70)
	cxcore.ConstI32(CONST_GL_DEPTH_TEST, "gl.DEPTH_TEST", 0x0B71)
	cxcore.ConstI32(CONST_GL_DEPTH_WRITEMASK, "gl.DEPTH_WRITEMASK", 0x0B72)
	cxcore.ConstI32(CONST_GL_DEPTH_CLEAR_VALUE, "gl.DEPTH_CLEAR_VALUE", 0x0B73)
	cxcore.ConstI32(CONST_GL_DEPTH_FUNC, "gl.DEPTH_FUNC", 0x0B74)
	cxcore.ConstI32(CONST_GL_STENCIL_TEST, "gl.STENCIL_TEST", 0x0B90)
	cxcore.ConstI32(CONST_GL_STENCIL_CLEAR_VALUE, "gl.STENCIL_CLEAR_VALUE", 0x0B91)
	cxcore.ConstI32(CONST_GL_STENCIL_FUNC, "gl.STENCIL_FUNC", 0x0B92)
	cxcore.ConstI32(CONST_GL_STENCIL_VALUE_MASK, "gl.STENCIL_VALUE_MASK", 0x0B93)
	cxcore.ConstI32(CONST_GL_STENCIL_FAIL, "gl.STENCIL_FAIL", 0x0B94)
	cxcore.ConstI32(CONST_GL_STENCIL_PASS_DEPTH_FAIL, "gl.STENCIL_PASS_DEPTH_FAIL", 0x0B95)
	cxcore.ConstI32(CONST_GL_STENCIL_PASS_DEPTH_PASS, "gl.STENCIL_PASS_DEPTH_PASS", 0x0B96)
	cxcore.ConstI32(CONST_GL_STENCIL_REF, "gl.STENCIL_REF", 0x0B97)
	cxcore.ConstI32(CONST_GL_STENCIL_WRITEMASK, "gl.STENCIL_WRITE_MASK", 0x0B98)
	cxcore.ConstI32(CONST_GL_DITHER, "gl.DITHER", 0x0BD0)
	cxcore.ConstI32(CONST_GL_BLEND, "gl.BLEND", 0x0BE2)
	cxcore.ConstI32(CONST_GL_SCISSOR_TEST, "gl.SCISSOR_TEST", 0x0C11)
	cxcore.ConstI32(CONST_GL_POLYGON_SMOOTH_HINT, "gl.POLYGON_SMOOTH_HINT", 0x0C53)
	cxcore.ConstI32(CONST_GL_TEXTURE_2D, "gl.TEXTURE_2D", 0x0DE1)
	cxcore.ConstI32(CONST_GL_TEXTURE_WIDTH, "gl.TEXTURE_WIDTH", 0x1000)
	cxcore.ConstI32(CONST_GL_TEXTURE_HEIGHT, "gl.TEXTURE_HEIGHT", 0x1001)
	cxcore.ConstI32(CONST_GL_DONT_CARE, "gl.DONT_CARE", 0x1100)
	cxcore.ConstI32(CONST_GL_BYTE, "gl.BYTE", 0x1400)
	cxcore.ConstI32(CONST_GL_UNSIGNED_BYTE, "gl.UNSIGNED_BYTE", 0x1401)
	cxcore.ConstI32(CONST_GL_SHORT, "gl.SHORT", 0x1402)
	cxcore.ConstI32(CONST_GL_UNSIGNED_SHORT, "gl.UNSIGNED_SHORT", 0x1403)
	cxcore.ConstI32(CONST_GL_INT, "gl.INT", 0x1404)
	cxcore.ConstI32(CONST_GL_UNSIGNED_INT, "gl.UNSIGNED_INT", 0x1405)
	cxcore.ConstI32(CONST_GL_FLOAT, "gl.FLOAT", 0x1406)
	cxcore.ConstI32(CONST_GL_INVERT, "gl.INVERT", 0x150A)
	cxcore.ConstI32(CONST_GL_TEXTURE, "gl.TEXTURE", 0x1702)
	cxcore.ConstI32(CONST_GL_COLOR, "gl.COLOR", 0x1800)
	cxcore.ConstI32(CONST_GL_DEPTH, "gl.DEPTH", 0x1801)
	cxcore.ConstI32(CONST_GL_STENCIL, "gl.STENCIL", 0x1802)
	cxcore.ConstI32(CONST_GL_STENCIL_INDEX, "gl.STENCIL_INDEX", 0x1901)
	cxcore.ConstI32(CONST_GL_DEPTH_COMPONENT, "gl.DEPTH_COMPONENT", 0x1902)
	cxcore.ConstI32(CONST_GL_RED, "gl.RED", 0x1903)
	cxcore.ConstI32(CONST_GL_RGB, "gl.RGB", 0x1907)
	cxcore.ConstI32(CONST_GL_RGBA, "gl.RGBA", 0x1908)
	cxcore.ConstI32(CONST_GL_KEEP, "gl.KEEP", 0x1E00)
	cxcore.ConstI32(CONST_GL_REPLACE, "gl.REPLACE", 0x1E01)
	cxcore.ConstI32(CONST_GL_INCR, "gl.INCR", 0x1E02)
	cxcore.ConstI32(CONST_GL_DECR, "gl.DECR", 0x1E03)
	cxcore.ConstI32(CONST_GL_NEAREST, "gl.NEAREST", 0x2600)
	cxcore.ConstI32(CONST_GL_LINEAR, "gl.LINEAR", 0x2601)
	cxcore.ConstI32(CONST_GL_NEAREST_MIPMAP_NEAREST, "gl.NEAREST_MIPMAP_NEAREST", 0x2700)
	cxcore.ConstI32(CONST_GL_LINEAR_MIPMAP_NEAREST, "gl.LINEAR_MIPMAP_NEAREST", 0x2701)
	cxcore.ConstI32(CONST_GL_NEAREST_MIPMAP_LINEAR, "gl.NEAREST_MIPMAP_LINEAR", 0x2702)
	cxcore.ConstI32(CONST_GL_LINEAR_MIPMAP_LINEAR, "gl.LINEAR_MIPMAP_LINEAR", 0x2703)
	cxcore.ConstI32(CONST_GL_TEXTURE_MAG_FILTER, "gl.TEXTURE_MAG_FILTER", 0x2800)
	cxcore.ConstI32(CONST_GL_TEXTURE_MIN_FILTER, "gl.TEXTURE_MIN_FILTER", 0x2801)
	cxcore.ConstI32(CONST_GL_TEXTURE_WRAP_S, "gl.TEXTURE_WRAP_S", 0x2802)
	cxcore.ConstI32(CONST_GL_TEXTURE_WRAP_T, "gl.TEXTURE_WRAP_T", 0x2803)
	cxcore.ConstI32(CONST_GL_REPEAT, "gl.REPEAT", 0x2901)

	// gl_1_1
	cxcore.ConstI32(CONST_GL_RGBA8, "gl.RGBA8", 0x8058)
	cxcore.ConstI32(CONST_GL_VERTEX_ARRAY, "gl.VERTEX_ARRAY", 0x8074)

	// gl_1_2
	cxcore.ConstI32(CONST_GL_TEXTURE_WRAP_R, "gl.TEXTURE_WRAP_R", 0x8072)
	cxcore.ConstI32(CONST_GL_CLAMP_TO_EDGE, "gl.CLAMP_TO_EDGE", 0x812F)

	// gl_1_3
	cxcore.ConstI32(CONST_GL_TEXTURE0, "gl.TEXTURE0", 0x84C0)
	cxcore.ConstI32(CONST_GL_MULTISAMPLE_ARB, "gl.MULTISAMPLE_ARB", 0x809D) // remove _ARB
	cxcore.ConstI32(CONST_GL_TEXTURE_CUBE_MAP, "gl.TEXTURE_CUBE_MAP", 0x8513)
	cxcore.ConstI32(CONST_GL_TEXTURE_CUBE_MAP_POSITIVE_X, "gl.TEXTURE_CUBE_MAP_POSITIVE_X", 0x8515)
	cxcore.ConstI32(CONST_GL_TEXTURE_CUBE_MAP_NEGATIVE_X, "gl.TEXTURE_CUBE_MAP_NEGATIVE_X", 0x8516)
	cxcore.ConstI32(CONST_GL_TEXTURE_CUBE_MAP_POSITIVE_Y, "gl.TEXTURE_CUBE_MAP_POSITIVE_Y", 0x8517)
	cxcore.ConstI32(CONST_GL_TEXTURE_CUBE_MAP_NEGATIVE_Y, "gl.TEXTURE_CUBE_MAP_NEGATIVE_Y", 0x8518)
	cxcore.ConstI32(CONST_GL_TEXTURE_CUBE_MAP_POSITIVE_Z, "gl.TEXTURE_CUBE_MAP_POSITIVE_Z", 0x8519)
	cxcore.ConstI32(CONST_GL_TEXTURE_CUBE_MAP_NEGATIVE_Z, "gl.TEXTURE_CUBE_MAP_NEGATIVE_Z", 0x851A)
	cxcore.ConstI32(CONST_GL_CLAMP_TO_BORDER, "gl.CLAMP_TO_BORDER", 0x812D)

	// gl_1_4
	cxcore.ConstI32(CONST_GL_DEPTH_COMPONENT16, "gl.DEPTH_COMPONENT16", 0x81A5)
	cxcore.ConstI32(CONST_GL_DEPTH_COMPONENT24, "gl.DEPTH_COMPONENT24", 0x81A6)
	cxcore.ConstI32(CONST_GL_DEPTH_COMPONENT32, "gl.DEPTH_COMPONENT32", 0x81A7)
	cxcore.ConstI32(CONST_GL_MIRRORED_REPEAT, "gl.MIRRORED_REPEAT", 0x8370)
	cxcore.ConstI32(CONST_GL_INCR_WRAP, "gl.INCR_WRAP", 0x8507)
	cxcore.ConstI32(CONST_GL_DECR_WRAP, "gl.DECR_WRAP", 0x8508)

	// gl_1_5
	cxcore.ConstI32(CONST_GL_ARRAY_BUFFER, "gl.ARRAY_BUFFER", 0x8892)
	cxcore.ConstI32(CONST_GL_ELEMENT_ARRAY_BUFFER, "gl.ELEMENT_ARRAY_BUFFER", 0x8893)
	cxcore.ConstI32(CONST_GL_STREAM_DRAW, "gl.STREAM_DRAW", 0x88E0)
	cxcore.ConstI32(CONST_GL_STREAM_READ, "gl.STREAM_READ", 0x88E1)
	cxcore.ConstI32(CONST_GL_STREAM_COPY, "gl.STREAM_COPY", 0x88E2)
	cxcore.ConstI32(CONST_GL_STATIC_DRAW, "gl.STATIC_DRAW", 0x88E4)
	cxcore.ConstI32(CONST_GL_STATIC_READ, "gl.STATIC_READ", 0x88E5)
	cxcore.ConstI32(CONST_GL_STATIC_COPY, "gl.STATIC_COPY", 0x88E6)
	cxcore.ConstI32(CONST_GL_DYNAMIC_DRAW, "gl.DYNAMIC_DRAW", 0x88E8)
	cxcore.ConstI32(CONST_GL_DYNAMIC_READ, "gl.DYNAMIC_READ", 0x88E9)
	cxcore.ConstI32(CONST_GL_DYNAMIC_COPY, "gl.DYNAMIC_COPY", 0x88EA)

	// gl_2_0
	cxcore.ConstI32(CONST_GL_STENCIL_BACK_FUNC, "gl.STENCIL_BACK_FUNC", 0x8800)
	cxcore.ConstI32(CONST_GL_STENCIL_BACK_FAIL, "gl.STENCIL_BACK_FAIL", 0x8801)
	cxcore.ConstI32(CONST_GL_STENCIL_BACK_PASS_DEPTH_FAIL, "gl.STENCIL_BACK_PASS_DEPTH_FAIL", 0x8802)
	cxcore.ConstI32(CONST_GL_STENCIL_BACK_PASS_DEPTH_PASS, "gl.STENCIL_BACK_PASS_DEPTH_PASS", 0x8803)
	cxcore.ConstI32(CONST_GL_FRAGMENT_SHADER, "gl.FRAGMENT_SHADER", 0x8B30)
	cxcore.ConstI32(CONST_GL_VERTEX_SHADER, "gl.VERTEX_SHADER", 0x8B31)
	cxcore.ConstI32(CONST_GL_COMPILE_STATUS, "gl.COMPILE_STATUS", 0x8B81)
	cxcore.ConstI32(CONST_GL_LINK_STATUS, "gl.LINK_STATUS", 0x8B82)
	cxcore.ConstI32(CONST_GL_INFO_LOG_LENGTH, "gl.INFO_LOG_LENGTH", 0x8B84)
	cxcore.ConstI32(CONST_GL_STENCIL_BACK_REF, "gl.STENCIL_BACK_REF", 0x8CA3)
	cxcore.ConstI32(CONST_GL_STENCIL_BACK_VALUE_MASK, "gl.STENCIL_BACK_VALUE_MASK", 0x8CA4)
	cxcore.ConstI32(CONST_GL_STENCIL_BACK_WRITEMASK, "gl.STENCIL_BACK_WRITEMASK", 0x8CA5)

	// gl_3_0
	cxcore.ConstI32(CONST_GL_RGBA16F, "gl.RGBA16F", 0x881A)
	cxcore.ConstI32(CONST_GL_RGB16F, "gl.RGB16F", 0x881B)
	cxcore.ConstI32(CONST_GL_DEPTH_COMPONENT32F, "gl.DEPTH_COMPONENT32F", 0x8CAC)
	cxcore.ConstI32(CONST_GL_DEPTH32F_STENCIL8, "gl.DEPTH32F_STENCIL8", 0x8CAD)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_UNDEFINED, "gl.FRAMEBUFFER_UNDEFINED", 0x8219)
	cxcore.ConstI32(CONST_GL_DEPTH_STENCIL_ATTACHMENT, "gl.DEPTH_STENCIL_ATTACHMENT", 0x821A)
	cxcore.ConstI32(CONST_GL_DEPTH_STENCIL, "gl.DEPTH_STENCIL", 0x84F9)
	cxcore.ConstI32(CONST_GL_UNSIGNED_INT_24_8, "gl.UNSIGNED_INT_24_8", 0x84FA)
	cxcore.ConstI32(CONST_GL_DEPTH24_STENCIL8, "gl.DEPTH24_STENCIL8", 0x88F0)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_COMPLETE, "gl.FRAMEBUFFER_COMPLETE", 0x8CD5)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_INCOMPLETE_ATTACHMENT, "gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT", 0x8CD6)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT, "gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT", 0x8CD7)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER, "gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER", 0x8CDB)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_INCOMPLETE_READ_BUFFER, "gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER", 0x8CDC)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_UNSUPPORTED, "gl.FRAMEBUFFER_UNSUPPORTED", 0x8CDD)
	cxcore.ConstI32(CONST_GL_COLOR_ATTACHMENT0, "gl.COLOR_ATTACHMENT0", 0x8CE0)
	cxcore.ConstI32(CONST_GL_DEPTH_ATTACHMENT, "gl.DEPTH_ATTACHMENT", 0x8D00)
	cxcore.ConstI32(CONST_GL_STENCIL_ATTACHMENT, "gl.STENCIL_ATTACHMENT", 0x8D20)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER, "gl.FRAMEBUFFER", 0x8D40)
	cxcore.ConstI32(CONST_GL_RENDERBUFFER, "gl.RENDERBUFFER", 0x8D41)
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_INCOMPLETE_MULTISAMPLE, "gl.FRAMEBUFFER_INCOMPLETE_MULTISAMPLE", 0x8D56)
	cxcore.ConstI32(CONST_GL_STENCIL_INDEX1, "gl.STENCIL_INDEX1", 0x8D46)
	cxcore.ConstI32(CONST_GL_STENCIL_INDEX4, "gl.STENCIL_INDEX4", 0x8D47)
	cxcore.ConstI32(CONST_GL_STENCIL_INDEX8, "gl.STENCIL_INDEX8", 0x8D48)
	cxcore.ConstI32(CONST_GL_STENCIL_INDEX16, "gl.STENCIL_INDEX16", 0x8D49)
	cxcore.ConstI32(CONST_GL_HALF_FLOAT, "gl.HALF_FLOAT", 0x140B)
	cxcore.ConstI32(CONST_GL_R8, "gl.R8", 0x8229)

	// gl_3_2
	cxcore.ConstI32(CONST_GL_FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS, "gl.FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS", 0x8DA8)

	// gl_4_4
	cxcore.ConstI32(CONST_GL_MIRROR_CLAMP_TO_EDGE, "gl.MIRROR_CLAMP_TO_EDGE", 0x8743)

	// Fixed pipeline. Deprecated ?
	cxcore.ConstI32(CONST_GL_POLYGON, "gl.POLYGON", 9)
	cxcore.ConstI32(CONST_GL_MODELVIEW, "gl.MODELVIEW", 5888)
	cxcore.ConstI32(CONST_GL_PROJECTION, "gl.PROJECTION", 5889)
	cxcore.ConstI32(CONST_GL_MODELVIEW_MATRIX, "gl.MODELVIEW_MATRIX", 2982)
	cxcore.ConstI32(CONST_GL_LIGHTING, "gl.LIGHTING", 2896)
	cxcore.ConstI32(CONST_GL_LIGHT0, "gl.LIGHT0", 16384)
	cxcore.ConstI32(CONST_GL_AMBIENT, "gl.AMBIENT", 4608)
	cxcore.ConstI32(CONST_GL_DIFFUSE, "gl.DIFFUSE", 4609)
	cxcore.ConstI32(CONST_GL_POSITION, "gl.POSITION", 4611)
	cxcore.ConstI32(CONST_GL_TEXTURE_ENV, "gl.TEXTURE_ENV", 8960)
	cxcore.ConstI32(CONST_GL_TEXTURE_ENV_MODE, "gl.TEXTURE_ENV_MODE", 8704)
	cxcore.ConstI32(CONST_GL_MODULATE, "gl.MODULATE", 8448)
	cxcore.ConstI32(CONST_GL_DECAL, "gl.DECAL", 8449)
	cxcore.ConstI32(CONST_GL_POINT_SMOOTH, "gl.POINT_SMOOTH", 2832)

	// glfw
	cxcore.ConstI32(CONST_GLFW_FALSE, "glfw.False", 0)
	cxcore.ConstI32(CONST_GLFW_TRUE, "glfw.True", 1)
	cxcore.ConstI32(CONST_GLFW_PRESS, "glfw.Press", 1)
	cxcore.ConstI32(CONST_GLFW_RELEASE, "glfw.Release", 0)
	cxcore.ConstI32(CONST_GLFW_REPEAT, "glfw.Repeat", 2)
	cxcore.ConstI32(CONST_GLFW_KEY_UNKNOWN, "glfw.KeyUnknown", -1)
	cxcore.ConstI32(CONST_GLFW_CURSOR, "glfw.Cursor", 208897)
	cxcore.ConstI32(CONST_GLFW_STICKY_KEYS, "glfw.StickyKeys", 208898)
	cxcore.ConstI32(CONST_GLFW_STICKY_MOUSE_BUTTONS, "glfw.StickyMouseButtons", 208899)
	cxcore.ConstI32(CONST_GLFW_CURSOR_NORMAL, "glfw.CursorNormal", 212993)
	cxcore.ConstI32(CONST_GLFW_CURSOR_HIDDEN, "glfw.CursorHidden", 212994)
	cxcore.ConstI32(CONST_GLFW_CURSOR_DISABLED, "glfw.CursorDisabled", 212995)
	cxcore.ConstI32(CONST_GLFW_RESIZABLE, "glfw.Resizable", 131075)
	cxcore.ConstI32(CONST_GLFW_CONTEXT_VERSION_MAJOR, "glfw.ContextVersionMajor", 139266)
	cxcore.ConstI32(CONST_GLFW_CONTEXT_VERSION_MINOR, "glfw.ContextVersionMinor", 139267)
	cxcore.ConstI32(CONST_GLFW_OPENGL_PROFILE, "glfw.OpenGLProfile", 139272)
	cxcore.ConstI32(CONST_GLFW_OPENGL_COREPROFILE, "glfw.OpenGLCoreProfile", 204801)
	cxcore.ConstI32(CONST_GLFW_OPENGL_FORWARD_COMPATIBLE, "glfw.OpenGLForwardCompat", 139270)
	cxcore.ConstI32(CONST_GLFW_MOUSE_BUTTON_LAST, "glfw.MouseButtonLast", 7)
	cxcore.ConstI32(CONST_GLFW_MOUSE_BUTTON_LEFT, "glfw.MouseButtonLeft", 0)
	cxcore.ConstI32(CONST_GLFW_MOUSE_BUTTON_RIGHT, "glfw.MouseButtonRight", 1)
	cxcore.ConstI32(CONST_GLFW_MOUSE_BUTTON_MIDDLE, "glfw.MouseButtonMiddle", 2)
	cxcore.ConstI32(CONST_GLFW_COCOA_RETINA_FRAMEBUFFER, "glfw.CocoaRetinaFramebuffer", 0x23001)
	cxcore.ConstI32(CONST_GLFW_SCALE_TO_MONITOR, "glfw.ScaleToMonitor", 0x2200C)

	// gltext
	cxcore.ConstI32(CONST_GLTEXT_LEFT_TO_RIGHT, "gltext.LeftToRight", 0)
	cxcore.ConstI32(CONST_GLTEXT_RIGHT_TO_LEFT, "gltext.RightToLeft", 1)
	cxcore.ConstI32(CONST_GLTEXT_TOP_TO_BOTTOM, "gltext.TopToBottom", 2)
	// openal
	cxcore.ConstI32(CONST_AL_INVALID_NAME, "al.InvalidName", 0xA001)
	cxcore.ConstI32(CONST_AL_INVALID_ENUM, "al.InvalidEnum", 0xA002)
	cxcore.ConstI32(CONST_AL_INVALID_VALUE, "al.InvalidValue", 0xA003)
	cxcore.ConstI32(CONST_AL_INVALID_OPERATION, "al.InvalidOperation", 0xA004)
	cxcore.ConstI32(CONST_AL_OUT_OF_MEMORY, "al.OutOfMemory", 0xA005)
	cxcore.ConstI32(CONST_AL_INVERSE_DISTANCE, "al.InverseDistance", 0xD001)
	cxcore.ConstI32(CONST_AL_INVERSE_DISTANCE_CLAMPED, "al.InverseDistanceClamped", 0xD002)
	cxcore.ConstI32(CONST_AL_LINEAR_DISTANCE, "al.LinearDistance", 0xD003)
	cxcore.ConstI32(CONST_AL_LINEAR_DISTANCE_CLAMPED, "al.LinearDistanceClamped", 0xD004)
	cxcore.ConstI32(CONST_AL_EXPONENT_DISTANCE, "al.ExponentDistance", 0xD005)
	cxcore.ConstI32(CONST_AL_EXPONENT_DISTANCE_CLAMPED, "al.ExponentDistanceClamped", 0xD006)
	cxcore.ConstI32(CONST_AL_INITIAL, "al.Initial", 0x1011)
	cxcore.ConstI32(CONST_AL_PLAYING, "al.Playing", 0x1012)
	cxcore.ConstI32(CONST_AL_PAUSED, "al.Paused", 0x1013)
	cxcore.ConstI32(CONST_AL_STOPPED, "al.Stopped", 0x1014)
	cxcore.ConstI32(CONST_AL_FORMAT_MONO_8, "al.FormatMono8", 0x1100)
	cxcore.ConstI32(CONST_AL_FORMAT_MONO_16, "al.FormatMono16", 0x1101)
	cxcore.ConstI32(CONST_AL_FORMAT_STEREO_8, "al.FormatStereo8", 0x1102)
	cxcore.ConstI32(CONST_AL_FORMAT_STEREO_16, "al.FormatStereo16", 0x1103)
	cxcore.ConstI32(CONST_AL_WAV_FORMAT_PCM, "al.WavFormatPCM", 1)
	cxcore.ConstI32(CONST_AL_WAV_FORMAT_IEEE_FLOAT, "al.WavFormatIEEEFloat", 3)
}
