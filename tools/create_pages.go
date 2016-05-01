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
	Title       string
	Posts       []PostContext
	Navigation  NavigationContext
	ArticleTree map[string][]string
	sitemap     *data.Sitemap
}

type ArticleNavContext struct {
	sitemap     *data.Sitemap
	articlePath string
}

func (articleNavContext ArticleNavContext) Title() string {
	return articleNavContext.sitemap.Metadata[articleNavContext.articlePath].Title
}

func (articleNavContext ArticleNavContext) Url() string {
	return articleNavContext.sitemap.Metadata[articleNavContext.articlePath].Url
}

func (articleNavContext ArticleNavContext) Children() []ArticleNavContext {
	childNodes := articleNavContext.sitemap.ArticleTree[articleNavContext.articlePath]
	children := make([]ArticleNavContext, len(childNodes))
	for i, subArticlePath := range childNodes {
		children[i] = ArticleNavContext{sitemap: articleNavContext.sitemap, articlePath: subArticlePath}
	}
	return children
}

func (siteContext SiteContext) BlogTitle() string { return data.BlogTitle }

func (siteContext SiteContext) BlogSubtitle() string { return data.BlogSubtitle }

func (siteContext SiteContext) GetRootArticles() []ArticleNavContext {
	rootArticle := ArticleNavContext{sitemap: siteContext.sitemap, articlePath: data.RootArticlePath}
	return rootArticle.Children()
}

func (siteContext SiteContext) ResolveUrl(entryPath string) string {
	return siteContext.sitemap.MapUrl(entryPath)
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
	pageFile := util.CreateOutputFile(outputPath)

	err := siteTemplate.Execute(pageFile, context)
	errPanic(err)

	fmt.Printf("%s created\n", outputPath)
}

func createArticle(sitemap *data.Sitemap, articlePath string) {
	articleMetadata := sitemap.Metadata[articlePath]

	context := SiteContext{
		Title:       articleMetadata.Title,
		Posts:       []PostContext{createPostContext(articleMetadata)},
		ArticleTree: sitemap.ArticleTree,
		Navigation: NavigationContext{NextPageUrl: nil,
			PreviousPageUrl: nil},
		sitemap: sitemap}

	// TODO: add navigation!!

	//fmt.Println(postTemplateContext)

	// Run the template
	outputPath := path.Join("http", articleMetadata.Url)
	executeSiteTemplate(outputPath, context)
}

func createPost(sitemap *data.Sitemap, postPath string) {
	postMetadata := sitemap.Metadata[postPath]

	context := SiteContext{
		Title:       postMetadata.Title,
		Posts:       []PostContext{createPostContext(postMetadata)},
		ArticleTree: sitemap.ArticleTree,
		Navigation: NavigationContext{NextPageUrl: sitemap.MaybePostUrl(sitemap.NextPostPath(postPath)),
			PreviousPageUrl: sitemap.MaybePostUrl(sitemap.PrevPostPath(postPath))},
		sitemap: sitemap}

	//fmt.Println(postTemplateContext)

	// Run the template
	outputPath := path.Join("http", postMetadata.Url)
	executeSiteTemplate(outputPath, context)
}

func createIndexPage(sitemap *data.Sitemap, index int) {
	page := sitemap.IndexPages[index]

	postContexts := make([]PostContext, 0, len(page.PostPaths))
	for _, postPath := range page.PostPaths {
		postContext := createPostContext(sitemap.Metadata[postPath])
		postContexts = append(postContexts, postContext)
	}

	context := SiteContext{
		Title:       page.Title,
		Posts:       postContexts,
		ArticleTree: sitemap.ArticleTree,
		Navigation: NavigationContext{
			NextPageUrl:     sitemap.IndexPages.MaybeNextIndexPage(index),
			PreviousPageUrl: sitemap.IndexPages.MaybePreviousIndexPage(index)},
		sitemap: sitemap}

	outputPath := "http" + page.Url
	executeSiteTemplate(outputPath, context)
}

func createSubArticles(sitemap *data.Sitemap, articlePath string) {
	for _, subArticlePath := range sitemap.ArticleTree[articlePath] {
		println(subArticlePath)
		createArticle(sitemap, subArticlePath)
		createSubArticles(sitemap, subArticlePath)
	}
}

func main() {
	flag.Parse()

	// Read in sitemap
	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)

	// Create a page for every post.
	for _, postPath := range sitemap.OrderedPosts {
		createPost(sitemap, postPath)
	}

	// Create a page for every article.
	createSubArticles(sitemap, data.RootArticlePath)

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
