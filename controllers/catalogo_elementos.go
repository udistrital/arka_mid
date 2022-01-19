package controllers

import (
	//"github.com/udistrital/acta_recibido_crud/models"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/models"
)

// CatalogoElementosController operations for Catalogo
type CatalogoElementosController struct {
	beego.Controller
}

// URLMapping ...
func (c *CatalogoElementosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.GetAll)
	c.Mapping("GetOne", c.GetOne)
}

// GetAll ...
// @Title GetCatalogoById
// @Description get ActaRecibido
// @Success 200 {}
// @Failure 404 not found resource
// @router /:id [get]
func (c *CatalogoElementosController) GetAll() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "CatalogoElementosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		if err == nil {
			err = fmt.Errorf("id must be > 0")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetAll - strconv.Atoi(idStr)",
			"err":     err,
			"status":  "400",
		})
	}

	if v, err := catalogoElementosHelper.GetCatalogoById(id); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetOne ...
// @Title GetCuentasSubgrupoById
// @Description Devuelve las todas las actas de recibido
// @Param	id		path 	int	true		"subgroup id"
// @Success 200 {object} models.Acta_recibido
// @Failure 403
// @router /cuentas_contables/:id [get]
func (c *CatalogoElementosController) GetOne() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "CatalogoElementosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	idStr := c.Ctx.Input.Param(":id")
	var id int
	if v, err := strconv.Atoi(idStr); err == nil && v > 0 {
		id = v
	} else {
		if err == nil {
			err = fmt.Errorf("id MUST be > 0")
		}
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "GetOne - strconv.Atoi(idStr)",
			"err":     err,
			"status":  "400",
		})
	}

	if v, err := catalogoElementosHelper.GetCuentasContablesSubgrupo(id); err != nil {
		panic(err)
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// Post ...
// @Title Create
// @Description Tr_cuentas_subgrupo su función es permitir una consulta de n cuentas asociadas a subgrupos
// @Param	body		body 	models.Tr_cuentas_subgrupo	true		"body for Tr_cuentas_subgrupo content (Recibe un arreglo de subgrupos, cada subgrupo debe llevar un Id definido)"
// @Success 201 {object} models.Tr_cuentas_subgrupo
// @Failure 403 body is empty
// @router / [post]
func (c *CatalogoElementosController) Post() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "CatalogoElementosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	//esto deberia ser un get ya que es una consulta y recibir de a un id
	//var arreglosubgrupos []models.Subgrupo
	var arreglosubgrupos []models.SubgrupoCuentasModelo
	fmt.Println("entra al post")
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &arreglosubgrupos); err == nil {
		if data, err := catalogoElementosHelper.GetTipoMovimiento(arreglosubgrupos); err == nil {
			c.Data["json"] = data
		} else {
			panic(err)
		}
		c.ServeJSON()
	} else {
		logs.Error(err)
		panic(map[string]interface{}{
			"funcion": "Post - json.Unmarshal(c.Ctx.Input.RequestBody, &arreglosubgrupos)",
			"err":     err,
			"status":  "400",
		})
	}

}
