package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	trasladoshelper "github.com/udistrital/arka_mid/helpers/trasladosHelper"
)

type TrasladosController struct {
	beego.Controller
}

// URLMapping ...
func (c *TrasladosController) URLMapping() {
	c.Mapping("GetTr", c.GetTraslado)
}

// GetTraslado ...
// @Title Get User
// @Description get Salida by id
// @Param	id		path 	string	true		"movimientoId del traslado en el api movimientos_arka_crud"
// @Success 200 {object} models.DetalleTraslado
// @Failure 404 not found resource
// @router /:id [get]
func (c *TrasladosController) GetTraslado() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "TrasladosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	logs.Info(idStr)
	var id int
	if idConv, err := strconv.Atoi(idStr); err == nil && idConv > 0 {
		id = idConv
	} else if err != nil {
		panic(err)
	} else {
		panic(map[string]interface{}{
			"funcion": "GetDetalleTraslado",
			"err":     "El ID debe ser mayor a 0",
			"status":  "400",
		})
	}
	fmt.Println(id)

	if v, err := trasladoshelper.GetDetalleTraslado(id); err == nil {
		c.Data["json"] = v
	} else {
		panic(err)
	}
	c.ServeJSON()
}
