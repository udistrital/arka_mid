package controllers

import (
	"errors"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
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
// @Success 200 {object} models.IdentificacionTercero
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TercerosController) GetOne() {

	defer errorCtrl.ErrorControlController(c.Controller, "TercerosController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un tercero válido")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": `GetOne - GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if v, err := terceros.GetNombreTerceroById(id); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}
