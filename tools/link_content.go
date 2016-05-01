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

func linkContent(sitemap *data.Sitemap, path string, metadata *data.Metadata) {
	pandocContent := string(util.ReadFile(metadata.PandocPath))
	linkedContent := action.LinkHtmlReferences(sitemap, path, pandocContent)
	util.WriteFile(metadata.ContentPath, []byte(linkedContent))
}

func main() {
	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)
	for path, metadata := range sitemap.Metadata {
		linkContent(sitemap, path, metadata)
	}
}
