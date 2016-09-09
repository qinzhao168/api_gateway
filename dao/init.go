package dao

import (
	"api_gateway/basis/db"
	l "api_gateway/basis/log"
)

var (
	engine *db.Engine
	log    = l.New("dao")
)

func init() {
	engine, err := db.New()
	if err != nil {
		log.Fatalf("init engine fail ,the reason is %s", err.Error())
	}
	if err = engine.Ping(); err != nil {
		log.Fatalf("access the mysql db fail ,the reason is %s", err.Error())
	}
	engine.Sync(new(Node), new(App), new(Container))
}
