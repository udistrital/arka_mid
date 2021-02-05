package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
)

// TercerosController operations for Terceros
type TercerosController struct {
	beego.Controller
}

// URLMapping ...
func (c *TercerosController) URLMapping() {
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetByTipo", c.GetByTipo)
	c.Mapping("GetTipos", c.GetTipos)
	c.Mapping("GetByTipoAndId", c.GetByTipoAndID)
}

// GetOne ...
// @Title GetOne
// @Description get Terceros by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Terceros
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TercerosController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	//id, _ := strconv.Atoi(idStr)
	v, err := tercerosHelper.GetNombreTerceroById(idStr)
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Terceros
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Terceros
// @Failure 403
// @router / [get]
/*
func (c *TercerosController) GetAll() {

}
*/

// GetTipos ...
// @Title GetAll
// @Description List the Tercero types that can be used to gather Terceros by {tipo}
// @Success 200 {object} []string
// @Failure 403
// @router /tipo/ [get]
func (c *TercerosController) GetTipos() {
	c.Data["json"] = tercerosHelper.GetTipos()
	c.ServeJSON()
}

// GetByTipo ...
// @Title GetAll
// @Description get Terceros with the specified {tipo}
// @Param	tipo	path 	string	true		"Tercero type available from /tipo/"
// @Success 200 {object} []models.Terceros
// @Failure 403
// @router /tipo/:tipo [get]
func (c *TercerosController) GetByTipo() {
	tipo := c.Ctx.Input.Param(":tipo")

	if helper, err := tercerosHelper.GetHelperTipo(tipo); err == nil {
		if v, err := helper(0); err == nil {
			c.Data["json"] = v
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GetByTipoAndId ...
// @Title GetAll
// @Description get Terceros with the specified {tipo} and {id} of a record in terceros table from Terceros CRUD API
// @Param	tipo	path 	string	true		"Tercero type available from /tipo/"
// @Param	id		path 	uint	true		"ID. MUST be greater than 0"
// @Success 200 {object} []models.Terceros
// @Failure 403
// @router /tipo/:tipo/:id [get]
func (c *TercerosController) GetByTipoAndID() {

	tipo := c.Ctx.Input.Param(":tipo")
	idQuery := c.Ctx.Input.Param(":id")
	var id int
	if i, err := strconv.Atoi(idQuery); err == nil && i > 0 {
		id = i
	} else if err != nil {
		panic(err)
	} else {
		err := fmt.Errorf("Wrong ID")
		panic(err)
	}

	if helper, err := tercerosHelper.GetHelperTipo(tipo); err == nil {
		if v, err := helper(id); err == nil {
			c.Data["json"] = v
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Terceros
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Terceros	true		"body for Terceros content"
// @Success 200 {object} models.Terceros
// @Failure 403 :id is not int
// @router /:id [put]
/*
func (c *TercerosController) Put() {

}
*/

// Delete ...
// @Title Delete
// @Description delete the Terceros
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
/*
func (c *TercerosController) Delete() {

}
*/
