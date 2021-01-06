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
// @Failure 404 "not found resource"
// @Failure 500 "Unknown API Error"
// @Failure 502 "External API Error"
// @router /solicitud/:id [get]
func (c *BajaController) GetSolicitud() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "BajaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	fmt.Println(idStr)

	if v, err := bajasHelper.TraerDetalle(id); err == nil {
		c.Data["json"] = v
	} else {
		panic(err)
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
