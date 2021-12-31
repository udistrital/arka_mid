package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/depreciacionHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// DepreciacionController operations for Depreciacion
type DepreciacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *DepreciacionController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// GetCorte ...
// @Title GetCorte
// @Description Actualiza el estado de las solicitudes una vez se registra la revision del comite de almacen
// @Param	body			 body 	models.InfoDepreciacion	false	"Informacion de la liquidacion de depreciacion"
// @Success 200 {object} []int
// @Failure 404 not found resource
// @router / [post]
func (c *DepreciacionController) Post() {

	defer errorctrl.ErrorControlController(c.Controller, "DepreciacionController")

	var v *models.InfoDepreciacion
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		panic(errorctrl.Error("Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)", err, "400"))
	} else {
		if v, err := depreciacionHelper.GenerarTrDepreciacion(v); err != nil {
			logs.Error(err)
			c.Data["system"] = err
			c.Abort("404")
		} else {
			c.Data["json"] = v
		}
	}
	c.ServeJSON()
}
