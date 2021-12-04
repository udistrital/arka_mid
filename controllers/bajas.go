package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/bajasHelper"
	"github.com/udistrital/arka_mid/models"
	// "github.com/udistrital/arka_mid/models"
)

// BajaController
type BajaController struct {
	beego.Controller
}

// URLMapping ...
func (c *BajaController) URLMapping() {
	c.Mapping("Get", c.GetElemento)
	c.Mapping("Put", c.Put)
	c.Mapping("GetElemento", c.GetDetalleElemento)
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

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BajaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("Se debe especificar una baja válida")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Put - GetInt(\":id\")",
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	var v *models.TrSoporteMovimiento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if respuesta, err := bajasHelper.ActualizarBaja(v, id); err == nil && respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
		} else {
			if err != nil {
				panic(err)
			}

			panic(map[string]interface{}{
				"funcion": "Put - bajasHelper.ActualizarBaja(v, id)",
				"err":     errors.New("No se obtuvo respuesta al actualizar la baja"),
				"status":  "404",
			})
		}
	} else {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Put - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
			"err":     err,
			"status":  "400",
		})
	}

	c.ServeJSON()
}

// GetElemento ...
// @Title Get User
// @Description get Salida by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Salida
// @Failure 404 not found resource
// @router /elemento/:id [get]
func (c *BajaController) GetElemento() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BajaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	logs.Info(idStr)
	var id int
	if idConv, err := strconv.Atoi(idStr); err == nil && idConv > 0 {
		id = idConv
	} else if err != nil {
		panic(err)
	} else {
		panic(map[string]interface{}{
			"funcion": "GetElemento",
			"err":     "El ID debe ser mayor a 0",
			"status":  "400",
		})
	}

	if v, err := bajasHelper.TraerDatosElemento(id); err == nil {
		c.Data["json"] = v
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// Getsolicitud...
// @Title Get User
// @Description get Salida by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Salida
// @Failure 404 "not found resource"
// @Failure 500 "Unknown API Error"
// @Failure 502 "External API Error"
// @router /solicitud/:id [get]
func (c *BajaController) GetSolicitud() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BajaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	fmt.Println(idStr)

	if v, err := bajasHelper.TraerDetalle(id); err == nil {
		c.Data["json"] = v
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description Consulta todas las bajas y permite filtrar por las que estan para revision de almacen o comite
// @Param	revComite	query 	bool	false	"Indica si se traen las bajas en espera de comite. Tiene prioridad sobre la revision de almacen"
// @Param	revAlmacen	query 	bool	false	"Indica si se traen las bajas pendientes por revisar"
// @Success 200 {object} []models.DetalleBaja
// @Failure 404 not found resource
// @router /solicitud/ [get]
func (c *BajaController) GetAll() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BajaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

	revComite := false
	revAlmacen := false
	if v, err := c.GetBool("revComite"); err != nil {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetAll - GetBool(\"revComite\")",
			"err":     err,
			"status":  "400",
		})
	} else if v == false {
		if v, err := c.GetBool("revAlmacen"); err != nil {
			logs.Error(err)
			panic(map[string]interface{}{
				"funcion": "GetAll - GetBool(\"revAlmacen\")",
				"err":     err,
				"status":  "400",
			})
		} else {
			revAlmacen = v
		}
	} else {
		revComite = v
	}

	if v, err := bajasHelper.GetAllSolicitudes(revComite, revAlmacen); err == nil {
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

// GetDetalleElemento ...
// @Title GetElemento
// @Description Get Info relacionada con el historial de un determinado elemento
// @Param	id	path	int	true	"id del elemento en el api acta_recibido_crud"
// @Success 200 {object} models.DetalleElementoBaja
// @Failure 404 not found resource
// @router /elemento/:id [get]
func (c *BajaController) GetDetalleElemento() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "BajaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("Se debe especificar un elemento válido")
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

	l, err := bajasHelper.GetDetalleElemento(id)
	if err != nil {
		panic(err)
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}
