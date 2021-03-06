package ago

import (
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

//Router Router include mux.Router
type Router struct {
	Router *mux.Router `inject:""`
}

//NewHandler New a handler with route array
func (router *Router) NewHandler(routes []*Route) *mux.Router {
	router.Router.StrictSlash(true)

	for _, route := range routes {

		route = route.Middleware(Logger)

		var handler http.Handler

		if route.HandlerFunc != nil {
			handler = route.HandlerFunc
		} else {
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		}

		if route.Controller != nil {
			handler = HTTPMiddleware(handler, route)
		}

		if len(route.Middlewares) > 0 {
			for _, m := range route.Middlewares {
				handler = m(handler)
			}
		}

		if route.Method == "*" || route.Method == "" {
			router.Router.
				Path(route.Pattern).
				Name(route.Method + ":" + route.MethodName).
				Handler(handler)
		} else {
			router.Router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Method + ":" + route.MethodName).
				Handler(handler)
		}
	}

	return router.Router
}

//HTTPMiddleware HTTP request middleware
func HTTPMiddleware(next http.Handler, route *Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		t := reflect.Indirect(reflect.ValueOf(route.Controller)).Type()
		controller := reflect.New(t)

		//Call Controller Init Method
		initMethod := controller.MethodByName("Init")
		vars := make([]reflect.Value, 3)
		vars[0] = reflect.ValueOf(w)
		vars[1] = reflect.ValueOf(r)
		vars[2] = reflect.ValueOf(mux.Vars(r))
		initMethod.Call(vars)

		//Call Controller Prepare Method
		vars = make([]reflect.Value, 0)
		method := controller.MethodByName("Prepare")
		method.Call(vars)

		methodName := route.MethodName

		//Call Controller RESTFul Methods
		if route.Method == "GET" {
			if len(methodName) == 0 {
				methodName = "Get"
			}
			method = controller.MethodByName(methodName)
			method.Call(vars)
		} else if route.Method == "POST" {
			if len(methodName) == 0 {
				methodName = "Post"
			}
			method = controller.MethodByName(methodName)
			method.Call(vars)
		} else if route.Method == "HEAD" {
			if len(methodName) == 0 {
				methodName = "Head"
			}
			method = controller.MethodByName(methodName)
			method.Call(vars)
		} else if route.Method == "DELETE" {
			if len(methodName) == 0 {
				methodName = "Delete"
			}
			method = controller.MethodByName(methodName)
			method.Call(vars)
		} else if route.Method == "PUT" {
			if len(methodName) == 0 {
				methodName = "Put"
			}
			method = controller.MethodByName(methodName)
			method.Call(vars)
		} else if route.Method == "PATCH" {
			if len(methodName) == 0 {
				methodName = "Patch"
			}
			method = controller.MethodByName(methodName)
			method.Call(vars)
		} else if route.Method == "OPTIONS" {
			if len(methodName) == 0 {
				methodName = "Options"
			}
			method = controller.MethodByName(methodName)
			method.Call(vars)
		}

		//Call Controller Finish Methods
		method = controller.MethodByName("Finish")
		method.Call(vars)
	})
}
