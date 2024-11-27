package accesslog

import (
	"megaease/access-log-go/accesslog/eventhub"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Backend eventhub.Config

		// ServiceName is the name of the service.
		ServiceName string

		// SkipPaths is an url path array which logs are not written. Optional.
		SkipPaths []string
		// Skipper is a function to skip logs based on provided Request. Optional.
		Skipper func(req *http.Request) bool
	}

	AccessLogMiddleware struct {
		hostname    string
		serviceName string

		backend eventhub.EventHub
		skip    map[string]struct{}
		skipper func(req *http.Request) bool
	}
)

func New(config *Config) (*AccessLogMiddleware, error) {
	backend, err := eventhub.New(&config.Backend)
	if err != nil {
		return nil, err
	}

	var skip map[string]struct{}
	if length := len(config.SkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range config.SkipPaths {
			skip[path] = struct{}{}
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Errorf("get hostname failed: %v", err)
		hostname = "unknown"
	}

	middleware := &AccessLogMiddleware{
		serviceName: config.ServiceName,
		hostname:    hostname,

		backend: backend,
		skip:    skip,
		skipper: config.Skipper,
	}
	return middleware, nil
}

func (m *AccessLogMiddleware) Close() {
	m.backend.Close()
}

func (m *AccessLogMiddleware) checkSkip(req *http.Request) bool {
	if _, ok := m.skip[req.URL.Path]; ok || (m.skipper != nil && m.skipper(req)) {
		return true
	}
	return false
}
