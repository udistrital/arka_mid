package controllers

import (
	"encoding/json"
	"errors"
	"strconv"

	// "github.com/udistrital/utils_oas/formatdata"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/models"
)

// SalidaController operations for Salida
type SalidaController struct {
	beego.Controller
}

// URLMapping ...
func (c *SalidaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetSalida", c.GetSalida)
	c.Mapping("GetSalidas", c.GetSalidas)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Post transaccion salidas asociadas a una entrada
// @Description Realiza la aprobacion de una salida en caso de especificarse un Id, de lo contrario, genera los consecutivos de las salidas y hace el respectivo registro en api movimientos_arka_crud
// @Param	salidaId	query 	string					false		"Id del movimiento que se desea aprobar"
// @Param	body		body 	models.SalidaGeneral	true		"Informacion de las salidas y elementos asociados a cada una de ellas. Se valida solo si el id es 0""
// @Success 200 {object} models.SalidaGeneral
// @Failure 403 body is empty
// @router / [post]
func (c *SalidaController) Post() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SalidaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()
	var salidaId int = 0

	if v, err := c.GetInt("salidaId"); err == nil {
		salidaId = v
	}
	if salidaId > 0 {
		if respuesta, err := salidaHelper.AprobarSalida(salidaId); err == nil && respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
		} else {
			if err == nil {
				panic(map[string]interface{}{
					"funcion": "Post - salidaHelper.AprobarSalida(salidaId)",
					"err":     errors.New("No se obtuvo respuesta al aprobar la salida"),
					"status":  "404",
				})
			}
			panic(err)
		}
	} else {
		var v models.SalidaGeneral
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			if respuesta, err := salidaHelper.PostTrSalidas(&v); err == nil && respuesta != nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = respuesta
			} else {
				status := "400"
				if err == nil {
					err = map[string]interface{}{
						"err": errors.New("No se obtuvo respuesta al registrar la(s) salida(s)"),
					}
					status = "404"
				}
				logs.Error(err)
				panic(map[string]interface{}{
					"funcion": "Post - salidaHelper.PostTrSalidas(&v)",
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
	}

	c.ServeJSON()
}

// GetSalida ...
// @Title Get User
// @Description get Salida by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.TrSalida
// @Failure 404 not found resource
// @router /:id [get]
func (c *SalidaController) GetSalida() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SalidaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		if err == nil {
			err = errors.New("id MUST be > 0")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetSalida - strconv.Atoi(idStr)",
			"err":     err,
			"status":  "400",
		})
	}
	if v, err := salidaHelper.GetSalida(id); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetSalidas ...
// @Title Get User
// @Description get Entradas
// @Param	tramite_only		query	bool false	"Retornar salidas únicamente en estado En Trámite"
// @Success 200 {object} []models.Movimiento
// @Failure 404 not found resource
// @router / [get]
func (c *SalidaController) GetSalidas() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SalidaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()
	var tramiteOnly bool

	if v, err := c.GetBool("tramite_only"); err == nil {
		tramiteOnly = v
	}
	if v, err := salidaHelper.GetSalidas(tramiteOnly); err == nil {
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

// Put ...
// @Title Put transaccion salidas generadas a partir de otra
// @Description genera los consecutivos de las nuevas salidas generadas y hace el put en el api movimientos_arka_crud
// @Param	id		path 	int						true	"Id de la salida original"
// @Param	body	body 	models.SalidaGeneral	true	"Informacion de las salidas y elementos asociados a cada una de ellas. Se valida solo si el id es 0""
// @Success 200 {object} models.SalidaGeneral
// @Failure 403 body is empty
// @router /:id [put]
func (c *SalidaController) Put() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SalidaController" + "/" + (localError["funcion"]).(string))
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
			err = errors.New("Se debe especificar una salida válida")
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

	var v models.SalidaGeneral
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Put - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
			"err":     err,
			"status":  "400",
		})
	}

	if respuesta, err := salidaHelper.PutTrSalidas(&v, id); err == nil && respuesta != nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "Put - salidaHelper.PutTrSalidas(&v, id)",
			"err":     errors.New("No se obtuvo respuesta al actualizar la salida"),
			"status":  "404",
		})
	}

	c.ServeJSON()
}
