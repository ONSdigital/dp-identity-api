Feature: Users

#   Create User
    Scenario: POST /v1/users and checking the response status 201
        When I POST "/v1/users"
            """
            {
                "forename": "smileons",
                "lastname": "bobbings",
                "email": "emailx@ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "201":
            """
            {
                "id": "123e4567-e89b-12d3-a456-426614174000",
                "forename": "smileons",
                "lastname": "bobbings",
                "email": "emailx@ons.gov.uk",
                "groups": [],
                "status": "FORCE_CHANGE_PASSWORD",
                "active": true,
                "status_notes": ""
            }
            """

    Scenario: POST /v1/users missing email and checking the response status 400
        When I POST "/v1/users"
            """
            {
                "forename": "smileons",
                "lastname": "bobbings",
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
                "lastname": "bobbings",
                "email": "emailx@ons.gov.uk"
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

    Scenario: POST /v1/users missing lastname and checking the response status 400
        When I POST "/v1/users"
            """
            {
                "forename": "smileons",
                "lastname": "",
                "email": "emailx@ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidSurname",
                        "description": "the submitted user's lastname could not be validated"
                    }
                ]
            }
            """

    Scenario: POST /v1/users and checking the response status 400
        When I POST "/v1/users"
            """
            {
                "forename": "",
                "lastname": "",
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
                        "description": "the submitted user's lastname could not be validated"
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
                "lastname": "bobbings",
                "email": "emailx@ons.gov.uk"
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
                "lastname": "bobbings",
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

#   List Users
    Scenario: GET /v1/users and checking the response status 200
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And a user with non-verified email "new_email@ons.gov.uk" and password "TeMpPassw0rd!"
        When I GET "/v1/users"
        Then I should receive the following JSON response with status "200":
            """
            {
                "users": [
                    {
                        "id": "aaaabbbbcccc",
                        "forename": "Bob",
                        "lastname": "Smith",
                        "email": "email@ons.gov.uk",
                        "groups": [],
                        "status": "CONFIRMED",
                        "active": true,
                        "status_notes": ""
                    },
                    {
                        "id": "aaaabbbbcccc",
                        "forename": "Bob",
                        "lastname": "Smith",
                        "email": "new_email@ons.gov.uk",
                        "groups": [],
                        "status": "FORCE_CHANGE_PASSWORD",
                        "active": true,
                        "status_notes": ""
                    }
                ],
                "count": 2
            }
            """

    Scenario: GET /v1/users unexpected server error and checking the response status 500
        Given a user with email "internal.error@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I GET "/v1/users"
        Then I should receive the following JSON response with status "500":
            """
            {
                "errors": [
                    {
                        "code": "InternalServerError",
                        "description": "Something went wrong"
                    }
                ]
            }
            """

#   Get User
    Scenario: GET /v1/users/{id} and checking the response status 200
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I GET "/v1/users/abcd1234"
        Then I should receive the following JSON response with status "200":
            """
            {
                "id": "abcd1234",
                "forename": "Bob",
                "lastname": "Smith",
                "email": "email@ons.gov.uk",
                "groups": [],
                "status": "CONFIRMED",
                "active": true,
                "status_notes": ""
            }
            """

    Scenario: GET /v1/users/{id} user not found and checking the response status 404
        When I GET "/v1/users/abcd1234"
        Then I should receive the following JSON response with status "404":
            """
            {
                "errors": [
                    {
                        "code": "UserNotFound",
                        "description": "the user could not be found"
                    }
                ]
            }
            """

    Scenario: GET /v1/users/{id} unexpected server error and checking the response status 500
        Given a user with username "abcd1234" and email "internal.error@ons.gov.uk" exists in the database
        When I GET "/v1/users/abcd1234"
        Then I should receive the following JSON response with status "500":
            """
            {
                "errors": [
                    {
                        "code": "InternalServerError",
                        "description": "Something went wrong"
                    }
                ]
            }
            """

#   Update User
    Scenario: PUT /v1/users/{id} to update users names and checking the response status 200
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "Changed",
                "lastname": "Names",
                "active": true,
                "status_notes": ""
            }
        """
        Then I should receive the following JSON response with status "200":
            """
            {
                "id": "abcd1234",
                "forename": "Changed",
                "lastname": "Names",
                "email": "email@ons.gov.uk",
                "groups": [],
                "status": "CONFIRMED",
                "active": true,
                "status_notes": ""
            }
            """

    Scenario: PUT /v1/users/{id} set user disabled and checking the response status 200
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "Bob",
                "lastname": "Smith",
                "active": false,
                "status_notes": "user disabled"
            }
        """
        Then I should receive the following JSON response with status "200":
            """
            {
                "id": "abcd1234",
                "forename": "Bob",
                "lastname": "Smith",
                "email": "email@ons.gov.uk",
                "groups": [],
                "status": "CONFIRMED",
                "active": false,
                "status_notes": "user disabled"
            }
            """

    Scenario: PUT /v1/users/{id} set user enabled and checking the response status 200
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        And user "abcd1234" active is "false"
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "Bob",
                "lastname": "Smith",
                "active": true,
                "status_notes": "user reactivated"
            }
        """
        Then I should receive the following JSON response with status "200":
            """
            {
                "id": "abcd1234",
                "forename": "Bob",
                "lastname": "Smith",
                "email": "email@ons.gov.uk",
                "groups": [],
                "status": "CONFIRMED",
                "active": true,
                "status_notes": "user reactivated"
            }
            """

    Scenario: PUT /v1/users/{id} set user disabled and change names and checking the response status 200
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "Changed",
                "lastname": "Names",
                "active": false,
                "status_notes": "user suspended"
            }
        """
        Then I should receive the following JSON response with status "200":
            """
            {
                "id": "abcd1234",
                "forename": "Changed",
                "lastname": "Names",
                "email": "email@ons.gov.uk",
                "groups": [],
                "status": "CONFIRMED",
                "active": false,
                "status_notes": "user suspended"
            }
            """

    Scenario: PUT /v1/users/{id} set user enabled and change names and checking the response status 200
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        And user "abcd1234" active is "false"
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "Changed",
                "lastname": "Names",
                "active": true,
                "status_notes": "user reactivated"
            }
        """
        Then I should receive the following JSON response with status "200":
            """
            {
                "id": "abcd1234",
                "forename": "Changed",
                "lastname": "Names",
                "email": "email@ons.gov.uk",
                "groups": [],
                "status": "CONFIRMED",
                "active": true,
                "status_notes": "user reactivated"
            }
            """

    Scenario: PUT /v1/users/{id} missing forename and checking the response status 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "",
                "lastname": "Smith",
                "active": true,
                "status_notes": ""
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

    Scenario: PUT /v1/users/{id} missing lastname and checking the response status 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "Bob",
                "lastname": "",
                "active": true,
                "status_notes": ""
            }
        """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidSurname",
                        "description": "the submitted user's lastname could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} invalid notes and checking the response status 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "Stan",
                "lastname": "Smith",
                "active": true,
                "status_notes": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"
            }
        """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidStatusNotes",
                        "description": "the status notes are too long"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} missing forename and lastname and checking the response status 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "",
                "lastname": "",
                "active": true,
                "status_notes": ""
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
                        "description": "the submitted user's lastname could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} missing forename and invalid notes and checking the response status 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "",
                "lastname": "Smith",
                "active": true,
                "status_notes": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"
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
                        "code": "InvalidStatusNotes",
                        "description": "the status notes are too long"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} missing lastname and invalid notes and checking the response status 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "Stan",
                "lastname": "",
                "active": true,
                "status_notes": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"
            }
        """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidSurname",
                        "description": "the submitted user's lastname could not be validated"
                    },
                    {
                        "code": "InvalidStatusNotes",
                        "description": "the status notes are too long"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} missing forename, lastname and invalid notes and checking the response status 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
        """
            {
                "forename": "",
                "lastname": "",
                "active": true,
                "status_notes": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"
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
                        "description": "the submitted user's lastname could not be validated"
                    },
                    {
                        "code": "InvalidStatusNotes",
                        "description": "the status notes are too long"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} user not found and checking the response status 404
        When I PUT "/v1/users/abcd1234"
            """
            {
                "forename": "Bob",
                "lastname": "Smith",
                "active": true,
                "status_notes": ""
            }
            """
        Then I should receive the following JSON response with status "404":
            """
            {
                "errors": [
                    {
                        "code": "UserNotFound",
                        "description": "the user could not be found"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} unexpected server error disabling user and checking the response status 500
        Given a user with username "abcd1234" and email "disable.internalerror@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
            """
            {
                "forename": "Bob",
                "lastname": "Smith",
                "active": false,
                "status_notes": "user suspended"
            }
            """
        Then I should receive the following JSON response with status "500":
            """
            {
                "errors": [
                    {
                        "code": "InternalServerError",
                        "description": "Something went wrong whilst disabling"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} unexpected server error enabling user and checking the response status 500
        Given a user with username "abcd1234" and email "enable.internalerror@ons.gov.uk" exists in the database
        And user "abcd1234" active is "false"
        When I PUT "/v1/users/abcd1234"
            """
            {
                "forename": "Bob",
                "lastname": "Smith",
                "active": true,
                "status_notes": "user reactivated"
            }
            """
        Then I should receive the following JSON response with status "500":
            """
            {
                "errors": [
                    {
                        "code": "InternalServerError",
                        "description": "Something went wrong whilst enabling"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} unexpected server error updating user and checking the response status 500
        Given a user with username "abcd1234" and email "update.internalerror@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
            """
            {
                "forename": "Bob",
                "lastname": "Smith",
                "active": true,
                "status_notes": ""
            }
            """
        Then I should receive the following JSON response with status "500":
            """
            {
                "errors": [
                    {
                        "code": "InternalServerError",
                        "description": "Something went wrong"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/{id} unexpected server error loading updated user and checking the response status 500
        Given a user with username "abcd1234" and email "internal.error@ons.gov.uk" exists in the database
        When I PUT "/v1/users/abcd1234"
            """
            {
                "forename": "Bob",
                "lastname": "Smith",
                "active": true,
                "status_notes": ""
            }
            """
        Then I should receive the following JSON response with status "500":
            """
            {
                "errors": [
                    {
                        "code": "InternalServerError",
                        "description": "Something went wrong"
                    }
                ]
            }
            """

#   Change password - auth challenge
    Scenario: PUT /v1/users/self/password and checking the response status 202
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "session": "auth-challenge-session"
            }
            """
        Then the HTTP status code should be "202"
        And the response header "Authorization" should be "Bearer accessToken"
        And the response header "ID" should be "idToken"
        And the response header "Refresh" should be "refreshToken"

    Scenario: PUT /v1/users/self/password missing type and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "UnknownRequestType",
                        "description": "unknown password change type received"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password new password required type with verification token and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "",
                "password": "Password2",
                "session": "auth-challenge-session"
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

    Scenario: PUT /v1/users/self/password missing password and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "session": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and password and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "",
                "password": "",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "",
                "password": "Password2",
                "session": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    },
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing password and session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "",
                "session": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and password and session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "",
                "password": "",
                "session": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    },
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password Cognito internal error
        Given an internal server error is returned from Cognito
        When I PUT "/v1/users/self/password"
        """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "internalerrorException",
                "session": "auth-challenge-session"
            }
        """
        Then I should receive the following JSON response with status "500":
        """
        {
            "errors": [
                {
                    "code": "InternalServerError",
                    "description": "Something went wrong"
                }
            ]
        }
        """

    Scenario: PUT /v1/users/self/password Cognito invalid password
        When I PUT "/v1/users/self/password"
        """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "invalidpassword",
                "session": "auth-challenge-session"
            }
        """
        Then I should receive the following JSON response with status "400":
        """
        {
            "errors": [
                {
                    "code": "InvalidPassword",
                    "description": "password does not meet requirements"
                }
            ]
        }
        """

    Scenario: PUT /v1/users/self/password Cognito user not found
        When I PUT "/v1/users/self/password"
        """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "Password",
                "session": "auth-challenge-session"
            }
        """
        Then the HTTP status code should be "202"

#   Change password - forgotten password
    Scenario: PUT /v1/users/self/password forgotten password type and checking the response status 202
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": "verification-token"
            }
            """
        Then the HTTP status code should be "202"

    Scenario: PUT /v1/users/self/password missing type and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "UnknownRequestType",
                        "description": "unknown password change type received"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password forgotten password type with challenge session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidToken",
                        "description": "the submitted token could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "",
                "password": "Password2",
                "verification_token": "verification-token"
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

    Scenario: PUT /v1/users/self/password missing password and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing verification token and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidToken",
                        "description": "the submitted token could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and password and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "",
                "password": "",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and verification token and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "",
                "password": "Password2",
                "verification_token": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    },
                    {
                        "code": "InvalidToken",
                        "description": "the submitted token could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and password and verification token and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "",
                "password": "",
                "verification_token": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    },
                    {
                        "code": "InvalidToken",
                        "description": "the submitted token could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password Cognito internal error
        Given an internal server error is returned from Cognito
        When I PUT "/v1/users/self/password"
        """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "internalerrorException",
                "verification_token": "verification-token"
            }
        """
        Then I should receive the following JSON response with status "500":
        """
        {
            "errors": [
                {
                    "code": "InternalServerError",
                    "description": "Something went wrong"
                }
            ]
        }
        """

    Scenario: PUT /v1/users/self/password Cognito invalid password
        When I PUT "/v1/users/self/password"
        """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "invalidpassword",
                "verification_token": "verification-token"
            }
        """
        Then I should receive the following JSON response with status "400":
        """
        {
            "errors": [
                {
                    "code": "InvalidPassword",
                    "description": "password does not meet requirements"
                }
            ]
        }
        """

    Scenario: PUT /v1/users/self/password Cognito user not found
        When I PUT "/v1/users/self/password"
        """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password",
                "verification_token": "verification-token"
            }
        """
        Then the HTTP status code should be "202"

#   Request password reset
    Scenario: POST /v1/password-reset and checking the response status 202
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I POST "/v1/password-reset"
            """
                {
                    "email": "email@ons.gov.uk"
                }
            """
        Then the HTTP status code should be "202"

    Scenario: POST /v1/password-reset missing email and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I POST "/v1/password-reset"
            """
                {
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

    Scenario: POST /v1/password-reset non ONS email address and checking the response status 202
        When I POST "/v1/password-reset"
        """
            {
                "email": "email@gmail.com"
            }
        """
        Then the HTTP status code should be "202"

    Scenario: POST /v1/password-reset Cognito internal error
        Given an internal server error is returned from Cognito
        When I POST "/v1/password-reset"
        """
            {
                "email": "internal.error@ons.gov.uk"
            }
        """
        Then I should receive the following JSON response with status "500":
        """
            {
                "errors": [
                    {
                        "code": "InternalServerError",
                        "description": "Something went wrong"
                    }
                ]
            }
        """

    Scenario: POST /v1/password-reset Cognito too many requests error
        Given an internal server error is returned from Cognito
        When I POST "/v1/password-reset"
        """
            {
                "email": "too.many@ons.gov.uk"
            }
        """
        Then I should receive the following JSON response with status "400":
        """
            {
                "errors": [
                    {
                        "code": "TooManyRequests",
                        "description": "Slow down"
                    }
                ]
            }
        """

    Scenario: POST /v1/password-reset Cognito user not found
        When I POST "/v1/password-reset"
        """
            {
                "email": "email@ons.gov.uk"
            }
        """
        Then the HTTP status code should be "202"

<<<<<<< HEAD
 #   List User Groups         
    Scenario: GET /v1/users/{id}/groups and checking the response status 200
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        And user "abcd1234" is a member of group "test-group"
        When I GET "/v1/users/abcd1234/groups"
        Then I should receive the following JSON response with status "200":
            """ 
            {
                "count":1,
                "usergroups":{
                    "Groups":[
                        {
                            "CreationDate":     null,
                            "Description":      "some Group Desciption",
                            "GroupName":        "test-group",
                            "LastModifiedDate": null,
                            "Precedence":           97,
                            "RoleArn":          null,
                            "UserPoolId":       "eu-west-18_73289nds8w932"
                        }
                    ],
                    "NextToken":null
                    }
                }
            
            """

Scenario: GET /v1/users/{id}/groups user not found and checking the response status 404
        When I GET "/v1/users/get-user-not-found"
        Then I should receive the following JSON response with status "404":
            """
            {
                "errors": [
                    {
                        "code": "UserNotFound",
                        "description": "the user could not be found"
                    }
                ]
            }
            """ 
     Scenario: GET /v1/users/{id}/groups and checking the response status 200
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I GET "/v1/users/abcd1234/groups"
        Then I should receive the following JSON response with status "200":
            """ 
            {
                "count":0,
                "usergroups":null,
                    "NextToken":null
                    }
                
            """
=======
# Get List User Groups
    Scenario: Get /v1/password-reset and checking the response status 202
        Given a user with username "abcd1234" exists in the database
        When I GET "/v1/users/test-user-1/groups"
        """
            {
                "usergroups": {
                    "Groups": [
                        {
                            "CreationDate": "2021-07-30T10:56:05.574Z",
                            "Description": "Test Group 1",
                            "GroupName": "test-group-1",
                            "LastModifiedDate": "2021-07-30T10:56:05.574Z",
                            "Precedence": 3,
                            "RoleArn": null,
                            "UserPoolId": "eu-west-1_Rnma9lp2q"
                        }
                    ],
                    "NextToken": null
                },
            "count": 1
        }
        """
        Then I should receive the following JSON response with status "200":
>>>>>>> ca1365f (Add Pagination)
