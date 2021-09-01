package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type IStorageInterface interface {
	Save(f *FileForStorage) (string, error)
	Read(filename string) ([]byte, error)
	Stats(filename string) (os.FileInfo, error)
	Exists(filename string) (bool, error)
	Delete(filename string) (bool, error)
	GetUploadURL() string
}

type FileForStorage struct {
	Content           []byte
	PatternForTheFile string
	Filename          string
}
type FsStorage struct {
	UploadPath string
	URLPath    string
}

func (s *FsStorage) GetUploadURL() string {
	return s.URLPath
}

func (s *FsStorage) Save(f *FileForStorage) (string, error) {
	tmpfile, err := ioutil.TempFile(s.UploadPath, f.PatternForTheFile)
	if err != nil {
		return "", err
	}
	_, err = tmpfile.Write(f.Content)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	return strings.Replace(tmpfile.Name(), s.UploadPath, "", 1), nil
}

func (s *FsStorage) Read(filename string) ([]byte, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s", s.UploadPath, filename))
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (s *FsStorage) Stats(filename string) (os.FileInfo, error) {
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", s.UploadPath, filename), os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (s *FsStorage) Exists(filename string) (bool, error) {
	filepath := fmt.Sprintf("%s/%s", s.UploadPath, filename)
	var err error
	if _, err = os.Stat(filepath); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func (s *FsStorage) Delete(filename string) (bool, error) {
	filepath := fmt.Sprintf("%s/%s", s.UploadPath, filename)
	err := os.Remove(filepath)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewFsStorage() IStorageInterface {
	return &FsStorage{UploadPath: CurrentConfig.GetPathToUploadDirectory(), URLPath: CurrentConfig.GetURLToUploadDirectory()}
}
