package controllers

import (
	"errors"
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
	c.Mapping("Get", c.GetTraslado)
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
				c.Abort("500")
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	if trasladoId, err := strconv.Atoi(idStr); err != nil || trasladoId == 0 {
		if err == nil {
			err = errors.New("El id del movimiento no puede ser cero")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetTraslado - strconv.Atoi(idStr)",
			"err":     err,
			"status":  "400",
		})
	} else {
		if v, err := trasladoshelper.GetDetalleTraslado(trasladoId); err == nil {
			c.Data["json"] = v
		} else {
			panic(err)
		}
		c.ServeJSON()
	}

}
