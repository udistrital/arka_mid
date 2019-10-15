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
	)

	beego.AddNamespace(ns)
}
