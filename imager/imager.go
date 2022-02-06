package imager

import (
	"errors"
	"github.com/disintegration/imaging"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
)

// only first existing rule will be used, detected in order as defined in the struct
type Rule struct {
	Resize    *ResizeRule    `json:"resize,omitempty"`
	Fill      *FillRule      `json:"fill,omitempty"`
	Fit       *FitRule       `json:"fit,omitempty"`
	Thumbnail *ThumbnailRule `json:"thumbnail,omitempty"`
}

func (r *Rule) Validate() error {
	var count int

	if r.Resize != nil {
		count++

		if r.Resize.Width < 0 || r.Resize.Height < 0 {
			return errors.New("invalid resize dimensions")
		}
		if r.Resize.Width == 0 && r.Resize.Height == 0 {
			return errors.New("empty resize dimensions")
		}
	}

	if r.Fill != nil {
		count++

		if r.Fill.Width < 0 || r.Fill.Height < 0 {
			return errors.New("invalid fill dimensions")
		}
		if r.Fill.Width == 0 && r.Fill.Height == 0 {
			return errors.New("empty fill dimensions")
		}
	}

	if r.Fit != nil {
		count++

		if r.Fit.Width < 0 || r.Fit.Height < 0 {
			return errors.New("invalid fit dimensions")
		}
		if r.Fit.Width == 0 && r.Fit.Height == 0 {
			return errors.New("empty fit dimensions")
		}
	}

	if r.Thumbnail != nil {
		count++

		if r.Thumbnail.Width < 0 || r.Thumbnail.Height < 0 {
			return errors.New("invalid thumbnail dimensions")
		}
		if r.Thumbnail.Width == 0 && r.Thumbnail.Height == 0 {
			return errors.New("empty thumbnail dimensions")
		}
	}

	if count == 0 {
		return errors.New("rule is empty")
	}

	if count > 1 {
		return errors.New("rule has multiple presets set")
	}

	return nil
}

type ResizeRule struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type FillRule struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Anchor string `json:"anchor"` // Center(c)-default, TopLeft(tl) Top(t), TopRight(tr), Left(l), Right(r), BottomLeft(bl), Bottom(b), BottomRight(br)
}

func (r *FillRule) GetAnchor() imaging.Anchor {
	switch r.Anchor {
	case "c":
		return imaging.Center
	case "tl":
		return imaging.TopLeft
	case "t":
		return imaging.Top
	case "tr":
		return imaging.TopRight
	case "l":
		return imaging.Left
	case "r":
		return imaging.Right
	case "bl":
		return imaging.BottomLeft
	case "b":
		return imaging.Bottom
	case "br":
		return imaging.BottomRight
	default:
		return imaging.Center
	}
}

type FitRule struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ThumbnailRule struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// NOTE: there is no native WebP encoder so if the source is WebP, the result will be PNG, not WebP.
// NOTE: once the function returns, the caller has to call Seek(0,0) on dst in order to read the data.
func Apply(src io.ReadSeeker, dst io.Writer, rules ...Rule) error {
	if src == nil {
		return errors.New("no source provided")
	}
	if dst == nil {
		return errors.New("no destination provided")
	}
	if len(rules) == 0 {
		return errors.New("no rules provided")
	}
	for k := range rules {
		if err := rules[k].Validate(); err != nil {
			return err
		}
	}

	buf := make([]byte, 512)
	read, err := src.Read(buf)
	if err != nil {
		return err
	}

	mime := http.DetectContentType(buf[:read])
	if _, err := src.Seek(0, 0); err != nil {
		return err
	}

	var img image.Image

	switch mime {
	case "image/jpeg":
		img, err = jpeg.Decode(src)
	case "image/gif":
		img, err = gif.Decode(src)
	case "image/png":
		img, err = png.Decode(src)
	case "image/webp":
		img, err = webp.Decode(src)
	default:
		return errors.New("unsupported mime type: " + mime)
	}

	if err != nil {
		return err
	}

	filter := imaging.Lanczos // best quality but slow

	for _, r := range rules {
		if r.Resize != nil {
			img = imaging.Resize(img, r.Resize.Width, r.Resize.Height, filter)
			continue
		}

		if r.Fill != nil {
			img = imaging.Fill(img, r.Fill.Width, r.Fill.Height, r.Fill.GetAnchor(), filter)
			continue
		}

		if r.Fit != nil {
			img = imaging.Fit(img, r.Fit.Width, r.Fit.Height, filter)
			continue
		}

		if r.Thumbnail != nil {
			img = imaging.Thumbnail(img, r.Thumbnail.Width, r.Thumbnail.Height, filter)
			continue
		}

		return errors.New("empty rule provided")
	}

	switch mime {
	case "image/jpeg":
		return jpeg.Encode(dst, img, nil)
	case "image/gif":
		return gif.Encode(dst, img, nil)
	case "image/png":
		return png.Encode(dst, img)
	case "image/webp":
		return png.Encode(dst, img) // there is no native webp encoder https://github.com/golang/go/issues/45121 and we don't want cgo here
	default:
		return errors.New("unsupported mime type: " + mime)
	}
}
