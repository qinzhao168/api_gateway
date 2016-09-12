package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_gateway/dao"
	"github.com/gorilla/mux"
)

func RegisterContainerHandler(router *mux.Router) {
	RegisterHttpHandler(router, "/container", HTTP_POST, InserContainer)
	RegisterHttpHandler(router, "/container", HTTP_DELETE, DeleteContainer)
	RegisterHttpHandler(router, "/container", HTTP_PUT, UpdateContainer)
	RegisterHttpHandler(router, "/container", HTTP_GET, QueryOneContainer)
	RegisterHttpHandler(router, "/container/list", HTTP_GET, QuerySetContainer)
}

func InserContainer(req *http.Request) (code string, ret interface{}) {
	container := &dao.Container{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(container); err != nil {
		code = StatusBadRequest
		ret = OK
		return
	}

	if err := container.Insert(); err != nil {
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	return
}

func DeleteContainer(req *http.Request) (code string, ret interface{}) {
	container := new(dao.Container)
	container.Id, _ = strconv.Atoi(req.FormValue("id"))
	err := container.Delete()
	if err != nil {
		code = StatusInternalServerError
		ret = OK
		return
	}

	code = StatusNoContent
	ret = OK
	return
}

func UpdateContainer(req *http.Request) (code string, ret interface{}) {
	container := &dao.Container{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(container); err != nil {
		code = StatusBadRequest
		ret = OK
		return
	}

	if err := container.Update(); err != nil {
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusCreated
	ret = OK
	return
}

func QueryOneContainer(req *http.Request) (code string, ret interface{}) {
	container := new(dao.Container)
	container.Id, _ = strconv.Atoi(req.FormValue("id"))
	container, err := container.QueryOne()
	if err != nil {
		code = StatusInternalServerError
		ret = OK
		return
	}

	code = StatusCreated
	ret = container
	return
}

func QuerySetContainer(req *http.Request) (code string, ret interface{}) {
	container := new(dao.Container)
	containers, err := container.QuerySet()
	if err != nil {
		code = StatusInternalServerError
		ret = OK
		return
	}

	code = StatusOK
	ret = containers
	return
}
