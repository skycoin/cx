// +build opengles

package gltext

import (
	"fmt"
	"golang.org/x/mobile/gl"
	"image"
)

var glctx gl.Context

func SetGLContext(ctx gl.Context) {
	glctx = ctx
}

// checkGLError returns an opengl error if one exists.
func checkGLError() error {
	errno := glctx.GetError()
	if errno == gl.NO_ERROR {
		return nil
	}
	return fmt.Errorf("GL error: %d", errno)
}

func (f *Font) loadGLFont(img *image.RGBA) error {
	if f.fixedPipeline {
		return fmt.Errorf("loadFont does Not support opengl fixed pipeline on mobile")
	}

	ib := img.Bounds()

	// Create the texture itself. It will contain all glyphs.
	// Individual glyph-quads display a subset of this texture.
	f.texture = glctx.CreateTexture().Value
	glctx.BindTexture(gl.TEXTURE_2D, gl.Texture{f.texture})
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	glctx.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	glctx.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int(ib.Dx()), int(ib.Dy()), /*0,*/
		gl.RGBA, gl.UNSIGNED_BYTE, img.Pix)

	for _, glyph := range f.config.Glyphs {
		// Update max glyph bounds.
		if glyph.Width > f.maxGlyphWidth {
			f.maxGlyphWidth = glyph.Width
		}

		if glyph.Height > f.maxGlyphHeight {
			f.maxGlyphHeight = glyph.Height
		}
	}

	return checkGLError()
}

func (f *Font) releaseGLFont() {
	glctx.DeleteTexture(gl.Texture{f.texture})
}

func (f *Font) printfGLFont(x, y float32, fs string, argv ...interface{}) error {
	return fmt.Errorf("Font.Printf is Not implemented on mobile")
}
