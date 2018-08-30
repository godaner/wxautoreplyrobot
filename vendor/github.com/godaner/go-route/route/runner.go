package route

import (
	"log"
	"net/http"
	"github.com/godaner/go-util"
)

func Start(addr string){
	log.Println("Server api is : ", Routes)
	log.Println("Server will start at : [ ", addr," ]")
	err := http.ListenAndServe(addr,DispatcherRouter{})
	go_util.CheckErr(err)
}
