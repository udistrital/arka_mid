package depreciacionHelper

import (
	"time"
)

// GetDeltaTiempo retorna el tiempo en años entre dos fechas
func GetDeltaTiempo(ref, fin time.Time) (prct float64) {

	ref = time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, time.UTC)
	fin = time.Date(fin.Year(), fin.Month(), fin.Day(), 0, 0, 0, 0, time.UTC)

	prct = fin.Sub(ref).Hours() / (24 * 365)

	return
}

// CalculaDp Genera el valor y el tiempo en años a depreciar
// ref: Fecha de referencia para determinar el tiempo por el cual se correra la depreciacion.
// dp: Valor calculado correspondiente a la depreciacion.
func CalculaDp(presente, residual, vUtil float64, ref, fCorte time.Time) (dp, deltaT float64) {

	if vUtil == 0 {
		return 0, 0
	}

	deltaT = GetDeltaTiempo(ref, fCorte.AddDate(0, 0, 1))
	if deltaT > vUtil {
		deltaT = vUtil
	}

	if residual > presente {
		presente = residual
	}

	dp = (presente - residual) * deltaT / vUtil
	return dp, deltaT

}

func dscTransaccionCierre() string {
	return "Mediciones posteriores almacén"
}
