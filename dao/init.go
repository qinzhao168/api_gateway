package dao

import (
	"api_gateway/basis/db"
	l "api_gateway/basis/log"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/rest"
)

var (
	engine    *db.Engine
	Clientset *kubernetes.Clientset
	log       = l.New("dao")
	err       error
)

func init() {
	initDB()
}

func initDB() {
	engine, err = db.New()
	if err != nil {
		log.Fatalf("init engine fail ,the reason is %s", err.Error())
	}
	if err = engine.Ping(); err != nil {
		log.Fatalf("access the mysql db fail ,the reason is %s", err.Error())
	}
	if err = engine.Sync(new(Node), new(App), new(Container)); err != nil {
		log.Fatalf("Sync fail :%s", err.Error())
	}
}

func initK8Sclient() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
