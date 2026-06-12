package main

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/udistrital/arka_mid/routers"
	"github.com/udistrital/arka_mid/utils_oas/apiStatus"
	"github.com/udistrital/arka_mid/utils_oas/auditoria"
	"github.com/udistrital/arka_mid/utils_oas/customErrorv2"
	"github.com/udistrital/arka_mid/utils_oas/security"
	"github.com/udistrital/arka_mid/utils_oas/xray"
)

func main() {
	allowedOrigins := []string{"*.udistrital.edu.co"}
	if beego.BConfig.RunMode == beego.DEV {
		allowedOrigins = []string{"*"}
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"DELETE", "GET", "OPTIONS", "POST", "PUT"},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"User-Agent",
			"X-Amzn-Trace-Id"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
	}))

	apiStatus.Init()
	auditoria.InitMiddleware()
	security.SetSecurityHeaders()
	xray.Init()

	beego.ErrorController(&customErrorv2.CustomErrorController{})
	beego.Run()
}
