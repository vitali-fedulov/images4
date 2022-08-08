package images4

import (
	"path"
	"testing"
)

func testPropSimilar(fA, fB string, isSimilar bool,
	t *testing.T) {
	p := path.Join("testdata", "proportions")
	imgA, err := Open(path.Join(p, fA))
	if err != nil {
		t.Error("Error opening image:", err)
	}
	imgB, err := Open(path.Join(p, fB))
	if err != nil {
		t.Error("Error opening image:", err)
	}
	iconA := Icon(imgA)
	iconB := Icon(imgB)

	if isSimilar == true {
		if !propSimilar(iconA, iconB) {
			t.Errorf("Expecting similarity of %v to %v.", fA, fB)
		}
	}
	if isSimilar == false {
		if propSimilar(iconA, iconB) {
			t.Errorf("Expecting non-similarity of %v to %v.", fA, fB)
		}
	}
}

func TestPropSimilar(t *testing.T) {
	testPropSimilar("100x130.png", "100x124.png", true, t)
	testPropSimilar("100x130.png", "100x122.png", false, t)
	testPropSimilar("130x100.png", "260x200.png", true, t)
	testPropSimilar("200x200.png", "260x200.png", false, t)
	testPropSimilar("130x100.png", "124x100.png", true, t)
	testPropSimilar("130x100.png", "122x100.png", false, t)
	testPropSimilar("130x100.png", "130x100.png", true, t)
	testPropSimilar("100x130.png", "130x100.png", false, t)
	testPropSimilar("124x100.png", "260x200.png", true, t)
	testPropSimilar("122x100.png", "260x200.png", false, t)
	testPropSimilar("100x124.png", "100x130.png", true, t)
}

func testEucSimilar(fA, fB string, isSimilar bool,
	t *testing.T) {
	p := path.Join("testdata", "euclidean")
	imgA, err := Open(path.Join(p, fA))
	if err != nil {
		t.Error("Error opening image:", err)
	}
	iconA := Icon(imgA)
	imgB, err := Open(path.Join(p, fB))
	if err != nil {
		t.Error("Error opening image:", err)
	}
	iconB := Icon(imgB)
	if isSimilar == true {
		if !eucSimilar(iconA, iconB) {
			t.Errorf("Expecting similarity of %v to %v.", fA, fB)
		}
	}
	if isSimilar == false {
		if eucSimilar(iconA, iconB) {
			t.Errorf("Expecting non-similarity of %v to %v.", fA, fB)
		}
	}
}

func TestEucSimilar(t *testing.T) {
	testEucSimilar("large.jpg", "distorted.jpg", true, t)
	testEucSimilar("large.jpg", "flipped.jpg", false, t)
	testEucSimilar("large.jpg", "small.jpg", true, t)
	testEucSimilar("small.gif", "small.jpg", true, t) // GIF test.
	testEucSimilar("uniform-black.png", "uniform-black.png", true, t)
	testEucSimilar("uniform-black.png", "uniform-white.png", false, t)
	testEucSimilar("uniform-green.png", "uniform-green.png", true, t)
	testEucSimilar("uniform-green.png", "uniform-white.png", false, t)
	testEucSimilar("uniform-white.png", "uniform-white.png", true, t)
}
