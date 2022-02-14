package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetAllElementosConsumo",
            Router: "/elementosconsumo/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetAllActas",
            Router: "/get_all_actas/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetElementosActa",
            Router: "/get_elementos_acta/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
        beego.ControllerComments{
            Method: "PutRevision",
            Router: "/aprobar",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
        beego.ControllerComments{
            Method: "GetElemento",
            Router: "/elemento/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
        beego.ControllerComments{
            Method: "GetDetalleElemento",
            Router: "/elemento/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/solicitud/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BajaController"],
        beego.ControllerComments{
            Method: "GetSolicitud",
            Router: "/solicitud/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
        beego.ControllerComments{
            Method: "GetAperturasKardex",
            Router: "/aperturas_kardex/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
        beego.ControllerComments{
            Method: "GetElementos",
            Router: "/elementos_sin_asignar/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
        beego.ControllerComments{
            Method: "GetAllExistencias",
            Router: "/existencias_kardex/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:BodegaConsumoController"],
        beego.ControllerComments{
            Method: "GetOneSolicitud",
            Router: "/solicitud/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/cuentas_contables/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"],
        beego.ControllerComments{
            Method: "GetAll2",
            Router: "/movimientos_kronos/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ElementoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ElementoController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
        beego.ControllerComments{
            Method: "AnularEntrada",
            Router: "/anular/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
        beego.ControllerComments{
            Method: "GetEncargadoElemento",
            Router: "/encargado/:placa",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
        beego.ControllerComments{
            Method: "GetMovimientos",
            Router: "/movimientos/:acta_recibido_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ParametrosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ParametrosController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ParametrosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ParametrosController"],
        beego.ControllerComments{
            Method: "PostAsignacionEspacioDependencia",
            Router: "/post_asignacion_espacio_fisico_dependencia/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:PolizasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:PolizasController"],
        beego.ControllerComments{
            Method: "GetAllElementosPoliza",
            Router: "/AllElementosPoliza",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
        beego.ControllerComments{
            Method: "GetSalidas",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
        beego.ControllerComments{
            Method: "GetSalida",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:SalidaController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"],
        beego.ControllerComments{
            Method: "GetTipos",
            Router: "/tipo/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"],
        beego.ControllerComments{
            Method: "GetByTipo",
            Router: "/tipo/:tipo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TercerosController"],
        beego.ControllerComments{
            Method: "GetByTipoAndID",
            Router: "/tipo/:tipo/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
        beego.ControllerComments{
            Method: "GetTraslado",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:TrasladosController"],
        beego.ControllerComments{
            Method: "GetElementosFuncionario",
            Router: "/funcionario/:funcionarioId",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
