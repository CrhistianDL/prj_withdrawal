package headers

import (
	"errors"
	"sync"
)

// ServicioBancario protege el estado de la cuenta durante peticiones concurrentes.
type ServicioBancario struct {
	mu      sync.Mutex
	cuenta  *Cuenta
	archivo string
}

func NuevoServicioBancario(nombre string, saldo float64, archivo string) (*ServicioBancario, error) {
	c, err := NewCuenta(nombre, saldo)
	if err != nil {
		return nil, err
	}
	if archivo == "" {
		return nil, errors.New("ruta de registro obligatoria")
	}
	return &ServicioBancario{cuenta: c, archivo: archivo}, nil
}
func (s *ServicioBancario) Estado() (string, float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cuenta.Nombre(), s.cuenta.ObtenerSaldo()
}
func (s *ServicioBancario) Operar(tipo string, monto float64) (RegistroTransaccion, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return ProcesarTransaccion(s.cuenta, tipo, monto, s.archivo)
}