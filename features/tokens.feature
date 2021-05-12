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
    Given an error is returned from Cognito 
    When I POST "/tokens"
    """
    {
        "email": "email@ons.gov.uk",
        "password": "Passw0rd!"
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
                    "error": "Invalid password",
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
                    "error": "Invalid email",
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
                    "error": "Invalid email",
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
                    "error": "Invalid password",
                    "message": "Unable to validate the password in the request",
                    "source": {
                        "field": "",
                        "param": ""
                    }
                },
                {
                    "error": "Invalid email",
                    "message": "Unable to validate the email in the request",
                    "source": {
                        "field": "",
                        "param": ""
                    }
                }
            ]
        }
        """
