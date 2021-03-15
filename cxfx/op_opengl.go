// +build cxfx

package cxfx

import (
	"bufio"
	//"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/skycoin/cx/cx"
)

type Texture struct {
	path   string
	width  int32
	height int32
	level  uint32
	pixels []float32
}

var gifs map[string]*gif.GIF = make(map[string]*gif.GIF, 0)
var textures map[string]Texture = make(map[string]Texture, 0)

func decodeImg(file *os.File, cpuCopy bool) (data []byte, width int32, height int32, pixels []float32) {
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	data = rgba.Pix
	width = int32(rgba.Rect.Size().X)
	height = int32(rgba.Rect.Size().Y)
	if cpuCopy {
		pixels = make([]float32, width*height*4)
		var x int32
		var y int32
		for y = 0; y < height; y++ {
			yoffset := y * width * 4
			for x = 0; x < width; x++ {
				var xoffset = yoffset + x*4
				color := rgba.At(int(x), int(y))
				r, g, b, a := color.RGBA()
				pixels[xoffset] = float32(r) / 65535.0
				pixels[xoffset+1] = float32(g) / 65535.0
				pixels[xoffset+2] = float32(b) / 65535.0
				pixels[xoffset+3] = float32(a) / 65535.0
			}
		}
	}
	return
}

const (
	HDR_NONE = iota
	HDR_32_RLE_RGBE
	MINLEN = 8
	MAXLEN = 0x7fff
	R      = 0
	G      = 1
	B      = 2
	E      = 3
)

func unpack(file *os.File, width int, line []byte) bool {
	if width < MINLEN || width > MAXLEN {
		return unpack_(file, width, line)
	}

	file.Read(line[:4])
	if line[R] != 2 {
		file.Seek(-4, io.SeekCurrent)
		return unpack_(file, width, line)
	}

	if line[G] != 2 || (line[B]&128) != 0 {
		return unpack_(file, width-1, line[4:])
	}

	var b [1]byte
	for i := 0; i < 4; i++ {
		for j := 0; j < width; {
			file.Read(b[:])
			var count int = int(b[0])
			if count > 128 {
				count &= 127
				file.Read(b[:])
				var value int = int(b[0])
				for c := 0; c < count; c++ {
					line[j+c+i] = byte(value)
				}
			} else {
				for c := 0; c < count; c++ {
					offset := j + c + i
					file.Read(line[offset : offset+1])
				}
			}
		}
	}
	return true
}

func unpack_(file *os.File, width int, line []byte) bool {
	var rshift uint
	var repeat [4]byte
	for width > 0 {
		file.Read(line[0:4])
		if line[R] == 1 && line[G] == 1 && line[B] == 1 {
			for i := line[E] << rshift; i > 0; i-- {
				copy(line[0:4], repeat[:])
				line = line[4:]
				width--
			}
			rshift += 8
		} else {
			copy(repeat[:], line[0:4])
			line = line[4:]
			width--
			rshift = 0
		}
	}
	return true
}

func decodeHdr(file *os.File) (data []byte, iwidth int32, iheight int32) {
	data = nil
	iwidth = 0
	iheight = 0

	var format int
	scanner := bufio.NewScanner(file)

	var pos int64
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		pos += int64(advance)
		return
	}

	scanner.Split(scanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "#?RADIANCE" {
		} else if strings.HasPrefix(line, "#") {
		} else if strings.HasPrefix(line, "FORMAT=") {
			var sformat string
			if n, err := fmt.Sscanf(line, "FORMAT=%s\n", &sformat); n != 1 && err != nil {
				fmt.Printf("Failed to scan format : err '%s'\n", err)
				return
			}
			if sformat == "32-bit_rle_rgbe" {
				format = HDR_32_RLE_RGBE
			}
		} else if len(line) == 0 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to scan : err %v\n", scanner.Err())
	}

	if format != HDR_32_RLE_RGBE {
		fmt.Printf("Invalid format %d\n", format)
		return
	}

	file.Seek(pos, 0)

	var width int
	var height int
	if n, err := fmt.Fscanf(file, "-Y %d +X %d\n", &width, &height); n != 2 || err != nil {
		fmt.Printf("Failed to scan width and height : err '%s'\n", err)
		return
	}

	iwidth = int32(width)
	iheight = int32(height)

	//var colors []float32 = make([]float32, width*height*3)
	var line []byte = make([]byte, width*4)
	data = make([]byte, width*height*3*4)

	for y := int(0); y < height; y++ {
		if unpack(file, width, line) == false {
			fmt.Printf("Failed to unpack line %d\n", y)
			return
		}

		yoffset := y /*(height - y - 1)*/ * width * 3 * 4
		for x := 0; x < width; x++ {
			loffset := x * 4
			exponent := math.Pow(2.0, float64(int(line[loffset+3])-128))
			xoffset := yoffset + x*3*4
			r := float32(exponent * float64(line[loffset]) / 256.0)
			g := float32(exponent * float64(line[loffset+1]) / 256.0)
			b := float32(exponent * float64(line[loffset+2]) / 256.0)

			WriteMemF32(data, xoffset, r)
			WriteMemF32(data, xoffset+4, g)
			WriteMemF32(data, xoffset+8, b)
		}
	}
	return
}

func uploadTexture(path string, target uint32, level uint32, cpuCopy bool) {
	file, err := CXOpenFile(path)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf("texture %q not found on disk: %v\n", path, err))
	}

	ext := filepath.Ext(path)
	var data []byte
	var internalFormat int32
	var inputFormat uint32
	var inputType uint32
	var width int32
	var height int32
	var pixels []float32
	if ext == ".png" || ext == ".jpeg" || ext == ".jpg" {
		internalFormat = cxglRGBA8
		inputFormat = cxglRGBA
		inputType = cxglUNSIGNED_BYTE
		data, width, height, pixels = decodeImg(file, cpuCopy)
		if cpuCopy {
		}
		if len(pixels) > 0 {
			var texture Texture
			texture.pixels = pixels
			texture.width = width
			texture.height = height
			texture.path = path
			texture.level = level
			textures[path] = texture
		}
	} else if ext == ".hdr" {
		internalFormat = cxglRGB16F
		inputFormat = cxglRGB
		inputType = cxglFLOAT
		data, width, height = decodeHdr(file)
	}

	if len(data) > 0 {
		cxglTexImage2D(
			target,
			int32(level),
			internalFormat,
			width,
			height,
			0,
			inputFormat,
			inputType,
			data)
	}
}

// gogl
func opGlNewTexture(expr *CXExpression, fp int) {
	var texture uint32
	cxglEnable(cxglTEXTURE_2D)
	cxglGenTextures(1, &texture)
	cxglBindTexture(cxglTEXTURE_2D, texture)
	cxglTexParameteri(cxglTEXTURE_2D, cxglTEXTURE_MIN_FILTER, cxglNEAREST)
	cxglTexParameteri(cxglTEXTURE_2D, cxglTEXTURE_MAG_FILTER, cxglNEAREST)
	cxglTexParameteri(cxglTEXTURE_2D, cxglTEXTURE_WRAP_S, cxglCLAMP_TO_EDGE)
	cxglTexParameteri(cxglTEXTURE_2D, cxglTEXTURE_WRAP_T, cxglCLAMP_TO_EDGE)

	uploadTexture(ReadStr(fp, expr.Inputs[0]), cxglTEXTURE_2D, 0, false)

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(texture))
}

func opGlNewTextureCube(expr *CXExpression, fp int) {
	var texture uint32
	cxglEnable(cxglTEXTURE_CUBE_MAP)
	cxglGenTextures(1, &texture)
	cxglBindTexture(cxglTEXTURE_CUBE_MAP, texture)
	cxglTexParameteri(cxglTEXTURE_CUBE_MAP, cxglTEXTURE_MIN_FILTER, cxglNEAREST)
	cxglTexParameteri(cxglTEXTURE_CUBE_MAP, cxglTEXTURE_MAG_FILTER, cxglNEAREST)
	cxglTexParameteri(cxglTEXTURE_CUBE_MAP, cxglTEXTURE_WRAP_S, cxglCLAMP_TO_EDGE)
	cxglTexParameteri(cxglTEXTURE_CUBE_MAP, cxglTEXTURE_WRAP_T, cxglCLAMP_TO_EDGE)
	cxglTexParameteri(cxglTEXTURE_CUBE_MAP, cxglTEXTURE_WRAP_R, cxglCLAMP_TO_EDGE)

	var faces []string = []string{"posx", "negx", "posy", "negy", "posz", "negz"}
	var pattern string = ReadStr(fp, expr.Inputs[0])
	var extension string = ReadStr(fp, expr.Inputs[1])
	for i := 0; i < 6; i++ {
		uploadTexture(fmt.Sprintf("%s%s%s", pattern, faces[i], extension), uint32(cxglTEXTURE_CUBE_MAP_POSITIVE_X+i), 0, false)
	}
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(texture))
}

func opCxReleaseTexture(expr *CXExpression, fp int) {
	textures[ReadStr(fp, expr.Inputs[0])] = Texture{}
}

func opCxTextureGetPixel(expr *CXExpression, fp int) {
	var r float32
	var g float32
	var b float32
	var a float32

	var x = ReadI32(fp, expr.Inputs[1])
	var y = ReadI32(fp, expr.Inputs[2])

	if texture, ok := textures[ReadStr(fp, expr.Inputs[0])]; ok {
		var yoffset = y * texture.width * 4
		var xoffset = yoffset + x*4
		pixels := texture.pixels
		r = pixels[xoffset]
		g = pixels[xoffset+1]
		b = pixels[xoffset+2]
		a = pixels[xoffset+3]
	}
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), r)
	WriteF32(GetFinalOffset(fp, expr.Outputs[1]), g)
	WriteF32(GetFinalOffset(fp, expr.Outputs[2]), b)
	WriteF32(GetFinalOffset(fp, expr.Outputs[3]), a)
}

func opGlUploadImageToTexture(expr *CXExpression, fp int) {
	uploadTexture(ReadStr(fp, expr.Inputs[0]), uint32(ReadI32(fp, expr.Inputs[1])), uint32(ReadI32(fp, expr.Inputs[2])), ReadBool(fp, expr.Inputs[3]))
}

func opGlNewGIF(expr *CXExpression, fp int) {
	path := ReadStr(fp, expr.Inputs[0])

	file, err := CXOpenFile(path)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf("file not found %q, %v", path, err))
	}

	reader := bufio.NewReader(file)
	gif, err := gif.DecodeAll(reader)
	if err != nil {
		panic(fmt.Sprintf("failed to decode file %q, %v", path, err))
	}

	gifs[path] = gif

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(len(gif.Image)))
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), int32(gif.LoopCount))
	WriteI32(GetFinalOffset(fp, expr.Outputs[2]), int32(gif.Config.Width))
	WriteI32(GetFinalOffset(fp, expr.Outputs[3]), int32(gif.Config.Height))
}

func opGlFreeGIF(expr *CXExpression, fp int) {
	gifs[ReadStr(fp, expr.Inputs[0])] = nil
}

func opGlGIFFrameToTexture(expr *CXExpression, fp int) {
	path := ReadStr(fp, expr.Inputs[0])
	frame := ReadI32(fp, expr.Inputs[1])
	texture := ReadI32(fp, expr.Inputs[2])

	gif := gifs[path]
	img := gif.Image[frame]
	delay := int32(gif.Delay[frame])
	disposal := int32(gif.Disposal[frame])

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	cxglBindTexture(cxglTEXTURE_2D, uint32(texture))
	cxglTexImage2D(
		cxglTEXTURE_2D,
		0,
		cxglRGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		cxglRGBA,
		cxglUNSIGNED_BYTE,
		rgba.Pix)

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), delay)
	WriteI32(GetFinalOffset(fp, expr.Outputs[1]), disposal)
}

func opGlAppend(expr *CXExpression, fp int) {
	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[0])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))

	inputSliceOffset := GetSliceOffset(fp, expr.Inputs[0])
	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	inp1 := expr.Inputs[1]
	obj := ReadMemory(GetFinalOffset(fp, inp1), inp1)

	objLen := int32(len(obj))
	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, inputSliceLen+objLen, 1))
	SliceCopyEx(outputSliceOffset, inputSliceOffset, inputSliceLen+objLen, 1)
	SliceAppendWriteByte(outputSliceOffset, obj, inputSliceLen)
	WriteI32(outputSlicePointer, outputSliceOffset)
}

// gl_1_0
func opGlCullFace(expr *CXExpression, fp int) {
	cxglCullFace(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlFrontFace(expr *CXExpression, fp int) {
	cxglFrontFace(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlHint(expr *CXExpression, fp int) {
	cxglHint(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlScissor(expr *CXExpression, fp int) {
	cxglScissor(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]))
}

func opGlTexParameteri(expr *CXExpression, fp int) {
	cxglTexParameteri(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		ReadI32(fp, expr.Inputs[2]))
}

func opGlTexImage2D(expr *CXExpression, fp int) {
	cxglTexImage2D(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]),
		ReadI32(fp, expr.Inputs[4]),
		ReadI32(fp, expr.Inputs[5]),
		uint32(ReadI32(fp, expr.Inputs[6])),
		uint32(ReadI32(fp, expr.Inputs[7])),
		ReadData(fp, expr.Inputs[8], -1))
}

func opGlClear(expr *CXExpression, fp int) {
	cxglClear(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlClearColor(expr *CXExpression, fp int) {
	cxglClearColor(
		ReadF32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]),
		ReadF32(fp, expr.Inputs[2]),
		ReadF32(fp, expr.Inputs[3]))
}

func opGlClearStencil(expr *CXExpression, fp int) {
	cxglClearStencil(ReadI32(fp, expr.Inputs[0]))
}

func opGlClearDepth(expr *CXExpression, fp int) {
	cxglClearDepth(ReadF64(fp, expr.Inputs[0]))
}

func opGlStencilMask(expr *CXExpression, fp int) {
	cxglStencilMask(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlColorMask(expr *CXExpression, fp int) {
	cxglColorMask(
		ReadBool(fp, expr.Inputs[0]),
		ReadBool(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadBool(fp, expr.Inputs[3]))
}

func opGlDepthMask(expr *CXExpression, fp int) {
	cxglDepthMask(ReadBool(fp, expr.Inputs[0]))
}

func opGlDisable(expr *CXExpression, fp int) {
	cxglDisable(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlEnable(expr *CXExpression, fp int) {
	cxglEnable(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlBlendFunc(expr *CXExpression, fp int) {
	cxglBlendFunc(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlStencilFunc(expr *CXExpression, fp int) {
	cxglStencilFunc(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])))
}

func opGlStencilOp(expr *CXExpression, fp int) {
	cxglStencilOp(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])))
}

func opGlDepthFunc(expr *CXExpression, fp int) {
	cxglDepthFunc(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlGetError(expr *CXExpression, fp int) {
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(cxglGetError()))
}

func opGlGetTexLevelParameteriv(expr *CXExpression, fp int) {
	var outValue int32 = 0
	cxglGetTexLevelParameteriv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])),
		&outValue)
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outValue)
}

func opGlDepthRange(expr *CXExpression, fp int) {
	cxglDepthRange(
		ReadF64(fp, expr.Inputs[0]),
		ReadF64(fp, expr.Inputs[1]))
}

func opGlViewport(expr *CXExpression, fp int) {
	cxglViewport(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]))
}

// gl_1_1
func opGlDrawArrays(expr *CXExpression, fp int) {
	cxglDrawArrays(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]))
}

func opGlDrawElements(expr *CXExpression, fp int) {
	cxglDrawElements(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])),
		ReadData(fp, expr.Inputs[3], -1))
}

func opGlBindTexture(expr *CXExpression, fp int) {
	cxglBindTexture(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlDeleteTextures(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteTextures(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
}

func opGlGenTextures(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenTextures(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

// gl_1_3
func opGlActiveTexture(expr *CXExpression, fp int) {
	cxglActiveTexture(uint32(ReadI32(fp, expr.Inputs[0])))
}

// gl_1_4
func opGlBlendFuncSeparate(expr *CXExpression, fp int) {
	cxglBlendFuncSeparate(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

// gl_1_5
func opGlBindBuffer(expr *CXExpression, fp int) {
	cxglBindBuffer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlDeleteBuffers(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteBuffers(
		ReadI32(fp, expr.Inputs[0]),
		&inpV1) // will panic if inp0 > 1
}

func opGlGenBuffers(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenBuffers(
		ReadI32(fp, expr.Inputs[0]),
		&inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

func opGlBufferData(expr *CXExpression, fp int) {
	cxglBufferData(
		uint32(ReadI32(fp, expr.Inputs[0])),
		int(ReadI32(fp, expr.Inputs[1])),
		ReadData(fp, expr.Inputs[2], -1),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

func opGlBufferSubData(expr *CXExpression, fp int) {
	cxglBufferSubData(
		uint32(ReadI32(fp, expr.Inputs[0])),
		int(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])),
		ReadData(fp, expr.Inputs[3], -1))
}

func opGlDrawBuffers(expr *CXExpression, fp int) {
	cxglDrawBuffers(
		ReadI32(fp, expr.Inputs[0]),
		ReadData_ui32(fp, expr.Inputs[1], TYPE_UI32))
}

func opGlStencilOpSeparate(expr *CXExpression, fp int) {
	cxglStencilOpSeparate(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

func opGlStencilFuncSeparate(expr *CXExpression, fp int) {
	cxglStencilFuncSeparate(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		ReadI32(fp, expr.Inputs[2]),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

func opGlStencilMaskSeparate(expr *CXExpression, fp int) {
	cxglStencilMaskSeparate(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlAttachShader(expr *CXExpression, fp int) {
	cxglAttachShader(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlBindAttribLocation(expr *CXExpression, fp int) {
	cxglBindAttribLocation(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		ReadStr(fp, expr.Inputs[2]))
}

func opGlCompileShader(expr *CXExpression, fp int) {
	shader := uint32(ReadI32(fp, expr.Inputs[0]))
	cxglCompileShader(shader)
}

func opGlCreateProgram(expr *CXExpression, fp int) {
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(cxglCreateProgram()))
}

func opGlCreateShader(expr *CXExpression, fp int) {
	outV0 := int32(cxglCreateShader(uint32(ReadI32(fp, expr.Inputs[0]))))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlDeleteProgram(expr *CXExpression, fp int) {
	cxglDeleteShader(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlDeleteShader(expr *CXExpression, fp int) {
	cxglDeleteShader(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlDetachShader(expr *CXExpression, fp int) {
	cxglDetachShader(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlEnableVertexAttribArray(expr *CXExpression, fp int) {
	cxglEnableVertexAttribArray(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlGetAttribLocation(expr *CXExpression, fp int) {
	outV0 := cxglGetAttribLocation(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadStr(fp, expr.Inputs[1]))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlGetProgramiv(expr *CXExpression, fp int) {
	outV0 := cxglGetProgramiv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlGetProgramInfoLog(expr *CXExpression, fp int) {
	outV0 := cxglGetProgramInfoLog(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]))
	WriteString(fp, outV0, expr.Outputs[0])
}

func opGlGetShaderiv(expr *CXExpression, fp int) {
	outV0 := cxglGetShaderiv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlGetShaderInfoLog(expr *CXExpression, fp int) {
	outV0 := cxglGetShaderInfoLog(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]))
	WriteString(fp, outV0, expr.Outputs[0])
}

func opGlGetUniformLocation(expr *CXExpression, fp int) {
	outV0 := cxglGetUniformLocation(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadStr(fp, expr.Inputs[1]))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlLinkProgram(expr *CXExpression, fp int) {
	program := uint32(ReadI32(fp, expr.Inputs[0]))
	cxglLinkProgram(program)
}

func opGlShaderSource(expr *CXExpression, fp int) {
	cxglShaderSource(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		ReadStr(fp, expr.Inputs[2]))
}

func opGlUseProgram(expr *CXExpression, fp int) {
	cxglUseProgram(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlUniform1f(expr *CXExpression, fp int) {
	cxglUniform1f(
		ReadI32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]))
}

func opGlUniform2f(expr *CXExpression, fp int) {
	cxglUniform2f(
		ReadI32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]),
		ReadF32(fp, expr.Inputs[2]))
}

func opGlUniform3f(expr *CXExpression, fp int) {
	cxglUniform3f(
		ReadI32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]),
		ReadF32(fp, expr.Inputs[2]),
		ReadF32(fp, expr.Inputs[3]))
}

func opGlUniform4f(expr *CXExpression, fp int) {
	cxglUniform4f(
		ReadI32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]),
		ReadF32(fp, expr.Inputs[2]),
		ReadF32(fp, expr.Inputs[3]),
		ReadF32(fp, expr.Inputs[4]))
}

func opGlUniform1i(expr *CXExpression, fp int) {
	cxglUniform1i(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]))
}

func opGlUniform2i(expr *CXExpression, fp int) {
	cxglUniform2i(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]))
}

func opGlUniform3i(expr *CXExpression, fp int) {
	cxglUniform3i(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]))
}

func opGlUniform4i(expr *CXExpression, fp int) {
	cxglUniform4i(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]),
		ReadI32(fp, expr.Inputs[4]))
}

func opGlUniform1fv(expr *CXExpression, fp int) {
	cxglUniform1fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData_f32(fp, expr.Inputs[2], TYPE_F32))
}

func opGlUniform2fv(expr *CXExpression, fp int) {
	cxglUniform2fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData_f32(fp, expr.Inputs[2], TYPE_F32))
}

func opGlUniform3fv(expr *CXExpression, fp int) {
	cxglUniform3fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData_f32(fp, expr.Inputs[2], TYPE_F32))
}

func opGlUniform4fv(expr *CXExpression, fp int) {
	cxglUniform4fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData_f32(fp, expr.Inputs[2], TYPE_F32))
}

func opGlUniform1iv(expr *CXExpression, fp int) {
	cxglUniform1iv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData_i32(fp, expr.Inputs[2], TYPE_I32))
}

func opGlUniform2iv(expr *CXExpression, fp int) {
	cxglUniform2iv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData_i32(fp, expr.Inputs[2], TYPE_I32))
}

func opGlUniform3iv(expr *CXExpression, fp int) {
	cxglUniform3iv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData_i32(fp, expr.Inputs[2], TYPE_I32))
}

func opGlUniform4iv(expr *CXExpression, fp int) {
	cxglUniform4iv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData_i32(fp, expr.Inputs[2], TYPE_I32))
}

func opGlUniformMatrix2fv(expr *CXExpression, fp int) {
	cxglUniformMatrix2fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadData_f32(fp, expr.Inputs[3], TYPE_F32))
}

func opGlUniformMatrix3fv(expr *CXExpression, fp int) {
	cxglUniformMatrix3fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadData_f32(fp, expr.Inputs[3], TYPE_F32))
}

func opGlUniformMatrix4fv(expr *CXExpression, fp int) {
	cxglUniformMatrix4fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadData_f32(fp, expr.Inputs[3], TYPE_F32))
}

func opGlUniformV4F(expr *CXExpression, fp int) {
	cxglUniform4fv(
		ReadI32(fp, expr.Inputs[0]),
		1,
		ReadData(fp, expr.Inputs[1], -1))
}

func opGlUniformM44F(expr *CXExpression, fp int) {
	cxglUniformMatrix4fv(
		ReadI32(fp, expr.Inputs[0]),
		1,
		ReadBool(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], -1))
}

func opGlUniformM44FV(expr *CXExpression, fp int) {
	cxglUniformMatrix4fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadData(fp, expr.Inputs[3], -1))
}

func opGlVertexAttribPointer(expr *CXExpression, fp int) {
	cxglVertexAttribPointer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])),
		ReadBool(fp, expr.Inputs[3]),
		ReadI32(fp, expr.Inputs[4]), 0)
}

func opGlVertexAttribPointerI32(expr *CXExpression, fp int) {
	cxglVertexAttribPointer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])),
		ReadBool(fp, expr.Inputs[3]),
		ReadI32(fp, expr.Inputs[4]),
		ReadI32(fp, expr.Inputs[5]))
}

func opGlClearBufferI(expr *CXExpression, fp int) {
	color := []int32{
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]),
		ReadI32(fp, expr.Inputs[4]),
		ReadI32(fp, expr.Inputs[5])}

	cxglClearBufferiv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		color)
}

func opGlClearBufferUI(expr *CXExpression, fp int) {
	color := []uint32{
		ReadUI32(fp, expr.Inputs[2]),
		ReadUI32(fp, expr.Inputs[3]),
		ReadUI32(fp, expr.Inputs[4]),
		ReadUI32(fp, expr.Inputs[5])}

	cxglClearBufferuiv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		color)
}

func opGlClearBufferF(expr *CXExpression, fp int) {
	color := []float32{
		ReadF32(fp, expr.Inputs[2]),
		ReadF32(fp, expr.Inputs[3]),
		ReadF32(fp, expr.Inputs[4]),
		ReadF32(fp, expr.Inputs[5])}

	cxglClearBufferfv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		color)
}

func opGlBindRenderbuffer(expr *CXExpression, fp int) {
	cxglBindRenderbuffer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlDeleteRenderbuffers(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteRenderbuffers(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
}

func opGlGenRenderbuffers(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenRenderbuffers(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

func opGlRenderbufferStorage(expr *CXExpression, fp int) {
	cxglRenderbufferStorage(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]))
}

func opGlBindFramebuffer(expr *CXExpression, fp int) {
	cxglBindFramebuffer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlDeleteFramebuffers(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteFramebuffers(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
}

func opGlGenFramebuffers(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenFramebuffers(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

func opGlCheckFramebufferStatus(expr *CXExpression, fp int) {
	outV0 := int32(cxglCheckFramebufferStatus(uint32(ReadI32(fp, expr.Inputs[0]))))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlFramebufferTexture2D(expr *CXExpression, fp int) {
	cxglFramebufferTexture2D(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])),
		uint32(ReadI32(fp, expr.Inputs[3])),
		ReadI32(fp, expr.Inputs[4]))
}

func opGlFramebufferRenderbuffer(expr *CXExpression, fp int) {
	cxglFramebufferRenderbuffer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

func opGlGenerateMipmap(expr *CXExpression, fp int) {
	cxglGenerateMipmap(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlBindVertexArray(expr *CXExpression, fp int) {
	inpV0 := uint32(ReadI32(fp, expr.Inputs[0]))
	if runtime.GOOS == "darwin" {
		cxglBindVertexArrayAPPLE(inpV0)
	} else {
		cxglBindVertexArray(inpV0)
	}
}

func opGlBindVertexArrayCore(expr *CXExpression, fp int) {
	cxglBindVertexArray(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlDeleteVertexArrays(expr *CXExpression, fp int) {
	inpV0 := ReadI32(fp, expr.Inputs[0])
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	if runtime.GOOS == "darwin" {
		cxglDeleteVertexArraysAPPLE(inpV0, &inpV1) // will panic if inp0 > 1
	} else {
		cxglDeleteVertexArrays(inpV0, &inpV1) // will panic if inp0 > 1
	}
}

func opGlDeleteVertexArraysCore(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteVertexArrays(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
}

func opGlGenVertexArrays(expr *CXExpression, fp int) {
	inpV0 := ReadI32(fp, expr.Inputs[0])
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	if runtime.GOOS == "darwin" {
		cxglGenVertexArraysAPPLE(inpV0, &inpV1) // will panic if inp0 > 1
	} else {
		cxglGenVertexArrays(inpV0, &inpV1) // will panic if inp0 > 1
	}
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

func opGlGenVertexArraysCore(expr *CXExpression, fp int) {
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenVertexArrays(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}
