package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/tax/v2/rental", PostRental).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func PostRental(w http.ResponseWriter, r *http.Request) {

	var rentalRequest RentalRequest
	_ = json.NewDecoder(r.Body).Decode(&rentalRequest)

	var rentalResponse = RentalResponse{ID: "1"}
	json.NewEncoder(w).Encode(rentalResponse)
}

type RentalRequest struct {
}

type RentalResponse struct {
	ID string `json:"id,omitempty"`
}
