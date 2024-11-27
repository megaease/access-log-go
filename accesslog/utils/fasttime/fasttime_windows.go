//go:build windows
// +build windows

package fasttime

import (
	"syscall"
	"time"
)

// Now is the current time
func Now() time.Time {
	var ft syscall.Filetime
	syscall.GetSystemTimeAsFileTime(&ft)
	return time.Unix(0, ft.Nanoseconds())
}

// NowUnixNano is the current Unix Nano time
func NowUnixNano() int64 {
	var ft syscall.Filetime
	syscall.GetSystemTimeAsFileTime(&ft)
	return ft.Nanoseconds()
}
