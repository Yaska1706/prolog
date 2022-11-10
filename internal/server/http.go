package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHttpServer(addr string) *http.Server {
	httpsrvr := newHttpServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrvr.handleConsume).Methods("GET")
	r.HandleFunc("/", httpsrvr.handleProduce).Methods("POST")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}
type ConsumeResponse struct {
	Record Record `json:"record"`
}

type httpServer struct {
	Log *Log
}

func newHttpServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func (h *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	off, err := h.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ProduceResponse{Offset: off}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := h.Log.Read(req.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ConsumeResponse{Record: record}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
