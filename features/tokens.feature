Feature: Tokens

Scenario: POST /v1/tokens successful login
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/v1/tokens"
    """
    {
        "email": "email@ons.gov.uk",
        "password": "Passw0rd!"
    }
    """
    Then the HTTP status code should be "201"
    And the response header "Authorization" should be "Bearer accessToken"
    And the response header "ID" should be "idToken"
    And the response header "Refresh" should be "refreshToken"

Scenario: POST /v1/tokens non-verified email successful login
    Given a user with non-verified email "new_email@ons.gov.uk" and password "TeMpPassw0rd!"
    When I POST "/v1/tokens"
    """
    {
        "email": "new_email@ons.gov.uk",
        "password": "TeMpPassw0rd!"
    }
    """
    Then I should receive the following JSON response with status "202":
    """
    {
        "new_password_required": "true",
        "session": "AYABeBBsY5be-this-is-a-test-session-id-string-123456789iuerhcfdisieo-end"
    }
    """

Scenario: POST /v1/tokens 401 - invalid credentials
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/v1/tokens"
    """
    {
        "email": "email1@ons.gov.uk",
        "password": "Passw0rd!"
    }
    """
    Then I should receive the following JSON response with status "401":
        """
        {
            "errors": [
                {
                    "code": "NotAuthorised",
                    "description": "Incorrect username or password."
                }
            ]
        }
        """

Scenario: POST /v1/tokens 403 - too many failed attempts
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/v1/tokens"
    """
    {
        "email": "email@ons.gov.uk",
        "password": "TooManyPasswordAttempts"
    }
    """
    Then I should receive the following JSON response with status "403":
        """
        {
            "errors": [
                {
                    "code": "TooManyFailedAttempts",
                    "description": "Password attempts exceeded"
                }
            ]
        }
        """

Scenario: POST /v1/tokens Cognito internal error
    Given an internal server error is returned from Cognito
    When I POST "/v1/tokens"
    """
    {
        "email": "email@ons.gov.uk",
        "password": "internalerrorException"
    }
    """
    Then I should receive the following JSON response with status "500":
    """
    {
        "errors": [
            {
                "code": "InternalServerError",
                "description": "Something went wrong"
            }
        ]
    }
    """

Scenario: POST /v1/tokens
    Given an error is returned from Cognito
    When I POST "/v1/tokens"
    """
    {
        "email": "email@ons.gov.uk",
        "password": "Passw0rd!"
    }
    """
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "InvalidField",
                "description": "A parameter was invalid"
            }
        ]
    }
    """

Scenario: POST /v1/tokens 400 - no password submitted
    When I POST "/v1/tokens"
        """
        {
            "email": "email@ons.gov.uk",
            "password": ""
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

Scenario: POST /v1/tokens 400 - email does not match regex
    When I POST "/v1/tokens"
        """
        {
            "email": "email",
            "password": "password"
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

Scenario: POST /v1/tokens 400 - no email submitted
    When I POST "/v1/tokens"
        """
        {
            "email": "",
            "password": "password"
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

Scenario: POST /v1/tokens 400 - no email or password submitted
    When I POST "/v1/tokens"
        """
        {
            "email": "",
            "password": ""
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

Scenario: DELETE /v1/tokens/self no Authorization header
    Given I am not authorised
    When I DELETE "/v1/tokens/self"
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "InvalidToken",
                "description": "no Authorization token was provided"
            }
        ]
    }
    """

Scenario: DELETE /v1/tokens/self Authorization header missing JWT
    Given I set the "Authorization" header to "Bearer"
    When I DELETE "/v1/tokens/self"
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "InvalidToken",
                "description": "the authorization token does not meet the required format"
            }
        ]
    }
    """

Scenario: DELETE /v1/tokens/self malformed Authorization header
    Given I set the "Authorization" header to "BearerSomeToken"
    When I DELETE "/v1/tokens/self"
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "InvalidToken",
                "description": "the authorization token does not meet the required format"
            }
        ]
    }
    """

Scenario: DELETE /v1/tokens/self Cognito internal error
    Given I set the "Authorization" header to "Bearer InternalError"
    When I DELETE "/v1/tokens/self"
    Then I should receive the following JSON response with status "500":
    """
    {
        "errors": [
            {
                "code": "InternalServerError",
                "description": "Something went wrong"
            }
        ]
    }
    """

Scenario: DELETE /v1/tokens/self access token not valid in Cognito
    Given I set the "Authorization" header to "Bearer xxxx.yyyy.zzzz"
    When I DELETE "/v1/tokens/self"
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "NotAuthorised",
                "description": "Access Token has been revoked"
            }
        ]
    }
    """

Scenario: DELETE /v1/tokens/self success
    Given I have an active session with access token "aaaa.bbbb.cccc"
    And I set the "Authorization" header to "Bearer aaaa.bbbb.cccc"
    When I DELETE "/v1/tokens/self"
    Then the HTTP status code should be "204"

Scenario: PUT /v1/tokens/self with no ID token
    Given I set the "ID" header to ""
    And I set the "Refresh" header to "aaaa.bbbb.cccc.dddd.eeee"
    When I PUT "/v1/tokens/self"
    """
    {}
    """
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "InvalidToken",
                "description": "no ID token was provided"
            }
        ]
    }
    """

Scenario: PUT /v1/tokens/self with no refresh token
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to ""
    When I PUT "/v1/tokens/self"
    """
    {}
    """
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "InvalidToken",
                "description": "no Refresh token was provided"
            }
        ]
    }
    """

Scenario: PUT /v1/tokens/self with no tokens
    Given I set the "ID" header to ""
    And I set the "Refresh" header to ""
    When I PUT "/v1/tokens/self"
    """
    {}
    """
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "InvalidToken",
                "description": "no Refresh token was provided"
            },
            {
                "code": "InvalidToken",
                "description": "no ID token was provided"
            }
        ]
    }
    """

Scenario: PUT /v1/tokens/self with badly formatted ID token
    Given I set the "ID" header to "zzzz.yyyy.xxxx"
    And I set the "Refresh" header to "aaaa.bbbb.cccc.dddd.eeee"
    When I PUT "/v1/tokens/self"
    """
    {}
    """
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "InvalidToken",
                "description": "the ID token could not be parsed"
            }
        ]
    }
    """

Scenario: PUT /v1/tokens/self internal Cognito error
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to "InternalError"
    When I PUT "/v1/tokens/self"
    """
    {}
    """
    Then the HTTP status code should be "500"

Scenario: PUT /v1/tokens/self with expired refresh token
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to "ExpiredToken"
    When I PUT "/v1/tokens/self"
    """
    {}
    """
    Then the HTTP status code should be "403"

Scenario: PUT /v1/tokens/self success
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to "aaaa.bbbb.cccc.dddd.eeee"
    When I PUT "/v1/tokens/self"
    """
    {}
    """
    Then the HTTP status code should be "201"
    And the response header "Authorization" should be "Bearer llll.mmmm.nnnn"
