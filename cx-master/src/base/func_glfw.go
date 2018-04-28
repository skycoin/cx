package base

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var windows map[string]*glfw.Window = make(map[string]*glfw.Window, 0)

func glfw_Init () error {
	if err := glfw.Init(); err == nil {
		return nil
	} else {
		return err
	}
}

func glfw_WindowHint (target *CXArgument, hint *CXArgument) error {
	if err := checkTwoTypes("glfw.WindowHint", "i32", "i32", target, hint); err == nil {
		var tgt int32
		var h int32

		encoder.DeserializeAtomic(*target.Value, &tgt)
		encoder.DeserializeAtomic(*hint.Value, &h)

		glfw.WindowHint(glfw.Hint(tgt), int(h))
		return nil
	} else {
		return err
	}
}

func glfw_SetInputMode (window, mode, value *CXArgument) error {
	if err := checkThreeTypes("glfw.SetInputMode", "str", "i32", "i32", window, mode, value); err == nil {
		var winName string
		var _mode int32
		var _value int32

		encoder.DeserializeRaw(*window.Value, &winName)
		encoder.DeserializeAtomic(*mode.Value, &_mode)
		encoder.DeserializeAtomic(*value.Value, &_value)

		windows[winName].SetInputMode(glfw.InputMode(_mode), int(_value))
		return nil
	} else {
		return err
	}
}

func glfw_GetCursorPos (window *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("glfw.GetCursorPos", "str", window); err == nil {
		var winName string
		encoder.DeserializeRaw(*window.Value, &winName)
		
		x, y := windows[winName].GetCursorPos()

		sX := encoder.Serialize(x)
		sY := encoder.Serialize(y)

		assignOutput(0, sX, "f64", expr, call)
		assignOutput(1, sY, "f64", expr, call)
		
		return nil
	} else {
		return err
	}
}

func glfw_CreateWindow (window, width, height, title *CXArgument) error {
	if err := checkThreeTypes("glfw.CreateWindow", "i32", "i32", "str", width, height, title); err == nil {
		var w int32
		var h int32
		var t string
		var winName string
		encoder.DeserializeRaw(*window.Value, &winName)
		encoder.DeserializeRaw(*title.Value, &t)

		encoder.DeserializeAtomic(*width.Value, &w)
		encoder.DeserializeAtomic(*height.Value, &h)

		if win, err := glfw.CreateWindow(int(w), int(h), t, nil, nil); err == nil {
			windows[winName] = win
		} else {
			return err
		}
		return nil
	} else {
		return err
	}
}

func glfw_MakeContextCurrent (window *CXArgument) error {
	if err := checkType("glfw.MakeContextCurrent", "str", window); err == nil {
		var winName string
		encoder.DeserializeRaw(*window.Value, &winName)

		windows[winName].MakeContextCurrent()
		return nil
	} else {
		return err
	}
}



func glfw_ShouldClose (window *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("glfw.ShouldClose", "str", window); err == nil {
		var winName string
		encoder.DeserializeRaw(*window.Value, &winName)

		var output []byte
		if windows[winName].ShouldClose() {
			output = encoder.Serialize(int32(1))
		} else {
			output = encoder.Serialize(int32(0))
		}
		
		assignOutput(0, output, "bool", expr, call)
		return nil
	} else {
		return err
	}
}

func glfw_GetFramebufferSize (window *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("glfw.GetFramebufferSize", "str", window); err == nil {
		var winName string
		encoder.DeserializeRaw(*window.Value, &winName)

		width, height := windows[winName].GetFramebufferSize()
		sWidth := encoder.Serialize(int32(width))
		sHeight := encoder.Serialize(int32(height))
		
		assignOutput(0, sWidth, "i32", expr, call)
		assignOutput(1, sHeight, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func glfw_PollEvents () error {
	glfw.PollEvents()
	return nil
}

func glfw_SwapBuffers (window *CXArgument) error {
	if err := checkType("glfw.SwapBuffers", "str", window); err == nil {
		var winName string
		encoder.DeserializeRaw(*window.Value, &winName)

		windows[winName].SwapBuffers()
		return nil
	} else {
		return err
	}
}

func glfw_GetTime (expr *CXExpression, call *CXCall) error {
	time := glfw.GetTime()
	sTime := encoder.Serialize(time)
	assignOutput(0, sTime, "f64", expr, call)
	return nil
}

func glfw_SetKeyCallback (window, fnName *CXArgument, expr *CXExpression, call *CXCall) error {
	var wName string
	var name string
	encoder.DeserializeRaw(*window.Value, &wName)
	encoder.DeserializeRaw(*fnName.Value, &name)
	
	callback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if fn, err := call.Program.GetFunction(name, expr.Package.Name); err == nil {

			var winName []byte
			for key, win := range windows {
				if w == win {
					winName = []byte(key)
					break
				}
			}

			sKey := encoder.Serialize(int32(key))
			sScancode := encoder.Serialize(int32(scancode))
			sAction := encoder.Serialize(int32(action))
			sModifierKey := encoder.Serialize(int32(mods))

			state := make([]*CXDefinition, len(fn.Inputs))
			state[0] = MakeDefinition(fn.Inputs[0].Name, &winName,fn.Inputs[0].Typ)
			state[1] = MakeDefinition(fn.Inputs[1].Name, &sKey,fn.Inputs[1].Typ)
			state[2] = MakeDefinition(fn.Inputs[2].Name, &sScancode,fn.Inputs[2].Typ)
			state[3] = MakeDefinition(fn.Inputs[3].Name, &sAction,fn.Inputs[3].Typ)
			state[4] = MakeDefinition(fn.Inputs[4].Name, &sModifierKey,fn.Inputs[4].Typ)
			
			subcall := MakeCall(fn, state, call, call.Package, call.Program)
			call.Program.CallStack = append(call.Program.CallStack, *subcall)
		}
	}

	windows[wName].SetKeyCallback(callback)
	return nil
}

func glfw_SetCursorPosCallback (window, fnName *CXArgument, expr *CXExpression, call *CXCall) error {
	var wName string
	var name string
	encoder.DeserializeRaw(*window.Value, &wName)
	encoder.DeserializeRaw(*fnName.Value, &name)
	
	callback := func(w *glfw.Window, xpos float64, ypos float64) {
		if fn, err := call.Program.GetFunction(name, expr.Package.Name); err == nil {
			
			var winName []byte
			for key, win := range windows {
				if w == win {
					winName = []byte(key)
					break
				}
			}

			sXpos := encoder.Serialize(xpos)
			sYpos := encoder.Serialize(ypos)

			state := make([]*CXDefinition, len(fn.Inputs))

			state[0] = MakeDefinition(fn.Inputs[0].Name, &winName,fn.Inputs[0].Typ)
			state[1] = MakeDefinition(fn.Inputs[1].Name, &sXpos,fn.Inputs[1].Typ)
			state[2] = MakeDefinition(fn.Inputs[2].Name, &sYpos,fn.Inputs[2].Typ)
			
			subcall := MakeCall(fn, state, call, call.Package, call.Program)
			call.Program.CallStack = append(call.Program.CallStack, *subcall)
		}
	}

	windows[wName].SetCursorPosCallback(callback)
	return nil
}

func glfw_SetShouldClose (window, value *CXArgument) error {
	if err := checkTwoTypes("glfw.SetShouldClose", "str", "bool", window, value); err == nil {
		var _window string
		var _value int32

		encoder.DeserializeRaw(*window.Value, &_window)
		encoder.DeserializeRaw(*value.Value, &_value)

		if _value > 0 {
			windows[_window].SetShouldClose(true)
		} else {
			windows[_window].SetShouldClose(false)
		}
		return nil
	} else {
		return err
	}
}

func glfw_SetMouseButtonCallback (window, fnName *CXArgument, expr *CXExpression, call *CXCall) error {
	var wName string
	var name string
	encoder.DeserializeRaw(*window.Value, &wName)
	encoder.DeserializeRaw(*fnName.Value, &name)
	
	callback := func(w *glfw.Window, key glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if fn, err := call.Program.GetFunction(name, expr.Package.Name); err == nil {

			var winName []byte
			for key, win := range windows {
				if w == win {
					winName = []byte(key)
					break
				}
			}

			sKey := encoder.Serialize(int32(key))
			sAction := encoder.Serialize(int32(action))
			sModifierKey := encoder.Serialize(int32(mods))

			state := make([]*CXDefinition, len(fn.Inputs))

			state[0] = MakeDefinition(fn.Inputs[0].Name, &winName,fn.Inputs[0].Typ)
			state[1] = MakeDefinition(fn.Inputs[1].Name, &sKey,fn.Inputs[1].Typ)
			state[2] = MakeDefinition(fn.Inputs[2].Name, &sAction,fn.Inputs[2].Typ)
			state[3] = MakeDefinition(fn.Inputs[3].Name, &sModifierKey,fn.Inputs[3].Typ)
			
			subcall := MakeCall(fn, state, call, call.Package, call.Program)
			call.Program.CallStack = append(call.Program.CallStack, *subcall)
		}
	}

	windows[wName].SetMouseButtonCallback(callback)
	return nil
}

func glfw_GetKey (window, key *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("glfw.GetKey", "str", "i32", window, key); err == nil {
		var wName string
		var _key int32
		
		encoder.DeserializeRaw(*window.Value, &wName)
		encoder.DeserializeAtomic(*key.Value, &_key)

		action := windows[wName].GetKey(glfw.Key(_key))
		output := encoder.SerializeAtomic(int32(action))
		assignOutput(0, output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func Bar () {
	//fmt.Println("glfw.CURSOR", glfw.CursorMode)
}
