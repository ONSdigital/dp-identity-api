Feature: Login

    Scenario: POST /login
        When I POST "/login" 
        """
        {
            "email":"email",
            "password":"password"
        }
        """
        Then the HTTP status code should be "400"
