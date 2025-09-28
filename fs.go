package rss

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"syscall"
)

type FSRSS struct {
	Name        string
	InternalRep *Node
	fs.Inode
}

func (f *FSRSS) OnAdd(ctx context.Context) {
	node := f.NewPersistentInode(ctx, &fsdir{
		InternalRep: f.InternalRep,
	}, fs.StableAttr{Mode: syscall.S_IFDIR})
	f.AddChild(f.Name, node, true)
}

type fsdir struct {
	InternalRep *Node
	fs.Inode
}

func (f *fsdir) OnAdd(ctx context.Context) {
	if f.InternalRep == nil {
		return
	}
	for _, child := range f.InternalRep.Children {
		var embed fs.InodeEmbedder
		attr := fs.StableAttr{Mode: syscall.S_IFDIR}
		embed = &fsdir{
			InternalRep: child,
		}
		if !child.Dir {
			embed = &fs.MemRegularFile{Data: []byte(child.Content)}
			attr = fs.StableAttr{}
		}
		inode := f.NewPersistentInode(ctx, embed, attr)
		f.AddChild(child.Name, inode, true)
	}
}
