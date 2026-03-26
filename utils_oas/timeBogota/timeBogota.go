package timebogota

import (
	"time"
)

var tiempoBogota time.Time

func TiempoBogota() time.Time {

	tiempoBogota = time.Now()
	loc, _ := time.LoadLocation("America/Bogota")
	tiempoBogota = tiempoBogota.In(loc)

	return tiempoBogota
}
