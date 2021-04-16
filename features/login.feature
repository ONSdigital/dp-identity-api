Feature: Login

    Scenario: POST /login
        When I POST "/login"
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