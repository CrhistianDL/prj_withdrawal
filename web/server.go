package web

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"prj_withdrawal/headers"
	"strconv"
	"strings"
)

//go:embed templates/* static/*
var assets embed.FS

type Server struct {
	servicio *headers.ServicioBancario
	registro string
	tpl      *template.Template
}

func New(nombre string, saldo float64, registro string) (*Server, error) {
	svc, err := headers.NuevoServicioBancario(nombre, saldo, registro)
	if err != nil {
		return nil, err
	}
	tpl, err := template.ParseFS(assets, "templates/index.html")
	if err != nil {
		return nil, err
	}
	return &Server{servicio: svc, registro: registro, tpl: tpl}, nil
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.index)
	mux.HandleFunc("GET /api/estado", s.estado)
	mux.HandleFunc("POST /api/transacciones", s.operar)
	mux.HandleFunc("GET /api/transacciones", s.historial)
	mux.Handle("GET /static/", http.FileServerFS(assets))
	return mux
}
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = s.tpl.Execute(w, nil)
}
func jsonOut(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
func (s *Server) estado(w http.ResponseWriter, r *http.Request) {
	n, b := s.servicio.Estado()
	jsonOut(w, 200, map[string]any{"titular": n, "saldo": b})
}
func (s *Server) operar(w http.ResponseWriter, r *http.Request) {
	// Limitar el cuerpo y rechazar campos desconocidos.
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	var req struct {
		Operacion string      `json:"operacion"`
		Monto     json.Number `json:"monto"`
	}
	d := json.NewDecoder(r.Body)
	d.UseNumber()
	d.DisallowUnknownFields()
	if err := d.Decode(&req); err != nil {
		jsonOut(w, 400, map[string]string{"error": "solicitud JSON invalida"})
		return
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		jsonOut(w, 400, map[string]string{"error": "solo se permite un objeto JSON"})
		return
	}
	if req.Operacion != "ABONO" && req.Operacion != "RETIRO" {
		jsonOut(w, 400, map[string]string{"error": "operacion invalida"})
		return
	}
	monto, err := strconv.ParseFloat(req.Monto.String(), 64)
	if err != nil {
		jsonOut(w, 400, map[string]string{"error": "monto invalido"})
		return
	}
	reg, err := s.servicio.Operar(req.Operacion, monto)
	if err != nil {
		jsonOut(w, 500, map[string]string{"error": "no se pudo guardar el registro; saldo sin cambios"})
		return
	}
	jsonOut(w, 200, reg)
}
func (s *Server) historial(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(s.registro)
	if errors.Is(err, os.ErrNotExist) {
		jsonOut(w, 200, []headers.RegistroTransaccion{})
		return
	}
	if err != nil {
		jsonOut(w, 500, map[string]string{"error": "no se pudo leer el historial"})
		return
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, 2<<20))
	if err != nil {
		jsonOut(w, 500, map[string]string{"error": "error al leer historial"})
		return
	}
	rows := []headers.RegistroTransaccion{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var row headers.RegistroTransaccion
		if json.Unmarshal([]byte(line), &row) != nil {
			jsonOut(w, 500, map[string]string{"error": "registro JSON invalido"})
			return
		}
		rows = append(rows, row)
	}
	// Mostrar solo las 30 operaciones mas recientes.
	if len(rows) > 30 {
		rows = rows[len(rows)-30:]
	}
	jsonOut(w, 200, rows)
}
func ValidarRegistro(path string) error {
	if filepath.Ext(path) != ".jsonl" {
		return fmt.Errorf("el registro debe usar extension .jsonl")
	}
	return nil
}