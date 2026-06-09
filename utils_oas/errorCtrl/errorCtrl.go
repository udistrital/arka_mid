package errorCtrl

import (
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// Error control for controller
func ErrorControlController(c web.Controller, controller string) {
	if err := recover(); err != nil {
		logs.Error(err)
		localError, ok := err.(map[string]interface{})
		if !ok {
			c.Abort(strconv.Itoa(http.StatusInternalServerError))
			return
		}
		appName, _ := web.AppConfig.String("appname")
		c.Data["mesaage"] = (appName + "/" + controller + "/" + (localError["funcion"]).(string))
		c.Data["data"] = normalizeErrorValue(localError["err"])
		if status, ok := localError["status"]; ok && status != nil {
			c.Abort(status.(string))
		} else {
			c.Abort(strconv.Itoa(http.StatusInternalServerError)) // Unhandled Error!
		}
	}
}

// Error control for functions
func ErrorControlFunction(funcion string, status string) {
	if err := recover(); err != nil {
		panic(Error(funcion, err, status))
	}
}

// Get a error with standard struct
func Error(funcion string, err interface{}, status string) (outputError map[string]interface{}) {
	switch localError := err.(type) {
	case map[string]interface{}:
		if fun, ok := localError["funcion"]; ok && fun != nil {
			funcion = funcion + "/" + fun.(string)
		}
		if er, ok := localError["err"]; ok && er != nil {
			err = er
		}
	}
	outputError = map[string]interface{}{"funcion": funcion, "err": normalizeErrorValue(err), "status": status}
	return outputError
}

func normalizeErrorValue(err interface{}) interface{} {
	switch val := err.(type) {
	case nil:
		return nil
	case error:
		return val.Error()
	case map[string]interface{}:
		normalized := make(map[string]interface{}, len(val))
		for key, item := range val {
			normalized[key] = normalizeErrorValue(item)
		}
		return normalized
	default:
		return err
	}
}
