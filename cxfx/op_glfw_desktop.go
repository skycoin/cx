// +build cxfx,!mobile

package cxfx

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	. "github.com/skycoin/cx/cx"
)

var windows map[string]*glfw.Window = make(map[string]*glfw.Window, 0)

func opGlfwFullscreen(expr *CXExpression, fp int) {
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

func opGlfwInit(expr *CXExpression, fp int) {
	glfw.Init()
	initialized = true
}

func opGlfwSwapBuffers(expr *CXExpression, fp int) {
	windows[ReadStr(fp, expr.Inputs[0])].SwapBuffers()
}

func opGlfwMakeContextCurrent(expr *CXExpression, fp int) {
	windows[ReadStr(fp, expr.Inputs[0])].MakeContextCurrent()
}

func opGlfwWindowHint(expr *CXExpression, fp int) {
	glfw.WindowHint(glfw.Hint(ReadI32(fp, expr.Inputs[0])), int(ReadI32(fp, expr.Inputs[1])))
}

func opGlfwSetInputMode(expr *CXExpression, fp int) {
	windows[ReadStr(fp, expr.Inputs[0])].SetInputMode(
		glfw.InputMode(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])))
}

func opGlfwGetCursorPos(expr *CXExpression, fp int) {
	x, y := windows[ReadStr(fp, expr.Inputs[0])].GetCursorPos()
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), x)
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), y)
}

func opGlfwGetKey(expr *CXExpression, fp int) {
	act := int32(windows[ReadStr(fp, expr.Inputs[0])].GetKey(glfw.Key(ReadI32(fp, expr.Inputs[1]))))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), act)
}

func opGlfwCreateWindow(expr *CXExpression, fp int) {
	if win, err := glfw.CreateWindow(
		int(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])),
		ReadStr(fp, expr.Inputs[3]), nil, nil); err == nil {
		windows[ReadStr(fp, expr.Inputs[0])] = win
	} else {
		panic(err)
	}
}

func opGlfwGetWindowContentScale(expr *CXExpression, fp int) {
	xscale, yscale := windows[ReadStr(fp, expr.Inputs[0])].GetContentScale()
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), xscale)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), yscale)
}

func opGlfwGetMonitorContentScale(expr *CXExpression, fp int) {
	xscale, yscale := glfw.GetPrimaryMonitor().GetContentScale()
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), xscale)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), yscale)
}

func opGlfwSetWindowPos(expr *CXExpression, fp int) {
	windows[ReadStr(fp, expr.Inputs[0])].SetPos(
		int(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])))
}

func opGlfwShouldClose(expr *CXExpression, fp int) {
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), windows[ReadStr(fp, expr.Inputs[0])].ShouldClose())
}

func opGlfwGetFramebufferSize(expr *CXExpression, fp int) {
	width, height := windows[ReadStr(fp, expr.Inputs[0])].GetFramebufferSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func opGlfwGetWindowPos(expr *CXExpression, fp int) {
	x, y := windows[ReadStr(fp, expr.Inputs[0])].GetPos()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(x))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(y))
}

func opGlfwGetWindowSize(expr *CXExpression, fp int) {
	width, height := windows[ReadStr(fp, expr.Inputs[0])].GetSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func opGlfwSwapInterval(expr *CXExpression, fp int) {
	glfw.SwapInterval(int(ReadI32(fp, expr.Inputs[0])))
}

func opGlfwPollEvents(expr *CXExpression, fp int) {
	if initialized {
		glfw.PollEvents()
	}
	PollEvents()
}

func opGlfwGetTime(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	WriteF64(GetFinalOffset(fp, out1), glfw.GetTime())
}

func glfwSetKeyCallback(expr *CXExpression, fp int) {
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			PushKeyboardEvent(ActionType(action), int32(key), int64(scancode), int32(mods))
		})
}

func glfwSetCursorPosCallback(expr *CXExpression, fp int, eventType EventType) {
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetCursorPosCallback(
		func(w *glfw.Window, xpos float64, ypos float64) {
			PushMouseEvent(eventType, ACTION_MOVE, 0, -1, 0, xpos, ypos)
		})
}

func glfwSetMouseButtonCallback(expr *CXExpression, fp int, eventType EventType) {
	window := ReadStr(fp, expr.Inputs[0])
	windows[window].SetMouseButtonCallback(
		func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			x, y := w.GetCursorPos()
			PushMouseEvent(eventType, ActionType(action), int32(key), -1, int32(mods), x, y)
		})
}

func opGlfwSetKeyCallback(expr *CXExpression, fp int) { // TODO : to deprecate
	glfwSetKeyCallback(expr, fp)
	appKeyboardCallback.Init(expr, fp)
}

func opGlfwSetCursorPosCallback(expr *CXExpression, fp int) { // TODO : to deprecate
	glfwSetCursorPosCallback(expr, fp, APP_CURSOR_POS)
	appCursorPositionCallback.Init(expr, fp)
}

func opGlfwSetMouseButtonCallback(expr *CXExpression, fp int) { // TODO : to deprecate
	glfwSetMouseButtonCallback(expr, fp, APP_MOUSE_BUTTON)
	appMouseButtonCallback.Init(expr, fp)
}

func opGlfwSetKeyboardCallback(expr *CXExpression, fp int) {
	glfwSetKeyCallback(expr, fp)
	appKeyboardCallback.InitEx(expr, fp)
}

func opGlfwSetMouseCallback(expr *CXExpression, fp int) {
	glfwSetCursorPosCallback(expr, fp, APP_MOUSE)
	glfwSetMouseButtonCallback(expr, fp, APP_MOUSE)
	appMouseCallback.InitEx(expr, fp)
}

func opGlfwSetFramebufferSizeCallback(expr *CXExpression, fp int) {
	window := ReadStr(fp, expr.Inputs[0])

	windows[window].SetFramebufferSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushFramebufferSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})
	appFramebufferSizeCallback.InitEx(expr, fp)
}

func opGlfwSetWindowSizeCallback(expr *CXExpression, fp int) {
	window := ReadStr(fp, expr.Inputs[0])

	windows[window].SetSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushWindowSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})
	appWindowSizeCallback.InitEx(expr, fp)
}

func opGlfwSetWindowPosCallback(expr *CXExpression, fp int) {
	window := ReadStr(fp, expr.Inputs[0])

	windows[window].SetPosCallback(
		func(w *glfw.Window, x int, y int) {
			PushWindowPositionEvent(float64(x), float64(y)) // TODO to deprecate, use float64
		})
	appWindowPosCallback.InitEx(expr, fp)
}

func opGlfwSetStartCallback(expr *CXExpression, fp int) {
	appStartCallback.InitEx(expr, fp)
	PushEvent(APP_START)
}

func opGlfwSetStopCallback(expr *CXExpression, fp int) {
	appStopCallback.InitEx(expr, fp)
}

func opGlfwSetShouldClose(expr *CXExpression, fp int) {
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
