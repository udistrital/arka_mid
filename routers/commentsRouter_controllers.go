package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetActasByTipo",
            Router: `/get_actas_recibido_tipo/:tipo`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetElementosActa",
            Router: `/get_elementos_acta/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetSoportesActa",
            Router: `/get_soportes_acta/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"],
        beego.ControllerComments{
            Method: "GetCatalogoById",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:CatalogoElementosController"],
        beego.ControllerComments{
            Method: "GetCuentasSubgrupoById",
            Router: `/cuentas_contables/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:EntradaController"],
        beego.ControllerComments{
            Method: "GetEntrada",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ParametrosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ParametrosController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
