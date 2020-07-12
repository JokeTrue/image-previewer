package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JokeTrue/image-previewer/pkg/logging"
)

// Logger middleware setup logger and logs all requests.
func Logger(logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.WithContext(r.Context()).WithFields(map[string]interface{}{
				"request":                   fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
				"request_method":            r.Method,
				"request_uri":               r.RequestURI,
				"request_proto":             r.Proto,
				"request_duration_ms":       int(time.Since(start).Seconds() * 1000),
				"real_ip":                   r.Header.Get("X-Real-IP"),
				"proxy_add_x_forwarded_for": r.Header.Get("X-Forwarded-For"),
				"remote_addr":               r.RemoteAddr,
				"http_referrer":             r.Referer(),
				"http_user_agent":           r.UserAgent(),
			}).Info("http request handled")
		})
	}
}
