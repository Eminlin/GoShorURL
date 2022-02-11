package server

import "net/http"

func WebRun() {
	http.Handle("/", http.NotFoundHandler())
	http.ListenAndServe(":1818", nil)
}
