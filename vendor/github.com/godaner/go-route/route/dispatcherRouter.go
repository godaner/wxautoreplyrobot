package route

import (
	"net/http"
	"log"
)


type DispatcherRouter struct {

}


func (dispatcherHandler DispatcherRouter)ServeHTTP(w http.ResponseWriter,r *http.Request){

	prepareRequest(r)

	//start route
	for _,route:= range Routes {
		if matchRoute(route,r) {//match request path and request method:post get delete put etc.
			route.Router(makeRouteResponse(w),makeRouteRequest(r))
			return
		}
	}

	err := "Method : "+r.Method + "#" + r.URL.Path + " is not exits !"
	w.Write([]byte(err))

}
func prepareRequest(r *http.Request) {
	//parse form
	r.ParseForm()
	//log
	log.Println("====> path is : ",r.URL.Path," ====>")
	log.Println(r)
	log.Println(r.Form)
}

func matchRoute( route Route,r *http.Request) bool{
	if r.URL.Path!=route.Path{
		return false
	}
	if route.Method == ANY{
		return true
	}
	return r.Method == route.Method
}
func makeRouteRequest(request *http.Request) RouteRequest {

	//get form param
	params:=map[string]interface{}{}

	for k,v :=range request.Form {
		if len(v)==1 {
			params[k] = v[0]
		}else{
			params[k] = v
		}
	}
	return RouteRequest{
		Request:request,
		Params:params,
	}
}
func makeRouteResponse(writer http.ResponseWriter) RouteResponse {
	return RouteResponse{writer}
}