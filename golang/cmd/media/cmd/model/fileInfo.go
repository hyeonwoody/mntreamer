package model

import "time"

type FileInfo struct {
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
	Path        string `json:"path"`
	ModTime     time.Time
	Size        int64  `json:"size"`
	UpdatedAt   string `json:"updatedAt"`
}
