package handlers

import (
	"net/http"
)

func UiHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "ui.html")
		},
	)
}
