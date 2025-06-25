@Users @UsersResetPassword
Feature: Users - Reset Password
  Scenario: POST /v1/password-reset for an existing user and checking the response status 202
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/v1/password-reset"
      """
      {
        "email": "email@ons.gov.uk"
      }
      """
    Then the HTTP status code should be "202"

  Scenario: POST /v1/password-reset missing email and checking the response status 400
    When I POST "/v1/password-reset"
      """
      {
        "email": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidEmail",
            "description": "the submitted email could not be validated"
          }
        ]
      }
      """

  Scenario: POST /v1/password-reset non ONS email address and checking the response status 202
    When I POST "/v1/password-reset"
      """
      {
        "email": "email@gmail.com"
      }
      """
    Then the HTTP status code should be "202"

  Scenario: POST /v1/password-reset Cognito internal error
    Given an internal server error is returned from Cognito
    When I POST "/v1/password-reset"
      """
      {
        "email": "internal.error@ons.gov.uk"
      }
      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: POST /v1/password-reset Cognito too many requests error
    Given an internal server error is returned from Cognito
    When I POST "/v1/password-reset"
      """
      {
        "email": "too.many@ons.gov.uk"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "TooManyRequests",
            "description": "Slow down"
          }
        ]
      }
      """

  Scenario: POST /v1/password-reset Cognito user not found
    When I POST "/v1/password-reset"
      """
      {
        "email": "email@ons.gov.uk"
      }
      """
    Then the HTTP status code should be "202"
