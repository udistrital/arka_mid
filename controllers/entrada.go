package controllers

import (
	"errors"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// EntradaController operations for Entrada
type EntradaController struct {
	beego.Controller
}

// URLMapping ...
func (c *EntradaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("PutAnular", c.PutAnular)
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
	defer errorCtrl.ErrorControlController(c.Controller, "EntradaController")

	logs.Info("==== INICIO EntradaController.Post ====")

	entradaId, errEntradaID := c.GetInt("entradaId", 0)
	etl, errETL := c.GetBool("etl", false)
	aprobar, errAprobar := c.GetBool("aprobar", false)

	logs.Info("Query params recibidos -> entradaId: %d, etl: %v, aprobar: %v", entradaId, etl, aprobar)
	logs.Info("Errores parseando query params -> entradaId: %v, etl: %v, aprobar: %v", errEntradaID, errETL, errAprobar)
	logs.Info("Body crudo recibido -> %s", string(c.Ctx.Input.RequestBody))

	if aprobar && entradaId > 0 {
		logs.Info("Entró en rama de APROBACIÓN. entradaId=%d", entradaId)

		var res models.ResultadoMovimiento
		if err := entradaHelper.AprobarEntrada(entradaId, &res); err != nil {
			logs.Error("Error en entradaHelper.AprobarEntrada(%d): %v", entradaId, err)
			panic(err)
		}

		logs.Info("Aprobación exitosa. Respuesta: %+v", res)
		c.Data["json"] = res

	} else if !aprobar {
		logs.Info("Entró en rama de REGISTRO/ACTUALIZACIÓN. aprobar=%v, entradaId=%d", aprobar, entradaId)

		var (
			v       models.TransaccionEntrada
			entrada models.ResultadoMovimiento
		)

		if err := utilsHelper.Unmarshal(string(c.Ctx.Input.RequestBody), &v); err != nil {
			logs.Error("Error deserializando body en models.TransaccionEntrada: %v", err)
			logs.Error("Body que falló: %s", string(c.Ctx.Input.RequestBody))
			panic(map[string]interface{}{
				"funcion": "Post - utilsHelper.Unmarshal(RequestBody, &v)",
				"err":     err,
				"status":  "400",
			})
		}

		logs.Info("Body deserializado correctamente: %+v", v)

		if entradaId > 0 {
			logs.Info("Entró en UPDATE de entrada. entradaId=%d", entradaId)

			if err := entradaHelper.UpdateEntrada(&v, entradaId, &entrada); err != nil {
				logs.Error("Error en entradaHelper.UpdateEntrada(&v, %d, &entrada): %v", entradaId, err)
				logs.Error("Payload usado para update: %+v", v)
				panic(err)
			}

			logs.Info("Update exitoso. Respuesta: %+v", entrada)

		} else if entradaId == 0 {
			logs.Info("Entró en REGISTRO de entrada nueva. etl=%v", etl)

			if err := entradaHelper.RegistrarEntrada(&v, etl, &entrada); err != nil {
				logs.Error("Error en entradaHelper.RegistrarEntrada(&v, %v, &entrada): %v", etl, err)
				logs.Error("Payload usado para registro: %+v", v)
				panic(err)
			}

			logs.Info("Registro exitoso. Respuesta: %+v", entrada)
		}

		c.Data["json"] = entrada
	} else {
		logs.Error("Caso no contemplado en Post. entradaId=%d, aprobar=%v, etl=%v", entradaId, aprobar, etl)
		panic(map[string]interface{}{
			"funcion": "Post - validación de flujo",
			"err":     errors.New("combinación de parámetros no válida"),
			"status":  "400",
		})
	}

	logs.Info("Respuesta final Post: %+v", c.Data["json"])
	logs.Info("==== FIN EntradaController.Post ====")
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Detalle de entrada por Id. Retorna la transaccion contable si la entrada ya fue aprobada
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DetalleEntrada
// @Failure 403 :id is empty
// @router /:id [get]
func (c *EntradaController) GetOne() {
	defer errorCtrl.ErrorControlController(c.Controller, "EntradaController")

	logs.Info("==== INICIO EntradaController.GetOne ====")
	logs.Info("Path param recibido :id -> %s", c.Ctx.Input.Param(":id"))

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una entrada válida")
		}
		logs.Error("Error obteniendo/parsing :id. Valor recibido=%s, error=%v", c.Ctx.Input.Param(":id"), err)
		panic(map[string]interface{}{
			"funcion": `GetOne - c.GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	logs.Info("Consultando detalle de entrada con id=%d", id)

	respuesta, err := entradaHelper.DetalleEntrada(id)
	if err == nil || respuesta != nil {
		logs.Info("DetalleEntrada respuesta -> err: %v, respuesta: %+v", err, respuesta)
		c.Data["json"] = respuesta
	} else {
		logs.Error("Error en entradaHelper.DetalleEntrada(%d): %v", id, err)

		if err != nil {
			panic(err)
		}

		panic(map[string]interface{}{
			"funcion": "GetOne - entradaHelper.DetalleEntrada(id)",
			"err":     errors.New("no se obtuvo respuesta al consultar la entrada"),
			"status":  "404",
		})
	}

	logs.Info("Respuesta final GetOne: %+v", c.Data["json"])
	logs.Info("==== FIN EntradaController.GetOne ====")
	c.ServeJSON()
}

// PutAnular ...
// @Title PutAnular
// @Description Anula una entrada, registra un movimiento de reversión contable y devuelve el acta a estado en verificación.
// @Param	id		path 	int								true	"Id de la entrada a anular"
// @Param	body	body	models.AnulacionEntradaRequest	false	"Observación de la anulación"
// @Success 200 {object} models.ResultadoAnulacionEntrada
// @Failure 400 the request contains incorrect syntax
// @router /:id/anular [put]
func (c *EntradaController) PutAnular() {
	defer errorCtrl.ErrorControlController(c.Controller, "EntradaController")

	var id int
	if v, err := c.GetInt(":id"); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("se debe especificar una entrada válida")
		}
		panic(map[string]interface{}{
			"funcion": `PutAnular - c.GetInt(":id")`,
			"err":     err,
			"status":  "400",
		})
	} else {
		id = v
	}

	var request models.AnulacionEntradaRequest
	if body := strings.TrimSpace(string(c.Ctx.Input.RequestBody)); body != "" {
		if err := utilsHelper.Unmarshal(body, &request); err != nil {
			panic(map[string]interface{}{
				"funcion": "PutAnular - utilsHelper.Unmarshal(RequestBody, &request)",
				"err":     err,
				"status":  "400",
			})
		}
	}

	var resultado models.ResultadoAnulacionEntrada
	if err := entradaHelper.AnularEntrada(id, &request, &resultado); err != nil {
		panic(err)
	}

	c.Data["json"] = resultado
	c.ServeJSON()
}
