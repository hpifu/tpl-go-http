package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/hconf"
	"github.com/hpifu/go-kit/hdef"
	"github.com/hpifu/go-kit/henv"
	"github.com/hpifu/go-kit/hflag"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/hpifu/go-kit/hrule"
	"github.com/hpifu/go-kit/logger"
	"github.com/hpifu/tpl-go-http/internal/service"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// AppVersion name
var AppVersion = "unknown"

type Options struct {
	Service struct {
		Port         string   `hflag:"usage: service port" hdef:":1234"`
		AllowOrigins []string `hflag:"usage: allow origins" hdef:"127.0.0.1"`
		CookieSecure bool     `hflag:"usage: http or https"`
		CookieDomain string   `hflag:"usage: cookie domain"`
	}
	Es struct {
		Uri string `hflag:"usage: elasticsearch address"`
	}
	Logger struct {
		Info   logger.Options
		Warn   logger.Options
		Access logger.Options
	}
}

func main() {
	version := hflag.Bool("v", false, "print current version")
	configfile := hflag.String("c", "configs/echo.json", "config file path")
	if err := hflag.Bind(&Options{}); err != nil {
		panic(err)
	}
	if err := hflag.Parse(); err != nil {
		panic(err)
	}
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	// load config
	options := &Options{}
	if err := hdef.SetDefault(options); err != nil {
		panic(err)
	}
	config, err := hconf.New("json", "local", *configfile)
	if err != nil {
		panic(err)
	}
	if err := config.Unmarshal(options); err != nil {
		panic(err)
	}
	if err := henv.NewHEnv("ECHO").Unmarshal(options); err != nil {
		panic(err)
	}
	if err := hflag.Unmarshal(options); err != nil {
		panic(err)
	}
	if err := hrule.Evaluate(options); err != nil {
		panic(err)
	}

	// init logger
	logs, err := logger.NewLoggerGroup([]*logger.Options{
		&options.Logger.Info, &options.Logger.Warn, &options.Logger.Access,
	})
	if err != nil {
		panic(err)
	}
	infoLog := logs[0]
	warnLog := logs[1]
	accessLog := logs[2]
	client, err := elastic.NewClient(
		elastic.SetURL(options.Es.Uri),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	hook, err := elogrus.NewAsyncElasticHook(client, "go-tech", logrus.InfoLevel, "go-tech-log")
	if err != nil {
		panic(err)
	}
	accessLog.Hooks.Add(hook)

	// init services
	svc := service.NewService(options.Service.CookieSecure, options.Service.CookieDomain)
	svc.SetLogger(infoLog, warnLog, accessLog)

	// init gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     options.Service.AllowOrigins,
		AllowMethods:     []string{"PUT", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// set handler
	d := hhttp.NewGinHttpDecorator(infoLog, warnLog, accessLog)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	r.GET("/echo", d.Decorate(svc.Echo))

	infoLog.Infof("%v init success, port [%v]", os.Args[0], options.Service.Port)

	// run server
	server := &http.Server{
		Addr:    options.Service.Port,
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// graceful quit
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	infoLog.Infof("%v shutdown ...", os.Args[0])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		warnLog.Errorf("%v shutdown fail or timeout", os.Args[0])
		return
	}
	for _, log := range logs {
		_ = log.Out.(*rotatelogs.RotateLogs).Close()
	}
}
