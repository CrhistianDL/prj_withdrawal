Feature: Integridad y registro de transacciones
  Como responsable de operaciones
  Quiero registrar cada retiro o abono con su resultado en JSON
  Para verificar los saldos y los intentos de operacion.

  Scenario: Registro de abono exitoso
    Given una cuenta con saldo 500.00
    When se abonan 200.00
    Then el saldo es 700.00 y el registro tiene estado EXITO

  Scenario: Retiro rechazado por saldo insuficiente
    Given una cuenta con saldo 700.00
    When se intenta retirar 900.00
    Then el saldo permanece en 700.00 y el registro tiene estado FALLO

  Scenario: Error de monto invalido
    Given una cuenta con saldo 700.00
    When se intenta abonar -5.00
    Then el saldo permanece en 700.00 y el registro tiene estado ERROR

  Scenario: Error al guardar el registro
    Given una cuenta con saldo 100.00 y una ruta de registro no escribible
    When se intenta abonar 50.00
    Then se informa el error y el saldo permanece en 100.00