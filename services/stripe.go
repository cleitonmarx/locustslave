package services

import (
	"time"

	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type StripeService struct {
	stripeKey string
	client    *client.API
}

func (s *StripeService) GetTokenID(credicardNumber string) (string, error) {
	token, err := s.client.Tokens.New(&stripe.TokenParams{
		Card: &stripe.CardParams{
			Number: credicardNumber,
			Month:  "12",
			Year:   fmt.Sprintf("%d", time.Now().AddDate(3, 0, 0).Year()),
			CVC:    "123",
		},
	})
	if err != nil {
		return "", err
	}
	return token.ID, nil
}

func NewStripeService(stripeKey string) *StripeService {
	return &StripeService{
		stripeKey: stripeKey,
		client:    client.New(stripeKey, nil),
	}
}
