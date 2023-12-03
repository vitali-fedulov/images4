package images4

import (
	"image"
	"image/color"
)

// Icon has square shape. Its pixels are uint16 values
// in 3 channels. uint16 is intentional to preserve color
// relationships from the full-size image. It is a 255-
// premultiplied color value in [0, 255] range.
type IconT struct {
	Pixels  []uint16    // Visual signature.
	ImgSize image.Point // Original image size.
}

// Icon generates a normalized image signature ("icon").
// Generated icons can then be stored in a database and used
// for comparison. Icon is the recommended function,
// vs less robust func IconNN.
func Icon(img image.Image) IconT {

	icon := IconNN(img)

	// Maximizing icon contrast. This to reflect on the human visual
	// experience, when high contrast (normalized) images are easier
	// to see. Normalization also compensates for considerable loss
	// of visual information during scarse resampling during
	// icon creation step.
	icon.normalize()

	return icon
}

// IconNN generates a NON-normalized image signature (icon).
// Icons made with IconNN can be used instead of icons made with
// func Icon, but mostly for experimental purposes, allowing
// better understand how the algorithm works, or performing
// less agressive customized normalization. Not for general use.
func IconNN(img image.Image) IconT {

	// Resizing to a large icon approximating average color
	// values of the source image. YCbCr space is used instead
	// of RGB for better results in image comparison.
	resImg, imgSize := ResizeByNearest(
		img, image.Point{resizedImgSize, resizedImgSize})
	largeIcon := sizedIcon(largeIconSize)
	var r, g, b, sumR, sumG, sumB uint32
	var yc, cb, cr float64
	// For each pixel of the largeIcon.
	for x := 0; x < largeIconSize; x++ {
		for y := 0; y < largeIconSize; y++ {
			sumR, sumG, sumB = 0, 0, 0
			// Sum over pixels of resImg.
			for m := 0; m < samples; m++ {
				for n := 0; n < samples; n++ {
					r, g, b, _ =
						resImg.At(
							x*samples+m, y*samples+n).RGBA()
					sumR += r >> 8
					sumG += g >> 8
					sumB += b >> 8
				}
			}
			Set(largeIcon, largeIconSize, image.Point{x, y},
				float64(sumR)*invSamplePixels2,
				float64(sumG)*invSamplePixels2,
				float64(sumB)*invSamplePixels2)
		}
	}

	// Box blur filter with resizing to the final icon of smaller size.

	icon := sizedIcon(IconSize)
	// Pixel positions in the final icon.
	var xd, yd int
	var c1, c2, c3, s1, s2, s3 float64

	// For pixels of source largeIcon with stride 2.
	for x := 1; x < largeIconSize-1; x += 2 {
		xd = x / 2
		for y := 1; y < largeIconSize-1; y += 2 {
			yd = y / 2
			// For each pixel of a 3x3 box.
			for n := -1; n <= 1; n++ {
				for m := -1; m <= 1; m++ {
					c1, c2, c3 =
						Get(largeIcon, largeIconSize,
							image.Point{x + n, y + m})
					s1, s2, s3 = s1+c1, s2+c2, s3+c3
				}
			}
			yc, cb, cr = yCbCr(
				s1*oneNinth, s2*oneNinth, s3*oneNinth)
			Set(icon, IconSize, image.Point{xd, yd},
				yc, cb, cr)
			s1, s2, s3 = 0, 0, 0
		}
	}

	icon.ImgSize = imgSize
	return icon
}

// EmptyIcon is an icon constructor in case you need an icon
// with nil values, for example for convenient error handling.
// Then you can use icon.Pixels == nil condition.
func EmptyIcon() (icon IconT) {
	icon = sizedIcon(IconSize)
	icon.Pixels = nil
	return icon
}

func sizedIcon(size int) (icon IconT) {
	icon.Pixels = make([]uint16, size*size*3)
	return icon
}

// ArrIndex gets a pixel position in 1D array from a point
// of 2D array. ch is color channel index (0 to 2).
func arrIndex(p image.Point, size, ch int) (index int) {
	return size*(ch*size+p.Y) + p.X
}

// Set places pixel values in an icon at a point.
// c1, c2, c3 are color values for each channel
// (RGB for example). Size is icon size.
// Exported to be used in package imagehash.
func Set(icon IconT, size int, p image.Point, c1, c2, c3 float64) {
	// Multiplication by 255 is basically encoding float64 as uint16.
	icon.Pixels[arrIndex(p, size, 0)] = uint16(c1 * 255)
	icon.Pixels[arrIndex(p, size, 1)] = uint16(c2 * 255)
	icon.Pixels[arrIndex(p, size, 2)] = uint16(c3 * 255)
}

// Get reads pixel values in an icon at a point.
// c1, c2, c3 are color values for each channel
// (RGB for example).
// Exported to be used in package imagehash.
func Get(icon IconT, size int, p image.Point) (c1, c2, c3 float64) {
	// Division by 255 is basically decoding uint16 into float64.
	c1 = float64(icon.Pixels[arrIndex(p, size, 0)]) * one255th
	c2 = float64(icon.Pixels[arrIndex(p, size, 1)]) * one255th
	c3 = float64(icon.Pixels[arrIndex(p, size, 2)]) * one255th
	return c1, c2, c3
}

// yCbCr transforms RGB components to YCbCr. This is a high
// precision version different from the Golang image library
// operating on uint8.
func yCbCr(r, g, b float64) (yc, cb, cr float64) {
	yc = 0.299000*r + 0.587000*g + 0.114000*b
	cb = 128 - 0.168736*r - 0.331264*g + 0.500000*b
	cr = 128 + 0.500000*r - 0.418688*g - 0.081312*b
	return yc, cb, cr
}

// Normalize stretches histograms for the 3 channels of an icon, so that
// min/max values of each channel are 0/255 correspondingly.
// Note: values of IconT are premultiplied by 255, thus having maximum
// value of sq255 constant corresponding to display color value of 255.
func (src IconT) normalize() {

	var c1Min, c2Min, c3Min, c1Max, c2Max, c3Max uint16
	c1Min, c2Min, c3Min = maxUint16, maxUint16, maxUint16
	c1Max, c2Max, c3Max = 0, 0, 0
	var scale float64
	var n int

	// Looking for extreme values.
	for n = 0; n < numPix; n++ {
		// Channel 1.
		if src.Pixels[n] > c1Max {
			c1Max = src.Pixels[n]
		}
		if src.Pixels[n] < c1Min {
			c1Min = src.Pixels[n]
		}
		// Channel 2.
		if src.Pixels[n+numPix] > c2Max {
			c2Max = src.Pixels[n+numPix]
		}
		if src.Pixels[n+numPix] < c2Min {
			c2Min = src.Pixels[n+numPix]
		}
		// Channel 3.
		if src.Pixels[n+2*numPix] > c3Max {
			c3Max = src.Pixels[n+2*numPix]
		}
		if src.Pixels[n+2*numPix] < c3Min {
			c3Min = src.Pixels[n+2*numPix]
		}
	}

	// Normalization.
	if c1Max != c1Min { // Must not divide by zero.
		scale = sq255 / (float64(c1Max) - float64(c1Min))
		for n = 0; n < numPix; n++ {
			src.Pixels[n] = uint16(
				(float64(src.Pixels[n]) - float64(c1Min)) *
					scale)
		}
	}
	if c2Max != c2Min { // Must not divide by zero.
		scale = sq255 / (float64(c2Max) - float64(c2Min))
		for n = 0; n < numPix; n++ {
			src.Pixels[n+numPix] = uint16(
				(float64(src.Pixels[n+numPix]) - float64(c2Min)) *
					scale)
		}
	}
	if c3Max != c3Min { // Must not divide by zero.
		scale = sq255 / (float64(c3Max) - float64(c3Min))
		for n = 0; n < numPix; n++ {
			src.Pixels[n+2*numPix] = uint16(
				(float64(src.Pixels[n+2*numPix]) - float64(c3Min)) *
					scale)
		}
	}

}

// ToRGBA transforms a sized icon to image.RGBA. This is
// an auxiliary function to visually evaluate an icon.
func (icon IconT) ToRGBA(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			r, g, b := Get(icon, size, image.Point{x, y})
			img.Set(x, y,
				color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
	return img
}

// Rotate rotates an icon by 90 degrees clockwise.
func Rotate90(icon IconT) IconT {

	var c1, c2, c3 float64
	rotated := sizedIcon(IconSize)
	for x := 0; x < IconSize; x++ {
		for y := 0; y < IconSize; y++ {
			c1, c2, c3 = Get(icon, IconSize, image.Point{y, IconSize - 1 - x})
			Set(rotated, IconSize, image.Point{x, y},
				c1, c2, c3)
		}
	}

	// Swap image sizes.
	rotated.ImgSize.X, rotated.ImgSize.Y = icon.ImgSize.Y, icon.ImgSize.X

	return rotated
}
