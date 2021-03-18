// +build cxfx,!mobile

package cxfx

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	. "github.com/skycoin/cx/cx"
)

var windows map[string]*glfw.Window = make(map[string]*glfw.Window, 0)

func opGlfwFullscreen(inputs []CXValue, outputs []CXValue) {
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

func opGlfwInit(inputs []CXValue, outputs []CXValue) {
	glfw.Init()
	initialized = true
}

func opGlfwSwapBuffers(inputs []CXValue, outputs []CXValue) {
	windows[inputs[0].Get_str()].SwapBuffers()
}

func opGlfwMakeContextCurrent(inputs []CXValue, outputs []CXValue) {
	windows[inputs[0].Get_str()].MakeContextCurrent()
}

func opGlfwWindowHint(inputs []CXValue, outputs []CXValue) {
	glfw.WindowHint(glfw.Hint(inputs[0].Get_i32()), int(inputs[1].Get_i32()))
}

func opGlfwSetInputMode(inputs []CXValue, outputs []CXValue) {
	windows[inputs[0].Get_str()].SetInputMode(
		glfw.InputMode(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()))
}

func opGlfwGetCursorPos(inputs []CXValue, outputs []CXValue) {
	x, y := windows[inputs[0].Get_str()].GetCursorPos()
	outputs[0].Set_f64(x)
	outputs[1].Set_f64(y)
}

func opGlfwGetKey(inputs []CXValue, outputs []CXValue) {
	act := int32(windows[inputs[0].Get_str()].GetKey(glfw.Key(inputs[1].Get_i32())))
	outputs[0].Set_i32(act)
}

func opGlfwCreateWindow(inputs []CXValue, outputs []CXValue) {
	if win, err := glfw.CreateWindow(
		int(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()),
		inputs[3].Get_str(), nil, nil); err == nil {
		windows[inputs[0].Get_str()] = win
	} else {
		panic(err)
	}
}

func opGlfwGetWindowContentScale(inputs []CXValue, outputs []CXValue) {
	xscale, yscale := windows[inputs[0].Get_str()].GetContentScale()
	outputs[0].Set_f32(xscale)
	outputs[1].Set_f32(yscale)
}

func opGlfwGetMonitorContentScale(inputs []CXValue, outputs []CXValue) {
	xscale, yscale := glfw.GetPrimaryMonitor().GetContentScale()
	outputs[0].Set_f32(xscale)
	outputs[1].Set_f32(yscale)
}

func opGlfwSetWindowPos(inputs []CXValue, outputs []CXValue) {
	windows[inputs[0].Get_str()].SetPos(
		int(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()))
}

func opGlfwShouldClose(inputs []CXValue, outputs []CXValue) {
	outputs[0].Set_bool(windows[inputs[0].Get_str()].ShouldClose())
}

func opGlfwGetFramebufferSize(inputs []CXValue, outputs []CXValue) {
	width, height := windows[inputs[0].Get_str()].GetFramebufferSize()
	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGlfwGetWindowPos(inputs []CXValue, outputs []CXValue) {
	x, y := windows[inputs[0].Get_str()].GetPos()
	outputs[0].Set_i32(int32(x))
	outputs[1].Set_i32(int32(y))
}

func opGlfwGetWindowSize(inputs []CXValue, outputs []CXValue) {
	width, height := windows[inputs[0].Get_str()].GetSize()
	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGlfwSwapInterval(inputs []CXValue, outputs []CXValue) {
	glfw.SwapInterval(int(inputs[0].Get_i32()))
}

func opGlfwPollEvents(inputs []CXValue, outputs []CXValue) {
	if initialized {
		glfw.PollEvents()
	}
	PollEvents()
}

func opGlfwGetTime(inputs []CXValue, outputs []CXValue) {
	outputs[0].Set_f64(glfw.GetTime())
}

func glfwSetKeyCallback(inputs []CXValue, outputs []CXValue) {
	window := inputs[0].Get_str()
	windows[window].SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			PushKeyboardEvent(ActionType(action), int32(key), int64(scancode), int32(mods))
		})
}

func glfwSetCursorPosCallback(inputs []CXValue, outputs []CXValue, eventType EventType) {
	window := inputs[0].Get_str()
	windows[window].SetCursorPosCallback(
		func(w *glfw.Window, xpos float64, ypos float64) {
			PushMouseEvent(eventType, ACTION_MOVE, 0, -1, 0, xpos, ypos)
		})
}

func glfwSetMouseButtonCallback(inputs []CXValue, outputs []CXValue, eventType EventType) {
	window := inputs[0].Get_str()
	windows[window].SetMouseButtonCallback(
		func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			x, y := w.GetCursorPos()
			PushMouseEvent(eventType, ActionType(action), int32(key), -1, int32(mods), x, y)
		})
}

func opGlfwSetKeyCallback(inputs []CXValue, outputs []CXValue) { // TODO : to deprecate
	glfwSetKeyCallback(inputs, outputs)
	appKeyboardCallback.Init(inputs, outputs)
}

func opGlfwSetCursorPosCallback(inputs []CXValue, outputs []CXValue) { // TODO : to deprecate
	glfwSetCursorPosCallback(inputs, outputs, APP_CURSOR_POS)
	appCursorPositionCallback.Init(inputs, outputs)
}

func opGlfwSetMouseButtonCallback(inputs []CXValue, outputs []CXValue) { // TODO : to deprecate
	glfwSetMouseButtonCallback(inputs, outputs, APP_MOUSE_BUTTON)
	appMouseButtonCallback.Init(inputs, outputs)
}

func opGlfwSetKeyboardCallback(inputs []CXValue, outputs []CXValue) {
	glfwSetKeyCallback(inputs, outputs)
	appKeyboardCallback.InitEx(inputs, outputs)
}

func opGlfwSetMouseCallback(inputs []CXValue, outputs []CXValue) {
	glfwSetCursorPosCallback(inputs, outputs, APP_MOUSE)
	glfwSetMouseButtonCallback(inputs, outputs, APP_MOUSE)
	appMouseCallback.InitEx(inputs, outputs)
}

func opGlfwSetFramebufferSizeCallback(inputs []CXValue, outputs []CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetFramebufferSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushFramebufferSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})
	appFramebufferSizeCallback.InitEx(inputs, outputs)
}

func opGlfwSetWindowSizeCallback(inputs []CXValue, outputs []CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushWindowSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})
	appWindowSizeCallback.InitEx(inputs, outputs)
}

func opGlfwSetWindowPosCallback(inputs []CXValue, outputs []CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetPosCallback(
		func(w *glfw.Window, x int, y int) {
			PushWindowPositionEvent(float64(x), float64(y)) // TODO to deprecate, use float64
		})
	appWindowPosCallback.InitEx(inputs, outputs)
}

func opGlfwSetStartCallback(inputs []CXValue, outputs []CXValue) {
	appStartCallback.InitEx(inputs, outputs)
	PushEvent(APP_START)
}

func opGlfwSetStopCallback(inputs []CXValue, outputs []CXValue) {
	appStopCallback.InitEx(inputs, outputs)
}

func opGlfwSetShouldClose(inputs []CXValue, outputs []CXValue) {
	shouldClose := inputs[1].Get_bool()
	windows[inputs[0].Get_str()].SetShouldClose(shouldClose)
}

func getWindowName(w *glfw.Window) []byte {
	for key, win := range windows {
		if w == win {
			return FromI32(int32(WriteStringData(key)))
		}
	}

	return nil
}
