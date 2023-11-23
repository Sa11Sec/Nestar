package main

import (
	"Nestar/config"
	"Nestar/middlewares"
	"Nestar/routes"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	//设置为正式版
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	Port := strconv.FormatUint(uint64(config.ServerSetting.Port), 10)
	// TLS 中间件
	r.Use(middlewares.TlSHandler(Port))
	// 加载模板
	r.LoadHTMLGlob("templates/**/*")
	// 设置静态文件
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	// 配置路由
	routes.SetupRoutes(r)

	// 开启服务
	err := r.RunTLS(":"+Port, "config/cert.pem", "config/key.pem")
	if err != nil {
		println(err)
	}
}
