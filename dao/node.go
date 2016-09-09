package dao

import (
	"basis/jsonx"
	"basis/log"
)

//Node node is server node
type Node struct {
	Ip          string
	Hostname    string
	Os          string
	CpuNum      int
	CpuFree     float64
	CpuUse      float64
	MemoryTotal int64
	MemoryFree  int64
	DiskTotal   int64
	DiskFree    int64
}

func (node *Node) String() string {
	nodeStr, err := jsonx.ToJson(node)
	if err != nil {
		log.New("").Errorf("node to string err :%s", err.Error())
		return ""
	}
	return nodeStr
}
