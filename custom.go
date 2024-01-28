package images4

// Threshold multiplication coefficients for func CustomSimilar.
// When all values equal 1.0 func CustomSimilar is equivalent
// to func Similar. By setting those values less than 1, similarity
// comparison becomes stricter (more precise). Values larger than 1
// will generalize more and show more false positives. When uncertain,
// setting all coefficients to 1.0 is the safe starting point.
type CustomCoefficients struct {
	Y    float64 // Luma (grayscale information).
	Cb   float64 // Chrominance b (color information).
	Cr   float64 // Chrominance r (color information).
	Prop float64 // Proportion tolerance (how similar are image borders).
}

// CustomSimilar is like Similar, except it allows changing default
// thresholds by multiplying them. The practically useful range of
// the coefficients is [0, 1.0), but can be equal or larger than 1
// if necessary. All coefficients set to 0 correspond to identical images,
// for example an image file copy. All coefficients equal to 1 make func
// CustomSimilar equivalent to func Similar.
func CustomSimilar(iconA, iconB IconT, coeff CustomCoefficients) bool {

	if !customPropSimilar(iconA, iconB, coeff) {
		return false
	}
	if !customEucSimilar(iconA, iconB, coeff) {
		return false
	}
	return true
}

func customPropSimilar(iconA, iconB IconT, coeff CustomCoefficients) bool {
	return PropMetric(iconA, iconB) <= thProp*coeff.Prop
}

func customEucSimilar(iconA, iconB IconT, coeff CustomCoefficients) bool {

	m1, m2, m3 := EucMetric(iconA, iconB)

	return m1 <= thY*coeff.Y &&
		m2 <= thCbCr*coeff.Cb &&
		m3 <= thCbCr*coeff.Cr
}

// Similar90270 works like Similar, but also considers rotations of ±90°.
// Those are rotations users might reasonably often do.
func CustomSimilar90270(iconA, iconB IconT, coeff CustomCoefficients) bool {

	if CustomSimilar(iconA, iconB, coeff) {
		return true
	}

	// iconB rotated 90 degrees.
	if CustomSimilar(iconA, Rotate90(iconB), coeff) {
		return true
	}

	// As if iconB was rotated 270 degrees.
	if CustomSimilar(Rotate90(iconA), iconB, coeff) {
		return true
	}

	return false
}
