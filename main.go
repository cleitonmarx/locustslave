package main

import (
	"flag"
	"time"

	"github.com/myzhan/boomer"
)

var parameters map[string]interface{}

func main() {

	buildParameters()

	if parameters["mode"] == "slave" {
		task := &boomer.Task{
			Weight: 10,
			Fn:     execTask,
		}
		boomer.Run(task)
	} else {
		execTask()
	}

}

func execTask() {
	ckService := NewCheckoutService(CheckoutOptions{
		ApiHost:       parameters["api_host"].(string),
		EventID:       parameters["event_id"].(int64),
		TicketPriceID: parameters["ticket_price_id"].(int64),
	})
	if createResponse, err := ckService.Create(); err == nil {
		time.Sleep(1 * time.Second)
		if _, err := ckService.Patch(createResponse); err == nil {
			time.Sleep(1 * time.Second)
			ckService.Confirm()
		}
	}
}

func buildParameters() {

	masterHost := flag.String("api_host", "127.0.0.1", "Host or IP address of Go-API server. Defaults to 127.0.0.1.")
	eventID := flag.Int64("event_id", 0, "Event id")
	ticketPriceID := flag.Int64("ticket_price_id", 0, "Host or IP address of Go-API server. Defaults to 127.0.0.1.")
	mode := flag.String("mode", "slave", "Defines the operation mode. Defaults to slave")

	flag.Parse()

	if *eventID == 0 {
		panic("event_id parameter was not informed")
	}

	if *ticketPriceID == 0 {
		panic("ticket_price_id parameter was not informed")
	}

	parameters = map[string]interface{}{}
	parameters["api_host"] = *masterHost
	parameters["event_id"] = *eventID
	parameters["ticket_price_id"] = *ticketPriceID
	parameters["mode"] = *mode
}
