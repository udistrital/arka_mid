package models

type TransaccionActaRecibido struct {
	ActaRecibido *ActaRecibido
	UltimoEstado *HistoricoActa
	SoportesActa *[]SoporteActa
	Elementos    []*Elemento
}
