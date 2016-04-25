// Creates _page html files in html/ for now using a template stored in templates/
package main

import (
	"flag"
	"fmt"
	"path"
	"text/template"

	"github.com/blackhc.github.io/generator/data"
	"github.com/blackhc.github.io/generator/util"
)

// Create pages from content and template
// For now read posts_metadata.json and use that for content
// and template/post.template as template
const siteTemplatePath = "templates/site.template"

var siteTemplate = loadSiteTemplate()

type NavigationContext struct {
	NextPageUrl     *string
	PreviousPageUrl *string
}

type PostContext struct {
	Title   string
	Date    string
	Content string
	Url     string
}

type SiteContext struct {
	BlogTitle    string
	BlogSubtitle string
	Title        string
	Posts        []PostContext
	Navigation   NavigationContext
}

func createPostContext(postMetadata *data.Metadata) PostContext {
	postContent, err := postMetadata.Content()
	errPanic(err)

	return PostContext{Title: postMetadata.Title, Date: postMetadata.Date.Format("Mon Jan _2 2006 15:04:05"),
		Content: postContent, Url: postMetadata.Url}
}

func loadSiteTemplate() *template.Template {
	// Create template from file
	siteTemplate, err := template.ParseFiles(siteTemplatePath)
	errPanic(err)

	// Set options (missingkey)
	siteTemplate.Option("missingkey=error")
	return siteTemplate
}

func executeSiteTemplate(outputPath string, context interface{}) {
	pageFile, err := util.CreateOutputFile(outputPath)
	errPanic(err)

	err = siteTemplate.Execute(pageFile, context)
	errPanic(err)

	fmt.Printf("%s created\n", outputPath)
}

func createPost(sitemap *data.Sitemap, postPath string) {
	postMetadata := sitemap.Posts[postPath]

	context := SiteContext{
		BlogTitle:    data.BlogTitle,
		BlogSubtitle: data.BlogSubtitle,
		Title:        postMetadata.Title,
		Posts:        []PostContext{createPostContext(postMetadata)},
		Navigation: NavigationContext{NextPageUrl: sitemap.MaybePostUrl(sitemap.NextPostPath(postPath)),
			PreviousPageUrl: sitemap.MaybePostUrl(sitemap.PrevPostPath(postPath))}}

	//fmt.Println(postTemplateContext)

	// Run the template
	outputPath := path.Join("http", postMetadata.Url)
	executeSiteTemplate(outputPath, context)
}

func createIndexPage(sitemap *data.Sitemap, index int) {
	page := sitemap.IndexPages[index]

	postContexts := make([]PostContext, 0, len(page.PostPaths))
	for _, postPath := range page.PostPaths {
		postContext := createPostContext(sitemap.Posts[postPath])
		postContexts = append(postContexts, postContext)
	}

	context := SiteContext{
		BlogTitle:    data.BlogTitle,
		BlogSubtitle: data.BlogSubtitle,
		Title:        page.Title,
		Posts:        postContexts,
		Navigation: NavigationContext{
			NextPageUrl:     sitemap.IndexPages.MaybeNextIndexPage(index),
			PreviousPageUrl: sitemap.IndexPages.MaybePreviousIndexPage(index)}}

	outputPath := "http" + page.Url
	executeSiteTemplate(outputPath, context)
}

func main() {
	flag.Parse()

	// Read in sitemap
	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)

	// Create a page for every post.
	for postPath, _ := range sitemap.Posts {
		createPost(sitemap, postPath)
	}

	// Create summary pages.
	indexPages := sitemap.IndexPages
	for index := 0; index < len(indexPages); index++ {
		createIndexPage(sitemap, index)
	}
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
