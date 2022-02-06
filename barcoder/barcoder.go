package barcoder

// Note: at this time, if the label is longer than the width of the barcode,
// the last character will be slightly cut-off(most likely due to the used font) and the barcode will not be centered.

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
)

func padImage(src image.Image, horizontal int, vertical int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx()+horizontal+horizontal, src.Bounds().Dy()+vertical+vertical))
	draw.Draw(img, image.Rect(horizontal, vertical, src.Bounds().Dx()+horizontal, src.Bounds().Dy()+vertical), src, src.Bounds().Min, draw.Over)
	return img
}

func subtitleImage(bc image.Image, label string, fontFace font.Face, topMargin int) image.Image {
	fontColor := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// Get the bounds of the string
	bounds, _ := font.BoundString(fontFace, label)

	widthTxt := int((bounds.Max.X - bounds.Min.X) / 64)
	heightTxt := int((bounds.Max.Y - bounds.Min.Y) / 64)

	// calc width and height
	width := widthTxt
	if bc.Bounds().Dx() > width {
		width = bc.Bounds().Dx()
	}
	height := heightTxt + bc.Bounds().Dy() + topMargin

	// create result img
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// draw the barcode
	draw.Draw(img, image.Rect(0, 0, bc.Bounds().Dx(), bc.Bounds().Dy()), bc, bc.Bounds().Min, draw.Over)

	// TextPt
	offsetY := bc.Bounds().Dy() + topMargin - int(bounds.Min.Y/64)
	offsetX := (width - widthTxt) / 2

	point := fixed.Point26_6{
		X: fixed.Int26_6(offsetX * 64),
		Y: fixed.Int26_6(offsetY * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fontColor),
		Face: fontFace,
		Dot:  point,
	}
	d.DrawString(label)
	return img
}

func EncodeUuid(uuid string, labels ...string) (image.Image, error) {
	bc, err := code128.Encode(uuid)
	if err != nil {
		return nil, err
	}

	scaled, err := barcode.Scale(bc, bc.Bounds().Dx(), 100)
	if err != nil {
		return nil, err
	}

	ttFont, err := truetype.Parse(gomono.TTF)
	if err != nil {
		return nil, err
	}

	faceOpts := &truetype.Options{
		DPI:  72,
		Size: 16,
	}
	fontFace := truetype.NewFace(ttFont, faceOpts)

	textMargin := 10
	var img image.Image
	img = scaled
	labels = append([]string{uuid}, labels...)
	for _, label := range labels {
		img = subtitleImage(img, label, fontFace, textMargin)
	}

	padding := 10
	padded := padImage(img, padding, padding)

	// set the background to white since the barcode uses it and we would end up with text
	// and padding with transparent background.
	final := image.NewRGBA(image.Rect(0, 0, padded.Bounds().Dx(), padded.Bounds().Dy()))
	draw.Draw(final, img.Bounds(), &image.Uniform{C: color.White}, image.ZP, draw.Src)

	return final, nil
}
