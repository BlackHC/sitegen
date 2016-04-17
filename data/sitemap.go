// Data model of the blogs sitemap
package data

import "sort"

// Relative path to the original markdown file (used as index for everything)
type PostPath string

type Sitemap struct {
	// Ordered list of posts.
	Posts        map[string]Metadata
	OrderedPosts []string
}

func (s *Sitemap) AddPost(postPath string, post Metadata) {
	s.Posts[postPath] = post
	s.OrderedPosts = append(s.OrderedPosts, postPath)
}

func NewSitemap() *Sitemap {
	return &Sitemap{Posts: map[string]Metadata{}, OrderedPosts: []string{}}
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
