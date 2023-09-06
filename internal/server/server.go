package server

import (
	"fmt"
	"html/template"
	"net/http"
)

// StartHTTPServer запускает HTTP-сервер.
func StartHTTPServer(addr string) error {
	return http.ListenAndServe(addr, nil)
}

func RenderOrderHTML(w http.ResponseWriter, order interface{}) {
	tmpl, err := template.ParseFiles("order.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering order: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering order: %v", err), http.StatusInternalServerError)
	}
}
