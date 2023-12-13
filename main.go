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
		http.Error(w, fmt.Sprintf("%v", err), http.StatusForbidden)
		return
	}
	if err := engine.HandlePath(r.URL.Path); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	handle := api.NewApiHandle(engine, r.Host)
	if err := handle.New(r.Method, r.Body); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	handle.RequestHeaders(r.Header)
	if err := handle.Do(); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	handle.ResponseHeaders(w)
	if err := handle.Write(w); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}
func main() {
	log.Printf("[%v] ListenAndServe :8080", os.Getenv("VERSION"))
	log.Println()
	if err := http.ListenAndServe(":8080", http.HandlerFunc(proxy)); err != nil {
		log.Panicln(err)
	}
}
