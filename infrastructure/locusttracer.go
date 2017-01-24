package infrastructure

import (
	"time"

	"github.com/cleitonmarx/locustslave/trace"
	"github.com/myzhan/boomer"
)

type LocustTracer struct {
}

func (t *LocustTracer) NewSpan() trace.Span {
	return &LocustSpan{
		labels:    map[string]interface{}{},
		startTime: now(),
	}
}

type LocustSpan struct {
	labels    map[string]interface{}
	startTime int64
}

func (s *LocustSpan) SetLabel(key string, value interface{}) {
	s.labels[key] = value
}
func (s *LocustSpan) Finish() {
	endTime := now()
	method := s.labels["method"].(string)
	endpoint := s.labels["endpoint"].(string)
	err := s.labels["error"]

	if err != nil {
		boomer.Events.Publish("request_failure", method, endpoint, float64(endTime-s.startTime), err.(error).Error())
	} else {
		contentLength := s.labels["contentLength"].(int64)
		boomer.Events.Publish("request_success", method, endpoint, float64(endTime-s.startTime), contentLength)
	}
}

func now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
