package controllers

import (
	//"github.com/udistrital/acta_recibido_crud/models"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// CatalogoElementosController operations for Catalogo
type CatalogoElementosController struct {
	beego.Controller
}

// URLMapping ...
func (c *CatalogoElementosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.GetAll)
	c.Mapping("Get", c.GetAll2)
	c.Mapping("GetOne", c.GetOne)
}

// GetAll ...
// @Title GetCatalogoById
// @Description get ActaRecibido
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *CatalogoElementosController) GetAll() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := catalogoElementosHelper.GetCatalogoById(id)
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetOne ...
// @Title GetCuentasSubgrupoById
// @Description Devuelve las todas las actas de recibido
// @Success 200 {object} models.Acta_recibido
// @Failure 403
// @router /cuentas_contables/:id [get]
func (c *CatalogoElementosController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := catalogoElementosHelper.GetCuentasContablesSubgrupo(id)
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll2 ...
// @Title GetMovimientosKronos
// @Description get ActaRecibido
// @Success 200 {}
// @Failure 404 not found resource
// @router /movimientos_kronos/
func (c *CatalogoElementosController) GetAll2() {
	v, err := catalogoElementosHelper.GetMovimientosKronos()
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}
