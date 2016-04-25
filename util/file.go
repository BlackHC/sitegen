package util

import (
	"io/ioutil"
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

func ReadFile(filepath string) []byte {
	data, err := ioutil.ReadFile(filepath)
	errPanic(err)
	return data
}

func WriteFile(filepath string, data []byte) {
	file, err := CreateOutputFile(filepath)
	if err == nil {
		_, err = file.Write(data)
	}
	errPanic(err)
}
