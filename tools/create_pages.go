// Creates _page html files in html/ for now using a template stored in templates/
package main

import (
	"flag"
	"os"
	"strings"
	"text/template"

	"github.com/blackhc.github.io/generator/data"
	"github.com/blackhc.github.io/generator/util"
)

// Create pages from content and template
// For now read posts_metadata.json and use that for content
// and template/post.template as template
const postTemplatePath = "templates/post.template"

type NavigationContext struct {
	NextPost     *string
	PreviousPost *string
}

type PostTemplateContext struct {
	Title      string
	Content    string
	Date       string
	Navigation NavigationContext
}

var postPath string

func init() {
	flag.StringVar(&postPath, "post", "posts/2007-12-02-geckos-here-and-geckos-there-geckos-everywhere.markdown", "Post name")
}

func MdToPageHtml(mdName string) string {
	blogHtmlPath := "html/" + strings.TrimPrefix(strings.TrimSuffix(mdName, ".markdown"), "posts/") + "_page.html"
	return blogHtmlPath
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	// Read in sitemap
	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)

	postMetadata := sitemap.Posts[postPath]

	postContent, err := postMetadata.Content()
	errPanic(err)

	postTemplateContext := PostTemplateContext{
		Title:   postMetadata.Title,
		Content: postContent, Date: postMetadata.Date.Format("Mon Jan _2 2006 15:04:05"),
		Navigation: NavigationContext{NextPost: sitemap.MaybePostUrl(sitemap.NextPostPath(postPath)),
			PreviousPost: sitemap.MaybePostUrl(sitemap.PrevPostPath(postPath))}}

	//fmt.Println(postTemplateContext)
	// Create template from file
	postTemplate, err := template.ParseFiles(postTemplatePath)
	errPanic(err)

	// Set options (missingkey)
	postTemplate.Option("missingkey=error")

	// Run the template
	// Write output to _blog.html
	pageFile, err := os.Create(MdToPageHtml(postPath))
	errPanic(err)

	err = postTemplate.Execute(pageFile, postTemplateContext)
	errPanic(err)
}
