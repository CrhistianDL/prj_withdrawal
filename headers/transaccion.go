package headers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type RegistroTransaccion struct {
	Fecha          string  `json:"fecha"`
	Titular        string  `json:"titular"`
	Operacion      string  `json:"operacion"`
	Monto          float64 `json:"monto"`
	SaldoAnterior  float64 `json:"saldo_anterior"`
	SaldoPosterior float64 `json:"saldo_posterior"`
	Estado         string  `json:"estado"`
	Detalle        string  `json:"detalle"`
}

// ProcesarTransaccion valida sobre una copia, registra el resultado y solo entonces
// confirma el nuevo saldo. Si falla la escritura, el saldo original no cambia.
// El registro es JSON Lines: un objeto JSON por linea.
func ProcesarTransaccion(c *Cuenta, operacion string, monto float64, archivo string) (RegistroTransaccion, error) {
	if c == nil {
		return RegistroTransaccion{}, errors.New("cuenta no inicializada")
	}
	if archivo == "" {
		return RegistroTransaccion{}, errors.New("ruta de registro obligatoria")
	}
	reg := RegistroTransaccion{Fecha: time.Now().Format(time.RFC3339Nano), Titular: c.Nombre(), Operacion: operacion, Monto: monto, SaldoAnterior: c.ObtenerSaldo(), SaldoPosterior: c.ObtenerSaldo()}
	copia := *c
	var ok bool
	var err error
	switch operacion {
	case "RETIRO":
		ok, err = copia.Retirar(monto)
	case "ABONO":
		ok, err = copia.Abonar(monto)
	default:
		err = fmt.Errorf("operacion desconocida: %s", operacion)
	}
	if err != nil {
		reg.Estado = "ERROR"
		reg.Detalle = err.Error()
	} else if !ok {
		reg.Estado = "FALLO"
		reg.Detalle = "saldo insuficiente"
	} else {
		reg.Estado = "EXITO"
		reg.Detalle = "operacion completada"
		reg.SaldoPosterior = copia.ObtenerSaldo()
	}
	bytes, marshalErr := json.Marshal(reg)
	if marshalErr != nil {
		return RegistroTransaccion{}, fmt.Errorf("serializar JSON: %w", marshalErr)
	}
	if err := os.MkdirAll(filepath.Dir(archivo), 0755); err != nil {
		return RegistroTransaccion{}, fmt.Errorf("crear directorio de registros: %w", err)
	}
	f, openErr := os.OpenFile(archivo, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if openErr != nil {
		return RegistroTransaccion{}, fmt.Errorf("abrir registro: %w", openErr)
	}
	if _, writeErr := f.Write(append(bytes, '\n')); writeErr != nil {
		f.Close()
		return RegistroTransaccion{}, fmt.Errorf("escribir registro: %w", writeErr)
	}
	if syncErr := f.Sync(); syncErr != nil {
		f.Close()
		return RegistroTransaccion{}, fmt.Errorf("sincronizar registro: %w", syncErr)
	}
	if closeErr := f.Close(); closeErr != nil {
		return RegistroTransaccion{}, fmt.Errorf("cerrar registro: %w", closeErr)
	}
	if reg.Estado == "EXITO" {
		c.saldoCentavos = copia.saldoCentavos
	}
	return reg, nil
}