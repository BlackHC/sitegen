package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func enumeratePosts(sitemap *data.Sitemap) {
	postsMetadata := map[string]interface{}{}
	err := util.ImportJson("raw_posts_metadata.json", &postsMetadata)
	errPanic(err)
	for postPath, metadataUncasted := range postsMetadata {
		metadata := metadataUncasted.(map[string]interface{})

		slugString := metadata["slug"].(string)
		title := metadata["title"].(string)
		postDateString := metadata["date"].(string)
		// TODO: support title -> slug conversion for new posts
		postDate, err := time.Parse("2006-01-02 15:04:05+00:00", postDateString)
		errPanic(err)

		postUrl := fmt.Sprintf("/%04d/%02d/%s/index.html", postDate.Year(), postDate.Month(), slugString)
		errPanic(err)

		contentPath := mdToContentHtml(postPath)
		postMetadata := data.Metadata{Title: title,
			Date: data.JSONTime{postDate}, Slug: slugString, Url: postUrl, ContentPath: contentPath}
		log.Println(postPath)
		sitemap.AddPost(postPath, postMetadata)
	}
	sitemap.OrderPosts()
}

func enumerateImages(sitemap *data.Sitemap) {
	filepath.Walk("html/images", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// TODO: move images out of html
			localPath := strings.TrimPrefix(path, "html/")
			sitemap.AddRemapping(localPath, "/"+localPath)
		}
		return err
	})
}

func main() {
	// Iterate through all md files and determine the output file.

	sitemap := data.NewSitemap()
	enumeratePosts(sitemap)
	enumerateImages(sitemap)

	util.ExportJson("sitemap.json", sitemap)
}
