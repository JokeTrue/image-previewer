package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/JokeTrue/image-previewer/pkg/logging"
)

var (
	ready     bool
	readyLock sync.RWMutex
	alive     bool
	aliveLock sync.RWMutex
)

func SetReady(state bool) {
	readyLock.Lock()
	defer readyLock.Unlock()
	ready = state
}

func IsReady() bool {
	readyLock.RLock()
	defer readyLock.RUnlock()
	return ready
}

func SetAlive(state bool) {
	aliveLock.Lock()
	defer aliveLock.Unlock()
	alive = state
}

func IsAlive() bool {
	aliveLock.RLock()
	defer aliveLock.RUnlock()
	return alive
}

func init() {
	SetAlive(true)
}

type HTTPServer struct {
	Server          *http.Server
	shutdownTimeout time.Duration
	done            chan struct{}
}

func NewHTTPServer(addr string, shutdownTimeout time.Duration, router http.Handler) *HTTPServer {
	r := http.NewServeMux()
	r.Handle("/", router)
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if IsAlive() {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	})
	r.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if IsReady() {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &HTTPServer{
		shutdownTimeout: shutdownTimeout,
		Server:          srv,
		done:            make(chan struct{}),
	}
}

func (s *HTTPServer) Start() {
	go func() {
		defer close(s.done)
		logging.Infof("starting http server on %s", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
			logging.WithError(err).Fatal("http server failure")
		}
		logging.Infof("http server on %s stopped listening", s.Server.Addr)
	}()
}

func (s *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	if err := s.Server.Shutdown(ctx); err != nil {
		logging.WithError(err).Error("http shutdown error")
	}
	logging.Infof("http server on %s stopped", s.Server.Addr)
	<-s.done
	cancel()
}
