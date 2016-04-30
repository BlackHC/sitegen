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
	linkContent(sitemap, postPath, postMetadata)
}

func linkArticle(sitemap *data.Sitemap, articlePath string) {
	articleMetadata := sitemap.Articles[articlePath]
	linkContent(sitemap, articlePath, articleMetadata)
}

func linkContent(sitemap *data.Sitemap, path string, metadata *data.Metadata) {
	pandocContent := string(util.ReadFile(metadata.PandocPath))
	linkedContent := action.LinkHtmlReferences(sitemap, path, pandocContent)
	util.WriteFile(metadata.ContentPath, []byte(linkedContent))
}

func main() {
	// Iterate through all md files and determine the output file.

	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)
	for postPath, _ := range sitemap.Posts {
		linkPost(sitemap, postPath)
	}
	for articlePath, _ := range sitemap.Articles {
		linkArticle(sitemap, articlePath)
	}
}
