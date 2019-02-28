package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/mux"
)

const AuthToken = "FE1CBBBE-653B-4594-A5AD-67B5256D2FEB"

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

func postRental(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")

	if AuthToken != authToken {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(rentalResponse)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/tax/v2/rental", postRental).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
