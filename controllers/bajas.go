package controllers

import (
	"encoding/json"
	"errors"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/helpers/bajasHelper"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/inventarioHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// BajaController
type BajaController struct {
	beego.Controller
}

// URLMapping ...
func (c *BajaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("GetSolicitud", c.GetSolicitud)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetDetalleElemento", c.GetDetalleElemento)
	c.Mapping("PutRevision", c.PutRevision)
}

// Post ...
// @Title Post
// @Description Registrar Baja. Crea el registro del soporte y crea el consecutivo
// @Param	body	body 	models.TrSoporteMovimiento	false	"Informacion de la baja"
// @Success 201	{object}	models.Movimiento
// @Failure 403 body is empty
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *BajaController) Post() {

	defer errorCtrl.ErrorControlController(c.Controller, "BajaController")

	var v *models.TrSoporteMovimiento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		panic(errorCtrl.Error("Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)", err, "400"))
	} else {
		if respuesta, err := bajasHelper.Post(v); err == nil && respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
			c.ServeJSON()
		} else {
			if err != nil {
				panic(err)
			}

			panic(errorCtrl.Error("Post", "No se obtuvo respuesta al registrar la baja", "404"))
		}

	}

}

// Put ...
// @Title Put
// @Description Update Baja. Actualiza los detalles de la baja y el documento
// @Param	id	path	int	true	"movimientoId de la baja en el api movimientos_arka_crud"
// @Success 201 {object} models.TrSoporteMovimiento
// @Failure 403 body is empty
// @Failure 400 the request contains incorrect syntax
// @router /:id [put]
func (c *BajaController) Put() {

	defer errorCtrl.ErrorControlController(c.Controller, "BajaController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una baja válida")
		}
		panic(errorCtrl.Error("Put - c.GetInt(\":id\")", err, "400"))
	} else {
		id = v
	}

	var v *models.TrSoporteMovimiento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		panic(errorCtrl.Error("Put - json.Unmarshal(c.Ctx.Input.RequestBody, &v)", err, "400"))
	} else {
		if respuesta, err := bajasHelper.Put(v, id); err == nil && respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
			c.ServeJSON()
		} else {
			if err != nil {
				panic(err)
			}

			panic(errorCtrl.Error("Put", "No se obtuvo respuesta al actualizar la baja", "404"))
		}
	}

}

// Getsolicitud...
// @Title Get User
// @Description consulta detalle de Baja
// @Param	id	path 	string	true	"Id de la baja en el api movimientos_arka_crud"
// @Success 200 {object} models.TrBaja
// @Failure 404 "not found resource"
// @Failure 500 "Unknown API Error"
// @Failure 502 "External API Error"
// @router /:id [get]
func (c *BajaController) GetSolicitud() {

	defer errorCtrl.ErrorControlController(c.Controller, "BajaController")

	var (
		id   int
		baja models.TrBaja
	)

	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una baja válida")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetSolicitud - GetInt(\":id\")",
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if err := bajasHelper.GetOne(id, &baja); err == nil {
		c.Data["json"] = baja
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description Consulta todas las bajas y permite filtrar por las que estan para revision de almacen o comite
// @Param	user	query	string	true	"Tercero que consulta las bajas"
// @Param	revComite	query 	bool	false	"Indica si se traen las bajas en espera de comite. Tiene prioridad sobre la revision de almacen"
// @Param	revAlmacen	query 	bool	false	"Indica si se traen las bajas pendientes por revisar"
// @Success 200 {object} []models.DetalleBaja
// @Failure 404 not found resource
// @router / [get]
func (c *BajaController) GetAll() {

	defer errorCtrl.ErrorControlController(c.Controller, "BajaController")

	var (
		revComite  bool
		revAlmacen bool
		terceroId  string
	)

	if v := c.GetString("user", ""); v == "" {
		panic(errorCtrl.Error(`GetAll - c.GetString("user", "")`, "Se debe indicar un usuario válido", "400"))
	} else {
		terceroId = v
	}

	if v, err := c.GetBool("revComite"); err != nil {
		panic(errorCtrl.Error(`GetAll - c.GetBool("revComite")`, err, "400"))
	} else {
		revComite = v
	}

	if v, err := c.GetBool("revAlmacen"); err != nil {
		panic(errorCtrl.Error(`GetAll - c.GetBool("revAlmacen")`, err, "400"))
	} else {
		revAlmacen = v
	}

	var bajas = make([]models.DetalleBaja, 0)
	if err := bajasHelper.GetAll(terceroId, revComite, revAlmacen, &bajas); err != nil {
		panic(err)
	} else {
		c.Data["json"] = bajas
	}

	c.ServeJSON()
}

// GetDetalleElemento ...
// @Title GetElemento
// @Description Get Info relacionada con el historial de un determinado elemento
// @Param	id	path	int	true	"id del elemento en el api acta_recibido_crud"
// @Success 200 {object} models.DetalleElementoBaja
// @Failure 404 not found resource
// @router /elemento/:id [get]
func (c *BajaController) GetDetalleElemento() {

	defer errorCtrl.ErrorControlController(c.Controller, "BajaController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un elemento válido")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetDetalleElemento - GetInt(\":id\")",
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	var elemento models.DetalleElementoBaja
	err := inventarioHelper.GetDetalleElemento(id, &elemento)
	if err != nil {
		panic(err)
	}

	c.Data["json"] = elemento
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description Realiza la transacciones contables y actualiza los movimientos
// @Param	body	body 	models.TrRevisionBaja	true	"Informacion de la revision"
// @Success 200 {object} []int
// @Failure 404 not found resource
// @router /aprobar [put]
func (c *BajaController) PutRevision() {

	defer errorCtrl.ErrorControlController(c.Controller, "BajaController")

	var trBaja *models.TrRevisionBaja
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &trBaja); err != nil {
		panic(errorCtrl.Error("PutRevision - json.Unmarshal(c.Ctx.Input.RequestBody, &trBaja)", err, "400"))
	}

	if !trBaja.Aprobacion {
		if ids, err := movimientosArka.PutRevision(trBaja); err != nil {
			panic(errorCtrl.Error("PutRevision - movimientosArkaHelper.PutRevision(trBaja)", err, "404"))
		} else {
			c.Data["json"] = ids
		}
	} else {
		var response models.ResultadoMovimiento
		if err := bajasHelper.AprobarBajas(trBaja, &response); err != nil {
			panic(errorCtrl.Error("PutRevision - bajasHelper.AprobarBajas(trBaja)", err, "404"))
		} else {
			c.Data["json"] = response
		}
	}

	c.ServeJSON()
}
