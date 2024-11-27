package accesslog

import (
	"megaease/access-log-go/accesslog/api"
	"megaease/access-log-go/accesslog/eventhub"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEchoAccessLogMiddleware(t *testing.T) {
	// setup access log middleware
	config := &Config{
		Backend: eventhub.Config{
			Type: eventhub.EventHubTypeMock,
		},
		ServiceName: "test",
		SkipPaths:   []string{"/healthz"},
	}
	middleware, err := New(config)
	assert.Nil(t, err)
	defer middleware.Close()
	mockHub := middleware.backend.(*eventhub.EventHubMock)

	// setup Echo router
	e := echo.New()
	e.Use(middleware.GetEchoMiddleWare())
	e.GET("/test/:testid", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test Passed")
	})
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	{
		req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Test Passed", rec.Body.String())
	}
	{
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	}
	assert.Equal(t, 1, len(mockHub.Logs))
	logEntry := mockHub.Logs[0]
	expectedLog := &api.AccessLog{
		Category:     "application",
		System:       "gpu-runtime",
		Type:         "access-log",
		Service:      "test",
		URL:          "/test/123",
		MatchURL:     "/test/:testid",
		Method:       "GET",
		Headers:      map[string]string{},
		StatusCode:   "200",
		ResponseSize: 11,
	}
	assert.Equal(t, expectedLog.Category, logEntry.Category)
	assert.Equal(t, expectedLog.System, logEntry.System)
	assert.Equal(t, expectedLog.Type, logEntry.Type)
	assert.Equal(t, expectedLog.Service, logEntry.Service)
	assert.Equal(t, expectedLog.URL, logEntry.URL)
	assert.Equal(t, expectedLog.MatchURL, logEntry.MatchURL)
	assert.Equal(t, expectedLog.Method, logEntry.Method)
	assert.Equal(t, expectedLog.Headers, logEntry.Headers)
	assert.Equal(t, expectedLog.StatusCode, logEntry.StatusCode)
	assert.Equal(t, expectedLog.ResponseSize, logEntry.ResponseSize)
}
