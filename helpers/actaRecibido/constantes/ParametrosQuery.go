package constantes

import (
	"fmt"
)

type opcionCriterioActas struct {
	Type string // Forma de procesamiento
	Map  string // Traducción
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
)

func OpcionesParametroListaActas(in string, out *opcionCriterioActas) error {
	opciones := map[string]opcionCriterioActas{
		ActaNumero:       {"historico", "ActaRecibidoId__contains"},                    // Consecutivo
		CreadaDesde:      {"historico", "ActaRecibidoId__FechaCreacion__gte"},          // Fecha de creacion, mayor o igual que
		CreadaHasta:      {"historico", "ActaRecibidoId__FechaCreacion__lte"},          // Fecha de creacion, menor o igual que
		ModificadaDesde:  {"historico", "FechaModificacion__gte"},                      // Fecha de modificacion, mayor o igual que
		ModificadaHasta:  {"historico", "FechaModificacion__lte"},                      // Fecha de modificacion, menor o igual que
		ModificadaPor:    {"terceros", "RevisorId__NombreCompleto__icontains"},         // NOTA: Aplicar el "contains" en terceros para traer los id y armar el __in!
		AuxiliarAsignado: {"terceros", "PersonaAsignadaId__NombreCompleto__icontains"}, // NOTA: Aplicar el "contains" en terceros para traer los id y armar el __in!
		EstadoActa:       {"historico", "EstadoActaId__in"},                            // Estado(s, sepadados por |)
		Ubicacion:        {"oikos", "EspacioFisicoId__Nombre__icontains"},              // NOTA: Aplicar el "contains" en oikos para traer los id y armar el __in!
		Observaciones:    {"historico", "Observaciones__icontains"},                    //
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
