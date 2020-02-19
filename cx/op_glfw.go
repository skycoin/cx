// +build opengl

package cxcore

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// declared in func_glfw.go
var windows map[string]*glfw.Window = make(map[string]*glfw.Window, 0)

func op_glfw_Fullscreen(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	window := windows[ReadStr(fp, expr.Inputs[0])]
	fullscreen := ReadBool(fp, expr.Inputs[1])
	x := ReadI32(fp, expr.Inputs[2])
	y := ReadI32(fp, expr.Inputs[3])
	w := ReadI32(fp, expr.Inputs[4])
	h := ReadI32(fp, expr.Inputs[5])

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	if fullscreen {
		window.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		window.SetMonitor(nil, int(x), int(y), int(w), int(h), mode.RefreshRate)
	}
	window.MakeContextCurrent()
}

func op_glfw_Init(prgrm *CXProgram) {
	glfw.Init()
}

func op_glfw_WindowHint(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	glfw.WindowHint(glfw.Hint(ReadI32(fp, inp1)), int(ReadI32(fp, inp2)))
}

func op_glfw_SetInputMode(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	windows[ReadStr(fp, inp1)].SetInputMode(glfw.InputMode(ReadI32(fp, inp2)), int(ReadI32(fp, inp3)))
}

func op_glfw_GetCursorPos(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1, out2 := expr.Inputs[0], expr.Outputs[0], expr.Outputs[1]
	x, y := windows[ReadStr(fp, inp1)].GetCursorPos()
	WriteF64(GetFinalOffset(fp, out1), x)
	WriteF64(GetFinalOffset(fp, out2), y)
}

func op_glfw_GetKey(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	act := int32(windows[ReadStr(fp, inp1)].GetKey(glfw.Key(ReadI32(fp, inp2))))

	WriteI32(GetFinalOffset(fp, out1), act)
}

func op_glfw_CreateWindow(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	if win, err := glfw.CreateWindow(int(ReadI32(fp, inp2)), int(ReadI32(fp, inp3)), ReadStr(fp, inp4), nil, nil); err == nil {
		windows[ReadStr(fp, inp1)] = win
	} else {
		panic(err)
	}
}

func op_glfw_SetWindowPos(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	windows[ReadStr(fp, inp1)].SetPos(int(ReadI32(fp, inp2)), int(ReadI32(fp, inp3)))
}

func op_glfw_MakeContextCurrent(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1 := expr.Inputs[0]
	windows[ReadStr(fp, inp1)].MakeContextCurrent()
}

func op_glfw_ShouldClose(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	WriteBool(GetFinalOffset(fp, out1), windows[ReadStr(fp, inp1)].ShouldClose())
}

func op_glfw_GetFramebufferSize(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1, out2 := expr.Inputs[0], expr.Outputs[0], expr.Outputs[1]
	width, height := windows[ReadStr(fp, inp1)].GetFramebufferSize()
	WriteI32(GetFinalOffset(fp, out1), int32(width))
	WriteI32(GetFinalOffset(fp, out2), int32(height))
}

func op_glfw_GetWindowPos(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	x, y := windows[ReadStr(fp, expr.Inputs[0])].GetPos()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(x))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(y))
}

func op_glfw_GetWindowSize(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	width, height := windows[ReadStr(fp, expr.Inputs[0])].GetSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func op_glfw_SwapInterval(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1 := expr.Inputs[0]
	glfw.SwapInterval(int(ReadI32(fp, inp1)))
}

func op_glfw_PollEvents(_ *CXProgram) {
	glfw.PollEvents()
}

func op_glfw_SwapBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1 := expr.Inputs[0]
	windows[ReadStr(fp, inp1)].SwapBuffers()
}

func op_glfw_GetTime(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	out1 := expr.Outputs[0]
	WriteF64(GetFinalOffset(fp, out1), glfw.GetTime())
}

func glfw_SetKeyCallback(expr *CXExpression, window string, functionName string, packageName string) {
	callback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		var inps [][]byte = make([][]byte, 5)
		inps[0] = GetWindowName(w)
		inps[1] = FromI32(int32(key))
		inps[2] = FromI32(int32(scancode))
		inps[3] = FromI32(int32(action))
		inps[4] = FromI32(int32(mods))
		PROGRAM.ccallback(expr, functionName, packageName, inps)
	}

	windows[window].SetKeyCallback(callback)
}

func op_glfw_SetKeyCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp0, inp1 := expr.Inputs[0], expr.Inputs[1]
	glfw_SetKeyCallback(expr, ReadStr(fp, inp0), ReadStr(fp, inp1), expr.Package.Name)
}

func op_glfw_SetKeyCallbackEx(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp0, inp1, inp2 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	glfw_SetKeyCallback(expr, ReadStr(fp, inp0), ReadStr(fp, inp1), ReadStr(fp, inp2))
}

func GetWindowName(w *glfw.Window) []byte {
	for key, win := range windows {
		if w == win {
			return FromI32(int32(newwriteObj(FromStr(key))))
		}
	}

	return nil
}

func glfw_SetCursorPosCallback(expr *CXExpression, window string, functionName string, packageName string) {
	callback := func(w *glfw.Window, xpos float64, ypos float64) {
		var inps [][]byte = make([][]byte, 3)
		inps[0] = GetWindowName(w)
		inps[1] = FromF64(xpos)
		inps[2] = FromF64(ypos)
		PROGRAM.ccallback(expr, functionName, packageName, inps)
	}

	windows[window].SetCursorPosCallback(callback)
}

func op_glfw_SetCursorPosCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp0, inp1 := expr.Inputs[0], expr.Inputs[1]
	glfw_SetCursorPosCallback(expr, ReadStr(fp, inp0), ReadStr(fp, inp1), expr.Package.Name)
}

func op_glfw_SetCursorPosCallbackEx(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp0, inp1, inp2 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	glfw_SetCursorPosCallback(expr, ReadStr(fp, inp0), ReadStr(fp, inp1), ReadStr(fp, inp2))
}

func op_glfw_SetFramebufferSizeCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	functionName := ReadStr(fp, expr.Inputs[1])
	packageName := ReadStr(fp, expr.Inputs[2])
	callback := func(w *glfw.Window, width int, height int) {
		var inps [][]byte = make([][]byte, 3)
		inps[0] = GetWindowName(w)
		inps[1] = FromI32(int32(width))
		inps[2] = FromI32(int32(height))
		PROGRAM.ccallback(expr, functionName, packageName, inps)
	}
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetFramebufferSizeCallback(callback)
}

func op_glfw_SetWindowPosCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	functionName := ReadStr(fp, expr.Inputs[1])
	packageName := ReadStr(fp, expr.Inputs[2])
	callback := func(w *glfw.Window, x int, y int) {
		var inps [][]byte = make([][]byte, 3)
		inps[0] = GetWindowName(w)
		inps[1] = FromI32(int32(x))
		inps[2] = FromI32(int32(y))
		PROGRAM.ccallback(expr, functionName, packageName, inps)
	}
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetPosCallback(callback)
}

func op_glfw_SetWindowSizeCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	functionName := ReadStr(fp, expr.Inputs[1])
	packageName := ReadStr(fp, expr.Inputs[2])
	callback := func(w *glfw.Window, width int, height int) {
		var inps [][]byte = make([][]byte, 3)
		inps[0] = GetWindowName(w)
		inps[1] = FromI32(int32(width))
		inps[2] = FromI32(int32(height))
		PROGRAM.ccallback(expr, functionName, packageName, inps)
	}
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetSizeCallback(callback)
}

func op_glfw_SetShouldClose(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	if ReadBool(fp, inp2) {
		windows[ReadStr(fp, inp1)].SetShouldClose(true)
	} else {
		windows[ReadStr(fp, inp1)].SetShouldClose(false)
	}
}

func glfw_SetMouseButtonCallback(expr *CXExpression, window string, functionName string, packageName string) {
	callback := func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		var inps [][]byte = make([][]byte, 4)
		inps[0] = GetWindowName(w)
		inps[1] = FromI32(int32(key))
		inps[2] = FromI32(int32(action))
		inps[3] = FromI32(int32(mods))
		PROGRAM.ccallback(expr, functionName, packageName, inps)
	}

	windows[window].SetMouseButtonCallback(callback)
}

func op_glfw_SetMouseButtonCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp0, inp1 := expr.Inputs[0], expr.Inputs[1]
	glfw_SetMouseButtonCallback(expr, ReadStr(fp, inp0), ReadStr(fp, inp1), expr.Package.Name)
}

func op_glfw_SetMouseButtonCallbackEx(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp0, inp1, inp2 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	glfw_SetMouseButtonCallback(expr, ReadStr(fp, inp0), ReadStr(fp, inp1), ReadStr(fp, inp2))
}

type Func_i32_i32 func(a int32, b int32)

var Functions_i32_i32 []Func_i32_i32

func op_glfw_func_i32_i32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	packageName := ReadStr(fp, inp1)
	functionName := ReadStr(fp, inp2)
	callback := func(a int32, b int32) {
		var inps [][]byte = make([][]byte, 2)
		inps[0] = FromI32(a)
		inps[1] = FromI32(b)
		PROGRAM.ccallback(expr, functionName, packageName, inps)
	}

	Functions_i32_i32 = append(Functions_i32_i32, callback)
	WriteI32(GetFinalOffset(fp, out1), int32(len(Functions_i32_i32)-1))
}

func op_glfw_call_i32_i32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	index := ReadI32(fp, inp1)
	count := int32(len(Functions_i32_i32))
	if index >= 0 && index < count {
		Functions_i32_i32[index](ReadI32(fp, inp2), ReadI32(fp, inp3))
	}
}
