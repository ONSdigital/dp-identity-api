Feature: users

    Scenario: POST /users and checking the response
        When I POST "/users"
        """
        {
            "email": "email@ons.gov.uk",
            "password": "password"
            "username":"smileons"
        }
        """
        Then I should receive the following JSON response with status "200":
        """
        {
           
            {
                "UserAttributes":{
                    Name  "email"
                    Value "email@ons.gov.uk"
                }
                "Username": "smileons"
            }

        }
        """
