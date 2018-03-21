package main

import (
	"log"
	"net/http"
	"flag"
	. "common"
	. "db"
	. "router"
	"os/signal"
	"syscall"
	"fmt"
	"os"
)

func main() {
	var env *string = flag.String("env", "env", "Use -env <prd,dev or test>")
	flag.Parse()
	log.Println("current env is", *env)
	router := NewRouter(AllRoutes())
	Init(*env)
	defer Close()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		//Close()
		server.Close()
		done <- true
	}()

	//本来直接使用http去监听端口并启动，但是这样没有任何返回值，没有办法优雅关闭http服务,现在使用server
	//log.Fatal(http.ListenAndServe(":8080", router))
	go log.Fatal(server.ListenAndServe());
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}

func Init(env string) {
	InitDb(env)
	InitZkPools()
}

func Close() {
	fmt.Println("begin release")

	Db.Close()
	ZkPoolsInstance.Close()
	fmt.Println("end release")

}
