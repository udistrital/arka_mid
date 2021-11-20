package trasladoshelper

// GetDetalle Consulta los funcionarios, ubicación y elementos asociados a un traslado
func GetDetalleTraslado(id int) (Elemento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDetalleTraslado", "err": err, "status": "502"}
			panic(outputError)
		}
	}()
	return nil, nil
}
