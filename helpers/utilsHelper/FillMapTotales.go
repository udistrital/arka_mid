package utilsHelper

func FillMapTotales(totales map[int]float64, subgrupo int, valor float64) {
	totales[subgrupo] += valor
}
