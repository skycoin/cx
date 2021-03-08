// +build cxfx,!mobile

package cxfx

import (
	"strings"

	"github.com/go-gl/gl/v3.2-compatibility/gl"

	. "github.com/skycoin/cx/cx"
)

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
func getCString(key string, value string) **uint8 {
	if cstrings, ok := cSources[key]; ok {
		return cstrings
	} else {
		cstrings, free := gl.Strs(value + string('\000'))
		freeFns[key] = &free
		cSources[key] = cstrings
		return cstrings
	}
}

func freeCString(key string) {
	(*freeFns[key])()
	delete(freeFns, key)
	delete(cSources, key)
}

func to_pf32(p *uint8) *float32 {
	return (*float32)(gl.Ptr(p))
}

func to_pi32(p *uint8) *int32 {
	return (*int32)(gl.Ptr(p))
}

func to_pui32(p *uint8) *uint32 {
	return (*uint32)(gl.Ptr(p))
}

func ito_pui8(i interface{}) *uint8 {
	return (*uint8)(gl.Ptr(i))
}

func ito_pui32(i interface{}) *uint32 {
	return to_pui32(ito_pui8(i))
}

func ito_pi32(i interface{}) *int32 {
	return to_pi32(ito_pui8(i))
}

func ito_pf32(i interface{}) *float32 {
	return to_pf32(ito_pui8(i))
}

// TODO: change to ReadData_i64
func readPtr(fp int, inp *CXArgument, dataType int) *uint8 {
	// return (*uint8)(gl.Ptr(ReadData(fp, inp, dataType)))
	return (*uint8)(gl.Ptr(ReadData_i32(fp, inp, dataType)))
}

func readF32Ptr(fp int, inp *CXArgument) *float32 {
	return to_pf32(readPtr(fp, inp, TYPE_F32))
}

func readI32Ptr(fp int, inp *CXArgument) *int32 {
	return to_pi32(readPtr(fp, inp, TYPE_I32))
}

func readUI32Ptr(fp int, inp *CXArgument) *uint32 {
	return to_pui32(readPtr(fp, inp, TYPE_UI32))
}

// gogl
func opGlInit(expr *CXExpression, fp int) {
	gl.Init()
}

func opGlDestroy(expr *CXExpression, fp int) {
	for k, _ := range cSources {
		freeCString(k)
	}
}

func opGlStrs(expr *CXExpression, fp int) {
	getCString(ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]))
}

func opGlFree(expr *CXExpression, fp int) {
	freeCString(ReadStr(fp, expr.Inputs[0]))
}

// cxgl_1_0
func cxglCullFace(mode uint32) {
	gl.CullFace(mode)
}

func cxglFrontFace(mode uint32) {
	gl.FrontFace(mode)
}

func cxglHint(target uint32, mode uint32) {
	gl.Hint(target, mode)
}

func cxglScissor(x int32, y int32, width int32, height int32) {
	gl.Scissor(x, y, width, height)
}

func cxglTexParameteri(target uint32, pname uint32, param int32) {
	gl.TexParameteri(target, pname, param)
}

func cxglTexImage2D(target uint32, level int32, internalFormat int32, width int32, height int32, border int32, format uint32, gltype uint32, pixels interface{}) {
	gl.TexImage2D(target, level, internalFormat, width, height, border, format, gltype, gl.Ptr(pixels))
}

func cxglClear(mask uint32) {
	gl.Clear(mask)
}

func cxglClearColor(red float32, green float32, blue float32, alpha float32) {
	gl.ClearColor(red, green, blue, alpha)
}

func cxglClearStencil(s int32) {
	gl.ClearStencil(s)
}

func cxglClearDepth(depth float64) {
	gl.ClearDepth(depth)
}

func cxglStencilMask(mask uint32) {
	gl.StencilMask(mask)
}

func cxglColorMask(red bool, green bool, blue bool, alpha bool) {
	gl.ColorMask(red, green, blue, alpha)
}

func cxglDepthMask(flag bool) {
	gl.DepthMask(flag)
}

func cxglDisable(cap uint32) {
	gl.Disable(cap)
}

func cxglEnable(cap uint32) {
	gl.Enable(cap)
}

func cxglBlendFunc(sfactor uint32, dfactor uint32) {
	gl.BlendFunc(sfactor, dfactor)
}

func cxglStencilFunc(glfunc uint32, ref int32, mask uint32) {
	gl.StencilFunc(glfunc, ref, mask)
}

func cxglStencilOp(fail uint32, zfail uint32, zpass uint32) {
	gl.StencilOp(fail, zfail, zpass)
}

func cxglDepthFunc(glfunc uint32) {
	gl.DepthFunc(glfunc)
}

func cxglGetError() uint32 {
	return gl.GetError()
}

func cxglGetTexLevelParameteriv(target uint32, level int32, pname uint32, params *int32) {
	gl.GetTexLevelParameteriv(target, level, pname, params)
}

func cxglDepthRange(n float64, f float64) {
	gl.DepthRange(n, f)
}

func cxglViewport(x int32, y int32, width int32, height int32) {
	gl.Viewport(x, y, width, height)
}

// cxgl_1_1
func cxglDrawArrays(mode uint32, first int32, count int32) {
	gl.DrawArrays(mode, first, count)
}

func cxglDrawElements(mode uint32, count int32, gltype uint32, indices interface{}) {
	gl.DrawElements(mode, count, gltype, gl.Ptr(indices))
}

func cxglBindTexture(target uint32, texture uint32) {
	gl.BindTexture(target, texture)
}

func cxglDeleteTextures(n int32, textures *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.DeleteTextures(n, textures)
}

func cxglGenTextures(n int32, textures *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.GenTextures(n, textures)
}

// cxgl_1_3
func cxglActiveTexture(texture uint32) {
	gl.ActiveTexture(texture)
}

// cxgl_1_4
func cxglBlendFuncSeparate(sfactorRGB uint32, dfactorRGB uint32, sfactorAlpha uint32, dfactorAlpha uint32) {
	gl.BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha)
}

// cxgl_1_5
func cxglBindBuffer(target uint32, buffer uint32) {
	gl.BindBuffer(target, buffer)
}

func cxglDeleteBuffers(n int32, buffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.DeleteBuffers(n, buffers)
}

func cxglGenBuffers(n int32, buffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.GenBuffers(n, buffers)
}

// cxgl_2_0
func cxglBufferData(target uint32, size int, data interface{}, usage uint32) {
	gl.BufferData(target, size, gl.Ptr(data), usage)
}

func cxglBufferSubData(target uint32, offset int, size int, data interface{}) {
	gl.BufferSubData(target, offset, size, gl.Ptr(data))
}

func cxglDrawBuffers(size int32, bufs interface{}) {
	gl.DrawBuffers(size, ito_pui32(bufs))
}

func cxglStencilOpSeparate(face uint32, sfail uint32, dpfail uint32, dppass uint32) {
	gl.StencilOpSeparate(face, sfail, dpfail, dppass)
}

func cxglStencilFuncSeparate(face uint32, glfunc uint32, ref int32, mask uint32) {
	gl.StencilFuncSeparate(face, glfunc, ref, mask)
}

func cxglStencilMaskSeparate(face uint32, mask uint32) {
	gl.StencilMaskSeparate(face, mask)
}

func cxglAttachShader(program uint32, shader uint32) {
	gl.AttachShader(program, shader)
}

func cxglBindAttribLocation(program uint32, index uint32, name string) {
	gl.BindAttribLocation(program, index, *getCString(name, name))
}

func cxglCompileShader(shader uint32) {
	gl.CompileShader(shader)
}

func cxglCreateProgram() uint32 {
	return gl.CreateProgram()
}

func cxglCreateShader(gltype uint32) uint32 {
	return gl.CreateShader(gltype)
}

func cgxlDeleteProgram(program uint32) {
	gl.DeleteProgram(program)
}

func cxglDeleteShader(shader uint32) {
	gl.DeleteShader(shader)
}

func cxglDetachShader(program uint32, shader uint32) {
	gl.DetachShader(program, shader)
}

func cxglEnableVertexAttribArray(index uint32) {
	gl.EnableVertexAttribArray(index)
}

func cxglGetAttribLocation(program uint32, name string) int32 {
	return gl.GetAttribLocation(program, *getCString(name, name))
}

func cxglGetProgramiv(program uint32, pname uint32) int32 {
	var params int32
	gl.GetProgramiv(program, pname, &params)
	return params
}

func cxglGetProgramInfoLog(program uint32, maxLen int32) string {
	log := strings.Repeat("\x00", int(maxLen+1))
	gl.GetProgramInfoLog(program, maxLen, nil, gl.Str(log))
	return log
}

func cxglGetShaderiv(shader uint32, pname uint32) int32 {
	var params int32
	gl.GetShaderiv(shader, pname, &params)
	return params
}

func cxglGetShaderInfoLog(shader uint32, maxLen int32) string {
	log := strings.Repeat("\x00", int(maxLen+1))
	gl.GetShaderInfoLog(shader, maxLen, nil, gl.Str(log))
	return log
}

func cxglGetUniformLocation(program uint32, name string) int32 {
	return gl.GetUniformLocation(program, *getCString(name, name))
}

func cxglLinkProgram(program uint32) {
	gl.LinkProgram(program)
}

func cxglShaderSource(shader uint32, count int32, glstring string) {
	if count > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.ShaderSource(shader, count, getCString(glstring, glstring), nil)
}

func cxglUseProgram(program uint32) {
	gl.UseProgram(program)
}

func cxglUniform1f(location int32, v0 float32) {
	gl.Uniform1f(location, v0)
}

func cxglUniform2f(location int32, v0 float32, v1 float32) {
	gl.Uniform2f(location, v0, v1)
}

func cxglUniform3f(location int32, v0 float32, v1 float32, v2 float32) {
	gl.Uniform3f(location, v0, v1, v2)
}

func cxglUniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	gl.Uniform4f(location, v0, v1, v2, v3)
}

func cxglUniform1i(location int32, v0 int32) {
	gl.Uniform1i(location, v0)
}

func cxglUniform2i(location int32, v0 int32, v1 int32) {
	gl.Uniform2i(location, v0, v1)
}

func cxglUniform3i(location int32, v0 int32, v1 int32, v2 int32) {
	gl.Uniform3i(location, v0, v1, v2)
}

func cxglUniform4i(location int32, v0 int32, v1 int32, v2 int32, v3 int32) {
	gl.Uniform4i(location, v0, v1, v2, v3)
}

func cxglUniform1fv(location int32, count int32, value interface{}) {
	gl.Uniform1fv(location, count, ito_pf32(value))
}

func cxglUniform2fv(location int32, count int32, value interface{}) {
	gl.Uniform2fv(location, count, ito_pf32(value))
}

func cxglUniform3fv(location int32, count int32, value interface{}) {
	gl.Uniform3fv(location, count, ito_pf32(value))
}

func cxglUniform4fv(location int32, count int32, value interface{}) {
	gl.Uniform4fv(location, count, ito_pf32(value))
}

func cxglUniform1iv(location int32, count int32, value interface{}) {
	gl.Uniform1iv(location, count, ito_pi32(value))
}

func cxglUniform2iv(location int32, count int32, value interface{}) {
	gl.Uniform2iv(location, count, ito_pi32(value))
}

func cxglUniform3iv(location int32, count int32, value interface{}) {
	gl.Uniform3iv(location, count, ito_pi32(value))
}

func cxglUniform4iv(location int32, count int32, value interface{}) {
	gl.Uniform4iv(location, count, ito_pi32(value))
}

func cxglUniformMatrix2fv(location int32, count int32, transpose bool, value interface{}) {
	gl.UniformMatrix2fv(location, count, transpose, ito_pf32(value))
}

func cxglUniformMatrix3fv(location int32, count int32, transpose bool, value interface{}) {
	gl.UniformMatrix3fv(location, count, transpose, ito_pf32(value))
}

func cxglUniformMatrix4fv(location int32, count int32, transpose bool, value interface{}) {
	gl.UniformMatrix4fv(location, count, transpose, ito_pf32(value))
}

func cxglVertexAttribPointer(index uint32, size int32, gltype uint32, normalized bool, stride int32, pointer int32) {
	gl.VertexAttribPointer(index, size, gltype, normalized, stride, gl.PtrOffset(int(pointer)))
}

// cxgl_3_0
func cxglClearBufferiv(buffer uint32, drawBuffer int32, value []int32) {
	gl.ClearBufferiv(buffer, drawBuffer, (*int32)(gl.Ptr(&value[0])))
}

func cxglClearBufferuiv(buffer uint32, drawBuffer int32, value []uint32) {
	gl.ClearBufferuiv(buffer, drawBuffer, (*uint32)(gl.Ptr(&value[0])))
}

func cxglClearBufferfv(buffer uint32, drawBuffer int32, value []float32) {
	gl.ClearBufferfv(buffer, drawBuffer, (*float32)(gl.Ptr(&value[0])))
}

func cxglBindRenderbuffer(buffer uint32, renderbuffer uint32) {
	gl.BindRenderbuffer(buffer, renderbuffer)
}

func cxglDeleteRenderbuffers(n int32, renderbuffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.DeleteRenderbuffers(n, renderbuffers)
}

func cxglGenRenderbuffers(n int32, renderbuffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.GenRenderbuffers(n, renderbuffers)
}

func cxglRenderbufferStorage(target uint32, internalFormat uint32, width int32, height int32) {
	gl.RenderbufferStorage(target, internalFormat, width, height)
}

func cxglBindFramebuffer(target uint32, framebuffer uint32) {
	gl.BindFramebuffer(target, framebuffer)
}

func cxglDeleteFramebuffers(n int32, framebuffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.DeleteFramebuffers(n, framebuffers)
}

func cxglGenFramebuffers(n int32, framebuffers *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.GenFramebuffers(n, framebuffers)
}

func cxglCheckFramebufferStatus(target uint32) uint32 {
	return gl.CheckFramebufferStatus(target)
}

func cxglFramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	gl.FramebufferTexture2D(target, attachment, textarget, texture, level)
}

func cxglFramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer uint32) {
	gl.FramebufferRenderbuffer(target, attachment, renderbuffertarget, renderbuffer)
}

func cxglGenerateMipmap(target uint32) {
	gl.GenerateMipmap(target)
}

func cxglBindVertexArrayAPPLE(array uint32) {
	gl.BindVertexArrayAPPLE(array)
}

func cxglBindVertexArray(array uint32) {
	gl.BindVertexArray(array)
}

func cxglDeleteVertexArraysAPPLE(n int32, arrays *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.DeleteVertexArraysAPPLE(n, arrays)
}

func cxglDeleteVertexArrays(n int32, arrays *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.DeleteVertexArrays(n, arrays)
}

func cxglGenVertexArraysAPPLE(n int32, arrays *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.GenVertexArraysAPPLE(n, arrays)
}

func cxglGenVertexArrays(n int32, arrays *uint32) {
	if n > 1 {
		panic(CX_RUNTIME_NOT_IMPLEMENTED)
	}
	gl.GenVertexArrays(n, arrays)
}

// gl_0_0
func opGlMatrixMode(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.MatrixMode(uint32(ReadI32(fp, inp1)))
}

func opGlRotatef(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Rotatef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3), ReadF32(fp, inp4))
}

func opGlTranslatef(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Translatef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func opGlLoadIdentity(expr *CXExpression, fp int) {
	gl.LoadIdentity()
}

func opGlPushMatrix(expr *CXExpression, fp int) {
	gl.PushMatrix()
}

func opGlPopMatrix(expr *CXExpression, fp int) {
	gl.PopMatrix()
}

func opGlEnableClientState(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.EnableClientState(uint32(ReadI32(fp, inp1)))
}

func opGlColor3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Color3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func opGlColor4f(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Color4f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3), ReadF32(fp, inp4))
}

func opGlBegin(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Begin(uint32(ReadI32(fp, inp1)))
}

func opGlEnd(expr *CXExpression, fp int) {
	gl.End()
}

func opGlNormal3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Normal3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func opGlVertex2f(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.Vertex2f(ReadF32(fp, inp1), ReadF32(fp, inp2))
}

func opGlVertex3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Vertex3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func opGlLightfv(expr *CXExpression, fp int) {
	// pointers
	panic("gl.Lightfv")
}

func opGlFrustum(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	gl.Frustum(ReadF64(fp, inp1), ReadF64(fp, inp2), ReadF64(fp, inp3), ReadF64(fp, inp4), ReadF64(fp, inp5), ReadF64(fp, inp6))
}

func opGlTexEnvi(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.TexEnvi(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), ReadI32(fp, inp3))
}

func opGlOrtho(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	gl.Ortho(ReadF64(fp, inp1), ReadF64(fp, inp2), ReadF64(fp, inp3), ReadF64(fp, inp4), ReadF64(fp, inp5), ReadF64(fp, inp6))
}

func opGlScalef(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Scalef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func opGlTexCoord2d(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.TexCoord2d(ReadF64(fp, inp1), ReadF64(fp, inp2))
}

func opGlTexCoord2f(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.TexCoord2f(ReadF32(fp, inp1), ReadF32(fp, inp2))
}
