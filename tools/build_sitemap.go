package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/blackhc.github.io/generator/data"
	"github.com/blackhc.github.io/generator/util"
)

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func mdToContentHtml(path string) string {
	contentHtmlPath := "html/" + strings.TrimPrefix(strings.TrimSuffix(path, ".markdown"), "posts/") + "_content.html"
	return contentHtmlPath
}

func main() {
	// Iterate through all md files and determine the output file.
	postsMetadata := map[string]interface{}{}
	err := util.ImportJson("raw_posts_metadata.json", &postsMetadata)
	errPanic(err)

	sitemap := data.NewSitemap()

	for postPath, metadataUncasted := range postsMetadata {
		metadata := metadataUncasted.(map[string]interface{})

		slugString := metadata["slug"].(string)
		title := metadata["title"].(string)
		postDateString := metadata["date"].(string)
		// TODO: support title -> slug conversion for new posts
		postDate, err := time.Parse("2006-01-02 15:04:05+00:00", postDateString)
		errPanic(err)

		postUrl := fmt.Sprintf("%04d/%02d/%s/index.html", postDate.Year(), postDate.Month(), slugString)
		errPanic(err)

		contentPath := mdToContentHtml(postPath)
		postMetadata := data.Metadata{Title: title,
			Date: data.JSONTime{postDate}, Slug: slugString, Url: postUrl, ContentPath: contentPath}
		log.Println(postPath)
		sitemap.AddPost(postPath, postMetadata)
	}

	sitemap.OrderPosts()

	util.ExportJson("sitemap.json", sitemap)
}
