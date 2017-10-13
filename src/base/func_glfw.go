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
