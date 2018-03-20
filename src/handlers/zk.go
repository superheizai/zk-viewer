package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	. "base"
	. "common"
	. "db"
	. "models"

	"github.com/julienschmidt/httprouter"
	"github.com/samuel/go-zookeeper/zk"
)

//list zk nodes
func Zks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	//zks, err := ListZks();
	zks, err := ArrayZks();
	if ( err != nil ) {
		writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
	} else {
		writeOKResponse(w, zks)
	}
}

//list zk nodes
func ZksNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	q := r.URL.Query();
	name := q.Get("name");
	ip := q.Get("ip");
	zks, err := FindZks(name, ip);
	if ( err != nil ) {
		writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
	} else {
		writeOKResponse(w, zks)
	}
}

func FindByName(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	ip := params.ByName("ip")

	zk, err := Get(name, ip)
	if ( err != nil ) {
		writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
	} else {
		writeOKResponse(w, zk)
	}
}

//list zk nodes
func CreateZk(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	server := &Server{}
	if err := populateModelFromHandler(w, r, params, server); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	writeOKResponse(w, Insert(server))
}

func UpdateZk(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	server := &Server{}
	if err := populateModelFromHandler(w, r, params, server); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	writeOKResponse(w, Update(server))
}

func OptionsRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length, Authorization, Accept,X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.WriteHeader(http.StatusOK)
}

func DeleteZk(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	server := &Server{}

	if err := populateModelFromHandler(w, r, params, server); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	log.Print("deleted id is ", server.Name)
	rows, err := Delete(server.Name, server.Id);
	if err != nil {
		log.Print(err)
		return
	}
	writeOKResponse(w, rows)
}

//list zk nodes
func ListPath(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	cluster := params.ByName("cluster")

	rawpath := r.URL.Query().Get("path");
	if (rawpath == "") {
		rawpath = "/";
	}
	path, _ := url.PathUnescape(rawpath);
	log.Print("cluster is ", cluster)
	server, exist := Cache[cluster];
	if ( !exist) {
		writeErrorResponse(w, ResponseErrorCode, NoClusterFoundError.Error())
		return
	}

	pool := ZkPoolsInstance.Get(server.Name)
	conn, _ := pool.Get()
	defer pool.Put(conn)
	fmt.Println(conn.State().String())

	exists, _, err := conn.Exists(path);
	if (err != nil) {
		log.Panic("error when exist", err)
		writeErrorResponse(w, 0, "");
		return;
	}
	if (exists) {
		zktreeNodeInfo := &ZkTreeNodeInfo{}

		paths, _, _ := conn.Children(path)
		dataB, stat, err := conn.Get(path);
		if (err != nil) {
			writeErrorResponse(w, 0, err.Error());
			return;
		}

		//zktreeNode := &ZkTreeNode{}

		zktreeNodeInfo.Name = filepath.Base(path)
		zktreeNodeInfo.Path = path

		content := Content{}
		content.Content = string(dataB[:])
		content.Stat = stat
		//cnt, err := json.Marshal(content);
		//if err != nil {
		//	log.Panic("marshal content errror", err)
		//	return
		//}
		zktreeNodeInfo.Content = string(dataB[:]);
		zktreeNodeInfo.Stat = stat;
		//zktreeNodeInfo.Content = string(cnt[:]);
		subNodes := []*ZkTreeNodeInfo{}
		for _, v := range paths {
			subNodeInfo := &ZkTreeNodeInfo{}

			subNodeInfo.Name = v
			if (path != "/") {
				subNodeInfo.Path = path + "/" + v
			} else {
				subNodeInfo.Path = "/" + v
			}
			subNodes = append(subNodes, subNodeInfo)
		}
		zktreeNodeInfo.Children = subNodes
		writeOKResponse(w, zktreeNodeInfo)

	} else {
		writeErrorResponse(w, 0, "path not exist");
	}

}

//list zk nodes
func ListZms(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	cluster := params.ByName("cluster")

	log.Print("cluster is ", cluster)
	server, exist := Cache[cluster];
	if ( !exist) {
		writeErrorResponse(w, ResponseErrorCode, NoClusterFoundError.Error())
		return
	}

	pool := ZkPoolsInstance.Get(server.Name)
	conn, _ := pool.Get()
	defer pool.Put(conn)
	//if err != nil {
	//	writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
	//	return
	//}
	readAndBuildResult(conn, w)
}

func readAndBuildResult(zkConn *ZkConn, w http.ResponseWriter) {
	clusters, _, _ := zkConn.Children("/zms/cluster")
	consumers, _, _ := zkConn.Children("/zms/consumergroup")
	topics, _, _ := zkConn.Children("/zms/topic")

	zk := &Zk{}
	zk.Clusters = clusters
	zk.Topics = topics
	zk.Consumers = consumers
	writeOKResponse(w, zk)
}

//list zk nodes
func ReadZkNode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	cluster := params.ByName("cluster")
	path := params.ByName("path")

	path, err := url.QueryUnescape(path)
	if (err != nil) {
		writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	log.Print("cluster is ", cluster)

	server, exist := Cache[cluster];
	if ( !exist) {
		writeErrorResponse(w, ResponseErrorCode, NoClusterFoundError.Error())
		return
	}
	pool := ZkPoolsInstance.Get(server.Name)
	conn, _ := pool.Get()

	exist, _, err = conn.Exists(path)
	if err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if (exist) {
		contentByte, _, _ := conn.Get(path);
		str := string(contentByte[:])
		writeOKResponse(w, str)
	} else {
		writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

}

func CreateZkNode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	zkNode := &ZkNodeInfo{}
	if err := populateModelFromHandler(w, r, params, zkNode); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}

	log.Print("cluster is ", zkNode.Path)

	b1 := []byte(zkNode.Content)
	//flags定义
	//0:永久，除非手动删除
	//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
	//zk.FlagSequence  = 2:会自动在节点后面添加序号
	//3:Ephemeral和Sequence，即，短暂且自动添加序号

	//conn.Create(zkNode.Path, b1, 0, zk.WorldACL(zk.PermAll))
	//str,err :=conn.Create(zkNode.Path, b1, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	str, err := CreateDepthNode(zkNode.Path, b1, 0, zkNode.Cluster)
	if err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	log.Print(str)
}

func CreateDepthNode(path string, data []byte, flag int32, cluster string) (string, error) {
	nodeArr := strings.Split(path, "/")
	addPath := ""

	server, exist := Cache[cluster];
	if ( !exist) {
		return "", errors.New("no cluster")
	}
	pool := ZkPoolsInstance.Get(server.Name)
	mConn, _ := pool.Get()

	defer mConn.Close()

	//mConn, _, err := zk.Connect([]string{"10.9.20.106:2181"}, time.Second)
	//if err != nil {
	//	return "", err
	//}
	existPath, _, _ := mConn.Exists(path)

	if existPath {
		mConn.Set(path, data, -1);
		return "", nil
	}
	for i := 0; i < len(nodeArr)-1; i++ {
		if len(nodeArr[i]) == 0 {
			continue
		}
		addPath += ("/" + nodeArr[i])
		exist, _, err := mConn.Exists(addPath)
		if err != nil {
			return "", err
		}
		if exist {
			continue
		}
		if _, cerr := mConn.Create(addPath, []byte(""), flag, zk.WorldACL(zk.PermAll)); cerr != nil {
			return "", cerr
		}
	}

	con, err := mConn.Create(path, data, flag, zk.WorldACL(zk.PermAll))
	return con, err
}

func DeleteZkNode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	idInfo := &ZkNodeInfo{}
	if err := populateModelFromHandler(w, r, params, idInfo); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}

	log.Print("cluster is ", idInfo.Path)

	server, exist := Cache[idInfo.Cluster];
	if ( !exist) {
		writeErrorResponse(w, ResponseErrorCode, NoClusterFoundError.Error())
		return
	}
	pool := ZkPoolsInstance.Get(server.Name)
	conn, _ := pool.Get()

	defer conn.Close()

	conn.Delete(idInfo.Path, -1);
}

//conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
//if err != nil {
//panic(err)
//}
//defer conn.Close()
////zk 包没有提供rmr命令，只能递归删除了
//if b, _, _ := conn.Exists("/demo"); b {
//log.Println("exists /demo")
//paths, _, _ := conn.Children("/demo")
//for _, p := range paths {
//conn.Delete("/demo/"+p, -1)
//}
//err = conn.Delete("/demo", -1)
//if err != nil {
//log.Println(err)
//} else {

// Handler for the books Create action
// POST /books

//func BookCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//	book := &Book{}
//	if err := populateModelFromHandler(w, r, params, book); err != nil {
//		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
//		return
//	}
//	Bookstore[book.ISDN] = book
//	writeOKResponse(w, book)
//}
//
//// Handler for the books index action
//// GET /books
//func BookIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	books := []*Book{}
//	for _, book := range Bookstore {
//		books = append(books, book)
//	}
//	writeOKResponse(w, books)
//}
//
//// Handler for the books Show action
//// GET /books/:isdn
//func BookShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//	isdn := params.ByName("isdn")
//	book, ok := Bookstore[isdn]
//	if !ok {
//		// No book with the isdn in the url has been found
//		writeErrorResponse(w, http.StatusNotFound, "Record Not Found")
//		return
//	}
//	writeOKResponse(w, book)
//}
//
//// Writes the response as a standard JSON response with StatusOK
//func writeOKResponse(w http.ResponseWriter, m interface{}) {
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(&JsonResponse{Data: m}); err != nil {
//		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
//	}
//}
//
//// Writes the error response as a Standard API JSON response with a response code
//func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(errorCode)
//	json.
//		NewEncoder(w).
//		Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})
//}
//
////Populates a model from the params in the Handler
//func populateModelFromHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
//	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
//	if err != nil {
//		return err
//	}
//	if err := r.Body.Close(); err != nil {
//		return err
//	}
//	if err := json.Unmarshal(body, model); err != nil {
//		return err
//	}
//	return nil
//}
