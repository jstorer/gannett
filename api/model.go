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
	currentDB.mu.RLock()
	defer currentDB.mu.RUnlock()
	sendchannel <- currentDB.Data
}

func readProduceItem(pCode string, sendchannel chan<- ProduceItem) {
	currentDB.mu.RLock()
	defer currentDB.mu.RUnlock()
	for _, item := range currentDB.Data {
		if item.ProduceCode == pCode {
			sendchannel <- item
			return
		}
	}
	sendchannel <- ProduceItem{}
}

func writeNewProduceItem(pItem ProduceItem, pItemChnl chan ProduceItem) {
	currentDB.mu.Lock()
	defer currentDB.mu.Unlock()

	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	for _, item := range currentDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			pItemChnl <- ProduceItem{}
			return
		}
	}
	currentDB.Data = append(currentDB.Data, pItem)
	pItemChnl <- pItem
}

func writeUpdateProduceItem(pCode string, pItem ProduceItem, pItemChnl chan ProduceItem) {
	currentDB.mu.Lock()
	defer currentDB.mu.Unlock()

	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	pCode = strings.ToUpper(pCode)

	for index, item := range currentDB.Data {
		if item.ProduceCode == pCode {
			for _, item := range currentDB.Data {
				if item.ProduceCode == pItem.ProduceCode && pCode != pItem.ProduceCode {
					pItemChnl <- ProduceItem{ProduceCode: "0"}
					return
				}
			}
			currentDB.Data[index].ProduceCode = pItem.ProduceCode
			currentDB.Data[index].Name = pItem.Name
			currentDB.Data[index].UnitPrice = pItem.UnitPrice
			pItemChnl <- pItem
			return
		}
	}

	pItemChnl <- ProduceItem{}
}

func removeProduceItem(pItem ProduceItem, pItemChnl chan ProduceItem) {
	currentDB.mu.Lock()
	defer currentDB.mu.Unlock()

	for index, item := range currentDB.Data {
		if item.ProduceCode == pItem.ProduceCode {
			currentDB.Data = append(currentDB.Data[:index], currentDB.Data[index+1:]...)
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
