package api

import (
	"net/http"
	"strconv"
	"time"
)

type AccessLog struct {
	// Not used yet
	Gid     string `json:"gid"`
	Tags    string `json:"tags"`    // Tags of the log
	TraceID string `json:"traceId"` // Trace ID for distributed tracing
	SpanID  string `json:"spanId"`  // Span ID for distributed tracing
	PSpanID string `json:"pspanId"` // Parent Span ID for distributed tracing

	// Default values
	Category  string `json:"category"`  // value: "application"
	System    string `json:"system"`    // value: gpu-runtime
	Type      string `json:"type"`      // value: "access-log"
	Timestamp int64  `json:"timestamp"` // Milliseconds

	// Service values
	Service  string `json:"service"`  // Service name
	HostName string `json:"hostName"` // HostName of the server
	HostIpv4 string `json:"hostIpv4"` // IPv4 address of the server

	// Request values
	URL      string            `json:"url"`      // Requested URL
	MatchURL string            `json:"matchUrl"` // Matched URL pattern
	ClientIP string            `json:"clientIp"` // Client's IP address
	Method   string            `json:"method"`   // HTTP method (e.g., GET, POST)
	Headers  map[string]string `json:"headers"`  // Request headers

	// Response values
	StatusCode   string `json:"statusCode"`   // HTTP status code (e.g., 200, 404)
	ResponseSize int64  `json:"responseSize"` // Size of the response in bytes
	RequestTime  int64  `json:"requestTime"`  // Duration of the request (e.g., "150ms")
}

func NewAccessLog() *AccessLog {
	return &AccessLog{
		Category:  "application",
		System:    "gpu-runtime",
		Type:      "access-log",
		Timestamp: time.Now().UnixMilli(),
	}
}

func (a *AccessLog) SetService(service, hostName, hostIpv4 string) {
	a.Service = service
	a.HostName = hostName
	a.HostIpv4 = hostIpv4
}

func (a *AccessLog) SetRequest(req *http.Request, matchURL, clientIP string) {
	a.URL = req.URL.Path
	a.MatchURL = matchURL
	a.ClientIP = clientIP
	a.Method = req.Method
	a.Headers = make(map[string]string)
	for k := range req.Header {
		a.Headers[k] = req.Header.Get(k)
	}
}

func (a *AccessLog) SetResponse(statusCode int, responseSize int64, requestTime int64) {
	a.StatusCode = strconv.Itoa(statusCode)
	a.ResponseSize = responseSize
	a.RequestTime = requestTime
}
