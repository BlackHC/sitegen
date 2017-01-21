package action

import (
	"regexp"
	"strings"
)

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Maps source urls to final urls on the server
type Remapper interface {
	AddRemapping(localPath string, finalUrl string)
	MapUrl(sourceUrl string) string
}

func LinkHtmlReferences(remaper Remapper, postPath string, content string) string {
	changedUrls := map[string]string{}

	urlRegExp, err := regexp.Compile(`(href|src)="(.*?)"`)
	errPanic(err)

	urls := urlRegExp.FindAllStringSubmatch(content, -1)
	for _, matches := range urls {
		url := matches[2]
		if strings.HasPrefix(url, "#") {
			continue
		}
		// TODO check relative paths using the "current" directory
		targetUrl := remaper.MapUrl(url)
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
