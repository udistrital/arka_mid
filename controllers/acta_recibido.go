package controllers

import (
	"fmt"

	"github.com/astaxie/beego/logs"

	"encoding/json"

	"github.com/astaxie/beego"

	// "github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/models"
)

// ActaRecibidoController operations for ActaRecibido
type ActaRecibidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ActaRecibidoController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Create
// @Description create Acta_recibido
// @Param	body		body 	models.Acta_recibido	true		"body for Acta_recibido content"
// @Success 201 {object} models.Acta_recibido
// @Failure 403 body is empty
// @router / [post]
func (c *ActaRecibidoController) Post() {

	var archivo map[string]interface{}

	// Alertas
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &archivo); err == nil {
	} else {
		fmt.Println("err reading multipartFile", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, "err reading file")
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}

	b, _ := archivo["archivo"]
	logs.Info(b)

	// Lectura del archivo
	// xlFile, err := xlsx.OpenBinary(b)

	// if err != nil {
	// 	fmt.Println("err reading file", err)
	// 	alerta.Type = "error"
	// 	alerta.Code = "400"
	// 	alertas = append(alertas, "err reading file")
	// 	alerta.Body = alertas
	// 	c.Data["json"] = alerta
	// 	c.ServeJSON()
	// 	return
	// }

	// for _, sheet := range xlFile.Sheets {
	// 	for _, row := range sheet.Rows {
	// 		for _, cell := range row.Cells {
	// 			text := cell.String()
	// 			fmt.Printf("%s\n", text)
	// 		}
	// 	}
	// }
	c.Data["json"] = alertas
	c.ServeJSON()
	// }
}
