package ffmpeg

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ImageAt_NegativeWidth(t *testing.T) {
	width := -1

	err := ImageAt("", 0, "", width, 1)

	assert.ErrorContains(t, err, fmt.Sprintf(ErrNegativeWidth, width))
}

func Test_ImageAt_NegativeHeight(t *testing.T) {
	height := -1

	err := ImageAt("", 0, "", 1, height)

	assert.ErrorContains(t, err, fmt.Sprintf(ErrNegativeHeight, height))
}

func Test_ImageAt_Success(t *testing.T) {
	width, height := 20, 60
	time := float64(3)

	err := ImageAt(testVideoPath, time, testImagePath, width, height)
	assert.Nil(t, err)
	assert.FileExists(t, testImagePath)

	os.Remove(testImagePath)
}

func Test_ScaleWidthByHeight(t *testing.T) {
	width, height, max := 100, 1000, 400

	newWidth := ScaleWidthByHeight(height, width, max)

	if newWidth != 40 {
		t.Errorf("calculated height was not 400 it was %v", newWidth)
	}
}

func Test_ScaleHeightByWidth(t *testing.T) {
	width, height, max := 1000, 100, 400

	newHeight := ScaleHeightByWidth(height, width, max)

	if newHeight != 40 {
		t.Errorf("calculated height was not 400 it was %v", newHeight)
	}
}

func Test_DetermineDimensions_WithNoWantedHeightOrWidth_ShouldReturnCurrentDimension(t *testing.T) {
	wantedDimension := Dimension{}
	currentDimension := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*currentDimension.Height = 666
	*currentDimension.Width = 666

	actualDimension := DetermineDimensions(wantedDimension, currentDimension)

	if *actualDimension.Height != *currentDimension.Height ||
		*actualDimension.Width != *currentDimension.Width {
		t.Errorf("value returned did not match the current dimension: %vx%v", *actualDimension.Width, *actualDimension.Height)
	}
}

func Test_DetermineDimensions_WithWantedHeightAndWidthDefined_ShouldReturnWantedDimension(t *testing.T) {
	wantedDimension := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*wantedDimension.Height = 420
	*wantedDimension.Width = 420
	currentDimension := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*currentDimension.Height = 666
	*currentDimension.Width = 666

	actualDimension := DetermineDimensions(wantedDimension, currentDimension)

	if *actualDimension.Height != *wantedDimension.Height ||
		*actualDimension.Width != *wantedDimension.Width {
		t.Errorf("value returned did not match the wanted dimension: %vx%v", *actualDimension.Width, *actualDimension.Height)
	}
}

func Test_DetermineDimensions_WithWantedHeightDefinedButNoWidth_ShouldReturnDimensionWithCalculatedWidth(t *testing.T) {
	wantedDimension := Dimension{
		Height: new(int),
	}
	*wantedDimension.Height = 69
	currentDimension := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*currentDimension.Height = 1000
	*currentDimension.Width = 2000

	actualDimension := DetermineDimensions(wantedDimension, currentDimension)

	expectedDimension := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*expectedDimension.Height = *wantedDimension.Height
	*expectedDimension.Width = *wantedDimension.Height * 2

	if *actualDimension.Height != *expectedDimension.Height ||
		*actualDimension.Width != *expectedDimension.Width {
		t.Errorf("value returned did not match the expected dimension(%vx%v): %vx%v",
			*expectedDimension.Width, *expectedDimension.Height,
			*actualDimension.Width, *actualDimension.Height)
	}
}

func Test_DetermineDimensions_WithWantedWidthDefinedButNoHeight_ShouldReturnDimensionWithCalculatedHeight(t *testing.T) {
	wantedDimension := Dimension{
		Width: new(int),
	}
	*wantedDimension.Width = 138
	currentDimension := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*currentDimension.Height = 1000
	*currentDimension.Width = 2000

	actualDimension := DetermineDimensions(wantedDimension, currentDimension)

	expectedDimension := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*expectedDimension.Height = int(float64(*wantedDimension.Width) / 2.0)
	*expectedDimension.Width = *wantedDimension.Width

	if *actualDimension.Height != *expectedDimension.Height ||
		*actualDimension.Width != *expectedDimension.Width {
		t.Errorf("value returned did not match the expected dimension(%vx%v): %vx%v",
			*expectedDimension.Width, *expectedDimension.Height,
			*actualDimension.Width, *actualDimension.Height)
	}
}

func Test_ScaleByMaxDimension_WithWidthGreaterThanMaxDimension_ShouldScaleHeightAndSetWidthToMaxDimension(t *testing.T) {
	maxDimension := 10

	c := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*c.Height = 100
	*c.Width = 200

	d := ScaleByMaxDimension(maxDimension, c)

	assert.Equal(t, 5, *d.Height, "height was not correctly scaled")
	assert.Equal(t, maxDimension, *d.Width, "width was not set to max dimension")
}

func Test_ScaleByMaxDimension_WithHeightGreaterThanMaxDimension_ShouldScaleWidthAndSetHightToMaxDimension(t *testing.T) {
	maxDimension := 10

	c := Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*c.Height = 200
	*c.Width = 100

	d := ScaleByMaxDimension(maxDimension, c)

	assert.Equal(t, 5, *d.Width, "width was not correctly scaled")
	assert.Equal(t, maxDimension, *d.Height, "height was not set to max dimension")
}
