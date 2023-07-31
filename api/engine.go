package api

import (
	"log"
	"os"
	"strings"
)

type ApiEngine struct {
	apiKey             string
	version            string
	hash               string
	isTransaction      bool
	isTransactionStart bool
	path               []string
}

func NewApiEngine() *ApiEngine {
	return &ApiEngine{}
}

func (e *ApiEngine) HandleAuthorization(authorization string) {
	e.apiKey = strings.TrimPrefix(authorization, "Bearer ")
}

func (e *ApiEngine) HandlePath(path string) {
	sections := strings.Split(path, "/")
	e.version = sections[1]
	e.hash = sections[2]
	isTransactionHandle := false
	if len(sections) > 3 && sections[3] == "transaction" {
		e.isTransaction = true
		switch sections[len(sections)-1] {
		case "start":
			e.isTransactionStart = true
			isTransactionHandle = true
		case "commit":
			isTransactionHandle = true
		case "rollback":
			isTransactionHandle = true
		}
	}
	if isTransactionHandle {
		e.path = sections[3:]
	} else {
		e.path = []string{""}
	}
}

func (e *ApiEngine) EngineUrl() string {
	params := []string{e.apiKey, strings.ReplaceAll(e.version, ".", "")}
	if e.isTransaction {
		params = append(params, "0")
	}
	return os.Getenv(strings.Join(params, "-"))
}
func (e *ApiEngine) Url() string {
	path := []string{e.EngineUrl()}
	path = append(path, e.path...)
	log.Printf("URL: %v", path)
	return strings.Join(path, "/")
}
