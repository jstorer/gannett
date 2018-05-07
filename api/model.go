package api

import (
	"net/url"
	"strings"
	"sync"
)

type ProduceItem struct {
	ProduceCode string `json:"produce_code"`
	Name        string `json:"name"`
	UnitPrice   string `json:"unit_price"`
}

type DBObject struct {
	mu   sync.RWMutex
	Data []ProduceItem
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

func writeNewProduceItem(pItem ProduceItem, pItemChnl chan ProduceItem) {
	ProduceDB.mu.Lock()
	defer ProduceDB.mu.Unlock()

	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	for _, item := range ProduceDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			pItemChnl <- ProduceItem{}
			return
		}
	}
	ProduceDB.Data = append(ProduceDB.Data, pItem)
	pItemChnl <- pItem
}

func writeUpdateProduceItem(pCode string, pItem ProduceItem, pItemChnl chan ProduceItem) {
	ProduceDB.mu.Lock()
	defer ProduceDB.mu.Unlock()
	println(pCode)

	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	for _, item := range ProduceDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			pItemChnl <- ProduceItem{ProduceCode: "0"}
			return
		}
	}
	for index, item := range ProduceDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			ProduceDB.Data[index].ProduceCode = pItem.ProduceCode
			ProduceDB.Data[index].Name = pItem.Name
			ProduceDB.Data[index].UnitPrice = pItem.UnitPrice
			pItemChnl <- pItem
			return
		}
	}

	pItemChnl <- ProduceItem{}
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

func (pItem *ProduceItem) validateProduceItem() url.Values {
	errs := url.Values{}
	//all required fields exist
	if pItem.ProduceCode == "" {
		errs.Add("produce_code", "produce field is required")
	}

	if pItem.Name == "" {
		errs.Add("name", "name field is required")
	}

	if pItem.UnitPrice == "" {
		errs.Add("unit_price", "unit price field is required")
	}
	//are all values valid formats
	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	if !IsValidProduceCode(pItem.ProduceCode) {
		errs.Add("produce_code", "invalid produce code format")
	}

	if !IsValidName(pItem.Name) {
		errs.Add("name", "invalid name format")
	}

	if !IsValidUnitPrice(pItem.UnitPrice) {
		errs.Add("unit_price", "invalid unit price format")
	}
	return errs
}
