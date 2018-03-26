package common

import (
	//"fmt"
	"log"
	"strings"

	//"strconv"
	"sync"
	"time"

	. "base"
	"db"

	"github.com/samuel/go-zookeeper/zk"
)

type ZkPool struct {
	Coons   chan *ZkConn // 连接池
	Core    int          //核心连接数
	Max     int          //最大连接数
	Mux     sync.Mutex   //锁
	Closed  bool
	Name    string
	Url     string
	current int // 当前分配出的连接数量，包括conns里面的和在使用未归还的

	busy          bool       // current status, 记录当前操作频繁与否;true时候，put操作直接入队；false时候，先看是否小于core，再入队，大于core，直接回收了
	threadhold    int        // status切换阈值，多少个每秒会做切换；
	durationtimes int        // 持续几次低过threadhold，会切换busy为false
	times         chan int64 // times channel,每个get请求会把自己的请求时间放到channel里面，统计每个时间点的请求信息
	marktime      int64      //此时在统计的标记时间点,秒级时间，time.unix()返回
	markCount     int        // markTime时间点上，被mark的次数
	markTimes     int        // 现在已经因为超过或者低于阈值，而累积的次数。当超过durationtimes，会触发busy值的改变
}

type ZkConn struct {
	*zk.Conn
}

func InitZkPool(core int, max int, name string) *ZkPool {

	zkPool := new(ZkPool)
	zkPool.Core = core
	zkPool.Max = max
	zkPool.Coons = make(chan *ZkConn, max)
	zkPool.Closed = false
	zkPool.Name = name
	svr, _ := db.Cache[name]
	zkPool.Url = svr.Ips

	zkPool.busy = true
	zkPool.threadhold = 10;
	zkPool.durationtimes = 2;
	zkPool.times = make(chan int64, 1000)

	zkPool.marktime = time.Now().Unix()
	zkPool.markCount = 0
	zkPool.markTimes = 0

	go zkPool.switchStatus()
	return zkPool

}

func (zkPool *ZkPool) switchStatus() {

	for ; ; {
		if (!zkPool.Closed) {
			cur := <-zkPool.times

			if zkPool.marktime == cur {
				zkPool.markCount++
				//在空闲时候，如果操作次数大于2倍阈值的时候，会直接转为忙状态，不会等待这一秒的统计周期结束
				if (!zkPool.busy && zkPool.markCount > 2*zkPool.threadhold) {
					zkPool.busy = true
				}
			} else {
				zkPool.marktime = cur
				lastMarkCount := zkPool.markCount
				zkPool.markCount = 1
				//新的时间窗到来，根据之前是忙还是空，来判断状态
				if (zkPool.busy) {
					if (lastMarkCount < zkPool.threadhold) {
						zkPool.markTimes++
						if (zkPool.markTimes >= zkPool.durationtimes) {
							zkPool.busy = false
							zkPool.markTimes = 0;
						}
					} else {
						//新的时间窗到来，但是其阈值没有超过设定阈值
						zkPool.markTimes = 0
					}
				} else {

					if (lastMarkCount >= zkPool.threadhold) {
						zkPool.markTimes++
						if (zkPool.markTimes >= zkPool.durationtimes) {
							zkPool.busy = true
							zkPool.markTimes = 0;
						}
					} else {
						//新的时间窗到来，但是其阈值没有超过设定阈值
						zkPool.markTimes = 0
					}
				}

			}

		}
	}

}

func (zkPool *ZkPool) Put(conn *ZkConn) {
	if zkPool.Closed {
		return
	}
	zkPool.Mux.Lock()
	defer zkPool.Mux.Unlock()

	if (zkPool.busy) {
		zkPool.Coons <- conn
	} else {
		if len(zkPool.Coons) < zkPool.Core {
			zkPool.Coons <- conn
		} else {
			zkPool.current--
			conn.Close()
		}
	}

}

func (zkPool *ZkPool) Get() (*ZkConn, error) {
	if zkPool.Closed {
		log.Fatal("zkPool has been Closed")
	}
	go zkPool.pinPoint(time.Now().Unix())
	select {
	case conn := <-zkPool.Coons:
		log.Println("exist conn is status " + conn.State().String())
		return conn, nil
	default:
		zkPool.Mux.Lock()
		defer zkPool.Mux.Unlock()
		if (zkPool.current < zkPool.Max) {
			conn, _, err := zk.Connect(strings.Split(zkPool.Url, ","), 5*time.Second)
			if (err != nil) {
				log.Fatal(err)
				return nil, err
			}
			log.Println("new conn is status " + conn.State().String())

			zkConn := ZkConn{conn}
			zkPool.current++
			return &zkConn, nil
		} else {
			for {
				select {
				case <-time.After(5 * time.Second):
					log.Panic("can't get connection %s after 5S", zkPool.Url)
					return nil, TimeOutError
				case conn := <-zkPool.Coons:
					log.Println("wait conn is status " + conn.State().String())

					return conn, nil
				}
			}
		}
	}

}

func (zkPool *ZkPool) pinPoint(cur int64) {
	zkPool.times <- cur
}
func (zkPool *ZkPool) Close() {
	zkPool.Mux.Lock()
	defer zkPool.Mux.Unlock()

	for conn := range zkPool.Coons {
		conn.Close()
	}
	close(zkPool.times)
	close(zkPool.Coons)
}
