@Users @UsersChangePassword
Feature: Users - Change Password
  Scenario: PUT /v1/users/self/password and checking the response status 202
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "session": "auth-challenge-session"
      }
      """
    Then the HTTP status code should be "202"
    And the response header "Authorization" should be "Bearer accessToken"
    And the response header "ID" should be "idToken"
    And the response header "Refresh" should be "refreshToken"

  Scenario: PUT /v1/users/self/password missing type and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "session": "auth-challenge-session"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "UnknownRequestType",
            "description": "unknown password change type received"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password new password required type with verification token and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "verification_token": "verification-token"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidChallengeSession",
            "description": "no valid auth challenge session was provided"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing email and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "",
        "password": "Password2",
        "session": "auth-challenge-session"
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

  Scenario: PUT /v1/users/self/password missing password and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "email@ons.gov.uk",
        "password": "",
        "session": "auth-challenge-session"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "the submitted password could not be validated"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing session and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "session": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidChallengeSession",
            "description": "no valid auth challenge session was provided"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing email and password and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "",
        "password": "",
        "session": "auth-challenge-session"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "the submitted password could not be validated"
          },
          {
            "code": "InvalidEmail",
            "description": "the submitted email could not be validated"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing email and session and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "",
        "password": "Password2",
        "session": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidEmail",
            "description": "the submitted email could not be validated"
          },
          {
            "code": "InvalidChallengeSession",
            "description": "no valid auth challenge session was provided"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing password and session and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "email@ons.gov.uk",
        "password": "",
        "session": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "the submitted password could not be validated"
          },
          {
            "code": "InvalidChallengeSession",
            "description": "no valid auth challenge session was provided"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing email and password and session and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "",
        "password": "",
        "session": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "the submitted password could not be validated"
          },
          {
            "code": "InvalidEmail",
            "description": "the submitted email could not be validated"
          },
          {
            "code": "InvalidChallengeSession",
            "description": "no valid auth challenge session was provided"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password Cognito internal error
    Given an internal server error is returned from Cognito
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "email@ons.gov.uk",
        "password": "internalerrorException",
        "session": "auth-challenge-session"
      }
      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: PUT /v1/users/self/password Cognito invalid password
    Given I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "email@ons.gov.uk",
        "password": "invalidpassword",
        "session": "auth-challenge-session"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "password does not meet requirements"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password Cognito user not found
    Given I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "NewPasswordRequired",
        "email": "email@ons.gov.uk",
        "password": "Password",
        "session": "auth-challenge-session"
      }
      """
    Then the HTTP status code should be "202"

  #   Change password - forgotten password
  Scenario: PUT /v1/users/self/password forgotten password type and checking the response status 202
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "verification_token": "verification-token"
      }
      """
    Then the HTTP status code should be "202"

  Scenario: PUT /v1/users/self/password missing type and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "verification_token": "verification-token"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "UnknownRequestType",
            "description": "unknown password change type received"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password forgotten password type with challenge session and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "session": "auth-challenge-session"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidToken",
            "description": "the submitted token could not be validated"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing email and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "",
        "password": "Password2",
        "verification_token": "verification-token"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidUserId",
            "description": "the user id was missing"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing password and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "",
        "verification_token": "verification-token"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "the submitted password could not be validated"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing verification token and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "verification_token": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidToken",
            "description": "the submitted token could not be validated"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing email and password and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "",
        "password": "",
        "verification_token": "verification-token"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "the submitted password could not be validated"
          },
          {
            "code": "InvalidUserId",
            "description": "the user id was missing"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing email and verification token and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "",
        "password": "Password2",
        "verification_token": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidUserId",
            "description": "the user id was missing"
          },
          {
            "code": "InvalidToken",
            "description": "the submitted token could not be validated"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password missing email and password and verification token and checking the response status 400
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "",
        "password": "",
        "verification_token": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "the submitted password could not be validated"
          },
          {
            "code": "InvalidUserId",
            "description": "the user id was missing"
          },
          {
            "code": "InvalidToken",
            "description": "the submitted token could not be validated"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password Cognito internal error for ForgottenPassword
    Given an internal server error is returned from Cognito
    And I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "internalerrorException",
        "verification_token": "verification-token"
      }
      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: PUT /v1/users/self/password Cognito invalid password for forgottenPassword
    Given I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "invalidpassword",
        "verification_token": "verification-token"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidPassword",
            "description": "password does not meet requirements"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password Cognito invalid token
    Given I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "verification_token": "invalid-token"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidCode",
            "description": "verification token does not meet requirements"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password Cognito expired token
    Given I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "Password2",
        "verification_token": "expired-token"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "ExpiredCode",
            "description": "verification token has expired"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/self/password Cognito user not found
    Given I am an admin user
    When I PUT "/v1/users/self/password"
      """
      {
        "type": "ForgottenPassword",
        "email": "email@ons.gov.uk",
        "password": "Password",
        "verification_token": "verification-token"
      }
      """
    Then the HTTP status code should be "202"
