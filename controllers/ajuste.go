package controllers

import (
	"encoding/json"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/ajustesHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AjusteController operations for Ajuste
type AjusteController struct {
	beego.Controller
}

// URLMapping ...
func (c *AjusteController) URLMapping() {
	c.Mapping("PostManual", c.PostManual)
	c.Mapping("PostAjuste", c.PostAjuste)
	c.Mapping("GetOneManual", c.GetOneManual)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("GetElementos", c.GetElementos)
	c.Mapping("GetOneAuto", c.GetOneAuto)
}

// PostManual ...
// @Title Create
// @Description create Ajuste
// @Param	body		body 	models.PreTrAjuste	true		"body for Ajuste content"
// @Success 201 {object} models.Movimiento
// @Failure 403 body is empty
// @router / [post]
func (c *AjusteController) PostManual() {

	defer errorctrl.ErrorControlController(c.Controller, "AjusteController")

	var v *models.PreTrAjuste
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		panic(errorctrl.Error("Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)", err, "400"))
	} else {
		if v, err := ajustesHelper.PostAjuste(v); err != nil {
			logs.Error(err)
			c.Data["system"] = err
			c.Abort("404")
		} else {
			c.Data["json"] = v
		}
	}
	c.ServeJSON()
}

// PostAjuste ...
// @Title Create
// @Description create Ajuste
// @Param	body		body 	[]models.DetalleElemento_	true		"body for Ajuste content"
// @Success 201 {object} models.DetalleAjusteAutomatico
// @Failure 403 body is empty
// @router /automatico [post]
func (c *AjusteController) PostAjuste() {

	defer errorctrl.ErrorControlController(c.Controller, "AjusteController")

	var v []*models.DetalleElemento_
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		panic(errorctrl.Error("Post - json.Unmarshal(c.Ctx.Input.RequestBody, &v)", err, "400"))
	} else {
		if v, err := ajustesHelper.GenerarAjusteAutomatico(v); err != nil {
			logs.Error(err)
			c.Data["system"] = err
			c.Abort("404")
		} else {
			c.Data["json"] = v
		}
	}
	c.ServeJSON()
}

// GetOneManual ...
// @Title GetOneManual
// @Description get Ajuste by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DetalleAjuste
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AjusteController) GetOneManual() {

	defer errorctrl.ErrorControlController(c.Controller, "AjusteController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un ajuste válido")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": `GetOneManual - c.GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if respuesta, err := ajustesHelper.GetDetalleAjuste(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "GetOne - ajustesHelper.GetDetalleAjuste(id)",
			"err":     errors.New("no se obtuvo respuesta al consultar el ajuste"),
			"status":  "404",
		})
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Ajuste
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Ajuste
// @Failure 403
// @router / [get]
func (c *AjusteController) GetAll() {

	defer errorctrl.ErrorControlController(c.Controller, "AjusteController")

	panic(map[string]interface{}{
		"funcion": "All",
		"err":     errors.New("no implementado"),
		"status":  "501",
	})
}

// Put ...
// @Title Put
// @Description update the Ajuste
// @Param	id		path 	string	true		"The id you want to update"
// @Success 200 {object} models.Movimiento
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AjusteController) Put() {

	defer errorctrl.ErrorControlController(c.Controller, "AjusteController - Unhandled Error!")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un ajuste válido")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": `Put - GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if respuesta, err := ajustesHelper.AprobarAjuste(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "Put - ajustesHelper.AprobarAjuste(id)",
			"err":     errors.New("no se obtuvo respuesta al aprobar el ajuste"),
			"status":  "404",
		})
	}

	c.ServeJSON()
}

// GetElementos ...
// @Title GetElementos
// @Description Retorna la lista de elementos asociados a un acta con su respectiva vida útil y valor residual iniciales
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} []models.DetalleElemento__
// @Failure 403 :id is empty
// @router /automatico/elementos/:id [get]
func (c *AjusteController) GetElementos() {

	defer errorctrl.ErrorControlController(c.Controller, "AjusteController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un acta válida")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": `GetElementos - c.GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if respuesta, err := ajustesHelper.GetDetalleElementosActa(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "GetElementos - ajustesHelper.GetDetalleElementosActa(id)",
			"err":     errors.New("no se obtuvo respuesta al consultar los elementos"),
			"status":  "404",
		})
	}

	c.ServeJSON()
}

// GetOneAuto ...
// @Title GetOneAuto
// @Description Retorna la lista de elementos asociados a un ajuste y su transaccion contable
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DetalleAjusteAutomatico
// @Failure 403 :id is empty
// @router /automatico/:id [get]
func (c *AjusteController) GetOneAuto() {

	defer errorctrl.ErrorControlController(c.Controller, "AjusteController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar un ajuste válido")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": `GetOneAuto - c.GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if respuesta, err := ajustesHelper.GetAjusteAutomatico(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "GetOneAuto - ajustesHelper.GetAjusteAutomatico(id)",
			"err":     errors.New("no se obtuvo respuesta al consultar el ajuste"),
			"status":  "404",
		})
	}

	c.ServeJSON()
}
