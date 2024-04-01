# Find similar images with Go (LATEST VERSION)

Resized and near duplicate image comparison. **No dependencies**. For search in very large image sets use [imagehash2](https://github.com/vitali-fedulov/imagehash2) as a fast pre-filtering step.

**Demo**: [similar pictures search and clustering](https://vitali-fedulov.github.io/similar.pictures/) (pure in-browser JS app served [from](https://github.com/vitali-fedulov/similar.pictures)).

Major (semantic) versions have their own repositories and are mutually incompatible:
| Major version | Repository | Comment |
| ----------- | ---------- | ----------|
| 4 | images4 - this | recommended; fast hash prefiltering (re)moved to [imagehash2](https://github.com/vitali-fedulov/imagehash2) |
| 3 | [images3](https://github.com/vitali-fedulov/images3) | good, but less optimized |
| 1, 2 | [images](https://github.com/vitali-fedulov/images) | good, legacy code |

[Go doc](https://pkg.go.dev/github.com/vitali-fedulov/images4) - for full code documentation.

## Example of comparing 2 images

```go
package main

import (
	"fmt"
	"github.com/vitali-fedulov/images4"
)

func main() {

	// Opening and decoding images. Silently discarding errors.
	img1, _ := images4.Open("1.jpg")
	img2, _ := images4.Open("2.jpg")

	// Icons are compact hash-like image representations.
	icon1 := images4.Icon(img1)
	icon2 := images4.Icon(img2)

	// Comparison. Images are not used directly.
	// Use func CustomSimilar for different precision.
	if images4.Similar(icon1, icon2) {
		fmt.Println("Images are similar")
	} else {
		fmt.Println("Not similar")
	}

}
```

## Main functions

- `Open` decodes JPEG, PNG and GIF. But other types can be opened with third-party decoders, because the input to func 'Icon' is Golang image.Image. [Example fork](https://github.com/Pineapples27/images4) (not mine) expanded with support of WEBP images.

- `Icon` produces an image hash-like struct called "icon", which will be used for comparision. Side note: name "hash" is reserved for true hash tables in related package for faster comparison [imagehash2](https://github.com/vitali-fedulov/imagehash2).

- `Similar` gives a verdict whether 2 images are similar with well-tested default thresholds. Rotations and mirrors are not taken in account.

- `CustomSimilar` is like 'Similar' above, but allows modifying the default thresholds by multiplication coefficients. When the coefficients equal 1.0, those two functions are equivalent. When the coefficients are less than 1.0, the comparison is more precise, down to 0.0 for identical images.

## Advanced functions

- `Similar90270` is a superset of 'Similar' by additional comparison to images rotated ±90°. Such rotations are relatively common, even by accident when taking pictures on mobile phones.

- `CustomSimilar90270` is a custom func for rotations as above with 'CustomSimilar'.

- `EucMetric` is an alternative to 'CustomSimilar' when you need to know metric values, for example to sort by similarity. [Example](https://github.com/egor-romanov/png2gif/blob/main/main.go#L450) (not mine) of custom similarity function.

- `PropMetric` is as above for image proportions.

- `DefaultThresholds` prints default thresholds used in func 'Similar' and 'Similar90270', as a starting point for selecting thresholds on 'EucMetric' and 'PropMetric'.

- `Rotate90` turns an icon 90° clockwise. This is useful for developing custom similarity function for rotated images with 'EucMetric' and 'PropMetric'. With the function you can also compare to images rotated 180° (by applying 'Rotate90' twice).

- `ResizeByNearest` is an image resizing function useful for fast identification of identical images and development of custom distance metrics not involving any of the above comparison functions.


## Algorithm

Images are resampled and resized to squares of fixed size called "icons". Euclidean distance between the icons is used to give the similarity verdict. Also image proportions are used to avoid matching images of distinct shape.

[Detailed explanation](https://vitali-fedulov.github.io/similar.pictures/algorithm-for-perceptual-image-comparison.html), also as a [PDF](https://github.com/vitali-fedulov/research/blob/main/Algorithm%20for%20perceptual%20image%20comparison.pdf).

## Speed and precision

**To considerably accelerate comparison in large image collections** (thousands and more), use hash-table pre-filtering with package [imagehash2](https://github.com/vitali-fedulov/imagehash2).

**To considerably accelerate image decoding** you can generate icons for embedded image thumbnails. Specifically, many JPEG images contain [EXIF thumbnails](https://vitali-fedulov.github.io/similar.pictures/jpeg-thumbnail-reader.html). Packages to read thumbnails: [1](https://github.com/dsoprea/go-exif) and [2](https://github.com/rwcarlsen/goexif). A note of caution: in rare cases there could be [issues](https://security.stackexchange.com/questions/116552/the-history-of-thumbnails-or-just-a-previous-thumbnail-is-embedded-in-an-image/201785#201785) with thumbnails not matching image content. EXIF standard specification: [1](https://www.media.mit.edu/pia/Research/deepview/exif.html) and [2](https://www.exif.org/Exif2-2.PDF).

**An alternative method to increase precision** instead of func 'CustomSimilar' is to generate icons for image sub-regions and compare those icons.
