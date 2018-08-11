package resizer

import (
	"io"
	"os"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
	"github.com/techfront/core/src/lib/media"
)

var uploadsDir string

func Setup(config map[string]string) {
	uploadsDir = config["uploads_dir"]
}

func Resize(path string, format string) (p string, err error) {
	switch format {
	case "thumbnail":
		p, err = resizeThumbnail(path)
	}

	return
}

func resizeThumbnail(path string) (string, error) {
	config := media.MediaConfig

	dirName := filepath.Dir(path)
	baseName := filepath.Base(path)
	extName := filepath.Ext(path)
	fileName := baseName[0:len(baseName)-len(extName)]

	width := int(config.ThumbnailsSize.Width)
	height := int(config.ThumbnailsSize.Height)

	r, err := os.Open(uploadsDir + path)
	if err != nil {
		return "", err
	}

	target := dirName + "/" + fileName + "_thumbnail" + config.ThumbnailsFormat

	w, err := os.Create(uploadsDir + target)
	if err != nil {
		return "", err
	}

	img, err := decodeImage(r, extName)
	if err != nil {
		return "", err
	}
	r.Close()

	var res image.Image

	if extName == ".png" {
		im := image.NewRGBA(img.Bounds())
		draw.Draw(im, im.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
		draw.Draw(im, im.Bounds(), img, img.Bounds().Min, draw.Over)
		if img.Bounds().Max.X >= width {
			res = resizeImage(im, uint(width), uint(height), resize.Bilinear)
		} else {
			res = resizeImage(im, uint(img.Bounds().Max.X), uint(img.Bounds().Max.Y), resize.Bilinear)
		}
	} else {
		if img.Bounds().Max.X >= width {
			res = resizeImage(img, uint(width), uint(height), resize.Bilinear)
		} else {
			res = resizeImage(img, uint(img.Bounds().Max.X), uint(img.Bounds().Max.Y), resize.Bilinear)
		}
	}

	if err := encodeImage(w, res, config.ThumbnailsFormat); err != nil {
		return "", err
	}
	w.Close()

	return target, nil
}

func resizeImage(img image.Image, width, height uint, interpolation resize.InterpolationFunction) image.Image {
	return resize.Resize(width, height, img, interpolation)
}

func decodeImage(r io.Reader, format string) (m image.Image, err error) {
	switch format {
	case ".jpg", ".jpeg":
		m, err = jpeg.Decode(r)
	case ".png":
		m, err = png.Decode(r)
	case ".gif":
		m, err = gif.Decode(r)
	case ".webp":
		m, err = webp.Decode(r)
	}

	return
}

func encodeImage(w io.Writer, m image.Image, format string) error {
	var err error

	switch format {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(w, m, &jpeg.Options{Quality: 95})
	case ".png":
		err = png.Encode(w, m)
	case ".gif":
		err = gif.Encode(w, m, &gif.Options{})
	}

	return err
}
