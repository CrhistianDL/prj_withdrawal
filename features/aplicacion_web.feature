Feature: Operaciones bancarias mediante aplicacion web
  Como cliente
  Quiero realizar abonos y retiros desde el navegador
  Para consultar mi saldo y el historial de operaciones

  Scenario: Abono exitoso desde la web
    Given una cuenta de Carlos con saldo 500.00
    When se registra un abono de 200.00
    Then el estado es EXITO
    And el saldo disponible es 700.00
    And se registra la transaccion en formato JSON

  Scenario: Retiro rechazado sin alterar el saldo
    Given una cuenta de Carlos con saldo 700.00
    When se solicita un retiro de 900.00
    Then el estado es FALLO
    And el saldo permanece en 700.00
    And se registra el intento en formato JSON

  Scenario: Monto invalido
    Given una cuenta de Carlos con saldo 700.00
    When se solicita un abono de -5.00 mediante API
    Then el estado es ERROR
    And el saldo permanece en 700.00