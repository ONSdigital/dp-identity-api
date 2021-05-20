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
                    "code": "NotAuthorizedException: Incorrect username or password.",
                    "description": "unautheticated user: Unable to autheticate request"
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
                    "code": "NotAuthorizedException: Password attempts exceeded",
                    "description": "exceeded the number of attemps to login in with the provided credentials"
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
                "code": "InternalErrorException",
                "description": "api endpoint POST login returned an error and failed to login to cognito"
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
                "code": "InvalidParameterException",
                "description": "something went wrong, and api endpoint POST login returned an error and failed to login to cognito. Please try again or contact an administrator."
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
                    "code": "invalid password",
                    "description": "Unable to validate the password in the request"
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
                    "code": "invalid email",
                    "description": "Unable to validate the email in the request"
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
                    "code": "invalid email",
                    "description": "Unable to validate the email in the request"
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
                    "code": "invalid password",
                    "description": "Unable to validate the password in the request"
                },
                {
                    "code": "invalid email",
                    "description": "Unable to validate the email in the request"
                }
            ]
        }
        """

Scenario: DELETE /tokens/self
    Given I am not authorised
    When I DELETE "/tokens/self"
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "invalid token",
                "description": "no Authorization token was provided"
            }
        ]
    }
    """

Scenario: DELETE /tokens/self
    Given I set the "Authorization" header to "Bearer"
    When I DELETE "/tokens/self"
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "invalid token",
                "description": "the provided token does not meet the required format"
            }
        ]
    }
    """

Scenario: DELETE /tokens/self
    Given I set the "Authorization" header to "BearerSomeToken"
    When I DELETE "/tokens/self"
    Then I should receive the following JSON response with status "400":
    """
    {
        "errors": [
            {
                "code": "invalid token",
                "description": "the provided token does not meet the required format"
            }
        ]
    }
    """

Scenario: DELETE /tokens/self
    Given I set the "Authorization" header to "Bearer InternalError"
    When I DELETE "/tokens/self"
    Then the HTTP status code should be "500"

Scenario: DELETE /tokens/self
    Given I set the "Authorization" header to "Bearer xxxx.yyyy.zzzz"
    When I DELETE "/tokens/self"
    Then the HTTP status code should be "400"

Scenario: DELETE /tokens/self
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
                "code": "invalid ID token",
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
                "code": "invalid refresh token",
                "description": "no refresh token was provided"
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
                "code": "invalid refresh token",
                "description": "no refresh token was provided"
            },
            {
                "code": "invalid ID token",
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
                "code": "invalid ID token",
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
