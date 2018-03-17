package db

import (
	"fmt"
	"testing"

	. "models"
)

func TestInsert(t *testing.T) {

	InitDb();

	fmt.Println("cache----------------")

	for key,value :=range Cache {
		fmt.Println(key)

		fmt.Println(value.Name)
		fmt.Println(value.Ips)
		fmt.Println(value.Id)

	}

	fmt.Println("cache end----------------")



	server1 := &Server{Ips: "10.9.20.106:2181", Name: "106", Version: "3.4.6"}
	server2 := &Server{Ips: "10.9.20.100:2181,10.9.20.101:2181,10.9.20.65:2181", Name: "dev", Version: "3.4.6"}

	id1 := Insert(server1);
	id2 := Insert(server2)
	fmt.Println(id1)
	fmt.Println(id2)
	es, err := Get("106")

	fmt.Println("cache----------------")

	for key,value :=range Cache {
		fmt.Println(key)

		 fmt.Println(value.Name)
		fmt.Println(value.Ips)
		fmt.Println(value.Id)

	}

	fmt.Println("cache end----------------")


	if (err != nil) {
		fmt.Println(es.Ips)
		fmt.Println(es.Name)
		fmt.Println(es.Version)
	} else {
		fmt.Println("zk----------------")

		fmt.Println(es.Ips)
		fmt.Println(es.Name)
		fmt.Println(es.Version)
	}

	fmt.Println("zks----------------")
	l, err := ListZks()

	if (err == nil) {
		for e := l.Front(); e != nil; e = e.Next() {
			t := e.Value.(*Server);
			fmt.Println(t.Ips)
			fmt.Println(t.Name)
			fmt.Println(t.Version)
		}

	}

	Delete(server1.Name)
	Delete(server2.Name)

	fmt.Println("deleted zks----------------")
	l, err = ListZks()

	if (err == nil) {
		for e := l.Front(); e != nil; e = e.Next() {
			t := e.Value.(*Server);
			fmt.Println(t.Ips)
			fmt.Println(t.Name)
			fmt.Println(t.Version)
		}

	}

	fmt.Println("cache----------------")

	for key,value :=range Cache {
		fmt.Println(key)

		fmt.Println(value.Name)
		fmt.Println(value.Ips)
		fmt.Println(value.Id)

	}

	fmt.Println("cache end----------------")


}
