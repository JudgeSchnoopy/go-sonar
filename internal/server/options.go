package server

import (
	"fmt"
	"time"
)

type Timeouts struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

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

func WithCustomPort(port int) ServerOption {
	return func(server *Server) {
		server.http.Addr = fmt.Sprintf(":%v", port)
	}
}

func WithCustomSchedule(interval time.Duration) ServerOption {
	return func(server *Server) {
		server.scheduledInterval = interval
	}
}
