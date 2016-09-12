package dao

import (
	"fmt"
	"time"

	"api_gateway/basis/errors"
	"api_gateway/basis/jsonx"
)

//
type Port struct {
	ExposePort int
}

func (port *Port) String() string {
	return fmt.Sprintf("[ExposePort:%d]", port.ExposePort)
}

//ContainerDto is a DTO of Container
type ContainerDto struct {
	Id      string
	Image   string
	AppName string
	Ports   []*Port
	Status  string
}

//Container is struct of docker container
type Container struct {
	Id      int       `json:"id" xorm:"pk not null varchar(12)"`
	Ip      string    `json:"ip" xorm:"varchar(256)"`
	Image   string    `json:"image" xorm:"varchar(1024)"`
	Name    string    `json:"name" xorm:"varchar(1024)"`
	Ports   string    `json:"ports" xorm:"varchar(1024)"`
	Mounts  string    `json:"mounts" xorm:"varchar(1024)"`
	Status  int       `json:"status" xorm:"int(1)"`
	Created time.Time `json:"created" xorm:"created"`
}

func (container *Container) String() string {
	containerStr, err := jsonx.ToJson(container)
	if err != nil {
		log.Errorf("container to string err :%s", err.Error())
		return ""
	}
	return containerStr
}

func (container *Container) Insert() error {
	_, err := engine.Insert(container)

	if err != nil {
		return err
	}

	return nil
}

func (container *Container) Delete() error {
	_, err := engine.Id(container.Id).Delete(container)

	if err != nil {
		return err
	}

	return nil
}

func (container *Container) Update() error {
	_, err := engine.Id(container.Id).Update(container)
	if err != nil {
		return err
	}

	return nil
}

func (container *Container) QueryOne() (*Container, error) {
	has, err := engine.Id(container.Id).Get(container)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("the query data not exist")
	}

	return container, nil
}

func (container *Container) QuerySet() ([]*Container, error) {
	containerSet := []*Container{}
	err := engine.Where("1 and 1 order by ip desc").Find(&containerSet)

	if err != nil {
		return nil, err
	}

	return containerSet, nil
}
