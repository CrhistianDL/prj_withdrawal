# CuentaGo — Aplicación bancaria web en Go

Aplicación local de demostración para consultar una cuenta, realizar abonos y retiros y revisar los registros de cada transacción. El proyecto conserva la lógica del sistema de terminal y la adapta a una interfaz web mediante la biblioteca estándar de Go.

## Requisitos

- Go **1.22 o superior** (recomendado Go 1.27 en Windows).
- Navegador web actualizado.
- No se requieren paquetes externos, base de datos ni Node.js.

## Iniciar la aplicación web

Abre una terminal dentro de `prj_withdrawal` y ejecuta:

```powershell
go run . -web -nombre Carlos -saldo 500
```

Abre **http://127.0.0.1:8080** en el navegador. Para cambiar el puerto:

```powershell
go run . -web -nombre Carlos -saldo 500 -puerto 8081
```

Para detener el servidor: `Ctrl+C` en la terminal. Al reiniciar, el saldo vuelve al saldo inicial indicado; el historial JSON se conserva en el archivo de registros.

## Cómo utilizarla

1. Consulta el titular y saldo inicial en la tarjeta superior.
2. Elige "Abonar" o "Retirar", escribe el monto y confirma.
3. Observa el estado de la operación y el saldo actualizado.
4. Revisa el historial de transacciones; muestra hasta las 30 operaciones más recientes.

### Ejemplo de pruebas manuales

Con titular `Carlos` y saldo inicial `500`:

| Operación | Monto | Estado esperado | Saldo posterior |
|---|---:|---|---:|
| Abono | 200 | EXITO | 700 |
| Retiro | 900 | FALLO | 700 |
| Abono | -5 | ERROR | 700 |

**Nota:** la interfaz HTML impide ingresar montos negativos, por lo que el caso `-5` se prueba mediante API o pruebas automatizadas. Para la captura web basta con demostrar el abono y el retiro con saldo insuficiente.

## Pruebas automáticas

```powershell
go test -count=1 -v ./...
```

Se ejecutan las pruebas existentes del dominio y las nuevas pruebas del servidor HTTP.

## Registros JSON

Cada transacción se registra como un objeto JSON por línea (JSON Lines), por defecto en `registros/transacciones.jsonl`. Campos: `fecha`, `titular`, `operacion`, `monto`, `saldo_anterior`, `saldo_posterior`, `estado` y `detalle`. Los estados son `EXITO`, `FALLO` y `ERROR`.

Para consultar el archivo en PowerShell:

```powershell
Get-Content .\registros\transacciones.jsonl
```

El historial es persistente, pero el saldo se mantiene en memoria y se reinicia al arrancar. Por ello, los registros de ejecuciones anteriores pueden permanecer visibles; para una demostración limpia, elimina el archivo de registros antes de iniciar la aplicación (con el servidor detenido).

## Estructura del proyecto

- `headers/`: entidad Cuenta, operaciones, validaciones, registro JSON y pruebas de negocio.
- `web/server.go`: rutas HTTP, servicio y consultas al historial.
- `web/templates/index.html`: página principal.
- `web/static/`: estilos e interacción JavaScript.
- `web/server_test.go`: pruebas de integración HTTP.
- `features/`: escenarios Gherkin de las funcionalidades.
- `main.go`: punto de entrada CLI y web.