package storage

import (
	"os"
	"io"
	"io/ioutil"
	"errors"
)

type LocalStorageConfiguration struct {
	Path string `json:"path"`
}

type LocalStorage struct {
	Configuration *LocalStorageConfiguration
}

func NewLocalStorage(Configuration *LocalStorageConfiguration)(StorageHandle) {
	return StorageHandle(&LocalStorage{
		Configuration: Configuration,
	})
}

func(Storage *LocalStorage) HasFile(file string) bool {
	_, err := os.Open(Storage.Configuration.Path + file)
	return err == nil
}

func (Storage *LocalStorage) GetFile(file string) (io.ReadCloser, error) {
	f, err := os.Open(Storage.Configuration.Path + file)

	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(f), nil
}

func (Storage *LocalStorage) DeleteFile(file string) error {
	if !Storage.HasFile(file) {
		return errors.New("File doesn't exist")
	}

	return os.Remove(Storage.Configuration.Path + file)
}

func (Storage *LocalStorage) WriteFile(buffer []byte, file string) (error) {
	if Storage.HasFile(file) {
		return errors.New("File already exists")
	}

	f, err := os.Open(Storage.Configuration.Path + file)
	defer f.Close()

	if err != nil {
		return err
	}

	_, err = f.Write(buffer)

	if err != nil {
		return err
	}

	err = f.Sync()

	if err != nil {
		return err
	}

	return nil
}