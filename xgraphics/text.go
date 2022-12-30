package xgraphics

import (
	"image"
	"image/color"
	"io"
	"io/ioutil"

	"github.com/BurntSushi/freetype-go/freetype"
	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"golang.org/x/image/font"
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

// Extents returns the *correct* max width and height extents of a string
// given a font. See freetype.MeasureString for the deets.
func Extents(font *truetype.Font, fontSize float64, text string) (int, int) {
	c := ftContext(font, fontSize)
	w, h, err := c.MeasureString(text)
	if err != nil {
		return 0, 0
	}
	return int(w / 256), int(h / 256)
}

// Returns the max width and height extents of a string given a font.
// This is calculated by determining the number of pixels in an "em" unit
// for the given font, and multiplying by the number of characters in 'text'.
// Since a particular character may be smaller than one "em" unit, this has
// a tendency to overestimate the extents.
// It is provided because I do not know how to calculate the precise extents
// using freetype-go.
// TODO: This does not currently account for multiple lines. It may never do so.
func TextMaxExtents(font *truetype.Font, fontSize float64,
	text string) (width int, height int) {

	c := ftContext(font, fontSize)
	emSquarePix := int(c.PointToFix32(fontSize) >> 8)
	return len(text) * emSquarePix, emSquarePix
}

// ftContext does the boiler plate to create a freetype context
func ftContext(font *truetype.Font, fontSize float64) *freetype.Context {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(fontSize)

	return c
}

// ParseFont reads a font file and creates a freetype.Font type
func ParseFont(fontReader io.Reader) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadAll(fontReader)
	if err != nil {
		return nil, err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}

// MustFont panics if err is not nil or if the font is nil.
func MustFont(font *truetype.Font, err error) *truetype.Font {
	if err != nil {
		panic(err)
	}
	if font == nil {
		panic("font is nil")
	}
	return font
}
