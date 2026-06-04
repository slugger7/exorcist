package ffmpeg

import (
	"fmt"
	"os"

	errs "github.com/slugger7/exorcist/apps/server/internal/errors"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type ConvertDto struct {
	InputFilePath      string
	OutputFilePath     string
	Dimension          Dimension
	ConstantRateFactor *int
	VariableBitrate    *int
	ForcePixelFormat   *string
}

func Convert(c ConvertDto) error {
	if *c.Dimension.Height <= 0 {
		return fmt.Errorf(ErrNegativeHeight, *c.Dimension.Height)
	}
	if *c.Dimension.Width <= 0 {
		return fmt.Errorf(ErrNegativeWidth, *c.Dimension.Width)
	}

	err := ffmpeg_go.Input(c.InputFilePath).Output(c.OutputFilePath,
		ffmpeg_go.KwArgs{"vf": fmt.Sprintf("scale=%v:%v", *c.Dimension.Width, *c.Dimension.Height)}).
		Run()
	if err != nil {
		str := err.Error()
		_ = str
		_ = os.Remove(c.OutputFilePath)
		return errs.BuildError(err, "error converting %v to %v", c.InputFilePath, c.OutputFilePath)
	}

	return nil
}
