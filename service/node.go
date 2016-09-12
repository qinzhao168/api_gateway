package service

import (
	"encoding/json"
	"net/http"

	"api_gateway/dao"
	"github.com/gorilla/mux"
)

func RegisterNodeHandler(router *mux.Router) {
	RegisterHttpHandler(router, "/node", HTTP_POST, InserNode)
	RegisterHttpHandler(router, "/node", HTTP_DELETE, DeleteNode)
	RegisterHttpHandler(router, "/node", HTTP_PUT, UpdateNode)
	RegisterHttpHandler(router, "/node", HTTP_GET, QueryOneNode)
	RegisterHttpHandler(router, "/node/list", HTTP_GET, QuerySetNode)
}

func InserNode(req *http.Request) (code string, ret interface{}) {
	node := &dao.Node{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(node); err != nil {
		code = StatusBadRequest
		ret = OK
		return
	}

	if err := node.Insert(); err != nil {
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	return
}

func DeleteNode(req *http.Request) (code string, ret interface{}) {
	node := new(dao.Node)
	node.Ip = req.FormValue("ip")
	err := node.Delete()
	if err != nil {
		code = StatusInternalServerError
		ret = OK
		return
	}

	code = StatusNoContent
	ret = OK
	return
}

func UpdateNode(req *http.Request) (code string, ret interface{}) {
	node := &dao.Node{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(node); err != nil {
		code = StatusBadRequest
		ret = OK
		return
	}

	if err := node.Update(); err != nil {
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusCreated
	ret = OK
	return
}

func QueryOneNode(req *http.Request) (code string, ret interface{}) {
	node := new(dao.Node)
	node.Ip = req.FormValue("ip")
	node, err := node.QueryOne()
	if err != nil {
		code = StatusInternalServerError
		ret = OK
		return
	}

	code = StatusOK
	ret = node
	return
}

func QuerySetNode(req *http.Request) (code string, ret interface{}) {
	node := new(dao.Node)
	nodes, err := node.QuerySet()
	if err != nil {
		code = StatusInternalServerError
		ret = OK
		return
	}

	code = StatusOK
	ret = nodes
	return
}
