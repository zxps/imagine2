package models

import (
	"database/sql"
	"encoding/json"
)

// File - file model
type File struct {
	ID           int       `db:"id" json:"id" goqu:"skipinsert"`
	Fullname     string    `db:"fullname" json:"fullname"`
	Name         string    `db:"name" json:"name"`
	SourceName   string    `db:"source_name" json:"source_name"`
	Ext          string    `db:"ext" json:"ext"`
	Path         string    `db:"path" json:"path"`
	Size         int64     `db:"size" json:"size"`
	Width        int       `db:"width" json:"width"`
	Height       int       `db:"height" json:"height"`
	Mime         string    `db:"mime" json:"mime"`
	Updated      string    `db:"updated" json:"updated"`
	Created      string    `db:"created" json:"created"`
	SourceFileID NullInt64 `db:"source_file_id" json:"source_file_id" goqu:"defaultifempty"`
	Tags         string    `db:"tags" json:"tags"`
	Source       string    `db:"source" json:"source"`
	SourceLabel  string    `db:"source_label" json:"source_label"`
}

// NullInt64 struct
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for null string
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid && v.Int64 > 0 {
		return json.Marshal(v.Int64)
	}

	return json.Marshal(nil)
}
