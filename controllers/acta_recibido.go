package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	//"github.com/udistrital/acta_recibido_crud/models"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
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

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ActaRecibidoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	if multipartFile, _, err := c.GetFile("archivo"); err == nil {
		if Archivo, err := actaRecibido.DecodeXlsx2Json(multipartFile); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = Archivo
		} else {
			panic(err)
		}
	} else {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Post",
			"err":     err,
			"status":  "400",
		})
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

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ActaRecibidoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	if l, err := actaRecibido.GetAllParametrosActa(); err != nil {
		panic(err)
	} else {
		c.Data["json"] = l
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
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ActaRecibidoController" + "/" + (localError["funcion"]).(string))
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
	} else {
		if err == nil {
			err = fmt.Errorf("the Id MUST be greater than 0 - Got:%v", idStr)
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetElementosActa - strconv.Atoi(idStr)",
			"err":     err,
			"status":  "400",
		})
	}
	// fmt.Printf("id: %v\n", id)

	if v, err := actaRecibido.GetElementos(id, nil); err == nil {
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

// GetAllElementosConsumo ...
// @Title GetAllElementosConsumo
// @Description Trae todos los elementos de consumo
// @Success 200 {object} models.Elemento
// @Failure 404 not found resource
// @router /elementosconsumo/ [get]
func (c *ActaRecibidoController) GetAllElementosConsumo() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ActaRecibidoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500")
			}
		}
	}()

	if l, err := actaRecibido.GetAllElementosConsumo(); err != nil {
		panic(err)
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
			err := errors.New("bad syntax. States MUST be comma separated")
			logs.Error(err)
			panic(map[string]interface{}{
				"funcion": "GetAllActas - c.GetString(\"states\")",
				"err":     err,
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
			err := errors.New("user not specified in parameter value")
			logs.Error(err)
			panic(map[string]interface{}{
				"funcion": "GetAllActas - c.GetString(\"u\")",
				"err":     err,
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
