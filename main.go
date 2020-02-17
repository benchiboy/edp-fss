// qmlServer project main.go
package main

import (
	"flag"
	"io"
	"log"

	"edp-fss/control/rpt"
	"edp-fss/service/dbcomm"
	"net/http"
	"os"

	goconf "github.com/pantsing/goconf"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	http_srv   *http.Server
	dbUrl      string
	listenPort int
	idleConns  int
	openConns  int
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)
	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "sms.log",
		MaxSize:    500, // megabytes
		MaxBackups: 50,
		MaxAge:     90, //days
	}))
	envConf := flag.String("env", "config-ci.json", "select a environment config file")
	flag.Parse()
	log.Println("config file ==", *envConf)
	c, err := goconf.New(*envConf)
	if err != nil {
		log.Fatalln("读配置文件出错", err)
	}

	//填充配置文件
	c.Get("/config/LISTEN_PORT", &listenPort)
	c.Get("/config/DB_URL", &dbUrl)
	c.Get("/config/OPEN_CONNS", &openConns)
	c.Get("/config/IDLE_CONNS", &idleConns)

	dbcomm.InitDB(dbUrl, openConns, openConns)

}

func go_WebServer() {
	log.Println("Listen Service start...")
	http.HandleFunc("/edp-fss/api/v1/qry_agrt", rpt.QryAgrt)
	http.HandleFunc("/edp-fss/api/v1/crt_rpt", rpt.CrtRptFile)
	http.HandleFunc("/edp-fss/api/v1/qry_rpt", rpt.QryRptStatus)
	http.HandleFunc("/edp-fss/api/v1/qry_rpt_file", rpt.QryRptFile)

	http_srv = &http.Server{
		Addr: ":8087",
	}
	log.Printf("listen:")
	if err := http_srv.ListenAndServe(); err != nil {
		log.Printf("listen: %s\n", err)
	}
}

func main() {

	go_WebServer()

}
