Feature: Tokens

Scenario: POST /tokens
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/tokens"
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

Scenario: POST /tokens
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/tokens"
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
                    "description": "Incorrect username or password"
                }
            ]
        }
        """

Scenario: POST /tokens
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/tokens"
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

Scenario: POST /tokens
    Given an internal server error is returned from Cognito
    When I POST "/tokens"
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

Scenario: POST /tokens
    Given an error is returned from Cognito
    When I POST "/tokens"
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

Scenario: POST /tokens
    When I POST "/tokens"
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

Scenario: POST /tokens
    When I POST "/tokens"
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

Scenario: POST /tokens
    When I POST "/tokens"
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

Scenario: POST /tokens
    When I POST "/tokens"
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

Scenario: POST /tokens
    Given I have an active session with access token "aaaa.bbbb.cccc"
    And a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/tokens"
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

Scenario: POST /tokens
    Given I am not authorised
    And a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/tokens"
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

Scenario: POST /tokens
    Given the AdminUserGlobalSignOut endpoint in cognito returns an internal server error
    And a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/tokens"
    """
    {
        "email": "internalservererror@ons.gov.uk",
        "password": "Passw0rd!"
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

Scenario: POST /tokens
    Given the AdminUserGlobalSignOut endpoint in cognito returns an internal server error
    And a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    When I POST "/tokens"
    """
    {
        "email": "clienterror@ons.gov.uk",
        "password": "Passw0rd!"
    }
    """
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "NotAuthorised",
                "description": "Something went wrong"
            }
        ]
    }
    """

Scenario: DELETE /tokens/self no Authorization header
    Given I am not authorised
    When I DELETE "/tokens/self"
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

Scenario: DELETE /tokens/self Authorization header missing JWT
    Given I set the "Authorization" header to "Bearer"
    When I DELETE "/tokens/self"
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

Scenario: DELETE /tokens/self malformed Authorization header
    Given I set the "Authorization" header to "BearerSomeToken"
    When I DELETE "/tokens/self"
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

Scenario: DELETE /tokens/self Cognito internal error
    Given I set the "Authorization" header to "Bearer InternalError"
    When I DELETE "/tokens/self"
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

Scenario: DELETE /tokens/self access token not valid in Cognito
    Given I set the "Authorization" header to "Bearer xxxx.yyyy.zzzz"
    When I DELETE "/tokens/self"
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

Scenario: DELETE /tokens/self success
    Given I have an active session with access token "aaaa.bbbb.cccc"
    And I set the "Authorization" header to "Bearer aaaa.bbbb.cccc"
    When I DELETE "/tokens/self"
    Then the HTTP status code should be "204"

Scenario: PUT /tokens/self with no ID token
    Given I set the "ID" header to ""
    And I set the "Refresh" header to "aaaa.bbbb.cccc.dddd.eeee"
    When I PUT "/tokens/self"
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

Scenario: PUT /tokens/self with no refresh token
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to ""
    When I PUT "/tokens/self"
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

Scenario: PUT /tokens/self with no tokens
    Given I set the "ID" header to ""
    And I set the "Refresh" header to ""
    When I PUT "/tokens/self"
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

Scenario: PUT /tokens/self with badly formatted ID token
    Given I set the "ID" header to "zzzz.yyyy.xxxx"
    And I set the "Refresh" header to "aaaa.bbbb.cccc.dddd.eeee"
    When I PUT "/tokens/self"
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

Scenario: PUT /tokens/self internal Cognito error
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to "InternalError"
    When I PUT "/tokens/self"
    """
    {}
    """
    Then the HTTP status code should be "500"

Scenario: PUT /tokens/self with expired refresh token
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to "ExpiredToken"
    When I PUT "/tokens/self"
    """
    {}
    """
    Then the HTTP status code should be "403"

Scenario: PUT /tokens/self success
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to "aaaa.bbbb.cccc.dddd.eeee"
    When I PUT "/tokens/self"
    """
    {}
    """
    Then the HTTP status code should be "201"
    And the response header "Authorization" should be "Bearer llll.mmmm.nnnn"
