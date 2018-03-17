package common

type ZkPools struct {
	pools map[string]*ZkPool
}

var ZkPoolsInstance *ZkPools

func InitZkPools() {

	ZkPoolsInstance = new(ZkPools)
	ZkPoolsInstance.pools = make(map[string]*ZkPool)
}

func (zKPools *ZkPools) Get(name string) *ZkPool {

	pool, ok := zKPools.pools[name]

	if (!ok) {
		pool = InitZkPool(5, 10, name)
		zKPools.pools[name] = pool

	}
	return pool
}

func (zKPools *ZkPools) Close() {
	for _, pool := range zKPools.pools {
		pool.Close()
	}
}
