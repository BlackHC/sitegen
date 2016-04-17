package util

import (
	"os"
	"path"
)

func CreateOutputFile(filepath string) (*os.File, error) {
	dir := path.Dir(filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	file, err := os.Create(filepath)
	return file, err
}

func WriteFile(filepath string, data []byte) error {
	file, err := CreateOutputFile(filepath)
	if err == nil {
		_, err = file.Write(data)
	}
	return err
}
