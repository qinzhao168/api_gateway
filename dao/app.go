package dao

import (
	"api_gateway/basis/errors"
	"api_gateway/basis/jsonx"
)

//App is struct of application
type App struct {
	Id            int    `json:"id" xorm:"pk not null autoincr int(11)"`
	Name          string `json:"name" xorm:"varchar(1024)"`
	Region        string `json:"region" xorm:"varchar(1024)"`
	Memory        int    `json:"memory" xorm:"int(11)"`
	Cpu           int    `json:"cpu" xorm:"int(11)"`
	InstanceCount int    `json:"instanceCount" xorm:"int(11)"`
	Envs          string `json:"envs" xorm:"varchar(1024)"`
	Ports         string `json:"ports" xorm:"varchar(1024)"`
	Image         string `json:"image" xorm:"varchar(1024)"`
	Status        int    `json:"status" xorm:"int(1)"`
	UserName      string `json:"userName" xorm:"varchar(1024)"`
	Remark        string `json:"remark" xorm:"varchar(1024)"`
}

func (app *App) String() string {
	appStr, err := jsonx.ToJson(app)
	if err != nil {
		log.Errorf("node to string err :%s", err.Error())
		return ""
	}
	return appStr
}

func (app *App) Insert() error {
	_, err := engine.Insert(app)

	if err != nil {
		return err
	}

	return nil
}
func (app *App) Delete() error {
	_, err := engine.Id(app.Id).Delete(app)

	if err != nil {
		return err
	}

	return nil
}
func (app *App) Update() error {
	_, err := engine.Id(app.Id).Update(app)
	if err != nil {
		return err
	}

	return nil
}
func (app *App) QueryOne() (*App, error) {
	has, err := engine.Id(app.Id).Get(app)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("the query data not exist")
	}

	return app, nil
}

func (app *App) QuerySet() ([]*App, error) {
	appSet := []*App{}
	err := engine.Where("1 and 1 order by id desc").Find(&appSet)

	if err != nil {
		return nil, err
	}

	return appSet, nil
}
