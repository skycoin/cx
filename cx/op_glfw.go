package base

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// declared in func_glfw.go
// var windows map[string]*glfw.Window = make(map[string]*glfw.Window, 0)

func op_glfw_Init (expr *CXExpression, stack *CXStack, fp int) {
	glfw.Init()
}

func op_glfw_WindowHint (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	glfw.WindowHint(glfw.Hint(ReadI32(stack, fp, inp1)), int(ReadI32(stack, fp, inp2)))
}

func op_glfw_SetInputMode (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	windows[ReadStr(stack, fp, inp1)].SetInputMode(glfw.InputMode(ReadI32(stack, fp, inp2)), int(ReadI32(stack, fp, inp3)))
}

func op_glfw_GetCursorPos (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1, out2 := expr.Inputs[0], expr.Outputs[0], expr.Outputs[1]
	x, y := windows[ReadStr(stack, fp, inp1)].GetCursorPos()
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, FromF64(x))
	WriteMemory(stack, GetFinalOffset(stack, fp, out2, MEM_WRITE), out2, FromF64(y))
}

func op_glfw_CreateWindow (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]
	if win, err := glfw.CreateWindow(int(ReadI32(stack, fp, inp2)), int(ReadI32(stack, fp, inp3)), ReadStr(stack, fp, inp4), nil, nil); err == nil {
		windows[ReadStr(stack, fp, inp1)] = win
	} else {
		panic(err)
	}
}

func op_glfw_MakeContextCurrent (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	windows[ReadStr(stack, fp, inp1)].MakeContextCurrent()
}

func op_glfw_ShouldClose (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	if windows[ReadStr(stack, fp, inp1)].ShouldClose() {
		WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, FromBool(true))
	} else {
		WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, FromBool(false))
	}
}

func op_glfw_GetFramebufferSize (expr *CXExpression, stack *CXStack, fp int) {
	inp1, out1, out2 := expr.Inputs[0], expr.Outputs[0], expr.Outputs[1]
	width, height := windows[ReadStr(stack, fp, inp1)].GetFramebufferSize()
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, FromI32(int32(width)))
	WriteMemory(stack, GetFinalOffset(stack, fp, out2, MEM_WRITE), out2, FromI32(int32(height)))
}

func op_glfw_PollEvents () {
	glfw.PollEvents()
}

func op_glfw_SwapBuffers (expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	windows[ReadStr(stack, fp, inp1)].SwapBuffers()
}

func op_glfw_GetTime (expr *CXExpression, stack *CXStack, fp int) {
	out1 := expr.Outputs[0]
	WriteMemory(stack, GetFinalOffset(stack, fp, out1, MEM_WRITE), out1, FromF64(glfw.GetTime()))
}

func op_glfw_SetKeyCallback (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]

	callback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if fn, err := expr.Program.GetFunction(ReadStr(stack, fp, inp2), expr.Package.Name); err == nil {
			_ = fn
			// var winName []byte
			// for key, win := range windows {
			// 	if w == win {
			// 		winName = []byte(key)
			// 		break
			// 	}
			// }

			// sKey := FromI32(int32(key))
			// sScancode := FromI32(int32(scancode))
			// sAction := FromI32(int32(action))
			// sModifierKey := FromI32(int32(mods))

			// MakeCall(fn, nil)




			

			// state := make([]*CXDefinition, len(fn.Inputs))
			// state[0] = MakeDefinition(fn.Inputs[0].Name, &winName,fn.Inputs[0].Typ)
			// state[1] = MakeDefinition(fn.Inputs[1].Name, &sKey,fn.Inputs[1].Typ)
			// state[2] = MakeDefinition(fn.Inputs[2].Name, &sScancode,fn.Inputs[2].Typ)
			// state[3] = MakeDefinition(fn.Inputs[3].Name, &sAction,fn.Inputs[3].Typ)
			// state[4] = MakeDefinition(fn.Inputs[4].Name, &sModifierKey,fn.Inputs[4].Typ)
			
			// subcall := MakeCall(fn, state, call, call.Module, call.Context)
			// call.Context.CallStack.Calls = append(call.Context.CallStack.Calls, subcall)
		}
	}

	windows[ReadStr(stack, fp, inp1)].SetKeyCallback(callback)
}

func op_glfw_SetCursorPosCallback (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]

	callback := func(w *glfw.Window, xpos float64, ypos float64) {
		if fn, err := expr.Program.GetFunction(ReadStr(stack, fp, inp2), expr.Package.Name); err == nil {
			// TODO
			_ = fn
		}
	}

	windows[ReadStr(stack, fp, inp1)].SetCursorPosCallback(callback)
}

func op_glfw_SetShouldClose (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	if ReadBool(stack, fp, inp2) {
		windows[ReadStr(stack, fp, inp1)].SetShouldClose(true)
	} else {
		windows[ReadStr(stack, fp, inp1)].SetShouldClose(false)
	}
}

func op_glfw_SetMouseButtonCallback (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]

	callback := func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if fn, err := expr.Program.GetFunction(ReadStr(stack, fp, inp2), expr.Package.Name); err == nil {
			// TODO
			_ = fn
		}
	}

	windows[ReadStr(stack, fp, inp1)].SetMouseButtonCallback(callback)
}
