package headers

import (
	"errors"
	"fmt"
	"math"
)

// Cuenta representa la entidad de negocio
type Cuenta struct {
	nombre        string
	saldoCentavos int64
}

func aCentavos(monto float64) (int64, error) {
	if math.IsNaN(monto) || math.IsInf(monto, 0) || monto < 0 || monto > 1e12 {
		return 0, errors.New("monto invalido")
	}
	centavos := math.Round(monto * 100)
	if math.Abs(monto*100-centavos) > 0.000001 {
		return 0, errors.New("solo se permiten dos decimales")
	}
	return int64(centavos), nil
}

func NewCuenta(nombre string, saldoInicial float64) (*Cuenta, error) {
	if nombre == "" {
		return nil, errors.New("el nombre es obligatorio")
	}
	centavos, err := aCentavos(saldoInicial)
	if err != nil {
		return nil, fmt.Errorf("saldo inicial: %w", err)
	}
	return &Cuenta{nombre: nombre, saldoCentavos: centavos}, nil
}
func (c *Cuenta) ObtenerSaldo() float64 { return float64(c.saldoCentavos) / 100 }
func (c *Cuenta) Nombre() string        { return c.nombre }
func (c *Cuenta) Retirar(cantidad float64) (bool, error) {
	centavos, err := aCentavos(cantidad)
	if err != nil {
		return false, err
	}
	if centavos == 0 {
		return false, errors.New("el retiro debe ser mayor que cero")
	}
	if centavos > c.saldoCentavos {
		return false, nil
	}
	c.saldoCentavos -= centavos
	return true, nil
}

// Un abono valido incrementa el saldo
func (c *Cuenta) Abonar(cantidad float64) (bool, error) {
	centavos, err := aCentavos(cantidad)
	if err != nil {
		return false, err
	}
	if centavos == 0 {
		return false, errors.New("el abono debe ser mayor que cero")
	}
	if centavos > math.MaxInt64-c.saldoCentavos {
		return false, errors.New("saldo excede el limite permitido")
	}
	c.saldoCentavos += centavos
	return true, nil
}