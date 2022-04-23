package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// EntradaController operations for Entrada
type EntradaController struct {
	beego.Controller
}

// URLMapping ...
func (c *EntradaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetEncargadoElemento", c.GetEncargadoElemento)
	c.Mapping("AnularEntrada", c.AnularEntrada)
	c.Mapping("GetMovimientos", c.GetMovimientos)
}

// Post ...
// @Title Post
// @Description Transaccion entrada. Estado de registro o aprobacion
// @Param	entradaId		 query 	string			false		"Id del movimiento que se desea aprobar"
// @Param	body			 body 	models.Entrada	false		"Detalles de la entrada. Se valida solo si el id es 0"
// @Success 201 {object} models.Entrada
// @Failure 403 body is empty
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *EntradaController) Post() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "EntradaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	var entradaId int = 0

	if v, err := c.GetInt("entradaId"); err == nil {
		entradaId = v
	}
	if entradaId > 0 {
		if respuesta, err := entradaHelper.AprobarEntrada(entradaId); err == nil && respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
		} else {
			if err == nil {
				panic(map[string]interface{}{
					"funcion": "Post - entradaHelper.AprobarEntrada(entradaId)",
					"err":     err,
					"status":  "400",
				})
			}
			panic(err)
		}
	} else {
		var v models.TransaccionEntrada
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			if respuesta, err := entradaHelper.RegistrarEntrada(&v); err == nil && respuesta != nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = respuesta
			} else if err != nil {
				panic(map[string]interface{}{
					"funcion": "Post - entradaHelper.RegistrarEntrada(v)",
					"err":     err,
					"status":  "400",
				})
			} else {
				panic(map[string]interface{}{
					"funcion": "Post - entradaHelper.RegistrarEntrada(v)",
					"err":     errors.New("no se obtuvo respuesta al registrar la entrada"),
					"status":  "404",
				})
			}
		} else {
			logs.Error(err)
			panic(map[string]interface{}{
				"funcion": "Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
				"err":     err,
				"status":  "400",
			})
		}
	}

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Detalle de entrada por Id. Retorna la transaccion contable si la entrada ya  fue aprobada
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} map[string]interface{}
// @Failure 403 :id is empty
// @router /:id [get]
func (c *EntradaController) GetOne() {

	defer errorctrl.ErrorControlController(c.Controller, "EntradaController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una entrada válida")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": `GetOne - c.GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if respuesta, err := entradaHelper.DetalleEntrada(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "GetOne - entradaHelper.DetalleEntrada(id)",
			"err":     errors.New("no se obtuvo respuesta al consultar la anetrada"),
			"status":  "404",
		})
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
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "EntradaController" + "/" + (localError["funcion"]).(string))
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
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "EntradaController" + "/" + (localError["funcion"]).(string))
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

// GetMovimientos ...
// @Title Get User
// @Description return movimientos asociados a un acta
// @Param	acta_recibido_id		path 	string	true		"The key for staticblock"
// @Success 200 {object]
// @Failure 404 not found resource
// @router /movimientos/:acta_recibido_id [get]
func (c *EntradaController) GetMovimientos() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "EntradaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":acta_recibido_id")
	actaId, _ := strconv.Atoi(idStr)
	if actaId > 0 {
		if v, err := entradaHelper.GetMovimientosByActa(actaId); err == nil {
			c.Data["json"] = v
			c.Ctx.Output.SetStatus(200)
		} else {
			logs.Error(err)
			panic(map[string]interface{}{
				"funcion": "GetMovimientosByActa",
				"err":     err,
				"status":  err["status"],
			})
		}
	} else {
		err := fmt.Errorf("{acta} no debe ser vacia")
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetMovimientosByActa",
			"err":     err,
			"status":  "404",
		})
	}
	c.ServeJSON()
}
