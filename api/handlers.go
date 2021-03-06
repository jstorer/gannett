//sets router and creates end points
package api

import "github.com/gorilla/mux"

//creates new router and sets end point function triggers
func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/produce", handleGetAllProduce).Methods("GET")
	router.HandleFunc("/api/produce/{produce_code}", handleGetProduceItem).Methods("GET")
	router.HandleFunc("/api/produce/{produce_code}", handleUpdateProduceItem).Methods("POST")
	router.HandleFunc("/api/produce", handleCreateProduceItem).Methods("POST")
	router.HandleFunc("/api/produce/{produce_code}", handleDeleteProduceItem).Methods("DELETE")
	return router
}
