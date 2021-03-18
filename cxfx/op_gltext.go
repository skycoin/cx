// +build cxfx

package cxfx

import (
	"unicode/utf8"

	"github.com/skycoin/gltext"

	. "github.com/skycoin/cx/cx"
	cxos "github.com/skycoin/cx/cxos"
)

var fonts map[string]*gltext.Font = make(map[string]*gltext.Font, 0)

func loadTrueType(expr *CXExpression, fp int, fixedPipeline bool) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]
	if file := cxos.ValidFile(ReadI32(fp, inp1)); file != nil {
		if theFont, err := gltext.LoadTruetype(file, ReadI32(fp, inp3), rune(ReadI32(fp, inp4)), rune(ReadI32(fp, inp5)), gltext.Direction(ReadI32(fp, inp6)), fixedPipeline); err == nil {
			fonts[ReadStr(fp, inp2)] = theFont
		}
	}
}

func opGltextLoadTrueType(expr *CXExpression, fp int) {
	loadTrueType(expr, fp, true)
}

func opGltextLoadTrueTypeCore(expr *CXExpression, fp int) {
	loadTrueType(expr, fp, false)
}

func opGltextPrintf(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]

	if err := fonts[ReadStr(fp, inp1)].Printf(ReadF32(fp, inp2), ReadF32(fp, inp3), ReadStr(fp, inp4)); err != nil {
		panic(err)
	}
}

func opGltextMetrics(expr *CXExpression, fp int) {
	inp1, inp2, out1, out2 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0], expr.Outputs[1]

	width, height := fonts[ReadStr(fp, inp1)].Metrics(ReadStr(fp, inp2))

	WriteI32(GetFinalOffset(fp, out1), int32(width))
	WriteI32(GetFinalOffset(fp, out2), int32(height))
}

func opGltextTexture(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	WriteI32(GetFinalOffset(fp, out1), int32(fonts[ReadStr(fp, inp1)].Texture()))
}

func opGltextNextGlyph(expr *CXExpression, fp int) { // refactor
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	out1, out2, out3, out4, out5, out6, out7 := expr.Outputs[0], expr.Outputs[1], expr.Outputs[2], expr.Outputs[3], expr.Outputs[4], expr.Outputs[5], expr.Outputs[6]
	font := fonts[ReadStr(fp, inp1)]
	str := ReadStr(fp, inp2)
	var index int = int(ReadI32(fp, inp3))
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

	WriteI32(GetFinalOffset(fp, out1), int32(runeValue-font.Low()))
	WriteI32(GetFinalOffset(fp, out2), int32(width))
	WriteI32(GetFinalOffset(fp, out3), int32(x))
	WriteI32(GetFinalOffset(fp, out4), int32(y))
	WriteI32(GetFinalOffset(fp, out5), int32(w))
	WriteI32(GetFinalOffset(fp, out6), int32(h))
	WriteI32(GetFinalOffset(fp, out7), int32(advance))
}

func opGltextGlyphBounds(expr *CXExpression, fp int) {
	inp1, out1, out2 := expr.Inputs[0], expr.Outputs[0], expr.Outputs[1]
	font := fonts[ReadStr(fp, inp1)]
	var maxGlyphWidth, maxGlyphHeight int = font.GlyphBounds()
	WriteI32(GetFinalOffset(fp, out1), int32(maxGlyphWidth))
	WriteI32(GetFinalOffset(fp, out2), int32(maxGlyphHeight))
}

func opGltextGlyphMetrics(expr *CXExpression, fp int) { // refactor
	inp1, inp2, out1, out2 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0], expr.Outputs[1]

	width, height := fonts[ReadStr(fp, inp1)].GlyphMetrics(uint32(ReadI32(fp, inp2)))

	WriteI32(GetFinalOffset(fp, out1), int32(width))
	WriteI32(GetFinalOffset(fp, out2), int32(height))
}

func opGltextGlyphInfo(expr *CXExpression, fp int) { // refactor
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]
	out1, out2, out3, out4, out5 := expr.Outputs[0], expr.Outputs[1], expr.Outputs[2], expr.Outputs[3], expr.Outputs[4]
	font := fonts[ReadStr(fp, inp1)]
	glyph := ReadI32(fp, inp2)
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

	WriteI32(GetFinalOffset(fp, out1), int32(x))
	WriteI32(GetFinalOffset(fp, out2), int32(y))
	WriteI32(GetFinalOffset(fp, out3), int32(w))
	WriteI32(GetFinalOffset(fp, out4), int32(h))
	WriteI32(GetFinalOffset(fp, out5), int32(advance))
}
