package controllers

import (
	"fmt"
	"io/ioutil"

	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
	"github.com/udistrital/arka_mid/models"
	// "github.com/udistrital/arka_mid/helpers/actaRecibido"
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

	// var archivo map[string]interface{}

	// Alertas
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})

	// if err := json.Unmarshal(c.Ctx.Input.RequestBody, &archivo); err == nil {
	// } else {
	// 	fmt.Println("err reading multipartFile", err)
	// 	alerta.Type = "error"
	// 	alerta.Code = "400"
	// 	alertas = append(alertas, "err reading file")
	// 	alerta.Body = alertas
	// 	c.Data["json"] = alerta
	// 	c.ServeJSON()
	// 	return
	// }

	// b, _ := archivo["archivo"]

	// Bytes, err := models.GetBytes(b)
	// if err != nil {
	// 	fmt.Println("no se pudo convertir")
	// }
	// fmt.Println(Bytes)

	// Lectura del archivo
	// xlFile, err := xlsx.OpenBinary(Bytes)

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

	multipartFile, _, err := c.GetFile("archivo")
	if err != nil {
		fmt.Println("err reading multipartFile", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, "err reading file")
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
	file, err := ioutil.ReadAll(multipartFile)
	if err != nil {
		fmt.Println("err reading file", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, "err reading file")
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}

	xlFile, err := xlsx.OpenBinary(file)
	if err != nil {
		fmt.Println("err reading file", err)
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, "err reading file")
		alerta.Body = alertas
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}

	Respuesta := make([]map[string]interface{}, 0)
	Elemento := make([]map[string]interface{}, 0)

	var hojas []string
	var campos []string
	var elementos [14]string
	for s, sheet := range xlFile.Sheets {

		if s == 0 {
			fmt.Println(sheet.Name)
			hojas = append(hojas, sheet.Name)
			for r, row := range sheet.Rows {
				if r == 0 {
					for _, cell := range row.Cells {
						campos = append(campos, cell.String())
					}
				} else {
					for i, cell := range row.Cells {
						elementos[i] = cell.String()
					}
					fmt.Println(elementos)
					if elementos[0] != "" {
						Elemento = append(Elemento, map[string]interface{}{
							"NivelInventariosId": elementos[0],
							"TipoBienId":         elementos[1],
							"SubgrupoCatalogoId": elementos[2],
							"Descripcion":        elementos[3],
							"Marca":              elementos[4],
							"Serie":              elementos[5],
							"Cantidad":           elementos[6],
							"UnidadMedida":       elementos[7],
							"ValorUnitario":      elementos[8],
							"Subtotal":           elementos[9],
							"Descuento":          elementos[10],
							"PorcentajeIvaId":    elementos[11],
							"ValorIva":           elementos[12],
							"ValorTotal":         elementos[13],
						})
					}
					for i := range row.Cells {
						elementos[i] = ""
					}
					fmt.Println(elementos)
				}
			}
		}
	}
	Respuesta = append(Respuesta, map[string]interface{}{
		"Hoja":      hojas,
		"Campos":    campos,
		"Elementos": Elemento,
	})
	c.Data["json"] = append(Respuesta)
	c.ServeJSON()
}
