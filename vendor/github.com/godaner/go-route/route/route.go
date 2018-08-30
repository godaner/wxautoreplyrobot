package route

//Router
//Regist routes first , then Start

type Route struct {
	Path string
	Method string //POST.GET.DELETE,PUT etc.
	Router func (w RouteResponse, r RouteRequest)
}

var (
	Routes []Route
)
const (
	POST="POST"
	GET="GET"
	PUT="PUT"
	DELETE="DELETE"
	ANY="ANY"
)


func MakeGetRoute(path string, router func (response RouteResponse, request RouteRequest)) Route{
	return Route{
		Method:GET,
		Path:path,
		Router:router,
	}
}
func MakePostRoute(path string, router func (response RouteResponse, request RouteRequest)) Route{
	return Route{
		Method:POST,
		Path:path,
		Router:router,
	}
}
func MakeDeleteRoute(path string, router func (response RouteResponse, request RouteRequest)) Route{
	return Route{
		Method:DELETE,
		Path:path,
		Router:router,
	}
}
func MakePutRoute(path string, router func (response RouteResponse, request RouteRequest)) Route{
	return Route{
		Method:PUT,
		Path:path,
		Router:router,
	}
}

func MakeAnyRoute(path string, router func (response RouteResponse, request RouteRequest)) Route{
	return Route{
		Method:ANY,
		Path:path,
		Router:router,
	}
}

func RegistRoutes(rs...Route) {
	Routes = rs
}
