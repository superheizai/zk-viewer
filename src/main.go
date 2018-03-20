package main

import (
	"log"
	"net/http"
	"flag"
	. "common"
	. "db"
	. "router"
)

func main() {
	var env *string = flag.String("env", "env", "Use -env <prd,dev or test>")
	flag.Parse()
	log.Println("current env is", *env)
	router := NewRouter(AllRoutes())
	Init(*env)
	defer Close()
	log.Fatal(http.ListenAndServe(":8080", router))

}

func Init(env string) {
	InitDb(env)
	InitZkPools()
}

func Close() {
	Db.Close()
	ZkPoolsInstance.Close()
}
