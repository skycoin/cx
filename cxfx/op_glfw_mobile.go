// +build cxfx,mobile

package cxfx

import (
	"golang.org/x/mobile/app"

	. "github.com/skycoin/cx/cx"
	//"golang.org/x/mobile/event/paint"
)

var goapp app.App

func SetGOApp(a app.App) {
	goapp = a
}

func opGlfwFullscreen(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwInit(expr *CXExpression, fp int) {
	goapp.CreateSurface()
}

func opGlfwSwapBuffers(expr *CXExpression, fp int) {
	goapp.SwapBuffers()
}

func opGlfwMakeContextCurrent(expr *CXExpression, fp int) {
	goapp.MakeCurrent()
}

func opGlfwWindowHint(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwSetInputMode(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwGetCursorPos(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), 0.0)
	WriteF64(GetFinalOffset(fp, expr.Outputs[1]), 0.0)
}

func opGlfwGetKey(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
}

func opGlfwCreateWindow(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwGetWindowContentScale(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), 1.0)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), 1.0)
}

func opGlfwGetMonitorContentScale(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), 1.0)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), 1.0)
}

func opGlfwSetWindowPos(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwShouldClose(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
}

func opGlfwGetFramebufferSize(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	width, height := goapp.GetWindowSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func opGlfwGetWindowPos(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), 0)
}

func opGlfwGetWindowSize(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	width, height := goapp.GetWindowSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func opGlfwSwapInterval(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwPollEvents(expr *CXExpression, fp int) {
	PollEvents()
}

func opGlfwGetTime(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), 0.0)
}

func opGlfwSetKeyCallback(expr *CXExpression, fp int) { // TODO : to deprecate
	appKeyboardCallback.Init(prgrm)
}

func opGlfwSetCursorPosCallback(expr *CXExpression, fp int) { // TODO : to deprecate
	appCursorPositionCallback.Init(prgrm)
}

func opGlfwSetMouseButtonCallback(expr *CXExpression, fp int) { // TODO : to deprecate
	appMouseButtonCallback.Init(prgrm)
}

func opGlfwSetKeyboardCallback(expr *CXExpression, fp int) {
	appKeyboardCallback.InitEx(prgrm)
}

func opGlfwSetMouseCallback(expr *CXExpression, fp int) {
	appMouseCallback.InitEx(prgrm)
}

func opGlfwSetFramebufferSizeCallback(expr *CXExpression, fp int) {
	appFramebufferSizeCallback.InitEx(prgrm)
}

func opGlfwSetWindowSizeCallback(expr *CXExpression, fp int) {
	appWindowSizeCallback.InitEx(prgrm)
}

func opGlfwSetWindowPosCallback(expr *CXExpression, fp int) {
	appWindowPosCallback.InitEx(prgrm)
}

func opGlfwSetStartCallback(expr *CXExpression, fp int) {
	appStartCallback.InitEx(prgrm)
}

func opGlfwSetStopCallback(expr *CXExpression, fp int) {
	appStopCallback.InitEx(prgrm)
}
func opGlfwSetShouldClose(expr *CXExpression, fp int) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}
