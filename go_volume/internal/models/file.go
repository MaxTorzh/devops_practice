package models

import (
	"os"
	"time"
)

type FileInfo struct {
	Name string `json:"name:"`
	Size int64 `json:"size"`
	ModTime time.Time `json:"mod_time"`
	IsDir bool `json:"is_dir"`
	Path string `json:"path,omitempty"`
}

type FileContent struct {
	Name string `json:"name"`
	Content string `json:"content"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Message string `json:"message"`
	Status int `json:"status"`
}

func FromFileInfo(path string, info os.FileInfo) FileInfo {
	return FileInfo{
		Name: info.Name(),
		Size: info.Size(),
		ModTime: info.ModTime(),
		IsDir: info.IsDir(),
		Path: path,
	}
}