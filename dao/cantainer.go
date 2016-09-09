package dao

import (
	"fmt"
	"time"

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
	Id      string    `json:"id" xorm:"pk not null varchar(12)"`
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
