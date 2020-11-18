package files

import (
	"bytes"
	"encoding/base64"
	"errors"
	"imagine2/utils"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// UploadFileFromRequest - ...
func UploadFileFromRequest(ctx *fasthttp.RequestCtx, name string) (*FilePartition, error) {
	handler, err := ctx.FormFile(name)
	if err != nil {
		log.Error("file upload error ", err.Error())
		return nil, err
	}

	p := NewFilePartition(handler.Filename)
	p.Generate()
	if !p.Success {
		log.Error("unable to generate new file p")
		return p, errors.New("unable to generate new file p")
	}

	err = p.Sync()
	if err != nil {
		log.Error("unable to synchronize file p: ", err.Error())
		return p, err
	}

	tempDir, _ := ioutil.TempDir("", "imagine2")
	tempFile, err := ioutil.TempFile(tempDir, "upload-*")
	if err != nil {
		log.Error("unable to create temp file", err.Error())
		return p, err
	}

	defer tempFile.Close()

	log.Info("temp file: ", tempFile.Name())

	err = fasthttp.SaveMultipartFile(handler, tempFile.Name())
	if err != nil {
		log.Error("unable to save temp uploaded file")
		os.Remove(tempFile.Name())
		return p, err
	}

	mime, err := mimetype.DetectFile(tempFile.Name())
	if err != nil {
		log.Error("unable to detect file mime info: ", err.Error())
		os.Remove(tempFile.Name())
		return p, err
	}

	p.Mime = mime.String()
	p.Extension = strings.Replace(mime.Extension(), ".", "", 1)
	p.Size = handler.Size

	log.Info("uploaded file mime: ", p.Mime)
	log.Info("uploaded file ext: ", p.Extension)

	if err := p.SaveFromFile(tempFile.Name()); err != nil {
		log.Info("unable to save file to p: ", err.Error())
		os.Remove(tempFile.Name())
		return p, err
	}

	os.Remove(tempFile.Name())

	return p, nil
}

// UploadFileFromBase64 - ...
func UploadFileFromBase64(p *FilePartition, data string) error {
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

	if err := p.SaveFromFile(tempFile.Name()); err != nil {
		os.Remove(tempFile.Name())
		return err
	}

	os.Remove(tempFile.Name())

	return nil
}

// WriteBytesToFilepath ...
func WriteBytesToFilepath(data []byte, targetFilepath string) error {
	targetFile, targetErr := os.Create(targetFilepath)
	if targetErr != nil {
		return targetErr
	}

	defer targetFile.Close()

	if _, err := targetFile.Write(data); err != nil {
		return err
	}

	return nil
}

// ExtractPathFromFilepath ...
func ExtractPathFromFilepath(filepath string) string {
	parts := strings.Split(filepath, string(os.PathSeparator))

	return strings.Join(parts[:len(parts)-1], string(os.PathSeparator))
}

// SynchronizePathFromFilename ...
func SynchronizePathFromFilename(filepath string) error {
	path := ExtractPathFromFilepath(filepath)

	if !utils.IsDirExists(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// FillPartitionByFilepath ...
func FillPartitionByFilepath(p *FilePartition, filepath string) {
	parts := strings.Split(filepath, string(os.PathSeparator))
	filename := parts[len(parts)-1]
	filenameParts := strings.Split(filename, ".")
	p.Name = filenameParts[0]
	p.Extension = filenameParts[1]
	p.Path = strings.Join(parts[:len(parts)-1], string(os.PathSeparator))
	p.Fullpath = GetFilespath() + string(os.PathSeparator) + p.Path
}
