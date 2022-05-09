package main

import (
	"fmt"
	"net/http"
	"resizeimage/src"
	"resizeimage/util"

	"github.com/gorilla/mux"
)

func main() {

	src.InitFromFile("config/config.toml")
	conf := src.GetConfig()

	src.InitLogger(conf.DBHost, conf.DBName, conf.DBCollection, conf.DBUsername, conf.DBPassword)

	util.InitAwsSession(conf.StorageRegion, conf.StorageAccessKeyID, conf.StorageAccessKeySecret)

	router := mux.NewRouter()
	src.InitRouter(router)

	endPoint := fmt.Sprintf("%s:%d", conf.ServerAddress, conf.ServerPort)
	fmt.Println("Starting the application...", endPoint)
	http.ListenAndServe(endPoint, router)
}
