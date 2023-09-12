package controllers

import (
	"errors"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
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
// @Param	id				path	int		true	"subgroupoId"
// @Param	movimientoId	query	string	false	TipoMovimientoId o SubtipoMovimientoId que se desea filtrar"
// @Success 200 {object} models.DetalleCuentasSubgrupo
// @Failure 403
// @router /cuentas_contables/:id [get]
func (c *CatalogoElementosController) GetOne() {

	defer errorCtrl.ErrorControlController(c.Controller, "CatalogoElementosController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una subgrupo válido")
		}
		panic(errorCtrl.Error(`GetOne - c.GetInt(":id")`, err, "400"))
	} else {
		id = v
	}

	movimientoId, err := c.GetInt("movimientoId", 0)
	if err != nil {
		panic(errorCtrl.Error(`GetOne - c.GetInt("movimientoId")`, err, "400"))
	}

	var cuentas = make([]models.DetalleCuentasSubgrupo, 0)
	if err := catalogoElementosHelper.GetCuentasContablesSubgrupo(id, movimientoId, &cuentas); err != nil {
		panic(err)
	} else {
		c.Data["json"] = cuentas
	}
	c.ServeJSON()
}
