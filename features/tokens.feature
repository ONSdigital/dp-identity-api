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

    Scenario: DELETE /tokens/self
        Given I am not authorised
        When I DELETE "/tokens/self"
        Then I should receive the following JSON response with status "400":
        """
        {
            "errors": [
                {
                    "error": "Invalid token",
                    "message": "The provided token does not correspond to an active session",
                    "source": {
                        "field": "",
                        "param": ""
                    }
                }
            ]
        }
        """
