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

const blogTitle = "BlackHC's Adventures in the Dev World"
const blogSubtitle = "Just another weblog"

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Output from pandoc
func mdToContentHtml(path string) string {
	contentHtmlPath := "html/" + strings.TrimPrefix(strings.TrimSuffix(path, ".markdown"), "posts/") + "_content.html"
	return contentHtmlPath
}

func mdToLinkedContentHtml(path string) string {
	contentHtmlPath := "html/" + strings.TrimPrefix(strings.TrimSuffix(path, ".markdown"), "posts/") + "_linked.html"
	return contentHtmlPath
}

func enumeratePosts(sitemap *data.Sitemap) {
	postsMetadata := map[string]interface{}{}
	util.ImportJson("raw_posts_metadata.json", &postsMetadata)

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

		pandocContentPath := mdToContentHtml(postPath)
		linkedContentPath := mdToLinkedContentHtml(postPath)
		postMetadata := data.Metadata{Title: title,
			Date: data.JSONTime{postDate}, Slug: slugString, Url: postUrl, PandocPath: pandocContentPath, ContentPath: linkedContentPath}
		log.Println(postPath)
		sitemap.AddPost(postPath, postMetadata)
	}
	sitemap.OrderPosts()
	sitemap.IndexPages = buildIndexPages(sitemap)
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

// TODO move this into its own file
func buildIndexPages(sitemap *data.Sitemap) data.IndexPages {
	postCount := len(sitemap.OrderedPosts)
	const postsPerPage = 10
	pageCount := (postCount + postsPerPage - 1) / postsPerPage
	if pageCount == 0 {
		pageCount = 1
	}
	indexPages := make([]*data.IndexPage, pageCount)
	for pageIndex := 0; pageIndex < pageCount; pageIndex++ {
		indexPage := data.IndexPage{}
		indexPage.Url = buildIndexPageUrl(pageIndex)
		postStartIndex := pageIndex * postsPerPage
		postEndIndex := postStartIndex + postsPerPage
		if postEndIndex >= postCount {
			postEndIndex = postCount
		}
		indexPage.PostPaths = sitemap.OrderedPosts[postStartIndex:postEndIndex]
		indexPages[pageIndex] = &indexPage
	}
	return indexPages
}

func buildIndexPageTitle(index int) string {
	title := blogTitle
	if index > 0 {
		title = fmt.Sprintf("%s | %d", blogTitle, index)
	}
	return title
}

func buildIndexPageUrl(index int) string {
	if index == 0 {
		return "/index.html"
	} else {
		return fmt.Sprintf("/pages/%d/index.html", index)
	}
}

func main() {
	// Iterate through all md files and determine the output file.

	sitemap := data.NewSitemap()
	enumeratePosts(sitemap)
	enumerateImages(sitemap)

	util.ExportJson("sitemap.json", sitemap)
}
