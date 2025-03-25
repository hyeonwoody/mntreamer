package model

type FileInfo struct {
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	UpdatedAt   string `json:"updatedAt"`
}
