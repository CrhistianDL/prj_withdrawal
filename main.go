package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"prj_withdrawal/headers"
	"prj_withdrawal/web"
	"strconv"
	"strings"
)

func leerMonto(sc *bufio.Scanner, prompt string) (float64, error) {
	fmt.Print(prompt)
	if !sc.Scan() {
		return 0, fmt.Errorf("entrada finalizada")
	}
	return strconv.ParseFloat(strings.TrimSpace(sc.Text()), 64)
}
func interactivo(archivo string) {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Print("Nombre del titular: ")
	if !sc.Scan() {
		return
	}
	nombre := strings.TrimSpace(sc.Text())
	saldo, err := leerMonto(sc, "Saldo inicial (S/): ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	c, err := headers.NewCuenta(nombre, saldo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for {
		fmt.Println("\n=== SISTEMA BANCARIO ===")
		fmt.Println("1. Retirar efectivo\n2. Abonar efectivo\n3. Consultar saldo\n4. Salir")
		fmt.Print("Seleccione una opcion: ")
		if !sc.Scan() {
			return
		}
		switch strings.TrimSpace(sc.Text()) {
		case "1", "2":
			op := strings.TrimSpace(sc.Text())
			m, e := leerMonto(sc, "Monto (S/): ")
			if e != nil {
				fmt.Println("Error:", e)
				continue
			}
			tipo := "ABONO"
			if op == "1" {
				tipo = "RETIRO"
			}
			reg, logErr := headers.ProcesarTransaccion(c, tipo, m, archivo)
			if logErr != nil {
				fmt.Println("Error de registro; saldo sin cambios:", logErr)
			} else {
				fmt.Printf("Estado: %s | %s\n", reg.Estado, reg.Detalle)
				mostrarJSON(reg)
			}
			fmt.Printf("Saldo actual: S/ %.2f\n", c.ObtenerSaldo())
		case "3":
			fmt.Printf("Titular: %s | Saldo: S/ %.2f\n", c.Nombre(), c.ObtenerSaldo())
		case "4":
			fmt.Println("Fin del programa")
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}
func mostrarJSON(reg headers.RegistroTransaccion) {
	b, _ := json.MarshalIndent(reg, "", "  ")
	fmt.Println(string(b))
}
func main() {
	nombre := flag.String("nombre", "", "Titular de cuenta")
	saldo := flag.Float64("saldo", 0, "Saldo inicial")
	retiro := flag.Float64("retiro", 0, "Monto de retiro")
	abono := flag.Float64("abono", 0, "Monto de abono")
	registro := flag.String("registro", "registros/transacciones.jsonl", "Archivo JSON Lines")
	webMode := flag.Bool("web", false, "Iniciar aplicacion web")
	puerto := flag.String("puerto", "8080", "Puerto local para aplicacion web")
	test := flag.String("test", "", "Escenario: retiro, abono o interactivo")
	flag.Parse()
	if *webMode {
		titular := *nombre
		if titular == "" {
			titular = "Carlos"
		}
		servidor, err := web.New(titular, *saldo, *registro)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		direccion := "127.0.0.1:" + *puerto
		fmt.Println("Aplicacion web disponible en http://" + direccion)
		if err := http.ListenAndServe(direccion, servidor.Handler()); err != nil {
			fmt.Fprintln(os.Stderr, "Error del servidor:", err)
			os.Exit(1)
		}
		return
	}
	if *test == "interactivo" || (*nombre == "" && *test == "") {
		interactivo(*registro)
		return
	}
	if *nombre == "" {
		fmt.Fprintln(os.Stderr, "Debe indicar -nombre")
		os.Exit(1)
	}
	switch *test {
	case "retiro":
		headers.TestEscenarioRetiroExito(*nombre, *saldo, *retiro)
	case "abono":
		headers.TestEscenarioAbono(*nombre, *saldo, *abono)
	case "":
		c, err := headers.NewCuenta(*nombre, *saldo)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if *abono != 0 {
			reg, e := headers.ProcesarTransaccion(c, "ABONO", *abono, *registro)
			if e != nil {
				fmt.Println("Error:", e)
				return
			}
			mostrarJSON(reg)
		}
		if *retiro != 0 {
			reg, e := headers.ProcesarTransaccion(c, "RETIRO", *retiro, *registro)
			if e != nil {
				fmt.Println("Error:", e)
				return
			}
			mostrarJSON(reg)
		}
		fmt.Printf("Saldo final: S/ %.2f\n", c.ObtenerSaldo())
	default:
		fmt.Fprintln(os.Stderr, "Test desconocido:", *test)
		os.Exit(1)
	}
}