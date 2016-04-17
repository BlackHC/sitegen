// Creates _page html files in html/ for now using a template stored in templates/
package main

import (
	"flag"
	"os"
	"path"
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

func MdToPageHtml(mdName string) string {
	blogHtmlPath := "html/" + strings.TrimPrefix(strings.TrimSuffix(mdName, ".markdown"), "posts/") + "_page.html"
	return blogHtmlPath
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func createOutputFile(filepath string) (*os.File, error) {
	dir := path.Dir(filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	file, err := os.Create(filepath)
	return file, err
}

func main() {
	flag.Parse()

	// Read in sitemap
	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)

	for postPath, postMetadata := range sitemap.Posts {
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
		// Write output to _page.html
		outputPath := path.Join("http", postMetadata.Url)
		println(outputPath)
		pageFile, err := createOutputFile(outputPath)
		errPanic(err)

		err = postTemplate.Execute(pageFile, postTemplateContext)
		errPanic(err)
	}
}
