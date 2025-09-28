package rss

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
)

type FSRSS struct {
	fs.Inode
}

func (f *FSRSS) OnAdd(ctx context.Context) {
	node := f.NewPersistentInode(ctx, &fs.MemRegularFile{
		Data: []byte("wah"),
	}, fs.StableAttr{})
	f.AddChild("example.txt", node, true)
}
