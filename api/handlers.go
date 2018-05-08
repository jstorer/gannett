package api

import "github.com/gorilla/mux"

func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/produce", handleGetAllProduce).Methods("GET")
	router.HandleFunc("/api/produce/{produce_code}", handleGetProduceItem).Methods("GET")
	router.HandleFunc("/api/produce/{produce_code}", handleUpdateProduceItem).Methods("POST")
	router.HandleFunc("/api/produce", handleCreateProduceItem).Methods("POST")
	router.HandleFunc("/api/produce/{produce_code}", handleDeleteProduceItem).Methods("DELETE")
	return router
}
