package rss

import "strings"

// probably unnecesary - but it makes it easier for me to translate Rss -> tree -> filesystem
// i'll eventually just do Rss->filesystem
type Node struct {
	Name     string
	Content  string
	Dir      bool
	Children []*Node
}

func (n *Node) add(child *Node) {
	n.Children = append(n.Children, child)
}

func FromRss(r Rss) *Node {
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

func FromAtom(a AtomFeed) *Node {
	root := &Node{
		Dir: true,
	}
	root.add(&Node{
		Name:    "title",
		Content: a.Title,
	})
	root.add(&Node{
		Name:    "description",
		Content: a.Subtitle,
	})
	
	var link string
	for _, l := range a.Link {
		if l.Rel == "alternate" || l.Rel == "" {
			link = l.Href
			break
		}
	}
	root.add(&Node{
		Name:    "link",
		Content: link,
	})
	
	articles := &Node{
		Name: "articles",
		Dir:  true,
	}
	for _, entry := range a.Entries {
		n := &Node{
			Name: sanitizeFilepath(entry.Title),
			Dir:  true,
		}
		n.add(&Node{
			Name:    "title",
			Content: entry.Title,
		})
		n.add(&Node{
			Name:    "description",
			Content: entry.Summary,
		})
		
		var entryLink string
		for _, l := range entry.Link {
			if l.Rel == "alternate" || l.Rel == "" {
				entryLink = l.Href
				break
			}
		}
		n.add(&Node{
			Name:    "link",
			Content: entryLink,
		})
		n.add(&Node{
			Name:    "guid",
			Content: entry.ID,
		})
		n.add(&Node{
			Name:    "date",
			Content: entry.Published,
		})
		
		content := entry.Content.Text
		if content == "" {
			content = entry.Summary
		}
		n.add(&Node{
			Name:    "content",
			Content: content,
		})
		articles.add(n)
	}
	root.add(articles)

	return root
}

func sanitizeFilepath(in string) string {
	removedSpaces := strings.Replace(in, " ", "_", -1)
	removedQuotes := strings.Replace(removedSpaces, "'", "", -1)
	removedDoubleQuotes := strings.Replace(removedQuotes, "\"", "", -1)
	removedQuestions := strings.Replace(removedDoubleQuotes, "?", "", -1)
	return removedQuestions
}
