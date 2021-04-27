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
	c.Mapping("Get", c.GetSalida)
	c.Mapping("GetAll", c.GetSalidas)
}

// Post ...
// @Title Create
// @Description create Salida
// @Param	body		body 	models.SalidaGeneral	true		"body for Salida content"
// @Success 200 {object} []models.TrSalida2
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

	// fmt.Printf("body: %v\n", c.Ctx.Input.RequestBody)
	var v models.SalidaGeneral
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		// fmt.Printf("valores: %#v\n", v)
		// formatdata.JsonPrint(v)
		if respuesta, err := salidaHelper.AddSalida(&v); err == nil && respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
		} else if err != nil {
			panic(err)
		} else {
			panic(map[string]interface{}{
				"funcion": "Post",
				"err":     "No hubo respuesta",
				"status":  "404",
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

	if v, err := salidaHelper.GetSalidas(); err == nil {
		c.Data["json"] = v
	} else {
		panic(err)
	}
	c.ServeJSON()

}
