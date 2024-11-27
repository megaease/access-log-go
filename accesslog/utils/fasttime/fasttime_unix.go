//go:build !windows
// +build !windows

package fasttime

import (
	"syscall"
	"time"
)

// Now is the current time
func Now() time.Time {
	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)
	return time.Unix(0, syscall.TimevalToNsec(tv))
}

// NowUnixNano is the current Unix Nano time
func NowUnixNano() int64 {
	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)
	return syscall.TimevalToNsec(tv)
}
