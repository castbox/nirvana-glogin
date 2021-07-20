package main

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gmoss/v2"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"glogin/config"
	"glogin/constant"
	"glogin/db"
	"glogin/handler/cgi"
	"glogin/internal/configure"
	glogin2 "glogin/pbs/glogin"
	"net/http"
	"runtime"
)

func main() {
	// 开发模式设置
	log.ResetToDevelopment()
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
	log.Infow("http server started")
	// 循环
	select {}
}

func initDB() {
	db.InitMongo()
}

func startHttp() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	pprof.Register(router)
	router.POST("/:service/:action", cgi.ServiceHandler)
	router.POST("/cfg", cgi.CfgHandler)
	router.POST("/token", cgi.TokenHandler)
	router.GET("/metrics", prometheusHandler())
	router.GET("/heart", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"err_code": constant.ErrCodeOk}) })
	//strWebPort := strconv.Itoa(int(config.GetAll().WebPort))
	if err := router.Run(config.GetAll().WebPort); err != nil {
		log.Fatalw("failed to run http server", "err", err)
		panic(err)
	}
}

func initConfig() {
	config.Init()
	configure.WatchPubCfg()
}

func startEngine() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	serviceReg := gmoss.NewServiceRegister()
	serviceReg.Regist(&glogin2.Login_serviceDesc, cgi.Login{})
	serviceReg.Regist(&glogin2.Bind_serviceDesc, cgi.Bind{})
	loop, err := gmoss.RegServer(constant.ClusterName, constant.ServerName, constant.Index, serviceReg)
	//loop, err := gmoss.RegServerByPath(serviceReg)
	if err != nil {
		log.Fatalw("failed to register serverByPath", "err", err)
		panic(err)
	}
	go loop()
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
