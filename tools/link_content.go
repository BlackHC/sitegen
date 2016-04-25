package main

import (
	"github.com/blackhc.github.io/generator/action"
	"github.com/blackhc.github.io/generator/data"
	"github.com/blackhc.github.io/generator/util"
)

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Returns the filename of the linked content
func linkPost(sitemap *data.Sitemap, postPath string) {
	postMetadata := sitemap.Posts[postPath]
	pandocContent := string(util.ReadFile(postMetadata.PandocPath))
	linkedContent := action.LinkHtmlReferences(sitemap, postPath, pandocContent)
	util.WriteFile(postMetadata.ContentPath, []byte(linkedContent))
}

func main() {
	// Iterate through all md files and determine the output file.

	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)
	for postPath, _ := range sitemap.Posts {
		linkPost(sitemap, postPath)
	}
}
