package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterNodeHandler(router *mux.Router) {
	RegisterHttpHandler(router, "/node", HTTP_POST, InserNode)
	RegisterHttpHandler(router, "/node", HTTP_DELETE, DeleteNode)
	RegisterHttpHandler(router, "/node", HTTP_PUT, UpdateNode)
	RegisterHttpHandler(router, "/node", HTTP_GET, QueryNode)
}

func InserNode(req *http.Request) (code string, ret interface{}) {
	return
}

func DeleteNode(req *http.Request) (code string, ret interface{}) {
	return
}

func UpdateNode(req *http.Request) (code string, ret interface{}) {
	return
}

func QueryNode(req *http.Request) (code string, ret interface{}) {
	return
}
