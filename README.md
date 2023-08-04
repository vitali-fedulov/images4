# Find similar images with Go (LATEST VERSION)

Resized images and near duplicates can be found with the package. **No dependencies**.

**Demo**: [similar pictures search and clustering](https://vitali-fedulov.github.io/similar.pictures/) (deployed [from](https://github.com/vitali-fedulov/similar.pictures)).

Major versions are semantic. They have own repositories and are mutually incompatible. The repositories will be kept available indefinitely.
| Major version | Repository | Comment |
| ----------- | ---------- | ----------|
| 4 | images4 - this | recommended |
| 3 | [images3](https://github.com/vitali-fedulov/images3) | good, but less optimized |
| 1, 2 | [images](https://github.com/vitali-fedulov/images) | good, legacy code |

## Example of comparing 2 images

```go
package main

import (
	"fmt"
	"github.com/vitali-fedulov/images4"
)

func main() {

	// Photos to compare.
	path1 := "1.jpg"
	path2 := "2.jpg"

	// Open files (discarding errors here).
	img1, _ := images4.Open(path1)
	img2, _ := images4.Open(path2)

	// Icons are compact image representations (image "hashes"). Name "hash" is reserved for "true" hashes in package imagehash.
	icon1 := images4.Icon(img1)
	icon2 := images4.Icon(img2)

	// Comparison. Images are not used directly. Icons are used instead, because they have tiny memory footprint and fast to compare. If you need to include images rotated right and left use func Similar90270.
	if images4.Similar(icon1, icon2) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```

## Main functions

- `Open` supports JPEG, PNG and GIF. But other image types can be used through third-party decoders, because input for func `Icon` is Golang `image.Image`. [Example fork](https://github.com/Pineapples27/images4) (not mine) expanded with support of WEBP images.

- `Icon` produces "image hashes" called "icons", which will be used for comparision.

- `Similar` gives a verdict whether 2 images are similar with well-tested default thresholds. To see the thresholds use `DefaultThresholds`. Rotations and mirrors are not taken in account.

- `Similar90270` is a superset of `Similar` by additional comparison to images rotated ±90°. Such rotations are relatively common, even by accident when taking pictures on mobile phones.

- `EucMetric` can be used instead of `Similar` when you need different precision or want to sort by similarity. [Example](https://github.com/egor-romanov/png2gif/blob/main/main.go#L450) (not mine) of custom similarity function.

- `PropMetric` allows customization of image proportion threshold.

- `DefaultThresholds` prints default thresholds used in func `Similar` and `Similar90270`, as a starting point for selecting thresholds on `EucMetric` and `PropMetric`.

- `Rotate90` turns an icon 90° clockwise. This is useful for developing custom similarity function for rotated images with `EucMetric` and `PropMetric`. With the function you can also compare to images rotated 180° (by applying `Rotate90` twice).

[Go doc](https://pkg.go.dev/github.com/vitali-fedulov/images4) for code reference.


## Algorithm

[Detailed explanation](https://vitali-fedulov.github.io/similar.pictures/algorithm-for-perceptual-image-comparison.html), also as a [PDF](https://github.com/vitali-fedulov/research/blob/main/Algorithm%20for%20perceptual%20image%20comparison.pdf).

Summary: Images are resized in a special way to squares of fixed size called "icons". Euclidean distance between the icons is used to give the similarity verdict. Also image proportions are used to avoid matching images of distinct shape.

## Customization suggestions

**To increase precision** you can either use your own thresholds in func `EucMetric` (and `PropMetric`) OR generate icons for image sub-regions and compare those icons.

**To speedup file processing** you may want to generate icons for available image thumbnails. Specifically, many JPEG images contain [EXIF thumbnails](https://vitali-fedulov.github.io/similar.pictures/jpeg-thumbnail-reader.html), you could considerably speedup the reads by using decoded thumbnails to feed into func `Icon`. External packages to read thumbnails: [1](https://github.com/dsoprea/go-exif) and [2](https://github.com/rwcarlsen/goexif). A note of caution: in rare cases there could be [issues](https://security.stackexchange.com/questions/116552/the-history-of-thumbnails-or-just-a-previous-thumbnail-is-embedded-in-an-image/201785#201785) with thumbnails not matching image content. EXIF standard specification: [1](https://www.media.mit.edu/pia/Research/deepview/exif.html) and [2](https://www.exif.org/Exif2-2.PDF).

**To search in very large image collections** (billions or more), use preliminary hash-table-based filtering with package [imagehash](https://github.com/vitali-fedulov/imagehash).
