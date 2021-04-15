package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/models"
)

// EntradaController operations for Entrada
type EntradaController struct {
	beego.Controller
}

// URLMapping ...
func (c *EntradaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.GetEncargadoElemento)
}

// Post ...
// @Title Post
// @Description create Entrada
// @Param	body		body 	models.Entrada	true		"body for Entrada content"
// @Success 201 {object} models.Entrada
// @Failure 403 body is empty
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *EntradaController) Post() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "EntradaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	var v models.Movimiento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if respuesta, err := entradaHelper.AddEntrada(v); err == nil && respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
		} else {
			if err == nil {
				panic(map[string]interface{}{
					"funcion": "Post - entradaHelper.AddEntrada(v)",
					"err":     err,
					"status":  "400",
				})
			}
			panic(err)
		}
	} else {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
			"err":     err,
			"status":  "400",
		})
	}
	c.ServeJSON()
}

// GetEntrada ...
// @Title Get User
// @Description get Entrada by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaEntrada
// @Failure 404 not found resource
// @router /:id [get]
func (c *EntradaController) GetEntrada() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := entradaHelper.GetEntrada(id)
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetEntradas ...
// @Title Get User
// @Description get Entradas
// @Success 200 {object} models.ConsultaEntrada
// @Failure 404 not found resource
// @router / [get]
func (c *EntradaController) GetEntradas() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "EntradaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	if v, err := entradaHelper.GetEntradas(); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetEncargadoElemento ...
// @Title Get User
// @Description get Entradas
// @Param	placa		path 	string	true		"The key for staticblock"
// @Success 200  {"funcionario": "string"}
// @Failure 404 not found resource
// @router /encargado/:placa [get]
func (c *EntradaController) GetEncargadoElemento() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "EntradaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	placa := c.Ctx.Input.Param(":placa")
	if placa == "" {
		err := fmt.Errorf("{placa} no debe ser vacia")
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetEncargadoElemento",
			"err":     err,
			"status":  "400",
		})
	}

	if funcionario, err := entradaHelper.GetEncargadoElemento(placa); err == nil {
		c.Data["json"] = funcionario
		c.Ctx.Output.SetStatus(200)
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// AnularEntrada ...
// @Title Get User
// @Description anular Entrada by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaEntrada
// @Failure 404 not found resource
// @router /anular/:id [get]
func (c *EntradaController) AnularEntrada() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "EntradaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if idStr != "" {
		if v, err := entradaHelper.AnularEntrada(id); err == nil {
			c.Data["json"] = v
			c.Ctx.Output.SetStatus(200)
		} else {
			logs.Error(err)
			panic(map[string]interface{}{
				"funcion": "AnularEntrada",
				"err":     err,
				"status":  err["status"],
			})
		}
	} else {
		panic(map[string]interface{}{
			"funcion": "AnularEntrada",
			"err":     "La entrada no puede ser nula",
			"status":  "404",
		})
	}
	c.ServeJSON()
}
