package api

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server     *httptest.Server
	reader     io.Reader
	produceUrl string
)

func init() {
	server = httptest.NewServer(Handlers())
	produceUrl = fmt.Sprintf("%s/api/produce", server.URL)
	Initialize(true)
}

func reinitTest() {
	currentDB.Data = []ProduceItem{
		{"A12T-4GH7-QPL9-3N4M", "Lettuce", "$3.46"},
		{"E5T6-9UI3-TH15-QR88", "Peach", "$2.99"},
		{"YRT6-72AS-K736-L4AR", "Green Pepper", "$0.79"},
		{"2222-2222-2222-2222", "Gala Apple", "$3.59"},
	}
}

func TestIsValidProduceCode(t *testing.T) {
	var testCodes = []struct {
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
	reinitTest()
	var getAllTests = []struct {
		desc         string
		method       string
		path         string
		statusCode   int
		expectedBody string
	}{
		{"all items payload from reinitTest()", "GET", produceUrl,
			200, `[{"produce_code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","unit_price":"$3.46"},{"produce_code":"E5T6-9UI3-TH15-QR88","name":"Peach","unit_price":"$2.99"},{"produce_code":"YRT6-72AS-K736-L4AR","name":"Green Pepper","unit_price":"$0.79"},{"produce_code":"2222-2222-2222-2222","name":"Gala Apple","unit_price":"$3.59"}]`},
	}

	for _, item := range getAllTests {
		request, err := http.NewRequest(item.method, item.path, nil)
		response, err := http.DefaultClient.Do(request)

		if err != nil {
			t.Error(err)
		}

		responseData, _ := ioutil.ReadAll(response.Body)
		assert.Equal(t, item.expectedBody, string(responseData), fmt.Sprintf("unexpected response for %s", item.desc))
		assert.Equal(t, item.statusCode, response.StatusCode, "unexpected status code")
	}

}

func TestGetProduceItem(t *testing.T) {
	reinitTest()
	var getItemTests = []struct {
		desc         string
		method       string
		path         string
		statusCode   int
		expectedBody string
	}{
		{"get existing item", "GET", fmt.Sprintf("%s/A12T-4GH7-QPL9-3N4M", produceUrl),
			200, `{"produce_code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","unit_price":"$3.46"}`},
		//
		{"invalid produce code", "GET", fmt.Sprintf("%s/ABCDe-1234-EFGH-5678", produceUrl),
			400, "error 400 - invalid produce code format\n"},
		//
		{"produce code does note exist", "GET", fmt.Sprintf("%s/ABCD-1234-EFGH-0000", produceUrl),
			404, "error 404 - produce code does not exist\n"},
	}

	for _, item := range getItemTests {
		request, err := http.NewRequest(item.method, item.path, nil)
		response, err := http.DefaultClient.Do(request)

		if err != nil {
			t.Error(err)
		}

		responseData, _ := ioutil.ReadAll(response.Body)
		assert.Equal(t, item.expectedBody, string(responseData), fmt.Sprintf("unexpected response for %s", item.desc))
		assert.Equal(t, item.statusCode, response.StatusCode, "unexpected status code")
	}

}

func TestUpdateProduceItem(t *testing.T) {
	reinitTest()
	var updateItemTests = []struct {
		desc         string
		method       string
		path         string
		statusCode   int
		produceCode  string
		name         string
		unitPrice    string
		expectedBody string
	}{
		{"produce code remains same", "POST", fmt.Sprintf("%s/A12T-4GH7-QPL9-3N4M", produceUrl),
			200, "A12T-4GH7-QPL9-3N4M", "Lettuce", "$3.46", `{"produce_code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","unit_price":"$3.46"}`},
		//
		{"push to db check", "POST", fmt.Sprintf("%s/A12T-4GH7-QPL9-3N4M", produceUrl),
			200, "A12T-4GH7-QPL9-1111", "Cheese", "$5.00", `{"produce_code":"A12T-4GH7-QPL9-1111","name":"Cheese","unit_price":"$5.00"}`},
		//
		{"updated code already exists", "POST", fmt.Sprintf("%s/E5T6-9UI3-TH15-QR88", produceUrl),
			409, "2222-2222-2222-2222", "Cheese", "$5.00", "error 409 - updated produce code value already exists\n"},
		//
		{"produce code doesn't exist to update", "PUT", fmt.Sprintf("%s/E5T6-9UI3-TH15-1111", produceUrl),
			404, "A12T-4GH7-QPL9-3N4A", "Cheese", "$5.00", "error 404 - produce code does not exist\n"},
		//
		{"invalid end point", "POST", fmt.Sprintf("%s/E5T6-9UI3-TH15-111", produceUrl),
			400, "", "", "", "error 400 - invalid produce code format\n"},
		//
		{"bad payload", "POST", fmt.Sprintf("%s/E5T6-9UI3-TH15-QR88", produceUrl),
			400, "A12T-4GH7-QPL9-3NM", "", "5.00",
			`{"validationError":{"name":["name field is required","invalid name format"],"produce_code":["invalid produce code format"],"unit_price":["invalid unit price format"]}}`},
	}

	for _, item := range updateItemTests {
		produceItemJson := `{"produce_code":"` + item.produceCode + `","name":"` + item.name + `","unit_price":"` + item.unitPrice + `"}`
		reader = strings.NewReader(produceItemJson)
		request, err := http.NewRequest(item.method, item.path, reader)
		response, err := http.DefaultClient.Do(request)

		if err != nil {
			t.Error(err)
		}

		responseData, _ := ioutil.ReadAll(response.Body)
		assert.Equal(t, item.expectedBody, string(responseData), fmt.Sprintf("unexpected response for %s", item.desc))
		assert.Equal(t, item.statusCode, response.StatusCode, "unexpected status code")
	}
}

func TestCreateProduceItem(t *testing.T) {
	reinitTest()
	var createItemTests = []struct {
		desc         string
		method       string
		path         string
		statusCode   int
		produceCode  string
		name         string
		unitPrice    string
		expectedBody string
	}{
		{"create item", "POST", produceUrl, 201, "1234-5678-90ab-cdef",
			"Cheese", "$9.99", `{"produce_code":"1234-5678-90AB-CDEF","name":"Cheese","unit_price":"$9.99"}`},
		//
		{"try to create duplicate code", "POST", produceUrl, 409, "1234-5678-90ab-cdef",
			"Cheese", "$9.99", "error 409 - produce code already exists\n"},
		//
		{"left produce code field empty", "POST", produceUrl, 400, "",
			"Cheese", "$4.60", `{"validationError":{"produce_code":["produce field is required","invalid produce code format"]}}`},
		//
		{"left name and unit field empty", "POST", produceUrl, 400, "1111-1111-1111-1111",
			"", "", `{"validationError":{"name":["name field is required","invalid name format"],"unit_price":["unit price field is required","invalid unit price format"]}}`},
		//
		{"all fields invalid", "POST", produceUrl, 400, "23aja-fafe-grge-sdf",
			"Ch!eese", "23.432", `{"validationError":{"name":["invalid name format"],"produce_code":["invalid produce code format"],"unit_price":["invalid unit price format"]}}`},
	}

	for _, item := range createItemTests {
		produceItemJson := `{"produce_code":"` + item.produceCode + `","name":"` + item.name + `","unit_price":"` + item.unitPrice + `"}`
		reader = strings.NewReader(produceItemJson)
		request, err := http.NewRequest(item.method, item.path, reader)
		response, err := http.DefaultClient.Do(request)

		if err != nil {
			t.Error(err)
		}

		responseData, _ := ioutil.ReadAll(response.Body)
		assert.Equal(t, item.expectedBody, string(responseData), fmt.Sprintf("unexpected response for %s", item.desc))
		assert.Equal(t, item.statusCode, response.StatusCode, "unexpected status code")
	}
}

func TestDeleteProduceItem(t *testing.T) {
	reinitTest()
	var deleteItemTests = []struct {
		desc         string
		method       string
		path         string
		statusCode   int
		expectedBody string
	}{
		{"delete item", "DELETE", fmt.Sprintf("%s/A12T-4GH7-QPL9-3N4M", produceUrl),
			200, `{"produce_code":"A12T-4GH7-QPL9-3N4M","name":"","unit_price":""}`},
		//
		{"invalid produce code", "DELETE", fmt.Sprintf("%s/ABCDe-1234-EFGH-5678", produceUrl),
			400, "error 400 - invalid produce code format\n"},
		//
		{"code does not exist", "DELETE", fmt.Sprintf("%s/A12T-4GH7-QPL9-3N4M", produceUrl),
			404, "error 404 - produce code not found.\n"},
	}

	for _, item := range deleteItemTests {
		request, err := http.NewRequest(item.method, item.path, nil)
		response, err := http.DefaultClient.Do(request)

		if err != nil {
			t.Error(err)
		}

		responseData, _ := ioutil.ReadAll(response.Body)
		assert.Equal(t, item.expectedBody, string(responseData), fmt.Sprintf("unexpected response for %s", item.desc))
		assert.Equal(t, item.statusCode, response.StatusCode, "unexpected status code")
	}
}

func TestValidateProduceItem(t *testing.T) {
	var validateProduceItemTests = []struct {
		desc           string
		produceCode    string
		name           string
		unitPrice      string
		expectedOutput string
	}{
		{"everything valid", "1111-1111-1111-1111", "milk", "$1.00", `{"validationError":{}}`},
		{"everything empty", "", "", "",
			`{"validationError":{"name":["name field is required","invalid name format"],"produce_code":["produce field is required","invalid produce code format"],"unit_price":["unit price field is required","invalid unit price format"]}}`},
		//
		{"everything invalid", "12fava-sdfw-eaav-va", "fj#@j", " 12.2",
			`{"validationError":{"name":["invalid name format"],"produce_code":["invalid produce code format"],"unit_price":["invalid unit price format"]}}`},
		//
	}

	for _, item := range validateProduceItemTests {
		var pItem ProduceItem
		pItem.ProduceCode = item.produceCode
		pItem.Name = item.name
		pItem.UnitPrice = item.unitPrice
		validErrs := pItem.validateProduceItem()
		err := map[string]interface{}{"validationError": validErrs}
		response, _ := json.Marshal(err)
		assert.Equal(t, item.expectedOutput, string(response), "unexpected output")
	}

}
