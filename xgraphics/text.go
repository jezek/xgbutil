package xgraphics

import (
	"image"
	"image/color"
	"io"
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Text takes an image and, using the freetype package, writes text in the
// position specified on to the image. A color.Color, a font size and a font
// must also be specified.
// Finally, the (x, y) coordinate advanced by the text extents is returned.
//
// Note that the ParseFont helper function can be used to get a *truetype.Font
// value without having to import freetype-go directly.
//
// If you need more control over the 'context' used to draw text (like the DPI),
// then you'll need to ignore this convenience method and use your own.
func (im *Image) Text(position fixed.Point26_6, clr color.Color, fontFace font.Face, text string) fixed.Point26_6 {
	// Create a solid color image
	textClr := image.NewUniform(clr)

	fontMetrics := fontFace.Metrics()
	drawer := font.Drawer{
		Dst:  im,
		Src:  textClr,
		Face: fontFace,
		Dot:  position.Add(fixed.Point26_6{X: 0, Y: fontMetrics.Height - fixed.I(fontMetrics.CaretSlope.Y*2)}),
	}
	drawer.DrawString(text)

	return drawer.Dot
}

// ParseFont reads a font file and creates a freetype.Font type
func ParseFont(fontReader io.Reader) (*opentype.Font, error) {
	var otf *opentype.Font
	var err error
	if readerAt, is := fontReader.(io.ReaderAt); is {
		otf, err = opentype.ParseReaderAt(readerAt)
	} else {
		fontBytes, err := ioutil.ReadAll(fontReader)
		if err != nil {
			return nil, err
		}
		otf, err = opentype.Parse(fontBytes)
	}

	if err != nil {
		return nil, err
	}
	return otf, nil
}

// MustFont panics if err is not nil or if the font is nil.
func MustFont(font *opentype.Font, err error) *opentype.Font {
	if err != nil {
		panic(err)
	}
	if font == nil {
		panic("font is nil")
	}
	return font
}
