package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type BodyStart struct {
	Id        string `json:"id"`
	DataProxy struct {
		Endpoint string `json:"endpoint"`
	} `json:"data-proxy"`
}

type ApiHandle struct {
	host   string
	engine *ApiEngine
	req    *http.Request
	resp   *http.Response
}

func NewApiHandle(engine *ApiEngine, host string) *ApiHandle {
	return &ApiHandle{engine: engine, host: host}
}

func (r *ApiHandle) New(method string, body io.ReadCloser) error {
	req, err := http.NewRequest(method, r.engine.Url(), body)
	if err != nil {
		return err
	}
	r.req = req
	return nil
}

func (r *ApiHandle) RequestHeaders(src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			r.req.Header.Add(k, v)
		}
	}
}

func (r *ApiHandle) Do() error {
	resp, err := http.DefaultClient.Do(r.req)
	if err != nil {
		return err
	}
	r.resp = resp
	return nil
}

func (r *ApiHandle) ResponseHeaders(w http.ResponseWriter) {
	w.WriteHeader(r.resp.StatusCode)
	for k, vv := range r.resp.Header {
		for _, v := range vv {
			w.Header().Set(k, v)
		}
	}
}

func (r *ApiHandle) Write(w http.ResponseWriter) error {
	defer r.resp.Body.Close()
	if r.engine.isTransactionStart {
		return r.WriteStart(w)
	} else {
		_, err := io.Copy(w, r.resp.Body)
		return err
	}
}

func (r *ApiHandle) WriteStart(w http.ResponseWriter) error {
	data := &BodyStart{}
	if err := json.NewDecoder(r.resp.Body).Decode(data); err != nil {
		return err
	}
	endpoint := []string{r.host}
	endpoint = append(endpoint, r.engine.version)
	endpoint = append(endpoint, "hash")
	endpoint = append(endpoint, "transaction")
	endpoint = append(endpoint, data.Id)
	log.Printf("endpoint: %v", endpoint)
	data.DataProxy.Endpoint = strings.Join(endpoint, "/")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}
