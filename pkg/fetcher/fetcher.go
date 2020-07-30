package fetcher

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/JokeTrue/image-previewer/pkg/utils"

	"github.com/JokeTrue/image-previewer/pkg/logging"
	"github.com/pkg/errors"
)

var (
	ErrNotSupportedContentType = errors.New("got not supported content type")
	ErrNotSupportedScheme      = errors.New("got not supported scheme")
	SupportedContentTypes      = []string{"image/jpeg"}
)

type Fetcher interface {
	Fetch(ctx context.Context, url string, header http.Header) ([]byte, error)
}

type HTTPFetcher struct {
	logger         logging.Logger
	transport      http.RoundTripper
	requestTimeout time.Duration
}

func NewFetcher(l logging.Logger, connectTimeout time.Duration, requestTimeout time.Duration) *HTTPFetcher {
	return &HTTPFetcher{
		logger:         l,
		requestTimeout: requestTimeout,
		transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: connectTimeout,
			}).DialContext,
		},
	}
}

func (f HTTPFetcher) Fetch(ctx context.Context, url string, header http.Header) ([]byte, error) {
	proxyRequest, err := prepareRequest(ctx, url, header)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare request")
	}
	responseBody, err := f.doRequest(proxyRequest)
	if err != nil {
		return nil, errors.Wrap(err, "error making request")
	}
	return responseBody, nil
}

func prepareRequest(ctx context.Context, rawURL string, header http.Header) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create proxy request")
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse image url")
	}
	if parsedURL.Scheme != "http" {
		return nil, ErrNotSupportedScheme
	}
	request.URL = parsedURL
	request.Header = header
	return request, nil
}

func (f *HTTPFetcher) doRequest(request *http.Request) ([]byte, error) {
	client := http.Client{
		Timeout:   f.requestTimeout,
		Transport: f.transport,
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			f.logger.WithError(errClose).Error("failed to close body")
		}
	}()

	if !utils.Contains(SupportedContentTypes, resp.Header.Get("Content-type")) {
		return nil, ErrNotSupportedContentType
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}
	return buff, nil
}
