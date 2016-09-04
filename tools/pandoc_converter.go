// Uses pandoc to convert markdown files in pages/ and posts/ to html snippets
package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blackhc.github.io/generator/util"
)

const basePath string = "html/"

func ConvertDirectory(directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".markdown") {
			log.Println(path)
			standaloneOutputPath := basePath + strings.TrimSuffix(path, ".markdown") + ".html"
			if err = util.CreateParentDirectory(standaloneOutputPath); err != nil {
				return err
			}
			standaloneCmd := exec.Command("pandoc", "--mathjax", "-f", "markdown", "-t", "html5", "-s", path, "--no-highlight", "-o", standaloneOutputPath)
			if standardOutput, standaloneErr := standaloneCmd.CombinedOutput(); standaloneErr != nil {
				log.Fatal(string(standardOutput))
			}
			contentOutputPath := basePath + strings.TrimSuffix(path, ".markdown") + "_content.html"
			if err = util.CreateParentDirectory(contentOutputPath); err != nil {
				return err
			}
			contentCmd := exec.Command("pandoc", "--mathjax", "-f", "markdown", "-t", "html5", path, "--no-highlight", "-o", contentOutputPath)
			if contentOutput, contentError := contentCmd.CombinedOutput(); contentError != nil {
				log.Fatal(contentOutput)
			}
		}
		return err
	})
}

func main() {
	err := ConvertDirectory("pages")
	if err != nil {
		panic(err)
	}
	// TODO: this is broken atm. atm all the converted posts are in html/*
	err = ConvertDirectory("posts")
	if err != nil {
		panic(err)
	}
}
