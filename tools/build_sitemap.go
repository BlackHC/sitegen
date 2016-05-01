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

// Output from pandoc
func mdToContentHtml(path string) string {
	contentHtmlPath := "html/" + strings.TrimPrefix(strings.TrimSuffix(path, ".markdown"), "posts/") + "_content.html"
	return contentHtmlPath
}

func mdToLinkedContentHtml(path string) string {
	contentHtmlPath := "html/" + strings.TrimPrefix(strings.TrimSuffix(path, ".markdown"), "posts/") + "_linked.html"
	return contentHtmlPath
}

func createWordPressDisqusIdentifier(wordPressId int, layout string) string {
	if layout == "page" {
		return fmt.Sprintf("%d http://blog.blackhc.net/?page_id=%d", wordPressId, wordPressId)
	} else if layout == "post" {
		return fmt.Sprintf("%d http://blog.blackhc.net/?p=%d", wordPressId, wordPressId)
	}
	panic("Unknown layout " + layout)
}

func createWordPressDisqusPageUrl(url string, layout string) string {
	if layout == "page" {
		return "http://blog.blackhc.net" + strings.TrimSuffix(url, "index.html")
	} else if layout == "post" {
		return "http://blog.blackhc.net" + strings.TrimSuffix(url, "index.html")
	}
	panic("Unknown layout " + layout)
}

func enumerateArticles(sitemap *data.Sitemap) {
	articlesMetadata := map[string]interface{}{}
	util.ImportJson("raw_pages_metadata.json", &articlesMetadata)

	for articlePath, metadataUncasted := range articlesMetadata {
		metadata := metadataUncasted.(map[string]interface{})

		slugString := metadata["slug"].(string)
		title := metadata["title"].(string)
		postDateString := metadata["date"].(string)
		// TODO: support title -> slug conversion for new posts
		postDate, err := time.Parse("2006-01-02 15:04:05+00:00", postDateString)
		errPanic(err)

		articleUrl := strings.TrimSuffix(strings.TrimPrefix(articlePath, "pages"), ".markdown") + ".html"
		errPanic(err)

		layout := metadata["layout"].(string)
		wordPressId := int(metadata["wordpress_id"].(float64))
		disqusId := createWordPressDisqusIdentifier(wordPressId, layout)
		disqusPageUrl := createWordPressDisqusPageUrl(articleUrl, layout)

		pandocContentPath := mdToContentHtml(articlePath)
		linkedContentPath := mdToLinkedContentHtml(articlePath)
		articleMetadata := data.Metadata{Title: title,
			DisqusId:      disqusId,
			DisqusPageUrl: disqusPageUrl,
			Date:          data.JSONTime{postDate},
			Slug:          slugString,
			Url:           articleUrl,
			PandocPath:    pandocContentPath,
			ContentPath:   linkedContentPath}
		log.Println(articlePath)
		sitemap.AddArticle(articlePath, articleMetadata)
	}
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

		layout := metadata["layout"].(string)
		wordPressId := int(metadata["wordpress_id"].(float64))
		disqusId := createWordPressDisqusIdentifier(wordPressId, layout)
		disqusPageUrl := createWordPressDisqusPageUrl(postUrl, layout)

		pandocContentPath := mdToContentHtml(postPath)
		linkedContentPath := mdToLinkedContentHtml(postPath)
		postMetadata := data.Metadata{Title: title,
			DisqusId:      disqusId,
			DisqusPageUrl: disqusPageUrl,
			Date:          data.JSONTime{postDate},
			Slug:          slugString,
			Url:           postUrl,
			PandocPath:    pandocContentPath,
			ContentPath:   linkedContentPath}
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
		indexPage.Title = buildIndexPageTitle(pageIndex)
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
	title := data.BlogTitle
	if index > 0 {
		title = fmt.Sprintf("%s | %d", data.BlogTitle, index)
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
	util.ImportJson("article_tree.json", &sitemap.ArticleTree)
	enumeratePosts(sitemap)
	enumerateArticles(sitemap)
	enumerateImages(sitemap)

	util.ExportJson("sitemap.json", sitemap)
}
