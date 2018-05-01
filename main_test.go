package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type ProduceTestCode struct {
	ProduceCode string
	Valid       bool
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
		assert.Equal(t, item.Valid, isValidProduceCode(item.ProduceCode), fmt.Sprintf("On Produce Code: `%s`", item.ProduceCode))
	}
}

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/produce", getAllProduce).Methods("GET")
	return router
}

func TestGetAllProduce(t *testing.T) {
	request := httptest.NewRequest("GET", "/produce", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "Ok response is expected")
}
