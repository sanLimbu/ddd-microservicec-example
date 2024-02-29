package recommendation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Rhymond/go-money"
)

type parternShipResponse struct {
	AvailableHotels []struct {
		Name         string `json:"name"`
		PriceInNight int    `json:"pricePerNight"`
	} `json:"availableHotels"`
}

type parternShipAdaptor struct {
	client *http.Client
	url    string
}

func NewPartnerShipAdaptor(client *http.Client, url string) (*parternShipAdaptor, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}
	if url == "" {
		return nil, errors.New("url cannot be empty")
	}
	return &parternShipAdaptor{client: client, url: url}, nil
}

func (p parternShipAdaptor) GetAvailability(ctx context.Context, tripStart time.Time, tripEnd time.Time, location string) ([]Option, error) {
	from := fmt.Sprintf("%d-%d-%d", tripStart.Year(), tripStart.Month(), tripStart.Day())
	to := fmt.Sprintf("%d-%d-%d", tripEnd.Year(), tripEnd.Month(), tripEnd.Day())

	url := fmt.Sprintf("%s/partnerships?location=%s&from=%s&to=%s", p.url, location, from, to)
	res, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call partnership: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad request to partnerships: %d", res.StatusCode)
	}

	var response parternShipResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("couldn't decode the response body of parnership: %w", err)
	}

	opts := make([]Option, len(response.AvailableHotels))
	for i, p := range response.AvailableHotels {
		opts[i] = Option{
			HotelName:     p.Name,
			Location:      location,
			PricePerNight: *money.New(int64(p.PriceInNight), "GBP"),
		}
	}
	return opts, nil
}
