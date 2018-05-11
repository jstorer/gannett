//tests for model.go
package api

import (
	"github.com/stretchr/testify/assert"
	"fmt"
	"testing"
	"encoding/json"
)

//test get all produce items
func TestGetAllProduceItems(t *testing.T) {
	reinitTest()
	pItemChnl := make(chan []ProduceItem)
	go getAllProduceItems(pItemChnl)
	allItems := <-pItemChnl
	assert.Equal(t, currentDB.Data, allItems, "DB not returning correct values")
}

//test getting a single produce item from server
func TestGetProduceItem(t *testing.T) {
	var getProduceItemTests = []struct {
		desc           string
		produceCode    string
		expectedOutput string
	}{
		{"produce code valid", "2222-2222-2222-2222", "2222-2222-2222-2222"},
		{"produce code invalid", "aji-ewfi-23ijf", ""},
		{"produce code does not exist", "1111-1111-1111-1111", ""},
	}
	for _, item := range getProduceItemTests {
		pItemChnl := make(chan ProduceItem)
		go getProduceItem(item.produceCode, pItemChnl)
		pItem := <-pItemChnl
		assert.Equal(t, item.expectedOutput, pItem.ProduceCode, fmt.Sprintf("unexpected output for %s", item.desc))
	}
}

//test creating a produce item on the create end point
func TestCreateProduceItem(t *testing.T) {
	var createProduceItemTests = []struct {
		desc           string
		pItem          ProduceItem
		expectedOutput ProduceItem
	}{
		{"valid produce item", ProduceItem{"1111-1111-1111-1111", "Bacon", "$1.23"}, ProduceItem{"1111-1111-1111-1111", "Bacon", "$1.23"}},
		{"produce code already exists", ProduceItem{"2222-2222-2222-2222", "Bacon", "$1.23"}, ProduceItem{"", "", ""}},
	}
	for _, item := range createProduceItemTests {
		reinitTest()
		pItemChanl := make(chan ProduceItem)
		go createProduceItem(item.pItem, pItemChanl)
		pItem := <-pItemChanl
		assert.Equal(t, item.expectedOutput, pItem, fmt.Sprintf("unexpeted output for %s", item.desc))
	}
}

//testing updating a produce item on the update end point
func TestUpdateProduceItem(t *testing.T) {
	var updateProduceItemTests = []struct {
		desc           string
		produceCode    string
		pItem          ProduceItem
		expectedOutput ProduceItem
	}{
		{"valid produce code", "2222-2222-2222-2222", ProduceItem{"1111-1111-1111-1111", "Bacon", "$1.23"}, ProduceItem{"1111-1111-1111-1111", "Bacon", "$1.23"}},
		{"updated code exists", "2222-2222-2222-2222", ProduceItem{"A12T-4GH7-QPL9-3N4M", "Bacon", "$1.23"}, ProduceItem{"0", "", ""}},
		{"produce code not found", "ABCD-2222-2222-2222", ProduceItem{"A12T-4GH7-QPL9-3N4M", "Bacon", "$1.23"}, ProduceItem{"", "", ""}},
	}

	for _, item := range updateProduceItemTests {
		reinitTest()
		pItemChnl := make(chan ProduceItem)
		go updateProduceItem(item.produceCode, item.pItem, pItemChnl)
		pItem := <-pItemChnl
		assert.Equal(t, item.expectedOutput, pItem, fmt.Sprintf("unexpeted output for %s", item.desc))
	}
}

//test deleting an item from the database on the delete end point
func TestDeleteProduceItem(t *testing.T) {
	var deleteProduceItemTests = []struct {
		desc           string
		produceCode    string
		expectedOutput ProduceItem
	}{
		{"valid produce code", "2222-2222-2222-2222", ProduceItem{"2222-2222-2222-2222", "Gala Apple", "$3.59"}},
		{"code does not exist", "ABCD-2222-2222-2222", ProduceItem{"", "", ""}},
	}
	for _, item := range deleteProduceItemTests {
		reinitTest()
		pItemChnl := make(chan ProduceItem)
		go deleteProduceItem(item.produceCode, pItemChnl)
		pItem := <-pItemChnl
		assert.Equal(t, item.expectedOutput, pItem, fmt.Sprintf("unexpeted output for %s", item.desc))
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
		assert.Equal(t, item.expectedOutput, string(response), fmt.Sprintf("unexpected output for %s", item.desc))
	}

}