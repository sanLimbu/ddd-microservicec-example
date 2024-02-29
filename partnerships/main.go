package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Res struct {
	AvailableHotels []struct {
		Name          string `json:"name"`
		PricePerNight int    `json:"pricePerNight"`
	} `json:"availableHotels"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 10

	sampleRes := Res{AvailableHotels: []struct {
		Name          string `json:"name"`
		PricePerNight int    `json:"pricePerNight"`
	}{
		{
			Name:          "some hotel",
			PricePerNight: 300,
		},
		{
			Name:          "some other hotel",
			PricePerNight: 30,
		},
		{
			Name:          "some third hotel",
			PricePerNight: 90,
		},
		{
			Name:          "some fourth hotel",
			PricePerNight: 80,
		},
	}}

	b, err := json.Marshal(sampleRes)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.
		Path("/parternships").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ran := rand.Intn(max - min + 1)
			if ran > 7 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(b)

		})
	log.Println("running")
	if err := http.ListenAndServe(":3031", r); err != nil {
		log.Fatal(err)
	}
}
