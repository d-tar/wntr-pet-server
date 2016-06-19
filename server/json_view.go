package server

import (
	"encoding/json"
	"net/http"
)

type JsonMavRenderable struct {
	data     interface{}
	httpCode int
}

func (v *JsonMavRenderable) Render(w http.ResponseWriter) error {
	b, err := json.Marshal(v.data)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type","application/json;charset=UTF-8")
	if v.httpCode != 0 {
		w.WriteHeader(v.httpCode)
	}
	w.Write(b)
	return nil
}

type BaseJsonResponse struct {
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
	Status string      `json:"status"`
}

func MavOk(data interface{}) MavRenderable {
	return &JsonMavRenderable{
		data: &BaseJsonResponse{
			Data:   data,
			Status: "ok",
		},
	}
}

func MavErr(error interface{}) MavRenderable {
	return &JsonMavRenderable{
		data: &BaseJsonResponse{
			Error:  error,
			Status: "err",
		},
		httpCode: 500,
	}
}
