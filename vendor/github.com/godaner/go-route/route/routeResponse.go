package route

import (
	"net/http"
	"encoding/json"
	"github.com/godaner/go-util"
)

type RouteResponse struct {
	ResponseWriter http.ResponseWriter
}

func (response RouteResponse ) WriteJsonStr(obj interface{}){
	json,_ :=json.Marshal(obj)
	jsonStr :=string(json)
	response.ResponseWriter.Write([]byte(jsonStr))
}
func (response RouteResponse) WriteString(str string){
	response.ResponseWriter.Write([]byte(str))
}
func (response RouteResponse) WriteTemplate(template string){
	go_util.WriteTemplate(response.ResponseWriter,template)
}