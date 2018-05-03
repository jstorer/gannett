package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type ProduceItem struct {
	ProduceCode string `json:"produce_code"`
	Name        string `json:"name"`
	UnitPrice   string `json:"unit_price"`
}

var ProduceDB DBObject

type DBObject struct {
	mu   sync.Mutex
	Data []ProduceItem
}

func Initialize() {
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
	match, _ := regexp.MatchString(`^[A-Za-z0-9]+.*[A-Za-z0-9]+$`, name)
	return match
}

func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func readAllProduceItems(sendchannel chan<- []ProduceItem) {
	sendchannel <- ProduceDB.Data
}

func readProduceItem(pCode string, sendchannel chan<- ProduceItem) {
	for _, item := range ProduceDB.Data {
		if item.ProduceCode == pCode {
			sendchannel <- item
		}
	}
	sendchannel <- ProduceItem{}
}

func writeProduceItem(pItem ProduceItem, pItemChnl chan ProduceItem) {
	ProduceDB.mu.Lock()
	for _, item := range ProduceDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			ProduceDB.mu.Unlock()
			pItemChnl <- ProduceItem{}
			return
		}
	}
	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	ProduceDB.Data = append(ProduceDB.Data, pItem)
	ProduceDB.mu.Unlock()
	pItemChnl <- pItem
}

func removeProduceItem(pItem ProduceItem, pItemChnl chan ProduceItem) {
	ProduceDB.mu.Lock()
	for index, item := range ProduceDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			ProduceDB.Data = append(ProduceDB.Data[:index], ProduceDB.Data[index+1:]...)
			ProduceDB.mu.Unlock()
			pItemChnl <- pItem
			return
		}
	}
	ProduceDB.mu.Unlock()
	pItemChnl <- ProduceItem{}
}

func getAllProduce(w http.ResponseWriter, _ *http.Request) {
	chnl := make(chan []ProduceItem)
	go readAllProduceItems(chnl)
	jsonResponse(w, http.StatusOK, <-chnl)
}

func getProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	params["produce_code"] = strings.ToUpper(params["produce_code"])
	if !IsValidProduceCode(params["produce_code"]) {
		http.Error(w, "Error 400 - Invalid Produce Code Format.", 400)
		return
	}

	chnl := make(chan ProduceItem)

	go readProduceItem(params["produce_code"], chnl)

	pItem := <-chnl

	//if produce item found
	if pItem.ProduceCode != "" {
		jsonResponse(w, http.StatusOK, pItem)
		return
	}

	//produce code does not exist
	http.Error(w, "Error 404 - Produce Code Does Not Exist.", 404)
	return
}

func createProduceItem(w http.ResponseWriter, r *http.Request) {
	var pItem ProduceItem

	_ = json.NewDecoder(r.Body).Decode(&pItem)

	//are required fields empty
	if pItem.ProduceCode == "" || pItem.Name == "" || pItem.UnitPrice == "" {
		http.Error(w, "Error 400 - A Produce Item Must Contain a Code, Name, And Unit Price.", 400)
		return
	}

	//are all values valid formats
	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	validCode := IsValidProduceCode(pItem.ProduceCode)
	validPrice := IsValidUnitPrice(pItem.UnitPrice)
	validName := IsValidName(pItem.Name)
	if !validCode || !validPrice || !validName {
		if !validCode && !validPrice && !validName {
			http.Error(w, "Error 400 - Invalid Produce Code, Name, And Unit Price Format.", 400)
			return
		}
		if !validCode && !validPrice {
			http.Error(w, "Error 400 - Invalid Produce Code And Unit Price Format.", 400)
			return
		}
		if !validCode && !validName {
			http.Error(w, "Error 400 - Invalid Produce Code And Name Format.", 400)
			return
		}
		if !validPrice && !validName {
			http.Error(w, "Error 400 - Name And Unit Price Format.", 400)
			return
		}
		if !validCode {
			http.Error(w, "Error 400 - Invalid Produce Code Format.", 400)
			return
		}
		if !validPrice {
			http.Error(w, "Error 400 - Invalid Unit Price Format.", 400)
			return
		}
		if !validName {
			http.Error(w, "Error 400 - Invalid Name Format.", 400)
			return
		}

	}

	pItemChnl := make(chan ProduceItem)
	go writeProduceItem(pItem, pItemChnl)
	pItem = <-pItemChnl

	if pItem.ProduceCode == "" {
		http.Error(w, "Error 409 - Produce Code Already Exists", 409)
		return
	}

	jsonResponse(w, http.StatusCreated, pItem)

}

func deleteProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	params["produce_code"] = strings.ToUpper(params["produce_code"])
	if !IsValidProduceCode(params["produce_code"]) {
		http.Error(w, "Error 400 - Invalid Produce Code Format.", 400)
		return
	}

	var pItem ProduceItem
	pItem.ProduceCode = params["produce_code"]
	pItemChnl := make(chan ProduceItem)
	go removeProduceItem(pItem, pItemChnl)
	pItem = <-pItemChnl

	if pItem.ProduceCode == "" {
		http.Error(w, "Error 404 - Produce Code Not Found.", 404)
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
