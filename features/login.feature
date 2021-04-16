Feature: Login

    Scenario: POST /login
        When I POST "/login"
        """
        {
            "email": "",
            "password": ""
        }
        """
        Then the HTTP status code should be "400"
        And I should receive the following response:
        """
        {
            "errors": [
                {
                    "error": "string, unchanging so devs can use this in code",
                    "message": "detailed explanation of error",
                    "source": {
                        "field": "reference to field like some.field or something",
                        "param": "query param causing issue"
                    }
                }
            ]
        }
        """

    Scenario: POST /login
        When I POST "/login"
        """
        {
            "email": "email",
            "password": "password"
        }
        """
        Then the HTTP status code should be "400"