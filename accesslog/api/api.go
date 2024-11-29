package api

import (
	"net/http"
	"strconv"
	"time"
)

// AccessLog is the structure of access log.
type AccessLog struct {
	// Not used yet
	Gid      string `json:"gid"`
	Tags     string `json:"tags"`
	TraceID  string `json:"trace_id"`
	SpanID   string `json:"span_id"`
	PSpanID  string `json:"pspan_id"`
	TenantID string `json:"tenant_id"`

	// Initial values
	Category  string `json:"category"`  // value: "application"
	System    string `json:"system"`    // value: gpu-runtime
	Type      string `json:"type"`      // value: "access-log"
	Timestamp int64  `json:"timestamp"` // Milliseconds
	Service   string `json:"service"`   // Service name
	HostName  string `json:"host_name"` // HostName of the server

	// Request values
	HostIpv4 string            `json:"host_ipv4"` // IPv4 address of the server
	URL      string            `json:"url"`       // Requested URL
	MatchURL string            `json:"match_url"` // Matched URL pattern
	ClientIP string            `json:"client_ip"` // Client's IP address
	Method   string            `json:"method"`    // HTTP method (e.g., GET, POST)
	Headers  map[string]string `json:"headers"`   // Request headers
	Queries  map[string]string `json:"queries"`   // Request queries

	// Response values
	StatusCode   string `json:"status_code"`   // HTTP status code (e.g., 200, 404)
	ResponseSize int64  `json:"response_size"` // Size of the response in bytes
	RequestTime  int64  `json:"request_time"`  // Duration of the request (e.g., "150ms")
}

// NewAccessLog creates a new AccessLog instance.
func NewAccessLog(service string, hostName string, tenantID string) *AccessLog {
	return &AccessLog{
		Category:  "application",
		System:    "gpu-runtime",
		Type:      "access-log",
		TenantID:  tenantID,
		Service:   service,
		HostName:  hostName,
		Timestamp: time.Now().UnixMilli(),
	}
}

// SetRequest sets the request information.
func (a *AccessLog) SetRequest(req *http.Request, matchURL, clientIP string, hostIP string) {
	a.URL = req.URL.Path
	a.MatchURL = matchURL
	a.ClientIP = clientIP
	a.Method = req.Method
	a.HostIpv4 = hostIP

	a.Headers = make(map[string]string)
	for k, v := range req.Header.Clone() {
		if len(v) > 0 {
			a.Headers[k] = v[0]
		} else {
			a.Headers[k] = ""
		}
	}

	a.Queries = make(map[string]string)
	for k, v := range req.URL.Query() {
		if len(v) > 0 {
			a.Queries[k] = v[0]
		} else {
			a.Queries[k] = ""
		}
	}
}

// SetResponse sets the response information.
func (a *AccessLog) SetResponse(statusCode int, responseSize int64, requestTime int64) {
	a.StatusCode = strconv.Itoa(statusCode)
	a.ResponseSize = responseSize
	a.RequestTime = requestTime
}
