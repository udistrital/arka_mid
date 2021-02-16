package controllers

import (
	"encoding/json"
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
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "SalidaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
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
			panic(map[string]interface{}{
				"funcion": "Post",
				"err":     err,
				"status":  "502",
			})
		} else {
			panic(map[string]interface{}{
				"funcion": "Post",
				"err":     "No hubo respuesta",
				"status":  "404",
			})
		}
	} else {
		panic(map[string]interface{}{
			"funcion": "Post",
			"err":     "JSON mal estructurado",
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
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := salidaHelper.GetSalida(id)
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
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
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "SalidaController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
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
