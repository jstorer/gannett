package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
)

var ProduceDB DBObject

func Initialize() {
	//init DB to indicated default values
	ProduceDB.Data = nil
	ProduceDB.Data = append(ProduceDB.Data, ProduceItem{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "$3.46"})
	ProduceDB.Data = append(ProduceDB.Data, ProduceItem{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "$2.99"})
	ProduceDB.Data = append(ProduceDB.Data, ProduceItem{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "$0.79"})
	ProduceDB.Data = append(ProduceDB.Data, ProduceItem{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "$3.59"})
}

func IsValidProduceCode(produceCode string) bool {
	match, _ := regexp.MatchString(`^[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}$`, produceCode)
	return match
}

func IsValidUnitPrice(unitPrice string) bool {
	match, _ := regexp.MatchString(`^\$(([1-9]\d{0,2}(,\d{3})*)|(([1-9]\d*)?\d))(\.\d\d?)?$`, unitPrice)
	return match
}

func IsValidName(name string) bool {
	match, _ := regexp.MatchString(`^\w+(?: \w+)*$`, name)
	return match
}

func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func getAllProduce(w http.ResponseWriter, _ *http.Request) {
	chnl := make(chan []ProduceItem)
	go readAllProduceItems(chnl)
	allItems := <-chnl
	jsonResponse(w, http.StatusOK, allItems)
}

func getProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	params["produce_code"] = strings.ToUpper(params["produce_code"])
	if !IsValidProduceCode(params["produce_code"]) {
		http.Error(w, "error 400 - invalid produce code format", 400)
		return
	}

	chnl := make(chan ProduceItem)

	go readProduceItem(params["produce_code"], chnl)

	pItem := <-chnl

	if pItem.ProduceCode != "" {
		jsonResponse(w, http.StatusOK, pItem)
		return
	}

	http.Error(w, "error 404 - produce code does not exist", 404)
	return
}

func createProduceItem(w http.ResponseWriter, r *http.Request) {
	var pItem ProduceItem

	if err := json.NewDecoder(r.Body).Decode(&pItem); err != nil {
		panic(err)
	}

	validErrs := pItem.validateProduceItem()
	if len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		jsonResponse(w, http.StatusBadRequest, err)
		return
	}

	pItemChnl := make(chan ProduceItem)
	go writeNewProduceItem(pItem, pItemChnl)
	pItem = <-pItemChnl

	if pItem.ProduceCode == "" {
		http.Error(w, "error 409 - produce code already exists", 409)
		return
	}

	jsonResponse(w, http.StatusCreated, pItem)

}

func updateProduceItem(w http.ResponseWriter, r *http.Request) {
	var pItem ProduceItem
	params := mux.Vars(r)

	//check if produce code format is valid
	params["produce_code"] = strings.ToUpper(params["produce_code"])
	if !IsValidProduceCode(params["produce_code"]) {
		http.Error(w, "error 400 - invalid produce code format", 400)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&pItem); err != nil {
		panic(err)
	}

	validErrs := pItem.validateProduceItem()
	if len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		jsonResponse(w, http.StatusBadRequest, err)
		return
	}

	pItemChnl := make(chan ProduceItem)
	go writeUpdateProduceItem(params["produce_code"], pItem, pItemChnl)
	pItem = <-pItemChnl

	if pItem.ProduceCode == "" {
		http.Error(w, "error 404 - produce code does not exist", 404)
		return
	}

	if pItem.ProduceCode == "0" {
		http.Error(w, "error 409 - updated produce code value already exists", 409)
		return
	}

	jsonResponse(w, http.StatusOK, pItem)

}

func deleteProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	params["produce_code"] = strings.ToUpper(params["produce_code"])
	if !IsValidProduceCode(params["produce_code"]) {
		http.Error(w, "error 400 - invalid produce code format", 400)
		return
	}

	var pItem ProduceItem
	pItem.ProduceCode = params["produce_code"]
	pItemChnl := make(chan ProduceItem)
	go removeProduceItem(pItem, pItemChnl)
	pItem = <-pItemChnl

	if pItem.ProduceCode == "" {
		http.Error(w, "error 404 - produce code not found.", 404)
		return
	}

	jsonResponse(w, http.StatusOK, pItem)
}
