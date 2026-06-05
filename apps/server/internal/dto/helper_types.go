package dto

import "github.com/slugger7/exorcist/apps/server/internal/ffmpeg"

type Dimension struct {
	Height *int `json:"height"`
	Width  *int `json:"width"`
}

func (d *Dimension) ToFfmpegDto() *ffmpeg.Dimension {
	v := ffmpeg.Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*v.Height = *d.Height
	*v.Width = *d.Width

	return &v
}

func (d *Dimension) FromFfmpegDto(m *ffmpeg.Dimension) *Dimension {
	if m.Height != nil {
		if d.Height == nil {
			d.Height = new(int)
		}
		*d.Height = *m.Height
	}

	if m.Width != nil {
		if d.Width == nil {
			d.Width = new(int)
		}
		*d.Width = *m.Width
	}

	return d
}
