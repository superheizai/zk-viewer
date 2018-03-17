package common

import (
	"fmt"
	"testing"
	"time"

	"db"
)

//本地启动
func TestZkPool_Put(t *testing.T) {
	db.InitDb()
	zkPool := InitZkPool(10, 20, "local");
	for i := 0; i < 320; i++ {

		con, _ := zkPool.Get()
		if (i > 100) {
			time.Sleep(10 * time.Millisecond)
		} else {
			time.Sleep(200 * time.Millisecond)

		}
		//time.Sleep(100 * time.Millisecond)

		zkPool.Put(con)
		fmt.Println(zkPool.busy)
		fmt.Println(zkPool.marktime)
		fmt.Println(zkPool.markCount)
		fmt.Println(zkPool.markTimes)
	}
	fmt.Println("==============================")
}
