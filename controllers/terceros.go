package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/terceros"
)

// TercerosController operations for Terceros
type TercerosController struct {
	beego.Controller
}

// URLMapping ...
func (c *TercerosController) URLMapping() {
	c.Mapping("GetOne", c.GetOne)
}

// GetOne ...
// @Title GetOne
// @Description get Terceros by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Terceros
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TercerosController) GetOne() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "TercerosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500")
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")

	if v, err := terceros.GetNombreTerceroById(idStr); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}
