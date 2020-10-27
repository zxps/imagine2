package files

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	widthAndHeightResizePattern = `/(?P<width>[\d]{1,3})(x|X)(?P<height>[\d]{1,3})/`
	widthResizePattern          = `/(?P<width>[\d]{1,3})(x)/`
	heightResizePattern         = `/x(?P<height>[\d]{1,3})/`
)

// FileTransform ...
type FileTransform struct {
	SourceFilepath string
	Filepath       string

	FullSourceFilepath string
	FullFilepath       string

	FullTransformFilepath string

	Width        int
	Height       int
	Proportional bool
	Resize       bool

	Source *FilePartition
	Target *FilePartition
}

// ExtractTransform ...
func ExtractTransform(filepath string) *FileTransform {
	t := &FileTransform{
		FullTransformFilepath: "",
		SourceFilepath:        filepath,
		FullFilepath:          filepath,
		Filepath:              filepath,
		Width:                 0,
		Height:                0,
		Resize:                false,
		Proportional:          true,
		Source:                &FilePartition{},
		Target:                &FilePartition{},
	}

	FillPartitionByFilepath(t.Source, filepath)
	FillPartitionByFilepath(t.Target, filepath)

	compiled := false

	r := regexp.MustCompile(widthAndHeightResizePattern)
	matches := r.FindStringSubmatch(filepath)
	if len(matches) > 0 {
		t.Width, _ = strconv.Atoi(matches[1])
		t.Height, _ = strconv.Atoi(matches[3])

		switch matches[2] {
		case "x":
			t.Proportional = true
			break
		case "X":
			t.Proportional = false
		}

		t.Resize = true
		t.Filepath = strings.Replace(filepath, matches[0], "/", 1)
		compiled = true
	}

	if !compiled {
		r := regexp.MustCompile(widthResizePattern)
		matches := r.FindStringSubmatch(filepath)
		if len(matches) > 1 {
			t.Width, _ = strconv.Atoi(matches[1])
			t.Filepath = strings.Replace(filepath, matches[0], "/", 1)
			t.Resize = true
			compiled = true
		}
	}

	if !compiled {
		r := regexp.MustCompile(heightResizePattern)
		matches := r.FindStringSubmatch(filepath)
		if len(matches) > 1 {
			t.Height, _ = strconv.Atoi(matches[1])
			t.Filepath = strings.Replace(filepath, matches[0], "/", 1)
			t.Resize = true
			compiled = true
		}
	}

	if compiled && t.SourceFilepath != t.Filepath {
		FillPartitionByFilepath(t.Source, t.Filepath)
		t.Target.Fullpath = GetCachePath() + string(os.PathSeparator) + t.Target.Path
	}

	return t
}
