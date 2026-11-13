package web

import (
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestWebOperacionesYRegistro(t *testing.T) {
	s, e := New("Carlos", 500, filepath.Join(t.TempDir(), "transacciones.jsonl"))
	if e != nil {
		t.Fatal(e)
	}
	h := s.Handler()
	for _, tc := range []struct {
		tipo   string
		monto  string
		estado string
		saldo  float64
	}{
		{"ABONO", "200", "EXITO", 700}, {"RETIRO", "900", "FALLO", 700}, {"ABONO", "-5", "ERROR", 700},
	} {
		req := httptest.NewRequest("POST", "/api/transacciones", strings.NewReader(`{"operacion":"`+tc.tipo+`","monto":`+tc.monto+`}`))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("status %d: %s", w.Code, w.Body.String())
		}
		var reg struct {
			Estado         string  `json:"estado"`
			SaldoPosterior float64 `json:"saldo_posterior"`
		}
		if e := json.Unmarshal(w.Body.Bytes(), &reg); e != nil {
			t.Fatal(e)
		}
		if reg.Estado != tc.estado || reg.SaldoPosterior != tc.saldo {
			t.Fatalf("got %+v want %s %.2f", reg, tc.estado, tc.saldo)
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/api/transacciones", nil))
	var regs []json.RawMessage
	if e := json.Unmarshal(w.Body.Bytes(), &regs); e != nil || len(regs) != 3 {
		t.Fatalf("registros=%d error=%v", len(regs), e)
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/api/estado", nil))
	var state struct {
		Saldo float64 `json:"saldo"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &state)
	if state.Saldo != 700 {
		t.Fatal(state)
	}
}
func TestWebPaginaYValidacion(t *testing.T) {
	s, e := New("Ana", 100, filepath.Join(t.TempDir(), "log.jsonl"))
	if e != nil {
		t.Fatal(e)
	}
	h := s.Handler()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "Gestión de cuenta") {
		t.Fatal("pagina web no disponible")
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("POST", "/api/transacciones", strings.NewReader(`{"operacion":"ABONO","monto":1.001}`)))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"estado":"ERROR"`) {
		t.Fatalf("respuesta: %d %s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("POST", "/api/transacciones", strings.NewReader(`{"operacion":"OTRA","monto":10}`)))
	if w.Code != 400 {
		t.Fatal("debe rechazar operacion desconocida")
	}
}