package rss

import (
	"context"
	"log"
	"github.com/hanwen/go-fuse/v2/fs"
	"syscall"
)

type FSRSS struct {
	Name    string
	Content *Content
	fs.Inode
}

func (f *FSRSS) OnAdd(ctx context.Context) {
	for _, feedName := range f.Content.ListFeeds() {
		node := f.NewPersistentInode(ctx, &fsdir{
			FeedName: feedName,
			Content:  f.Content,
		}, fs.StableAttr{Mode: syscall.S_IFDIR})
		f.AddChild(feedName, node, true)
	}
}

type fsdir struct {
	FeedName    string
	Content     *Content
	InternalRep *Node // cached representation
	fs.Inode
}

func (f *fsdir) OnAdd(ctx context.Context) {
	// Load the feed data on-demand
	if f.InternalRep == nil {
		if f.FeedName != "" {
			// This is a feed root directory - load the feed
			node, err := f.Content.GetNode(f.FeedName)
			if err != nil {
				log.Printf("Failed to load feed %s: %v", f.FeedName, err)
				return
			}
			f.InternalRep = node
		} else {
			// This shouldn't happen in normal operation
			return
		}
	}
	
	for _, child := range f.InternalRep.Children {
		var embed fs.InodeEmbedder
		attr := fs.StableAttr{Mode: syscall.S_IFDIR}
		embed = &fsdir{
			InternalRep: child,
			Content:     f.Content,
		}
		if !child.Dir {
			embed = &fs.MemRegularFile{Data: []byte(child.Content)}
			attr = fs.StableAttr{}
		}
		inode := f.NewPersistentInode(ctx, embed, attr)
		f.AddChild(child.Name, inode, true)
	}
}
