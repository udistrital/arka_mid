// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/acta_recibido",
			beego.NSInclude(
				&controllers.ActaRecibidoController{},
			),
		),
		beego.NSNamespace("/entrada",
			beego.NSInclude(
				&controllers.EntradaController{},
			),
		),
		beego.NSNamespace("/elemento",
			beego.NSInclude(
				&controllers.ElementoController{},
			),
		),
		beego.NSNamespace("/parametros_soporte",
			beego.NSInclude(
				&controllers.ParametrosController{},
			),
		),
		beego.NSNamespace("/catalogo_elementos",
			beego.NSInclude(
				&controllers.CatalogoElementosController{},
			),
		),
		beego.NSNamespace("/salida",
			beego.NSInclude(
				&controllers.SalidaController{},
			),
		),
		beego.NSNamespace("/terceros",
			beego.NSInclude(
				&controllers.TercerosController{},
			),
		),
		beego.NSNamespace("/bodega_consumo",
			beego.NSInclude(
				&controllers.BodegaConsumoController{},
			),
		),
		beego.NSNamespace("/bajas_elementos",
			beego.NSInclude(
				&controllers.BajaController{},
			),
		),
		beego.NSNamespace("/traslados",
			beego.NSInclude(
				&controllers.TrasladosController{},
			),
		),
		beego.NSNamespace("/polizas",
			beego.NSInclude(
				&controllers.PolizasController{},
			),
		),
	)

	beego.AddNamespace(ns)
}
