package image

import (
	"imagine2/files"
	"io/ioutil"

	"github.com/h2non/bimg"
)

// TransformFile ...
func TransformFile(t *files.FileTransform) error {
	if t.Resize {
		buffer, err := ioutil.ReadFile(t.Source.GetFilepath())

		if err != nil {
			return err
		}

		image := bimg.NewImage(buffer)

		dataResized, err := image.Resize(t.Width, t.Height)

		if err != nil {
			return err
		}

		return files.WriteBytesToFilepath(dataResized, t.Target.GetFilepath())
	}

	return nil
}
