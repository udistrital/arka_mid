package movimientosArkaHelper

// GetFormatoTipoMovimientoIdByCodigoAbreviacion Consulta el Id de un FormatoTipoMovimiento según el Codigo de abrebiación del mismo
func GetFormatoTipoMovimientoIdByCodigoAbreviacion(id *int, codigoAbreviacion string) (outputError map[string]interface{}) {

	query := "query=CodigoAbreviacion:" + codigoAbreviacion
	if fm, err := GetAllFormatoTipoMovimiento(query); err != nil {
		return err
	} else {
		*id = fm[0].Id
	}

	return
}
