// +build opengl

package cxcore

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

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

	glfw.WindowHint(glfw.Hint(ReadI32(fp, expr.Inputs[0])), int(ReadI32(fp, expr.Inputs[1])))
}

func op_glfw_SetInputMode(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	windows[ReadStr(fp, expr.Inputs[0])].SetInputMode(
		glfw.InputMode(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])))
}

func op_glfw_GetCursorPos(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	x, y := windows[ReadStr(fp, expr.Inputs[0])].GetCursorPos()
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), x)
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), y)
}

func op_glfw_GetKey(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	act := int32(windows[ReadStr(fp, expr.Inputs[0])].GetKey(glfw.Key(ReadI32(fp, expr.Inputs[1]))))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), act)
}

func op_glfw_CreateWindow(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	if win, err := glfw.CreateWindow(
		int(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])),
		ReadStr(fp, expr.Inputs[3]), nil, nil); err == nil {
		windows[ReadStr(fp, expr.Inputs[0])] = win
	} else {
		panic(err)
	}
}

func op_glfw_GetWindowContentScale(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	xscale, yscale := windows[ReadStr(fp, expr.Inputs[0])].GetContentScale()
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), xscale)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), yscale)
}

func op_glfw_GetMonitorContentScale(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	xscale, yscale := glfw.GetPrimaryMonitor().GetContentScale()
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), xscale)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), yscale)
}

func op_glfw_SetWindowPos(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	windows[ReadStr(fp, expr.Inputs[0])].SetPos(
		int(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])))
}

func op_glfw_MakeContextCurrent(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	windows[ReadStr(fp, expr.Inputs[0])].MakeContextCurrent()
}

func op_glfw_ShouldClose(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), windows[ReadStr(fp, expr.Inputs[0])].ShouldClose())
}

func op_glfw_GetFramebufferSize(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	width, height := windows[ReadStr(fp, expr.Inputs[0])].GetFramebufferSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
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

	glfw.SwapInterval(int(ReadI32(fp, expr.Inputs[0])))
}

func op_glfw_PollEvents(_ *CXProgram) {
	glfw.PollEvents()
}

func op_glfw_SwapBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	windows[ReadStr(fp, expr.Inputs[0])].SwapBuffers()
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

	glfw_SetKeyCallback(expr, ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]), expr.Package.Name)
}

func op_glfw_SetKeyCallbackEx(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	glfw_SetKeyCallback(expr, ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]), ReadStr(fp, expr.Inputs[2]))
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

	glfw_SetCursorPosCallback(expr, ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]), expr.Package.Name)
}

func op_glfw_SetCursorPosCallbackEx(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	glfw_SetCursorPosCallback(expr, ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]), ReadStr(fp, expr.Inputs[2]))
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

	shouldClose := ReadBool(fp, expr.Inputs[1])
	windows[ReadStr(fp, expr.Inputs[0])].SetShouldClose(shouldClose)
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

	glfw_SetMouseButtonCallback(expr, ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]), expr.Package.Name)
}

func op_glfw_SetMouseButtonCallbantckEx(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	glfw_SetMouseButtonCallback(expr, ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]), ReadStr(fp, expr.Inputs[2]))
}
