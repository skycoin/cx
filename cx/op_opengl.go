package base

import (
	"fmt"
	"os"
	"github.com/go-gl/gl/v2.1/gl"
	"image"
	"image/draw"
	_ "image/png"
	_ "image/jpeg"
	_ "image/gif"
	"runtime"
	"unsafe"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// declared in func_opengl.go
var freeFns map[string]*func() = make(map[string]*func(), 0)
var cSources map[string]**uint8 = make(map[string]**uint8, 0)

func op_gl_Init() {
	gl.Init()
}

func op_gl_GetError(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI32(int32(gl.GetError()))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_BindAttribLocation(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	xstr := cSources[ReadStr(fp, inp3)]
	gl.BindAttribLocation(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), *xstr)
}

func op_gl_GetAttribLocation(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	xstr := cSources[ReadStr(fp, inp2)]
	outB1 := FromI32(gl.GetAttribLocation(uint32(ReadI32(fp, inp1)), *xstr))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_CullFace(expr *CXExpression, fp int) {
    inp1 := expr.Inputs[0]
    gl.CullFace(uint32(ReadI32(fp, inp1)))
}

func op_gl_CreateProgram(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI32(int32(gl.CreateProgram()))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_LinkProgram(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.LinkProgram(uint32(ReadI32(fp, inp1)))
}

func op_gl_Clear(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Clear(uint32(ReadI32(fp, inp1)))
}

func op_gl_UseProgram(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.UseProgram(uint32(ReadI32(fp, inp1)))
}

func op_gl_BindBuffer(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.BindBuffer(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_Viewport(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Viewport(ReadI32(fp, inp1), ReadI32(fp, inp2), ReadI32(fp, inp3), ReadI32(fp, inp4))
}

func op_gl_BindVertexArray(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	if runtime.GOOS == "darwin" {
		gl.BindVertexArrayAPPLE(uint32(ReadI32(fp, inp1)))
	} else {
		gl.BindVertexArray(uint32(ReadI32(fp, inp1)))
	}
}

func op_gl_EnableVertexAttribArray(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.EnableVertexAttribArray(uint32(ReadI32(fp, inp1)))
}

func op_gl_VertexAttribPointer(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	gl.VertexAttribPointer(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), uint32(ReadI32(fp, inp3)), ReadBool(fp, inp4), ReadI32(fp, inp5), unsafe.Pointer(uintptr(ReadI32(fp, inp6))))
}

func op_gl_DrawArrays(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.DrawArrays(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), ReadI32(fp, inp3))
}

func op_gl_GenBuffers(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	tmp := uint32(ReadI32(fp, inp2))
	gl.GenBuffers(ReadI32(fp, inp1), &tmp)
	outB1 := FromI32(int32(tmp))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_BufferData(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.BufferData(uint32(ReadI32(fp, inp1)), int(ReadI32(fp, inp2)), gl.Ptr(ReadF32A(fp, inp3)), uint32(ReadI32(fp, inp4)))
}

func op_gl_GenVertexArrays(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	tmp := uint32(ReadI32(fp, inp2))
	if runtime.GOOS == "darwin" {
		gl.GenVertexArraysAPPLE(ReadI32(fp, inp1), &tmp)
	} else {
		gl.GenVertexArrays(ReadI32(fp, inp1), &tmp)
	}
	outB1 := FromI32(int32(tmp))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_CreateShader(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(gl.CreateShader(uint32(ReadI32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_gl_Strs(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	dsSource := ReadStr(fp, inp1)
	fnName := ReadStr(fp, inp2)

	csources, free := gl.Strs(dsSource + string('\000'))

	freeFns[fnName] = &free
	cSources[fnName] = csources
}

func op_gl_Free(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fnName := ReadStr(fp, inp1)

	(*freeFns[fnName])()
	delete(freeFns, fnName)
	delete(cSources, fnName)
}

func op_gl_ShaderSource(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	xstr := cSources[ReadStr(fp, inp3)]
	gl.ShaderSource(uint32(ReadI32(fp, inp1)), ReadI32(fp, inp2), xstr, nil)
}

func op_gl_CompileShader(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	shad := uint32(ReadI32(fp, inp1))
	gl.CompileShader(shad)

	var status int32
	gl.GetShaderiv(shad, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shad, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shad, logLength, nil, gl.Str(log))

		fmt.Printf("failed to compile: %v", log)
	}
}

func op_gl_GetShaderiv(expr *CXExpression, fp int) {
	// pointers
	panic("gl.GetShaderiv")
}

func op_gl_AttachShader(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.AttachShader(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_LoadIdentity() {
	gl.LoadIdentity()
}

func op_gl_PushMatrix() {
	gl.PushMatrix()
}

func op_gl_PopMatrix() {
	gl.PopMatrix()
}

func op_gl_End() {
	gl.End()
}

func op_gl_Rotatef(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Rotatef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3), ReadF32(fp, inp4))
}

func op_gl_Translatef(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Translatef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_MatrixMode(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.MatrixMode(uint32(ReadI32(fp, inp1)))
}

func op_gl_EnableClientState(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.EnableClientState(uint32(ReadI32(fp, inp1)))
}

func op_gl_BindTexture(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.BindTexture(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_Ortho(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	gl.Ortho(ReadF64(fp, inp1), ReadF64(fp, inp2), ReadF64(fp, inp3), ReadF64(fp, inp4), ReadF64(fp, inp5), ReadF64(fp, inp6))
}

func op_gl_Color3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Color3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_Color4f(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.Color4f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3), ReadF32(fp, inp4))
}

func op_gl_Begin(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Begin(uint32(ReadI32(fp, inp1)))
}

func op_gl_Normal3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Normal3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_TexCoord2f(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.TexCoord2f(ReadF32(fp, inp1), ReadF32(fp, inp2))
}

func op_gl_Vertex2f(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.Vertex2f(ReadF32(fp, inp1), ReadF32(fp, inp2))
}

func op_gl_Vertex3f(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Vertex3f(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_Enable(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Enable(uint32(ReadI32(fp, inp1)))
}

func op_gl_Disable(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.Disable(uint32(ReadI32(fp, inp1)))
}

func op_gl_ClearColor(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	gl.ClearColor(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3), ReadF32(fp, inp4))
}

func op_gl_ClearDepth(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.ClearDepth(ReadF64(fp, inp1))
}

func op_gl_DepthFunc(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.DepthFunc(uint32(ReadI32(fp, inp1)))
}

func op_gl_Lightfv(expr *CXExpression, fp int) {
	// pointers
	panic("gl.Lightfv")
}

func op_gl_Frustum(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	gl.Frustum(ReadF64(fp, inp1), ReadF64(fp, inp2), ReadF64(fp, inp3), ReadF64(fp, inp4), ReadF64(fp, inp5), ReadF64(fp, inp6))
}

func op_gl_NewTexture(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	file := ReadStr(fp, inp1)

	imgFile, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("texture %q not found on disk: %v\n", file, err))
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	
	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	outB1 := encoder.SerializeAtomic(int32(texture))

	WriteMemory(out1Offset, outB1)
}

func op_gl_DepthMask(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	gl.DepthMask(ReadBool(fp, inp1))
}

func op_gl_Scalef(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.Scalef(ReadF32(fp, inp1), ReadF32(fp, inp2), ReadF32(fp, inp3))
}

func op_gl_TexCoord2d(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.TexCoord2d(ReadF64(fp, inp1), ReadF64(fp, inp2))
}

func op_gl_TexEnvi(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	gl.TexEnvi(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)), ReadI32(fp, inp3))
}

func op_gl_BlendFunc(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.BlendFunc(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}

func op_gl_Hint(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	gl.Hint(uint32(ReadI32(fp, inp1)), uint32(ReadI32(fp, inp2)))
}
            
