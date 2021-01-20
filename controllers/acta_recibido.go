package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	//"github.com/udistrital/acta_recibido_crud/models"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/actaRecibidoHelper"
)

// ActaRecibidoController operations for ActaRecibido
type ActaRecibidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ActaRecibidoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetElementosActa", c.GetElementosActa)
	c.Mapping("GetElementosConsumo", c.GetAllElementosConsumo)
}

// Post ...
// @Title Create
// @Description create Acta_recibido
// @Param	archivo	formData  file	true	"body for Acta_recibido content"
// @Success 201 {}
// @Failure 403 body is empty
// @router / [post]
func (c *ActaRecibidoController) Post() {
	fmt.Println(c.GetFile("archivo"))
	if multipartFile, _, err := c.GetFile("archivo"); err == nil {
		if Archivo, err := actaRecibido.DecodeXlsx2Json(multipartFile); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = Archivo
		} else {
			c.Data["system"] = err
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get ActaRecibido
// @Success 200 {object} models.ActaRecibido
// @Failure 404 not found resource
// @router / [get]
func (c *ActaRecibidoController) GetAll() {

	fmt.Println("hola")
	l, err := actaRecibido.GetAllParametrosActa()
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// GetActasByTipo ...
// @Title GetActasRecibidoTipo
// @Description Devuelve las todas las actas de recibido
// @Param	id		path 	string	true		"id del acta"
// @Success 200 {object} models.Acta_recibido
// @Failure 403
// @router /get_actas_recibido_tipo/:tipo [get]
func (c *ActaRecibidoController) GetActasByTipo() {
	tipoStr := c.Ctx.Input.Param(":tipo")
	tipo, _ := strconv.Atoi(tipoStr)
	v, err := actaRecibidoHelper.GetActasRecibidoTipo(tipo)
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetElementosActa ...
// @Title Get Elementos
// @Description get Elementos by id
// @Param	id		path 	int	true		"id del acta"
// @Success 200 {object} []models.Elemento
// @Success 204 Empty response (Due to Act not found or without elements)
// @Failure 400 Wrong ID (MUST be greater than 0)
// @Failure 404 not found resource
// @Failure 500 Internal Error
// @Failure 502 Error with external API
// @router /get_elementos_acta/:id [get]
func (c *ActaRecibidoController) GetElementosActa() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "ActaRecibidoController" + "/" + (localError["funcion"]).(string))
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
	if idTest, err := strconv.Atoi(idStr); err == nil && idTest > 0 {
		id = idTest
	} else if err != nil {
		panic(map[string]interface{}{
			"funcion": "GetElementosActa",
			"err":     err,
			"status":  "400",
		})
	} else {
		panic(map[string]interface{}{
			"funcion": "GetElementosActa",
			"err":     "The Id MUST be greater than 0",
			"status":  "400",
		})
	}
	// fmt.Printf("id: %v\n", id)

	if v, err := actaRecibido.GetElementos(id); err == nil {
		c.Data["json"] = v
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GetSoportesActa ...
// @Title Get Soportes
// @Description get Soportes by id
// @Param	body	body 	models.Entrada	true
// @Success 200 {object} []models.AsignacionEspacioFisicoDependencia
// @Failure 404 not found resource
// @router /get_soportes_acta/:id [get]
func (c *ActaRecibidoController) GetSoportesActa() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := actaRecibidoHelper.GetSoportes(id)
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

// GetAllElementosConsumo ...
// @Title GetAllElementosConsumo
// @Description Trae todos los elementos de consumo
// @Success 200 {object} models.Elemento
// @Failure 404 not found resource
// @router /elementosconsumo/ [get]
func (c *ActaRecibidoController) GetAllElementosConsumo() {

	fmt.Println("hola hola")
	l, err := actaRecibido.GetAllElementosConsumo()
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// GetAllActas ...
// @Title Get All Actas
// @Description get ActaRecibido
// @Param	states	query	string	false	"If specified, returns only acts with the specified state(s) from ACTA_RECIBIDO_SERVICE / estado_acta, separated by commas"
// @Param u query string false "WSO2 User. When specified, acts will be filtered upon the available roles for the specified user"
// @Success 200 {object} []models.ActaRecibido
// @Failure 400 "Wrong IDs"
// @Failure 404 "not found resource"
// @Failure 500 "Unknown API Error"
// @Failure 502 "External API Error"
// @router /get_all_actas/ [get]
func (c *ActaRecibidoController) GetAllActas() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ActaRecibidoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Unhandled Error!
			}
		}
	}()

	var reqStates []string
	var WSO2user string

	if v := c.GetString("states"); v != "" {
		valido := false
		states := strings.Split(v, ",")
		for _, state := range states {
			state = strings.TrimSpace(state)
			if state != "" {
				reqStates = append(reqStates, state)
				valido = true
			}
		}

		if !valido {
			panic(map[string]interface{}{
				"funcion": "GetAllActas",
				"err":     "Bad syntax. Acts MUST be comma separated",
				"status":  "400",
			})
		}
	}
	// fmt.Print("ESTADOS SOLICITADOS: ")
	// fmt.Println(reqStates)

	if v := c.GetString("u"); v != "" {
		valido := false
		user := strings.TrimSpace(v)
		if user != "" {
			WSO2user = v
			valido = true
		}
		if !valido {
			panic(map[string]interface{}{
				"funcion": "GetAllActas",
				"err":     "Bad syntax",
				"status":  "400",
			})
		}
	}

	if l, err := actaRecibido.GetAllActasRecibidoActivas(reqStates, WSO2user); err == nil {
		// fmt.Print("DATA FINAL: ")
		// fmt.Println(l)
		c.Data["json"] = l
	} else {
		panic(err)
	}

	c.ServeJSON()
}
