package actaRecibido

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"
	"github.com/udistrital/utils_oas/request"
)

// GetAllActasRecibido ...
func GetAllActasRecibido() (historicoActa interface{}, outputError map[string]interface{}) {
	// if idUser != 0 { // (1) error parametro
	if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("actaRecibidoService")+"historico_acta?query=ActaRecibidoId.Activo:True", &historicoActa); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			// c.Data["json"] = response
			// c.ServeJSON()
			return historicoActa, nil
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetAllActasRecibido:GetAllActasRecibido", "Error": response.Status}
			return outputError, nil
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return outputError, nil
	}
	// c.ServeJSON()
	// } else {
	// 	logs.Info("Error (1) Parametro")
	// 	outputError = map[string]interface{}{"Function": "FuncionalidadMidController:getUserAgora", "Error": "null parameter"}
	// 	return nil, outputError
	// }
}

// GetActasRecibidoTipo ...
func GetAllParametrosActa() (Parametros map[string]interface{}, outputError map[string]interface{}) {

	if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("actaRecibidoService")+"historico_acta", &Parametros); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			return Parametros, nil
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetAllActasRecibido:GetAllActasRecibido", "Error": response.Status}
			return outputError, nil
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return outputError, nil
	}
}

// "PostDecodeXlsx2Json ..."
func DecodeXlsx2Json(c multipart.File) (Archivo []map[string]interface{}, outputError map[string]interface{}) {

	file, err := ioutil.ReadAll(c)
	if err != nil {
		fmt.Println("err reading file", err)
		logs.Info("Error (1) error de recepcion")
		outputError = map[string]interface{}{"Function": "PostDecodeXlsx2Json", "Error": 400}
		return nil, outputError
	}
	xlFile, err := xlsx.OpenBinary(file)
	if err != nil {
		fmt.Println("err reading file", err)
		logs.Info("Error (1) error de recepcion")
		outputError = map[string]interface{}{"Function": "PostDecodeXlsx2Json", "Error": 400}
		return nil, outputError
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
	return Respuesta, nil
}
