package data

type IndexPages []*IndexPage

type IndexPage struct {
	Title     string
	Url       string
	PostPaths []string
}

func (info IndexPages) MaybePreviousIndexPage(index int) *string {
	if index > 0 {
		url := info[index-1].Url
		return &url
	}
	return nil
}

func (info IndexPages) MaybeNextIndexPage(index int) *string {
	if index+1 < len(info) {
		url := info[index+1].Url
		return &url
	}
	return nil
}
