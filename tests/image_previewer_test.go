package tests

import (
	"context"
	_ "image/jpeg"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	client *http.Client
}

func NewTestSuite() *TestSuite {
	return &TestSuite{client: http.DefaultClient}
}

func (s TestSuite) DoRequest(t *testing.T, url string, width, height int) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://image-previewer:8080", nil)
	assert.NoError(t, err)

	q := req.URL.Query()
	q.Add("width", strconv.FormatInt(int64(width), 10))
	q.Add("height", strconv.FormatInt(int64(height), 10))
	q.Add("url", url)
	req.URL.RawQuery = q.Encode()

	res, err := s.client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	return res, b, err
}

func TestFill(t *testing.T) {
	s := NewTestSuite()

	url := "http://nginx:80/gopher.jpg"
	width, height := 333, 666

	// nolint
	res, body, err := s.DoRequest(t, url, width, height)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, len(body) > 0)
	assert.True(t, res.Header.Get("Content-Type") == "image/jpeg")
}

func TestServerDoesntExist(t *testing.T) {
	s := NewTestSuite()

	url := "http://not_exist.com/gopher.jpg"
	width, height := 333, 666

	// nolint
	res, _, err := s.DoRequest(t, url, width, height)
	assert.NoError(t, err)

	require.Equal(t, http.StatusBadGateway, res.StatusCode)
}

func TestCropNotImage(t *testing.T) {
	s := NewTestSuite()

	url := "http://ngingx:80/text.txt"
	width, height := 333, 666

	// nolint
	res, _, err := s.DoRequest(t, url, width, height)
	assert.NoError(t, err)

	require.Equal(t, http.StatusBadGateway, res.StatusCode)
}

func TestURLWrongScheme(t *testing.T) {
	s := NewTestSuite()

	url := "ftp://ngingx:80/gopher.jpg"
	width, height := 333, 666

	// nolint
	res, body, err := s.DoRequest(t, url, width, height)
	assert.NoError(t, err)

	require.Equal(t, http.StatusBadGateway, res.StatusCode)
	require.False(t, strings.Contains(string(body), "got not supported scheme"))
}
