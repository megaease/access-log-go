package api

import (
	"net/http"
	"strconv"
	"time"
)

// AccessLog is the structure of access log.
type AccessLog struct {
	// Not used yet
	Gid     string `json:"gid"`
	Tags    string `json:"tags"`
	TraceID string `json:"traceId"`
	SpanID  string `json:"spanId"`
	PSpanID string `json:"pspanId"`

	// Initial values
	Category  string `json:"category"`  // value: "application"
	System    string `json:"system"`    // value: gpu-runtime
	Type      string `json:"type"`      // value: "access-log"
	Timestamp int64  `json:"timestamp"` // Milliseconds
	Service   string `json:"service"`   // Service name

	// Request values
	HostName string            `json:"hostName"` // HostName of the server
	HostIpv4 string            `json:"hostIpv4"` // IPv4 address of the server
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

// NewAccessLog creates a new AccessLog instance.
func NewAccessLog(service string) *AccessLog {
	return &AccessLog{
		Category:  "application",
		System:    "gpu-runtime",
		Type:      "access-log",
		Service:   service,
		Timestamp: time.Now().UnixMilli(),
	}
}

// SetRequest sets the request information.
func (a *AccessLog) SetRequest(req *http.Request, matchURL, clientIP string) {
	a.URL = req.URL.Path
	a.MatchURL = matchURL
	a.ClientIP = clientIP
	a.Method = req.Method
	a.HostName = req.Host
	a.HostIpv4 = req.RemoteAddr

	a.Headers = make(map[string]string)
	for k, v := range req.Header.Clone() {
		if len(v) > 0 {
			a.Headers[k] = v[0]
		} else {
			a.Headers[k] = ""
		}
	}
}

// SetResponse sets the response information.
func (a *AccessLog) SetResponse(statusCode int, responseSize int64, requestTime int64) {
	a.StatusCode = strconv.Itoa(statusCode)
	a.ResponseSize = responseSize
	a.RequestTime = requestTime
}
