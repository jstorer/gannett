package api

import "github.com/gorilla/mux"

func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/produce", getAllProduce).Methods("GET")
	router.HandleFunc("/api/produce/{produce_code}", getProduceItem).Methods("GET")
	router.HandleFunc("/api/produce/{produce_code}", updateProduceItem).Methods("PUT")
	router.HandleFunc("/api/produce", createProduceItem).Methods("POST")
	router.HandleFunc("/api/produce/{produce_code}", deleteProduceItem).Methods("DELETE")
	return router
}
