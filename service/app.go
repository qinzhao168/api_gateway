package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api_gateway/dao"
	"github.com/gorilla/mux"
)

func RegisterAppHandler(router *mux.Router) {
	RegisterHttpHandler(router, "/app", HTTP_POST, InserApp)
	RegisterHttpHandler(router, "/app", HTTP_DELETE, DeleteApp)
	RegisterHttpHandler(router, "/app", HTTP_PUT, UpdateApp)
	RegisterHttpHandler(router, "/app", HTTP_GET, QueryOneApp)
	RegisterHttpHandler(router, "/app/list", HTTP_GET, QuerySetApp)
}

func InserApp(req *http.Request) (code string, ret interface{}) {
	app := &dao.App{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(app); err != nil {
		code = StatusBadRequest
		ret = OK
		return
	}

	if err := app.Insert(); err != nil {
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusCreated
	ret = OK
	return
}

func DeleteApp(req *http.Request) (code string, ret interface{}) {
	app := new(dao.App)
	app.Id, _ = strconv.Atoi(req.FormValue("id"))
	err := app.Delete()
	if err != nil {
		code = StatusInternalServerError
		ret = OK
		return
	}

	code = StatusNoContent
	ret = OK
	return
}

func UpdateApp(req *http.Request) (code string, ret interface{}) {
	app := &dao.App{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(app); err != nil {
		code = StatusBadRequest
		ret = OK
		return
	}

	if err := app.Update(); err != nil {
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusCreated
	ret = OK
	return
}

func QueryOneApp(req *http.Request) (code string, ret interface{}) {
	app := new(dao.App)
	app.Id, _ = strconv.Atoi(req.FormValue("id"))
	app, err := app.QueryOne()
	if err != nil {
		code = StatusNotFound
		ret = OK
		return
	}

	code = StatusOK
	ret = app
	return
}

func QuerySetApp(req *http.Request) (code string, ret interface{}) {
	app := new(dao.App)
	apps, err := app.QuerySet()
	if err != nil {
		code = StatusNotFound
		ret = OK
		return
	}

	code = StatusOK
	ret = apps
	return
}
