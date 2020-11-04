package image

import (
	"imagine2/files"
	"io/ioutil"

	"github.com/h2non/bimg"
)

// ProcessImageOptions ...
type ProcessImageOptions struct {
	Normalize bool
}

// ProcessImagePartition ...
func ProcessImagePartition(p *files.FilePartition, options ProcessImageOptions) error {
	buffer, err := ioutil.ReadFile(p.GetFilepath())

	if err != nil {
		return err
	}

	image := bimg.NewImage(buffer)

	size, err := image.Size()

	if err != nil {
		return err
	}

	p.ImageWidth = size.Width
	p.ImageHeight = size.Height

	if p.ImageWidth > 0 && p.ImageHeight > 0 {
		if options.Normalize {
			data, err := image.Process(bimg.Options{
				Width:  p.ImageWidth,
				Height: p.ImageHeight,
				//Quality: 90,
			})

			if err == nil {
				image = bimg.NewImage(data)
				ioutil.WriteFile(p.GetFilepath(), data, 0777)

				p.Size = int64(image.Length())
			}
		}
	}

	return nil
}
