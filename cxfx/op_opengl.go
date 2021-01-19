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
func opGlNewTexture(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opGlNewTextureCube(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opCxReleaseTexture(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	textures[ReadStr(fp, expr.Inputs[0])] = Texture{}
}

func opCxTextureGetPixel(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opGlUploadImageToTexture(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	uploadTexture(ReadStr(fp, expr.Inputs[0]), uint32(ReadI32(fp, expr.Inputs[1])), uint32(ReadI32(fp, expr.Inputs[2])), ReadBool(fp, expr.Inputs[3]))
}

func opGlNewGIF(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opGlFreeGIF(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	gifs[ReadStr(fp, expr.Inputs[0])] = nil
}

func opGlGIFFrameToTexture(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opGlAppend(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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
func opGlCullFace(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglCullFace(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlFrontFace(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglFrontFace(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlHint(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglHint(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlScissor(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglScissor(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]))
}

func opGlTexParameteri(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglTexParameteri(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		ReadI32(fp, expr.Inputs[2]))
}

func opGlTexImage2D(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opGlClear(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglClear(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlClearColor(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglClearColor(
		ReadF32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]),
		ReadF32(fp, expr.Inputs[2]),
		ReadF32(fp, expr.Inputs[3]))
}

func opGlClearStencil(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglClearStencil(ReadI32(fp, expr.Inputs[0]))
}

func opGlClearDepth(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglClearDepth(ReadF64(fp, expr.Inputs[0]))
}

func opGlStencilMask(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglStencilMask(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlColorMask(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglColorMask(
		ReadBool(fp, expr.Inputs[0]),
		ReadBool(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadBool(fp, expr.Inputs[3]))
}

func opGlDepthMask(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDepthMask(ReadBool(fp, expr.Inputs[0]))
}

func opGlDisable(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	cxglDisable(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlEnable(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglEnable(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlBlendFunc(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBlendFunc(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlStencilFunc(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglStencilFunc(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])))
}

func opGlStencilOp(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglStencilOp(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])))
}

func opGlDepthFunc(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDepthFunc(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlGetError(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(cxglGetError()))
}

func opGlGetTexLevelParameteriv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	var outValue int32 = 0
	cxglGetTexLevelParameteriv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])),
		&outValue)
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outValue)
}

func opGlDepthRange(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDepthRange(
		ReadF64(fp, expr.Inputs[0]),
		ReadF64(fp, expr.Inputs[1]))
}

func opGlViewport(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglViewport(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]))
}

// gl_1_1
func opGlDrawArrays(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDrawArrays(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]))
}

func opGlDrawElements(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDrawElements(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])),
		ReadData(fp, expr.Inputs[3], -1))
}

func opGlBindTexture(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBindTexture(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlDeleteTextures(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteTextures(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
}

func opGlGenTextures(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenTextures(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

// gl_1_3
func opGlActiveTexture(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglActiveTexture(uint32(ReadI32(fp, expr.Inputs[0])))
}

// gl_1_4
func opGlBlendFuncSeparate(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBlendFuncSeparate(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

// gl_1_5
func opGlBindBuffer(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBindBuffer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlDeleteBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteBuffers(
		ReadI32(fp, expr.Inputs[0]),
		&inpV1) // will panic if inp0 > 1
}

func opGlGenBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenBuffers(
		ReadI32(fp, expr.Inputs[0]),
		&inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

func opGlBufferData(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBufferData(
		uint32(ReadI32(fp, expr.Inputs[0])),
		int(ReadI32(fp, expr.Inputs[1])),
		ReadData(fp, expr.Inputs[2], -1),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

func opGlBufferSubData(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBufferSubData(
		uint32(ReadI32(fp, expr.Inputs[0])),
		int(ReadI32(fp, expr.Inputs[1])),
		int(ReadI32(fp, expr.Inputs[2])),
		ReadData(fp, expr.Inputs[3], -1))
}

func opGlDrawBuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDrawBuffers(
		ReadI32(fp, expr.Inputs[0]),
		ReadData(fp, expr.Inputs[1], TYPE_UI32))
}

func opGlStencilOpSeparate(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglStencilOpSeparate(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

func opGlStencilFuncSeparate(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglStencilFuncSeparate(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		ReadI32(fp, expr.Inputs[2]),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

func opGlStencilMaskSeparate(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglStencilMaskSeparate(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlAttachShader(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglAttachShader(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlBindAttribLocation(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBindAttribLocation(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		ReadStr(fp, expr.Inputs[2]))
}

func opGlCompileShader(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	shader := uint32(ReadI32(fp, expr.Inputs[0]))
	cxglCompileShader(shader)
}

func opGlCreateProgram(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(cxglCreateProgram()))
}

func opGlCreateShader(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := int32(cxglCreateShader(uint32(ReadI32(fp, expr.Inputs[0]))))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlDeleteProgram(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDeleteShader(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlDeleteShader(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDeleteShader(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlDetachShader(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglDetachShader(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlEnableVertexAttribArray(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglEnableVertexAttribArray(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlGetAttribLocation(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := cxglGetAttribLocation(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadStr(fp, expr.Inputs[1]))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlGetProgramiv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := cxglGetProgramiv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlGetProgramInfoLog(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := cxglGetProgramInfoLog(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(outV0))
}

func opGlGetShaderiv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := cxglGetShaderiv(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlGetShaderInfoLog(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := cxglGetShaderInfoLog(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]))
	WriteObject(GetFinalOffset(fp, expr.Outputs[0]), FromStr(outV0))
}

func opGlGetUniformLocation(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := cxglGetUniformLocation(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadStr(fp, expr.Inputs[1]))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlLinkProgram(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	program := uint32(ReadI32(fp, expr.Inputs[0]))
	cxglLinkProgram(program)
}

func opGlShaderSource(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglShaderSource(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		ReadStr(fp, expr.Inputs[2]))
}

func opGlUseProgram(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUseProgram(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlUniform1f(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform1f(
		ReadI32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]))
}

func opGlUniform2f(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform2f(
		ReadI32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]),
		ReadF32(fp, expr.Inputs[2]))
}

func opGlUniform3f(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform3f(
		ReadI32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]),
		ReadF32(fp, expr.Inputs[2]),
		ReadF32(fp, expr.Inputs[3]))
}

func opGlUniform4f(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform4f(
		ReadI32(fp, expr.Inputs[0]),
		ReadF32(fp, expr.Inputs[1]),
		ReadF32(fp, expr.Inputs[2]),
		ReadF32(fp, expr.Inputs[3]),
		ReadF32(fp, expr.Inputs[4]))
}

func opGlUniform1i(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform1i(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]))
}

func opGlUniform2i(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform2i(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]))
}

func opGlUniform3i(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform3i(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]))
}

func opGlUniform4i(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform4i(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]),
		ReadI32(fp, expr.Inputs[4]))
}

func opGlUniform1fv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform1fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], TYPE_F32))
}

func opGlUniform2fv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform2fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], TYPE_F32))
}

func opGlUniform3fv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform3fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], TYPE_F32))
}

func opGlUniform4fv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform4fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], TYPE_F32))
}

func opGlUniform1iv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform1iv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], TYPE_I32))
}

func opGlUniform2iv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform2iv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], TYPE_I32))
}

func opGlUniform3iv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform3iv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], TYPE_I32))
}

func opGlUniform4iv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform4iv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], TYPE_I32))
}

func opGlUniformMatrix2fv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniformMatrix2fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadData(fp, expr.Inputs[3], TYPE_F32))
}

func opGlUniformMatrix3fv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniformMatrix3fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadData(fp, expr.Inputs[3], TYPE_F32))
}

func opGlUniformMatrix4fv(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniformMatrix4fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadData(fp, expr.Inputs[3], TYPE_F32))
}

func opGlUniformV4F(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniform4fv(
		ReadI32(fp, expr.Inputs[0]),
		1,
		ReadData(fp, expr.Inputs[1], -1))
}

func opGlUniformM44F(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniformMatrix4fv(
		ReadI32(fp, expr.Inputs[0]),
		1,
		ReadBool(fp, expr.Inputs[1]),
		ReadData(fp, expr.Inputs[2], -1))
}

func opGlUniformM44FV(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglUniformMatrix4fv(
		ReadI32(fp, expr.Inputs[0]),
		ReadI32(fp, expr.Inputs[1]),
		ReadBool(fp, expr.Inputs[2]),
		ReadData(fp, expr.Inputs[3], -1))
}

func opGlVertexAttribPointer(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglVertexAttribPointer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])),
		ReadBool(fp, expr.Inputs[3]),
		ReadI32(fp, expr.Inputs[4]), 0)
}

func opGlVertexAttribPointerI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglVertexAttribPointer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		ReadI32(fp, expr.Inputs[1]),
		uint32(ReadI32(fp, expr.Inputs[2])),
		ReadBool(fp, expr.Inputs[3]),
		ReadI32(fp, expr.Inputs[4]),
		ReadI32(fp, expr.Inputs[5]))
}

func opGlClearBufferI(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opGlClearBufferUI(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opGlClearBufferF(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

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

func opGlBindRenderbuffer(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBindRenderbuffer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlDeleteRenderbuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteRenderbuffers(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
}

func opGlGenRenderbuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenRenderbuffers(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

func opGlRenderbufferStorage(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglRenderbufferStorage(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		ReadI32(fp, expr.Inputs[2]),
		ReadI32(fp, expr.Inputs[3]))
}

func opGlBindFramebuffer(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBindFramebuffer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])))
}

func opGlDeleteFramebuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteFramebuffers(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
}

func opGlGenFramebuffers(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenFramebuffers(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

func opGlCheckFramebufferStatus(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0 := int32(cxglCheckFramebufferStatus(uint32(ReadI32(fp, expr.Inputs[0]))))
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opGlFramebufferTexture2D(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglFramebufferTexture2D(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])),
		uint32(ReadI32(fp, expr.Inputs[3])),
		ReadI32(fp, expr.Inputs[4]))
}

func opGlFramebufferRenderbuffer(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglFramebufferRenderbuffer(
		uint32(ReadI32(fp, expr.Inputs[0])),
		uint32(ReadI32(fp, expr.Inputs[1])),
		uint32(ReadI32(fp, expr.Inputs[2])),
		uint32(ReadI32(fp, expr.Inputs[3])))
}

func opGlGenerateMipmap(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglGenerateMipmap(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlBindVertexArray(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV0 := uint32(ReadI32(fp, expr.Inputs[0]))
	if runtime.GOOS == "darwin" {
		cxglBindVertexArrayAPPLE(inpV0)
	} else {
		cxglBindVertexArray(inpV0)
	}
}

func opGlBindVertexArrayCore(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	cxglBindVertexArray(uint32(ReadI32(fp, expr.Inputs[0])))
}

func opGlDeleteVertexArrays(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV0 := ReadI32(fp, expr.Inputs[0])
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	if runtime.GOOS == "darwin" {
		cxglDeleteVertexArraysAPPLE(inpV0, &inpV1) // will panic if inp0 > 1
	} else {
		cxglDeleteVertexArrays(inpV0, &inpV1) // will panic if inp0 > 1
	}
}

func opGlDeleteVertexArraysCore(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglDeleteVertexArrays(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
}

func opGlGenVertexArrays(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV0 := ReadI32(fp, expr.Inputs[0])
	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	if runtime.GOOS == "darwin" {
		cxglGenVertexArraysAPPLE(inpV0, &inpV1) // will panic if inp0 > 1
	} else {
		cxglGenVertexArrays(inpV0, &inpV1) // will panic if inp0 > 1
	}
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}

func opGlGenVertexArraysCore(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inpV1 := uint32(ReadI32(fp, expr.Inputs[1]))
	cxglGenVertexArrays(ReadI32(fp, expr.Inputs[0]), &inpV1) // will panic if inp0 > 1
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(inpV1))
}
