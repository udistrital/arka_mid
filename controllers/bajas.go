package controllers

import (
	// "encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/bajasHelper"
	// "github.com/udistrital/arka_mid/models"
)

// BajaController operations for Salida
type BajaController struct {
	beego.Controller
}

// URLMapping ...
func (c *BajaController) URLMapping() {
	c.Mapping("Get", c.GetElemento)
}

// GetElemento ...
// @Title Get User
// @Description get Salida by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Salida
// @Failure 404 not found resource
// @router /elemento/:id [get]
func (c *BajaController) GetElemento() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	fmt.Println(idStr)
	v, err := bajasHelper.TraerDatosElemento(id)
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

// Getsolicitud...
// @Title Get User
// @Description get Salida by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Salida
// @Failure 404 not found resource
// @router /solicitud/:id [get]
func (c *BajaController) GetSolicitud() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	fmt.Println(idStr)
	v, err := bajasHelper.TraerDetalle(id)
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

// GetAll ...
// @Title Get All
// @Description get Baja
// @Success 200 {object} models.Baja
// @Failure 404 not found resource
// @router /solicitud/ [get]
func (c *BajaController) GetAll() {

	fmt.Println("hola")
	l, err := bajasHelper.GetAllSolicitudes()
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
