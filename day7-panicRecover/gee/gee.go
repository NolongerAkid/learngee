package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
		router        *router
		groups        []*Routegroup //store all groups
		htmlTemplates *template.Template//for html render
		funcMap       template.FuncMap//for html render
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
func (group *Routegroup) Group(prefix string) *Routegroup {
	engine := group.engine
	newGroup := &Routegroup{
		prefix:      group.prefix + prefix,
		middlewares: nil,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *Routegroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//GET defines the method to add GET request
func (group *Routegroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//POST defines the method to add POST request
func (group *Routegroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//Create static handler
func (group *Routegroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		//Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *Routegroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	//Register GET handlers
	group.GET(urlPattern, handler)

}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap){
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string){
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

//Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}


func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)

}

//Use is defined to add middleware to the group
func (group *Routegroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
