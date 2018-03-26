package models

import (
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

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
	Cluster string `json:"cluster"`
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
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Content  string  `json:"content"`
	Stat     *ZkStat `json:"stat"`
	Children []*ZkTreeNodeInfo
}
type ChildrenSlice []*ZkTreeNodeInfo

func (c ChildrenSlice) Len() int {
	return len(c)
}
func (c ChildrenSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ChildrenSlice) Less(i, j int) bool {
	return (strings.Compare(c[i].Path, c[j].Path) == -1)
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

type ZkStat struct {
	Czxid          int64  // The zxid of the change that caused this znode to be created.
	Mzxid          int64  // The zxid of the change that last modified this znode.
	Ctime          string // The time in milliseconds from epoch when this znode was created.
	Mtime          string // The time in milliseconds from epoch when this znode was last modified.
	Version        int32  // The number of changes to the data of this znode.
	Cversion       int32  // The number of changes to the children of this znode.
	Aversion       int32  // The number of changes to the ACL of this znode.
	EphemeralOwner int64  // The session id of the owner of this znode if the znode is an ephemeral node. If it is not an ephemeral node, it will be zero.
	DataLength     int32  // The length of the data field of this znode.
	NumChildren    int32  // The number of children of this znode.
	Pzxid          int64  // last modified children
}
