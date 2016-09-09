package dao

import (
	"basis/jsonx"
	"basis/log"
)

//App is struct of application
type App struct {
	Id            int
	Name          string
	Region        string
	Memory        int
	Cpu           int
	InstanceCount int
	Envs          string
	Ports         string
	Image         string
	Status        int
	UserName      string
	Remark        string
}

func (app *App) String() string {
	appStr, err := jsonx.ToJson(app)
	if err != nil {
		log.New("").Errorf("node to string err :%s", err.Error())
		return ""
	}
	return appStr
}
