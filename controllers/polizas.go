package controllers

import (
	"errors"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	polizasHelper "github.com/udistrital/arka_mid/helpers/polizasHelper"
)

// PolizasController operations for Polizas
type PolizasController struct {
	beego.Controller
}

// URLMapping ...
func (c *PolizasController) URLMapping() {
	// c.Mapping("Post", c.Post)
	// c.Mapping("GetOne", c.GetPolizaId)
	c.Mapping("GetAll", c.GetAllElementosPoliza)
	//c.Mapping("Put", c.Put)
}

// Post ...
// @Title Create
// @Description Registro de una póliza
// @Param	body		body 	models.Elemento_campo	true		"Contenido para registrar una póliza"
// @Success 201 {object} models.Elemento_campo
// @Failure 403 body is empty
// @router / [post]
//func (c *PolizasController) Post() {

// 	defer func() {
// 		if err := recover(); err != nil {
// 			logs.Error(err)
// 			localError := err.(map[string]interface{})
// 			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "PolizasController" + "/" + (localError["funcion"]).(string))
// 			c.Data["data"] = (localError["err"])
// 			if status, ok := localError["status"]; ok {
// 				c.Abort(status.(string))
// 			} else {
// 				c.Abort("500")
// 			}
// 		}
// 	}()

// 	var v models.Elemento_campo
// 	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
// 		logs.Error(err)
// 		panic(map[string]interface{}{
// 			"funcion": "Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)",
// 			"err":     err,
// 			"status":  "400",
// 		})
// 	}

// 	if respuesta, err := polizasHelper.RegistrarPoliza(&v); err == nil && respuesta != nil {
// 		c.Ctx.Output.SetStatus(201)
// 		c.Data["json"] = respuesta
// 	} else {
// 		if err != nil {
// 			panic(err)
// 		}

// 		panic(map[string]interface{}{
// 			"funcion": "polizasHelper.RegistrarPoliza(&v)",
// 			"err":     errors.New("No se obtuvo respuesta al registrar la póliza"),
// 			"status":  "404",
// 		})
// 	}

// 	c.ServeJSON()

// }

// GetOne ...
// @Title GetOne
// @Description get Polizas by id
// @Param	id		path 	string	true		"Id de poliza a consultar"
// @Success 200 {object} models.Elemento_campo
// @Failure 403 :id is empty
// @router /poliza/:id [get]
// func (c *PolizasController) GetPolizaId() {

// 	defer errorctrl.ErrorControlController(c.Controller, "PolizasController - Unhandled Error!")
// 	var id int
// 	if v, err := c.GetInt(":ElementoCampoId"); err != nil || v <= 0 {
// 		if err == nil {
// 			err = errors.New("Se debe especificar un número valido")
// 		}
// 		panic(errorctrl.Error("GetPolizaId - c.GetInt(\":ElementoCampoId\")", err, "400"))
// 	} else {
// 		id = v
// 	}

// 	if respuesta, err := polizasHelper.GetPoliza(id); err == nil || respuesta != nil {
// 		c.Data["json"] = respuesta
// 	} else {
// 		if err != nil {
// 			panic(err)
// 		}
// 		panic(errorctrl.Error("GetPolizaId - polizaHelper.GetPolizaId(id)", err, "404"))
// 	}

// 	c.ServeJSON()

// }

// GetAll ...
// @Title GetAll
// @Description get AllElementosParaPoliza
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	int		false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	int		false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Elemento
// @Failure 404 not found resource
// @router /ActasPoliza [get]
func (c *PolizasController) GetAllElementosPoliza() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "PolizasController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	var fields []string                 //filtra un valor en concreto
	var order []string                  //ordenar valores
	var query = make(map[string]string) //Hace una busqueda de valores
	var limit int = 10                  //limita la lista para poder paginarla
	var offset int                      //comienza la lista

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	if l, err := polizasHelper.GetElementosPoliza(offset, limit, fields, order, query); err != nil {
		panic(err)
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()

}

// Put ...
// @Title Put
// @Description update the Polizas
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Polizas	true		"body for Polizas content"
// @Success 200 {object} models.Polizas
// @Failure 403 :id is not int
// @router /:id [put]
//func (c *PolizasController) Put() {

//}
