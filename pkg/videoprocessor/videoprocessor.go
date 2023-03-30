package videoprocessor

import (
	"errors"
	"fmt"

	"github.com/3d0c/gmf"
)

type VideoProcessor struct {
	ImagePath string
}

func NewVideoProcessor(imagePath string) *VideoProcessor {
	return &VideoProcessor{
		ImagePath: imagePath,
	}
}

func (vp *VideoProcessor) OverlayImage(inputVideoPath, outputVideoPath string) error {
	gmf.InitFFmpeg()

	inputCtx, err := gmf.NewInputCtx(inputVideoPath)
	if err != nil {
		return err
	}
	defer inputCtx.CloseInputAndRelease()

	outputCtx, err := gmf.NewOutputCtx(outputVideoPath)
	if err != nil {
		return err
	}
	defer outputCtx.CloseOutputAndRelease()

	srcVideoStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		return err
	}

	codec, err := gmf.FindEncoder(gmf.AV_CODEC_ID_H264)
	if err != nil {
		return err
	}

	codecCtx := gmf.NewCodecCtx(codec)
	if codecCtx == nil {
		return errors.New("unable to allocate codec context")
	}
	defer gmf.Release(codecCtx)

	// TODO: Configure codec context with necessary settings
	// ...

	// TODO: Initialize and configure filter graph to overlay image
	// ...

	if err := outputCtx.WriteHeader(); err != nil {
		return err
	}

	// TODO: Process and write packets
	// ...

	if err := outputCtx.WriteTrailer(); err != nil {
		return err
	}

	fmt.Println("Video processed successfully:", outputVideoPath)
	return nil
}
