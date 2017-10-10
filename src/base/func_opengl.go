package base

import (
	"fmt"
	// "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/gl/v2.1/gl"
	"runtime"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var freeFns map[string]*func() = make(map[string]*func(), 0)
var cSources map[string]**uint8 = make(map[string]**uint8, 0)

func gl_Init (expr *CXExpression, call *CXCall) error {
	if err := gl.Init(); err == nil {
		return nil
	} else {
		return err
	}
}

func gl_CreateProgram (expr *CXExpression, call *CXCall) error {
	prog := gl.CreateProgram()
	output := encoder.Serialize(int32(prog))

	for _, def := range call.State {
		if def.Name == expr.OutputNames[0].Name {
			def.Value = &output
			return nil
		}
	}
	call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, "i32"))
	return nil
}

func gl_LinkProgram (progId *CXArgument) error {
	if err := checkType("gl.LinkProgram", "i32", progId); err == nil {
		var id int32

		encoder.DeserializeAtomic(*progId.Value, &id)
		
		gl.LinkProgram(uint32(id))
		
		return nil
	} else {
		return err
	}
}

func gl_Clear (mask *CXArgument) error {
	if err := checkType("gl.Clear", "i32", mask); err == nil {
		var m int32

		encoder.DeserializeAtomic(*mask.Value, &m)
		
		gl.Clear(uint32(m))
		return nil
	} else {
		return err
	}
}

func gl_UseProgram (progId *CXArgument) error {
	if err := checkType("gl.Clear", "i32", progId); err == nil {
		var id int32

		encoder.DeserializeAtomic(*progId.Value, &id)
		
		gl.UseProgram(uint32(id))
		return nil
	} else {
		return err
	}
}

func gl_BindBuffer (target, buffer *CXArgument) error {
	if err := checkTwoTypes("gl.BindBuffer", "i32", "i32", target, buffer); err == nil {
		var tgt int32
		var buf int32

		encoder.DeserializeAtomic(*target.Value, &tgt)
		encoder.DeserializeAtomic(*buffer.Value, &buf)
		
		gl.BindBuffer(uint32(tgt), uint32(buf))
		return nil
	} else {
		return err
	}
}

func gl_BindVertexArray (array *CXArgument) error {
	if err := checkType("gl.BindVertexArray", "i32", array); err == nil {
		var arr int32

		encoder.DeserializeAtomic(*array.Value, &arr)

		gl.BindVertexArray(uint32(arr))
		return nil
	} else {
		return err
	}
}

func gl_EnableVertexAttribArray (index *CXArgument) error {
	if err := checkType("gl.EnableVertexAttribArray", "i32", index); err == nil {
		var idx int32

		encoder.DeserializeAtomic(*index.Value, &idx)

		gl.EnableVertexAttribArray(uint32(idx))
		return nil
	} else {
		return err
	}
}

func gl_VertexAttribPointer (index, size, xtype, normalized, stride *CXArgument) error {
	if err := checkFiveTypes("gl.VertexAttribPointer", "i32", "i32", "i32", "bool", "i32", index, size, xtype, normalized, stride); err == nil {
		var idx int32
		var siz int32
		var xtyp int32
		var norm int32 //and later to bool
		var strid int32

		encoder.DeserializeAtomic(*index.Value, &idx)
		encoder.DeserializeAtomic(*size.Value, &siz)
		encoder.DeserializeAtomic(*xtype.Value, &xtyp)
		encoder.DeserializeAtomic(*xtype.Value, &xtyp)
		encoder.DeserializeAtomic(*stride.Value, &strid)

		var normal bool
		if norm == 1 {
			normal = true
		} else {
			normal = false
		}

		gl.VertexAttribPointer(uint32(idx), int32(siz), uint32(xtyp), normal, int32(strid), nil) // fix nil
		return nil
	} else {
		return err
	}
}

func gl_DrawArrays (mode, first, count *CXArgument) error {
	if err := checkThreeTypes("gl.DrawArrays", "i32", "i32", "i32", mode, first, count); err == nil {
		var mod int32
		var fst int32
		var cnt int32

		encoder.DeserializeAtomic(*mode.Value, &mod)
		encoder.DeserializeAtomic(*first.Value, &fst)
		encoder.DeserializeAtomic(*count.Value, &cnt)

		gl.DrawArrays(uint32(mod), fst, cnt)
		return nil
	} else {
		return err
	}
}

func gl_GenBuffers (n, buffers *CXArgument) error {
	if err := checkTwoTypes("gl.GenBuffers", "i32", "i32", n, buffers); err == nil {
		var _n int32
		var bufs int32

		encoder.DeserializeAtomic(*n.Value, &_n)
		encoder.DeserializeAtomic(*buffers.Value, &bufs)

		tmp := uint32(bufs)
		
		gl.GenBuffers(_n, &tmp)

		*buffers.Value = encoder.Serialize(tmp)
		//*buffers.Value = encoder.Serialize(tmp)
		
		return nil
	} else {
		return err
	}
}

func gl_BufferData (target, size, data, usage *CXArgument) error {
	if err := checkFourTypes("gl.BufferData", "i32", "i32", "[]f32", "i32", target, size, data, usage); err == nil {
		var tgt int32
		var siz int32
		var dat []float32
		var usag int32

		encoder.DeserializeAtomic(*target.Value, &tgt)
		encoder.DeserializeAtomic(*size.Value, &siz)
		encoder.DeserializeRaw(*data.Value, &dat)
		encoder.DeserializeAtomic(*usage.Value, &usag)

		gl.BufferData(uint32(tgt), int(siz), gl.Ptr(dat), uint32(usag))
		return nil
	} else {
		return err
	}
}

func gl_GenVertexArrays (n, arrays *CXArgument) error {
	if err := checkTwoTypes("gl.GenVertexArrays", "i32", "i32", n, arrays); err == nil {
		var _n int32
		var arrs int32

		encoder.DeserializeAtomic(*n.Value, &_n)
		encoder.DeserializeAtomic(*arrays.Value, &arrs)

		tmp := uint32(arrs)

		if runtime.GOOS == "darwin" {
			gl.GenVertexArraysAPPLE(_n, &tmp)
		} else {
			gl.GenVertexArrays(_n, &tmp)
		}

		*arrays.Value = encoder.Serialize(tmp)
		
		return nil
	} else {
		return err
	}
}

func gl_CreateShader (xtype *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("gl.CreateShader", "i32", xtype); err == nil {
		var xtyp int32

		encoder.DeserializeAtomic(*xtype.Value, &xtyp)

		shader := gl.CreateShader(uint32(xtyp))
		sShader := encoder.Serialize(int32(shader))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &sShader
				return nil
			}
		}
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sShader, "i32"))

		return nil
	} else {
		return err
	}
}

func gl_Strs (source, freeFn *CXArgument) error {
	if err := checkTwoTypes("gl.Strs", "str", "str", source, freeFn); err == nil {
		fnName := string(*freeFn.Value)
		dsSource := string(*source.Value)
		
		csources, free := gl.Strs(dsSource)

		freeFns[fnName] = &free
		cSources[fnName] = csources
		
		return nil
	} else {
		return err
	}
}

func gl_Free (fnName *CXArgument) error {
	if err := checkType("gl.Free", "str", fnName); err == nil {
		fnName := string(*fnName.Value)

		(*freeFns[fnName])()
		delete(freeFns, fnName)
		delete(cSources, fnName)

		return nil
	} else {
		return err
	}
}

func gl_ShaderSource (shader, count, xstring *CXArgument) error {
	if err := checkThreeTypes("gl.ShaderSource", "i32", "i32", "str", shader, count, xstring); err == nil {
		var shad int32
		var cnt int32

		encoder.DeserializeAtomic(*shader.Value, &shad)
		encoder.DeserializeAtomic(*count.Value, &cnt)

		xstrin := string(*xstring.Value)
		xstr := cSources[xstrin]

		gl.ShaderSource(uint32(shad), cnt, xstr, nil)
		return nil
	} else {
		return err
	}
}

func gl_CompileShader (shader *CXArgument) error {
	if err := checkType("gl.CompileShader", "i32", shader); err == nil {
		var shad int32

		encoder.DeserializeAtomic(*shader.Value, &shad)

		gl.CompileShader(uint32(shad))

		var status int32
		gl.GetShaderiv(uint32(shad), gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetShaderiv(uint32(shad), gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(uint32(shad), logLength, nil, gl.Str(log))

			fmt.Printf("failed to compile: %v", log)
		}
		
		return nil
	} else {
		return err
	}
}

func gl_GetShaderiv (shader, pname, params *CXArgument) error {
	if err := checkThreeTypes("gl.GetShaderiv", "i32", "i32", "i32", shader, pname, params); err == nil {
		var shad int32
		var pnam int32
		var param int32

		encoder.DeserializeAtomic(*shader.Value, &shad)
		encoder.DeserializeAtomic(*pname.Value, &pnam)
		encoder.DeserializeAtomic(*params.Value, &param)

		gl.GetShaderiv(uint32(shad), uint32(pnam), &param)

		*params.Value = encoder.Serialize(param)
		return nil
	} else {
		return err
	}
}

func gl_AttachShader (program, shader *CXArgument) error {
	if err := checkTwoTypes("gl.AttachShader", "i32", "i32", program, shader); err == nil {
		var prog int32
		var shad int32

		encoder.DeserializeAtomic(*program.Value, &prog)
		encoder.DeserializeAtomic(*shader.Value, &shad)

		gl.AttachShader(uint32(prog), uint32(shad))
		return nil
	} else {
		return err
	}
}
