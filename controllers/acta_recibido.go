package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	//"github.com/udistrital/acta_recibido_crud/models"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/actaRecibidoHelper"
)

// ActaRecibidoController operations for ActaRecibido
type ActaRecibidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ActaRecibidoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetElementosActa", c.GetElementosActa)
	c.Mapping("GetElementosConsumo", c.GetAllElementosConsumo)
}

// Post ...
// @Title Create
// @Description create Acta_recibido
// @Param	archivo	formData  file	true	"body for Acta_recibido content"
// @Success 201 {}
// @Failure 403 body is empty
// @router / [post]
func (c *ActaRecibidoController) Post() {
	fmt.Println(c.GetFile("archivo"))
	if multipartFile, _, err := c.GetFile("archivo"); err == nil {
		if Archivo, err := actaRecibido.DecodeXlsx2Json(multipartFile); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = Archivo
		} else {
			c.Data["system"] = err
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get ActaRecibido
// @Success 200 {object} models.ActaRecibido
// @Failure 404 not found resource
// @router / [get]
func (c *ActaRecibidoController) GetAll() {

	fmt.Println("hola")
	l, err := actaRecibido.GetAllParametrosActa()
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// GetActasByTipo ...
// @Title GetActasRecibidoTipo
// @Description Devuelve las todas las actas de recibido
// @Param	id		path 	string	true		"id del acta"
// @Success 200 {object} models.Acta_recibido
// @Failure 403
// @router /get_actas_recibido_tipo/:tipo [get]
func (c *ActaRecibidoController) GetActasByTipo() {
	tipoStr := c.Ctx.Input.Param(":tipo")
	tipo, _ := strconv.Atoi(tipoStr)
	v, err := actaRecibidoHelper.GetActasRecibidoTipo(tipo)
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetElementosActa ...
// @Title Get Elementos
// @Description get Elementos by id
// @Param	id		path 	string	true		"id del acta"
// @Success 200 {object} models.Elemento
// @Failure 404 not found resource
// @router /get_elementos_acta/:id [get]
func (c *ActaRecibidoController) GetElementosActa() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := actaRecibidoHelper.GetElementos(id)
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

// GetSoportesActa ...
// @Title Get Soportes
// @Description get Soportes by id
// @Param	body	body 	models.Entrada	true
// @Success 200 {object} []models.AsignacionEspacioFisicoDependencia
// @Failure 404 not found resource
// @router /get_soportes_acta/:id [get]
func (c *ActaRecibidoController) GetSoportesActa() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := actaRecibidoHelper.GetSoportes(id)
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

// GetAllElementosConsumo ...
// @Title GetAllElementosConsumo
// @Description Trae todos los elementos de consumo
// @Success 200 {object} models.Elemento
// @Failure 404 not found resource
// @router /elementosconsumo/ [get]
func (c *ActaRecibidoController) GetAllElementosConsumo() {

	fmt.Println("hola hola")
	l, err := actaRecibido.GetAllElementosConsumo()
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}
