// Data model of the blogs sitemap
package data

import (
	"fmt"
	"net/url"
	"sort"
)

const BlogTitle = "BlackHC's Adventures"
const BlogSubtitle = "... in the Dev World"

const RootArticlePath = "pages/index.markdown"

type Sitemap struct {
	// Map of PostPath to Metadata
	Metadata     map[string]*Metadata
	OrderedPosts []string
	// Tree that starts in pages/index.markdown and contains the article paths of subnodes in order.
	ArticleTree map[string][]string
	IndexPages  IndexPages
	// Linker will replace links to key with value
	Remappings map[string]string
}

func (s *Sitemap) AddPost(postPath string, post Metadata) {
	s.Metadata[postPath] = &post
	s.OrderedPosts = append(s.OrderedPosts, postPath)
}

func (s *Sitemap) AddArticle(articlePath string, article Metadata) {
	// ArticleTree is managed separately.
	s.Metadata[articlePath] = &article
}

func (s *Sitemap) AddRemapping(localPath string, finalUrl string) {
	s.Remappings[localPath] = finalUrl
}

func NewSitemap() *Sitemap {
	return &Sitemap{
		Metadata:     map[string]*Metadata{},
		OrderedPosts: []string{},
		ArticleTree:  map[string][]string{},
		Remappings:   map[string]string{},
	}
}

func (s Sitemap) GetPostByIndex(i int) *Metadata {
	return s.Metadata[s.OrderedPosts[i]]
}

func (s *Sitemap) GetIndex(postPath string) int {
	for i, p := range s.OrderedPosts {
		if p == postPath {
			return i
		}
	}
	panic("Post not found")
}

func (s *Sitemap) NextPostPath(postPath string) *string {
	i := s.GetIndex(postPath)
	if i+1 < len(s.OrderedPosts) {
		return &s.OrderedPosts[i+1]
	}
	return nil
}

func (s *Sitemap) PrevPostPath(postPath string) *string {
	i := s.GetIndex(postPath)
	if i > 0 {
		return &s.OrderedPosts[i-1]
	}
	return nil
}

type ByDate Sitemap

func (a ByDate) Len() int { return len(a.OrderedPosts) }
func (a ByDate) Swap(i, j int) {
	a.OrderedPosts[i], a.OrderedPosts[j] = a.OrderedPosts[j], a.OrderedPosts[i]
}
func (a ByDate) Less(i, j int) bool {
	return Sitemap(a).GetPostByIndex(i).Date.After(Sitemap(a).GetPostByIndex(j).Date.Time)
}

func (s Sitemap) OrderPosts() {
	sort.Sort(ByDate(s))
}

// TODO: where does this go?
func (s *Sitemap) MaybePostUrl(postPath *string) *string {
	if postPath != nil {
		url := s.Metadata[*postPath].Url
		return &url
	}
	return nil
}

func (s Sitemap) MapUrl(sourceUrl string) string {
	if metadata, found := s.Metadata[sourceUrl]; found {
		return metadata.Url
	}
	if targetUrl, found := s.Remappings[sourceUrl]; found {
		return targetUrl
	}
	if parsedUrl, err := url.Parse(sourceUrl); err != nil || !parsedUrl.IsAbs() {
		fmt.Println(sourceUrl + " not found in remappings/posts!")
	}
	return sourceUrl
}
