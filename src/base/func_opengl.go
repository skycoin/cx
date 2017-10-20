package base

import (
	"fmt"
	// "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/gl/v2.1/gl"
	"runtime"
	"strings"
	"os"
	"image"
	"image/draw"
	_ "image/png"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var freeFns map[string]*func() = make(map[string]*func(), 0)
var cSources map[string]**uint8 = make(map[string]**uint8, 0)

func gl_Init () error {
	if err := gl.Init(); err == nil {
		return nil
	} else {
		return err
	}
}

func gl_CreateProgram (expr *CXExpression, call *CXCall) error {
	prog := gl.CreateProgram()
	output := encoder.Serialize(int32(prog))

	assignOutput(&output, "i32", expr, call)
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

// uses pointers. change after implementing cx pointers
func gl_GenBuffers (n, buffers *CXArgument) error {
	if err := checkTwoTypes("gl.GenBuffers", "i32", "i32", n, buffers); err == nil {
		var _n int32
		var bufs int32

		encoder.DeserializeAtomic(*n.Value, &_n)
		encoder.DeserializeAtomic(*buffers.Value, &bufs)

		tmp := uint32(bufs)
		
		gl.GenBuffers(_n, &tmp)

		*buffers.Value = encoder.Serialize(tmp)
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

// uses pointers. change after implementing cx pointers
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

		assignOutput(&sShader, "i32", expr, call)

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

// uses pointers. change after implementing cx pointers
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

func gl_LoadIdentity () error {
	gl.LoadIdentity()
	return nil
}

func gl_PushMatrix () error {
	gl.PushMatrix()
	return nil
}

func gl_PopMatrix () error {
	gl.PopMatrix()
	return nil
}

func gl_Rotatef (angle, x, y, z *CXArgument) error {
	if err := checkFourTypes("gl.Rotatef", "f32", "f32", "f32", "f32", angle, x, y, z); err == nil {
		var dsA float32
		var dsX float32
		var dsY float32
		var dsZ float32

		encoder.DeserializeRaw(*angle.Value, &dsA)
		encoder.DeserializeRaw(*x.Value, &dsX)
		encoder.DeserializeRaw(*y.Value, &dsY)
		encoder.DeserializeRaw(*z.Value, &dsZ)

		gl.Rotatef(dsA, dsX, dsY, dsZ)
		return nil
	} else {
		return err
	}
}

func gl_Translatef (x, y, z *CXArgument) error {
	if err := checkThreeTypes("gl.Translatef", "f32", "f32", "f32", x, y, z); err == nil {
		var dsX float32
		var dsY float32
		var dsZ float32

		encoder.DeserializeRaw(*x.Value, &dsX)
		encoder.DeserializeRaw(*y.Value, &dsY)
		encoder.DeserializeRaw(*z.Value, &dsZ)

		gl.Translatef(dsX, dsY, dsZ)
		return nil
	} else {
		return err
	}
}

func gl_MatrixMode (mode *CXArgument) error {
	if err := checkType("gl.MatrixMode", "i32", mode); err == nil {
		var mod int32

		encoder.DeserializeAtomic(*mode.Value, &mod)

		gl.MatrixMode(uint32(mod))
		return nil
	} else {
		return err
	}
}

func gl_EnableClientState (array *CXArgument) error {
	if err := checkType("gl.EnableClientState", "i32", array); err == nil {
		var arr int32

		encoder.DeserializeAtomic(*array.Value, &arr)

		gl.EnableClientState(uint32(arr))
		return nil
	} else {
		return err
	}
}

func gl_BindTexture (target, texture *CXArgument) error {
	if err := checkTwoTypes("gl.BindTexture", "i32", "i32", target, texture); err == nil {
		var tgt int32
		var tex int32

		encoder.DeserializeAtomic(*target.Value, &tgt)
		encoder.DeserializeAtomic(*texture.Value, &tex)

		gl.BindTexture(uint32(tgt), uint32(tex))
		return nil
	} else {
		return err
	}
}

func gl_Color4f (red, green, blue, alpha *CXArgument) error {
	if err := checkFourTypes("gl.Color4f", "f32", "f32", "f32", "f32", red, green, blue, alpha); err == nil {
		var r float32
		var g float32
		var b float32
		var a float32

		encoder.DeserializeRaw(*red.Value, &r)
		encoder.DeserializeRaw(*green.Value, &g)
		encoder.DeserializeRaw(*blue.Value, &b)
		encoder.DeserializeRaw(*alpha.Value, &a)

		gl.Color4f(r, g, b, a)
		return nil
	} else {
		return err
	}
}

func gl_Begin (mode *CXArgument) error {
	if err := checkType("gl.Begin", "i32", mode); err == nil {
		var mod int32

		encoder.DeserializeAtomic(*mode.Value, &mod)

		gl.Begin(uint32(mod))
		return nil
	} else {
		return err
	}
}

func gl_End () error {
	gl.End()
	return nil
}

func gl_Normal3f (nx, ny, nz *CXArgument) error {
	if err := checkThreeTypes("gl.Normal3f", "f32", "f32", "f32", nx, ny, nz); err == nil {
		var x float32
		var y float32
		var z float32

		encoder.DeserializeRaw(*nx.Value, &x)
		encoder.DeserializeRaw(*ny.Value, &y)
		encoder.DeserializeRaw(*nz.Value, &z)

		gl.Normal3f(x, y, z)
		return nil
	} else {
		return err
	}
}

func gl_TexCoord2f (s, t *CXArgument) error {
	if err := checkTwoTypes("gl.TexCoord2f", "f32", "f32", s, t); err == nil {
		var _s float32
		var _t float32

		encoder.DeserializeRaw(*s.Value, &_s)
		encoder.DeserializeRaw(*t.Value, &_t)

		gl.TexCoord2f(_s, _t)
		return nil
	} else {
		return err
	}
}

func gl_Vertex3f (nx, ny, nz *CXArgument) error {
	if err := checkThreeTypes("gl.Vertex3f", "f32", "f32", "f32", nx, ny, nz); err == nil {
		var x float32
		var y float32
		var z float32

		encoder.DeserializeRaw(*nx.Value, &x)
		encoder.DeserializeRaw(*ny.Value, &y)
		encoder.DeserializeRaw(*nz.Value, &z)

		gl.Vertex3f(x, y, z)
		return nil
	} else {
		return err
	}
}

func gl_Enable (cap *CXArgument) error {
	if err := checkType("gl.Enable", "i32", cap); err == nil {
		var c int32

		encoder.DeserializeAtomic(*cap.Value, &c)

		gl.Enable(uint32(c))
		return nil
	} else {
		return err
	}
}

func gl_Disable (cap *CXArgument) error {
	if err := checkType("gl.Disable", "i32", cap); err == nil {
		var c int32

		encoder.DeserializeAtomic(*cap.Value, &c)

		gl.Disable(uint32(c))
		return nil
	} else {
		return err
	}
}

func gl_ClearColor (red, green, blue, alpha *CXArgument) error {
	if err := checkFourTypes("gl.Color4f", "f32", "f32", "f32", "f32", red, green, blue, alpha); err == nil {
		var r float32
		var g float32
		var b float32
		var a float32

		encoder.DeserializeRaw(*red.Value, &r)
		encoder.DeserializeRaw(*green.Value, &g)
		encoder.DeserializeRaw(*blue.Value, &b)
		encoder.DeserializeRaw(*alpha.Value, &a)

		gl.ClearColor(r, g, b, a)
		return nil
	} else {
		return err
	}
}

func gl_ClearDepth (depth *CXArgument) error {
	if err := checkType("gl.ClearDepth", "f64", depth); err == nil {
		var d float64

		encoder.DeserializeRaw(*depth.Value, &d)

		gl.ClearDepth(d)
		return nil
	} else {
		return err
	}
}

func gl_DepthFunc (xfunc *CXArgument) error {
	if err := checkType("gl.DepthFunc", "i32", xfunc); err == nil {
		var xfn int32

		encoder.DeserializeRaw(*xfunc.Value, &xfn)

		gl.DepthFunc(uint32(xfn))
		return nil
	} else {
		return err
	}
}

// uses pointers. change after implementing cx pointers
func gl_Lightfv (light, pname, params *CXArgument) error {
	if err := checkThreeTypes("gl.Lightfv", "i32", "i32", "f32", light, pname, params); err == nil {
		var ligh int32
		var pnam int32
		var param float32

		encoder.DeserializeAtomic(*light.Value, &ligh)
		encoder.DeserializeAtomic(*pname.Value, &pnam)
		encoder.DeserializeRaw(*params.Value, &param)

		gl.Lightfv(uint32(ligh), uint32(pnam), &param)

		*params.Value = encoder.Serialize(param)
		return nil
	} else {
		return err
	}
}

func gl_Frustum (left, right, bottom, top, zNear, zFar *CXArgument) error {
	if err := checkSixTypes("gl.Frustum", "f64", "f64", "f64", "f64", "f64", "f64", left, right, bottom, top, zNear, zFar); err == nil {
		var l float64
		var r float64
		var b float64
		var t float64
		var zN float64
		var zF float64

		encoder.DeserializeRaw(*left.Value, &l)
		encoder.DeserializeRaw(*right.Value, &r)
		encoder.DeserializeRaw(*bottom.Value, &b)
		encoder.DeserializeRaw(*top.Value, &t)
		encoder.DeserializeRaw(*zNear.Value, &zN)
		encoder.DeserializeRaw(*zFar.Value, &zF)

		gl.Frustum(l, r, b, t, zN, zF)
		return nil
	} else {
		return err
	}
}

func newTexture(file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("texture %q not found on disk: %v\n", file, err))
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	//gl.TexEnvi(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

func gl_NewTexture (file *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("gl.NewTexture", "str", file); err == nil {
		name := string(*file.Value)

		texture := newTexture(name)
		output := encoder.Serialize(int32(texture))
		
		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}

func gl_DepthMask (flag *CXArgument) error {
	if err := checkType("gl.DepthMask", "bool", flag); err == nil {
		var f bool = false
		if (*flag.Value)[0] == 1 {
			f = true
		}

		gl.DepthMask(f)
		return nil
	} else {
		return err
	}
}

func gl_TexEnvi (target, pname, param *CXArgument) error {
	if err := checkThreeTypes("gl.TexEnvi", "i32", "i32", "i32", target, pname, param); err == nil {
		var _target int32
		var _pname int32
		var _param int32

		encoder.DeserializeAtomic(*target.Value, &_target)
		encoder.DeserializeAtomic(*pname.Value, &_pname)
		encoder.DeserializeAtomic(*param.Value, &_param)

		gl.TexEnvi(uint32(_target), uint32(_pname), _param)
		return nil
	} else {
		return err
	}
}

func gl_BlendFunc (sfactor, dfactor *CXArgument) error {
	if err := checkTwoTypes("gl.BlendFunc", "i32", "i32", sfactor, dfactor); err == nil {
		var _sfactor int32
		var _dfactor int32

		encoder.DeserializeAtomic(*sfactor.Value, &_sfactor)
		encoder.DeserializeAtomic(*dfactor.Value, &_dfactor)

		gl.BlendFunc(uint32(_sfactor), uint32(_dfactor))
		return nil
	} else {
		return err
	}
}

func gl_Hint (target, mode *CXArgument) error {
	if err := checkTwoTypes("gl.Hint", "i32", "i32", target, mode); err == nil {
		var _target int32
		var _mode int32

		encoder.DeserializeAtomic(*target.Value, &_target)
		encoder.DeserializeAtomic(*mode.Value, &_mode)

		gl.Hint(uint32(_target), uint32(_mode))
		return nil
	} else {
		return err
	}
}

//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
func Foo () {
	fmt.Println("gl.DITHER", gl.DITHER)
	fmt.Println("gl.POINT_SMOOTH", gl.POINT_SMOOTH)
	fmt.Println("gl.LINE_SMOOTH", gl.LINE_SMOOTH)
	fmt.Println("gl.POLYGON_SMOOTH", gl.POLYGON_SMOOTH)
	fmt.Println("gl.POINT_SMOOTH", gl.POINT_SMOOTH)
	fmt.Println("gl.DONT_CARE", gl.DONT_CARE)
	fmt.Println("gl.POLYGON_SMOOTH_HINT", gl.POLYGON_SMOOTH_HINT)
	fmt.Println("gl.MULTISAMPLE_ARB", gl.MULTISAMPLE_ARB)
	
	// fmt.Println("gl.SRC_ALPHA", gl.SRC_ALPHA)
	// fmt.Println("gl.ONE_MINUS_SRC_ALPHA", gl.ONE_MINUS_SRC_ALPHA)
	
	// fmt.Println("gl.TEXTURE_ENV", gl.TEXTURE_ENV)
	// fmt.Println("gl.TEXTURE_ENV_MODE", gl.TEXTURE_ENV_MODE)
	// fmt.Println("gl.MODULATE", gl.MODULATE)
	// fmt.Println("gl.DECAL", gl.DECAL)
	// fmt.Println("gl.BLEND", gl.BLEND)
	// fmt.Println("gl.REPLACE", gl.REPLACE)
	
	// fmt.Println("gl.BLEND", gl.BLEND)
	// fmt.Println("gl.DEPTH_TEST", gl.DEPTH_TEST)
	// fmt.Println("gl.LIGHTING", gl.LIGHTING)
	// fmt.Println("gl.LEQUAL", gl.LEQUAL)
	// fmt.Println("gl.LIGHT0", gl.LIGHT0)
	// fmt.Println("gl.AMBIENT", gl.AMBIENT)
	// fmt.Println("gl.DIFFUSE", gl.DIFFUSE)
	// fmt.Println("gl.POSITION", gl.POSITION)
}
