package images4

import (
	"image"
	"path"
	"reflect"
	"testing"
)

const (
	testDir1 = "testdata"
	testDir2 = "resample"
)

func TestResizeByNearest(t *testing.T) {
	testDir := path.Join(testDir1, testDir2)
	tables := []struct {
		inFile     string
		srcX, srcY int
		outFile    string
		dstX, dstY int
	}{
		{"original.png", 533, 400,
			"nearest100x100.png", 100, 100},
		{"nearest100x100.png", 100, 100,
			"nearest533x400.png", 533, 400},
	}

	for _, table := range tables {
		inImg, err := Open(path.Join(testDir, table.inFile))
		if err != nil {
			t.Error("Cannot decode", path.Join(testDir, table.inFile))
		}
		outImg, err := Open(path.Join(testDir, table.outFile))
		if err != nil {
			t.Error("Cannot decode", path.Join(testDir, table.outFile))
		}
		resampled, srcSize := ResizeByNearest(inImg,
			image.Point{table.dstX, table.dstY})
		if !reflect.DeepEqual(
			outImg.(*image.RGBA), &resampled) ||
			table.srcX != srcSize.X ||
			table.srcY != srcSize.Y {
			t.Errorf(
				"Resample data do not match for %s and %s.",
				table.inFile, table.outFile)
		}
	}
}
