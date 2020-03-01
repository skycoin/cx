// +build cxfx,mobile

package cxfx

import (
	. "github.com/SkycoinProject/cx/cx"
	"golang.org/x/mobile/app"
	//"golang.org/x/mobile/event/paint"
)

var goapp app.App

func SetGOApp(a app.App) {
	goapp = a
}

func opGlfwFullscreen(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwInit(prgrm *CXProgram) {
	goapp.CreateSurface()
}

func opGlfwSwapBuffers(prgrm *CXProgram) {
	goapp.SwapBuffers()
}

func opGlfwMakeContextCurrent(prgrm *CXProgram) {
	goapp.MakeCurrent()
}

func opGlfwWindowHint(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwSetInputMode(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwGetCursorPos(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), 0.0)
	WriteF64(GetFinalOffset(fp, expr.Outputs[1]), 0.0)
}

func opGlfwGetKey(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
}

func opGlfwCreateWindow(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwGetWindowContentScale(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), 1.0)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), 1.0)
}

func opGlfwGetMonitorContentScale(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), 1.0)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), 1.0)
}

func opGlfwSetWindowPos(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwShouldClose(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
}

func opGlfwGetFramebufferSize(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	width, height := goapp.GetWindowSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func opGlfwGetWindowPos(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), 0)
}

func opGlfwGetWindowSize(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	width, height := goapp.GetWindowSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func opGlfwSwapInterval(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func opGlfwPollEvents(prgrm *CXProgram) {
	PollEvents()
}

func opGlfwGetTime(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), 0.0)
}

func opGlfwSetKeyCallback(prgrm *CXProgram) {
	appKeyCallback.Init(prgrm)
}

func opGlfwSetKeyCallbackEx(prgrm *CXProgram) {
	appKeyCallback.InitEx(prgrm)
}

func opGlfwSetCursorPosCallback(prgrm *CXProgram) {
	appCursorPositionCallback.Init(prgrm)
}

func opGlfwSetCursorPosCallbackEx(prgrm *CXProgram) {
	appCursorPositionCallback.InitEx(prgrm)
}

func opGlfwSetMouseButtonCallback(prgrm *CXProgram) {
	appMouseButtonCallback.Init(prgrm)
}

func opGlfwSetMouseButtonCallbackEx(prgrm *CXProgram) {
	appMouseButtonCallback.InitEx(prgrm)
}

func opGlfwSetTouchCallback(prgrm *CXProgram) {
	appTouchCallback.InitEx(prgrm)
}

func opGlfwSetFramebufferSizeCallback(prgrm *CXProgram) {
	appFramebufferSizeCallback.InitEx(prgrm)
}

func opGlfwSetWindowSizeCallback(prgrm *CXProgram) {
	appWindowSizeCallback.InitEx(prgrm)
}

func opGlfwSetWindowPosCallback(prgrm *CXProgram) {
	appWindowPosCallback.InitEx(prgrm)
}

func opGlfwSetStartCallback(prgrm *CXProgram) {
	appStartCallback.InitEx(prgrm)
}

func opGlfwSetStopCallback(prgrm *CXProgram) {
	appStopCallback.InitEx(prgrm)
}
func opGlfwSetShouldClose(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}
