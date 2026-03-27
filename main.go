package main

import (
	_ "github.com/udistrital/arka_mid/routers"

	"github.com/udistrital/arka_mid/utils_oas/apiStatus"
	"github.com/udistrital/arka_mid/utils_oas/auditoria"
	"github.com/udistrital/arka_mid/utils_oas/customErrorv2"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
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
	web.InsertFilter("*", web.BeforeExec, SecurityHeaders)
	apiStatus.Init()
	auditoria.InitMiddleware()
	web.Run()
}

func SecurityHeaders(ctx *context.Context) {
	ctx.Output.Header("Clear-Site-Data", "'cache', 'cookies', 'storage', 'executionContexts'")
	ctx.Output.Header("Cross-Origin-Embedder-Policy", "require-corp")
	ctx.Output.Header("Cross-Origin-Opener-Policy", "same-origin")
	ctx.Output.Header("Cross-Origin-Resource-Policy", "same-origin")
	ctx.Output.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	ctx.Output.Header("Referrer-Policy", "no-referrer")
	ctx.Output.Header("Server", "")
	ctx.Output.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	ctx.Output.Header("X-Content-Type-Options", "nosniff")
	ctx.Output.Header("X-Frame-Options", "DENY")
	ctx.Output.Header("X-Permitted-Cross-Domain-Policies", "none")
}
