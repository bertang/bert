package router

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
)

var (
	rootRouter = &router{prefix: "/"}
)

type RouteFunc func(iRouter IRouter)
type IRouter interface {
	Group(string, RouteFunc)
	Middleware(handlers ...context.Handler) IRouter
	Prefix(prefix string) IRouter
	Handler(controller interface{})
	Services(services ...interface{}) IRouter
}

//RegisterRouter 向框架中注入
func RegisterRouter(app *iris.Application) {
	//开始注册
	root := mvc.New(app.Party(rootRouter.prefix))
	for k := range rootRouter.routers {
		register(rootRouter.routers[k], root)
	}

}

//循环获取
func register(r *router, app *mvc.Application) {
	//创建party
	pMvc := app.Party(r.prefix)
	if len(r.middlewares) > 0 {
		pMvc.Router.Use(r.middlewares...)
	}

	if len(r.services) > 0 {
		pMvc.Register(r.services...)
	}

	for k := range r.handler {
		pMvc.Handle(r.handler[k])
	}

	//子路由
	for k := range r.routers {
		register(r.routers[k], pMvc.Party("/"))
	}
}

//路由
type router struct {
	parent      *router
	prefix      string
	middlewares []context.Handler
	services    []interface{}
	handler     []interface{}
	routers     []*router
}

//分组
func (r *router) Group(prefix string, routeFunc RouteFunc) {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	subRouter := &router{prefix: prefix}

	subRouter.parent = r
	r.routers = append(r.routers, subRouter)
	routeFunc(subRouter)
}

//Middleware 中间件
func (r *router) Middleware(handlers ...context.Handler) IRouter {
	rSub := &router{middlewares: handlers}
	r.routers = append(r.routers, rSub)
	return rSub
}

//Prefix 路由名称
func (r *router) Prefix(name string) IRouter {
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}
	//添加判断避免冲突
	rSub := &router{prefix: name}
	r.routers = append(r.routers, rSub)
	return rSub
}

//Services 需要注入到controller 的service
func (r *router) Services(services ...interface{}) IRouter {
	rSub := &router{services: services}
	r.routers = append(r.routers, rSub)
	return rSub
}

//Handler 处理controller
func (r *router) Handler(controller interface{}) {
	r.handler = append(r.handler, controller)
}

//Group 路由分组
func Group(prefix string, routerFunc RouteFunc) {
	r := &router{prefix: prefix}

	r.Group(prefix, routerFunc)
}

//Middleware 中间件
func Middleware(handlers ...context.Handler) IRouter {
	newRouter := newSubRouter()
	return newRouter.Middleware(handlers...)
}

//路由前缀，用于整个controller的路径前缀
func Prefix(name string) IRouter {
	newRouter := newSubRouter()
	return newRouter.Prefix(name)
}

func newSubRouter() *router {
	newRouter := &router{prefix: "/"}
	rootRouter.routers = append(rootRouter.routers, newRouter)
	return newRouter
}

//Services 注册service
func Services(services ...interface{}) IRouter {
	newRouter := newSubRouter()
	newRouter.services = append(newRouter.services, services...)

	return newRouter
}

//Handler 处理函数
func Handler(controller interface{}) {
	rootRouter.Handler(controller)
}
