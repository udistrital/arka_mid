package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
)

// ParametrosController operations for Parametros
type ParametrosController struct {
	beego.Controller
}

// URLMapping ...
func (c *ParametrosController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}

// GetAll ...
// @Title Get All
// @Description get ActaRecibido
// @Success 200 {object} models.ActaRecibido
// @Failure 404 not found resource
// @router / [get]
func (c *ParametrosController) GetAll() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ParametrosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	if l, err := actaRecibido.GetAllParametrosSoporte(); err != nil {
		panic(err)
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}
