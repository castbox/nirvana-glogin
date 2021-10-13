package main

import (
	"flag"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gmoss/v3"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"glogin/config"
	"glogin/constant"
	"glogin/db"
	"glogin/handler/cgi"
	"glogin/handler/moss"
	"glogin/internal/configure"
	glogin2 "glogin/pbs/glogin"
	"net/http"
	"os"
	"runtime"
)

func main() {
	// 开发模式设置
	initDevelopment()
	// 注册服务
	startEngine()
	log.Infow("server started")
	// 初始化配置
	initConfig()
	log.Infow("config init ok", "config", config.GetAll())
	// 初始化DB
	initDB()
	// 启动http服务
	go startHttp()
	log.Infow("glogin server started")
	// 循环
	select {}
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
	router.POST("/anti/:action", cgi.AutiHandler)
	router.POST("/ids/:action", cgi.IdsHandler)
	router.POST("/cfg", cgi.CfgHandler)
	router.POST("/token", cgi.TokenHandler)
	router.GET("/metrics", prometheusHandler())
	router.GET("/heart", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"err_code": constant.ErrCodeOk}) })
	if err := router.Run(config.WebPort()); err != nil {
		log.Fatalw("failed to run http server", "err", err)
		panic(err)
	}
}

func initConfig() {
	config.Init()
	configure.WatchPubCfg()
}

func startEngine() {
	gmoss.RegSignal(func() {
		log.Infow("glogin server stop")
		os.Exit(0)
	})
	runtime.GOMAXPROCS(runtime.NumCPU())
	serviceReg := gmoss.NewServiceRegister()
	serviceReg.Regist(&glogin2.Login_serviceDesc, cgi.Login{})
	serviceReg.Regist(&glogin2.Bind_serviceDesc, cgi.Bind{})
	serviceReg.Regist(&glogin2.Gmp_serviceDesc, moss.Gmp{})
	err := gmoss.RunServerByPath(serviceReg)
	//	err := gmoss.RunServer(constant.ClusterName, constant.ServerName, constant.Index, serviceReg)
	if err != nil {
		log.Fatalw("failed to register serverByPath", "err", err)
		panic(err)
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
