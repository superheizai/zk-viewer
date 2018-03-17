package models

import "github.com/samuel/go-zookeeper/zk"

type Zk struct {
	// The main identifier for the Book. This will be unique.
	Clusters  []string `json:"clusters"`
	Topics    []string `json:"topics"`
	Consumers []string `json:"consumers"`
}

type ZkNode struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ZkNodeInfo struct {
	ZkNode
	Cluster string
}

//type ZkTreeNode struct {
//	*zk.Stat
//	Name    string `json:"name"`
//	Path    string `json:"path"`
//	Content string `json:"content"`
//}

type Content struct {
	*zk.Stat
	Content string
}

type ZkTreeNodeInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Content  string `json:"content"`
	Children []*ZkTreeNodeInfo
}

type Server struct {
	Id      int64  `json:"id"`
	Ips     string `json:"ips"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type DeleteInfo struct {
	Id int64 `json:"id"`
}
