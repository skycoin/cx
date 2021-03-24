// +build cxfx,mobile

package cxfx

import (
	"golang.org/x/mobile/app"

	"github.com/skycoin/cx/cx"
	//"golang.org/x/mobile/event/paint"
)

var goapp app.App

func SetGOApp(a app.App) {
	goapp = a
}

func opGlfwFullscreen(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwInit(inputs []CXValue, outputs []CXValue) {
	goapp.CreateSurface()
}

func opGlfwSwapBuffers(inputs []CXValue, outputs []CXValue) {
	goapp.SwapBuffers()
}

func opGlfwMakeContextCurrent(inputs []CXValue, outputs []CXValue) {
	goapp.MakeCurrent()
}

func opGlfwWindowHint(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwSetInputMode(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwGetCursorPos(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_f64(0.0)
	outputs[1].Set_f64(0.0)
}

func opGlfwGetKey(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_i32(0)
}

func opGlfwCreateWindow(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwGetWindowContentScale(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_f32(1.0)
	outputs[1].Set_f32(1.0)
}

func opGlfwGetMonitorContentScale(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_f32(1.0)
	outputs[1].Set_f32(1.0)
}

func opGlfwSetWindowPos(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwShouldClose(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_bool(false)
}

func opGlfwGetFramebufferSize(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	width, height := goapp.GetWindowSize()
	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGlfwGetWindowPos(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_i32(0)
	outputs[1].Set_i32(0)
}

func opGlfwGetWindowSize(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	width, height := goapp.GetWindowSize()
	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGlfwSwapInterval(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwPollEvents(inputs []CXValue, outputs []CXValue) {
	PollEvents()
}

func opGlfwGetTime(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	outputs[0].Set_f64(0.0)
}

func opGlfwSetKeyCallback(inputs []CXValue, outputs []CXValue) { // TODO : to deprecate
	appKeyboardCallback.Init(inputs, outputs)
}

func opGlfwSetCursorPosCallback(inputs []CXValue, outputs []CXValue) { // TODO : to deprecate
	appCursorPositionCallback.Init(inputs, outputs)
}

func opGlfwSetMouseButtonCallback(inputs []CXValue, outputs []CXValue) { // TODO : to deprecate
	appMouseButtonCallback.Init(inputs, outputs)
}

func opGlfwSetKeyboardCallback(inputs []CXValue, outputs []CXValue) {
	appKeyboardCallback.InitEx(inputs, outputs)
}

func opGlfwSetMouseCallback(inputs []CXValue, outputs []CXValue) {
	appMouseCallback.InitEx(inputs, outputs)
}

func opGlfwSetFramebufferSizeCallback(inputs []CXValue, outputs []CXValue) {
	appFramebufferSizeCallback.InitEx(inputs, outputs)
}

func opGlfwSetWindowSizeCallback(inputs []CXValue, outputs []CXValue) {
	appWindowSizeCallback.InitEx(inputs, outputs)
}

func opGlfwSetWindowPosCallback(inputs []CXValue, outputs []CXValue) {
	appWindowPosCallback.InitEx(inputs, outputs)
}

func opGlfwSetStartCallback(inputs []CXValue, outputs []CXValue) {
	appStartCallback.InitEx(inputs, outputs)
}

func opGlfwSetStopCallback(inputs []CXValue, outputs []CXValue) {
	appStopCallback.InitEx(inputs, outputs)
}
func opGlfwSetShouldClose(inputs []CXValue, outputs []CXValue) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}
