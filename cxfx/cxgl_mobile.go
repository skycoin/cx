// +build cxfx,mobile

package cxfx

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/skycoin/gltext"
	"golang.org/x/mobile/gl"

	"github.com/skycoin/cx/cx"
)

func Ptr(data interface{}) (unsafe.Pointer, uintptr) {
	if data == nil {
		return unsafe.Pointer(nil), 0
	}
	var addr unsafe.Pointer
	var size uintptr
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Ptr:
		e := v.Elem()
		size = e.Type().Size()
		switch e.Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			addr = unsafe.Pointer(e.UnsafeAddr())
		default:
			panic(fmt.Errorf("unsupported pointer to type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", e.Kind()))
		}
	case reflect.Uintptr:
		addr = unsafe.Pointer(v.Pointer())
		size = v.Type().Size()
	case reflect.Slice:
		addr = unsafe.Pointer(v.Index(0).UnsafeAddr())
		size = uintptr(v.Len()) * v.Type().Elem().Size()
	default:
		panic(fmt.Errorf("unsupported type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", v.Type()))
	}
	return addr, size
}

func ByteSlice(i interface{}) []byte {
	ptr, size := Ptr(i)
	h := reflect.SliceHeader{uintptr(ptr), int(size), int(size)}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Float32Slice(i interface{}) []float32 {
	ptr, size := Ptr(i)
	floatSize := int(size / unsafe.Sizeof(float32(0)))
	h := reflect.SliceHeader{uintptr(ptr), int(floatSize), int(floatSize)}
	return *(*[]float32)(unsafe.Pointer(&h))
}

func Int32Slice(i interface{}) []int32 {
	ptr, size := Ptr(i)
	intSize := int(size / unsafe.Sizeof(int32(0)))
	h := reflect.SliceHeader{uintptr(ptr), int(intSize), int(intSize)}
	return *(*[]int32)(unsafe.Pointer(&h))
}

func EnumSlice(i interface{}) []gl.Enum {
	ptr, size := Ptr(i)
	enumSize := int(size / unsafe.Sizeof(gl.Enum(0)))
	h := reflect.SliceHeader{uintptr(ptr), enumSize, enumSize}
	return *(*[]gl.Enum)(unsafe.Pointer(&h))
}

const (
	cxglCLAMP_TO_EDGE               = gl.CLAMP_TO_EDGE
	cxglNEAREST                     = gl.NEAREST
	cxglRGB                         = gl.RGB
	cxglRGBA                        = gl.RGBA
	cxglRGBA8                       = gl.RGBA8
	cxglRGB16F                      = gl.RGB16F
	cxglTEXTURE_2D                  = gl.TEXTURE_2D
	cxglTEXTURE_CUBE_MAP            = gl.TEXTURE_CUBE_MAP
	cxglTEXTURE_CUBE_MAP_POSITIVE_X = gl.TEXTURE_CUBE_MAP_POSITIVE_X
	cxglTEXTURE_MAG_FILTER          = gl.TEXTURE_MAG_FILTER
	cxglTEXTURE_MIN_FILTER          = gl.TEXTURE_MIN_FILTER
	cxglTEXTURE_WRAP_R              = gl.TEXTURE_WRAP_R
	cxglTEXTURE_WRAP_S              = gl.TEXTURE_WRAP_S
	cxglTEXTURE_WRAP_T              = gl.TEXTURE_WRAP_T
	cxglUNSIGNED_BYTE               = gl.UNSIGNED_BYTE
	cxglFLOAT                       = gl.FLOAT
)

// gogl
var glctx gl.Context

func SetGLContext(ctx gl.Context) {
	glctx = ctx
	gltext.SetGLContext(ctx)
}

// gogl
func opGlInit(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlDestroy(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlStrs(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlFree(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

// cxgl_1_0
func cxglCullFace(mode uint32) {
	glctx.CullFace(gl.Enum(mode))
}

func cxglFrontFace(mode uint32) {
	glctx.FrontFace(gl.Enum(mode))
}

func cxglHint(target uint32, mode uint32) {
	glctx.Hint(gl.Enum(target), gl.Enum(mode))
}

func cxglScissor(x int32, y int32, width int32, height int32) {
	glctx.Scissor(x, y, width, height)
}

func cxglTexParameteri(target uint32, pname uint32, param int32) {
	glctx.TexParameteri(gl.Enum(target), gl.Enum(pname), int(param))
}

func cxglTexImage2D(target uint32, level int32, internalFormat int32, width int32, height int32, border int32, format uint32, gltype uint32, pixels interface{}) {
	glctx.TexImage2D(gl.Enum(target), int(level), int(internalFormat), int(width), int(height), int(border), gl.Enum(format), gl.Enum(gltype), ByteSlice((pixels)))
}

func cxglClear(mask uint32) {
	glctx.Clear(gl.Enum(mask))
}

func cxglClearColor(red float32, green float32, blue float32, alpha float32) {
	glctx.ClearColor(red, green, blue, alpha)
}

func cxglClearStencil(s int32) {
	glctx.ClearStencil(int(s))
}

func cxglClearDepth(depth float64) {
	glctx.ClearDepthf(float32(depth))
}

func cxglStencilMask(mask uint32) {
	glctx.StencilMask(mask)
}

func cxglColorMask(red bool, green bool, blue bool, alpha bool) {
	glctx.ColorMask(red, green, blue, alpha)
}

func cxglDepthMask(flag bool) {
	glctx.DepthMask(flag)
}

func cxglDisable(cap uint32) {
	glctx.Disable(gl.Enum(cap))
}

func cxglEnable(cap uint32) {
	glctx.Enable(gl.Enum(cap))
}

func cxglBlendFunc(sfactor uint32, dfactor uint32) {
	glctx.BlendFunc(gl.Enum(sfactor), gl.Enum(dfactor))
}

func cxglStencilFunc(glfunc uint32, ref int32, mask uint32) {
	glctx.StencilFunc(gl.Enum(glfunc), int(ref), mask)
}

func cxglStencilOp(fail uint32, zfail uint32, zpass uint32) {
	glctx.StencilOp(gl.Enum(fail), gl.Enum(zfail), gl.Enum(zpass))
}

func cxglDepthFunc(glfunc uint32) {
	glctx.DepthFunc(gl.Enum(glfunc))
}

func cxglGetError() uint32 {
	return uint32(glctx.GetError())
}

func cxglGetTexLevelParameteriv(target uint32, level int32, pname uint32, params *int32) {
	var sparams []int32 = make([]int32, 1)
	glctx.(gl.Context31).GetTexLevelParameteriv(sparams, gl.Enum(target), int(level), gl.Enum(pname))
	*params = sparams[0] // TODO : remove hardcode
}

func cxglDepthRange(n float64, f float64) {
	glctx.DepthRangef(float32(n), float32(f))
}

func cxglViewport(x int32, y int32, width int32, height int32) {
	glctx.Viewport(int(x), int(y), int(width), int(height))
}

// cxgl_1_1
func cxglDrawArrays(mode uint32, first int32, count int32) {
	glctx.DrawArrays(gl.Enum(mode), int(first), int(count))
}

func cxglDrawElements(mode uint32, count int32, gltype uint32, indices interface{}) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED) // last param can be either offset or pointer
	glctx.DrawElements(gl.Enum(mode), int(count), gl.Enum(gltype), 0)
}

func cxglBindTexture(target uint32, texture uint32) {
	glctx.BindTexture(gl.Enum(target), gl.Texture{texture})
}

func cxglDeleteTextures(n int32, textures *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.DeleteTexture(gl.Texture{*textures})
}

func cxglGenTextures(n int32, textures *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	*textures = glctx.CreateTexture().Value
}

// cxgl_1_3
func cxglActiveTexture(texture uint32) {
	glctx.ActiveTexture(gl.Enum(texture))
}

// cxgl_1_4
func cxglBlendFuncSeparate(sfactorRGB uint32, dfactorRGB uint32, sfactorAlpha uint32, dfactorAlpha uint32) {
	glctx.BlendFuncSeparate(gl.Enum(sfactorRGB), gl.Enum(dfactorRGB), gl.Enum(sfactorAlpha), gl.Enum(dfactorAlpha))
}

// cxgl_1_5
func cxglBindBuffer(target uint32, buffer uint32) {
	glctx.BindBuffer(gl.Enum(target), gl.Buffer{buffer})
}

func cxglDeleteBuffers(n int32, buffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.DeleteBuffer(gl.Buffer{*buffers})
}

func cxglGenBuffers(n int32, buffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	*buffers = glctx.CreateBuffer().Value
}

// cxgl_2_0
func cxglBufferData(target uint32, size int, data interface{}, usage uint32) {
	slice := ByteSlice(data)
	if slice == nil || len(slice) == 0 {
		glctx.BufferInit(gl.Enum(target), size, gl.Enum(usage))
	} else {
		glctx.BufferData(gl.Enum(target), ByteSlice(data), gl.Enum(usage))
	}
}

func cxglBufferSubData(target uint32, offset int, size int, data interface{}) {
	glctx.BufferSubData(gl.Enum(target), offset, ByteSlice(data))
}

func cxglDrawBuffers(size int32, bufs interface{}) {
	glctx.(gl.Context3).DrawBuffers(int(size), EnumSlice(bufs))
}

func cxglStencilOpSeparate(face uint32, sfail uint32, dpfail uint32, dppass uint32) {
	glctx.StencilOpSeparate(gl.Enum(face), gl.Enum(sfail), gl.Enum(dpfail), gl.Enum(dppass))
}

func cxglStencilFuncSeparate(face uint32, glfunc uint32, ref int32, mask uint32) {
	glctx.StencilFuncSeparate(gl.Enum(face), gl.Enum(glfunc), int(ref), mask)
}

func cxglStencilMaskSeparate(face uint32, mask uint32) {
	glctx.StencilMaskSeparate(gl.Enum(face), mask)
}

func cxglAttachShader(program uint32, shader uint32) {
	glctx.AttachShader(gl.Program{true, program}, gl.Shader{shader})
}

func cxglBindAttribLocation(program uint32, index uint32, name string) {
	glctx.BindAttribLocation(gl.Program{true, program}, gl.Attrib{uint(index)}, name)
}

func cxglCompileShader(shader uint32) {
	glctx.CompileShader(gl.Shader{shader})
}

func cxglCreateProgram() uint32 {
	program := glctx.CreateProgram()
	return program.Value
}

func cxglCreateShader(gltype uint32) uint32 {
	return glctx.CreateShader(gl.Enum(gltype)).Value
}

func cxglDeleteProgram(program uint32) {
	glctx.DeleteProgram(gl.Program{true, program})
}

func cxglDeleteShader(shader uint32) {
	glctx.DeleteShader(gl.Shader{shader})
}

func cxglDetachShader(program uint32, shader uint32) {
	glctx.DetachShader(gl.Program{true, program}, gl.Shader{shader})
}

func cxglEnableVertexAttribArray(index uint32) {
	glctx.EnableVertexAttribArray(gl.Attrib{uint(index)})
}

func cxglGetAttribLocation(program uint32, name string) int32 {
	return int32(glctx.GetAttribLocation(gl.Program{true, program}, name).Value)
}

func cxglGetProgramiv(program uint32, pname uint32) int32 {
	return int32(glctx.GetProgrami(gl.Program{true, program}, gl.Enum(pname)))
}

func cxglGetProgramInfoLog(program uint32, maxLen int32) string {
	ll := glctx.GetProgramInfoLog(gl.Program{true, program})
	return ll
}

func cxglGetShaderiv(shader uint32, pname uint32) int32 {
	return int32(glctx.GetShaderi(gl.Shader{shader}, gl.Enum(pname)))
}

func cxglGetShaderInfoLog(shader uint32, maxLen int32) string {
	return glctx.GetShaderInfoLog(gl.Shader{shader})
}

func cxglGetUniformLocation(program uint32, name string) int32 {
	return glctx.GetUniformLocation(gl.Program{true, program}, name).Value
}

func cxglLinkProgram(program uint32) {
	glctx.LinkProgram(gl.Program{true, program})
}

func cxglShaderSource(shader uint32, count int32, glstring string) {
	if count > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.ShaderSource(gl.Shader{shader}, glstring)
}

func cxglUseProgram(program uint32) {
	glctx.UseProgram(gl.Program{true, program})
}

func cxglUniform1f(location int32, v0 float32) {
	glctx.Uniform1f(gl.Uniform{location}, v0)
}

func cxglUniform2f(location int32, v0 float32, v1 float32) {
	glctx.Uniform2f(gl.Uniform{location}, v0, v1)
}

func cxglUniform3f(location int32, v0 float32, v1 float32, v2 float32) {
	glctx.Uniform3f(gl.Uniform{location}, v0, v1, v2)
}

func cxglUniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	glctx.Uniform4f(gl.Uniform{location}, v0, v1, v2, v3)
}

func cxglUniform1i(location int32, v0 int32) {
	glctx.Uniform1i(gl.Uniform{location}, int(v0))
}

func cxglUniform2i(location int32, v0 int32, v1 int32) {
	glctx.Uniform2i(gl.Uniform{location}, int(v0), int(v1))
}

func cxglUniform3i(location int32, v0 int32, v1 int32, v2 int32) {
	glctx.Uniform3i(gl.Uniform{location}, v0, v1, v2)
}

func cxglUniform4i(location int32, v0 int32, v1 int32, v2 int32, v3 int32) {
	glctx.Uniform4i(gl.Uniform{location}, v0, v1, v2, v3)
}

func cxglUniform1fv(location int32, count int32, value interface{}) {
	glctx.Uniform1fv(gl.Uniform{location}, Float32Slice(value))
}

func cxglUniform2fv(location int32, count int32, value interface{}) {
	glctx.Uniform2fv(gl.Uniform{location}, Float32Slice(value))
}

func cxglUniform3fv(location int32, count int32, value interface{}) {
	glctx.Uniform3fv(gl.Uniform{location}, Float32Slice(value))
}

func cxglUniform4fv(location int32, count int32, value interface{}) {
	glctx.Uniform4fv(gl.Uniform{location}, Float32Slice(value))
}

func cxglUniform1iv(location int32, count int32, value interface{}) {
	glctx.Uniform1iv(gl.Uniform{location}, Int32Slice(value))
}

func cxglUniform2iv(location int32, count int32, value interface{}) {
	glctx.Uniform2iv(gl.Uniform{location}, Int32Slice(value))
}

func cxglUniform3iv(location int32, count int32, value interface{}) {
	glctx.Uniform3iv(gl.Uniform{location}, Int32Slice(value))
}

func cxglUniform4iv(location int32, count int32, value interface{}) {
	glctx.Uniform4iv(gl.Uniform{location}, Int32Slice(value))
}

func cxglUniformMatrix2fv(location int32, count int32, transpose bool, value interface{}) {
	if transpose {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.UniformMatrix2fv(gl.Uniform{location}, Float32Slice(value))
}

func cxglUniformMatrix3fv(location int32, count int32, transpose bool, value interface{}) {
	if transpose {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.UniformMatrix3fv(gl.Uniform{location}, Float32Slice(value))
}

func cxglUniformMatrix4fv(location int32, count int32, transpose bool, value interface{}) {
	if transpose {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.UniformMatrix4fv(gl.Uniform{location}, Float32Slice(value))
}

func cxglVertexAttribPointer(index uint32, size int32, gltype uint32, normalized bool, stride int32, pointer int32) {
	glctx.VertexAttribPointer(gl.Attrib{uint(index)}, int(size), gl.Enum(gltype), normalized, int(stride), int(pointer))
}

// cxgl_3_0
func cxglClearBufferiv(buffer uint32, drawBuffer int32, value []int32) {
	glctx.(gl.Context3).ClearBufferiv(gl.Enum(buffer), int(drawBuffer), value)
}

func cxglClearBufferuiv(buffer uint32, drawBuffer int32, value []uint32) {
	glctx.(gl.Context3).ClearBufferuiv(gl.Enum(buffer), int(drawBuffer), value)
}

func cxglClearBufferfv(buffer uint32, drawBuffer int32, value []float32) {
	glctx.(gl.Context3).ClearBufferfv(gl.Enum(buffer), int(drawBuffer), value)
}

func cxglBindRenderbuffer(buffer uint32, renderbuffer uint32) {
	glctx.BindRenderbuffer(gl.Enum(buffer), gl.Renderbuffer{renderbuffer})
}

func cxglDeleteRenderbuffers(n int32, renderbuffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.DeleteRenderbuffer(gl.Renderbuffer{*renderbuffers})
}

func cxglGenRenderbuffers(n int32, renderbuffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	*renderbuffers = glctx.CreateRenderbuffer().Value
}

func cxglRenderbufferStorage(target uint32, internalFormat uint32, width int32, height int32) {
	glctx.RenderbufferStorage(gl.Enum(target), gl.Enum(internalFormat), int(width), int(height))
}

func cxglBindFramebuffer(target uint32, framebuffer uint32) {
	glctx.BindFramebuffer(gl.Enum(target), gl.Framebuffer{framebuffer})
}

func cxglDeleteFramebuffers(n int32, framebuffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.DeleteFramebuffer(gl.Framebuffer{*framebuffers})
}

func cxglGenFramebuffers(n int32, framebuffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	*framebuffers = glctx.CreateFramebuffer().Value
}

func cxglCheckFramebufferStatus(target uint32) uint32 {
	return uint32(glctx.CheckFramebufferStatus(gl.Enum(target)))
}

func cxglFramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	glctx.FramebufferTexture2D(gl.Enum(target), gl.Enum(attachment), gl.Enum(textarget), gl.Texture{texture}, int(level))
}

func cxglFramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer uint32) {
	glctx.FramebufferRenderbuffer(gl.Enum(target), gl.Enum(attachment), gl.Enum(renderbuffertarget), gl.Renderbuffer{renderbuffer})
}

func cxglGenerateMipmap(target uint32) {
	glctx.GenerateMipmap(gl.Enum(target))
}

func cxglBindVertexArrayAPPLE(array uint32) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
	//glctx.BindVertexArrayAPPLE(array)
}

func cxglBindVertexArray(array uint32) {
	glctx.BindVertexArray(gl.VertexArray{array})
}

func cxglDeleteVertexArraysAPPLE(n int32, arrays *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
	//glctx.DeleteVertexArraysAPPLE(n, arrays)
}

func cxglDeleteVertexArrays(n int32, arrays *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	glctx.DeleteVertexArray(gl.VertexArray{*arrays})
}

func cxglGenVertexArraysAPPLE(n int32, arrays *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
	//glctx.GenVertexArraysAPPLE(n, arrays)
}

func cxglGenVertexArrays(n int32, arrays *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	*arrays = glctx.CreateVertexArray().Value
}

// gl_0_0
func opGlMatrixMode(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlRotatef(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlTranslatef(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlLoadIdentity(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlPushMatrix(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlPopMatrix(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlEnableClientState(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlColor3f(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlColor4f(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlBegin(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlEnd(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlNormal3f(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlVertex2f(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlVertex3f(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlLightfv(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlFrustum(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlTexEnvi(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlOrtho(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlScalef(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlTexCoord2d(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlTexCoord2f(inputs []CXValue, outputs []CXValue) {
	panic(CX_RUNTIME_NOT_IMPLEMENTED)
}
