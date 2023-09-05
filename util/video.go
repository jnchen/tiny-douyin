package util

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ReadSingleFrameAsBytes(videoPath string, frameNum int) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{
			fmt.Sprintf("gte(n,%d)", frameNum),
		}).
		Output("pipe:", ffmpeg.KwArgs{
			"vframes": 1,
			"format":  "image2",
			"vcodec":  "mjpeg",
		}).
		WithOutput(buf).
		Run(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
