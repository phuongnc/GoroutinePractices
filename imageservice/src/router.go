package src

import "github.com/gorilla/mux"

func InitRouter(router *mux.Router) {
	router.HandleFunc("/resize-image", ResizeImage).Methods("POST")
}
