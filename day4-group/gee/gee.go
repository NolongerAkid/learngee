package gee

import (
	"log"
	"net/http"
)

type (
	Routegroup struct {
		prefix      string
		middlewares []HandlerFunc //support middleware
		parent      *Routegroup   //supprot nesting
		engine      *Engine       //all groups share a Engine instance
	}
	Engine struct {
		*Routegroup
		router *router
		groups []*Routegroup //store all groups
	}
)

// new is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.Routegroup = &Routegroup{engine: engine}
	engine.groups = []*Routegroup{engine.Routegroup}
	return engine
}

//Group is defined to create a new RouterGroup
//remember all groups share the same Engine instance
func(group *Routegroup) Group(prefix string) *Routegroup{
	engine := group.engine
	newGroup := &Routegroup{
		prefix:      group.prefix + prefix,
		middlewares: nil,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups,newGroup)
	return newGroup
}

func(group *Routegroup) addRoute(method string,comp string,handler HandlerFunc){
	pattern := group.prefix +comp
	log.Printf("Route %4s - %s",method,pattern)
	group.engine.router.addRoute(method,pattern,handler)
}

//GET defines the method to add GET request
func(group *Routegroup) GET(pattern string,handler HandlerFunc){
	group.addRoute("GET",pattern,handler)
}

//POST defines the method to add POST request
func(group *Routegroup) POST(pattern string,handler HandlerFunc){
	group.addRoute("POST",pattern,handler)
}

//Run defines the method to start a http server
func(engine *Engine) Run(addr string)(err error){
	return http.ListenAndServe(addr,engine)
}

func(engine *Engine) ServeHTTP(w http.ResponseWriter,req *http.Request){
	c := newContext(w,req)
	engine.router.handle(c)

}