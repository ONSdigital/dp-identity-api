@Users @UsersSetPassword
Feature: Users - Set Password
  Scenario: POST /v1/users/{id}/password to set the passsword for unconfirmed user
    Given a user with username "abcd1234" exists in the database and is unconfirmed
    And I am an admin user
    When I POST "/v1/users/abcd1234/password"
    """
    """
    Then the HTTP status code should be "202"

  Scenario: POST /v1/users/{id}/password to set the passsword for confirmed user
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I POST "/v1/users/abcd1234/password"
    """
    """
    Then the HTTP status code should be "403"

  Scenario: POST /v1/users/{id}/password to set the passsword for unknown user
    Given I am an admin user
    When I POST "/v1/users/abcd1234/password"
    """
    """
    Then the HTTP status code should be "404"

  Scenario: POST /v1/users/{id}/password to set the passsword for unconfirmed user when not admin
    Given a user with username "abcd1234" exists in the database and is unconfirmed
    When I POST "/v1/users/abcd1234/password"
    """
    """
    Then the HTTP status code should be "401"
