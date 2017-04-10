package storage

import (
	"io"
)

const LOCAL_STORAGE = "local"
const CLOUD_STORAGE = "cloud"

type StorageHandle interface {
	HasFile(file string) bool
	GetFile(file string) (io.ReadCloser, error)
	DeleteFile(file string) error
	WriteFile(buffer []byte, file string) error
}