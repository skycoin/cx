// +build cxfx

package cxfx

import (
	"bufio"
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cx/util"
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
	//"bytes"
	"unsafe"
)

type Texture struct {
	path   string
	width  int32
	height int32
	level  uint32
	pixels []float32
}

func Slice_ui8_ToPtr(value []uint8) unsafe.Pointer {
	count := len(value)
	if count == 0 {
		return unsafe.Pointer(nil)
	}
	return unsafe.Pointer(&value[0])
}

func Slice_ui32_ToPtr(value []uint32) unsafe.Pointer {
	count := len(value)
	if count == 0 {
		return unsafe.Pointer(nil)
	}
	return unsafe.Pointer(&value[0])
}

func Slice_i32_ToPtr(value []int32) unsafe.Pointer {
	count := len(value)
	if count == 0 {
		return unsafe.Pointer(nil)
	}
	return unsafe.Pointer(&value[0])
}

func Slice_f32_ToPtr(value []float32) unsafe.Pointer {
	count := len(value)
	if count == 0 {
		return unsafe.Pointer(nil)
	}
	return unsafe.Pointer(&value[0])
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

func unpack(file *os.File, width types.Pointer, line []byte) bool {
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
	for i := types.Pointer(0); i < 4; i++ {
		for j := types.Pointer(0); j < width; {
			file.Read(b[:])
			count := types.Pointer(b[0])
			if count > 128 {
				count &= 127
				file.Read(b[:])
				var value int = int(b[0])
				for c := types.Pointer(0); c < count; c++ {
					line[j+c+i] = byte(value)
				}
			} else {
				for c := types.Pointer(0); c < count; c++ {
					offset := j + c + i
					file.Read(line[offset : offset+1])
				}
			}
		}
	}
	return true
}

func unpack_(file *os.File, width types.Pointer, line []byte) bool {
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

func decodeHdr(file *os.File) (data []byte, i32Width int32, i32Height int32) {
	data = nil
	i32Width = 0
	i32Height = 0

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

	var iwidth int
	var iheight int
	if n, err := fmt.Fscanf(file, "-Y %d +X %d\n", &iwidth, &iheight); n != 2 || err != nil {
		fmt.Printf("Failed to scan width and height : err '%s'\n", err)
		return
	}

	i32Width = int32(iwidth)
	i32Height = int32(iheight)

	width := types.Cast_int_to_ptr(iwidth)
	height := types.Cast_int_to_ptr(iheight)

	//var colors []float32 = make([]float32, width*height*3)
	var line []byte = make([]byte, width*4)
	data = make([]byte, width*height*3*4)

	for y := types.Pointer(0); y < height; y++ {
		if unpack(file, width, line) == false {
			fmt.Printf("Failed to unpack line %d\n", y)
			return
		}

		yoffset := y /*(height - y - 1)*/ * width * 3 * 4
		for x := types.Pointer(0); x < width; x++ {
			loffset := x * 4
			exponent := math.Pow(2.0, float64(int(line[loffset+3])-128))
			xoffset := yoffset + x*3*4
			r := float32(exponent * float64(line[loffset]) / 256.0)
			g := float32(exponent * float64(line[loffset+1]) / 256.0)
			b := float32(exponent * float64(line[loffset+2]) / 256.0)

			types.Write_f32(data, xoffset, r)
			types.Write_f32(data, xoffset+4, g)
			types.Write_f32(data, xoffset+8, b)
		}
	}
	return
}

func uploadTexture(path string, target uint32, level uint32, cpuCopy bool) {
	file, err := util.CXOpenFile(path)
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
func opGlNewTexture(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var texture uint32
	cxglEnable(cxglTEXTURE_2D)
	cxglGenTextures(1, &texture)
	cxglBindTexture(cxglTEXTURE_2D, texture)
	cxglTexParameteri(cxglTEXTURE_2D, cxglTEXTURE_MIN_FILTER, cxglNEAREST)
	cxglTexParameteri(cxglTEXTURE_2D, cxglTEXTURE_MAG_FILTER, cxglNEAREST)
	cxglTexParameteri(cxglTEXTURE_2D, cxglTEXTURE_WRAP_S, cxglCLAMP_TO_EDGE)
	cxglTexParameteri(cxglTEXTURE_2D, cxglTEXTURE_WRAP_T, cxglCLAMP_TO_EDGE)

	uploadTexture(inputs[0].Get_str(), cxglTEXTURE_2D, 0, false)

	outputs[0].Set_i32(int32(texture))
}

func opGlNewTextureCube(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
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
	var pattern string = inputs[0].Get_str()
	var extension string = inputs[1].Get_str()
	for i := 0; i < 6; i++ {
		uploadTexture(fmt.Sprintf("%s%s%s", pattern, faces[i], extension), uint32(cxglTEXTURE_CUBE_MAP_POSITIVE_X+i), 0, false)
	}
	outputs[0].Set_i32(int32(texture))
}

func opCxReleaseTexture(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	textures[inputs[0].Get_str()] = Texture{}
}

func opCxTextureGetPixel(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var r float32
	var g float32
	var b float32
	var a float32

	var x = inputs[1].Get_i32()
	var y = inputs[2].Get_i32()

	if texture, ok := textures[inputs[0].Get_str()]; ok {
		var yoffset = y * texture.width * 4
		var xoffset = yoffset + x*4
		pixels := texture.pixels
		r = pixels[xoffset]
		g = pixels[xoffset+1]
		b = pixels[xoffset+2]
		a = pixels[xoffset+3]
	}
	outputs[0].Set_f32(r)
	outputs[1].Set_f32(g)
	outputs[2].Set_f32(b)
	outputs[3].Set_f32(a)
}

func opGlUploadImageToTexture(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	uploadTexture(inputs[0].Get_str(), uint32(inputs[1].Get_i32()), uint32(inputs[2].Get_i32()), inputs[3].Get_bool())
}

func opGlNewGIF(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	path := inputs[0].Get_str()

	file, err := util.CXOpenFile(path)
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

	outputs[0].Set_i32(int32(len(gif.Image)))
	outputs[1].Set_i32(int32(gif.LoopCount))
	outputs[2].Set_i32(int32(gif.Config.Width))
	outputs[3].Set_i32(int32(gif.Config.Height))
}

func opGlFreeGIF(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	gifs[inputs[0].Get_str()] = nil
}

func opGlGIFFrameToTexture(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	path := inputs[0].Get_str()
	frame := inputs[1].Get_i32()
	texture := inputs[2].Get_i32()

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

	outputs[0].Set_i32(delay)
	outputs[1].Set_i32(disposal)
}

func opGlAppend(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputSlicePointer := outputs[0].Offset
	outputSliceOffset := types.Read_ptr(prgrm.Memory, outputSlicePointer)

	inputSliceOffset := ast.GetSliceOffset(inputs[0].FramePointer, inputs[0].Arg)
	var inputSliceLen types.Pointer
	if inputSliceOffset != 0 {
		inputSliceLen = ast.GetSliceLen(inputSliceOffset)
	}

	obj := inputs[1].Get_bytes()

	objLen := types.Cast_int_to_ptr(len(obj))
	outputSliceOffset = ast.SliceResizeEx(outputSliceOffset, inputSliceLen+objLen, 1)
	ast.SliceCopyEx(outputSliceOffset, inputSliceOffset, inputSliceLen+objLen, 1)
	ast.SliceAppendWriteByte(outputSliceOffset, obj, inputSliceLen)
	outputs[0].Set_ptr(outputSliceOffset)
}

// gl_1_0
func opGlCullFace(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglCullFace(uint32(inputs[0].Get_i32()))
}

func opGlFrontFace(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglFrontFace(uint32(inputs[0].Get_i32()))
}

func opGlHint(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglHint(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlScissor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglScissor(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_i32(),
		inputs[3].Get_i32())
}

func opGlTexParameteri(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglTexParameteri(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		inputs[2].Get_i32())
}

func opGlTexImage2D(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglTexImage2D(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		inputs[2].Get_i32(),
		inputs[3].Get_i32(),
		inputs[4].Get_i32(),
		inputs[5].Get_i32(),
		uint32(inputs[6].Get_i32()),
		uint32(inputs[7].Get_i32()),
		inputs[8].GetSlice_bytes())
}

func opGlClear(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglClear(uint32(inputs[0].Get_i32()))
}

func opGlClearColor(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglClearColor(
		inputs[0].Get_f32(),
		inputs[1].Get_f32(),
		inputs[2].Get_f32(),
		inputs[3].Get_f32())
}

func opGlClearStencil(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglClearStencil(inputs[0].Get_i32())
}

func opGlClearDepth(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglClearDepth(inputs[0].Get_f64())
}

func opGlStencilMask(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglStencilMask(uint32(inputs[0].Get_i32()))
}

func opGlColorMask(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglColorMask(
		inputs[0].Get_bool(),
		inputs[1].Get_bool(),
		inputs[2].Get_bool(),
		inputs[3].Get_bool())
}

func opGlDepthMask(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDepthMask(inputs[0].Get_bool())
}

func opGlDisable(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDisable(uint32(inputs[0].Get_i32()))
}

func opGlEnable(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglEnable(uint32(inputs[0].Get_i32()))
}

func opGlBlendFunc(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBlendFunc(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlStencilFunc(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglStencilFunc(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		uint32(inputs[2].Get_i32()))
}

func opGlStencilOp(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglStencilOp(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		uint32(inputs[2].Get_i32()))
}

func opGlDepthFunc(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDepthFunc(uint32(inputs[0].Get_i32()))
}

func opGlGetError(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_i32(int32(cxglGetError()))
}

func opGlGetTexLevelParameteriv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	var outValue int32 = 0
	cxglGetTexLevelParameteriv(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		uint32(inputs[2].Get_i32()),
		&outValue)
	outputs[0].Set_i32(outValue)
}

func opGlDepthRange(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDepthRange(
		inputs[0].Get_f64(),
		inputs[1].Get_f64())
}

func opGlViewport(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglViewport(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_i32(),
		inputs[3].Get_i32())
}

// gl_1_1
func opGlDrawArrays(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDrawArrays(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		inputs[2].Get_i32())
}

func opGlDrawElements(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDrawElements(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		uint32(inputs[2].Get_i32()),
		nil)
}

func opGlBindTexture(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBindTexture(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlDeleteTextures(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglDeleteTextures(inputs[0].Get_i32(), &inpV1) // will panic if inp0 > 1
}

func opGlGenTextures(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglGenTextures(inputs[0].Get_i32(), &inpV1) // will panic if inp0 > 1
	outputs[0].Set_i32(int32(inpV1))
}

// gl_1_3
func opGlActiveTexture(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglActiveTexture(uint32(inputs[0].Get_i32()))
}

// gl_1_4
func opGlBlendFuncSeparate(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBlendFuncSeparate(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		uint32(inputs[2].Get_i32()),
		uint32(inputs[3].Get_i32()))
}

// gl_1_5
func opGlBindBuffer(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBindBuffer(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlDeleteBuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglDeleteBuffers(
		inputs[0].Get_i32(),
		&inpV1) // will panic if inp0 > 1
}

func opGlGenBuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglGenBuffers(
		inputs[0].Get_i32(),
		&inpV1) // will panic if inp0 > 1
	outputs[0].Set_i32(int32(inpV1))
}

func opGlBufferData(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBufferData(
		uint32(inputs[0].Get_i32()),
		int(inputs[1].Get_i32()),
		inputs[2].GetSlice_bytes(),
		uint32(inputs[3].Get_i32()))
}

func opGlBufferSubData(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBufferSubData(
		uint32(inputs[0].Get_i32()),
		int(inputs[1].Get_i32()),
		int(inputs[2].Get_i32()),
		inputs[3].GetSlice_bytes())
}

func opGlDrawBuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDrawBuffers(
		inputs[0].Get_i32(),
		inputs[1].GetSlice_bytes())
}

func opGlStencilOpSeparate(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglStencilOpSeparate(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		uint32(inputs[2].Get_i32()),
		uint32(inputs[3].Get_i32()))
}

func opGlStencilFuncSeparate(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglStencilFuncSeparate(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		inputs[2].Get_i32(),
		uint32(inputs[3].Get_i32()))
}

func opGlStencilMaskSeparate(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglStencilMaskSeparate(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlAttachShader(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglAttachShader(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlBindAttribLocation(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBindAttribLocation(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		inputs[2].Get_str())
}

func opGlCompileShader(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	shader := uint32(inputs[0].Get_i32())
	cxglCompileShader(shader)
}

func opGlCreateProgram(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_i32(int32(cxglCreateProgram()))
}

func opGlCreateShader(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(cxglCreateShader(uint32(inputs[0].Get_i32())))
	outputs[0].Set_i32(outV0)
}

func opGlDeleteProgram(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDeleteShader(uint32(inputs[0].Get_i32()))
}

func opGlDeleteShader(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDeleteShader(uint32(inputs[0].Get_i32()))
}

func opGlDetachShader(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglDetachShader(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlEnableVertexAttribArray(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglEnableVertexAttribArray(uint32(inputs[0].Get_i32()))
}

func opGlGetAttribLocation(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := cxglGetAttribLocation(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_str())
	outputs[0].Set_i32(outV0)
}

func opGlGetProgramiv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := cxglGetProgramiv(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
	outputs[0].Set_i32(outV0)
}

func opGlGetProgramInfoLog(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := cxglGetProgramInfoLog(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32())
	outputs[0].Set_str(outV0)
}

func opGlGetShaderiv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := cxglGetShaderiv(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
	outputs[0].Set_i32(outV0)
}

func opGlGetShaderInfoLog(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := cxglGetShaderInfoLog(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32())
	outputs[0].Set_str(outV0)
}

func opGlGetUniformLocation(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := cxglGetUniformLocation(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_str())
	outputs[0].Set_i32(outV0)
}

func opGlLinkProgram(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	program := uint32(inputs[0].Get_i32())
	cxglLinkProgram(program)
}

func opGlShaderSource(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglShaderSource(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		inputs[2].Get_str())
}

func opGlUseProgram(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUseProgram(uint32(inputs[0].Get_i32()))
}

func opGlUniform1f(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform1f(
		inputs[0].Get_i32(),
		inputs[1].Get_f32())
}

func opGlUniform2f(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform2f(
		inputs[0].Get_i32(),
		inputs[1].Get_f32(),
		inputs[2].Get_f32())
}

func opGlUniform3f(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform3f(
		inputs[0].Get_i32(),
		inputs[1].Get_f32(),
		inputs[2].Get_f32(),
		inputs[3].Get_f32())
}

func opGlUniform4f(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform4f(
		inputs[0].Get_i32(),
		inputs[1].Get_f32(),
		inputs[2].Get_f32(),
		inputs[3].Get_f32(),
		inputs[4].Get_f32())
}

func opGlUniform1i(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform1i(
		inputs[0].Get_i32(),
		inputs[1].Get_i32())
}

func opGlUniform2i(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform2i(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_i32())
}

func opGlUniform3i(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform3i(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_i32(),
		inputs[3].Get_i32())
}

func opGlUniform4i(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform4i(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_i32(),
		inputs[3].Get_i32(),
		inputs[4].Get_i32())
}

func opGlUniform1fv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform1fv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].GetSlice_bytes())
}

func opGlUniform2fv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform2fv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].GetSlice_bytes())
}

func opGlUniform3fv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform3fv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].GetSlice_bytes())
}

func opGlUniform4fv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform4fv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].GetSlice_bytes())
}

func opGlUniform1iv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform1iv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].GetSlice_bytes())
}

func opGlUniform2iv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform2iv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].GetSlice_bytes())
}

func opGlUniform3iv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform3iv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].GetSlice_bytes())
}

func opGlUniform4iv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform4iv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].GetSlice_bytes())
}

func opGlUniformMatrix2fv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniformMatrix2fv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_bool(),
		inputs[3].GetSlice_bytes())
}

func opGlUniformMatrix3fv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniformMatrix3fv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_bool(),
		inputs[3].GetSlice_bytes())
}

func opGlUniformMatrix4fv(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniformMatrix4fv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_bool(),
		inputs[3].GetSlice_bytes())
}

func opGlUniformV4F(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniform4fv(
		inputs[0].Get_i32(),
		1,
		inputs[1].Get_bytes())
}

func opGlUniformM44F(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniformMatrix4fv(
		inputs[0].Get_i32(),
		1,
		inputs[1].Get_bool(),
		inputs[2].Get_bytes())
}

func opGlUniformM44FV(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglUniformMatrix4fv(
		inputs[0].Get_i32(),
		inputs[1].Get_i32(),
		inputs[2].Get_bool(),
		inputs[3].GetSlice_bytes())
}

func opGlVertexAttribPointer(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglVertexAttribPointer(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		uint32(inputs[2].Get_i32()),
		inputs[3].Get_bool(),
		inputs[4].Get_i32(), 0)
}

func opGlVertexAttribPointerI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglVertexAttribPointer(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		uint32(inputs[2].Get_i32()),
		inputs[3].Get_bool(),
		inputs[4].Get_i32(),
		inputs[5].Get_i32())
}

func opGlClearBufferI(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	color := [4]int32{
		inputs[2].Get_i32(),
		inputs[3].Get_i32(),
		inputs[4].Get_i32(),
		inputs[5].Get_i32()}

	cxglClearBufferiv(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		color[:])
}

func opGlClearBufferUI(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	color := [4]uint32{
		inputs[2].Get_ui32(),
		inputs[3].Get_ui32(),
		inputs[4].Get_ui32(),
		inputs[5].Get_ui32()}

	cxglClearBufferuiv(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		color[:])
}

func opGlClearBufferF(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	color := [4]float32{
		inputs[2].Get_f32(),
		inputs[3].Get_f32(),
		inputs[4].Get_f32(),
		inputs[5].Get_f32()}

	cxglClearBufferfv(
		uint32(inputs[0].Get_i32()),
		inputs[1].Get_i32(),
		color[:])
}

func opGlBindRenderbuffer(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBindRenderbuffer(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlDeleteRenderbuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglDeleteRenderbuffers(inputs[0].Get_i32(), &inpV1) // will panic if inp0 > 1
}

func opGlGenRenderbuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglGenRenderbuffers(inputs[0].Get_i32(), &inpV1) // will panic if inp0 > 1
	outputs[0].Set_i32(int32(inpV1))
}

func opGlRenderbufferStorage(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglRenderbufferStorage(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		inputs[2].Get_i32(),
		inputs[3].Get_i32())
}

func opGlBindFramebuffer(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBindFramebuffer(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()))
}

func opGlDeleteFramebuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglDeleteFramebuffers(inputs[0].Get_i32(), &inpV1) // will panic if inp0 > 1
}

func opGlGenFramebuffers(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglGenFramebuffers(inputs[0].Get_i32(), &inpV1) // will panic if inp0 > 1
	outputs[0].Set_i32(int32(inpV1))
}

func opGlCheckFramebufferStatus(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := int32(cxglCheckFramebufferStatus(uint32(inputs[0].Get_i32())))
	outputs[0].Set_i32(outV0)
}

func opGlFramebufferTexture2D(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglFramebufferTexture2D(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		uint32(inputs[2].Get_i32()),
		uint32(inputs[3].Get_i32()),
		inputs[4].Get_i32())
}

func opGlFramebufferRenderbuffer(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglFramebufferRenderbuffer(
		uint32(inputs[0].Get_i32()),
		uint32(inputs[1].Get_i32()),
		uint32(inputs[2].Get_i32()),
		uint32(inputs[3].Get_i32()))
}

func opGlGenerateMipmap(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglGenerateMipmap(uint32(inputs[0].Get_i32()))
}

func opGlBindVertexArray(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := uint32(inputs[0].Get_i32())
	if runtime.GOOS == "darwin" {
		cxglBindVertexArrayAPPLE(inpV0)
	} else {
		cxglBindVertexArray(inpV0)
	}
}

func opGlBindVertexArrayCore(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	cxglBindVertexArray(uint32(inputs[0].Get_i32()))
}

func opGlDeleteVertexArrays(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i32()
	inpV1 := uint32(inputs[1].Get_i32())
	if runtime.GOOS == "darwin" {
		cxglDeleteVertexArraysAPPLE(inpV0, &inpV1) // will panic if inp0 > 1
	} else {
		cxglDeleteVertexArrays(inpV0, &inpV1) // will panic if inp0 > 1
	}
}

func opGlDeleteVertexArraysCore(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglDeleteVertexArrays(inputs[0].Get_i32(), &inpV1) // will panic if inp0 > 1
}

func opGlGenVertexArrays(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_i32()
	inpV1 := uint32(inputs[1].Get_i32())
	if runtime.GOOS == "darwin" {
		cxglGenVertexArraysAPPLE(inpV0, &inpV1) // will panic if inp0 > 1
	} else {
		cxglGenVertexArrays(inpV0, &inpV1) // will panic if inp0 > 1
	}
	outputs[0].Set_i32(int32(inpV1))
}

func opGlGenVertexArraysCore(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV1 := uint32(inputs[1].Get_i32())
	cxglGenVertexArrays(inputs[0].Get_i32(), &inpV1) // will panic if inp0 > 1
	outputs[0].Set_i32(int32(inpV1))
}
