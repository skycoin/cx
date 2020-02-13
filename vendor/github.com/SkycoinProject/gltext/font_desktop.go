// +build opengl

package gltext

import (
	"fmt"
	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"image"
	"unsafe"
)

// checkGLError returns an opengl error if one exists.
func checkGLError() error {
	errno := gl.GetError()
	if errno == gl.NO_ERROR {
		return nil
	}
	return fmt.Errorf("GL error: %d", errno)
}

func (f *Font) loadGLFont(img *image.RGBA) error {
	ib := img.Bounds()

	// Create the texture itself. It will contain all glyphs.
	// Individual glyph-quads display a subset of this texture.
	gl.GenTextures(1, &f.texture)
	gl.BindTexture(gl.TEXTURE_2D, f.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(ib.Dx()), int32(ib.Dy()), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	if f.fixedPipeline {
		// Create display lists for each glyph.
		f.listbase = gl.GenLists(int32(len(f.config.Glyphs)))
	}

	texWidth := float32(ib.Dx())
	texHeight := float32(ib.Dy())

	for index, glyph := range f.config.Glyphs {
		// Update max glyph bounds.
		if glyph.Width > f.maxGlyphWidth {
			f.maxGlyphWidth = glyph.Width
		}

		if glyph.Height > f.maxGlyphHeight {
			f.maxGlyphHeight = glyph.Height
		}

		if f.fixedPipeline {
			// Quad width/height
			vw := float32(glyph.Width)
			vh := float32(glyph.Height)

			// Texture coordinate offsets.
			tx1 := float32(glyph.X) / texWidth
			ty1 := float32(glyph.Y) / texHeight
			tx2 := (float32(glyph.X) + vw) / texWidth
			ty2 := (float32(glyph.Y) + vh) / texHeight

			// Advance width (or height if we render top-to-bottom)
			adv := float32(glyph.Advance)
			gl.NewList(f.listbase+uint32(index), gl.COMPILE)
			{
				gl.Begin(gl.QUADS)
				{
					gl.TexCoord2f(tx1, ty2)
					gl.Vertex2f(0, 0)
					gl.TexCoord2f(tx2, ty2)
					gl.Vertex2f(vw, 0)
					gl.TexCoord2f(tx2, ty1)
					gl.Vertex2f(vw, vh)
					gl.TexCoord2f(tx1, ty1)
					gl.Vertex2f(0, vh)
				}
				gl.End()

				switch f.config.Dir {
				case LeftToRight:
					gl.Translatef(adv, 0, 0)
				case RightToLeft:
					gl.Translatef(-adv, 0, 0)
				case TopToBottom:
					gl.Translatef(0, -adv, 0)
				}
			}
			gl.EndList()
		}
	}

	return checkGLError()
}

func (f *Font) releaseGLFont() {
	gl.DeleteTextures(1, &f.texture)
	gl.DeleteLists(f.listbase, int32(len(f.config.Glyphs)))
}

func (f *Font) printfGLFont(x, y float32, fs string, argv ...interface{}) error {
	indices := []rune(fmt.Sprintf(fs, argv...))

	if len(indices) == 0 {
		return nil
	}

	// Runes form display list indices.
	// For this purpose, they need to be offset by -FontConfig.Low
	low := f.config.Low
	for i := range indices {
		indices[i] -= low
	}

	var vp [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &vp[0])

	gl.PushAttrib(gl.TRANSFORM_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(float64(vp[0]), float64(vp[2]), float64(vp[1]), float64(vp[3]), 0, 1)
	gl.PopAttrib()

	gl.PushAttrib(gl.LIST_BIT | gl.CURRENT_BIT | gl.ENABLE_BIT | gl.TRANSFORM_BIT)
	{
		gl.MatrixMode(gl.MODELVIEW)
		gl.Disable(gl.LIGHTING)
		gl.Disable(gl.DEPTH_TEST)
		gl.Enable(gl.BLEND)
		gl.Enable(gl.TEXTURE_2D)

		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
		gl.BindTexture(gl.TEXTURE_2D, f.texture)
		gl.ListBase(f.listbase)

		var mv [16]float32
		gl.GetFloatv(gl.MODELVIEW_MATRIX, &mv[0])

		gl.PushMatrix()
		{
			gl.LoadIdentity()

			mgw := float32(f.maxGlyphWidth)
			mgh := float32(f.maxGlyphHeight)

			switch f.config.Dir {
			case LeftToRight, TopToBottom:
				gl.Translatef(x, float32(vp[3])-y-mgh, 0)
			case RightToLeft:
				gl.Translatef(x-mgw, float32(vp[3])-y-mgh, 0)
			}

			gl.MultMatrixf(&mv[0])
			gl.CallLists(int32(len(indices)), gl.UNSIGNED_INT, unsafe.Pointer(&indices[0]))
		}
		gl.PopMatrix()
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
	gl.PopAttrib()

	gl.PushAttrib(gl.TRANSFORM_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.PopAttrib()
	return checkGLError()
}
