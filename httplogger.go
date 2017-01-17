package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

type HTTPLogger struct {
	log *log.Logger
}

func NewLogger() *HTTPLogger {
	return &HTTPLogger{
		log: log.New(os.Stderr, "log - ", log.LstdFlags),
	}
}

func (l *HTTPLogger) LogRequest(req *http.Request) {
	l.log.Printf(
		"Request %s %s",
		req.Method,
		req.URL.String(),
	)
}

func (l *HTTPLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		l.log.Println(err)
	} else {
		l.log.Printf(
			"Response method=%s status=%d durationMs=%d %s",
			req.Method,
			res.StatusCode,
			duration,
			req.URL.String(),
		)
	}
}
