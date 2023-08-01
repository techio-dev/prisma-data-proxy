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
	if err := engine.HandleAuthorization(r.Header.Get("authorization")); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusForbidden)
		return
	}
	if err := engine.HandlePath(r.URL.Path); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	handle := api.NewApiHandle(engine, r.Host)
	if err := handle.New(r.Method, r.Body); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	handle.RequestHeaders(r.Header)
	if err := handle.Do(); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	handle.ResponseHeaders(w)
	if err := handle.Write(w); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}
func main() {
	log.Println("[0.1.3] ListenAndServe :8080")
	if err := http.ListenAndServe(":8080", http.HandlerFunc(proxy)); err != nil {
		log.Panicln(err)
	}
}
