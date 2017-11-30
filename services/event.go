package services

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cleitonmarx/locustslave/trace"
)

type EventService struct {
	tracer     trace.Tracer
	httpClient *http.Client
	ApiHost    string
}

func (s *EventService) GetEvent(eventID int64) (*bytes.Buffer, error) {
	span := s.tracer.NewSpan()
	defer span.Finish()
	span.SetLabel("method", "GET")
	span.SetLabel("endpoint", "v2/event")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s/v2/event/%d", s.ApiHost, eventID))
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			respBody := getBodyBuffer(resp)
			err = fmt.Errorf("GetEvent: Received status %d (%s)", resp.StatusCode, respBody.String())
		}
		span.SetLabel("error", err)
		return nil, err
	}

	span.SetLabel("contentLength", resp.ContentLength)
	respBody := getBodyBuffer(resp)
	return respBody, nil
}

func (s *EventService) GetTicketPrice(eventID int64) (*bytes.Buffer, error) {
	span := s.tracer.NewSpan()
	defer span.Finish()
	span.SetLabel("method", "GET")
	span.SetLabel("endpoint", "v2/ticket_price")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s/v2/ticket_price?filter[event_id]=%d&page[limit]=20&page[offset]=0", s.ApiHost, eventID))
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			respBody := getBodyBuffer(resp)
			err = fmt.Errorf("GetEvent: Received status %d (%s)", resp.StatusCode, respBody.String())
		}
		span.SetLabel("error", err)
		return nil, err
	}

	span.SetLabel("contentLength", resp.ContentLength)
	respBody := getBodyBuffer(resp)
	return respBody, nil
}

func (s *EventService) GetCountry() (*bytes.Buffer, error) {
	span := s.tracer.NewSpan()
	defer span.Finish()
	span.SetLabel("method", "GET")
	span.SetLabel("endpoint", "v2/country/1")

	resp, err := s.httpClient.Get(fmt.Sprintf("%s/v2/country/1", s.ApiHost))
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			respBody := getBodyBuffer(resp)
			err = fmt.Errorf("GetCountry: Received status %d (%s)", resp.StatusCode, respBody.String())
		}
		span.SetLabel("error", err)
		return nil, err
	}

	span.SetLabel("contentLength", resp.ContentLength)
	respBody := getBodyBuffer(resp)
	return respBody, nil
}

func NewEventService(apiHost string, client *http.Client, tracer trace.Tracer) *EventService {
	return &EventService{
		httpClient: client,
		ApiHost:    apiHost,
		tracer:     tracer,
	}
}
