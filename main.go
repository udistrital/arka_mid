package main

import (
	_ "github.com/udistrital/arka_mid/routers"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/udistrital/arka_mid/utils_oas/apiStatus"
	"github.com/udistrital/arka_mid/utils_oas/auditoria"
	"github.com/udistrital/arka_mid/utils_oas/customErrorv2"
	"github.com/udistrital/arka_mid/utils_oas/security"
)

func main() {
	allowedOrigins := []string{"*.udistrital.edu.co"}
	if web.BConfig.RunMode == web.DEV {
		allowedOrigins = []string{"*"}
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"accept", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
	}))

	apiStatus.Init()
	auditoria.InitMiddleware()
	security.SetSecurityHeaders()

	web.ErrorController(&customErrorv2.CustomErrorController{})
	web.Run()
}
