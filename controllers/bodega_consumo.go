package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/helpers/bodegaConsumoHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// BodegaConsumoController operations for Bodega-Consumo
type BodegaConsumoController struct {
	beego.Controller
}

// URLMapping ...
func (c *BodegaConsumoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOneSolicitud", c.GetOneSolicitud)
	c.Mapping("GetAllSolicitud", c.GetAllSolicitud)
	c.Mapping("GetElementos", c.GetElementos)
	c.Mapping("GetAperturasKardex", c.GetAperturasKardex)
	c.Mapping("GetAllExistencias", c.GetAllExistencias)
}

// Post ...
// @Title Post
// @Description Genera el movimiento correspondiente a una solicitud de bodega de consumo. Asigna el consecutivo correspondiente.
// @Param	body		body 	models.FormatoSolicitudBodega	true		"Detalle de la solicitud a la bodega de consumo."
// @Success 200 {object} models.Movimiento
// @Failure 404 not found resource
// @router /solicitud/ [post]
func (c *BodegaConsumoController) Post() {

	defer errorctrl.ErrorControlController(c.Controller, "BodegaConsumoController")

	var (
		v         models.FormatoSolicitudBodega
		solicitud models.Movimiento
	)

	if err := utilsHelper.Unmarshal(string(c.Ctx.Input.RequestBody), &v); err != nil {
		panic(errorctrl.Error("Post - utilsHelper.Unmarshal(string(c.Ctx.Input.RequestBody), &v)", err, "400"))
	}

	if err := bodegaConsumoHelper.PostSolicitud(&v, &solicitud); err != nil {
		panic(err)
	}

	c.Data["json"] = solicitud
	c.ServeJSON()

}

// GetOneSolicitud ...
// @Title GetOneSolicitud
// @Description get Bodega-Consumo by id
// @Param	id		path 	uint	true		"MovimientoId from Movimientos Arka CRUD"
// @Success 200 {object} models.BodegaConsumoSolicitud
// @Failure 400 "Wrong parameter (ID MUST be > 0)"
// @Failure 404 "Not found"
// @Failure 500 "Internal Error"
// @Failure 502 "Error with external API"
// @router /solicitud/:id [get]
func (c *BodegaConsumoController) GetOneSolicitud() {

	defer errorctrl.ErrorControlController(c.Controller, "BodegaConsumoController")

	idStr := c.Ctx.Input.Param(":id")
	var id int
	if idConv, err := strconv.Atoi(idStr); err == nil && idConv > 0 {
		id = idConv
	} else {
		if err == nil {
			err = fmt.Errorf("ID MUST be an integer > 0")
		}
		panic(map[string]interface{}{
			"funcion": "GetOneSolicitud",
			"err":     err,
			"status":  "400",
		})
	}

	if v, err := bodegaConsumoHelper.GetSolicitudById(id); err == nil {
		c.Data["json"] = v
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GetAllSolicitud ...
// @Title GetAllSolicitud
// @Description get Lista solicitudes de elementos de la bodega de consumo.
// @Param	user			query	string	true	"Tercero que consulta las solicitudes a bodega de consumo"
// @Param	tramite_only	query	bool	false	"Retornar solo las solicitudes en estado pendiente"
// @Success 200 {object} []models.DetalleSolicitudBodega
// @router /solicitud/ [get]
func (c *BodegaConsumoController) GetAllSolicitud() {

	defer errorctrl.ErrorControlController(c.Controller, "BodegaConsumoController")

	var (
		tramiteOnly bool
		err         error
		user        string
	)
	if tramiteOnly, err = c.GetBool("tramite_only", false); err != nil {
		panic(err)
	}

	if v := c.GetString("user", ""); v == "" {
		panic(errorctrl.Error(`GetAllSolicitud - c.GetString("user", "")`, "Se debe indicar un usuario válido", "400"))
	} else {
		user = v
	}

	solicitudes := make([]models.DetalleSolicitudBodega, 0)
	if err := bodegaConsumoHelper.GetAllSolicitudes(user, tramiteOnly, &solicitudes); err != nil {
		panic(err)
	}

	if solicitudes != nil {
		c.Data["json"] = solicitudes
	} else {
		c.Data["json"] = []interface{}{}
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Bodega-Consumo
// @Success 200 {object} []models.ElementoSinAsignar
// @Failure 403
// @router /elementos_sin_asignar/ [get]
func (c *BodegaConsumoController) GetElementos() {

	defer errorctrl.ErrorControlController(c.Controller, "BodegaConsumoController")

	if v, err := bodegaConsumoHelper.GetElementosSinAsignar(); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAperturasKardex ...
// @Title GetAll
// @Description get Bodega-Consumo
// @Success 200 {object} []models.ElementoAperturaKardex
// @Failure 403
// @router /aperturas_kardex/ [get]
func (c *BodegaConsumoController) GetAperturasKardex() {

	defer errorctrl.ErrorControlController(c.Controller, "BodegaConsumoController")

	if v, err := bodegaConsumoHelper.GetAperturasKardex(); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()

}

// GetAll ...
// @Title GetAll
// @Description get Bodega-Consumo
// @Success 200 {object} []models.ExistenciasKardex
// @Failure 403
// @router /existencias_kardex/ [get]
func (c *BodegaConsumoController) GetAllExistencias() {

	defer errorctrl.ErrorControlController(c.Controller, "BodegaConsumoController")

	if v, err := bodegaConsumoHelper.GetExistenciasKardex(); err != nil {
		panic(err)
	} else {
		if v == nil {
			v = []map[string]interface{}{}
		}

		c.Data["json"] = v
		c.ServeJSON()
	}
}
