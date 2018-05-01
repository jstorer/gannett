package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ProduceItem struct {
	ProduceCode string `json:"producecode"`
	Name        string `json:"name"`
	UnitPrice   string `json:"unitprice"`
}

var produceDB []ProduceItem

func isValidProduceCode(produceCode string) bool {
	match, _ := regexp.MatchString(`[\d\w]{4}\-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}`, produceCode)
	return match
}

func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func getAllProduce(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, produceDB)
}

func getProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	params["producecode"] = strings.ToUpper(params["producecode"])
	if !isValidProduceCode(params["producecode"]) {
		http.Error(w, "Error 400 - Invalid Produce Code Format.", 400)
		return
	}

	//check if produce code exists
	for _, item := range produceDB {
		if item.ProduceCode == params["producecode"] {
			jsonResponse(w, http.StatusOK, item)
			return
		}
	}

	//produce code does not exist
	http.Error(w, "Error 404 - Produce Code Does Not Exist.", 404)

}

func createProduceItem(w http.ResponseWriter, r *http.Request) {
	var produceItem ProduceItem

	_ = json.NewDecoder(r.Body).Decode(&produceItem)

	produceItem.ProduceCode = strings.ToUpper(produceItem.ProduceCode)

	//check if produce code already in DB
	for _, item := range produceDB {
		if item.ProduceCode == produceItem.ProduceCode {
			http.Error(w, "Error 409 - Produce Code Already Exists", 409)
			return
		}
	}

	produceItem.ProduceCode = strings.ToUpper(produceItem.ProduceCode)
	produceDB = append(produceDB, produceItem)
	jsonResponse(w, http.StatusOK, produceItem)

}

func deleteProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	params["producecode"] = strings.ToUpper(params["producecode"])
	if !isValidProduceCode(params["producecode"]) {
		http.Error(w, "Error 400 - Invalid Produce Code Format.", 400)
		return
	}

	//check if item to delete exists
	for index, item := range produceDB {
		if item.ProduceCode == params["producecode"] {
			produceDB = append(produceDB[:index], produceDB[index+1:]...)
			jsonResponse(w, http.StatusOK, produceDB)
			return
		}
	}

	http.Error(w, "Error 404 - Produce Code Not Found.", 404)
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
