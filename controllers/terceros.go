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
	c.Mapping("GetTipos", c.GetTipos)
	c.Mapping("GetByTipo", c.GetByTipo)
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

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "TercerosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500")
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")

	if v, err := tercerosHelper.GetNombreTerceroById(idStr); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetTipos ...
// @Title GetAll
// @Description List the Tercero types that can be used to gather Terceros by {tipo}
// @Success 200 {object} []string
// @router /tipo/ [get]
func (c *TercerosController) GetTipos() {

	// Puede que ni sea necesario en este controlador, pero se coloca por lineamiento...
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "TercerosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	if v, err := tercerosHelper.GetTipos(); err == nil {
		c.Data["json"] = v
		c.ServeJSON()
	} else {
		panic(err)
	}
}

// GetByTipo ...
// @Title GetAll
// @Description get Terceros with the specified {tipo}
// @Param	tipo	path 	string	true		"Tercero type available from /tipo/"
// @Success 200 {object} []map[string]interface{}
// @Failure 500 Internal Error
// @Failure 501 {tipo} Not Implemented
// @Failure 502 Error with external API
// @router /tipo/:tipo [get]
func (c *TercerosController) GetByTipo() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "TercerosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

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
// @Success 200 {object} []map[string]interface{}
// @Failure 400 Wrong ID
// @Failure 404 ID with {tipo} Not Found
// @Failure 500 Internal Error
// @Failure 501 {tipo} Not Implemented
// @Failure 502 Error with external API
// @router /tipo/:tipo/:id [get]
func (c *TercerosController) GetByTipoAndID() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "TercerosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	tipo := c.Ctx.Input.Param(":tipo")
	idQuery := c.Ctx.Input.Param(":id")
	var id int
	if i, err := strconv.Atoi(idQuery); err == nil && i > 0 {
		id = i
	} else {
		if err == nil {
			err = fmt.Errorf("ID MUST be greater than 0 - Got: %d", i)
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetByTipoAndID - strconv.Atoi(idQuery)",
			"err":     err,
			"status":  "400",
		})
	}

	if helper, err := tercerosHelper.GetHelperTipo(tipo); err == nil {
		if v, err := helper(id); err == nil {
			if len(v) == 0 {
				err := fmt.Errorf("no se encontró un Tercero tipo '%s' con id '%d'", tipo, id)
				panic(map[string]interface{}{
					"funcion": "GetByTipoAndID - len(v) == 0",
					"err":     err,
					"status":  "404",
				})
			}
			c.Data["json"] = v
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
