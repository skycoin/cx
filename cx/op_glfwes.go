// +build opengles

package cxcore

import (
	"golang.org/x/mobile/app"
	//"golang.org/x/mobile/event/paint"
)

var goapp app.App

func SetGOApp(a app.App) {
	goapp = a
}

func op_glfw_Fullscreen(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_Init(prgrm *CXProgram) {
	goapp.CreateSurface()
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SwapBuffers(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	//goapp.Publish()
	goapp.SwapBuffers()
}

func op_glfw_MakeContextCurrent(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	goapp.MakeCurrent()
}

func op_glfw_WindowHint(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetInputMode(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_GetCursorPos(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), 0.0)
	WriteF64(GetFinalOffset(fp, expr.Outputs[1]), 0.0)
}

func op_glfw_GetKey(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
}

func op_glfw_CreateWindow(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_GetWindowContentScale(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), 1.0)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), 1.0)
}

func op_glfw_GetMonitorContentScale(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), 1.0)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), 1.0)
}

func op_glfw_SetWindowPos(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_ShouldClose(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
}

func op_glfw_GetFramebufferSize(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	width, height := goapp.GetWindowSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func op_glfw_GetWindowPos(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), 0)
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), 0)
}

func op_glfw_GetWindowSize(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	width, height := goapp.GetWindowSize()
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(height))
}

func op_glfw_SwapInterval(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_PollEvents(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_GetTime(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), 0.0)
}

func glfw_SetKeyCallback(expr *CXExpression, window string, functionName string, packageName string) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetKeyCallback(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetKeyCallbackEx(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func glfw_SetCursorPosCallback(expr *CXExpression, window string, functionName string, packageName string) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetCursorPosCallback(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetCursorPosCallbackEx(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetFramebufferSizeCallback(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetWindowPosCallback(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetWindowSizeCallback(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetShouldClose(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func glfw_SetMouseButtonCallback(expr *CXExpression, window string, functionName string, packageName string) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetMouseButtonCallback(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}

func op_glfw_SetMouseButtonCallbackEx(prgrm *CXProgram) {
	//panic(CX_RUNTIME_NOT_IMPLEMENTED)
}
