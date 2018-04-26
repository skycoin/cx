package base

import (
	"github.com/go-gl/gltext"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var fonts map[string]*gltext.Font = make(map[string]*gltext.Font, 0)

func gltext_LoadTrueType (font, file, scale, low, high, dir *CXArgument) error {
	if err := checkSixTypes("gltext.LoadTrueType", "str", "str", "i32", "i32", "i32", "i32", font, file, scale, low, high, dir); err == nil {
		var _font string
		var _file string
		var _scale int32
		var _low int32
		var _high int32
		var _dir int32

		encoder.DeserializeRaw(*font.Value, &_font)
		encoder.DeserializeRaw(*file.Value, &_file)
		encoder.DeserializeAtomic(*scale.Value, &_scale)
		encoder.DeserializeAtomic(*low.Value, &_low)
		encoder.DeserializeAtomic(*high.Value, &_high)
		encoder.DeserializeAtomic(*dir.Value, &_dir)

		if theFont, err := gltext.LoadTruetype(openFiles[_file], _scale, rune(_low), rune(_high), gltext.Direction(_dir)); err == nil {
			fonts[_font] = theFont
		} else {
			return err
		}
		
		return nil
	} else {
		return err
	}
}

func gltext_Printf (font, x, y, fs *CXArgument) error {
	if err := checkFourTypes("gltext.Printf", "str", "f32", "f32", "str", font, x, y, fs); err == nil {
		var _font string
		var _x float32
		var _y float32
		var _fs string

		encoder.DeserializeRaw(*font.Value, &_font)
		encoder.DeserializeRaw(*x.Value, &_x)
		encoder.DeserializeRaw(*y.Value, &_y)
		encoder.DeserializeRaw(*fs.Value, &_fs)

		if err := fonts[_font].Printf(_x, _y, _fs); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}
