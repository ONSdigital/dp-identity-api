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

    Scenario: POST /users and checking the response status 400
        When I POST "/users"
        """
        {
            "email": "",
            "username":"smileons"
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
                        "field": "validating email",
                        "param": "error validating email"
                    }
                }
            ]
        }
        """

    Scenario: POST /users and checking the response status 400
        When I POST "/users"
        """
        {
            "email": "email@ons.gov.uk",
            "username":""
        }
        """
        Then I should receive the following JSON response with status "400":
        """
        {
            "errors": [
                {
                    "error": "invalid username",
                    "message": "Unable to validate the username in the request",
                    "source": {
                        "field": "validating username",
                        "param": "error validating username"
                    }
                }
            ]
        }
        """

    Scenario: POST /users and checking the response status 400
        When I POST "/users"
        """
        {
            "email": "",
            "username":""
        }
        """
        Then I should receive the following JSON response with status "400":
        """
        {
            "errors": [
                {
                    "error": "invalid username",
                    "message": "Unable to validate the username in the request",
                    "source": {
                        "field": "validating username",
                        "param": "error validating username"
                    }
                },
                {
                    "error": "invalid email",
                    "message": "Unable to validate the email in the request",
                    "source": {
                        "field": "validating email",
                        "param": "error validating email"
                    }
                }
            ]
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
                        "field": "unmarshalling",
                        "param": "error unmarshalling request body"
                    }
                }
            ]  
        }
        """

    Scenario: POST /users with existing user and checking the response status 400
        When I POST "/users"
        """
        {
            "email": "email@ons.gov.uk",
            "username":"bob"
        }
        """
        Then I should receive the following JSON response with status "400":
        """
        {
            "errors": [
                {
                    "error": "UsernameExistsException: User account already exists",
                    "message": "UsernameExistsException: User account already exists",
                    "source": {
                        "field": "create new user pool user",
                        "param": "error creating new user pool user"
                    }
                }
            ]
        }
        """

    Scenario: POST /users unexpected server error and checking the response status 400
        When I POST "/users"
        """
        {
            "email": "email@ons.gov.uk",
            "username":"mike"
        }
        """
        Then I should receive the following JSON response with status "500":
        """
        {
            "errors": [
                {
                    "error": "InternalErrorException",
                    "message": "Failed to create new user in user pool",
                    "source": {
                        "field": "create new user pool user",
                        "param": "error creating new user pool user"
                    }
                }
            ]
        }
       """

    Scenario: POST /users duplicate email found and checking the response status 400
        When I POST "/users"
        """
        {
            "email": "email@ext.ons.gov.uk",
            "username":"mike"
        }
        """
        Then I should receive the following JSON response with status "400":
        """
        {
            "errors": [
                {
                    "error": "duplicate email",
                    "message": "duplicate email address found",
                    "source": {
                        "field": "duplicate email address check",
                        "param": "error checking duplicate email address"
                    }
                }
            ]
        }
        """
