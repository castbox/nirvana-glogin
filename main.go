package main

import (
	"flag"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/castbox/nirvana-gcore/glog"
	"github.com/castbox/nirvana-kite"
	"glogin/config"
	"glogin/constant"
	"glogin/db"
	"glogin/handler/cgi"
	"glogin/handler/moss"
	"glogin/internal/configure"
	glogin2 "glogin/pbs/glogin"
	"net/http"
)

func main() {
	// 开发模式设置
	initDevelopment()
	log.Infow("server started")
	// 初始化配置
	initConfig()
	// 初始化DB
	initDB()
	// 启动http服务
	go startHttp()
	log.Infow("glogin server started")
	// RPC服务
	kiteServe()
}

func initDevelopment() {
	development := flag.Int("dev", 0, "是否处于开发模式，决定日志等级。")
	flag.Parse()
	if development == nil || *development == 0 {
		log.ResetToProduction()
	} else {
		log.ResetToDevelopment()
	}
}

func initDB() {
	db.InitMongo()
}

func startHttp() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	pprof.Register(router)
	router.POST("/:service/:action", cgi.ServiceHandler)
	//router.POST("/anti/:action", cgi.AutiHandler)
	router.POST("/ids/:action", cgi.IdsHandler)
	router.POST("/cfg", cgi.CfgHandler)
	router.POST("/agreement", cgi.CheckHandler)
	router.POST("/token", cgi.TokenHandler)
	router.POST("/set_vsn", cgi.SetVsn)
	router.GET("/get_vsn", cgi.GetVsn)
	router.GET("/metrics", prometheusHandler())
	router.GET("/heart", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"err_code": constant.ErrCodeOk}) })
	if err := router.Run(config.WebPort()); err != nil {
		log.Fatalw("failed to run http server", "err", err)
		panic(err)
	}
}

func initConfig() {
	config.WatchStaticConfig()
	configure.WatchDynamicPubDir()
}

func kiteServe() {
	glogin2.RegLoginServer(cgi.Login{})
	glogin2.RegBindServer(cgi.Bind{})
	glogin2.RegGmpServer(moss.Gmp{})
	// 本地测试
	//kite.StartServer(ProcessStop{}, &kite.Destination{
	//	Cluster:      constant.ClusterName,
	//	ServiceName:  constant.ServerName,
	//	ServiceIndex: constant.Index,
	//})
	// 部署运行
	kite.Serve(ProcessStop{})
}

// ProcessStop 执行服务退出流程
type ProcessStop struct {
}

func (p ProcessStop) Stop() {

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
