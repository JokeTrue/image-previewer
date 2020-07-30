package main

import (
	"flag"
	"io/ioutil"
	"os"
	"time"

	fetcherPkg "github.com/JokeTrue/image-previewer/pkg/fetcher"
	transformerPkg "github.com/JokeTrue/image-previewer/pkg/transformer"

	"github.com/JokeTrue/image-previewer/pkg/middleware"

	"github.com/JokeTrue/image-previewer/pkg/app"

	"github.com/JokeTrue/image-previewer/pkg/service"
	"github.com/NYTimes/gziphandler"
	"github.com/justinas/alice"

	"github.com/JokeTrue/image-previewer/pkg/logging"
	lru "github.com/hashicorp/golang-lru"
)

var (
	appName         = "image-previewer"
	addr            string
	connectTimeout  time.Duration
	requestTimeout  time.Duration
	shutdownTimeout time.Duration
	cacheDir        string
	cacheSize       int
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "App addr")
	flag.DurationVar(&connectTimeout, "connect-timeout", 25*time.Second, "Сonnection timeout")
	flag.DurationVar(&requestTimeout, "request-timeout", 25*time.Second, "Request timeout")
	flag.DurationVar(&shutdownTimeout, "shutdown-timeout", 30*time.Second, "Graceful shutdown timeout")
	flag.StringVar(&cacheDir, "cache-dir", "", "Path to Cache dir")
	flag.IntVar(&cacheSize, "cache-size", 5, "Size of cache")
}

func main() {
	flag.Parse()

	// 1. Setup required Units
	logger := logging.DefaultLogger
	fetcher := fetcherPkg.NewFetcher(logger, connectTimeout, requestTimeout)
	cropper := transformerPkg.NewCropper()

	// 2. If cacheDir isn't provided, then use Temporary Dir
	if cacheDir == "" {
		var err error
		cacheDir, err = ioutil.TempDir("", "")
		if err != nil {
			logger.WithError(err).Fatal(err)
		}
		defer func() {
			if err := os.RemoveAll(cacheDir); err != nil {
				logger.WithError(err).Error("failed to remove cache dir")
			}
		}()
	}

	// 3. Setup Cache
	cache, err := lru.NewWithEvict(cacheSize, func(key interface{}, value interface{}) {
		if path, ok := value.(string); ok {
			defer func() {
				err := os.Remove(path)
				logger.WithError(err).Fatal("failed to remove item from cache")
			}()
		}
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to setup cache")
	}

	// 4. Setup Application
	application := app.NewApplication(cacheDir, logger, fetcher, cropper, cache)
	srv := service.NewHTTPServer(addr, shutdownTimeout, alice.New(
		gziphandler.GzipHandler,
		middleware.Logger(logger),
	).Then(application.Run()))

	// 5. Run Application
	service.Run(srv, appName)
}
