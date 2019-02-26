package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
)

type RequestLineItem struct {
	LineItemID  int     `json:"lineItemId"`
	GrossAmount float64 `json:"grossAmount"`
	Quantity    int     `json:"quantity"`
	Category    string  `json:"category"`
	SubCategory string  `json:"subCategory"`
}

type RentalRequest struct {
	TransactionDate string            `json:"transactionDate"`
	StoreNumber     string            `json:"storeNumber"`
	LineItems       []RequestLineItem `json:"lineItems"`
}

type ResponseLineItem struct {
	LineItemID int     `json:"lineItemId"`
	TaxAmount  float64 `json:"taxAmount"`
}

type RentalResponse struct {
	TotalTaxAmount float64            `json:"totalTaxAmount"`
	LineItems      []ResponseLineItem `json:"lineItems"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/tax/v2/rental", PostRental).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func PostRental(w http.ResponseWriter, r *http.Request) {

	var rentalRequest RentalRequest
	_ = json.NewDecoder(r.Body).Decode(&rentalRequest)

	resLines := make([]ResponseLineItem, 0)

	var totalTax = 0.00
	for _, reqLine := range rentalRequest.LineItems {
		lineTax := math.Round(reqLine.GrossAmount*float64(reqLine.Quantity)*0.06*100) / 100
		lineNum := reqLine.LineItemID
		newLine := ResponseLineItem{LineItemID: lineNum, TaxAmount: lineTax}
		resLines = append(resLines, newLine)
		totalTax += lineTax
	}

	var rentalResponse RentalResponse
	rentalResponse.LineItems = resLines
	rentalResponse.TotalTaxAmount = math.Round(totalTax*100) / 100
	_ = json.NewEncoder(w).Encode(rentalResponse)
}
