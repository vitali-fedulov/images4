package images4

import (
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

// Open opens and decodes an image file for a given path.
func Open(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, err
}

// resizeByNearest resizes an image to the destination size
// with the nearest neighbour method. It also returns the source
// image size.
func resizeByNearest(
	src image.Image, dstSize image.Point) (
	dst image.RGBA, srcSize image.Point) {
	// Original image size.
	xMax, xMin := src.Bounds().Max.X, src.Bounds().Min.X
	yMax, yMin := src.Bounds().Max.Y, src.Bounds().Min.Y
	srcX := xMax - xMin
	srcY := yMax - yMin
	xScale := float64(srcX) / float64(dstSize.X)
	yScale := float64(srcY) / float64(dstSize.Y)

	// Destination rectangle.
	outRect := image.Rectangle{
		image.Point{0, 0}, image.Point{dstSize.X, dstSize.Y}}
	// Color model of uint8 per color.
	dst = *image.NewRGBA(outRect)
	var (
		r, g, b, a uint32
	)
	for y := 0; y < dstSize.Y; y++ {
		for x := 0; x < dstSize.X; x++ {
			r, g, b, a = src.At(
				int(float64(x)*xScale+float64(xMin)),
				int(float64(y)*yScale+float64(yMin))).RGBA()
			dst.Set(x, y, color.RGBA{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(a >> 8)})
		}
	}
	return dst, image.Point{srcX, srcY}
}

// SaveToPNG encodes and saves image.RGBA to a file.
func SaveToPNG(img *image.RGBA, path string) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		png.Encode(destFile, img)
	}
}

// SaveToJPG encodes and saves image.RGBA to a file.
func SaveToJPG(img *image.RGBA, path string, quality int) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		jpeg.Encode(destFile, img, &jpeg.Options{Quality: quality})
	}
}
