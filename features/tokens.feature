Feature: Tokens

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
                    "error": "Invalid token",
                    "message": "No Authorization token was provided",
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
                    "error": "Invalid token",
                    "message": "The provided token does not meet the required format",
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
                    "error": "Invalid token",
                    "message": "The provided token does not meet the required format",
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
