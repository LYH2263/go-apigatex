package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/LYH2263/go-apigatex"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "closed": s.gw.Closed()})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.gw.Stats())
}

func (s *Server) handleRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.gw.ListRoutes())
	case http.MethodPost:
		var spec apigatex.RouteSpec
		if err := readJSON(r, &spec); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.gw.AddRoute(spec); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, spec)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRouteID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/routes/")
	id = strings.Trim(id, "/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		spec, ok := s.gw.GetRoute(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, spec)
	case http.MethodDelete:
		if err := s.gw.RemoveRoute(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type tryBody struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (s *Server) handleTry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var tb tryBody
	if err := readJSON(r, &tb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if tb.Method == "" {
		tb.Method = http.MethodGet
	}
	res, err := s.gw.TryRequest(r.Context(), tb.Method, tb.Path, tb.Headers, []byte(tb.Body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	// 去掉 /api/proxy 前缀，走网关 ServeHTTP
	r2 := r.Clone(r.Context())
	r2.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/proxy")
	if r2.URL.Path == "" {
		r2.URL.Path = "/"
	}
	s.gw.ServeHTTP(w, r2)
}

func readJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	return dec.Decode(dst)
}
