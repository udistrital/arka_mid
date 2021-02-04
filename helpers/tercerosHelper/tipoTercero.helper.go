package tercerosHelper

var diccionarioTipoHelper = map[string](func() ([]map[string]interface{}, map[string]interface{})){
	"funcionarioPlanta": GetFuncionariosPlanta,
}

func GetTipos() (tercero []string) {
	for k := range diccionarioTipoHelper {
		tercero = append(tercero, k)
	}
	return tercero
}

func GetHelperTipo(tipo string) (helper func() ([]map[string]interface{}, map[string]interface{}), outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/CertificacionDocumentosAprobados",
				"err":     err,
				"status":  "400",
			}
			panic(outputError)
		}
	}()

	return diccionarioTipoHelper[tipo], nil
}
