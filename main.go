package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type ProduceItem struct {
	ProduceCode string `json:"producecode"`
	Name        string `json:"name"`
	UnitPrice   string `json:"unitprice"`
}

var produceDB []ProduceItem

func getAllProduce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produceDB)
}

func getProduceItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range produceDB {
		if item.ProduceCode == params["producecode"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createProduceItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var produceItem ProduceItem
	_ = json.NewDecoder(r.Body).Decode(&produceItem)
	for _, item := range produceDB {
		if item.ProduceCode == produceItem.ProduceCode {
			println("item already exists.")
			return
		}
	}
	produceItem.ProduceCode = strings.ToUpper(produceItem.ProduceCode)
	produceDB = append(produceDB, produceItem)
	json.NewEncoder(w).Encode(produceItem)
}

func deleteProduceItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range produceDB {
		if item.ProduceCode == params["producecode"] {
			produceDB = append(produceDB[:index], produceDB[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(produceDB)
	}
}

func main() {
	router := mux.NewRouter()
	produceDB = append(produceDB, ProduceItem{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "$3.46"})
	produceDB = append(produceDB, ProduceItem{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "$2.99"})
	produceDB = append(produceDB, ProduceItem{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "$0.79"})
	produceDB = append(produceDB, ProduceItem{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "$3.59"})

	router.HandleFunc("/produce", getAllProduce).Methods("GET")
	router.HandleFunc("/produce/{producecode}", getProduceItem).Methods("GET")
	router.HandleFunc("/produce", createProduceItem).Methods("POST")
	router.HandleFunc("/produce/{producecode}", deleteProduceItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}
