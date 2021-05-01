package models

type TransaccionActaRecibido struct {
	ActaRecibido *ActaRecibido
	UltimoEstado *HistoricoActa
	SoportesActa []*TransaccionSoporteActa
}
