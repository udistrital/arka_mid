package controllers

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// CatalogoElementosController operations for Catalogo
type CatalogoElementosController struct {
	beego.Controller
}

// URLMapping ...
func (c *CatalogoElementosController) URLMapping() {
	c.Mapping("GetOne", c.GetOne)
}

// GetOne ...
// @Title GetCuentasSubgrupoById
// @Description Devuelve el detalle de la última cuenta de cada movimiento requerido y subgrupo determinado
// @Param	id		path 	int	true		"subgroupoId"
// @Success 200 {object} models.DetalleCuentasSubgrupo
// @Failure 403
// @router /cuentas_contables/:id [get]
func (c *CatalogoElementosController) GetOne() {

	defer errorctrl.ErrorControlController(c.Controller, "CatalogoElementosController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("Se debe especificar una subgrupo válido")
		}
		panic(errorctrl.Error(`GetOne - c.GetInt(":id")`, err, "400"))
	} else {
		id = v
	}

	if v, err := catalogoElementosHelper.GetCuentasContablesSubgrupo(id); err != nil {
		panic(err)
	} else {
		if v == nil {
			v = []*models.DetalleCuentasSubgrupo{}
		}
		c.Data["json"] = v
	}
	c.ServeJSON()
}
