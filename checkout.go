package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/ernesto-jimenez/httplogger"
	"github.com/myzhan/boomer"
)

type CheckoutOptions struct {
	ApiHost       string
	EventID       int64
	TicketPriceID int64
}

type CheckoutService struct {
	httpClient *http.Client
	CheckoutID string
	Options    CheckoutOptions
}

func (s *CheckoutService) Create() (*bytes.Buffer, error) {
	startTime := now()
	requestBody := bytes.NewBufferString(
		fmt.Sprintf(`{
					"data":{
						"attributes":{
							"event_id":%d,
							"tickets":[
								{
									"ticket_price":{
										"ticket_price_id":%d
									}
								}
							]
						},
						"type":"checkout"
					}
				}`, s.Options.EventID, s.Options.TicketPriceID),
	)

	resp, err := s.httpClient.Post(fmt.Sprintf("https://%s/v2/checkout", s.Options.ApiHost), "application/json", requestBody)
	if err != nil {
		boomer.Events.Publish("request_failure", "POST", "v2/checkout", 0.0, err.Error())
		return nil, err
	}

	responseBody := bytes.NewBufferString("")
	io.Copy(responseBody, resp.Body)
	defer resp.Body.Close()

	endTime := now()
	boomer.Events.Publish("request_success", "POST", "v2/checkout", float64(endTime-startTime), resp.ContentLength)

	return responseBody, nil
}

func (s *CheckoutService) Patch(createResponse *bytes.Buffer) (*bytes.Buffer, error) {
	startTime := now()

	jsonContainer, _ := gabs.ParseJSON(createResponse.Bytes())
	jsonContainer.SetP("Addison", "data.attributes.invoice.first_name")
	jsonContainer.SetP("White", "data.attributes.invoice.last_name")
	jsonContainer.SetP("sofiawilson@test.com", "data.attributes.invoice.email")
	jsonContainer.Path("data.attributes.tickets").Index(0).Set("Addison", "first_name")
	jsonContainer.Path("data.attributes.tickets").Index(0).Set("White", "last_name")
	jsonContainer.Path("data.attributes.tickets").Index(0).Set("sofiawilson@test.com", "email")
	s.CheckoutID = jsonContainer.Path("data.id").Data().(string)

	requestBody := bytes.NewBufferString(jsonContainer.String())

	request, _ := http.NewRequest("PATCH", fmt.Sprintf("https://%s/v2/checkout/%s", s.Options.ApiHost, s.CheckoutID), requestBody)
	resp, err := s.httpClient.Do(request)
	if err != nil {
		boomer.Events.Publish("request_failure", "PATCH", "v2/checkout", 0.0, err.Error())
		return nil, err
	}

	responsePatch := bytes.NewBufferString("")
	io.Copy(responsePatch, resp.Body)
	defer resp.Body.Close()

	endTime := now()
	boomer.Events.Publish("request_success", "PATCH", "v2/checkout", float64(endTime-startTime), resp.ContentLength)
	return responsePatch, nil
}

func (s *CheckoutService) Confirm() (*bytes.Buffer, error) {
	startTime := now()

	request, _ := http.NewRequest("POST", fmt.Sprintf("https://%s/v2/checkout/%s/confirm", s.Options.ApiHost, s.CheckoutID), nil)
	resp, err := s.httpClient.Do(request)
	if err != nil {
		boomer.Events.Publish("request_failure", "POST", "v2/checkout/confirm", 0.0, err.Error())
		return nil, err
	}

	responseConfirm := bytes.NewBufferString("")
	io.Copy(responseConfirm, resp.Body)
	defer resp.Body.Close()

	endTime := now()
	boomer.Events.Publish("request_success", "PATCH", "v2/checkout/confirm", float64(endTime-startTime), resp.ContentLength)
	return responseConfirm, nil
}

func now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func NewCheckoutService(options CheckoutOptions) *CheckoutService {

	transport := httplogger.NewLoggedTransport(&http.Transport{
		DisableKeepAlives: true,
	}, NewLogger())

	return &CheckoutService{
		httpClient: &http.Client{Transport: transport},
		Options:    options,
	}
}
