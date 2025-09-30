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
	mount     = flag.String("mount", "", "target to mount (required)")
)

func main() {
	flag.Parse()
	if mount == nil || *mount == "" {
		log.Println("-mount must be set")
		os.Exit(1)
	}
	if !rss.Exists(*mount) {
		err := os.MkdirAll(*mount, 0755)
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

	cache := rss.NewCache(*cacheDir)
	content := &rss.Content{
		Config: cfg,
		Cache:  *cache,
	}
	
	// Validate feeds at startup to provide early feedback
	validateFeeds(content)
	
	opts := &fs.Options{}
	raw := fs.NewNodeFS(&rss.FSRSS{
		Content: content,
	}, opts)
	server, err := fuse.NewServer(raw, *mount, &opts.MountOptions)
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

// validateFeeds checks that all configured feeds can be loaded
// This provides early feedback about feed problems at startup
func validateFeeds(content *rss.Content) {
	for _, feedName := range content.ListFeeds() {
		_, err := content.GetNode(feedName)
		if err != nil {
			log.Printf("Warning: Failed to load feed %s: %v", feedName, err)
		}
	}
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
