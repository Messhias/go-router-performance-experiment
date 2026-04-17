Feature: BDD bootstrap
  Scenario: Feature suite executes and reports failure
    Given the BDD suite is configured
    When I run the feature tests
    Then I should see a failing scenario report