package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/pkg/errors"

	"github.com/JokeTrue/image-previewer/pkg/fetcher"
	"github.com/JokeTrue/image-previewer/pkg/logging"
	"github.com/JokeTrue/image-previewer/pkg/lru"
	"github.com/JokeTrue/image-previewer/pkg/transformer"
	"github.com/JokeTrue/image-previewer/pkg/utils"
)

type Application struct {
	cacheDir    string
	logger      logging.Logger
	fetcher     fetcher.Fetcher
	transformer transformer.Transformer
	cache       lru.Cache
}

func NewApplication(cacheDir string, l logging.Logger, f fetcher.Fetcher, t transformer.Transformer, c lru.Cache) *Application {
	return &Application{cacheDir: cacheDir, logger: l, fetcher: f, transformer: t, cache: c}
}

func (a *Application) Run() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		url := params.Get("url")
		rawWidth, rawHeight := params.Get("width"), params.Get("height")

		ctx := logging.ContextWithFields(r.Context(), logging.Fields{
			"url":     url,
			"width":   rawWidth,
			"height":  rawHeight,
			"headers": r.Header,
		})
		ctxLogger := logging.WithContext(ctx)

		width, err := strconv.Atoi(rawWidth)
		if err != nil {
			ctxLogger.WithError(err).Error("failed to parse width")
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		height, err := strconv.Atoi(rawHeight)
		if err != nil {
			ctxLogger.WithError(err).Error("failed to parse height")
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		img, err := a.handle(ctx, url, r.Header, width, height)
		if err != nil {
			ctxLogger.WithError(err).Error("failed to handle request")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "image/jpeg")
		if _, err := w.Write(img); err != nil {
			ctxLogger.WithError(err).Error("failed to write response")
		}
	})
	return r
}

func (a *Application) handle(ctx context.Context, url string, header http.Header, width, height int) ([]byte, error) {
	// 1. Try to find image in Cache
	cacheKey := utils.GetMD5Hash(fmt.Sprintf("%s|%d|%d", url, width, height))
	if imgPath, found := a.cache.Get(cacheKey); found {
		img, err := ioutil.ReadFile(imgPath.(string))
		return img, err
	}

	// 2. If not found in Cache, then try to fetch Image
	img, err := a.fetcher.Fetch(ctx, url, header)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch image")
	}

	// 2. Transform Image
	img, err = a.transformer.Crop(img, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "failed to crop image")
	}

	// 3. Save transformed Image
	imgPath := path.Join(a.cacheDir, cacheKey+".jpeg")
	err = ioutil.WriteFile(imgPath, img, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save image")
	}

	// 4. Add Image path to Cache
	a.cache.Add(cacheKey, imgPath)

	return img, nil
}
