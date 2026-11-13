package headers

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRegistroExitoFalloError(t *testing.T) {
	c, _ := NewCuenta("Carlos", 500)
	archivo := filepath.Join(t.TempDir(), "transacciones.jsonl")
	casos := []struct {
		op     string
		monto  float64
		estado string
		saldo  float64
	}{
		{"ABONO", 200, "EXITO", 700},
		{"RETIRO", 900, "FALLO", 700},
		{"ABONO", -5, "ERROR", 700},
		{"RETIRO", 100, "EXITO", 600},
	}
	for _, caso := range casos {
		r, e := ProcesarTransaccion(c, caso.op, caso.monto, archivo)
		if e != nil || r.Estado != caso.estado || c.ObtenerSaldo() != caso.saldo {
			t.Fatalf("caso %+v: registro=%+v error=%v saldo=%v", caso, r, e, c.ObtenerSaldo())
		}
	}
	f, e := os.Open(archivo)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		var r RegistroTransaccion
		if e := json.Unmarshal(sc.Bytes(), &r); e != nil {
			t.Fatal(e)
		}
		n++
	}
	if e := sc.Err(); e != nil {
		t.Fatal(e)
	}
	if n != 4 {
		t.Fatalf("se esperaban 4 registros, hay %d", n)
	}
}
func TestSinCambioSiFallaRegistro(t *testing.T) {
	c, _ := NewCuenta("Ana", 100)
	archivo := t.TempDir() // un directorio no se puede abrir como archivo
	_, err := ProcesarTransaccion(c, "ABONO", 50, archivo)
	if err == nil || c.ObtenerSaldo() != 100 {
		t.Fatalf("saldo cambiado o falta error: %v", err)
	}
}
func TestOperacionDesconocida(t *testing.T) {
	c, _ := NewCuenta("Ana", 100)
	r, e := ProcesarTransaccion(c, "OTRA", 10, filepath.Join(t.TempDir(), "log.jsonl"))
	if e != nil || r.Estado != "ERROR" || c.ObtenerSaldo() != 100 {
		t.Fatalf("resultado inesperado: %+v %v", r, e)
	}
}