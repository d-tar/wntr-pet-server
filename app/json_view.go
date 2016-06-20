package app

import (
	"github.com/d-tar/wntr/webmvc"
	"net/http"
)

type BaseJsonResponse struct {
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
	Status string      `json:"status"`
}

func NewCustomJsonView() webmvc.WebView {
	var q webmvc.WebViewFunc = RenderJsonView
	return q
}

func RenderJsonView(mav webmvc.WebResult, w http.ResponseWriter, r *http.Request) error {
	var newMav webmvc.WebResult
	if mav.HttpCode() == 200 {
		newMav = webmvc.WebOk(&BaseJsonResponse{
			Data:   mav.Model(),
			Status: "ok",
		})
	} else {
		newMav = webmvc.WebErr(&BaseJsonResponse{
			Error:  mav.Model(),
			Status: "err",
		})
	}

	return webmvc.RenderJsonView(newMav, w, r)
}
