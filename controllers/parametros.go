package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/models"
	"encoding/json"

)

// ParametrosController operations for Parametros
type ParametrosController struct {
	beego.Controller
}

// URLMapping ...
func (c *ParametrosController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Post", c.Post)
	c.Mapping("PostAsignacionEspacioDependencia", c.PostAsignacionEspacioDependencia)
}

// GetAll ...
// @Title Get All
// @Description get ActaRecibido
// @Success 200 {object} models.ActaRecibido
// @Failure 404 not found resource
// @router / [get]
func (c *ParametrosController) GetAll() {
	l, err := actaRecibido.GetAllParametrosSoporte()
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

// PostAsignacionEspacioFisicoDependencia ...
// @Title Post Soportes
// @Description get Soportes by id
// @Param	body		body 	{}	true		"body for content"
// @Success 201 {object} []models.AsignacionEspacioFisicoDependencia
// @Failure 404 not found resource
// @Failure 400 the request contains incorrect syntax
// @router /post_asignacion_espacio_fisico_dependencia/ [post]
func (c *ParametrosController) PostAsignacionEspacioDependencia() {

	var v models.GetSedeDependencia
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if res, err := actaRecibido.GetAsignacionSedeDependencia(v); err == nil {
			c.Ctx.Output.SetStatus(201)
			if res == nil {
				res = append(res, map[string]interface{}{})
			}
			c.Data["json"] = res
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
