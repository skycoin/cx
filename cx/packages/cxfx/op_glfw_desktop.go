// +build cxfx,!mobile

package cxfx

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
)

var windows map[string]*glfw.Window = make(map[string]*glfw.Window, 0)

func opGlfwFullscreen(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
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

func opGlfwInit(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	glfw.Init()
	initialized = true
}

func opGlfwSwapBuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	windows[inputs[0].Get_str()].SwapBuffers()
}

func opGlfwMakeContextCurrent(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	windows[inputs[0].Get_str()].MakeContextCurrent()
}

func opGlfwWindowHint(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	glfw.WindowHint(glfw.Hint(inputs[0].Get_i32()), int(inputs[1].Get_i32()))
}

func opGlfwSetInputMode(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	windows[inputs[0].Get_str()].SetInputMode(
		glfw.InputMode(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()))
}

func opGlfwGetCursorPos(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	x, y := windows[inputs[0].Get_str()].GetCursorPos()
	outputs[0].Set_f64(x)
	outputs[1].Set_f64(y)
}

func opGlfwGetKey(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	act := int32(windows[inputs[0].Get_str()].GetKey(glfw.Key(inputs[1].Get_i32())))
	outputs[0].Set_i32(act)
}

func opGlfwCreateWindow(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	if win, err := glfw.CreateWindow(
		int(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()),
		inputs[3].Get_str(), nil, nil); err == nil {
		windows[inputs[0].Get_str()] = win
	} else {
		panic(err)
	}
}

func opGlfwGetWindowContentScale(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	xscale, yscale := windows[inputs[0].Get_str()].GetContentScale()
	outputs[0].Set_f32(xscale)
	outputs[1].Set_f32(yscale)
}

func opGlfwGetMonitorContentScale(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	xscale, yscale := glfw.GetPrimaryMonitor().GetContentScale()
	outputs[0].Set_f32(xscale)
	outputs[1].Set_f32(yscale)
}

func opGlfwSetWindowPos(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	windows[inputs[0].Get_str()].SetPos(
		int(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()))
}

func opGlfwShouldClose(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_bool(windows[inputs[0].Get_str()].ShouldClose())
}

func opGlfwGetFramebufferSize(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	width, height := windows[inputs[0].Get_str()].GetFramebufferSize()
	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGlfwGetWindowPos(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	x, y := windows[inputs[0].Get_str()].GetPos()
	outputs[0].Set_i32(int32(x))
	outputs[1].Set_i32(int32(y))
}

func opGlfwGetWindowSize(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	width, height := windows[inputs[0].Get_str()].GetSize()
	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGlfwSwapInterval(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	glfw.SwapInterval(int(inputs[0].Get_i32()))
}

func opGlfwPollEvents(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	if initialized {
		glfw.PollEvents()
	}
	PollEvents()
}

func opGlfwGetTime(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_f64(glfw.GetTime())
}

func glfwSetKeyCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	window := inputs[0].Get_str()
	windows[window].SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			PushKeyboardEvent(ActionType(action), int32(key), int64(scancode), int32(mods))
		})
}

func glfwSetCursorPosCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue, eventType EventType) {
	window := inputs[0].Get_str()
	windows[window].SetCursorPosCallback(
		func(w *glfw.Window, xpos float64, ypos float64) {
			PushMouseEvent(eventType, ACTION_MOVE, 0, -1, 0, xpos, ypos)
		})
}

func glfwSetMouseButtonCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue, eventType EventType) {
	window := inputs[0].Get_str()
	windows[window].SetMouseButtonCallback(
		func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			x, y := w.GetCursorPos()
			PushMouseEvent(eventType, ActionType(action), int32(key), -1, int32(mods), x, y)
		})
}

func opGlfwSetKeyCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // TODO : to deprecate
	glfwSetKeyCallback(prgrm, inputs, outputs)
	appKeyboardCallback.Init(prgrm, inputs, outputs)
}

func opGlfwSetCursorPosCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // TODO : to deprecate
	glfwSetCursorPosCallback(prgrm, inputs, outputs, APP_CURSOR_POS)
	appCursorPositionCallback.Init(prgrm, inputs, outputs)
}

func opGlfwSetMouseButtonCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // TODO : to deprecate
	glfwSetMouseButtonCallback(prgrm, inputs, outputs, APP_MOUSE_BUTTON)
	appMouseButtonCallback.Init(prgrm, inputs, outputs)
}

func opGlfwSetKeyboardCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	glfwSetKeyCallback(prgrm, inputs, outputs)
	appKeyboardCallback.InitEx(prgrm, inputs, outputs)
}

func opGlfwSetMouseCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	glfwSetCursorPosCallback(prgrm, inputs, outputs, APP_MOUSE)
	glfwSetMouseButtonCallback(prgrm, inputs, outputs, APP_MOUSE)
	appMouseCallback.InitEx(prgrm, inputs, outputs)
}

func opGlfwSetFramebufferSizeCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetFramebufferSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushFramebufferSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})
	appFramebufferSizeCallback.InitEx(prgrm, inputs, outputs)
}

func opGlfwSetWindowSizeCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetSizeCallback(
		func(w *glfw.Window, width int, height int) {
			PushWindowSizeEvent(float64(width), float64(height)) // TODO : to deprecate, use float64
		})
	appWindowSizeCallback.InitEx(prgrm, inputs, outputs)
}

func opGlfwSetWindowPosCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	window := inputs[0].Get_str()

	windows[window].SetPosCallback(
		func(w *glfw.Window, x int, y int) {
			PushWindowPositionEvent(float64(x), float64(y)) // TODO to deprecate, use float64
		})
	appWindowPosCallback.InitEx(prgrm, inputs, outputs)
}

func opGlfwSetStartCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	appStartCallback.InitEx(prgrm, inputs, outputs)
	PushEvent(APP_START)
}

func opGlfwSetStopCallback(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	appStopCallback.InitEx(prgrm, inputs, outputs)
}

func opGlfwSetShouldClose(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	shouldClose := inputs[1].Get_bool()
	windows[inputs[0].Get_str()].SetShouldClose(shouldClose)
}

func getWindowName(w *glfw.Window) []byte {
	for key, win := range windows {
		if w == win {
			var windowHeapPtr = types.AllocWrite_obj_data(ast.PROGRAM.Memory, []byte(key))
			var windowName [types.POINTER_SIZE]byte
			types.Write_ptr(windowName[:], 0, windowHeapPtr)
			return windowName[:]
		}
	}

	return nil
}
