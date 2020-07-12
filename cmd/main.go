package main

import (
	"flag"
	"io/ioutil"
	"os"
	"time"

	cachePkg "github.com/JokeTrue/image-previewer/pkg/cache"
	fetcherPkg "github.com/JokeTrue/image-previewer/pkg/fetcher"
	transformerPkg "github.com/JokeTrue/image-previewer/pkg/transformer"

	"github.com/JokeTrue/image-previewer/pkg/middleware"

	"github.com/JokeTrue/image-previewer/pkg/app"

	"github.com/JokeTrue/image-previewer/pkg/service"
	"github.com/NYTimes/gziphandler"
	"github.com/justinas/alice"

	"github.com/JokeTrue/image-previewer/pkg/logging"
)

var (
	appName         = "image-previewer"
	addr            = flag.String("addr", ":8080", "App addr")
	connectTimeout  = flag.Duration("connect-timeout", 25*time.Second, "Сonnection timeout")
	requestTimeout  = flag.Duration("request-timeout", 25*time.Second, "Request timeout")
	shutdownTimeout = flag.Duration("shutdown-timeout", 30*time.Second, "Graceful shutdown timeout")
	cacheDir        = flag.String("cache-dir", "", "Path to Cache dir")
	cacheSize       = flag.Int("cache-size", 5, "Size of cache")
)

func main() {
	flag.Parse()

	// 1. Setup required Units
	logger := logging.DefaultLogger
	fetcher := fetcherPkg.NewFetcher(logger, *connectTimeout, *requestTimeout)
	transformer := transformerPkg.NewTransformer()

	// 2. If cacheDir isn't provided, then use Temporary Dir
	if *cacheDir == "" {
		var err error
		*cacheDir, err = ioutil.TempDir("", "")
		if err != nil {
			logger.Fatal(err)
		}
		defer os.RemoveAll(*cacheDir)
	}

	// 3. Setup Cache
	cache, err := cachePkg.NewCacheWithEvict(*cacheSize, func(key interface{}, value interface{}) {
		if path, ok := value.(string); ok {
			os.Remove(path)
		}
	})
	if err != nil {
		logger.Fatal(err)
	}

	// 4. Setup Application
	application := app.NewApplication(*cacheDir, logger, fetcher, transformer, cache)
	srv := service.NewHTTPServer(*addr, *shutdownTimeout, alice.New(
		gziphandler.GzipHandler,
		middleware.Logger(logger),
	).Then(application.Run()))

	// 5. Run Application
	service.Run(srv, appName)
}
