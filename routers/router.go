// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"

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
		beego.NSNamespace("/depreciacion",
			beego.NSInclude(
				&controllers.DepreciacionController{},
			),
		),
		beego.NSNamespace("/ajustes",
			beego.NSInclude(
				&controllers.AjusteController{},
			),
		),
		beego.NSNamespace("/inmuebles",
			beego.NSInclude(
				&controllers.InmueblesController{},
			),
		),
		beego.NSNamespace("/avaluo",
			beego.NSInclude(
				&controllers.AvaluoController{},
			),
		),
	)

	beego.AddNamespace(ns)
}
