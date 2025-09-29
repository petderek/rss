package main

import (
	"flag"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/petderek/rss"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	configDir = flag.String("config", defaultConfig(), "directory with config files")
	cacheDir  = flag.String("cache", defaultCacheDir(), "directory with cache files")
	fuseDir   = flag.String("mount", "", "target to mount (required)")
)

func main() {
	flag.Parse()
	if fuseDir == nil || *fuseDir == "" {
		log.Println("-mount must be set")
		os.Exit(1)
	}
	if !rss.Exists(*fuseDir) {
		err := os.MkdirAll(*fuseDir, 0755)
		if err != nil {
			log.Println("directory doesn't exist and can't be created: ", err)
			os.Exit(1)
		}
	}
	cfg, err := loadSubscriptions(*configDir)
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	opts := &fs.Options{}
	raw := fs.NewNodeFS(&rss.FSRSS{
		InternalRep: createRep(cfg),
	}, opts)
	server, err := fuse.NewServer(raw, *fuseDir, &opts.MountOptions)
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

func loadSubscriptions(path string) (rss.SubscriptionConfig, error) {
	simple, err := rss.FromFile(filepath.Join(path, "subscriptions.cfg"))
	if err != nil {
		return nil, err
	}
	return rss.ToSubscription(simple), nil
}

func createRep(config rss.SubscriptionConfig) []*rss.Node {
	rep := []*rss.Node{}
	cache := rss.NewCache(*cacheDir)
	content := rss.Content{
		config,
		*cache,
	}
	for k, _ := range config {
		data, err := content.GetFeed(k)
		if err != nil {
			log.Printf("Failed to get feed for %s: %v", k, err)
			// TODO failed status
			continue // Skip this feed instead of failing entirely
		}
		in, err := rss.ToInternal(data)
		if err != nil {
			log.Printf("Failed to parse feed for %s: %v", k, err)
			// TODO failed status
			continue // Skip this feed instead of failing entirely
		}
		in.Name = k
		rep = append(rep, in)
	}
	return rep
}

func defaultConfig() string {
	config, err := os.UserConfigDir()
	if err != nil {
		config = "."
	}
	return filepath.Join(config, "refuss")
}

func defaultCacheDir() string {
	cache, err := os.UserCacheDir()
	if err != nil {
		cache = "."
	}
	return filepath.Join(cache, "refuss")
}
