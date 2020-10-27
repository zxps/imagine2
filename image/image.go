package image

import (
	"imagine2/files"
	"io/ioutil"

	"github.com/h2non/bimg"
)

// LoadImageInfo ...
func LoadImageInfo(p *files.FilePartition) error {
	buffer, err := ioutil.ReadFile(p.GetFilepath())

	if err != nil {
		return err
	}

	size, err := bimg.NewImage(buffer).Size()

	if err != nil {
		return err
	}

	p.ImageWidth = size.Width
	p.ImageHeight = size.Height

	return nil
}
