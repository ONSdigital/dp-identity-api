Feature: Users

    Scenario: POST /users and checking the response status 201
        When I POST "/users"
            """
            {
                "forename": "smileons",
                "surname": "bobbings",
                "email": "email@ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "201":
            """
            {
                "User": {
                    "Attributes": [
                        {
                            "Name": "sub",
                            "Value": "f0cf8dd9-755c-4caf-884d-b0c56e7d0704"
                        },
                        {
                            "Name": "name",
                            "Value": "smileons"
                        },
                        {
                            "Name": "family_name",
                            "Value": "bobbings"
                        },
                        {
                            "Name": "email",
                            "Value": "email@ons.gov.uk"
                        }
                    ],
                    "Enabled": null,
                    "MFAOptions": null,
                    "UserCreateDate": null,
                    "UserLastModifiedDate": null,
                    "UserStatus": "FORCE_CHANGE_PASSWORD",
                    "Username": "123e4567-e89b-12d3-a456-426614174000"
                }
            }
            """

    Scenario: POST /users missing email and checking the response status 400
        When I POST "/users"
            """
            {
                "forename": "smileons",
                "surname": "bobbings",
                "email": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "invalid email",
                        "description": "Unable to validate the email in the request"
                    }
                ]
            }
            """

    Scenario: POST /users missing forename and checking the response status 400
        When I POST "/users"
            """
            {
                "forename": "",
                "surname": "bobbings",
                "email": "email@ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "invalid forename",
                        "description": "Unable to validate the user's forename in the request"
                    }
                ]
            }
            """

    Scenario: POST /users missing surname and checking the response status 400
        When I POST "/users"
            """
            {
                "forename": "smileons",
                "surname": "",
                "email": "email@ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "invalid surname",
                        "description": "Unable to validate the user's surname in the request"
                    }
                ]
            }
            """

    Scenario: POST /users and checking the response status 400
        When I POST "/users"
            """
            {
                "forename": "",
                "surname": "",
                "email": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "invalid forename",
                        "description": "Unable to validate the user's forename in the request"
                    },
                    {
                        "code": "invalid surname",
                        "description": "Unable to validate the user's surname in the request"
                    },
                    {
                        "code": "invalid email",
                        "description": "Unable to validate the email in the request"
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
                        "code": "unexpected end of JSON input",
                        "description": "api endpoint POST user returned an error unmarshalling request body"
                    }
                ]
            }
            """

    Scenario: POST /users unexpected server error and checking the response status 500
        When I POST "/users"
            """
            {
                "forename": "bob",
                "surname": "bobbings",
                "email": "email@ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "500":
            """
            {
                "errors": [
                    {
                        "code": "InternalErrorException",
                        "description": "Failed to create new user in user pool"
                    }
                ]
            }
            """

    Scenario: POST /users duplicate email found and checking the response status 400
        When I POST "/users"
            """
            {
                "forename": "bob",
                "surname": "bobbings",
                "email": "email@ext.ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "duplicate email",
                        "description": "duplicate email address found"
                    }
                ]
            }
            """
