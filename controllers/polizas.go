package controllers

import (
	"errors"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	polizasHelper "github.com/udistrital/arka_mid/helpers/polizasHelper"
	"github.com/udistrital/utils_oas/errorctrl"
)

// PolizasController operations for Polizas
type PolizasController struct {
	beego.Controller
}

// URLMapping ...
func (c *PolizasController) URLMapping() {
	c.Mapping("GetAllElementosPoliza", c.GetAllElementosPoliza)
}

// GetAll ...
// @Title GetAll
// @Description get AllElementosParaPoliza
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ... POR IMPLEMENTAR"
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	int		false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	int		false	"Start position of result set. Must be an integer"
// @Success 200 {object} []models.Elemento
// @Failure 404 not found resource
// @router /AllElementosPoliza [get]
func (c *PolizasController) GetAllElementosPoliza() {

	defer errorctrl.ErrorControlController(c.Controller, "PolizasController")

	var fields []string                 //filtra un valor en concreto
	var sortby []string                 //
	var order []string                  //ordenar valores
	var query = make(map[string]string) //Hace una busqueda de valores
	var limit int = 10                  //limita la lista para poder paginarla
	var offset int                      //comienza la lista

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt("limit", limit); err == nil {
		limit = v
	} else {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetAllElementosPoliza - c.GetInt(\"limit\")",
			"err":     err,
			"status":  "400",
		})
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt("offset", offset); err == nil {
		offset = v
	} else {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetAllElementosPoliza - c.GetInt(\"offset\")",
			"err":     err,
			"status":  "400",
		})
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
				c.Data["json"] = errors.New("error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	if l, err := polizasHelper.GetElementosPoliza(offset, limit, fields, order, query, sortby); err != nil {
		panic(err)
	} else {
		if l == nil {
			c.Data["json"] = []interface{}{}
		} else {
			c.Data["json"] = l
		}
	}
	c.ServeJSON()
}
