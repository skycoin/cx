// +build cxfx,!mobile

package cxfx

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/skycoin/cx/cx"
)

var windows map[string]*glfw.Window = make(map[string]*glfw.Window, 0)

func opGlfwFullscreen(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	window := windows[inputs[0].Get_str()]
	fullscreen := inputs[1].Get_bool()
	x := inputs[2].Get_i32()
	y := inputs[3].Get_i32()
	w := inputs[4].Get_i32()
	h := inputs[5].Get_i32()

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

func opGlfwInit(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	glfw.Init()
	initialized = true
}

func opGlfwSwapBuffers(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	windows[inputs[0].Get_str()].SwapBuffers()
}

func opGlfwMakeContextCurrent(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	windows[inputs[0].Get_str()].MakeContextCurrent()
}

func opGlfwWindowHint(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	glfw.WindowHint(glfw.Hint(inputs[0].Get_i32()), int(inputs[1].Get_i32()))
}

func opGlfwSetInputMode(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	windows[inputs[0].Get_str()].SetInputMode(
		glfw.InputMode(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()))
}

func opGlfwGetCursorPos(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	x, y := windows[inputs[0].Get_str()].GetCursorPos()
	outputs[0].Set_f64(x)
	outputs[1].Set_f64(y)
}

func opGlfwGetKey(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	act := int32(windows[inputs[0].Get_str()].GetKey(glfw.Key(inputs[1].Get_i32())))
	outputs[0].Set_i32(act)
}

func opGlfwCreateWindow(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	if win, err := glfw.CreateWindow(
		int(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()),
		inputs[3].Get_str(), nil, nil); err == nil {
		windows[inputs[0].Get_str()] = win
	} else {
		panic(err)
	}
}

func opGlfwGetWindowContentScale(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	xscale, yscale := windows[inputs[0].Get_str()].GetContentScale()
	outputs[0].Set_f32(xscale)
	outputs[1].Set_f32(yscale)
}

func opGlfwGetMonitorContentScale(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	xscale, yscale := glfw.GetPrimaryMonitor().GetContentScale()
	outputs[0].Set_f32(xscale)
	outputs[1].Set_f32(yscale)
}

func opGlfwSetWindowPos(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	windows[inputs[0].Get_str()].SetPos(
		int(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()))
}

func opGlfwShouldClose(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	outputs[0].Set_bool(windows[inputs[0].Get_str()].ShouldClose())
}

func opGlfwGetFramebufferSize(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	width, height := windows[inputs[0].Get_str()].GetFramebufferSize()
	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGlfwGetWindowPos(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	x, y := windows[inputs[0].Get_str()].GetPos()
	outputs[0].Set_i32(int32(x))
	outputs[1].Set_i32(int32(y))
}

func opGlfwGetWindowSize(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	width, height := windows[inputs[0].Get_str()].GetSize()
	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGlfwSwapInterval(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	glfw.SwapInterval(int(inputs[0].Get_i32()))
}

func opGlfwPollEvents(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	if initialized {
		glfw.PollEvents()
	}
	PollEvents()
}

func opGlfwGetTime(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	outputs[0].Set_f64(glfw.GetTime())
}

func glfwSetKeyCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	window := inputs[0].Get_str()
	windows[window].SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			PushKeyboardEvent(ActionType(action), int32(key), int64(scancode), int32(mods))
		})
}

func glfwSetCursorPosCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue, eventType EventType) {
	window := inputs[0].Get_str()
	windows[window].SetCursorPosCallback(
		func(w *glfw.Window, xpos float64, ypos float64) {
			PushMouseEvent(eventType, ACTION_MOVE, 0, -1, 0, xpos, ypos)
		})
}

func glfwSetMouseButtonCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue, eventType EventType) {
	window := inputs[0].Get_str()
	windows[window].SetMouseButtonCallback(
		func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			x, y := w.GetCursorPos()
			PushMouseEvent(eventType, ActionType(action), int32(key), -1, int32(mods), x, y)
		})
}

func opGlfwSetKeyCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) { // TODO : to deprecate
	glfwSetKeyCallback(inputs, outputs)
	appKeyboardCallback.Init(inputs, outputs)
}

func opGlfwSetCursorPosCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) { // TODO : to deprecate
	glfwSetCursorPosCallback(inputs, outputs, APP_CURSOR_POS)
	appCursorPositionCallback.Init(inputs, outputs)
}

func opGlfwSetMouseButtonCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) { // TODO : to deprecate
	glfwSetMouseButtonCallback(inputs, outputs, APP_MOUSE_BUTTON)
	appMouseButtonCallback.Init(inputs, outputs)
}

func opGlfwSetKeyboardCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	glfwSetKeyCallback(inputs, outputs)
	appKeyboardCallback.InitEx(inputs, outputs)
}

func opGlfwSetMouseCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	glfwSetCursorPosCallback(inputs, outputs, APP_MOUSE)
	glfwSetMouseButtonCallback(inputs, outputs, APP_MOUSE)
	appMouseCallback.InitEx(inputs, outputs)
}

func opGlfwSetFramebufferSizeCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetFramebufferSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushFramebufferSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})
	appFramebufferSizeCallback.InitEx(inputs, outputs)
}

func opGlfwSetWindowSizeCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushWindowSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})
	appWindowSizeCallback.InitEx(inputs, outputs)
}

func opGlfwSetWindowPosCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetPosCallback(
		func(w *glfw.Window, x int, y int) {
			PushWindowPositionEvent(float64(x), float64(y)) // TODO to deprecate, use float64
		})
	appWindowPosCallback.InitEx(inputs, outputs)
}

func opGlfwSetStartCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	appStartCallback.InitEx(inputs, outputs)
	PushEvent(APP_START)
}

func opGlfwSetStopCallback(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	appStopCallback.InitEx(inputs, outputs)
}

func opGlfwSetShouldClose(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	shouldClose := inputs[1].Get_bool()
	windows[inputs[0].Get_str()].SetShouldClose(shouldClose)
}

func getWindowName(w *glfw.Window) []byte {
	for key, win := range windows {
		if w == win {
			return cxcore.FromI32(int32(cxcore.WriteStringData(key)))
		}
	}

	return nil
}
