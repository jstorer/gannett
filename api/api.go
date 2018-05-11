//contains handler functions and helper functions related to the handler functions
package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
)

var prodDB DBObject     //DBObject for production database
var testDB DBObject     //DBObject for testing database
var currentDB *DBObject //DBObject to represent the currently used database

//Select which database to use depending on if tests are being run by taking in a bool variable
func Initialize(isTesting bool) {
	if isTesting {
		currentDB = &testDB
		currentDB.Data = []ProduceItem{
			{"A12T-4GH7-QPL9-3N4M", "Lettuce", "$3.46"},
			{"E5T6-9UI3-TH15-QR88", "Peach", "$2.99"},
			{"YRT6-72AS-K736-L4AR", "Green Pepper", "$0.79"},
			{"2222-2222-2222-2222", "Gala Apple", "$3.59"},
		}
	} else {
		currentDB = &prodDB
		currentDB.Data = []ProduceItem{
			{"A12T-4GH7-QPL9-3N4M", "Lettuce", "$3.46"},
			{"E5T6-9UI3-TH15-QR88", "Peach", "$2.99"},
			{"YRT6-72AS-K736-L4AR", "Green Pepper", "$0.79"},
			{"TQ4C-VV6T-75ZX-1RMR", "Gala Apple", "$3.59"},
		}
	}
}

//This function sends a request to the database to fetch all produce items through a goroutine then returns them on a
// channel then finally returns them in JSON format with a 200 status code.
func handleGetAllProduce(w http.ResponseWriter, _ *http.Request) {
	pItemSliceChnl := make(chan []ProduceItem)
	go getAllProduceItems(pItemSliceChnl) //get all items from DB
	allItems := <-pItemSliceChnl
	jsonResponse(w, http.StatusOK, allItems)
}

//This function first retrieves the produce code from the URL and determines if it is valid. If it is not valid it
//triggers a status 400 error. If it is valid it fires a goroutine to fetch that particular item and waits for a
//response via a channel. If the database returned an item it is displayed in JSON along with a 200 status code. If
//it is not found a 404 status code is triggered.
func handleGetProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	if !isValidProduceCode(params["produce_code"]) {
		http.Error(w, "error 400 - invalid produce code format", 400)
		return
	}

	pItemChnl := make(chan ProduceItem)

	go getProduceItem(params["produce_code"], pItemChnl) // get item of corresponding code from DB

	pItem := <-pItemChnl // wait for channel to return data and store it in pItem

	//if produce code not found
	if pItem.ProduceCode == "" {
		http.Error(w, "error 404 - produce code does not exist", 404)
		return

	}
	//else produce code is found
	jsonResponse(w, http.StatusOK, pItem)
	return
}

//This function first parses the JSON body request (if body does not contain valid JSON status code 400 is triggered) into a `ProduceItem` type then checks to see that all fields are
//valid and filled in by calling the `ProduceItem` method `validateProduceItem()`. If validation fails a status code
//400 is triggered along with a JSON response of the errors. If the `ProduceItem` is valid a goroutine is triggered
//to create an item with the data passed back through a channel. If the produce code already exists in the data a
//status code 409 is triggered if not a 201 status code is triggered with the JSON of the `ProduceItem` returned.
func handleCreateProduceItem(w http.ResponseWriter, r *http.Request) {
	var pItem ProduceItem

	err := json.NewDecoder(r.Body).Decode(&pItem) //get request body and decode into JSON format

	//if unable to put the body into JSON format
	if err != nil {
		http.Error(w, "error 400 - invalid JSON syntax", http.StatusBadRequest)
		return
	}

	validErrs := pItem.validateProduceItem() //check that the JSON is a valid produce item

	//if len(validErrs) > 0 then errors occurred, display them to let the user know what they are
	if len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		jsonResponse(w, http.StatusBadRequest, err)
		return
	}

	pItemChnl := make(chan ProduceItem)
	go createProduceItem(pItem, pItemChnl) //attempt to add item to the database
	pItem = <-pItemChnl                    //wait for channel to return data and store it in pItem

	if pItem.ProduceCode == "" {
		http.Error(w, "error 409 - produce code already exists", 409)
		return
	}

	jsonResponse(w, http.StatusCreated, pItem)

}

//This function checks if the produce code passed in from the URL is valid and the body contains valid JSON, if it does not a status code 400 is triggered.
//If it is the JSON from the request body is placed into a `ProduceItem`. This JSON is then validated by calling the
//`ProduceItem` method `validationProduceItem()`. Upon validation success a go routine is called and passes the updated
//item back through a channel. If the produce code was not found a status code 404 is triggered or if the changed
//produce code already exists a status 409 is triggered. Otherwise a status 200 is triggered and the updated item
//contents are returned as a JSON.
func handleUpdateProduceItem(w http.ResponseWriter, r *http.Request) {
	var pItem ProduceItem
	params := mux.Vars(r)

	//check if produce code format is valid
	params["produce_code"] = strings.ToUpper(params["produce_code"])
	if !isValidProduceCode(params["produce_code"]) {
		http.Error(w, "error 400 - invalid produce code format", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&pItem) //get request body and decode into JSON format

	//if unable to put the body into JSON format
	if err != nil {
		http.Error(w, "error 400 - invalid JSON syntax", http.StatusBadRequest)
		return
	}

	//if len(validErrs) > 0 then errors occurred, display them to let the user know what they are
	validErrs := pItem.validateProduceItem()
	if len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		jsonResponse(w, http.StatusBadRequest, err)
		return
	}

	pItemChnl := make(chan ProduceItem)
	go updateProduceItem(params["produce_code"], pItem, pItemChnl) //update item of given produce code in DB
	pItem = <-pItemChnl                                            //wait for channel to return data and store in pItem

	//produce code not found
	if pItem.ProduceCode == "" {
		http.Error(w, "error 404 - produce code does not exist", 404)
		return
	}
	//new produce code value already exists
	if pItem.ProduceCode == "0" {
		http.Error(w, "error 409 - updated produce code value already exists", 409)
		return
	}

	//item updated successfully, display its new contents
	jsonResponse(w, http.StatusOK, pItem)

}

//This function first checks if the produce code passed in from the URL is valid, if it is not a status code 400 is
//triggered. If the produce code is valid a goroutine is triggered and passes the produce item back through a channel.
//If the code was not found a status 404 is triggered, if it was found a status 200 is triggered and the deleted produce
//item is returned as a JSON.
func handleDeleteProduceItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//check if produce code format is valid
	if !isValidProduceCode(params["produce_code"]) {
		http.Error(w, "error 400 - invalid produce code format", 400)
		return
	}

	var pItem ProduceItem
	pItemChnl := make(chan ProduceItem)
	go deleteProduceItem(params["produce_code"], pItemChnl) //delete item from DB
	pItem = <-pItemChnl                                     //wait for item to return on channel

	//if code not found
	if pItem.ProduceCode == "" {
		http.Error(w, "error 404 - produce code not found.", 404)
		return
	}

	//code found
	jsonResponse(w, http.StatusOK, pItem)
}

//This function accepts a produce string and validates via the regex expression "^[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}$"
//to determine if it is valid or not and returns true if valid or false if not. This expression checks that code is
//four groups of four alphanumeric characters.
func isValidProduceCode(produceCode string) bool {
	match, _ := regexp.MatchString(`^[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}$`, produceCode)
	return match
}

//This function accepts a unit price string and validates via the regex expression
//"^\$(([1-9]\d{0,2}(,\d{3})*)|(([1-9]\d*)?\d))(\.\d\d?)?$" which requires a dollar sign followed by numbers with or
//without correct comma seperation but not incorrect comma seperation and at most 2 trailing decimals. If valid returns
//true and if not valid returns false.
func isValidUnitPrice(unitPrice string) bool {
	match, _ := regexp.MatchString(`^\$(([1-9]\d{0,2}(,\d{3})*)|(([1-9]\d*)?\d))(\.\d\d?)?$`, unitPrice)
	return match
}

//This function accepts a name string and validates via the regex expression `\w+(?: \w+)*$` which will allow no
//whitespace before or trailing whitespace after a single space set of alphanumeric characters. If valid returns
//true if not valid returns false.
func isValidName(name string) bool {
	match, _ := regexp.MatchString(`^\w+(?: \w+)*$`, name)
	return match
}

//This function accepts a ResponseWriter, status code, and payload and then JSON encodes it via json.Marshall()
//and then writes the corresponding body and headers in JSON format to be displayed.
func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload) //convert payload to JSON

	//if things got this far and invalid JSON syntax is presented something unexpected has happened
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
