package util

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
)

func ReadVideoSingleFrame(videoPath string, frameNum int) (image.Image, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{
			fmt.Sprintf("gte(n,%d)", frameNum),
		}).
		Output("pipe:", ffmpeg.KwArgs{
			"vframes": 1,
			"format":  "image2",
			"vcodec":  "mjpeg",
		}).
		WithOutput(buf).
		Run()
	if err != nil {
		return nil, err
	}
	img, err := imaging.Decode(buf)
	if err != nil {
		return nil, err
	}
	return img, nil
}
