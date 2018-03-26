package router

import (
	. "handlers"

	"github.com/julienschmidt/httprouter"
)

/*
Define all the routes here.
A new Route entry passed to the routes slice will be automatically
translated to a handler with the NewRouter() function
*/
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func AllRoutes() Routes {
	routes := Routes{
		Route{"Index", "GET", "/", Index},
		Route{"BookIndex", "GET", "/books", BookIndex},
		Route{"Bookshow", "GET", "/books/:isdn", BookShow},
		Route{"Bookshow", "POST", "/books", BookCreate},

		Route{"Children", "GET", "/api/zk/cluster/:cluster", ListPath},
		//Route{"UpdateZkNode", "POST", "/api/zk/cluster/:cluster/createZkNode/", CreateZkNode},
		//Route{"DeleteZkNode", "POST", "/api/zk/cluster/:cluster/deleteZkNode/", DeleteZkNode},
		//Route{"ReadZkNode", "GET", "/api/zk/cluster/:cluster/path/:path/node", ReadZkNode},
		//

		Route{"CreateZkNode", "POST", "/api/zks/createZkNode", CreateZkNode},
		//Route{"UpdateZkNode", "POST", "/api/zk/cluster/:cluster/updateZkNode/", UpdateZkNode},
		Route{"DeleteZkNode", "POST", "/api/zks/deleteZkNode", DeleteZkNode},
		Route{"ReadZkNode", "GET", "/api/zks/cluster/:cluster/path/:path/node", ReadZkNode},

		Route{"zks", "GET", "/api/zks", ZksNew},
		Route{"createZk", "POST", "/api/zks/createZk", CreateZk},
		Route{"updateZk", "POST", "/api/zks/updateZk", UpdateZk},
		Route{"deleteZk", "POST", "/api/zks/deleteZk", DeleteZk},
		//Route{"options", "OPTIONS", "/api", OptionsRequest},
		//Route{"deleteZk", "POST", "/api/zks/deleteZk", DeleteZk},

		Route{"zkByName", "GET", "/api/zk/name/:name", FindByName},
		Route{"Zms", "GET", "/api/zms/cluster/:cluster", ListZms},
	}
	return routes
}
