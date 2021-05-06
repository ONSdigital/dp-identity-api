Feature: users

    Scenario: POST /users and checking the response
        When I POST "/users"
        """
        {
            "email": "email@ons.gov.uk",
            "username":"Foo_Bar"
        }
        """
        Then I should receive the following JSON response with status "201":
        """
        {
            "User": {
                "Attributes": null,
                "Enabled": null,
                "MFAOptions": null,
                "UserCreateDate": null,
                "UserLastModifiedDate": null,
                "UserStatus": "UNCONFIRMED",
                "Username": "Foo_Bar"
            }
        }
        """
