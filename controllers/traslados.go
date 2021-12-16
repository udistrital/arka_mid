package controllers

import (
	"encoding/json"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	trasladoshelper "github.com/udistrital/arka_mid/helpers/trasladosHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

type TrasladosController struct {
	beego.Controller
}

// URLMapping ...
func (c *TrasladosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.GetTraslado)
	c.Mapping("GetElementosFuncionario", c.GetElementosFuncionario)
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
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
			"err":     err,
			"status":  "400",
		})
	}

	if respuesta, err := trasladoshelper.RegistrarTraslado(&v); err == nil && respuesta != nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "trasladoshelper.RegistrarTraslado(&v)",
			"err":     errors.New("No se obtuvo respuesta al registrar el traslado"),
			"status":  "404",
		})
	}

	c.ServeJSON()
}

// GetTraslado ...
// @Title Get User
// @Description get Traslado by id
// @Param	id	path	int	true	"movimientoId del traslado en el api movimientos_arka_crud"
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

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("Se debe especificar un traslado válido")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetTraslado - GetInt(\":id\")",
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if respuesta, err := trasladoshelper.GetDetalleTraslado(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "GetTraslado - trasladoshelper.GetDetalleTraslado(id)",
			"err":     errors.New("No se obtuvo respuesta al consultar el traslado"),
			"status":  "404",
		})
	}

	c.ServeJSON()

}

// GetElementosFuncionario ...
// @Title Get Elementos
// @Description get Elementos by Tercero Origen
// @Param	funcionarioId	path	int	true	"tercero_id del funcionario"
// @Success 200 {object} models.DetalleElementoPlaca
// @Failure 404 not found resource
// @router /funcionario/:funcionarioId [get]
func (c *TrasladosController) GetElementosFuncionario() {

	defer errorctrl.ErrorControlController(c.Controller, "TrasladosController - Unhandled Error!")
	var id int
	if v, err := c.GetInt(":funcionarioId"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("Se debe especificar un tercero válido")
		}
		panic(errorctrl.Error("GetElementosFuncionario - c.GetInt(\":funcionarioId\")", err, "400"))
	} else {
		id = v
	}

	if respuesta, err := trasladoshelper.GetElementosFuncionario(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}
		panic(errorctrl.Error("GetElementosFuncionario - trasladoshelper.GetElementosFuncionario(id)", err, "404"))
	}

	c.ServeJSON()

}

// GetAll ...
// @Title Get All
// @Description Consulta todos los traslados, permitiendo filtrar por las que estan pendientes de ser revisados
// @Param	tramiteOnly	query 	bool	false	"Indica si se requieren los traslados en estado en tramite"
// @Success 200 {object} []models.DetalleTrasladoLista
// @Failure 404 not found resource
// @router / [get]
func (c *TrasladosController) GetAll() {

	defer errorctrl.ErrorControlController(c.Controller, "TrasladosController")

	var tramiteOnly bool
	if v, err := c.GetBool("tramiteOnly", false); err != nil {
		panic(errorctrl.Error("GetAll - c.GetBool(\"tramiteOnly\", false)", err, "400"))
	} else {
		tramiteOnly = v
	}

	if v, err := trasladoshelper.GetAllTraslados(tramiteOnly); err == nil {
		if v != nil {
			c.Data["json"] = v
		} else {
			c.Data["json"] = []interface{}{}
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
