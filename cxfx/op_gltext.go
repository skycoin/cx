// +build cxfx

package cxfx

import (
	"unicode/utf8"

	"github.com/skycoin/gltext"

	. "github.com/skycoin/cx/cx"
	cxos "github.com/skycoin/cx/cxos"
)

var fonts map[string]*gltext.Font = make(map[string]*gltext.Font, 0)

func loadTrueType(inputs []CXValue, outputs []CXValue, fixedPipeline bool) {
	if file := cxos.ValidFile(inputs[0].Get_i32()); file != nil {
		if theFont, err := gltext.LoadTruetype(file,
            inputs[2].Get_i32(), rune(inputs[3].Get_i32()), rune(inputs[4].Get_i32()),
            gltext.Direction(inputs[5].Get_i32()), fixedPipeline); err == nil {
			fonts[inputs[1].Get_str()] = theFont
		}
	}
}

func opGltextLoadTrueType(inputs []CXValue, outputs []CXValue) {
	loadTrueType(inputs, outputs, true)
}

func opGltextLoadTrueTypeCore(inputs []CXValue, outputs []CXValue) {
	loadTrueType(inputs, outputs, false)
}

func opGltextPrintf(inputs []CXValue, outputs []CXValue) {
	if err := fonts[inputs[0].Get_str()].Printf(inputs[1].Get_f32(), inputs[2].Get_f32(), inputs[3].Get_str()); err != nil {
		panic(err)
	}
}

func opGltextMetrics(inputs []CXValue, outputs []CXValue) {
	width, height := fonts[inputs[0].Get_str()].Metrics(inputs[1].Get_str())

	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGltextTexture(inputs []CXValue, outputs []CXValue) {
	outputs[0].Set_i32(int32(fonts[inputs[0].Get_str()].Texture()))
}

func opGltextNextGlyph(inputs []CXValue, outputs []CXValue) { // refactor
	font := fonts[inputs[0].Get_str()]
	str := inputs[1].Get_str()
	var index int = int(inputs[2].Get_i32())
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

	outputs[0].Set_i32(int32(runeValue-font.Low()))
	outputs[1].Set_i32(int32(width))
	outputs[2].Set_i32(int32(x))
	outputs[3].Set_i32(int32(y))
	outputs[4].Set_i32(int32(w))
	outputs[5].Set_i32(int32(h))
	outputs[6].Set_i32(int32(advance))
}

func opGltextGlyphBounds(inputs []CXValue, outputs []CXValue) {
	font := fonts[inputs[0].Get_str()]
	var maxGlyphWidth, maxGlyphHeight int = font.GlyphBounds()
	outputs[0].Set_i32(int32(maxGlyphWidth))
	outputs[1].Set_i32(int32(maxGlyphHeight))
}

func opGltextGlyphMetrics(inputs []CXValue, outputs []CXValue) { // refactor
	width, height := fonts[inputs[0].Get_str()].GlyphMetrics(uint32(inputs[1].Get_i32()))

	outputs[0].Set_i32(int32(width))
	outputs[1].Set_i32(int32(height))
}

func opGltextGlyphInfo(inputs []CXValue, outputs []CXValue) { // refactor
	font := fonts[inputs[0].Get_str()]
	glyph := inputs[1].Get_i32()
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

	outputs[0].Set_i32(int32(x))
	outputs[1].Set_i32(int32(y))
	outputs[2].Set_i32(int32(w))
	outputs[3].Set_i32(int32(h))
	outputs[4].Set_i32(int32(advance))
}
