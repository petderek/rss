package rss

import "strings"

// probably unnecesary - but it makes it easier for me to translate rss -> tree -> filesystem
// i'll eventually just do rss->filesystem
type Node struct {
	Name     string
	Content  string
	Dir      bool
	Children []*Node
}

func (n *Node) add(child *Node) {
	n.Children = append(n.Children, child)
}

func fromRss(r rss) *Node {
	root := &Node{
		Dir: true,
	}
	root.add(&Node{
		Name:    "title",
		Content: r.Channel.Title,
	})
	root.add(&Node{
		Name:    "description",
		Content: r.Channel.Description,
	})
	root.add(&Node{
		Name:    "link",
		Content: r.Channel.Link,
	})
	articles := &Node{
		Name: "articles",
		Dir:  true,
	}
	for _, article := range r.Channel.Items {
		n := &Node{
			Name: sanitizeFilepath(article.Title),
			Dir:  true,
		}
		n.add(&Node{
			Name:    "title",
			Content: article.Title,
		})
		n.add(&Node{
			Name:    "description",
			Content: article.Description,
		})
		n.add(&Node{
			Name:    "link",
			Content: article.Link,
		})
		n.add(&Node{
			Name:    "guid",
			Content: article.Guid,
		})
		n.add(&Node{
			Name:    "date",
			Content: article.PublishDate,
		})
		n.add(&Node{
			Name:    "content",
			Content: article.Content,
		})
		articles.add(n)
	}
	root.add(articles)

	return root
}

func sanitizeFilepath(in string) string {
	removedSpaces := strings.Replace(in, " ", "_", -1)
	removedQuestions := strings.Replace(removedSpaces, "?", "", -1)
	return removedQuestions
}
