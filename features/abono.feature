# US-01b - Escenarios propuestos para extender codeLAB-02.
Feature: Abono en efectivo
  Como cliente con una cuenta activa
  Quiero abonar efectivo a mi cuenta
  Para incrementar mi saldo disponible

  Scenario: Abono exitoso
    Given Carlos tiene 500.00 en su cuenta
    When Carlos abona 200.00
    Then el abono se realiza con exito
    And el saldo de Carlos es 700.00

  Scenario: Abono rechazado por monto negativo
    Given Maria tiene 500.00 en su cuenta
    When Maria intenta abonar -100.00
    Then el abono es rechazado
    And el saldo de Maria permanece en 500.00

  Scenario: Abono rechazado por monto cero
    Given Maria tiene 500.00 en su cuenta
    When Maria intenta abonar 0.00
    Then el abono es rechazado
    And el saldo de Maria permanece en 500.00