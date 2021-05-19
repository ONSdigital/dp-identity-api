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
                    "error": "NotAuthorizedException: Incorrect username or password.",
                    "message": "unautheticated user: Unable to autheticate request",
                    "source": {
                        "field": "",
                        "param": ""
                    }
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
                    "error": "NotAuthorizedException: Password attempts exceeded",
                    "message": "exceeded the number of attemps to login in with the provided credentials",
                    "source": {
                        "field": "",
                        "param": ""
                    }
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
                "error": "InternalErrorException",
                "message": "api endpoint POST login returned an error and failed to login to cognito",
                "source": {
                    "field": "",
                    "param": ""
                }
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
                "error": "InvalidParameterException",
                "message": "something went wrong, and api endpoint POST login returned an error and failed to login to cognito. Please try again or contact an administrator.",
                "source": {
                    "field": "",
                    "param": ""
                }
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
                    "error": "invalid password",
                    "message": "Unable to validate the password in the request",
                    "source": {
                        "field": "",
                        "param": ""
                    }
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
                    "error": "invalid email",
                    "message": "Unable to validate the email in the request",
                    "source": {
                        "field": "",
                        "param": ""
                    }
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
                    "error": "invalid email",
                    "message": "Unable to validate the email in the request",
                    "source": {
                        "field": "",
                        "param": ""
                    }
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
                    "error": "invalid password",
                    "message": "Unable to validate the password in the request",
                    "source": {
                        "field": "",
                        "param": ""
                    }
                },
                {
                    "error": "invalid email",
                    "message": "Unable to validate the email in the request",
                    "source": {
                        "field": "",
                        "param": ""
                    }
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
                "error": "invalid token",
                "message": "no Authorization token was provided",
                "source": {
                    "field": "",
                    "param": ""
                }
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
                "error": "invalid token",
                "message": "the provided token does not meet the required format",
                "source": {
                    "field": "",
                    "param": ""
                }
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
                "error": "invalid token",
                "message": "the provided token does not meet the required format",
                "source": {
                    "field": "",
                    "param": ""
                }
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
                "error": "invalid ID token",
                "message": "no ID token was provided",
                "source": {
                    "field": "",
                    "param": ""
                }
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
                "error": "invalid refresh token",
                "message": "no refresh token was provided",
                "source": {
                    "field": "",
                    "param": ""
                }
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
                "error": "invalid refresh token",
                "message": "no refresh token was provided",
                "source": {
                    "field": "",
                    "param": ""
                }
            },
            {
                "error": "invalid ID token",
                "message": "no ID token was provided",
                "source": {
                    "field": "",
                    "param": ""
                }
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
                "error": "invalid ID token",
                "message": "the ID token could not be parsed",
                "source": {
                    "field": "",
                    "param": ""
                }
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

Scenario: PUT /tokens/self with tokens from different users
    Given I have a valid ID header for user "test@ons.gov.uk"
    And I set the "Refresh" header to "AnotherUser"
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
