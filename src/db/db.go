package db

import (
	"database/sql"
	"strconv"
	"time"

	"container/list"

	. "base"
	"models"

	//"fmt"
	"log"
	//"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

var Cache map[string]*models.Server = make(map[string]*models.Server)

//func InitDb() {
//
//	server := new(models.Server)
//	server.Name = "local"
//	server.Ips = "localhost:2181"
//	server.Id = 12345
//	server.Version = "3.4.6"
//	Cache["local"] = server
//
//}

func InitDb() {

	log.Print("init method executed")

	var err error
	Db, err = sql.Open("mysql", "root:root@tcp(10.9.20.111)/rock?charset=utf8&parseTime=true")

	if err != nil {
		log.Fatal("connect to mysql failed")
		return;
	}
	Db.SetMaxIdleConns(100)
	err = Db.Ping()
	if err != nil {
		log.Panic(err)
		return
	}
	zks, err := ListZks();
	if ( err != nil ) {
		log.Fatal("start failed for: init cache failed", err)
	} else {
		for e := zks.Front(); e != nil; e = e.Next() {
			t := e.Value.(*models.Server);
			Cache[t.Name] = t
		}

	}
}

func Get(nameStr string, ip string) (*models.Server, error) {

	//rows, err := Db.Query("select id,name,ips from zookeeper where name ='test'");
	row := Db.QueryRow("select id,ips,name,version from zookeeper where name =?", nameStr);
	if (row == nil) {
		return nil, NoReulstFoundError
	}
	var server models.Server
	row.Scan(&server.Id, &server.Ips, &server.Name, &server.Version)

	return &server, nil

}

func ListZks() (*list.List, error) {

	//rows, err := Db.Query("select id,name,ips from zookeeper where name ='test'");

	rows, err := Db.Query("select id,ips,name,version from zookeeper ");
	if err != nil {
		return nil, err
	}

	var servers = list.New()

	defer rows.Close()

	for rows.Next() { //开始循环
		var server models.Server
		rerr := rows.Scan(&server.Id, &server.Ips, &server.Name, &server.Version)
		if rerr == nil {

			servers.PushBack(&server)
		}
	}
	return servers, nil

}

func makeStmt(name string, ip string) (stmt *sql.Stmt, err error) {

	if (name != "%%" && ip != "%%") {
		stmt, err = Db.Prepare("select id,ips,name,version from zookeeper where name like ? and ips like ?");
	} else if (name != "%%") {
		stmt, err = Db.Prepare("select id,ips,name,version from zookeeper where name like ? ");
	} else if (ip != "%%") {
		stmt, err = Db.Prepare("select id,ips,name,version from zookeeper where ips like ?");
	} else {
		stmt, err = Db.Prepare("select id,ips,name,version from zookeeper ");
	}
	return
}

func FindZks(name string, ip string) ([]*models.Server, error) {

	stmt, err := makeStmt(name, ip)
	if err != nil {
		return nil, err
	}

	//rows, err := Db.Query("select id,name,ips from zookeeper ");
	rows, err := stmt.Query("%"+name+"%", "%"+ip+"%");
	if err != nil {
		return nil, err
	}

	var servers = make([]*models.Server, 0)
	defer stmt.Close()
	defer rows.Close()

	for rows.Next() { //开始循环
		var server models.Server
		rerr := rows.Scan(&server.Id, &server.Ips, &server.Name, &server.Version)
		if rerr == nil {
			servers = append(servers, &server)
		}
	}
	return servers, nil

}

func ArrayZks() ([]*models.Server, error) {

	return FindZks("", "");
	////rows, err := Db.Query("select id,name,ips from zookeeper ");
	//rows, err := Db.Query("select id,ips,name,version from zookeeper ");
	//if err != nil {
	//	return nil, err
	//}
	//
	//var servers = make([]*models.Server, 0)
	//
	//defer rows.Close()
	//
	//for rows.Next() { //开始循环
	//	var server models.Server
	//	rerr := rows.Scan(&server.Id, &server.Ips, &server.Name, &server.Version)
	//	if rerr == nil {
	//		servers = append(servers, &server)
	//	}
	//}
	//return servers, nil

}
func Delete(name string, id int64) (string, error) {

	stmt, err := Db.Prepare("delete from zookeeper where id = ?")
	if err != nil {
		return "", err
	}

	res, err := stmt.Exec(id);
	//rows, err := common.Db.Query("select id,name,ips from zookeeper where name =?", "test");
	if err != nil {
		//log.Panic("delete failed,", strconv.FormatInt(id, 10))
		log.Panic("delete failed,", id)
	}

	deleteId, _ := res.RowsAffected()
	delete(Cache, name)
	return strconv.FormatInt(deleteId, 10), nil;
}

func Insert(server *models.Server) int64 {

	stmt, err := Db.Prepare("insert into zookeeper(ips,name,version,create_date) VALUES (?,?,?,?)");
	res, err := stmt.Exec(server.Ips, server.Name, server.Version, time.Now());
	if err != nil {
		log.Fatal(err)
	}

	id, _ := res.LastInsertId()

	server.Id = id
	Cache[server.Name] = server
	return id
}

func Update(server *models.Server) int64 {

	stmt, err := Db.Prepare("update zookeeper set ips=?,name=?,version=? where id=?");
	res, err := stmt.Exec(server.Ips, server.Name, server.Version, server.Id);
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.RowsAffected()
	server.Id = id
	Cache[server.Name] = server
	return id
}
