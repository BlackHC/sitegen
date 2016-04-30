package util

import (
	"io/ioutil"
	"os"
	"path"
)

func CreateParentDirectory(filepath string) error {
	dir := path.Dir(filepath)
	return os.MkdirAll(dir, 0755)
}

func CreateOutputFile(filepath string) *os.File {
	if err := CreateParentDirectory(filepath); err != nil {
		errPanic(err)
	}

	if file, err := os.Create(filepath); err != nil {
		errPanic(err)
	}
	return file
}

func ReadFile(filepath string) []byte {
	data, err := ioutil.ReadFile(filepath)
	errPanic(err)
	return data
}

func WriteFile(filepath string, data []byte) {
	file := CreateOutputFile(filepath)
	_, err = file.Write(data)
	errPanic(err)
}
