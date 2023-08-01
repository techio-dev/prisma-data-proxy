package api

import (
	"errors"
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

func (e *ApiEngine) HandleAuthorization(authorization string) error {
	if authorization == "" {
		return errors.New("authorization not found")
	}
	e.apiKey = strings.TrimPrefix(authorization, "Bearer ")
	return nil
}

func (e *ApiEngine) HandlePath(path string) error {
	sections := strings.Split(path, "/")
	if len(sections) <= 3 {
		return errors.New("parse url path failure")
	}
	e.version = sections[1]
	e.hash = sections[2]
	isTransactionHandle := false
	if sections[3] == "transaction" {
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
	return nil
}

func (e *ApiEngine) EngineUrl() string {
	params := []string{e.apiKey, strings.ReplaceAll(e.version, ".", "")}
	if e.isTransaction {
		params = append(params, "tx")
	}
	log.Printf("ENV: %v", params)
	return os.Getenv(strings.Join(params, "-"))
}
func (e *ApiEngine) Url() (string, error) {
	engineUrl := e.EngineUrl()
	if engineUrl == "" {
		return "", errors.New("not found engine url")
	}
	path := []string{engineUrl}
	path = append(path, e.path...)
	log.Printf("URL: %v", path)
	return strings.Join(path, "/"), nil
}
