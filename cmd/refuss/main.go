package main

import (
	_ "embed"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/petderek/rss"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	HACK_URL  = "https://feeds.npr.org/1019/rss.xml"
	HACK_NAME = "npr"
	HACK_DIR  = "~/refuss"
)

//go:embed hack.xml
var HackContent []byte

func main() {
	rep, err := rss.ToInternal(HackContent)
	if err != nil {
		log.Fatal(err)
	}

	opts := &fs.Options{}
	raw := fs.NewNodeFS(&rss.FSRSS{
		Name:        HACK_NAME,
		InternalRep: rep,
	}, opts)
	server, err := fuse.NewServer(raw, replaceHome(HACK_DIR), &opts.MountOptions)
	if err != nil {
		log.Fatal(err)
	}
	go server.Serve()
	if err := server.WaitMount(); err != nil {
		log.Fatal(err)
	}

	go handleSignals(server)
	server.Wait()
}

func handleSignals(server *fuse.Server) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	for {
		<-shutdown
		log.Println("received shutdown signal")
		err := server.Unmount()
		if err != nil {
			log.Println("error calling unmount: ", err)
		}
	}
}

func replaceHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("unable to get homedir: ", err)
	}
	return strings.Replace(path, "~", home, 1)
}
