Feature: Tokens

Scenario: POST /tokens
    Given a user exists in the database
    When I post "/tokens"
        """
        {
            "email": "email@ons.gov.uk",
            "password": "Passw0rd!"
        }
        """
    Then the HTTP status code should be "201"
    And the response header "Authorization" should be "VALUE"
    And the response header "ID" should be "VALUE"
    And the response header "Refresh" should be "VALUE"

Scenario: POST /tokens
    Given a user does not exists in the database
    When I post "/tokens"
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
                 "error": "Unable to autheticate the request",
                    "message": "Unautheticated user",
                    "source": {
                        "field": "",
                        "param": ""
                    }
                }
            ]
        }
        """
    
Scenario: POST /tokens
    Given a user does not exists in the database
    When I post "/tokens"
        """
        {
            "email": "email1@ons.gov.uk",
            "password": "Passw0rd!"
        }
        """
    Then I should receive the following JSON response with status "403":
        """
        {
            "errors": [
                {
                    "error": "Forbidden",
                    "message": "Too many login attempts",
                    "source": {
                        "field": "",
                        "param": ""
                    }
                }
            ]
        }
        """

Scenario: POST /tokens
    Given Cognito has an internal server error 
        When I post "/tokens"
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
                        "error": "Internal Server Error",
                        "message": "Amazon Cognito has encountered an internal error",
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
