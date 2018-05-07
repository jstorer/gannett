package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"io/ioutil"
)


var (
	server     *httptest.Server
	reader     io.Reader
	produceUrl string
)

func init() {
	server = httptest.NewServer(Handlers())
	produceUrl = fmt.Sprintf("%s/api/produce", server.URL)
	TestingMode = true
}

func initTest() {
	TestDB.Data = nil
	TestDB.Data = []ProduceItem{
		{"A12T-4GH7-QPL9-3N4M","Lettuce","$3.46"},
		{"E5T6-9UI3-TH15-QR88","Peach","$2.99"},
		{"YRT6-72AS-K736-L4AR","Green Pepper","$0.79"},
		{"TQ4C-VV6T-75ZX-1RMR","Gala Apple","$3.59"},
	}
}

func TestIsValidProduceCode(t *testing.T) {
	var testCodes = []struct{
		value string
		valid bool
	}{
		{"abcd-1234-z9y8-q123", true},
		{"abcd-1234-Z9Y8-q12W", true},
		{"1111-2222-3333-4444", true},
		{"aaaa-bbbb-cccc-dddd", true},
		{"AAAA-BBBB-CCCC-DDDD", true},
		{"", false},
		{"a", false},
		{"ab", false},
		{"abc", false},
		{"abcd", false},
		{"abcd-", false},
		{"abcd-1", false},
		{"abcd-12", false},
		{"abcd-123", false},
		{"abcd-1234", false},
		{"abcd-1234-", false},
		{"abcd-1234-z", false},
		{"abcd-1234-z9", false},
		{"abcd-1234-z9y", false},
		{"abcd-1234-z9y8", false},
		{"abcd-1234-z9y8-q", false},
		{"abcd-1234-z9y8-q1", false},
		{"abcd-1234-z9y8-q13", false},
		{"abcd-1234-z9y8-q123a", false},
		{"abcd-1234-z9y80-q123", false},
		{"abcd-12340-z9y8-q123", false},
		{"abcd0-1234-z9y8-q123", false},
		{"abcd-1234-z9y8-q123-", false},
		{"abcd-1234-z9y8-q123-abcd", false},
	}

	for _, item := range testCodes {
		assert.Equal(t, item.valid, IsValidProduceCode(item.value), fmt.Sprintf("On Produce Code: `%s`", item.value))
	}
}

func TestIsValidUnitPrice(t *testing.T) {
	var testPrices = []struct {
		value string
		valid bool
	}{
		{"$0.1", true},
		{"$0.01", true},
		{"$1", true},
		{"$1.0", true},
		{"$1.00", true},
		{"$0.10", true},
		{"$4231", true},
		{"$4,006", true},
		{"$4,000.5", true},
		{"$4,000.93", true},
		{"$4,000,001.23", true},
		{"$4,000.00", true},
		{"", false},
		{"0", false},
		{"5", false},
		{"5.23", false},
		{"5.2", false},
		{"$", false},
		{"$01", false},
		{"$01.50", false},
		{"$21,12345", false},
		{"$5.123", false},
	}

	for _, item := range testPrices {
		assert.Equal(t, item.valid, IsValidUnitPrice(item.value), fmt.Sprintf("On Unit Price: `%s`", item.value))
	}

}

func TestIsValidName(t *testing.T) {
	var testNames = []struct {
		value string
		valid bool
	}{
		{"A", true},
		{"A A", true},
		{"Milk", true},
		{"Cow Milk", true},
		{"Cow Milk 12", true},
		{"cow milk 12", true},
		{"", false},
		{" ", false},
		{" a", false},
		{"m!lk", false},
	}
	for _, item := range testNames {
		assert.Equal(t, item.valid, IsValidName(item.value), fmt.Sprintf("On Unit Price: `%s`", item.value))
	}
}

func TestGetAllProduce(t *testing.T) {
	initTest()
	var getAllTests = []struct{
		desc string
		method string
		path string
		statusCode int
		expectedBody string
	}{
		{"all items payload from initTest()","GET",produceUrl,200,`[{"produce_code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","unit_price":"$3.46"},{"produce_code":"E5T6-9UI3-TH15-QR88","name":"Peach","unit_price":"$2.99"},{"produce_code":"YRT6-72AS-K736-L4AR","name":"Green Pepper","unit_price":"$0.79"},{"produce_code":"TQ4C-VV6T-75ZX-1RMR","name":"Gala Apple","unit_price":"$3.59"}]`},
	}

	for _, item := range getAllTests{
		request, err := http.NewRequest(item.method, item.path, nil)
		response, err := http.DefaultClient.Do(request)

		if err != nil {
			t.Error(err)
		}

		responseData,_:=ioutil.ReadAll(response.Body)
		assert.Equal(t,item.expectedBody,string(responseData),fmt.Sprintf("unexpected response for %s",item.desc))
		assert.Equal(t,item.statusCode,response.StatusCode,"unexpected status code")
	}

}

func TestGetProduceItem(t *testing.T) {
	initTest()
	//VALID TEST----------------------------------------------------------
	//Get item
	ProduceDB.Data = append(ProduceDB.Data, ProduceItem{ProduceCode: "ABCD-1234-EFGH-5678", Name: "Black Beans", UnitPrice: "$2.25"})
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/ABCD-1234-EFGH-5678", produceUrl), nil)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("200 OK expected but %d returned", response.StatusCode)
	}

	//INVALID TESTS------------------------------------------------------
	//Invalid Produce Code
	request, err = http.NewRequest("GET", fmt.Sprintf("%s/ABCDe-1234-EFGH-5678", produceUrl), nil)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Produce code valid but does not exist
	request, err = http.NewRequest("GET", fmt.Sprintf("%s/ABCD-1234-EFGH-0000", produceUrl), nil)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 404 {
		t.Errorf("404 Not found expected but %d returned", response.StatusCode)
	}

}

func TestUpdateProduceItem(t *testing.T) {
	initTest()

	//valid
	//unchanging produce code
	produceItemJson := `{"produce_code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","unit_price":"$3.46"}`
	reader = strings.NewReader(produceItemJson)
	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/A12T-4GH7-QPL9-3N4M", produceUrl), reader)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("200 Created expected but %d returned", response.StatusCode)
	}

	//check if values push to DB
	produceItemJson = `{"produce_code":"A12T-4GH7-QPL9-1111","name":"Cheese","unit_price":"$5.00"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("PUT", fmt.Sprintf("%s/A12T-4GH7-QPL9-3N4M", produceUrl), reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if ProduceDB.Data[0].ProduceCode != "A12T-4GH7-QPL9-1111" || ProduceDB.Data[0].Name != "Cheese" || ProduceDB.Data[0].UnitPrice != "$5.00" {
		t.Errorf("failed to push update to DB")
	}

	if response.StatusCode != 200 {
		t.Errorf("200 Created expected but %d returned", response.StatusCode)
	}

	//invalid
	//updated code already exists
	initTest()
	produceItemJson = `{"produce_code":"A12T-4GH7-QPL9-3N4M","name":"Cheese","unit_price":"$5.00"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("PUT", fmt.Sprintf("%s/E5T6-9UI3-TH15-QR88", produceUrl), reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 409 {
		t.Errorf("409 conflict expected but %d returned", response.StatusCode)
	}

	//Produce code doesnt exist to update
	produceItemJson = `{"produce_code":"A12T-4GH7-QPL9-3N4A","name":"Cheese","unit_price":"$5.00"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("PUT", fmt.Sprintf("%s/E5T6-9UI3-TH15-1111", produceUrl), reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 404 {
		t.Errorf("404 produce code not found expected but %d returned", response.StatusCode)
	}

	//invalid produce code end point
	produceItemJson = `{"produce_code":"A12T-4GH7-QPL9-3N4A","name":"Cheese","unit_price":"$5.00"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("PUT", fmt.Sprintf("%s/E5T6-9UI3-TH15-111", produceUrl), reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 bad request expected but %d returned", response.StatusCode)
	}

	//bad payload
	produceItemJson = `{"produce_code":"A12T-4GH7-QPL9-3NM","name":"Cheese","unit_price":"$5.00"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("PUT", fmt.Sprintf("%s/E5T6-9UI3-TH15-QR88", produceUrl), reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 bad request expected but %d returned", response.StatusCode)
	}

}

func TestCreateProduceItem(t *testing.T) {
	initTest()
	//VALID TEST----------------------------------------------------------
	//Check if valid produce item created
	produceItemJson := `{"produce_code":"abcd-1234-EFGH-5I6J","name":"Cheese","unit_price":"$9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err := http.NewRequest("POST", produceUrl, reader)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 201 {
		t.Errorf("201 Created expected but %d returned", response.StatusCode)
	}

	produceItemJson = `{"produce_code":"abCd-12D4-2FGH-5i6J","name":"Lemon Grass","unit_price":"$2.1"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 201 {
		t.Errorf("201 Created expected but %d returned", response.StatusCode)
	}

	//INVALID TESTS------------------------------------------------------
	//left produce code blank
	produceItemJson = `{"produce_code":"","name":"Cheese","unit_price":"$9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//left name blank
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6J","name":"","unit_price":"$9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//left unit price blank
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6J","name":"Cheese","unit_price":""}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//All bad formats
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6Ja","name":"$Cheese","unit_price":"$09.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Bad produce code and price
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6Ja","name":"Cheese","unit_price":"9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Bad produce code and name
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6Ja","name":" Cheese","unit_price":"$9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Bad price and name format
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6J","name":"@Cheese","unit_price":"$09.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Bad produce code format
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6Ja","name":"Cheese","unit_price":"$9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Bad unit price format
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6J","name":"Cheese","unit_price":"9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Bad name format
	produceItemJson = `{"produce_code":"abcd-1234-EFGH-5I6J","name":" Cheese","unit_price":"$9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Produce code (case insensitive) already exists in DB
	produceItemJson = `{"produce_code":"ABCD-1234-EFGH-5I6J","name":"Cheese","unit_price":"$9.99"}`
	reader = strings.NewReader(produceItemJson)
	request, err = http.NewRequest("POST", produceUrl, reader)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 409 {
		t.Errorf("409 Item already exists expected but %d returned", response.StatusCode)
	}

}

func TestDeleteProduceItem(t *testing.T) {
	initTest()
	//VALID TEST----------------------------------------------------------
	//Delete item
	ProduceDB.Data = append(ProduceDB.Data, ProduceItem{ProduceCode: "ABCD-1234-EFGH-5678", Name: "Black Beans", UnitPrice: "$2.25"})
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/ABCD-1234-EFGH-5678", produceUrl), nil)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("200 OK expected but %d returned", response.StatusCode)
	}

	//INVALID TESTS------------------------------------------------------
	//Invalid Produce Code
	request, err = http.NewRequest("DELETE", fmt.Sprintf("%s/ABCDe-1234-EFGH-5678", produceUrl), nil)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("400 Bad request expected but %d returned", response.StatusCode)
	}

	//Produce code valid but does not exist, checking item just added also verifies it was removed from DB
	request, err = http.NewRequest("DELETE", fmt.Sprintf("%s/ABCD-1234-EFGH-5678", produceUrl), nil)
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 404 {
		t.Errorf("404 Not found expected but %d returned", response.StatusCode)
	}

}

func TestValidateProduceItem(t *testing.T) {
	//Valid
	var pItem ProduceItem
	pItem.ProduceCode = "1111-1111-1111-1111"
	pItem.Name = "milk"
	pItem.UnitPrice = "$1.00"
	errs := pItem.validateProduceItem()

	if len(errs) != 0 {
		t.Errorf("0 errors expected but %d returned", len(errs))
	}

	//Invalid all empty fields
	pItem = ProduceItem{}
	pItem.ProduceCode = ""
	pItem.Name = ""
	pItem.UnitPrice = ""
	errs = pItem.validateProduceItem()

	if len(errs) != 3 {
		t.Errorf("3 errors expected but %d returned", len(errs))
	}

	//Invalid name and unit price empty fields
	pItem = ProduceItem{}
	pItem.ProduceCode = "aaaa-bbbb-1111-2222"
	pItem.Name = ""
	pItem.UnitPrice = ""
	errs = pItem.validateProduceItem()

	if len(errs) != 2 {
		t.Errorf("2 errors expected but %d returned", len(errs))
	}

	//Invalid unit price empty fields
	pItem = ProduceItem{}
	pItem.ProduceCode = "aaaa-bbbb-1111-2222"
	pItem.Name = "milk"
	pItem.UnitPrice = ""
	errs = pItem.validateProduceItem()

	if len(errs) != 1 {
		t.Errorf("1 error expected but %d returned", len(errs))
	}

	//Invalid produce code format
	pItem = ProduceItem{}
	pItem.ProduceCode = "1111-1111-1111"
	pItem.Name = "milk"
	pItem.UnitPrice = "$1.00"
	errs = pItem.validateProduceItem()

	if len(errs) != 1 {
		t.Errorf("1 errors expected but %d returned", len(errs))
	}

	//Invalid name format
	pItem = ProduceItem{}
	pItem.ProduceCode = "1111-1111-1111-1111"
	pItem.Name = "m@ilk"
	pItem.UnitPrice = "$1.00"
	errs = pItem.validateProduceItem()

	if len(errs) != 1 {
		t.Errorf("1 errors expected but %d returned", len(errs))
	}

	//Invalid unit price format
	pItem = ProduceItem{}
	pItem.ProduceCode = "1111-1111-1111-1111"
	pItem.Name = "milk"
	pItem.UnitPrice = "1.00"
	errs = pItem.validateProduceItem()

	if len(errs) != 1 {
		t.Errorf("1 errors expected but %d returned", len(errs))
	}
}
