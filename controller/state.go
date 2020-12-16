package controller

import (
	"net/http"
	"time"
)

type (
	RequestFunc func(req *http.Request) (*http.Response, error)
	Status      string
)

const (
	ConnectionStatusOpen     = "Open"
	ConnectionStatusHalfOpen = "HalfOpen"
	ConnectionStatusClose    = "Close"
)

//Custom Circuit Breaker
type ConnectionHandler struct {
	Send          RequestFunc
	CurrentStatus Status
	RetryCount    int
	MaxRetryCount int
	Timeout       time.Duration
	Expire        time.Time
}

func (s *ConnectionHandler) Close() {

}

func (s *ConnectionHandler) Open() {

}

func (s *ConnectionHandler) HalfOpen() {

}
