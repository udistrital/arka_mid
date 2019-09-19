package main

import (
	_ "github.com/udistrital/arka_mid/routers"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	"github.com/udistrital/utils_oas/customerror"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	// "github.com/udistrital/utils_oas/responseformat"
)

func main() {
	// beego.BConfig.RecoverFunc = responseformat.GlobalResponseHandler
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	beego.ErrorController(&customerror.CustomErrorController{})
	apistatus.Init()
	beego.Run()
}
