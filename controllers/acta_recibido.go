package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// ActaRecibidoController operations for ActaRecibido
type ActaRecibidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ActaRecibidoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetParametros", c.GetParametros)
	c.Mapping("GetElementosActa", c.GetElementosActa)
	c.Mapping("GetAllActas", c.GetAllActas)
}

// Post ...
// @Title Create
// @Description create Acta_recibido
// @Param	archivo	formData  file	true	"body for Acta_recibido content"
// @Success 201 {object} models.CargaMasivaElementosActa
// @Failure 403 body is empty
// @router / [post]
func (c *ActaRecibidoController) Post() {

	defer errorCtrl.ErrorControlController(c.Controller, "ActaRecibidoController")

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

// GetParametros ...
// @Title Consulta de valores paramétricos
// @Description Consulta a tablas paramétricas de los APIs acta_recibido_crud y catalogo_elementos_crud
// @Success 200 {object} models.ActaRecibido
// @Failure 404 not found resource
// @router / [get]
func (c *ActaRecibidoController) GetParametros() {

	defer errorCtrl.ErrorControlController(c.Controller, "ActaRecibidoController")

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
// @router /elementos/:id [get]
func (c *ActaRecibidoController) GetElementosActa() {

	defer errorCtrl.ErrorControlController(c.Controller, "ActaRecibidoController")

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

// GetAllActas ...
// @Title Get All Actas
// @Description get ActaRecibido
// @Param 	u					query	string	false	"WSO2 User. When specified, acts will be filtered upon the available roles for the specified user"
// @Param	limit				query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset				query	string	false	"Start position of result set. Must be an integer"
// @Param	Id					query	string	false	"Id para utilizar en query __in"
// @Param	TipoActaId			query	string	false	"Tipos de acta para utilizar en query __in"
// @Param	UnidadEjecutoraId	query	string	false	"Unidad ejecutora por las actas que se desean filtrar"
// @Param	EstadoActaId		query	string	false	"Estado del acta"
// @Param	FechaCreacion		query	string	false	"Fecha creación del acta: __in"
// @Param	FechaModificacion	query	string	false	"Fecha modificación del acta: __in"
// @Param	FechaVistoBueno		query	string	false	"Fecha aprobación del acta: __in"
// @Param	sortby				query	string	false	"Columna por la que se ordenan los resultados"
// @Param	order				query	string	false	"Orden de los resultados de acuerdo a la columna indicada"
// @Success 200 {object} []models.ActaRecibido
// @Failure 400 "Wrong IDs"
// @Failure 404 "not found resource"
// @Failure 500 "Unknown API Error"
// @Failure 502 "External API Error"
// @router /get_all_actas/ [get]
func (c *ActaRecibidoController) GetAllActas() {

	defer errorCtrl.ErrorControlController(c.Controller, "ActaRecibidoController")

	var WSO2user string
	var limit int64 = 10
	var offset int64

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

	reqStates := []string{}
	estados := c.GetString("EstadoActaId")
	if len(estados) > 0 {
		reqStates = strings.Split(estados, ",")
	}

	tipos := c.GetString("TipoActaId")
	unidadEjecutora := c.GetString("UnidadEjecutoraId")
	id := c.GetString("Id")
	creacion := c.GetString("FechaCreacion")
	modificacion := c.GetString("FechaModificacion")
	vistoBueno := c.GetString("FechaVistoBueno")
	sortby := c.GetString("sortby")
	order := c.GetString("order")

	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}

	if l, t, err := actaRecibido.GetAllActasRecibidoActivas(WSO2user, id, tipos, reqStates, creacion, modificacion, vistoBueno, unidadEjecutora, sortby, order, limit, offset); err == nil {
		c.Ctx.Output.Header("x-total-count", t)
		if l == nil {
			l = []map[string]interface{}{}
		}
		c.Data["json"] = l
	} else {
		panic(err)
	}

	c.ServeJSON()
}
