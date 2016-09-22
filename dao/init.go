package dao

import (
	"api_gateway/basis/db"
	l "api_gateway/basis/log"
	"flag"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"
)

var (
	engine    *db.Engine
	Clientset *kubernetes.Clientset
	log       = l.New("dao")
	err       error
)

func init() {
	initDB()
	initK8Sclient()
}

func initDB() {
	engine, err = db.New()
	if err != nil {
		log.Fatalf("init engine fail ,the reason is %s", err.Error())
	}
	if err = engine.Ping(); err != nil {
		log.Fatalf("access the mysql db fail ,the reason is %s", err.Error())
	}
	if err = engine.Sync(new(App)); err != nil {
		log.Fatalf("Sync fail :%s", err.Error())
	}

	engine.ShowSQL(true)
}

func initK8Sclient() {
	kubeconfig := flag.String("kubeconfig", "./etc/config", "absolute path to the kubeconfig file")
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("192.168.1.113:8080", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
