Feature: Users

    Scenario: POST /users and checking the response status 201
        When I POST "/users"
        """
        {
            "email": "email@ons.gov.uk",
            "username":"smileons"
        }
        """
        Then I should receive the following JSON response with status "201":
        """
        {
            "User":{
                "Attributes":null, 
                "Enabled":null, 
                "MFAOptions":null, 
                "UserCreateDate":null, 
                "UserLastModifiedDate":null,
                "UserStatus":"FORCE_CHANGE_PASSWORD", 
                "Username":"smileons"
            }
        }
        """

 Scenario: POST /users and checking the response status 500
        When I POST "/users"
        """
       
        """
        Then I should receive the following JSON response with status "500":
        """
        {
            "errors": [
                {
                    "error": "unexpected end of JSON input",
                    "message": "api endpoint POST user returned an error unmarshalling request body",
                    "source": {
                        "field": "",
                        "param": ""
                    }
                }
            ]  
        }
        """