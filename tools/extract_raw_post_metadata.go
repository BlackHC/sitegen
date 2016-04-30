///bin/true; exec /usr/bin/env go run "$0" "$@"

// Extracts raw post metadata
package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/blackhc.github.io/generator/util"

	"gopkg.in/yaml.v2"
)

func extractMetadata(filename string) interface{} {
	byteContent, _ := ioutil.ReadFile(filename)

	m := map[string]interface{}{}
	if err := yaml.Unmarshal(byteContent, &m); err != nil {
		panic(err)
	}
	return m
}

func main() {
	posts, err := ioutil.ReadDir("posts")
	if err != nil {
		panic(err)
	}

	postMetadata := map[string]interface{}{}

	for _, post := range posts {
		postFilename := "posts/" + post.Name()
		postMetadata[postFilename] = extractMetadata(postFilename)
	}

	util.ExportJson("raw_posts_metadata.json", postMetadata)

	// Now walk all pages
	pageMetadata := map[string]interface{}{}
	err = filepath.Walk("pages", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".markdown") {
			pageMetadata[path] = extractMetadata(path)
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	util.ExportJson("raw_pages_metadata.json", pageMetadata)
}
