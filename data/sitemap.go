// Data model of the blogs sitemap
package data

import (
	"fmt"
	"net/url"
	"sort"
)

// Relative path to the original markdown file (used as index for everything)
type PostPath string

type Sitemap struct {
	// Ordered list of posts.
	Posts        map[string]Metadata
	OrderedPosts []string
	// Linker will replace links to key with value
	Remappings map[string]string
}

func (s *Sitemap) AddPost(postPath string, post Metadata) {
	s.Posts[postPath] = post
	s.OrderedPosts = append(s.OrderedPosts, postPath)
}

func (s *Sitemap) AddRemapping(localPath string, finalUrl string) {
	s.Remappings[localPath] = finalUrl
}

func NewSitemap() *Sitemap {
	return &Sitemap{Posts: map[string]Metadata{}, OrderedPosts: []string{}, Remappings: map[string]string{}}
}

func (s Sitemap) GetPostByIndex(i int) Metadata {
	return s.Posts[s.OrderedPosts[i]]
}

func (s Sitemap) GetIndex(postPath string) int {
	for i, p := range s.OrderedPosts {
		if p == postPath {
			return i
		}
	}
	panic("Post not found")
}

func (s Sitemap) NextPostPath(postPath string) *string {
	i := s.GetIndex(postPath)
	if i+1 < len(s.OrderedPosts) {
		return &s.OrderedPosts[i+1]
	}
	return nil
}

func (s Sitemap) PrevPostPath(postPath string) *string {
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
	return Sitemap(a).GetPostByIndex(i).Date.Before(Sitemap(a).GetPostByIndex(j).Date.Time)
}

func (s Sitemap) OrderPosts() {
	sort.Sort(ByDate(s))
}

// TODO: where does this go?
func (s *Sitemap) MaybePostUrl(postPath *string) *string {
	if postPath != nil {
		url := s.Posts[*postPath].Url
		return &url
	}
	return nil
}

func (s Sitemap) MapUrl(sourceUrl string) string {
	if postMetadata, found := s.Posts[sourceUrl]; found {
		return postMetadata.Url
	}
	if targetUrl, found := s.Remappings[sourceUrl]; found {
		return targetUrl
	}
	if parsedUrl, err := url.Parse(sourceUrl); err != nil && !parsedUrl.IsAbs() {
		fmt.Println(sourceUrl + " not found in remappings/posts!")
	}
	return sourceUrl
}
