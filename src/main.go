package main

import (
	"log"
	"net/http"

	. "common"
	. "db"
	. "router"
)

func main() {

	router := NewRouter(AllRoutes())
	Init()
	defer Close()
	log.Fatal(http.ListenAndServe(":8080", router))

}

func Init() {
	InitDb()
	InitZkPools()
}

func Close() {
	Db.Close()
	ZkPoolsInstance.Close()
}
