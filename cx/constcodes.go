package base

// constant codes
const (
	// opengl
	CONST_GL_FALSE = iota
	CONST_GL_TRUE
	CONST_GL_INVALID_ENUM
	CONST_GL_INVALID_VALUE
	CONST_GL_INVALID_OPERATION
	CONST_GL_STACK_OVERFLOW
	CONST_GL_STACK_UNDERFLOW
	CONST_GL_OUT_OF_MEMORY
	CONST_GL_QUADS
	CONST_GL_COLOR_BUFFER_BIT
	CONST_GL_DEPTH_BUFFER_BIT
	CONST_GL_ARRAY_BUFFER
	CONST_GL_FLOAT
	CONST_GL_TRIANGLES
	CONST_GL_POLYGON
	CONST_GL_VERTEX_SHADER
	CONST_GL_FRAGMENT_SHADER
	CONST_GL_MODELVIEW
	CONST_GL_TEXTURE_2D
	CONST_GL_PROJECTION
	CONST_GL_TEXTURE
	CONST_GL_COLOR
	CONST_GL_LESS
	CONST_GL_MODELVIEW_MATRIX
	CONST_GL_VERTEX_ARRAY
	CONST_GL_STREAM_DRAW
	CONST_GL_STREAM_READ
	CONST_GL_STREAM_COPY
	CONST_GL_STATIC_DRAW
	CONST_GL_STATIC_READ
	CONST_GL_STATIC_COPY
	CONST_GL_DYNAMIC_DRAW
	CONST_GL_DYNAMIC_READ
	CONST_GL_DYNAMIC_COPY
	CONST_GL_BLEND
	CONST_GL_DEPTH_TEST
	CONST_GL_LIGHTING
	CONST_GL_LEQUAL
	CONST_GL_LIGHT0
	CONST_GL_AMBIENT
	CONST_GL_DIFFUSE
	CONST_GL_POSITION
	CONST_GL_TEXTURE_ENV
	CONST_GL_TEXTURE_ENV_MODE
	CONST_GL_MODULATE
	CONST_GL_DECAL
	CONST_GL_REPLACE
	CONST_GL_SRC_ALPHA
	CONST_GL_ONE_MINUS_SRC_ALPHA
	CONST_GL_DITHER
	CONST_GL_POINT_SMOOTH
	CONST_GL_LINE_SMOOTH
	CONST_GL_POLYGON_SMOOTH
	CONST_GL_DONT_CARE
	CONST_GL_POLYGON_SMOOTH_HINT
	CONST_GL_MULTISAMPLE_ARB

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
)

// For the parser. These shouldn't be used in the runtime for performance reasons
var ConstNames map[int]string = map[int]string{
	// opengl
	CONST_GL_FALSE:               "gl.FALSE",
	CONST_GL_TRUE:                "gl.TRUE",
	CONST_GL_INVALID_ENUM:        "gl.INVALID_ENUM",
	CONST_GL_INVALID_VALUE:       "gl.INVALID_VALUE",
	CONST_GL_INVALID_OPERATION:   "gl.INVALID_OPERATION",
	CONST_GL_STACK_OVERFLOW:      "gl.STACK_OVERFLOW",
	CONST_GL_STACK_UNDERFLOW:     "gl.STACK_UNDERFLOW",
	CONST_GL_OUT_OF_MEMORY:       "gl.OUT_OF_MEMORY",
	CONST_GL_QUADS:               "gl.QUADS",
	CONST_GL_COLOR_BUFFER_BIT:    "gl.COLOR_BUFFER_BIT",
	CONST_GL_DEPTH_BUFFER_BIT:    "gl.DEPTH_BUFFER_BIT",
	CONST_GL_ARRAY_BUFFER:        "gl.ARRAY_BUFFER",
	CONST_GL_FLOAT:               "gl.FLOAT",
	CONST_GL_TRIANGLES:           "gl.TRIANGLES",
	CONST_GL_POLYGON:             "gl.POLYGON",
	CONST_GL_VERTEX_SHADER:       "gl.VERTEX_SHADER",
	CONST_GL_FRAGMENT_SHADER:     "gl.FRAGMENT_SHADER",
	CONST_GL_MODELVIEW:           "gl.MODELVIEW",
	CONST_GL_TEXTURE_2D:          "gl.TEXTURE_2D",
	CONST_GL_PROJECTION:          "gl.PROJECTION",
	CONST_GL_TEXTURE:             "gl.TEXTURE",
	CONST_GL_COLOR:               "gl.COLOR",
	CONST_GL_LESS:                "gl.LESS",
	CONST_GL_MODELVIEW_MATRIX:    "gl.MODELVIEW_MATRIX",
	CONST_GL_VERTEX_ARRAY:        "gl.VERTEX_ARRAY",
	CONST_GL_STREAM_DRAW:         "gl.STREAM_DRAW",
	CONST_GL_STREAM_READ:         "gl.STREAM_READ",
	CONST_GL_STREAM_COPY:         "gl.STREAM_COPY",
	CONST_GL_STATIC_DRAW:         "gl.STATIC_DRAW",
	CONST_GL_STATIC_READ:         "gl.STATIC_READ",
	CONST_GL_STATIC_COPY:         "gl.STATIC_COPY",
	CONST_GL_DYNAMIC_DRAW:        "gl.DYNAMIC_DRAW",
	CONST_GL_DYNAMIC_READ:        "gl.DYNAMIC_READ",
	CONST_GL_DYNAMIC_COPY:        "gl.DYNAMIC_COPY",
	CONST_GL_BLEND:               "gl.BLEND",
	CONST_GL_DEPTH_TEST:          "gl.DEPTH_TEST",
	CONST_GL_LIGHTING:            "gl.LIGHTING",
	CONST_GL_LEQUAL:              "gl.LEQUAL",
	CONST_GL_LIGHT0:              "gl.LIGHT0",
	CONST_GL_AMBIENT:             "gl.AMBIENT",
	CONST_GL_DIFFUSE:             "gl.DIFFUSE",
	CONST_GL_POSITION:            "gl.POSITION",
	CONST_GL_TEXTURE_ENV:         "gl.TEXTURE_ENV",
	CONST_GL_TEXTURE_ENV_MODE:    "gl.TEXTURE_ENV_MODE",
	CONST_GL_MODULATE:            "gl.MODULATE",
	CONST_GL_DECAL:               "gl.DECAL",
	CONST_GL_REPLACE:             "gl.REPLACE",
	CONST_GL_SRC_ALPHA:           "gl.SRC_ALPHA",
	CONST_GL_ONE_MINUS_SRC_ALPHA: "gl.ONE_MINUS_SRC_ALPHA",
	CONST_GL_DITHER:              "gl.DITHER",
	CONST_GL_POINT_SMOOTH:        "gl.POINT_SMOOTH",
	CONST_GL_LINE_SMOOTH:         "gl.LINE_SMOOTH",
	CONST_GL_POLYGON_SMOOTH:      "gl.POLYGON_SMOOTH",
	CONST_GL_DONT_CARE:           "gl.DONT_CARE",
	CONST_GL_POLYGON_SMOOTH_HINT: "gl.POLYGON_SMOOTH_HINT",
	CONST_GL_MULTISAMPLE_ARB:     "gl.MULTISAMPLE_ARB",

	// glfw
	CONST_GLFW_FALSE:                     "glfw.False",
	CONST_GLFW_TRUE:                      "glfw.True",
	CONST_GLFW_PRESS:                     "glfw.Press",
	CONST_GLFW_RELEASE:                   "glfw.Release",
	CONST_GLFW_REPEAT:                    "glfw.Repeat",
	CONST_GLFW_KEY_UNKNOWN:               "glfw.KeyUnknown",
	CONST_GLFW_CURSOR:                    "glfw.Cursor",
	CONST_GLFW_STICKY_KEYS:               "glfw.StickyKeys",
	CONST_GLFW_STICKY_MOUSE_BUTTONS:      "glfw.StickyMouseButtons",
	CONST_GLFW_CURSOR_NORMAL:             "glfw.CursorNormal",
	CONST_GLFW_CURSOR_HIDDEN:             "glfw.CursorHidden",
	CONST_GLFW_CURSOR_DISABLED:           "glfw.CursorDisabled",
	CONST_GLFW_RESIZABLE:                 "glfw.Resizable",
	CONST_GLFW_CONTEXT_VERSION_MAJOR:     "glfw.ContextVersionMajor",
	CONST_GLFW_CONTEXT_VERSION_MINOR:     "glfw.ContextVersionMinor",
	CONST_GLFW_OPENGL_PROFILE:            "glfw.Opengl.Profile",
	CONST_GLFW_OPENGL_COREPROFILE:        "glfw.Opengl.Coreprofile",
	CONST_GLFW_OPENGL_FORWARD_COMPATIBLE: "glfw.Opengl.ForwardCompatible",
	CONST_GLFW_MOUSE_BUTTON_LAST:         "glfw.MouseButtonLast",
	CONST_GLFW_MOUSE_BUTTON_LEFT:         "glfw.MouseButtonLeft",
	CONST_GLFW_MOUSE_BUTTON_RIGHT:        "glfw.MouseButtonRight",
	CONST_GLFW_MOUSE_BUTTON_MIDDLE:       "glfw.MouseButtonMiddle",
}

// For the parser. These shouldn't be used in the runtime for performance reasons
var ConstCodes map[string]int = map[string]int{
	// opengl
	"gl.FALSE":               CONST_GL_FALSE,
	"gl.TRUE":                CONST_GL_TRUE,
	"gl.INVALID_ENUM":        CONST_GL_INVALID_ENUM,
	"gl.INVALID_VALUE":       CONST_GL_INVALID_VALUE,
	"gl.INVALID_OPERATION":   CONST_GL_INVALID_OPERATION,
	"gl.STACK_OVERFLOW":      CONST_GL_STACK_OVERFLOW,
	"gl.STACK_UNDERFLOW":     CONST_GL_STACK_UNDERFLOW,
	"gl.OUT_OF_MEMORY":       CONST_GL_OUT_OF_MEMORY,
	"gl.QUADS":               CONST_GL_QUADS,
	"gl.COLOR_BUFFER_BIT":    CONST_GL_COLOR_BUFFER_BIT,
	"gl.DEPTH_BUFFER_BIT":    CONST_GL_DEPTH_BUFFER_BIT,
	"gl.ARRAY_BUFFER":        CONST_GL_ARRAY_BUFFER,
	"gl.FLOAT":               CONST_GL_FLOAT,
	"gl.TRIANGLES":           CONST_GL_TRIANGLES,
	"gl.POLYGON":             CONST_GL_POLYGON,
	"gl.VERTEX_SHADER":       CONST_GL_VERTEX_SHADER,
	"gl.FRAGMENT_SHADER":     CONST_GL_FRAGMENT_SHADER,
	"gl.MODELVIEW":           CONST_GL_MODELVIEW,
	"gl.TEXTURE_2D":          CONST_GL_TEXTURE_2D,
	"gl.PROJECTION":          CONST_GL_PROJECTION,
	"gl.TEXTURE":             CONST_GL_TEXTURE,
	"gl.COLOR":               CONST_GL_COLOR,
	"gl.LESS":                CONST_GL_LESS,
	"gl.MODELVIEW_MATRIX":    CONST_GL_MODELVIEW_MATRIX,
	"gl.VERTEX_ARRAY":        CONST_GL_VERTEX_ARRAY,
	"gl.STREAM_DRAW":         CONST_GL_STREAM_DRAW,
	"gl.STREAM_READ":         CONST_GL_STREAM_READ,
	"gl.STREAM_COPY":         CONST_GL_STREAM_COPY,
	"gl.STATIC_DRAW":         CONST_GL_STATIC_DRAW,
	"gl.STATIC_READ":         CONST_GL_STATIC_READ,
	"gl.STATIC_COPY":         CONST_GL_STATIC_COPY,
	"gl.DYNAMIC_DRAW":        CONST_GL_DYNAMIC_DRAW,
	"gl.DYNAMIC_READ":        CONST_GL_DYNAMIC_READ,
	"gl.DYNAMIC_COPY":        CONST_GL_DYNAMIC_COPY,
	"gl.BLEND":               CONST_GL_BLEND,
	"gl.DEPTH_TEST":          CONST_GL_DEPTH_TEST,
	"gl.LIGHTING":            CONST_GL_LIGHTING,
	"gl.LEQUAL":              CONST_GL_LEQUAL,
	"gl.LIGHT0":              CONST_GL_LIGHT0,
	"gl.AMBIENT":             CONST_GL_AMBIENT,
	"gl.DIFFUSE":             CONST_GL_DIFFUSE,
	"gl.POSITION":            CONST_GL_POSITION,
	"gl.TEXTURE_ENV":         CONST_GL_TEXTURE_ENV,
	"gl.TEXTURE_ENV_MODE":    CONST_GL_TEXTURE_ENV_MODE,
	"gl.MODULATE":            CONST_GL_MODULATE,
	"gl.DECAL":               CONST_GL_DECAL,
	"gl.REPLACE":             CONST_GL_REPLACE,
	"gl.SRC_ALPHA":           CONST_GL_SRC_ALPHA,
	"gl.ONE_MINUS_SRC_ALPHA": CONST_GL_ONE_MINUS_SRC_ALPHA,
	"gl.DITHER":              CONST_GL_DITHER,
	"gl.POINT_SMOOTH":        CONST_GL_POINT_SMOOTH,
	"gl.LINE_SMOOTH":         CONST_GL_LINE_SMOOTH,
	"gl.POLYGON_SMOOTH":      CONST_GL_POLYGON_SMOOTH,
	"gl.DONT_CARE":           CONST_GL_DONT_CARE,
	"gl.POLYGON_SMOOTH_HINT": CONST_GL_POLYGON_SMOOTH_HINT,
	"gl.MULTISAMPLE_ARB":     CONST_GL_MULTISAMPLE_ARB,

	// glfw
	"glfw.False":                    CONST_GLFW_FALSE,
	"glfw.True":                     CONST_GLFW_TRUE,
	"glfw.Press":                    CONST_GLFW_PRESS,
	"glfw.Release":                  CONST_GLFW_RELEASE,
	"glfw.Repeat":                   CONST_GLFW_REPEAT,
	"glfw.KeyUnknown":               CONST_GLFW_KEY_UNKNOWN,
	"glfw.Cursor":                   CONST_GLFW_CURSOR,
	"glfw.StickyKeys":               CONST_GLFW_STICKY_KEYS,
	"glfw.StickyMouseButtons":       CONST_GLFW_STICKY_MOUSE_BUTTONS,
	"glfw.CursorNormal":             CONST_GLFW_CURSOR_NORMAL,
	"glfw.CursorHidden":             CONST_GLFW_CURSOR_HIDDEN,
	"glfw.CursorDisabled":           CONST_GLFW_CURSOR_DISABLED,
	"glfw.Resizable":                CONST_GLFW_RESIZABLE,
	"glfw.ContextVersionMajor":      CONST_GLFW_CONTEXT_VERSION_MAJOR,
	"glfw.ContextVersionMinor":      CONST_GLFW_CONTEXT_VERSION_MINOR,
	"glfw.Opengl.Profile":           CONST_GLFW_OPENGL_PROFILE,
	"glfw.Opengl.Coreprofile":       CONST_GLFW_OPENGL_COREPROFILE,
	"glfw.Opengl.ForwardCompatible": CONST_GLFW_OPENGL_FORWARD_COMPATIBLE,
	"glfw.MouseButtonLast":          CONST_GLFW_MOUSE_BUTTON_LAST,
	"glfw.MouseButtonLeft":          CONST_GLFW_MOUSE_BUTTON_LEFT,
	"glfw.MouseButtonRight":         CONST_GLFW_MOUSE_BUTTON_RIGHT,
	"glfw.MouseButtonMiddle":        CONST_GLFW_MOUSE_BUTTON_MIDDLE,
}

var Constants map[int]CXConstant = map[int]CXConstant{
	// opengl
	CONST_GL_FALSE: CXConstant{Type: TYPE_I32, Value: FromI32(0)},

	CONST_GL_TRUE:                CXConstant{Type: TYPE_I32, Value: FromI32(1)},
	CONST_GL_INVALID_ENUM:        CXConstant{Type: TYPE_I32, Value: FromI32(1280)},
	CONST_GL_INVALID_VALUE:       CXConstant{Type: TYPE_I32, Value: FromI32(1281)},
	CONST_GL_INVALID_OPERATION:   CXConstant{Type: TYPE_I32, Value: FromI32(1282)},
	CONST_GL_STACK_OVERFLOW:      CXConstant{Type: TYPE_I32, Value: FromI32(1283)},
	CONST_GL_STACK_UNDERFLOW:     CXConstant{Type: TYPE_I32, Value: FromI32(1284)},
	CONST_GL_OUT_OF_MEMORY:       CXConstant{Type: TYPE_I32, Value: FromI32(1285)},
	CONST_GL_QUADS:               CXConstant{Type: TYPE_I32, Value: FromI32(7)},
	CONST_GL_COLOR_BUFFER_BIT:    CXConstant{Type: TYPE_I32, Value: FromI32(16384)},
	CONST_GL_DEPTH_BUFFER_BIT:    CXConstant{Type: TYPE_I32, Value: FromI32(256)},
	CONST_GL_ARRAY_BUFFER:        CXConstant{Type: TYPE_I32, Value: FromI32(34962)},
	CONST_GL_FLOAT:               CXConstant{Type: TYPE_I32, Value: FromI32(5126)},
	CONST_GL_TRIANGLES:           CXConstant{Type: TYPE_I32, Value: FromI32(4)},
	CONST_GL_POLYGON:             CXConstant{Type: TYPE_I32, Value: FromI32(9)},
	CONST_GL_VERTEX_SHADER:       CXConstant{Type: TYPE_I32, Value: FromI32(35633)},
	CONST_GL_FRAGMENT_SHADER:     CXConstant{Type: TYPE_I32, Value: FromI32(35632)},
	CONST_GL_MODELVIEW:           CXConstant{Type: TYPE_I32, Value: FromI32(5888)},
	CONST_GL_TEXTURE_2D:          CXConstant{Type: TYPE_I32, Value: FromI32(3553)},
	CONST_GL_PROJECTION:          CXConstant{Type: TYPE_I32, Value: FromI32(5889)},
	CONST_GL_TEXTURE:             CXConstant{Type: TYPE_I32, Value: FromI32(5890)},
	CONST_GL_COLOR:               CXConstant{Type: TYPE_I32, Value: FromI32(6144)},
	CONST_GL_LESS:                CXConstant{Type: TYPE_I32, Value: FromI32(513)},
	CONST_GL_MODELVIEW_MATRIX:    CXConstant{Type: TYPE_I32, Value: FromI32(2982)},
	CONST_GL_VERTEX_ARRAY:        CXConstant{Type: TYPE_I32, Value: FromI32(32884)},
	CONST_GL_STREAM_DRAW:         CXConstant{Type: TYPE_I32, Value: FromI32(35040)},
	CONST_GL_STREAM_READ:         CXConstant{Type: TYPE_I32, Value: FromI32(35041)},
	CONST_GL_STREAM_COPY:         CXConstant{Type: TYPE_I32, Value: FromI32(35042)},
	CONST_GL_STATIC_DRAW:         CXConstant{Type: TYPE_I32, Value: FromI32(35044)},
	CONST_GL_STATIC_READ:         CXConstant{Type: TYPE_I32, Value: FromI32(35045)},
	CONST_GL_STATIC_COPY:         CXConstant{Type: TYPE_I32, Value: FromI32(35046)},
	CONST_GL_DYNAMIC_DRAW:        CXConstant{Type: TYPE_I32, Value: FromI32(35048)},
	CONST_GL_DYNAMIC_READ:        CXConstant{Type: TYPE_I32, Value: FromI32(35049)},
	CONST_GL_DYNAMIC_COPY:        CXConstant{Type: TYPE_I32, Value: FromI32(35050)},
	CONST_GL_BLEND:               CXConstant{Type: TYPE_I32, Value: FromI32(3042)},
	CONST_GL_DEPTH_TEST:          CXConstant{Type: TYPE_I32, Value: FromI32(2929)},
	CONST_GL_LIGHTING:            CXConstant{Type: TYPE_I32, Value: FromI32(2896)},
	CONST_GL_LEQUAL:              CXConstant{Type: TYPE_I32, Value: FromI32(515)},
	CONST_GL_LIGHT0:              CXConstant{Type: TYPE_I32, Value: FromI32(16384)},
	CONST_GL_AMBIENT:             CXConstant{Type: TYPE_I32, Value: FromI32(4608)},
	CONST_GL_DIFFUSE:             CXConstant{Type: TYPE_I32, Value: FromI32(4609)},
	CONST_GL_POSITION:            CXConstant{Type: TYPE_I32, Value: FromI32(4611)},
	CONST_GL_TEXTURE_ENV:         CXConstant{Type: TYPE_I32, Value: FromI32(8960)},
	CONST_GL_TEXTURE_ENV_MODE:    CXConstant{Type: TYPE_I32, Value: FromI32(8704)},
	CONST_GL_MODULATE:            CXConstant{Type: TYPE_I32, Value: FromI32(8448)},
	CONST_GL_DECAL:               CXConstant{Type: TYPE_I32, Value: FromI32(8449)},
	CONST_GL_REPLACE:             CXConstant{Type: TYPE_I32, Value: FromI32(7681)},
	CONST_GL_SRC_ALPHA:           CXConstant{Type: TYPE_I32, Value: FromI32(770)},
	CONST_GL_ONE_MINUS_SRC_ALPHA: CXConstant{Type: TYPE_I32, Value: FromI32(771)},
	CONST_GL_DITHER:              CXConstant{Type: TYPE_I32, Value: FromI32(3024)},
	CONST_GL_POINT_SMOOTH:        CXConstant{Type: TYPE_I32, Value: FromI32(2832)},
	CONST_GL_LINE_SMOOTH:         CXConstant{Type: TYPE_I32, Value: FromI32(2848)},
	CONST_GL_POLYGON_SMOOTH:      CXConstant{Type: TYPE_I32, Value: FromI32(2881)},
	CONST_GL_DONT_CARE:           CXConstant{Type: TYPE_I32, Value: FromI32(4352)},
	CONST_GL_POLYGON_SMOOTH_HINT: CXConstant{Type: TYPE_I32, Value: FromI32(3155)},
	CONST_GL_MULTISAMPLE_ARB:     CXConstant{Type: TYPE_I32, Value: FromI32(32925)},

	// glfw
	CONST_GLFW_FALSE:                     CXConstant{Type: TYPE_I32, Value: FromI32(0)},
	CONST_GLFW_TRUE:                      CXConstant{Type: TYPE_I32, Value: FromI32(1)},
	CONST_GLFW_PRESS:                     CXConstant{Type: TYPE_I32, Value: FromI32(1)},
	CONST_GLFW_RELEASE:                   CXConstant{Type: TYPE_I32, Value: FromI32(0)},
	CONST_GLFW_REPEAT:                    CXConstant{Type: TYPE_I32, Value: FromI32(2)},
	CONST_GLFW_KEY_UNKNOWN:               CXConstant{Type: TYPE_I32, Value: FromI32(-1)},
	CONST_GLFW_CURSOR:                    CXConstant{Type: TYPE_I32, Value: FromI32(208897)},
	CONST_GLFW_STICKY_KEYS:               CXConstant{Type: TYPE_I32, Value: FromI32(208898)},
	CONST_GLFW_STICKY_MOUSE_BUTTONS:      CXConstant{Type: TYPE_I32, Value: FromI32(208899)},
	CONST_GLFW_CURSOR_NORMAL:             CXConstant{Type: TYPE_I32, Value: FromI32(212993)},
	CONST_GLFW_CURSOR_HIDDEN:             CXConstant{Type: TYPE_I32, Value: FromI32(212994)},
	CONST_GLFW_CURSOR_DISABLED:           CXConstant{Type: TYPE_I32, Value: FromI32(212995)},
	CONST_GLFW_RESIZABLE:                 CXConstant{Type: TYPE_I32, Value: FromI32(131075)},
	CONST_GLFW_CONTEXT_VERSION_MAJOR:     CXConstant{Type: TYPE_I32, Value: FromI32(139266)},
	CONST_GLFW_CONTEXT_VERSION_MINOR:     CXConstant{Type: TYPE_I32, Value: FromI32(139267)},
	CONST_GLFW_OPENGL_PROFILE:            CXConstant{Type: TYPE_I32, Value: FromI32(139272)},
	CONST_GLFW_OPENGL_COREPROFILE:        CXConstant{Type: TYPE_I32, Value: FromI32(204801)},
	CONST_GLFW_OPENGL_FORWARD_COMPATIBLE: CXConstant{Type: TYPE_I32, Value: FromI32(139270)},
	CONST_GLFW_MOUSE_BUTTON_LAST:         CXConstant{Type: TYPE_I32, Value: FromI32(7)},
	CONST_GLFW_MOUSE_BUTTON_LEFT:         CXConstant{Type: TYPE_I32, Value: FromI32(0)},
	CONST_GLFW_MOUSE_BUTTON_RIGHT:        CXConstant{Type: TYPE_I32, Value: FromI32(1)},
	CONST_GLFW_MOUSE_BUTTON_MIDDLE:       CXConstant{Type: TYPE_I32, Value: FromI32(2)},
}
