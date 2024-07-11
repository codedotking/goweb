package router

import (
	"net/http"
	"path"
	"strings"

	"github.com/techiehe/goweb/context"
)

// HandlerFunc 路由处理器
type HandlerFunc func(*context.Context.Context)

// HandlersChain 路由处理器链
type HandlersChain []HandlerFunc

// Last 返回路由处理器链最后一个处理器
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// IRoute 定义路由接口
type IRoute interface {
	Use(...HandlerFunc) IRoute

	Handle(string, string, ...HandlerFunc) IRoute
	Any(string, ...HandlerFunc) IRoute
	GET(string, ...HandlerFunc) IRoute
	POST(string, ...HandlerFunc) IRoute
	DELETE(string, ...HandlerFunc) IRoute
	PATCH(string, ...HandlerFunc) IRoute
	PUT(string, ...HandlerFunc) IRoute
	OPTIONS(string, ...HandlerFunc) IRoute
	HEAD(string, ...HandlerFunc) IRoute
	Match([]string, string, ...HandlerFunc) IRoute

	StaticFile(string, string) IRoute
	StaticFileFS(string, string, http.FileSystem) IRoute
	Static(string, string) IRoute
	StaticFS(string, http.FileSystem) IRoute
}

// RouteInfo 路由信息
type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

// RouteTree 路由树
type RouteTree struct {
	Handlers HandlersChain
	basePath string
	// engine   *Engine
	root bool
}

// 断言 是否实现了 IRoute 接口
var _ IRoute = (*RouteTree)(nil)

// Use adds middleware to the group, see example code in GitHub.
func (group *RouteTree) Use(middleware ...HandlerFunc) IRoute {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (group *RouteTree) Group(relativePath string, handlers ...HandlerFunc) *RouteTree {
	return &RouteTree{
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		// engine:   group.engine,
	}
}

// BasePath returns the base path of router group.
// For example, if v := router.Group("/rest/n/v1/api"), v.BasePath() is "/rest/n/v1/api".
func (group *RouteTree) BasePath() string {
	return group.basePath
}

func (group *RouteTree) handle(httpMethod, relativePath string, handlers HandlersChain) IRoute {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in GitHub.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (group *RouteTree) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoute {
	if matched := regEnLetter.MatchString(httpMethod); !matched {
		panic("http method " + httpMethod + " is not valid")
	}
	return group.handle(httpMethod, relativePath, handlers)
}

// POST is a shortcut for router.Handle("POST", path, handlers).
func (group *RouteTree) POST(relativePath string, handlers ...HandlerFunc) IRoute {
	return group.handle(http.MethodPost, relativePath, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handlers).
func (group *RouteTree) GET(relativePath string, handlers ...HandlerFunc) IRoute {
	return group.handle(http.MethodGet, relativePath, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handlers).
func (group *RouteTree) DELETE(relativePath string, handlers ...HandlerFunc) IRoute {
	return group.handle(http.MethodDelete, relativePath, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handlers).
func (group *RouteTree) PATCH(relativePath string, handlers ...HandlerFunc) IRoute {
	return group.handle(http.MethodPatch, relativePath, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handlers).
func (group *RouteTree) PUT(relativePath string, handlers ...HandlerFunc) IRoute {
	return group.handle(http.MethodPut, relativePath, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handlers).
func (group *RouteTree) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoute {
	return group.handle(http.MethodOptions, relativePath, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handlers).
func (group *RouteTree) HEAD(relativePath string, handlers ...HandlerFunc) IRoute {
	return group.handle(http.MethodHead, relativePath, handlers)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (group *RouteTree) Any(relativePath string, handlers ...HandlerFunc) IRoute {
	for _, method := range anyMethods {
		group.handle(method, relativePath, handlers)
	}

	return group.returnObj()
}

// Match registers a route that matches the specified methods that you declared.
func (group *RouteTree) Match(methods []string, relativePath string, handlers ...HandlerFunc) IRoute {
	for _, method := range methods {
		group.handle(method, relativePath, handlers)
	}

	return group.returnObj()
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func (group *RouteTree) StaticFile(relativePath, filepath string) IRoute {
	return group.staticFileHandler(relativePath, func(c *context.Context) {
		c.File(filepath)
	})
}

// StaticFileFS works just like `StaticFile` but a custom `http.FileSystem` can be used instead..
// router.StaticFileFS("favicon.ico", "./resources/favicon.ico", Dir{".", false})
// Gin by default uses: gin.Dir()
func (group *RouteTree) StaticFileFS(relativePath, filepath string, fs http.FileSystem) IRoute {
	return group.staticFileHandler(relativePath, func(c *context.Context) {
		c.FileFromFS(filepath, fs)
	})
}

func (group *RouteTree) staticFileHandler(relativePath string, handler HandlerFunc) IRoute {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
	return group.returnObj()
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//
//	router.Static("/static", "/var/www")
func (group *RouteTree) Static(relativePath, root string) IRoute {
	return group.StaticFS(relativePath, Dir(root, false))
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default uses: gin.Dir()
func (group *RouteTree) StaticFS(relativePath string, fs http.FileSystem) IRoute {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := group.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
	return group.returnObj()
}

func (group *RouteTree) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := group.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *context.Context) {
		if _, noListing := fs.(*OnlyFilesFS); noListing {
			c.Writer.WriteHeader(http.StatusNotFound)
		}

		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.handlers = group.engine.noRoute
			// Reset index
			c.index = -1
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (group *RouteTree) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	assert1(finalSize < int(abortIndex), "too many handlers")
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouteTree) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func (group *RouteTree) returnObj() IRoute {
	if group.root {
		return group.engine
	}
	return group
}
