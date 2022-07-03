package controllers

import (
	"encoding/json"
	"errors"

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
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Put", c.Put)
}

// GetCorte ...
// @Title GetCorte
// @Description Actualiza el estado de las solicitudes una vez se registra la revision del comite de almacen
// @Param	body			 body 	models.InfoDepreciacion	false	"Informacion de la liquidacion de depreciacion"
// @Success 200 {object} models.ResultadoMovimiento
// @Failure 404 not found resource
// @router / [post]
func (c *DepreciacionController) Post() {

	defer errorctrl.ErrorControlController(c.Controller, "DepreciacionController")

	var v *models.InfoDepreciacion
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		panic(errorctrl.Error("Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)", err, "400"))
	} else {
		var resultado models.ResultadoMovimiento
		if err := depreciacionHelper.GenerarCierre(v, &resultado); err != nil {
			logs.Error(err)
			c.Data["system"] = err
			c.Abort("404")
		} else {
			c.Data["json"] = resultado
		}
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get Info Depreciacion
// @Description get Depreciacion by id
// @Param	id	path	int	true	"movimientoId de la depreciacion en el api movimientos_arka_crud"
// @Success 200 {object} models.ResultadoMovimiento
// @Failure 404 not found resource
// @router /:id [get]
func (c *DepreciacionController) GetOne() {

	defer errorctrl.ErrorControlController(c.Controller, "DepreciacionController - Unhandled Error!")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una depreciación válida")
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

	var detalle models.ResultadoMovimiento
	if err := depreciacionHelper.GetCierre(id, &detalle); err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = detalle
	}

	c.ServeJSON()

}

// Put ...
// @Title Put
// @Description update the ElementosMovimiento
// @Param	id		path 	string	true		"The id you want to update"
// @Success 200 {object} models.ResultadoMovimiento
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *DepreciacionController) Put() {
	defer errorctrl.ErrorControlController(c.Controller, "DepreciacionController - Unhandled Error!")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una depreciación válida")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": `Put - GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	var detalle models.ResultadoMovimiento
	if err := depreciacionHelper.AprobarDepreciacion(id, &detalle); err == nil {
		c.Data["json"] = detalle
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "Put - depreciacionHelper.AprobarDepreciacion(id)",
			"err":     errors.New("no se obtuvo respuesta al consultar la depreciación"),
			"status":  "404",
		})
	}

	c.ServeJSON()
}
