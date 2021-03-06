package server

import (
	"fmt"
	"time"
)

// ServerOption overrides a default server value
type ServerOption func(*Server)

// Timeouts are customizable timeout settings for Sonar
type Timeouts struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

// WithCustomTimeouts sets server timeout values
func WithCustomTimouts(to Timeouts) ServerOption {
	return func(server *Server) {
		if to.WriteTimeout != 0 {
			server.http.WriteTimeout = to.WriteTimeout
		}
		if to.ReadTimeout != 0 {
			server.http.ReadTimeout = to.ReadTimeout
		}
	}
}

// WithCustomPort sets the Sonar port
func WithCustomPort(port int) ServerOption {
	return func(server *Server) {
		server.http.Addr = fmt.Sprintf(":%v", port)
	}
}

// WithCustomSchedule overwrites the default check-in schedule interval
func WithCustomSchedule(interval time.Duration) ServerOption {
	return func(server *Server) {
		server.scheduledInterval = interval
	}
}
