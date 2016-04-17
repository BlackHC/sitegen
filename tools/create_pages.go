// Creates _page html files in html/ for now using a template stored in templates/
package main

import (
	"flag"
	"fmt"
	"path"
	"regexp"
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

func linkContent(s *data.Sitemap, postPath string, content string) string {
	changedUrls := map[string]string{}

	urlRegExp, err := regexp.Compile(`(href|src)=\"(.*?)\"`)
	errPanic(err)

	urls := urlRegExp.FindAllStringSubmatch(content, -1)
	for _, matches := range urls {
		url := matches[2]
		// TODO check relative paths using the "current" directory
		targetUrl := s.MapUrl(url)
		if url != targetUrl {
			changedUrls[url] = targetUrl
			content = strings.Replace(content, matches[0], strings.Replace(matches[0], url, targetUrl, -1), -1)
		}
	}
	// if len(changedUrls) > 0 {
	// 	fmt.Printf("%v has been changed during linking:\n\t%+v\n\n", postPath, changedUrls)
	// }
	return content
}

func main() {
	flag.Parse()

	// Read in sitemap
	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)

	for postPath, postMetadata := range sitemap.Posts {
		postContent, err := postMetadata.Content()
		errPanic(err)

		linkedContent := linkContent(sitemap, postPath, postContent)

		postTemplateContext := PostTemplateContext{
			Title:   postMetadata.Title,
			Content: linkedContent, Date: postMetadata.Date.Format("Mon Jan _2 2006 15:04:05"),
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

		pageFile, err := util.CreateOutputFile(outputPath)
		errPanic(err)

		err = postTemplate.Execute(pageFile, postTemplateContext)
		errPanic(err)

		fmt.Printf("%s created\n", outputPath)
	}
}
