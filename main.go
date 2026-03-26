package main

import (
	_ "github.com/udistrital/arka_mid/routers"

	"github.com/udistrital/arka_mid/utils_oas/apiStatus"
	"github.com/udistrital/arka_mid/utils_oas/auditoria"
	"github.com/udistrital/arka_mid/utils_oas/customErrorv2"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {

	AllowedOrigins := []string{"*.udistrital.edu.co"}
	if web.BConfig.RunMode == "dev" {
		AllowedOrigins = []string{"*"}
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: AllowedOrigins,
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
	}))
	web.ErrorController(&customErrorv2.CustomErrorController{})
	apiStatus.Init()
	auditoria.InitMiddleware()
	web.Run()
}
