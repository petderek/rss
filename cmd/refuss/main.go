package main

import (
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
	HACK_DIR    = "~/refuss"
	HACK_CACHE  = "~/refuss-cache"
	HACK_CONFIG = "~/refuss-config/subscriptions.cfg"
)

func main() {
	cfg, err := loadcfg(replaceHome(HACK_CONFIG))
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}
	
	rep, err := createRep(cfg)
	if err != nil {
		log.Printf("Failed to create representation: %v", err)
		os.Exit(1)
	}
	
	opts := &fs.Options{}
	raw := fs.NewNodeFS(&rss.FSRSS{
		InternalRep: rep,
	}, opts)
	server, err := fuse.NewServer(raw, replaceHome(HACK_DIR), &opts.MountOptions)
	if err != nil {
		log.Printf("Failed to create FUSE server: %v", err)
		os.Exit(1)
	}
	go server.Serve()
	if err := server.WaitMount(); err != nil {
		log.Printf("Failed to mount filesystem: %v", err)
		os.Exit(1)
	}

	log.Println("RSS filesystem mounted successfully")
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
		log.Printf("unable to get homedir: %v", err)
		os.Exit(1)
	}
	return strings.Replace(path, "~", home, 1)
}

func loadcfg(path string) (rss.SubscriptionConfig, error) {
	simple, err := rss.FromFile(path)
	if err != nil {
		return nil, err
	}
	return rss.ToSubscription(simple), nil
}

func createRep(config rss.SubscriptionConfig) ([]*rss.Node, error) {
	rep := []*rss.Node{}
	cache := rss.NewCache(replaceHome(HACK_CACHE))
	content := rss.Content{
		config,
		*cache,
	}
	for k, _ := range config {
		data, err := content.GetFeed(k)
		if err != nil {
			log.Printf("Failed to get feed for %s: %v", k, err)
			continue // Skip this feed instead of failing entirely
		}
		in, err := rss.ToInternal(data)
		if err != nil {
			log.Printf("Failed to parse feed for %s: %v", k, err)
			continue // Skip this feed instead of failing entirely
		}
		in.Name = k
		rep = append(rep, in)
	}
	return rep, nil
}
