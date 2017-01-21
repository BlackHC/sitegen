// Creates an article_tree.json that creates a tree structure of articles.
// The root is stored in the pages/index.markdown.
package main

import (
	"os"
	"path"

	"github.com/blackhc.github.io/generator/data"
	"github.com/blackhc.github.io/generator/util"
)

type ArticleTree map[string][]string

func main() {
	articleTree := ArticleTree{}
	if _, err := os.Stat("article_tree.json"); !os.IsNotExist(err) {
		util.ImportJson("article_tree.json", &articleTree)
	}

	articleMetadata := map[string]interface{}{}
	util.ImportJson("raw_pages_metadata.json", &articleMetadata)
	for articlePath, _ := range articleMetadata {
		articleTree.registerWithParent(articlePath)
	}

	util.ExportJson("article_tree.json", articleTree)
}

func (articleTree *ArticleTree) registerWithParent(nodePath string) {
	if nodePath == data.RootArticlePath {
		return
	}

	parentDir := path.Dir(path.Dir(nodePath))
	parentPath := parentDir + "/index.markdown"
	articleNode, contained := (*(*map[string][]string)(articleTree))[parentPath]
	if contained {
		var foundArticle = false
		for _, subArticle := range articleNode {
			if subArticle == nodePath {
				foundArticle = true
				break
			}
		}
		if !foundArticle {
			(*(*map[string][]string)(articleTree))[parentPath] = append((*(*map[string][]string)(articleTree))[parentPath], nodePath)
		}
	} else {
		(*(*map[string][]string)(articleTree))[parentPath] = []string{nodePath}
		articleTree.registerWithParent(parentPath)
	}
}
