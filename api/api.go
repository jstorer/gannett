package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"net/url"
)

type ProduceItem struct {
	ProduceCode string `json:"produce_code"`
	Name        string `json:"name"`
	UnitPrice   string `json:"unit_price"`
}

var ProduceDB DBObject

type DBObject struct {
	mu   sync.RWMutex
	Data []ProduceItem
}


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

func readAllProduceItems(sendchannel chan<- []ProduceItem) {
	ProduceDB.mu.RLock()
	defer ProduceDB.mu.RUnlock()
	sendchannel <- ProduceDB.Data
}

func readProduceItem(pCode string, sendchannel chan<- ProduceItem) {
	ProduceDB.mu.RLock()
	defer ProduceDB.mu.RUnlock()
	for _, item := range ProduceDB.Data {
		if item.ProduceCode == pCode {
			sendchannel <- item
			return
		}
	}
	sendchannel <- ProduceItem{}
}

func writeProduceItem(pItem ProduceItem, pItemChnl chan ProduceItem) {
	ProduceDB.mu.Lock()
	defer ProduceDB.mu.Unlock()
	for _, item := range ProduceDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			pItemChnl <- ProduceItem{}
			return
		}
	}
	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	ProduceDB.Data = append(ProduceDB.Data, pItem)
	pItemChnl <- pItem
}

func removeProduceItem(pItem ProduceItem, pItemChnl chan ProduceItem) {
	ProduceDB.mu.Lock()
	defer ProduceDB.mu.Unlock()

	for index, item := range ProduceDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			ProduceDB.Data = append(ProduceDB.Data[:index], ProduceDB.Data[index+1:]...)
			pItemChnl <- pItem
			return
		}
	}
	pItemChnl <- ProduceItem{}
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

func (pItem *ProduceItem) validateProduceItem() url.Values{
	errs := url.Values{}
	//all required fields exist
	if pItem.ProduceCode == "" {
		errs.Add("produce_code","produce field is required")
	}

	if pItem.Name == ""{
		errs.Add("name","name field is required")
	}

	if pItem.UnitPrice ==""{
		errs.Add("unit_price","unit price field is required")
	}
	//are all values valid formats
	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	if !IsValidProduceCode(pItem.ProduceCode) {
		errs.Add("produce_code","invalid produce code format")
	}

	if !IsValidName(pItem.Name){
		errs.Add("name","invalid name format")
	}

	if !IsValidUnitPrice(pItem.UnitPrice){
		errs.Add("unit_price", "invalid unit price format")
	}
	return errs
}

func createProduceItem(w http.ResponseWriter, r *http.Request) {
	var pItem ProduceItem
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&pItem); err != nil{
		panic(err)
	}

	validErrs := pItem.validateProduceItem()
	if len(validErrs) > 0 {
		err:=map[string]interface{}{"validationError": validErrs}
		jsonResponse(w,http.StatusBadRequest,err)
		return
	}

	pItemChnl := make(chan ProduceItem)
	go writeProduceItem(pItem, pItemChnl)
	pItem = <-pItemChnl

	if pItem.ProduceCode == "" {
		http.Error(w, "error 409 - produce code already exists", 409)
		return
	}

	jsonResponse(w, http.StatusCreated, pItem)

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

func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/produce", getAllProduce).Methods("GET")
	router.HandleFunc("/produce/{produce_code}", getProduceItem).Methods("GET")
	router.HandleFunc("/produce", createProduceItem).Methods("POST")
	router.HandleFunc("/produce/{produce_code}", deleteProduceItem).Methods("DELETE")
	return router
}
