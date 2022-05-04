package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/models"
)

// ParametrosController operations for Parametros
type ParametrosController struct {
	beego.Controller
}

// URLMapping ...
func (c *ParametrosController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("PostAsignacionEspacioDependencia", c.PostAsignacionEspacioDependencia)
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

// PostAsignacionEspacioFisicoDependencia ...
// @Title Post Soportes
// @Description get Soportes by id
// @Param	body		body 	models.GetSedeDependencia	true		"body for content"
// @Success 201 {object} []models.AsignacionEspacioFisicoDependencia
// @Failure 404 not found resource
// @Failure 400 the request contains incorrect syntax
// @router /post_asignacion_espacio_fisico_dependencia/ [post]
func (c *ParametrosController) PostAsignacionEspacioDependencia() {

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

	var v models.GetSedeDependencia
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if res, err := actaRecibido.GetAsignacionSedeDependencia(v); err == nil {
			c.Ctx.Output.SetStatus(201)
			if res == nil {
				res = append(res, map[string]interface{}{})
			}
			c.Data["json"] = res
		} else {
			panic(err)
		}
	} else {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "PostAsignacionEspacioDependencia - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
			"err":     err,
			"status":  "500",
		})
	}
	c.ServeJSON()
}
