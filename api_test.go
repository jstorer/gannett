package main

import (
	"fmt"
	"github.com/jstorer/gannett/api"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ProduceTestCode struct {
	ProduceCode string
	Valid       bool
}

type UnitPriceTestValues struct {
	UnitPrice string
	Valid     bool
}

var (
	server     *httptest.Server
	reader     io.Reader
	produceUrl string
)

func init() {
	server = httptest.NewServer(api.Handlers())
	produceUrl = fmt.Sprintf("%s/produce", server.URL)
}


func TestIsValidProduceCode(t *testing.T) {
	var testCodes []ProduceTestCode

	//invalid produce code formats
	testCodes = append(testCodes, ProduceTestCode{"", false})
	testCodes = append(testCodes, ProduceTestCode{"a", false})
	testCodes = append(testCodes, ProduceTestCode{"ab", false})
	testCodes = append(testCodes, ProduceTestCode{"abc", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-12", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-123", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y8", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y8-q", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y8-q1", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y8-q13", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y8-q123a", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y80-q123", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-12340-z9y8-q123", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd0-1234-z9y8-q123", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y8-q123-", false})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y8-q123-abcd", false})

	//valid codes
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-z9y8-q123", true})
	testCodes = append(testCodes, ProduceTestCode{"abcd-1234-Z9Y8-q12W", true})
	testCodes = append(testCodes, ProduceTestCode{"1111-2222-3333-4444", true})
	testCodes = append(testCodes, ProduceTestCode{"aaaa-bbbb-cccc-dddd", true})
	testCodes = append(testCodes, ProduceTestCode{"AAAA-BBBB-CCCC-DDDD", true})

	for _, item := range testCodes {
		assert.Equal(t, item.Valid, api.IsValidProduceCode(item.ProduceCode), fmt.Sprintf("On Produce Code: `%s`", item.ProduceCode))
	}
}

func TestIsValidUnitPrice(t *testing.T) {
	var testPrices []UnitPriceTestValues

	//invalid prices
	testPrices = append(testPrices, UnitPriceTestValues{"", false})
	testPrices = append(testPrices, UnitPriceTestValues{"0", false})
	testPrices = append(testPrices, UnitPriceTestValues{"5", false})
	testPrices = append(testPrices, UnitPriceTestValues{"5.23", false})
	testPrices = append(testPrices, UnitPriceTestValues{"5.2", false})
	testPrices = append(testPrices, UnitPriceTestValues{"$", false})
	testPrices = append(testPrices, UnitPriceTestValues{"$01", false})
	testPrices = append(testPrices, UnitPriceTestValues{"$01.50", false})
	testPrices = append(testPrices, UnitPriceTestValues{"$21,12345", false})
	testPrices = append(testPrices, UnitPriceTestValues{"$5.123", false})

	//valid prices
	testPrices = append(testPrices, UnitPriceTestValues{"$0.1", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$0.01", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$1", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$1.0", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$1.00", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$0.10", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$4231", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$4,006", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$4,000.5", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$4,000.93", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$4,000,001.23", true})
	testPrices = append(testPrices, UnitPriceTestValues{"$4,000.00", true})

	for _, item := range testPrices {
		assert.Equal(t, item.Valid, api.IsValidUnitPrice(item.UnitPrice), fmt.Sprintf("On Unit Price: `%s`", item.UnitPrice))
	}

}

func TestGetAllProduce(t *testing.T) {
	api.Initialize()
	request, err := http.NewRequest("GET", produceUrl, nil)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("200 OK expected but %d returned", response.StatusCode)
	}
}

func TestGetProduceItem(t *testing.T) {
	api.Initialize()
	//VALID TEST----------------------------------------------------------
	//Get item
	api.ProduceDB.Data = append(api.ProduceDB.Data, api.ProduceItem{ProduceCode: "ABCD-1234-EFGH-5678", Name: "Black Beans", UnitPrice: "$2.25"})
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

func TestCreateProduceItem(t *testing.T) {
	api.Initialize()
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
	api.Initialize()
	//VALID TEST----------------------------------------------------------
	//Delete item
	api.ProduceDB.Data = append(api.ProduceDB.Data, api.ProduceItem{ProduceCode: "ABCD-1234-EFGH-5678", Name: "Black Beans", UnitPrice: "$2.25"})
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
