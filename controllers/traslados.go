package controllers

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	trasladoshelper "github.com/udistrital/arka_mid/helpers/trasladosHelper"
	"github.com/udistrital/arka_mid/models"
)

type TrasladosController struct {
	beego.Controller
}

// URLMapping ...
func (c *TrasladosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.GetTraslado)
}

// Post ...
// @Title Post Traslado
// @Description Genera el consecutivo y hace el respectivo registro en api movimientos_arka_crud
// @Param	body		body 	models.Movimiento	true		"Informacion de las salidas y elementos asociados a cada una de ellas. Se valida solo si el id es 0""
// @Success 200 {object} models.Movimiento
// @Failure 403 body is empty
// @router / [post]
func (c *TrasladosController) Post() {

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

	var v models.Movimiento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if respuesta, err := trasladoshelper.RegistrarTraslado(&v); err == nil && respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
		} else {
			status := "400"
			if err == nil {
				err = map[string]interface{}{
					"err": errors.New("No se obtuvo respuesta al registrar el traslado"),
				}
				status = "404"
			}
			logs.Error(err)
			panic(map[string]interface{}{
				"funcion": "trasladoshelper.RegistrarTraslado(&v)",
				"err":     err,
				"status":  status,
			})
		}
	} else {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
			"err":     err,
			"status":  "400",
		})
	}

	c.ServeJSON()
}

// GetTraslado ...
// @Title Get User
// @Description get Traslado by id
// @Param	id		path 	string	true		"movimientoId del traslado en el api movimientos_arka_crud"
// @Success 200 {object} models.TrTraslado
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
