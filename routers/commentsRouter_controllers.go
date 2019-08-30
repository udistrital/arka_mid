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
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetAllActasRecibido",
            Router: `/get_actas_recibido/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/arka_mid/controllers:ActaRecibidoController"],
        beego.ControllerComments{
            Method: "GetActasRecibidoTipo",
            Router: `/get_actas_recibido_tipo/:id`,
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

}
