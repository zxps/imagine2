package files

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"imagine2/config"
	"imagine2/models"
	"imagine2/utils"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

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

// Sync ...
func (p *FilePartition) Sync() error {
	if p.Synchronized {
		return nil
	}
	if utils.IsDirExists(p.Fullpath) {
		return nil
	}

	err := os.MkdirAll(p.Fullpath, os.ModePerm)
	if err != nil {
		return err
	}

	p.Synchronized = true

	return nil
}

// SaveBase64 ...
func (p *FilePartition) SaveBase64(data string) error {
	tempDir, _ := ioutil.TempDir("", "imagine2")
	tempFile, err := ioutil.TempFile(tempDir, "upload-*")
	if err != nil {
		return err
	}

	defer tempFile.Close()

	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return errors.New("invalid file or image from base64 data")
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[idx+8:]))
	buff := bytes.Buffer{}
	_, err = buff.ReadFrom(reader)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(tempFile.Name(), buff.Bytes(), 0777)

	if err != nil {
		os.Remove(tempFile.Name())
		return err
	}

	mime, err := mimetype.DetectFile(tempFile.Name())
	if err != nil {
		os.Remove(tempFile.Name())
		return err
	}

	p.Mime = mime.String()
	p.Extension = strings.Replace(mime.Extension(), ".", "", 1)

	stats, _ := tempFile.Stat()
	p.Size = stats.Size()

	err = p.SaveFromFile(tempFile.Name())
	if err != nil {
		os.Remove(tempFile.Name())
		return err
	}

	os.Remove(tempFile.Name())

	return nil
}

// Generate ...
func (p *FilePartition) Generate() {
	pathSeparator := string(os.PathSeparator)
	p.Fullpath = GetFilespath()

	t := time.Now()

	p.Path = t.Format(config.Context.Service.GeneratorPathTimePattern)
	p.Path += strconv.Itoa(
		utils.RandInt(
			config.Context.Service.GeneratorMinFileDirIndex,
			config.Context.Service.GeneratorMaxFileDirIndex),
	)

	p.Fullpath = p.Fullpath + pathSeparator + p.Path

	for count := 500; count > 0; count-- {
		token := make([]byte, config.Context.Service.GeneratorFilenameLength)

		rand.Read(token)

		hash := sha1.New()
		hash.Write(token)

		p.Name = base64.URLEncoding.EncodeToString(hash.Sum(nil))

		for i := 0; i < len(removeFilenameChars); i++ {
			p.Name = strings.Replace(p.Name, removeFilenameChars[i], "", 1)
		}

		p.Name = p.Name[:config.Context.Service.GeneratorFilenameLength]
		p.Name = strings.ToLower(p.Name)

		if !utils.IsFileExists(p.Fullpath + pathSeparator + p.Name) {
			p.Success = true
			break
		}
	}
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
