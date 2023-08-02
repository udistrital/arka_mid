package controllers

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// EntradaController operations for Entrada
type EntradaController struct {
	beego.Controller
}

// URLMapping ...
func (c *EntradaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
}

// Post ...
// @Title Post
// @Description Transaccion entrada. Estado de registro o aprobacion
// @Param	entradaId	query	string						false	"Id del movimiento que se desea aprobar"
// @Param	etl			query	bool						false	"Indica si la entrada se registra a partir del ETL"
// @Param	aprobar		query	bool						false	"Indica si la entrada se debe aprobar"
// @Param	body		body	models.TransaccionEntrada	false	"Detalles de la entrada. Se valida solo si el id es 0"
// @Success 201 {object} models.Movimiento
// @Failure 403 body is empty
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *EntradaController) Post() {

	defer errorctrl.ErrorControlController(c.Controller, "EntradaController")

	var (
		entradaId int
		etl       bool
		aprobar   bool
	)

	if v, err := c.GetInt("entradaId", 0); err == nil {
		entradaId = v
	}

	if v, err := c.GetBool("etl", false); err == nil {
		etl = v
	}

	if v, err := c.GetBool("aprobar", false); err == nil {
		aprobar = v
	}

	if aprobar && entradaId > 0 {

		var res models.ResultadoMovimiento
		if err := entradaHelper.AprobarEntrada(entradaId, &res); err != nil {
			if err == nil {
				panic(map[string]interface{}{
					"funcion": "Post - entradaHelper.AprobarEntrada(entradaId)",
					"err":     errors.New("no se obtuvo respuesta al aprobar la entrada."),
					"status":  "400",
				})
			}
			panic(err)
		} else {
			c.Data["json"] = res
		}
	} else if !aprobar {

		var (
			v       models.TransaccionEntrada
			entrada models.ResultadoMovimiento
		)

		if err := utilsHelper.Unmarshal(string(c.Ctx.Input.RequestBody), &v); err != nil {
			panic(err)
		}

		if entradaId > 0 {
			if err := entradaHelper.UpdateEntrada(&v, entradaId, &entrada); err != nil {
				panic(err)
			}
		} else if entradaId == 0 {
			if err := entradaHelper.RegistrarEntrada(&v, etl, &entrada); err != nil {
				panic(err)
			}
		}

		c.Data["json"] = entrada
	}

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Detalle de entrada por Id. Retorna la transaccion contable si la entrada ya  fue aprobada
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DetalleEntrada
// @Failure 403 :id is empty
// @router /:id [get]
func (c *EntradaController) GetOne() {

	defer errorctrl.ErrorControlController(c.Controller, "EntradaController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una entrada válida")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": `GetOne - c.GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	if respuesta, err := entradaHelper.DetalleEntrada(id); err == nil || respuesta != nil {
		c.Data["json"] = respuesta
	} else {
		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "GetOne - entradaHelper.DetalleEntrada(id)",
			"err":     errors.New("no se obtuvo respuesta al consultar la anetrada"),
			"status":  "404",
		})
	}

	c.ServeJSON()
}
