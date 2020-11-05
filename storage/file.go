package storage

import (
	"imagine2/files"
	"imagine2/models"
	"imagine2/utils"

	"github.com/doug-martin/goqu/v9"

	// Mysql dialect
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"

	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// GetFileByID - get file
func GetFileByID(id int) (*models.File, error) {
	file := &models.File{}

	err := DB.Get(file, `
		SELECT * FROM `+getFilesTable()+` 
		WHERE id = ? 
		LIMIT 1
	`, id)

	if err != nil {
		logrus.Warning(err)
	}

	return file, err
}

// GetFiles ...
func GetFiles(offset int, limit int) (*[]models.File, error) {
	return nil, nil
}

// GetFile - Retrieve file by path and name
func GetFile(path, name string) (*models.File, error) {
	file := &models.File{}

	err := DB.Get(file, `
		SELECT * FROM `+getFilesTable()+` 
		WHERE path = ? 
			AND name = ? 
		LIMIT 1
	`, path, name)

	if err != nil {
		logrus.Warning(err)
	}

	return file, err
}

// SaveFile ...
func SaveFile(file *models.File) error {
	return nil
}

// DeleteFile ...
func DeleteFile(id int) error {
	dialect := goqu.Dialect("mysql")
	sql, _, err := dialect.Delete(getFilesTable()).Where(goqu.Ex{"id": id}).ToSQL()

	if err != nil {
		return err
	}

	_, err = DB.Exec(sql)

	return err
}

// CreateFromPartition ...
func CreateFromPartition(p *files.FilePartition) (*models.File, error) {
	file := &models.File{}

	file.Name = p.Name
	file.Fullname = p.Name + "." + p.Extension
	file.SourceName = p.SourceName
	file.Path = p.Path
	file.Ext = p.Extension
	file.Mime = p.Mime
	file.Size = p.Size
	file.Height = p.ImageHeight
	file.Width = p.ImageWidth
	file.Tags = ""
	file.Created = utils.StorageTimestamp()
	file.Updated = utils.StorageTimestamp()

	if file.SourceName == "" {
		file.SourceName = file.Fullname
	}

	dialect := goqu.Dialect("mysql")
	sql, _, err := dialect.Insert(getFilesTable()).Rows(file).ToSQL()

	if err != nil {
		return nil, err
	}

	result, err := DB.Exec(sql)

	if err == nil {
		lastInsertID, err := result.LastInsertId()
		if err == nil {
			file.ID = int(lastInsertID)
		}
	}

	return file, err
}
