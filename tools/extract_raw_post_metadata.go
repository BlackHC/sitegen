///bin/true; exec /usr/bin/env go run "$0" "$@"

// Extracts raw post metadata
package main

import (
	"io/ioutil"

	"github.com/blackhc.github.io/generator/util"

	"gopkg.in/yaml.v2"
)

func main() {
	posts, err := ioutil.ReadDir("posts")
	if err != nil {
		panic(err)
	}

	postMetadata := map[string]interface{}{}

	for _, post := range posts {
		postFilename := "posts/" + post.Name()
		byteContent, _ := ioutil.ReadFile(postFilename)

		m := map[string]interface{}{}
		if err := yaml.Unmarshal(byteContent, &m); err != nil {
			panic(err)
		}
		postMetadata[postFilename] = m
	}

	util.ExportJson("raw_posts_metadata.json", postMetadata)
}
