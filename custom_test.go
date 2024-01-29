package images4

import (
	"path"
	"testing"
)

func TestCustomSimilar(t *testing.T) {

	// Proportions test.

	i1, _ := Open(path.Join("testdata", "euclidean", "distorted.jpg"))
	i2, _ := Open(path.Join("testdata", "euclidean", "large.jpg"))

	icon1 := Icon(i1)
	icon2 := Icon(i2)

	if Similar(icon1, icon2) {
		t.Errorf("distorted.jpg is NOT similar to large.jpg")
	}

	if !CustomSimilar(icon1, icon2, &CustomCoefficients{1, 1, 1, 10}) {
		t.Errorf("distorted.jpg IS similar to large.jpg, assuming proportion differences are widely tolerated.")
	}

	// Euclidean tests.

	i1, _ = Open(path.Join("testdata", "custom", "1.jpg"))
	i2, _ = Open(path.Join("testdata", "custom", "2.jpg"))

	icon1 = Icon(i1)
	icon2 = Icon(i2)

	if !Similar(icon1, icon2) {
		t.Errorf("1.jpg is GENERALLY similar to 2.jpg")
	}

	// Luma.
	if CustomSimilar(icon1, icon2, &CustomCoefficients{0, 1, 1, 1}) {
		t.Errorf("1.jpg is NOT IDENTICAL to 2.jpg")
	}

	// Luma.
	if CustomSimilar(icon1, icon2, &CustomCoefficients{0.4, 1, 1, 1}) {
		t.Errorf("1.jpg is similar to 2.jpg, BUT NOT VERY SIMILAR")
	}

	// Chrominance b.
	if CustomSimilar(icon1, icon2, &CustomCoefficients{1, 0.1, 1, 1}) {
		t.Errorf("1.jpg is similar to 2.jpg, BUT NOT VERY SIMILAR")
	}

	// Chrominance c.
	if CustomSimilar(icon1, icon2, &CustomCoefficients{1, 1, 0.1, 1}) {
		t.Errorf("1.jpg is similar to 2.jpg, BUT NOT VERY SIMILAR")
	}

	// Image comparison to itself (or its own copy).

	if !CustomSimilar(icon1, icon1, &CustomCoefficients{0, 0, 0, 0}) {
		t.Errorf("1.jpg IS IDENTICAL to itself")
	}

	if !CustomSimilar(icon1, icon1, &CustomCoefficients{0.5, 0.5, 0.5, 0.5}) {
		t.Errorf("1.jpg IS IDENTICAL to itself")
	}

	if !CustomSimilar(icon1, icon1, &CustomCoefficients{1, 1, 1, 1}) {
		t.Errorf("1.jpg IS IDENTICAL to itself")
	}

}
