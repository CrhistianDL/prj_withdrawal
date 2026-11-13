package headers

import "fmt"

// Escenarios demostrativos Given/When/Then para la CLI.
func TestEscenarioRetiroExito(nombre string, saldo, monto float64) {
	c, err := NewCuenta(nombre, saldo)
	if err != nil {
		fmt.Println("Given invalido:", err)
		return
	}
	fmt.Printf("Given: %s tiene S/ %.2f\n", nombre, saldo)
	fmt.Printf("When: intenta retirar S/ %.2f\n", monto)
	ok, err := c.Retirar(monto)
	fmt.Printf("Then: aprobado=%t, saldo=S/ %.2f, error=%v\n", ok, c.ObtenerSaldo(), err)
}
func TestEscenarioAbono(nombre string, saldo, monto float64) {
	c, err := NewCuenta(nombre, saldo)
	if err != nil {
		fmt.Println("Given invalido:", err)
		return
	}
	fmt.Printf("Given: %s tiene S/ %.2f\n", nombre, saldo)
	fmt.Printf("When: intenta abonar S/ %.2f\n", monto)
	ok, err := c.Abonar(monto)
	fmt.Printf("Then: aprobado=%t, saldo=S/ %.2f, error=%v\n", ok, c.ObtenerSaldo(), err)
}