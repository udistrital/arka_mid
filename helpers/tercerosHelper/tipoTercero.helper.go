package tercerosHelper

var diccionarioTipoHelper = map[string](func(int) ([]map[string]interface{}, map[string]interface{})){
	"funcionarioPlanta": GetFuncionariosPlanta,
	"jefes":             GetJefes,
}

// GetTipos retorna la lista de tipos que pueden ser usados con GetHelperTipo
func GetTipos() (tercero []string) {
	for k := range diccionarioTipoHelper {
		tercero = append(tercero, k)
	}
	return tercero
}

// GetHelperTipo trae los terceros con el criterio especificado.
// El criterio debe ser alguno de los valores retornados por GetTipos
func GetHelperTipo(tipo string) (helper func(int) ([]map[string]interface{}, map[string]interface{}), outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetHelperTipo",
				"err":     err,
				"status":  "400",
			}
			panic(outputError)
		}
	}()

	return diccionarioTipoHelper[tipo], nil
}
