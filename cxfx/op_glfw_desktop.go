// +build cxfx,!mobile

package cxfx

import (
	. "github.com/SkycoinProject/cx/cx"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var windows map[string]*glfw.Window = make(map[string]*glfw.Window, 0)

func opGlfwFullscreen(prgrm *CXProgram) {
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

var initialized bool

func opGlfwInit(prgrm *CXProgram) {
	glfw.Init()
	initialized = true
}

func opGlfwSwapBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	windows[ReadStr(fp, expr.Inputs[0])].SwapBuffers()
}

func opGlfwMakeContextCurrent(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	windows[ReadStr(fp, expr.Inputs[0])].MakeContextCurrent()
}

func opGlfwWindowHint(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	glfw.WindowHint(glfw.Hint(ReadI32(fp, expr.Inputs[0])), int(ReadI32(fp, expr.Inputs[1])))
}

func opGlfwSetInputMode(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	windows[ReadStr(fp, expr.Inputs[0])].SetInputMode(
		glfw.InputMode(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])))
}

func opGlfwGetCursorPos(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	x, y := windows[ReadStr(fp, expr.Inputs[0])].GetCursorPos()
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), x)
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), y)
}

func opGlfwGetKey(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	act := int32(windows[ReadStr(fp, expr.Inputs[0])].GetKey(glfw.Key(ReadI32(fp, expr.Inputs[1]))))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), act)
}

func opGlfwCreateWindow(prgrm *CXProgram) {
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

func opGlfwGetWindowContentScale(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	xscale, yscale := windows[ReadStr(fp, expr.Inputs[0])].GetContentScale()
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), xscale)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), yscale)
}

func opGlfwGetMonitorContentScale(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	xscale, yscale := glfw.GetPrimaryMonitor().GetContentScale()
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), xscale)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), yscale)
}

func opGlfwSetWindowPos(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	windows[ReadStr(fp, expr.Inputs[0])].SetPos(
		int(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])))
}

func opGlfwShouldClose(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), windows[ReadStr(fp, expr.Inputs[0])].ShouldClose())
}

func opGlfwGetFramebufferSize(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	width, height := windows[ReadStr(fp, expr.Inputs[0])].GetFramebufferSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func opGlfwGetWindowPos(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	x, y := windows[ReadStr(fp, expr.Inputs[0])].GetPos()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(x))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(y))
}

func opGlfwGetWindowSize(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	width, height := windows[ReadStr(fp, expr.Inputs[0])].GetSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func opGlfwSwapInterval(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	glfw.SwapInterval(int(ReadI32(fp, expr.Inputs[0])))
}

func opGlfwPollEvents(_ *CXProgram) {
	if initialized {
		glfw.PollEvents()
	}
	PollEvents()
}

func opGlfwGetTime(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	out1 := expr.Outputs[0]
	WriteF64(GetFinalOffset(fp, out1), glfw.GetTime())
}

func glfwSetKeyCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			PushKeyEvent(int32(key), int32(scancode), int32(action), int32(mods))
		})
}

func glfwSetCursorPosCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetCursorPosCallback(
		func(w *glfw.Window, xpos float64, ypos float64) {
			PushCursorPositionEvent(xpos, ypos)
		})
}

func glfwSetMouseButtonCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetMouseButtonCallback(
		func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			PushMouseButtonEvent(int32(key), int32(action), int32(mods))
		})
}

func opGlfwSetKeyCallback(prgrm *CXProgram) {
	glfwSetKeyCallback(prgrm)
	appKeyCallback.Init(prgrm)
}

func opGlfwSetKeyCallbackEx(prgrm *CXProgram) {
	glfwSetKeyCallback(prgrm)
	appKeyCallback.InitEx(prgrm)
}

func opGlfwSetCursorPosCallback(prgrm *CXProgram) {
	glfwSetCursorPosCallback(prgrm)
	appCursorPositionCallback.Init(prgrm)
}

func opGlfwSetCursorPosCallbackEx(prgrm *CXProgram) {
	glfwSetCursorPosCallback(prgrm)
	appCursorPositionCallback.InitEx(prgrm)
}

func opGlfwSetMouseButtonCallback(prgrm *CXProgram) {
	glfwSetMouseButtonCallback(prgrm)
	appMouseButtonCallback.Init(prgrm)
}

func opGlfwSetMouseButtonCallbackEx(prgrm *CXProgram) {
	glfwSetMouseButtonCallback(prgrm)
	appMouseButtonCallback.InitEx(prgrm)
}

func opGlfwSetTouchCallback(prgrm *CXProgram) {
	appTouchCallback.InitEx(prgrm)
}

func opGlfwSetFramebufferSizeCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	window := ReadStr(fp, expr.Inputs[0])

	windows[window].SetFramebufferSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushFramebufferSizeEvent(int32(width), int32(height))
		})
	appFramebufferSizeCallback.InitEx(prgrm)
}

func opGlfwSetWindowSizeCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	window := ReadStr(fp, expr.Inputs[0])

	windows[window].SetSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushWindowSizeEvent(int32(width), int32(height))
		})
	appWindowSizeCallback.InitEx(prgrm)
}

func opGlfwSetWindowPosCallback(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	window := ReadStr(fp, expr.Inputs[0])

	windows[window].SetPosCallback(
		func(w *glfw.Window, x int, y int) {
			PushWindowPositionEvent(int32(x), int32(y))
		})
	appWindowPosCallback.InitEx(prgrm)
}

func opGlfwSetStartCallback(prgrm *CXProgram) {
	appStartCallback.InitEx(prgrm)
	PushEvent(APP_START)
}

func opGlfwSetStopCallback(prgrm *CXProgram) {
	appStopCallback.InitEx(prgrm)
}

func opGlfwSetShouldClose(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	shouldClose := ReadBool(fp, expr.Inputs[1])
	windows[ReadStr(fp, expr.Inputs[0])].SetShouldClose(shouldClose)
}

func getWindowName(w *glfw.Window) []byte {
	for key, win := range windows {
		if w == win {
			return FromI32(int32(NewWriteObj(FromStr(key))))
		}
	}

	return nil
}
