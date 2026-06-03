package ffmpeg

type ConvertDto struct {
	InputFilePath      string
	OutputFilePath     string
	Dimension          Dimension
	ConstantRateFactor *int
	VariableBitrate    *int
	ForcePixelFormat   *string
}

func Convert(c ConvertDto) error {
	panic("not implemented")
}
