package main

import (
	"blog-server/config"
	_ "blog-server/dao"
	"blog-server/pkg/middleware/logger"
	"blog-server/router"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	port, mode string
)

func init() {
	flag.StringVar(&port, "port", "3000", "server listening on, default 3000")
	flag.StringVar(&mode, "mode", "debug", "server running mode, default debug mode")
}

func main() {
	port := config.GetServerCfg().Port
	flag.Parse()
	gin.SetMode(mode)
	r := gin.Default()
	router.Init(r)
	r.StaticFS("/images", http.Dir("./static/images"))
	r.Use(logger.LoggerToFile())
	err := r.Run(port)
	if err != nil {
		log.Fatalf("Start server: %+v", err)
	}
}
