package httpx

import (
	"net/http"
	"net/http/httputil"

	l "api_gateway/basis/log"
)

var log = l.New("basis/httpx")

func DumpRequest(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Info("dumpReq err:", err)
	} else {
		log.Info("dumpReq ok:", string(dump))
	}
}
