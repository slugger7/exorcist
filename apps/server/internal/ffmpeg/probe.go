package ffmpeg

import (
	"encoding/json"
	"errors"

	ffmpegGo "github.com/u2takey/ffmpeg-go"
)

type Stream struct {
	Height    *int   `json:"height"`
	Width     *int   `json:"width"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Duration string `json:"duration"`
}

type Probe struct {
	Format  *Format  `json:"format"`
	Streams []Stream `json:"streams"`
}

func UnmarshalProbeData(probeData string) (*Probe, error) {
	var data *Probe
	err := json.Unmarshal([]byte(probeData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UnmarshalledProbe(path string) (*Probe, error) {
	probeData, err := ffmpegGo.Probe(path)
	if err != nil {
		return nil, err
	}

	data, err := UnmarshalProbeData(probeData)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetDimensions(streams []Stream) (*Dimension, error) {
	for _, v := range streams {
		if v.CodecType == "video" {
			return &Dimension{Height: v.Height, Width: v.Width}, nil
		}
	}

	return nil, errors.New("could not extract the height and width from the probe data streams")
}
