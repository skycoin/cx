// +build cxfx

package cxfx

import (
	"github.com/skycoin/cx/cx/ast"
	"unicode/utf8"

	"github.com/skycoin/gltext"

	"github.com/skycoin/cx/cx/packages/cxos"
)

var fonts map[string]*gltext.Font = make(map[string]*gltext.Font, 0)

func loadTrueType(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue, fixedPipeline bool) {
	if file := cxos.ValidFile(inputs[0].Get_i32(prgrm)); file != nil {
		if theFont, err := gltext.LoadTruetype(file,
			inputs[2].Get_i32(prgrm), rune(inputs[3].Get_i32(prgrm)), rune(inputs[4].Get_i32(prgrm)),
			gltext.Direction(inputs[5].Get_i32(prgrm)), fixedPipeline); err == nil {
			fonts[inputs[1].Get_str(prgrm)] = theFont
		}
	}
}

func opGltextLoadTrueType(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	loadTrueType(prgrm, inputs, outputs, true)
}

func opGltextLoadTrueTypeCore(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	loadTrueType(prgrm, inputs, outputs, false)
}

func opGltextPrintf(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	if err := fonts[inputs[0].Get_str(prgrm)].Printf(inputs[1].Get_f32(prgrm), inputs[2].Get_f32(prgrm), inputs[3].Get_str(prgrm)); err != nil {
		panic(err)
	}
}

func opGltextMetrics(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	width, height := fonts[inputs[0].Get_str(prgrm)].Metrics(inputs[1].Get_str(prgrm))

	outputs[0].Set_i32(prgrm, int32(width))
	outputs[1].Set_i32(prgrm, int32(height))
}

func opGltextTexture(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_i32(prgrm, int32(fonts[inputs[0].Get_str(prgrm)].Texture()))
}

func opGltextNextGlyph(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // refactor
	font := fonts[inputs[0].Get_str(prgrm)]
	str := inputs[1].Get_str(prgrm)
	var index int = int(inputs[2].Get_i32(prgrm))
	var runeValue rune = -1
	var width int = -1
	var x int = 0
	var y int = 0
	var w int = 0
	var h int = 0
	var advance int = 0
	if index < len(str) {
		runeValue, width = utf8.DecodeRuneInString(str[index:])
		g := font.Glyphs()[runeValue-font.Low()]
		x = g.X
		y = g.Y
		w = g.Width
		h = g.Height
		advance = g.Advance
	}

	outputs[0].Set_i32(prgrm, int32(runeValue-font.Low()))
	outputs[1].Set_i32(prgrm, int32(width))
	outputs[2].Set_i32(prgrm, int32(x))
	outputs[3].Set_i32(prgrm, int32(y))
	outputs[4].Set_i32(prgrm, int32(w))
	outputs[5].Set_i32(prgrm, int32(h))
	outputs[6].Set_i32(prgrm, int32(advance))
}

func opGltextGlyphBounds(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	font := fonts[inputs[0].Get_str(prgrm)]
	var maxGlyphWidth, maxGlyphHeight int = font.GlyphBounds()
	outputs[0].Set_i32(prgrm, int32(maxGlyphWidth))
	outputs[1].Set_i32(prgrm, int32(maxGlyphHeight))
}

func opGltextGlyphMetrics(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // refactor
	width, height := fonts[inputs[0].Get_str(prgrm)].GlyphMetrics(uint32(inputs[1].Get_i32(prgrm)))

	outputs[0].Set_i32(prgrm, int32(width))
	outputs[1].Set_i32(prgrm, int32(height))
}

func opGltextGlyphInfo(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // refactor
	font := fonts[inputs[0].Get_str(prgrm)]
	glyph := inputs[1].Get_i32(prgrm)
	var x int = 0
	var y int = 0
	var w int = 0
	var h int = 0
	var advance int = 0
	g := font.Glyphs()[glyph]
	x = g.X
	y = g.Y
	w = g.Width
	h = g.Height
	advance = g.Advance

	outputs[0].Set_i32(prgrm, int32(x))
	outputs[1].Set_i32(prgrm, int32(y))
	outputs[2].Set_i32(prgrm, int32(w))
	outputs[3].Set_i32(prgrm, int32(h))
	outputs[4].Set_i32(prgrm, int32(advance))
}
