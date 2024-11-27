package accesslog

import (
	"megaease/access-log-go/accesslog/api"
	"megaease/access-log-go/accesslog/eventhub"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinAccessLogMiddleware(t *testing.T) {
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

	// setup gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.GetGinMiddleware())
	router.GET("/test/:testid", func(c *gin.Context) {
		c.String(http.StatusOK, "Test Passed")
	})
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	{
		req, _ := http.NewRequest(http.MethodGet, "/test/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Test Passed", w.Body.String())
	}
	{
		req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
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