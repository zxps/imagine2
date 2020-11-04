package files

import (
	"imagine2/models"
	"imagine2/utils"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	widthAndHeightResizePattern = `/(?P<width>[\d]{1,3})(x|X)(?P<height>[\d]{1,3})/`
	widthResizePattern          = `/(?P<width>[\d]{1,3})(x)/`
	heightResizePattern         = `/x(?P<height>[\d]{1,3})/`
)

// TransformedPartition ...
type TransformedPartition struct {
	Name         string `json:"name"`
	FullFilepath string `json:"full_filepath"`
	Filesize     int64  `json:"filesize"`
}

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

// InvalidateFileTransforms ...
func InvalidateFileTransforms(file models.File) int {
	transforms, _ := GetAllFileTransforms(file)
	resultCount := int(0)
	for _, transform := range transforms {
		err := os.Remove(transform.FullFilepath)
		if err == nil {
			resultCount++
		}
	}

	return resultCount
}

// GetAllFileTransforms ...
func GetAllFileTransforms(file models.File) ([]TransformedPartition, error) {
	cachePath := GetCachePath()
	sep := string(os.PathSeparator)

	transformPath := cachePath + string(os.PathSeparator) + file.Path

	dirs, err := ioutil.ReadDir(transformPath)
	if err != nil {
		return nil, err
	}

	var result []TransformedPartition = make([]TransformedPartition, 0)
	for _, d := range dirs {
		transformFilePath := cachePath + sep + file.Path + sep + d.Name() + sep + file.Fullname
		if utils.IsFileExists(transformFilePath) {
			size := int64(0)
			stats, err := os.Stat(transformFilePath)
			if err == nil {
				size = stats.Size()
			}

			result = append(result, TransformedPartition{
				Name:         d.Name(),
				FullFilepath: transformFilePath,
				Filesize:     size,
			})
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Filesize > result[j].Filesize
	})

	return result, nil
}
