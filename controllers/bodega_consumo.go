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
// @Param	id		path 	uint	true		"MovimientoId from Movimientos Arka CRUD"
// @Success 200 {object} models.BodegaConsumoSolicitud
// @Failure 400 "Wrong parameter (ID MUST be > 0)"
// @Failure 404 "Not found"
// @Failure 500 "Internal Error"
// @Failure 502 "Error with external API"
// @router /solicitud/:id [get]
func (c *BodegaConsumoController) GetOneSolicitud() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BodegaConsumoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

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
	// logs.Info(fmt.Sprintf("id: %d", id))

	if v, err := bodegaConsumoHelper.GetSolicitudById(id); err == nil {
		c.Data["json"] = v
	} else {
		panic(err)
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

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BodegaConsumoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

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
// @Success 200 {object} models.Bodega-Consumo
// @Failure 403
// @router /aperturas_kardex/ [get]
func (c *BodegaConsumoController) GetAperturasKardex() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BodegaConsumoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

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
// @Success 200 {object} models.Bodega-Consumo
// @Failure 403
// @router /existencias_kardex/ [get]
func (c *BodegaConsumoController) GetAllExistencias() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BodegaConsumoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500")
			}
		}
	}()

	if v, err := bodegaConsumoHelper.GetExistenciasKardex(); err != nil {
		panic(err)
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
/*
func (c *BodegaConsumoController) Put() {

}
*/

// Delete ...
// @Title Delete
// @Description delete the Bodega-Consumo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
/*
func (c *BodegaConsumoController) Delete() {

}
*/
