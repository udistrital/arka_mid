package constantes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	e "github.com/udistrital/utils_oas/errorctrl"
)

type opcionCriterioActas struct {
	Handler string // Forma de procesamiento
	Type    string // Tipo de dato
	Map     string // Traducción
}

const (
	// Parametros de filtrado permitidos desde el query
	// Criterio --> parametro en la URL
	ActaNumero       = "Id__contains"
	CreadaDesde      = "FechaCreacion__gte"
	CreadaHasta      = "FechaCreacion__lte"
	ModificadaDesde  = "FechaModificacion__gte"
	ModificadaHasta  = "FechaModificacion__lte"
	ModificadaPor    = "ModificadaPor__icontains"
	AuxiliarAsignado = "ContratistaAsignado__icontains"
	EstadoActa       = "EstadoActaId__in"
	Ubicacion        = "UbicacionId__icontains"
	Observaciones    = "Observaciones__icontains"

	// Otras
	SeparadorEstados = "|"
	FormatoFecha     = time.RFC3339
	FechaCero        = "0001-01-01T00:00:00Z"
)

func OpcionesParametroListaActas(in string, out *opcionCriterioActas) error {
	opciones := map[string]opcionCriterioActas{
		ActaNumero:       {"historico", "int", "ActaRecibidoId__contains"},                       // Consecutivo
		CreadaDesde:      {"historico", "date", "ActaRecibidoId__FechaCreacion__gte"},            // Fecha de creacion, mayor o igual que
		CreadaHasta:      {"historico", "date", "ActaRecibidoId__FechaCreacion__lte"},            // Fecha de creacion, menor o igual que
		ModificadaDesde:  {"historico", "date", "FechaModificacion__gte"},                        // Fecha de modificacion, mayor o igual que
		ModificadaHasta:  {"historico", "date", "FechaModificacion__lte"},                        // Fecha de modificacion, menor o igual que
		ModificadaPor:    {"terceros", "string", "RevisorId__NombreCompleto__icontains"},         // NOTA: Aplicar el "contains" en terceros para traer los id y armar el __in!
		AuxiliarAsignado: {"terceros", "string", "PersonaAsignadaId__NombreCompleto__icontains"}, // NOTA: Aplicar el "contains" en terceros para traer los id y armar el __in!
		EstadoActa:       {"historico", "string_arr", "EstadoActaId__in"},                        // Estado(s, sepadados por |)
		Ubicacion:        {"oikos", "string", "EspacioFisicoId__Nombre__icontains"},              // NOTA: Aplicar el "contains" en oikos para traer los id y armar el __in!
		Observaciones:    {"historico", "string", "Observaciones__icontains"},                    //
	}
	ok := false
	if *out, ok = opciones[in]; ok {
		return nil
	}
	return fmt.Errorf("'%s' is not an allowed param", in)
}

func ParametroValidoListaActas(parametro string) bool {
	var config opcionCriterioActas
	err := OpcionesParametroListaActas(parametro, &config)
	return err == nil
}

func ParserParametrosListaActas(in map[string]string, out map[string]interface{}) (outputError map[string]interface{}) {
	const funcion = "ParserParametrosListaActas - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	for k, v := range in {
		var clean interface{}
		if err := ParseParametro(k, v, &clean); err != nil {
			return err
		}
		out[k] = clean
	}
	return
}

func ParseParametro(inK, inV string, out *interface{}) (outputError map[string]interface{}) {
	const funcion = "ParseParametro"
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	var (
		opciones opcionCriterioActas
		err      error
	)
	if err = OpcionesParametroListaActas(inK, &opciones); err != nil {
		return e.Error(funcion+"OpcionesParametroListaActas(k, &opciones)",
			err, fmt.Sprint(http.StatusBadRequest))
	}
	switch opciones.Type {
	case "string":
		*out = inV
	case "string_arr":
		*out = strings.Split(inV, SeparadorEstados)
	case "int":
		var aux int
		if aux, err = strconv.Atoi(inV); err != nil {
			return e.Error(funcion+"strconv.Atoi(v)",
				err, fmt.Sprint(http.StatusBadRequest))
		}
		*out = aux
	case "date":
		var t time.Time
		if t, err = time.Parse(FormatoFecha, inV); err != nil {
			return e.Error(funcion+"time.Parse(FormatoFecha, v)",
				err, fmt.Sprint(http.StatusBadRequest))
		}
		*out = t
	default:
		err = fmt.Errorf("'%s' type processing not implemented", opciones.Type)
		logs.Error(err)
		return e.Error(funcion+"unhandled param",
			err, fmt.Sprint(http.StatusNotImplemented))
	}
	return
}
