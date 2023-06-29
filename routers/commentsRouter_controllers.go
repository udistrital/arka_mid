package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
		beego.ControllerComments{
			Method:           "GetParametros",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
		beego.ControllerComments{
			Method:           "GetElementosActa",
			Router:           "/elementos/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
		beego.ControllerComments{
			Method:           "GetAllActas",
			Router:           "/get_all_actas/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"],
		beego.ControllerComments{
			Method:           "PostManual",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"],
		beego.ControllerComments{
			Method:           "GetOneManual",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"],
		beego.ControllerComments{
			Method:           "PostAjuste",
			Router:           "/automatico",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"],
		beego.ControllerComments{
			Method:           "GetOneAuto",
			Router:           "/automatico/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AjusteController"],
		beego.ControllerComments{
			Method:           "GetElementos",
			Router:           "/automatico/elementos/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:AvaluoController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           "/:id",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
		beego.ControllerComments{
			Method:           "GetSolicitud",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
		beego.ControllerComments{
			Method:           "PutRevision",
			Router:           "/aprobar",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
		beego.ControllerComments{
			Method:           "GetDetalleElemento",
			Router:           "/elemento/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
		beego.ControllerComments{
			Method:           "GetAperturasKardex",
			Router:           "/aperturas_kardex/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
		beego.ControllerComments{
			Method:           "GetElementos",
			Router:           "/elementos_sin_asignar/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
		beego.ControllerComments{
			Method:           "GetAllExistencias",
			Router:           "/existencias_kardex/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/solicitud/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
		beego.ControllerComments{
			Method:           "GetAllSolicitud",
			Router:           "/solicitud/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
		beego.ControllerComments{
			Method:           "GetOneSolicitud",
			Router:           "/solicitud/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/cuentas_contables/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:DepreciacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:DepreciacionController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:DepreciacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:DepreciacionController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:DepreciacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:DepreciacionController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
		beego.ControllerComments{
			Method:           "GetMovimientos",
			Router:           "/movimientos/:acta_recibido_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:InmueblesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:InmueblesController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:InmueblesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:InmueblesController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:InmueblesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:InmueblesController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:InmueblesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:InmueblesController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:PolizasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:PolizasController"],
		beego.ControllerComments{
			Method:           "GetAllElementosPoliza",
			Router:           "/AllElementosPoliza",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
		beego.ControllerComments{
			Method:           "GetSalidas",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
		beego.ControllerComments{
			Method:           "GetSalida",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
		beego.ControllerComments{
			Method:           "GetElementos",
			Router:           "/elementos",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
		beego.ControllerComments{
			Method:           "GetTraslado",
			Router:           "/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
		beego.ControllerComments{
			Method:           "GetElementosFuncionario",
			Router:           "/funcionario/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
