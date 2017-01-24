package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/cleitonmarx/locustslave/trace"
)

type CheckoutService struct {
	tracer     trace.Tracer
	httpClient *http.Client
	ApiHost    string
}

func (s *CheckoutService) Create(eventID, ticketPriceID int64) (*bytes.Buffer, error) {
	span := s.tracer.NewSpan()
	defer span.Finish()
	span.SetLabel("method", "POST")
	span.SetLabel("endpoint", "v2/checkout")

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
				}`, eventID, ticketPriceID),
	)

	resp, err := s.httpClient.Post(fmt.Sprintf("%s/v2/checkout", s.ApiHost), "application/json", requestBody)
	if err != nil || resp.StatusCode != http.StatusCreated {
		if err == nil {
			respBody := getBodyBuffer(resp)
			err = fmt.Errorf("Create: Received status %d (%s)", resp.StatusCode, respBody.String())
		}
		span.SetLabel("error", err)
		return nil, err
	}
	span.SetLabel("contentLength", resp.ContentLength)
	respBody := getBodyBuffer(resp)
	return respBody, nil
}

func (s *CheckoutService) Patch(checkoutID string, createResponse *bytes.Buffer) (*bytes.Buffer, error) {
	span := s.tracer.NewSpan()
	defer span.Finish()
	span.SetLabel("method", "PATCH")
	span.SetLabel("endpoint", "v2/checkout")

	jsonContainer, _ := gabs.ParseJSON(createResponse.Bytes())
	jsonContainer.SetP("Addison", "data.attributes.invoice.first_name")
	jsonContainer.SetP("White", "data.attributes.invoice.last_name")
	jsonContainer.SetP("sofiawilson@test.com", "data.attributes.invoice.email")
	jsonContainer.Path("data.attributes.tickets").Index(0).Set("Addison", "first_name")
	jsonContainer.Path("data.attributes.tickets").Index(0).Set("White", "last_name")
	jsonContainer.Path("data.attributes.tickets").Index(0).Set("sofiawilson@test.com", "email")

	requestBody := bytes.NewBufferString(jsonContainer.String())

	request, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/v2/checkout/%s", s.ApiHost, checkoutID), requestBody)
	resp, err := s.httpClient.Do(request)
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			respBody := getBodyBuffer(resp)
			err = fmt.Errorf("Patch: Received status %d (%s)", resp.StatusCode, respBody.String())
		}
		span.SetLabel("error", err)
		return nil, err
	}

	span.SetLabel("contentLength", resp.ContentLength)
	respBody := getBodyBuffer(resp)
	return respBody, nil
}

func (s *CheckoutService) Pay(checkoutID, cardTokenID string, bodyBuffer *bytes.Buffer) (*bytes.Buffer, error) {
	span := s.tracer.NewSpan()
	defer span.Finish()
	span.SetLabel("method", "POST")
	span.SetLabel("endpoint", "v2/checkout/payment")

	jsonContainer, _ := gabs.ParseJSON(bodyBuffer.Bytes())
	jsonContainer.SetP(cardTokenID, "data.attributes.payment.source.stripe_card_token")

	requestBody := bytes.NewBufferString(jsonContainer.String())

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/checkout/%s/payment", s.ApiHost, checkoutID), requestBody)
	resp, err := s.httpClient.Do(request)
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			respBody := getBodyBuffer(resp)
			err = fmt.Errorf("Pay: Received status %d (%s)", resp.StatusCode, respBody.String())
		}
		span.SetLabel("error", err)
		return nil, err
	}

	span.SetLabel("contentLength", resp.ContentLength)
	respBody := getBodyBuffer(resp)
	return respBody, nil

}

func (s *CheckoutService) Confirm(checkoutID string) (*bytes.Buffer, error) {
	span := s.tracer.NewSpan()
	defer span.Finish()
	span.SetLabel("method", "POST")
	span.SetLabel("endpoint", "v2/checkout/confirm")

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/checkout/%s/confirm", s.ApiHost, checkoutID), nil)
	resp, err := s.httpClient.Do(request)
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			respBody := getBodyBuffer(resp)
			err = fmt.Errorf("Confirm: Received status %d (%s)", resp.StatusCode, respBody.String())
		}
		span.SetLabel("error", err)
		return nil, err
	}

	span.SetLabel("contentLength", resp.ContentLength)
	respBody := getBodyBuffer(resp)
	return respBody, nil
}

func getBodyBuffer(r *http.Response) *bytes.Buffer {
	responseBody := bytes.NewBufferString("")
	io.Copy(responseBody, r.Body)
	defer r.Body.Close()
	return responseBody
}

func NewCheckoutService(apiHost string, client *http.Client, tracer trace.Tracer) *CheckoutService {
	return &CheckoutService{
		httpClient: client,
		ApiHost:    apiHost,
		tracer:     tracer,
	}
}
