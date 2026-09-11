package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
)

// AuditoriaController maneja eventos funcionales que deben quedar
// registrados explícitamente en los logs de arka_mid.
type AuditoriaController struct {
	beego.Controller
}

// URLMapping ...
func (c *AuditoriaController) URLMapping() {
	c.Mapping("PostEvento", c.PostEvento)
}

// PostEvento ...
// @Title Registrar evento de auditoría
// @Description Registra en los logs de arka_mid una operación funcional relevante. El usuario es obtenido desde el token Authorization por utils_oas/v2.
// @Param body body models.EventoAuditoria true "Evento de auditoría"
// @Success 204
// @Failure 400 solicitud inválida
// @router /evento [post]
func (c *AuditoriaController) PostEvento() {
	var evento models.EventoAuditoria

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &evento); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Message": "El cuerpo del evento de auditoría no es válido",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(evento.Evento) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Message": "El campo evento es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(evento.RequestOriginal.Method) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Message": "El método de la petición original es obligatorio",
		}
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(evento.RequestOriginal.Endpoint) == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Message": "El endpoint de la petición original es obligatorio",
		}
		c.ServeJSON()
		return
	}

	/*
		El usuario se obtiene desde Authorization por el middleware.
	*/

	c.Ctx.Input.SetData("evento", evento.Evento)

	c.Ctx.Input.SetData(
		"request_original",
		evento.RequestOriginal,
	)

	c.Ctx.Input.SetData(
		"resultado_original",
		evento.ResultadoOriginal,
	)

	c.Ctx.Output.SetStatus(http.StatusNoContent)
}
