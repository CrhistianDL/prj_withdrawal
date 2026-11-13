package headers

import "testing"

func TestRetiroExitoso(t *testing.T) {
	c, _ := NewCuenta("Alice", 234.56)
	ok, err := c.Retirar(200)
	if err != nil || !ok || c.ObtenerSaldo() != 34.56 {
		t.Fatalf("retiro: ok=%t err=%v saldo=%.2f", ok, err, c.ObtenerSaldo())
	}
}
func TestRetiroInsuficiente(t *testing.T) {
	c, _ := NewCuenta("Hamza", 198.76)
	ok, err := c.Retirar(200)
	if err != nil || ok || c.ObtenerSaldo() != 198.76 {
		t.Fatalf("retiro insuficiente: ok=%t err=%v saldo=%.2f", ok, err, c.ObtenerSaldo())
	}
}
func TestAbonoExitoso(t *testing.T) {
	c, _ := NewCuenta("Carlos", 500)
	ok, err := c.Abonar(200)
	if err != nil || !ok || c.ObtenerSaldo() != 700 {
		t.Fatalf("abono: ok=%t err=%v saldo=%.2f", ok, err, c.ObtenerSaldo())
	}
}
func TestAbonoInvalido(t *testing.T) {
	for _, m := range []float64{-100, 0, 1.001} {
		c, _ := NewCuenta("Maria", 500)
		ok, err := c.Abonar(m)
		if err == nil || ok || c.ObtenerSaldo() != 500 {
			t.Fatalf("monto %v: ok=%t err=%v saldo=%.2f", m, ok, err, c.ObtenerSaldo())
		}
	}
}
func TestSecuenciaAbonoRetiro(t *testing.T) {
	c, _ := NewCuenta("Luis", 100)
	a, e := c.Abonar(50)
	r, e2 := c.Retirar(120)
	if !a || !r || e != nil || e2 != nil || c.ObtenerSaldo() != 30 {
		t.Fatalf("secuencia saldo %.2f", c.ObtenerSaldo())
	}
}
func TestSaldoInicialInvalido(t *testing.T) {
	if _, err := NewCuenta("A", -1); err == nil {
		t.Fatal("se acepto saldo inicial negativo")
	}
}