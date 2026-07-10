package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/utils_oas/xray"
)

var global string

func SetHeader(h string) {
	global = h

}

func GetHeader() (h string) {
	return global
}
func SendJson(urlp string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		json.NewEncoder(b).Encode(datajson)
	}
	//proxyUrl, err := url.Parse("http://10.20.4.15:3128")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	client := &http.Client{}
	req, err := http.NewRequest(trequest, urlp, b)
	if err != nil {
		logs.Error("Error creating request. ", err)
		return err
	}

	//Se intenta acceder a cabecera, si no existe, se realiza peticion normal.
	defer func() {
		//Catch
		if r := recover(); r != nil {

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logs.Error("Error reading response. ", err)
				return
			}
			if resp == nil {
				return
			}
			defer resp.Body.Close()
			if err := decodeJSONResponse(resp, target); err != nil {
				logs.Error("Error decoding response. ", err)
			}
		}
	}()

	//try
	header := GetHeader()
	req.Header.Set("Authorization", header)

	resp, err := client.Do(req)
	if err != nil {
		logs.Error("Error reading response. ", err)
		return err
	}
	if resp == nil {
		return nil
	}
	defer resp.Body.Close()
	return decodeJSONResponse(resp, target)
}

func GetJsonWSO2(urlp string, target interface{}) error {
	b := new(bytes.Buffer)
	//proxyUrl, err := url.Parse("http://10.20.4.15:3128")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlp, b)
	if err != nil {
		logs.Error("Error creating request. ", err)
		return err
	}

	req.Header.Set("Accept", "application/json")
	r, err := client.Do(req)
	//r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		logs.Error("error", err)
		return err
	}
	defer r.Body.Close()

	return decodeJSONResponse(r, target)
}

func GetJson(urlp string, target interface{}) error {

	req, err := http.NewRequest("GET", urlp, nil)
	if err != nil {
		logs.Error("Error reading request. ", err)
		return err
	}

	//Se intenta acceder a cabecera, si no existe, se realiza peticion normal.

	defer func() {
		//Catch
		if r := recover(); r != nil {

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logs.Error("Error reading response. ", err)
				return
			}
			if resp == nil {
				return
			}
			defer resp.Body.Close()
			if err := decodeJSONResponse(resp, target); err != nil {
				logs.Error("Error decoding response. ", err)
			}
		}
	}()

	//try
	header := GetHeader()
	req.Header.Set("Authorization", header)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logs.Error("Error reading response. ", err)
		return err
	}
	if resp == nil {
		return nil
	}
	defer resp.Body.Close()
	return decodeJSONResponse(resp, target)
}

func GetJsonTest(url string, target interface{}) (response *http.Response, err error) {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	response, err = myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return response, json.NewDecoder(response.Body).Decode(target)
}

// execRequest executes req using the provided HTTP client, wrapping the call
// with an X-Ray subsegment via the local xray package.
// The caller is responsible for closing resp.Body on success.
func execRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	subseg := xray.BeginSegmentSec(req)
	resp, err := client.Do(req)
	xray.CloseSubsegment(subseg, resp, err)
	return resp, err
}

func decodeJSONResponse(resp *http.Response, target interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("http %d: %s", resp.StatusCode, resumirBodyError(body))
	}

	if target == nil || len(bytes.TrimSpace(body)) == 0 {
		return nil
	}

	return json.Unmarshal(body, target)
}

func resumirBodyError(body []byte) string {
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		return "respuesta vacía"
	}

	if len(msg) > 1024 {
		return msg[:1024]
	}

	return msg
}
