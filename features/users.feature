Feature: Users

    Scenario: POST /v1/users and checking the response status 201
        When I POST "/v1/users"
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

    Scenario: POST /v1/users missing email and checking the response status 400
        When I POST "/v1/users"
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
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    }
                ]
            }
            """

    Scenario: POST /v1/users missing forename and checking the response status 400
        When I POST "/v1/users"
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
                        "code": "InvalidForename",
                        "description": "the submitted user's forename could not be validated"
                    }
                ]
            }
            """

    Scenario: POST /v1/users missing surname and checking the response status 400
        When I POST "/v1/users"
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
                        "code": "InvalidSurname",
                        "description": "the submitted user's surname could not be validated"
                    }
                ]
            }
            """

    Scenario: POST /v1/users and checking the response status 400
        When I POST "/v1/users"
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
                        "code": "InvalidForename",
                        "description": "the submitted user's forename could not be validated"
                    },
                    {
                        "code": "InvalidSurname",
                        "description": "the submitted user's surname could not be validated"
                    },
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    }
                ]
            }
            """

    Scenario: POST /v1/users and checking the response status 500
        When I POST "/v1/users"
            """

            """
        Then I should receive the following JSON response with status "500":
            """
            {
                "errors": [
                    {
                        "code": "JSONUnmarshalError",
                        "description": "failed to unmarshal the request body"
                    }
                ]
            }
            """

    Scenario: POST /v1/users unexpected server error and checking the response status 500
        When I POST "/v1/users"
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
                        "code": "InternalServerError",
                        "description": "Failed to create new user in user pool"
                    }
                ]
            }
            """

    Scenario: POST /v1/users duplicate email found and checking the response status 400
        When I POST "/v1/users"
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
                        "code": "InvalidEmail",
                        "description": "account using email address found"
                    }
                ]
            }
            """
