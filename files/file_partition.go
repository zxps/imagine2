package files

import (
	"crypto/sha1"
	"encoding/base64"
	"imagine2/config"
	"imagine2/models"
	"imagine2/utils"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const fileNameLength = 8
const minFileDir = 0
const maxFileDir = 5

var removeFilenameChars = [5]string{"-", "/", "_", "."}

// FilePartition ...
type FilePartition struct {
	Success      bool
	Synchronized bool
	SourceName   string
	Name         string
	Path         string
	Fullpath     string
	Extension    string
	Mime         string
	Size         int64
	ImageWidth   int
	ImageHeight  int
}

// NewFilePartition ...
func NewFilePartition(sourceName string) *FilePartition {
	p := &FilePartition{}

	p.SourceName = sourceName
	p.Success = false
	p.Synchronized = false
	p.Size = 0
	p.ImageHeight = 0
	p.ImageWidth = 0

	return p
}

// GetPartitionFromFile ...
func GetPartitionFromFile(file *models.File) *FilePartition {
	p := &FilePartition{}
	p.ImageHeight = file.Height
	p.ImageWidth = file.Width
	p.Name = file.Name
	p.Extension = file.Ext
	p.Mime = file.Mime
	p.Size = file.Size
	p.SourceName = file.SourceName
	p.Fullpath = GetFilespath() + string(os.PathSeparator) + file.Path
	return p
}

// Copy ...
func (p *FilePartition) Copy() *FilePartition {
	c := &FilePartition{}
	c = p
	return c
}

// SaveFromFile ...
func (p *FilePartition) SaveFromFile(filepath string) error {
	sourceFile, sourceErr := os.Open(filepath)
	if sourceErr != nil {
		return sourceErr
	}

	defer sourceFile.Close()

	targetFilepath := p.Fullpath + string(os.PathSeparator) + p.Name + "." + p.Extension
	targetFile, targetErr := os.Create(targetFilepath)
	if targetErr != nil {
		return targetErr
	}

	defer targetFile.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := sourceFile.Read(buffer)

		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := targetFile.Write(buffer[:n]); err != nil {
			return err
		}
	}

	return nil
}

// GetFilepath ...
func (p *FilePartition) GetFilepath() string {
	return p.Fullpath + string(os.PathSeparator) + p.Name + "." + p.Extension
}

// Generate ...
func (p *FilePartition) Generate() {
	pathSeparator := string(os.PathSeparator)
	p.Fullpath = GetFilespath()

	t := time.Now()

	p.Path = t.Format("2006/02/01")
	p.Path += strconv.Itoa(utils.RandInt(minFileDir, maxFileDir))

	p.Fullpath = p.Fullpath + pathSeparator + p.Path

	for count := 500; count > 0; count-- {
		token := make([]byte, fileNameLength)

		rand.Read(token)

		hash := sha1.New()
		hash.Write(token)

		p.Name = base64.URLEncoding.EncodeToString(hash.Sum(nil))

		for i := 0; i < len(removeFilenameChars); i++ {
			p.Name = strings.Replace(p.Name, removeFilenameChars[i], "", 1)
		}

		p.Name = p.Name[:fileNameLength]
		p.Name = strings.ToLower(p.Name)

		if !utils.IsFileExists(p.Fullpath + pathSeparator + p.Name) {
			p.Success = true
			break
		}
	}
}

// SyncFilePartition ...
func SyncFilePartition(p *FilePartition) (bool, error) {
	if p.Synchronized {
		return true, nil
	}
	if utils.IsDirExists(p.Fullpath) {
		return true, nil
	}

	err := os.MkdirAll(p.Fullpath, os.ModePerm)
	if err != nil {
		return false, err
	}
	p.Synchronized = true

	return true, nil
}

// GetFilespath ...
func GetFilespath() string {
	targetPath := config.Context.Service.FilesPath
	pathSeparator := string(os.PathSeparator)
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = config.Context.RootPath
		if !strings.HasSuffix(targetPath, pathSeparator) {
			targetPath += pathSeparator
		}

		targetPath += config.Context.Service.FilesPath
	}

	return targetPath
}

// GetCachePath ...
func GetCachePath() string {
	targetPath := config.Context.Service.CachePath
	pathSeparator := string(os.PathSeparator)
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = config.Context.RootPath
		if !strings.HasSuffix(targetPath, pathSeparator) {
			targetPath += pathSeparator
		}

		targetPath += config.Context.Service.CachePath
	}

	return targetPath
}
