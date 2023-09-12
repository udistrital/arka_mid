package consecutivos

import (
	"fmt"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// Genera un consecutivo con el año actual y para un contextoId determinado
func Get(contexto string, descripcion string, data *models.Consecutivo) (outputError map[string]interface{}) {

	funcion := "Get - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	contextoId, err := beego.AppConfig.Int(contexto)
	if err != nil {
		eval := "beego.AppConfig.Int(contexto)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	year := time.Now().Year()
	*data = models.Consecutivo{
		ContextoId:  contextoId,
		Year:        year,
		Descripcion: descripcion,
		Activo:      true,
	}

	if err := Post(&data); err != nil {
		return err
	}

	return

}

// Le da formato a un consecutivo, para un prefijo indicado, un formato determinado para el número del consecutivo. Se toma el año como el sufijo.
func Format(format, prefix string, consecutivo *models.Consecutivo) (consFormat string) {

	funcion := "Format"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	consecutivo_ := fmt.Sprintf(format, consecutivo.Consecutivo)
	suffix := fmt.Sprintf("%04d", consecutivo.Year)
	return prefix + "-" + consecutivo_ + "-" + suffix
}
