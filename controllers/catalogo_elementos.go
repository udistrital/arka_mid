package controllers

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// CatalogoElementosController operations for CatalogoElementos
type CatalogoElementosController struct {
	beego.Controller
}

// URLMapping ...
func (c *CatalogoElementosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetCatalogoById", c.GetCatalogoById)
}

// GetCatalogoById ...
// @Title GetCatalogoById
// @Description Devuelve el catalogo de elementos
// @Success 200 {object} models.CatalogoElementos
// @Failure 403
// @router /:id  [get]
func (c *CatalogoElementosController) GetCatalogoById() {
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

// GetCuentasSubgrupoById ...
// @Title GetCuentasSubgrupoById
// @Description Devuelve las cuentas contables asosciadas a un subgrupo
// @Success 200 {object} models.CuentasSubgrupo
// @Failure 403
// @router /cuentas_contables/:id  [get]
func (c *CatalogoElementosController) GetCuentasSubgrupoById() {
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
