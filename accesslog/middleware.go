package accesslog

import (
	"fmt"
	"net/http"

	"github.com/megaease/access-log-go/accesslog/eventhub"
)

type (
	// Config is the configuration of access log middleware.
	Config struct {
		// Backend is the configuration of event hub.
		Backend eventhub.Config

		// ServiceName is the name of the service.
		ServiceName string
		// HostName is the host name of the service. Optional.
		HostName string

		// SkipPaths is an url path array which logs are not written. Optional.
		SkipPaths []string
		// Skipper is a function to skip logs based on provided Request. Optional.
		Skipper func(req *http.Request) bool
	}

	// AccessLogMiddleware is the middleware of access log.
	AccessLogMiddleware struct {
		serviceName string
		hostName    string

		backend eventhub.EventHub
		skip    map[string]struct{}
		skipper func(req *http.Request) bool
	}
)

func (c *Config) validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("serviceName is required")
	}
	if c.HostName == "" {
		return fmt.Errorf("hostName is required")
	}
	return nil
}

// New creates a new AccessLogMiddleware.
func New(config *Config) (*AccessLogMiddleware, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

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

	middleware := &AccessLogMiddleware{
		serviceName: config.ServiceName,
		hostName:    config.HostName,

		backend: backend,
		skip:    skip,
		skipper: config.Skipper,
	}
	return middleware, nil
}

// Close closes the AccessLogMiddleware.
func (m *AccessLogMiddleware) Close() {
	m.backend.Close()
}

// checkSkip checks if the request should be skipped.
func (m *AccessLogMiddleware) checkSkip(req *http.Request) bool {
	if _, ok := m.skip[req.URL.Path]; ok || (m.skipper != nil && m.skipper(req)) {
		return true
	}
	return false
}
