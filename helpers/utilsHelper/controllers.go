// Agregar utilidades en relacion a controladores

package utilsHelper

import (
	"errors"
	"strings"
)

const (
	DefaultQueryKeySeparator  = ":"
	DefaultQueryTermSeparator = ","
)

// Separa los elementos de un query serializado en un string
// y los guarda en un objeto query especificado
//
// Por ejemplo
//
//   const input = "t1:v1,t2:v2"
//   var output map[string]string
//   if err := utilsHelper.CustomQuerySplit(input); err != nil {
//     panic(err)
//   }
//
// Despues de lo anterior, si se serializara en JSON se obtendría:
// {"t1":"v1","t2":"v2"}
func CustomQuerySplit(raw, termSeparator, keySeparator string,
	query *map[string]string) error {
	if termSeparator == "" {
		termSeparator = DefaultQueryTermSeparator
	}
	if keySeparator == "" {
		keySeparator = DefaultQueryKeySeparator
	}
	if termSeparator == keySeparator {
		return errors.New("termSeparator is the same keySeparator")
	}
	// Mayormente basado en el codigo generado por Beego para el query
	for _, cond := range strings.Split(raw, termSeparator) {
		kv := strings.SplitN(cond, keySeparator, 2)
		if len(kv) != 2 {
			return errors.New("invalid query key/value pair")
		}
		k, v := kv[0], kv[1]
		(*query)[k] = v
	}
	return nil
}

func QuerySplit(raw string, query *map[string]string) error {
	return CustomQuerySplit(raw, "", "", query)
}
