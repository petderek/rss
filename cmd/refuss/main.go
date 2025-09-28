package main

import (
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/petderek/rss"
	"log"
	"os"
	"strings"
	"time"
)

const (
	HACK_URL  = "https://feeds.npr.org/1019/rss.xml"
	HACK_NAME = "npr"
	HACK_DIR  = "~/refuss"
)

// go:embed hack.xml
var HACK_CONTENT string

func main() {
	opts := &fs.Options{}
	raw := fs.NewNodeFS(&rss.FSRSS{}, opts)
	server, err := fuse.NewServer(raw, replaceHome(HACK_DIR), &opts.MountOptions)
	if err != nil {
		log.Fatal(err)
	}
	go server.Serve()
	if err := server.WaitMount(); err != nil {
		log.Fatal(err)
	}
	log.Println("waiting")
	time.Sleep(30 * time.Second)
	log.Println("stopping")
	err = server.Unmount()
	if err != nil {
		log.Println("unable to unmount: ", err)
	}
	log.Println("done")
}

func replaceHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("unable to get homedir: ", err)
	}
	return strings.Replace(path, "~", home, 1)
}
