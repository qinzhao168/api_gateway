package model

import (
	"fmt"
	"time"

	"basis/jsonx"
	"basis/log"
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
	Id      string
	Ip      string
	Image   string
	Name    string
	Ports   string
	Mounts  string
	Status  int
	Created time.Time
}

func (container *Container) String() string {
	containerStr, err := jsonx.ToJson(container)
	if err != nil {
		log.New("").Errorf("container to string err :%s", err.Error())
		return ""
	}
	return containerStr
}
