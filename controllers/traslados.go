package controllers

import (
	"encoding/json"
	"errors"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	trasladoshelper "github.com/udistrital/arka_mid/helpers/trasladosHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

type TrasladosController struct {
	beego.Controller
}

// URLMapping ...
func (c *TrasladosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetTraslado", c.GetTraslado)
	c.Mapping("GetElementosFuncionario", c.GetElementosFuncionario)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Post Traslado
// @Description Genera el consecutivo y hace el respectivo registro en api movimientos_arka_crud
// @Param	body		body 	models.Movimiento	true		"Informacion de las salidas y elementos asociados a cada una de ellas. Se valida solo si el id es 0""
// @Success 200 {object} models.Movimiento
// @Failure 403 body is empty
// @router / [post]
func (c *TrasladosController) Post() {

	defer errorCtrl.ErrorControlController(c.Controller, "TrasladosController")

	var v models.Movimiento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
			"err":     err,
			"status":  "400",
		})
	}

	err := trasladoshelper.Post(&v)
	if err != nil {
		panic(err)
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = v
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

	defer errorCtrl.ErrorControlController(c.Controller, "TrasladosController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un traslado válido")
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

	if respuesta, err := trasladoshelper.GetOne(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "GetTraslado - trasladoshelper.GetDetalleTraslado(id)",
			"err":     errors.New("no se obtuvo respuesta al consultar el traslado"),
			"status":  "404",
		})
	}

	c.ServeJSON()

}

// GetElementosFuncionario ...
// @Title Get Elementos
// @Description get Elementos by Tercero Origen
// @Param	tercero_id	path	int	true	"tercero_id del funcionario"
// @Success 200 {object} models.InventarioTercero
// @Failure 404 not found resource
// @router /funcionario/:tercero_id [get]
func (c *TrasladosController) GetElementosFuncionario() {

	defer errorCtrl.ErrorControlController(c.Controller, "TrasladosController")

	var (
		id         int
		inventario models.InventarioTercero
	)

	if v, err := c.GetInt(":tercero_id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un tercero válido")
		}
		panic(errorCtrl.Error(`GetElementosFuncionario - c.GetInt(":tercero_id")`, err, "400"))
	} else {
		id = v
	}

	if err := trasladoshelper.GetElementosTercero(id, &inventario); err != nil {
		panic(errorCtrl.Error("GetElementosFuncionario - trasladoshelper.GetElementosTercero(id, &inventario)", err, "404"))
	} else {
		c.Data["json"] = inventario
	}

	c.ServeJSON()

}

// GetAll ...
// @Title Get All
// @Description Consulta todos los traslados, permitiendo filtrar por las que estan pendientes de ser revisados
// @Param	user	query	string	false	"Tercero que consulta los traslados"
// @Param	confirmar	query	bool	false	"Consulta los traslados que están pendientes por ser confirmados por el tercero que consulta."
// @Param	aprobar	query	bool	false	"Consulta los traslados que están pendientes por ser aprobados por almacén."
// @Success 200 {object} []models.DetalleTrasladoLista
// @Failure 404 not found resource
// @router / [get]
func (c *TrasladosController) GetAll() {

	defer errorCtrl.ErrorControlController(c.Controller, "TrasladosController")

	var (
		terceroId string
		confirmar bool
		aprobar   bool
		traslados []*models.DetalleTrasladoLista
	)

	if v := c.GetString("user", ""); v == "" {
		panic(errorCtrl.Error(`GetAll - c.GetString("user", "")`, "Se debe indicar un usuario válido", "400"))
	} else {
		terceroId = v
	}

	if v, err := c.GetBool("confirmar", false); err != nil {
		panic(errorCtrl.Error(`GetAll - c.GetBool("confirmar", false)`, err, "400"))
	} else {
		confirmar = v
	}

	if v, err := c.GetBool("aprobar", false); err != nil {
		panic(errorCtrl.Error(`GetAll - c.GetBool("aprobar", false)`, err, "400"))
	} else {
		aprobar = v
	}

	if err := trasladoshelper.GetAll(terceroId, confirmar, aprobar, &traslados); err != nil {
		panic(err)
	}

	if traslados != nil {
		c.Data["json"] = traslados
	} else {
		c.Data["json"] = []interface{}{}
	}

	c.ServeJSON()
}

// Put ...
// @Title Put Aprobar traslado
// @Description Actualiza el estado del traslado y genera la transacción contable correspondiente.
// @Param	id		path 	int						true	"Id del traslado"
// @Success 200 {object} models.MvtoArkaMasTransaccion
// @Failure 403 body is empty
// @router /:id [put]
func (c *TrasladosController) Put() {

	defer errorCtrl.ErrorControlController(c.Controller, "TrasladosController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un traslado válido")
		}
		panic(errorCtrl.Error(`Put - c.GetInt(":id")`, err, "400"))
	} else {
		id = v
	}

	var res models.ResultadoMovimiento
	if err := trasladoshelper.AprobarTraslado(id, &res); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = res
	} else {
		panic(errorCtrl.Error("Put - trasladoshelper.AprobarTraslado(id)", err, "404"))

	}

	c.ServeJSON()
}
