package images4

// Similar returns similarity verdict based on Euclidean
// and proportion similarity.
func Similar(iconA, iconB IconT) bool {

	if !propSimilar(iconA, iconB) {
		return false
	}
	if !eucSimilar(iconA, iconB) {
		return false
	}
	return true
}

// propSimilar gives a similarity verdict for image A and B based on
// their height and width. When proportions are similar, it returns
// true.
func propSimilar(iconA, iconB IconT) bool {
	return PropMetric(iconA, iconB) < thProp
}

// PropMetric gives image proportion similarity metric for image A
// and B. The smaller the metric the more similar are images by their
// x-y size.
func PropMetric(iconA, iconB IconT) (m float64) {

	// Filtering is based on rescaling a narrower side of images to 1,
	// then cutting off at threshold of a longer image vs shorter image.
	xA, yA := float64(iconA.ImgSize.X), float64(iconA.ImgSize.Y)
	xB, yB := float64(iconB.ImgSize.X), float64(iconB.ImgSize.Y)

	if xA <= yA { // x to 1.
		yA = yA / xA
		yB = yB / xB
		if yA > yB {
			m = (yA - yB) / yA
		} else {
			m = (yB - yA) / yB
		}
	} else { // y to 1.
		xA = xA / yA
		xB = xB / yB
		if xA > xB {
			m = (xA - xB) / xA
		} else {
			m = (xB - xA) / xB
		}
	}
	return m
}

// eucSimilar gives a similarity verdict for image A and B based
// on Euclidean distance between pixel values of their icons.
// When the distance is small, the function returns true.
// iconA and iconB are generated with the Icon function.
// eucSimilar wraps EucMetric with well-tested thresholds.
func eucSimilar(iconA, iconB IconT) bool {

	m1, m2, m3 := EucMetric(iconA, iconB)

	return m1 < thY && // Luma as most sensitive.
		m2 < thCbCr &&
		m3 < thCbCr
}

// EucMetric returns Euclidean distances between 2 icons.
// These are 3 metrics corresponding to each color channel.
// Distances are squared, not to waste CPU on square root calculations.
// Note: color channels of icons are YCbCr (not RGB).
func EucMetric(iconA, iconB IconT) (m1, m2, m3 float64) {

	var cA, cB uint16
	for i := 0; i < numPix; i++ {
		// Channel 1.
		cA = iconA.Pixels[i]
		cB = iconB.Pixels[i]
		m1 += ((float64(cA) - float64(cB)) * one255th2 * (float64(cA) - float64(cB)))
		// Channel 2.
		cA = iconA.Pixels[i+numPix]
		cB = iconB.Pixels[i+numPix]
		m2 += ((float64(cA) - float64(cB)) * one255th2 * (float64(cA) - float64(cB)))
		// Channel 3.
		cA = iconA.Pixels[i+2*numPix]
		cB = iconB.Pixels[i+2*numPix]
		m3 += ((float64(cA) - float64(cB)) * one255th2 * (float64(cA) - float64(cB)))

	}

	return m1, m2, m3
}
