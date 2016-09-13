package service

import (
	"github.com/gorilla/mux"
)

func RegisterDeploymentHandler(router *mux.Router) {
	RegisterHttpHandler(router, "/deployment", HTTP_POST, Create)                  //创建应用
	RegisterHttpHandler(router, "/deployment", HTTP_DELETE, Delete)                //删除应用
	RegisterHttpHandler(router, "/deployment", HTTP_PUT, Update)                   //更新状态  包括 启动 停止  更新应用配置
	RegisterHttpHandler(router, "/deployment/scale", HTTP_PUT, Scale)              //更新container个数
	RegisterHttpHandler(router, "/deployment/redployment", HTTP_POST, ReDeloyment) //重新部署
}

func Create(req *http.Request) (code string, ret interface{}) {
	return
}

func Delete(req *http.Request) (code string, ret interface{}) {
	return
}

func Update(req *http.Request) (code string, ret interface{}) {
	return
}

func Scale(req *http.Request) (code string, ret interface{}) {
	return
}

func ReDeloyment(req *http.Request) (code string, ret interface{}) {
	return
}
