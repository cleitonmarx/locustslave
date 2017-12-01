package main

import (
	"bytes"
	"flag"
	"net/http"
	"time"

	"strconv"

	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/cleitonmarx/locustslave/infrastructure"
	"github.com/cleitonmarx/locustslave/services"
	"github.com/myzhan/boomer"
)

var parameters map[string]interface{}

func main() {

	buildParameters()

	if parameters["mode"] == "slave" {
		task := &boomer.Task{}
		switch parameters["task"].(string) {
		case "completeCheckout":
			task.Fn = completeCheckout
		case "creationCheckout":
			task.Fn = creationCheckout
		case "ticketPrice":
			task.Fn = ticketPrice
		case "country":
			task.Fn = country

		}
		boomer.Run(task)
	} else {
		completeCheckout()
	}

}

func completeCheckout() {
	var (
		err        error
		respBuffer *bytes.Buffer
		tokenId    string
	)

	apiHost := parameters["api_host"].(string)
	eventID := parameters["event_id"].(int64)
	ticketPriceID := parameters["ticket_price_id"].(int64)
	stripeKey := parameters["stripeKey"].(string)

	tracer := &infrastructure.LocustTracer{}
	httpClient := &http.Client{Timeout: 120 * time.Second}
	ckService := services.NewCheckoutService(apiHost, httpClient, tracer)
	evService := services.NewEventService(apiHost, httpClient, tracer)
	stService := services.NewStripeService(stripeKey)

	if _, err = evService.GetEvent(eventID); err == nil {
		if _, err = evService.GetTicketPrice(eventID); err == nil {
			if respBuffer, err = ckService.Create(eventID, ticketPriceID); err == nil {
				jsonContainer, _ := gabs.ParseJSON(respBuffer.Bytes())
				checkoutID := jsonContainer.Path("data.id").Data().(string)
				//time.Sleep(1 * time.Second)
				if respBuffer, err = ckService.Patch(checkoutID, respBuffer); err == nil {
					//time.Sleep(1 * time.Second)
					jsonContainer, _ = gabs.ParseJSON(respBuffer.Bytes())
					total, _ := strconv.ParseFloat(jsonContainer.Path("data.attributes.order_summary.total").Data().(string), 32)
					//fmt.Println("Total:", total)
					if total > 0 {
						if tokenId, err = stService.GetTokenID("4242424242424242"); err == nil {
							_, err = ckService.Pay(checkoutID, tokenId, respBuffer)
						} else {
							fmt.Println("Error stripe:", err)
						}
					}
					if err == nil {
						ckService.Confirm(checkoutID)
						//time.Sleep(1 * time.Second)
					}
				}
			}
		}
	}
	if err != nil {
		time.Sleep(5 * time.Second)
	}
}

func creationCheckout() {
	var (
		err error
	)

	apiHost := parameters["api_host"].(string)
	eventID := parameters["event_id"].(int64)
	ticketPriceID := parameters["ticket_price_id"].(int64)

	tracer := &infrastructure.LocustTracer{}
	httpClient := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}, Timeout: 15 * time.Second}
	ckService := services.NewCheckoutService(apiHost, httpClient, tracer)

	_, err = ckService.Create(eventID, ticketPriceID)

	if err != nil {
		time.Sleep(5 * time.Second)
	}
}

func ticketPrice() {
	var (
		err error
	)

	apiHost := parameters["api_host"].(string)
	eventID := parameters["event_id"].(int64)

	tracer := &infrastructure.LocustTracer{}
	httpClient := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}, Timeout: 15 * time.Second}
	evService := services.NewEventService(apiHost, httpClient, tracer)

	_, err = evService.GetTicketPrice(eventID)
	if err != nil {
		time.Sleep(5 * time.Second)
	}
}

func country() {
	var (
		err error
	)

	apiHost := parameters["api_host"].(string)

	tracer := &infrastructure.LocustTracer{}
	httpClient := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}, Timeout: 15 * time.Second}
	evService := services.NewEventService(apiHost, httpClient, tracer)

	_, err = evService.GetCountry()
	if err != nil {
		time.Sleep(5 * time.Second)
	}
}

func buildParameters() {

	//masterHost := flag.String("master-host", "127.0.0.1", "Host or IP address of locust master for distributed load testing. Defaults to 127.0.0.1.")
	//masterPort := flag.Int("master-port", 5557, "The port to connect to that is used by the locust master for distributed load testing. Defaults to 5557.")

	apiHost := flag.String("api_host", "127.0.0.1", "Host or IP address of Go-API server. Defaults to 127.0.0.1.")
	eventID := flag.Int64("event_id", 0, "Event id")
	ticketPriceID := flag.Int64("ticket_price_id", 0, "Host or IP address of Go-API server. Defaults to 127.0.0.1.")
	stripeKey := flag.String("stripeKey", "", "")
	mode := flag.String("mode", "slave", "Defines the operation mode. Defaults to slave")
	task := flag.String("task", "completeCheckout", "Defines the task mode. Defaults to completeCheckout")

	flag.Parse()

	if *eventID == 0 {
		panic("event_id parameter was not informed")
	}

	if *ticketPriceID == 0 {
		panic("ticket_price_id parameter was not informed")
	}

	parameters = map[string]interface{}{}
	parameters["api_host"] = *apiHost
	parameters["event_id"] = *eventID
	parameters["ticket_price_id"] = *ticketPriceID
	parameters["mode"] = *mode
	parameters["task"] = *task
	parameters["stripeKey"] = *stripeKey
	// parameters["master_host"] = *masterHost
	// parameters["master_port"] = *masterPort

}
