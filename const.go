package images4

const (

	// Icon parameters.

	// Image resolution of the icon is very small
	// (11x11 pixels), therefore original image details
	// are lost in downsampling, except when source images
	// have very low resolution (e.g. favicons or simple
	// logos). This is useful from the privacy perspective
	// if you are to use generated icons in a large searchable
	// database.
	IconSize = 11 // Exported to be used in package imagehash.
	// Resampling rate defines how much information
	// (how many pixels) from the source image are used
	// to generate an icon. Too few will produce worse
	// comparisons. Too many will consume too much compute.
	samples = 12

	// Similarity parameters.

	// Cutoff value for color distance.
	colorDiff = 50
	// Cutoff coefficient for Euclidean distance (squared).
	euclCoeff = 0.2
	// Coefficient of sensitivity for Cb/Cr channels vs Y.
	chanCoeff = 2

	// Similarity thresholds.

	// Euclidean distance threshold (squared) for Y-channel.
	thY = float64(IconSize*IconSize) * float64(colorDiff*colorDiff) * euclCoeff
	// Euclidean distance threshold (squared) for Cb and Cr channels.
	thCbCr = thY * chanCoeff
	// Proportion similarity threshold (5%).
	thProp = 0.05

	// Auxiliary constants.

	numPix           = IconSize * IconSize
	largeIconSize    = IconSize*2 + 1
	resizedImgSize   = largeIconSize * samples
	invSamplePixels2 = 1 / float64(samples*samples)
	oneNinth         = 1 / float64(9)
	one255th         = 1 / float64(255)
	one255th2        = one255th * one255th
	sq255            = 255 * 255
	maxUint16        = 65535
)
