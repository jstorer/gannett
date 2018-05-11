//Contains functions that manipulate the database or methods attached to created data types.
package api

import (
	"net/url"
	"strings"
	"sync"
)

//type to store a produce item
type ProduceItem struct {
	ProduceCode string `json:"produce_code"`
	Name        string `json:"name"`
	UnitPrice   string `json:"unit_price"`
}

//type to represent a database with a mutex to assist in preventing race conditions
type DBObject struct {
	mu   sync.RWMutex
	Data []ProduceItem
}

//return all items from the database on a channel, used RLock since only reading done.
func getAllProduceItems(allItemsChnl chan []ProduceItem) {
	currentDB.mu.RLock()
	defer currentDB.mu.RUnlock()
	allItemsChnl <- currentDB.Data
}

//returns a single produce item on a channel based on the given produce code
//if the item is not found an empty item is returned on the channel.
//RLock is used since only read operations done here
func getProduceItem(pCode string, pItemChnl chan ProduceItem) {
	currentDB.mu.RLock()
	defer currentDB.mu.RUnlock()
	for _, item := range currentDB.Data {
		if item.ProduceCode == pCode {
			pItemChnl <- item
			return
		}
	}
	pItemChnl <- ProduceItem{}
}

//creates a new produce item in the database by checking if the given produce code already exists.
//if the code exists an empty item is returned to the channel. If the code does not exist the item is
//appended to the database and returned on the channel
func createProduceItem(pItem ProduceItem, pItemChnl chan ProduceItem) {
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

//updates an item in the database of the given produce code. If the produce code given does not exist an empty item
//is returned on the channel. If the code exists but the new code being updated already exists in the database a produce code
// "0" is returned. If the item is able to be created the new contents overwrite the old ones at the given index and return
//the new produce item on the channel.
func updateProduceItem(pCode string, pItem ProduceItem, pItemChnl chan ProduceItem) {
	currentDB.mu.Lock()
	defer currentDB.mu.Unlock()

	pItem.ProduceCode = strings.ToUpper(pItem.ProduceCode)
	pCode = strings.ToUpper(pCode)

	for index, item := range currentDB.Data {
		if item.ProduceCode == pCode {
			for _, item := range currentDB.Data {
				if item.ProduceCode == pItem.ProduceCode && pCode != pItem.ProduceCode {
					pItemChnl <- ProduceItem{ProduceCode: "0", Name: "", UnitPrice: ""}
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

//deletes an item from the database based on the incoming produce code. If the produce code is not found an
//empty item is returned. If the code is found it is then removed from the database.
func deleteProduceItem(pCode string, pItemChnl chan ProduceItem) {
	currentDB.mu.Lock()
	defer currentDB.mu.Unlock()
	var pItem ProduceItem

	for index, item := range currentDB.Data {
		if item.ProduceCode == pCode {
			currentDB.Data = append(currentDB.Data[:index], currentDB.Data[index+1:]...)
			pItem.ProduceCode = item.ProduceCode
			pItem.Name = item.Name
			pItem.UnitPrice = item.UnitPrice
			pItemChnl <- pItem
			return
		}
	}
	pItemChnl <- ProduceItem{}
}

//checks that produce item fields are populated as intended and in the correct format.
func (pItem *ProduceItem) validateProduceItem() url.Values {
	errs := url.Values{} //store errors

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
	if !isValidProduceCode(pItem.ProduceCode) {
		errs.Add("produce_code", "invalid produce code format")
	}

	if !isValidName(pItem.Name) {
		errs.Add("name", "invalid name format")
	}

	if !isValidUnitPrice(pItem.UnitPrice) {
		errs.Add("unit_price", "invalid unit price format")
	}

	return errs
}
