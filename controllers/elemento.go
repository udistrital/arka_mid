package controllers

import (
	"encoding/json"
	"fmt"

	// "github.com/udistrital/utils_oas/formatdata"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/models"
)

// SalidaController operations for Salida
type ElementoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ElementoController) URLMapping() {
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Create
// @Description create Salida
// @Param	body		body 	models.SalidaGeneral	true		"body for Salida content"
// @Success 200 {object} []models.TrSalida2
// @Failure 403 body is empty
// @router / [post]
func (c *ElementoController) Put() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "ElementoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	// fmt.Printf("body: %v\n", c.Ctx.Input.RequestBody)
	var v models.Elemento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		panic(map[string]interface{}{
			"funcion": "Post",
			"err":     "JSON mal estructurado",
			"status":  "400",
		})
	}

	fmt.Printf("valores")
	// formatdata.JsonPrint(v)
	if respuesta, err := salidaHelper.AsignarPlaca(&v); err == nil && respuesta != nil {
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

	c.ServeJSON()
}
