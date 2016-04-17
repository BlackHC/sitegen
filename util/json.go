package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

func ImportJson(filePath string, data interface{}) error {
	byteData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteData, data)
	return err
}

func ExportJson(filePath string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, jsonData, "", "  ")

	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = out.WriteTo(outputFile)
	return err
}
