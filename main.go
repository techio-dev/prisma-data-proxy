package main

import (
	"fmt"
	"log"
	"net/http"

	"techio.dev/prisma/data-proxy/api"
)

func proxy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	engine := api.NewApiEngine()
	engine.HandleAuthorization(r.Header.Get("authorization"))
	engine.HandlePath(r.URL.Path)
	handle := api.NewApiHandle(engine, r.Host)
	if err := handle.New(r.Method, r.Body); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 400)
	}
	handle.RequestHeaders(r.Header)
	if err := handle.Do(); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
	}
	handle.ResponseHeaders(w)
	if err := handle.Write(w); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
	}
}
func main() {
	log.Println("[0.1.2] ListenAndServe :8080")
	if err := http.ListenAndServe(":8080", http.HandlerFunc(proxy)); err != nil {
		log.Panicln(err)
	}
}
