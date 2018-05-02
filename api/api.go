package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
)

type ProduceItem struct {
	ProduceCode string `json:"producecode"`
	Name        string `json:"name"`
	UnitPrice   string `json:"unitprice"`
}

var ProduceDB []ProduceItem

func Initialize() {
	ProduceDB = nil
	ProduceDB = append(ProduceDB, ProduceItem{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "$3.46"})
	ProduceDB = append(ProduceDB, ProduceItem{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "$2.99"})
	ProduceDB = append(ProduceDB, ProduceItem{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "$0.79"})
	ProduceDB = append(ProduceDB, ProduceItem{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "$3.59"})
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

func getAllProduce(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, ProduceDB)
}

func getProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	params["producecode"] = strings.ToUpper(params["producecode"])
	if !IsValidProduceCode(params["producecode"]) {
		http.Error(w, "Error 400 - Invalid Produce Code Format.", 400)
		return
	}

	//check if produce code exists, if exists output json and return
	for _, item := range ProduceDB {
		if item.ProduceCode == params["producecode"] {
			jsonResponse(w, http.StatusOK, item)
			return
		}
	}

	//produce code does not exist
	http.Error(w, "Error 404 - Produce Code Does Not Exist.", 404)
	return
}

func createProduceItem(w http.ResponseWriter, r *http.Request) {
	var produceItem ProduceItem

	_ = json.NewDecoder(r.Body).Decode(&produceItem)

	//are required fields empty
	if produceItem.ProduceCode == "" || produceItem.Name == "" || produceItem.UnitPrice == "" {
		http.Error(w, "Error 400 - A Produce Item Must Contain a Code, Name, And Unit Price.", 400)
		return
	}

	//are all values valid formats
	produceItem.ProduceCode = strings.ToUpper(produceItem.ProduceCode)
	validCode := IsValidProduceCode(produceItem.ProduceCode)
	validPrice := IsValidUnitPrice(produceItem.UnitPrice)
	validName := IsValidName(produceItem.Name)
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

	//check if produce code already in DB
	for _, item := range ProduceDB {
		if item.ProduceCode == produceItem.ProduceCode {
			http.Error(w, "Error 409 - Produce Code Already Exists", 409)
			return
		}
	}

	produceItem.ProduceCode = strings.ToUpper(produceItem.ProduceCode)
	ProduceDB = append(ProduceDB, produceItem)
	jsonResponse(w, http.StatusCreated, produceItem)

}

func deleteProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	params["producecode"] = strings.ToUpper(params["producecode"])
	if !IsValidProduceCode(params["producecode"]) {
		http.Error(w, "Error 400 - Invalid Produce Code Format.", 400)
		return
	}

	//check if item to delete exists
	for index, item := range ProduceDB {
		if item.ProduceCode == params["producecode"] {
			ProduceDB = append(ProduceDB[:index], ProduceDB[index+1:]...)
			jsonResponse(w, http.StatusOK, ProduceDB)
			return
		}
	}

	http.Error(w, "Error 404 - Produce Code Not Found.", 404)
}

func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/produce", getAllProduce).Methods("GET")
	router.HandleFunc("/produce/{producecode}", getProduceItem).Methods("GET")
	router.HandleFunc("/produce", createProduceItem).Methods("POST")
	router.HandleFunc("/produce/{producecode}", deleteProduceItem).Methods("DELETE")
	return router
}
