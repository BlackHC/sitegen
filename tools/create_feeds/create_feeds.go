// Creates _page html files in html/ for now using a template stored in templates/
package main

import (
	"flag"
	"path"
	"time"

	"github.com/blackhc.github.io/generator/data"
	"github.com/blackhc.github.io/generator/util"
	"github.com/gorilla/feeds"
)

func createPostItem(sitemap *data.Sitemap, postPath string) *feeds.Item{
	postMetadata := sitemap.Metadata[postPath]

	textContent := util.ReadFile(postMetadata.ContentPath)
	return &feeds.Item{
		Title: postMetadata.Title,
		Link: &feeds.Link{Href: data.BlogDomainUrl + postMetadata.Url},
		Description: string(textContent),
		Author: &feeds.Author{Name: "Andreas 'BlackHC' Kirsch"},
		Created: postMetadata.Date.Time,
	}
}

func createFeed(sitemap *data.Sitemap) *feeds.Feed {
	feed := &feeds.Feed{
		Title: data.BlogTitle + " " + data.BlogSubtitle,
		Link: &feeds.Link{Href: data.BlogDomainUrl},
		Description: "",
		Author: &feeds.Author{Name: "Andreas 'BlackHC' Kirsch"},
		Created: time.Now(),
	}

	// Create a page for every post.
	feed.Items = []*feeds.Item{}
	for _, postPath := range sitemap.OrderedPosts {
		feed.Items = append(feed.Items, createPostItem(sitemap, postPath))
	}

	return feed
}



func main() {
	flag.Parse()

	// Read in sitemap
	sitemap := data.NewSitemap()
	util.ImportJson("sitemap.json", &sitemap)

	feed := createFeed(sitemap)

	rssOutputPath := path.Join("http/", "feed")
	rssContent, err := feed.ToRss()
	errPanic(err)
	util.WriteFile(rssOutputPath, []byte(rssContent))

	atomOutputPath := path.Join("http/", "atom.xml")
	atomContent, err := feed.ToAtom()
	errPanic(err)
	util.WriteFile(atomOutputPath, []byte(atomContent))
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
