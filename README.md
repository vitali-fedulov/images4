# Find similar images with Go (LATEST VERSION)

Near duplicates and resized images can be found with the package. No dependencies.

**Demo**: [similar pictures search and clustering](https://vitali-fedulov.github.io/similar.pictures/) (deployed [from](https://github.com/vitali-fedulov/similar.pictures)).

**Semantic versions**:
- v4 (/images4) - this repository, latest recommended,
- [v3](https://github.com/vitali-fedulov/images3) (/images3),
- [v1/v2](https://github.com/vitali-fedulov/images) (/images).

All versions will be kept available indefinitely.

Release note (v4): simplified func `Icon`; more than 2x reduction of icon memory footprint; removal of all dependencies; removal of hashes (a separate new package [imagehash](https://github.com/vitali-fedulov/imagehash) can be used for fast large scale preliminary search); fixed GIF support; new func `IconNN`. The main improvements in v4 are package simplification and memory footprint reduction.

### Key functions

`Open` supports JPEG, PNG and GIF. But other image types can be used through third-party decoders, because input for func `Icon` is Golang `image.Image`.

`Icon` produces "image hashes" called "icons", which will be used for comparision.

`Similar` gives a verdict whether 2 images are similar with well-tested default thresholds.

`EucMetric` can be used instead, when you need different precision or want to sort by similarity. Func `PropMetric` can be used for customization of image proportion threshold.

[Go doc](https://pkg.go.dev/github.com/vitali-fedulov/images4) for code reference.

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

	// Icons are compact image representations (image "hashes").
	// Name "hash" is not used intentionally.
	icon1 := images4.Icon(img1)
	icon2 := images4.Icon(img2)

	// Comparison.
	// Images are not used directly. Icons are used instead,
	// because they have tiny memory footprint and fast to compare.
	if images4.Similar(icon1, icon2) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```

## Algorithm

[Detailed explanation](https://vitali-fedulov.github.io/similar.pictures/algorithm-for-perceptual-image-comparison.html), also as a [PDF](https://github.com/vitali-fedulov/research/blob/main/Algorithm%20for%20perceptual%20image%20comparison.pdf).

Summary: Images are resized in a special way to squares of fixed size called "icons". Euclidean distance between the icons is used to give the similarity verdict. Also image proportions are used to avoid matching images of distinct shape.

## Customization suggestions

**To increase precision** you can either use your own thresholds in func `EucMetric` (and `PropMetric`) OR generate icons for image sub-regions and compare those icons.

**To speedup file processing** you may want to generate icons for available image thumbnails. Specifically, many JPEG images contain [EXIF thumbnails](https://vitali-fedulov.github.io/similar.pictures/jpeg-thumbnail-reader.html), you could considerably speedup the reads by using decoded thumbnails to feed into func `Icon`. External packages to read thumbnails: [1](https://github.com/dsoprea/go-exif) and [2](https://github.com/rwcarlsen/goexif). A note of caution: in rare cases there could be [issues](https://security.stackexchange.com/questions/116552/the-history-of-thumbnails-or-just-a-previous-thumbnail-is-embedded-in-an-image/201785#201785) with thumbnails not matching image content. EXIF standard specification: [1](https://www.media.mit.edu/pia/Research/deepview/exif.html) and [2](https://www.exif.org/Exif2-2.PDF).

**To search in very large image collections** (billions or more), use preliminary hash-table-based filtering with package [imagehash](https://github.com/vitali-fedulov/imagehash).
