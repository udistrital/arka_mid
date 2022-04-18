package utilsHelper

func FillMapTotales(totales map[int]float64, subgrupo int, valor float64) {
	if _, ok := totales[subgrupo]; ok {
		totales[subgrupo] += valor
	} else {
		totales[subgrupo] = valor
	}
}
