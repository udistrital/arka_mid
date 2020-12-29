package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/bodegaConsumoHelper"
)

// BodegaConsumoController operations for Bodega-Consumo
type BodegaConsumoController struct {
	beego.Controller
}

// URLMapping ...
func (c *BodegaConsumoController) URLMapping() {
	c.Mapping("GetOne", c.GetOneSolicitud)
	c.Mapping("GetAll", c.GetElementos)
	c.Mapping("Get", c.GetAllExistencias)
	

}

// GetOneSolicitud ...
// @Title GetOneSolicitud
// @Description get Bodega-Consumo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object}{"Id": int,"FechaCreacion": date,"Observacion": string,"Elementos": {"Id": int,"Nombre":string,"Marca": string,"Serie": string,"CantidadDisponible": int,"CantidadSolicitada": int,	"ValorUnitario": float,} }
// @Failure 403 :id is empty
// @router /solicitud/:id [get]
func (c *BodegaConsumoController) GetOneSolicitud() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	fmt.Println("id: ", id)
	v, err := bodegaConsumoHelper.GetSolicitudById(id)
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Bodega-Consumo
// @Success 200 {object} models.Bodega-Consumo
// @Failure 403
// @router /elementos_sin_asignar/ [get]
func (c *BodegaConsumoController) GetElementos() {

	v, err := bodegaConsumoHelper.GetElementosSinAsignar()
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAperturasKardex ...
// @Title GetAll
// @Description get Bodega-Consumo
// @Success 200 {object} models.Bodega-Consumo
// @Failure 403
// @router /aperturas_kardex/ [get]
func (c *BodegaConsumoController) GetAperturasKardex() {

	v, err := bodegaConsumoHelper.GetAperturasKardex()
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()


}

// GetAll ...
// @Title GetAll
// @Description get Bodega-Consumo
// @Success 200 {object} models.Bodega-Consumo
// @Failure 403
// @router /existencias_kardex/ [get]
func (c *BodegaConsumoController) GetAllExistencias() {

	v, err := bodegaConsumoHelper.GetExistenciasKardex()
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Bodega-Consumo
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Bodega-Consumo	true		"body for Bodega-Consumo content"
// @Success 200 {object} models.Bodega-Consumo
// @Failure 403 :id is not int
// @router /:id [put]
func (c *BodegaConsumoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Bodega-Consumo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *BodegaConsumoController) Delete() {

}