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

func glfw_CreateWindow (window, width, height, title *CXArgument) error {
	if err := checkThreeTypes("glfw.CreateWindow", "i32", "i32", "str", width, height, title); err == nil {
		var w int32
		var h int32
		var t string = string(*title.Value)
		var winName string = string(*window.Value)

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
		var winName string = string(*window.Value)

		windows[winName].MakeContextCurrent()
		return nil
	} else {
		return err
	}
}

func glfw_ShouldClose (window *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("glfw.ShouldClose", "str", window); err == nil {
		winName := string(*window.Value)

		var output []byte
		if windows[winName].ShouldClose() {
			output = encoder.Serialize(int32(1))
		} else {
			output = encoder.Serialize(int32(0))
		}
		
		assignOutput(&output, "bool", expr, call)
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
		var winName string = string(*window.Value)

		windows[winName].SwapBuffers()
		return nil
	} else {
		return err
	}
}



func glfw_SetKeyCallback (window, fnName *CXArgument, expr *CXExpression, call *CXCall) error {
	wName := string(*window.Value)
	name := string(*fnName.Value)
	
	callback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if fn, err := call.Context.GetFunction(name, expr.Module.Name); err == nil {

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
			
			subcall := MakeCall(fn, state, call, call.Module, call.Context)
			call.Context.CallStack.Calls = append(call.Context.CallStack.Calls, subcall)
		}
	}

	windows[wName].SetKeyCallback(callback)
	return nil
}
